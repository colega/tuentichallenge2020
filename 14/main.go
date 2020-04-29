package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

var learnRegex = regexp.MustCompile(`LEARN \{servers: \[([\d,]+)\], secret_owner: (\d+)\}`)

var (
	registryLock sync.RWMutex
	registry     = map[int]*client{}
)

func main() {
	c := color.FgGreen
	for i := 0; i < 5; i++ {
		conn, err := net.Dial("tcp", "52.49.91.111:2092")
		assertNoError(err)
		c := &client{color: color.New(c + color.Attribute(i%5)), leader: i == 0}
		c.start(conn)
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Kill, os.Interrupt)
	<-ch
}

type promise struct {
	round, from, n, v int
}

type client struct {
	conn     net.Conn
	serverID int
	scanner  *bufio.Scanner

	promises chan promise

	sync.Mutex
	serversList []int
	secretOwner int

	serversUpdated chan []int

	color         *color.Color
	roundFinished chan struct{}

	leader bool
}

func (c *client) start(conn net.Conn) {
	c.conn = conn
	c.scanner = bufio.NewScanner(conn)
	c.log("CONNECTING")
	for c.scanner.Scan() {
		line := c.scanner.Text()
		if strings.Contains(line, "SERVER ID: ") {
			var err error
			c.serverID, err = strconv.Atoi(strings.Replace(line, "SERVER ID: ", "", -1))
			assertNoError(err)
			break
		}
	}
	c.log("CONNECTED")
	c.promises = make(chan promise, 10)
	c.serversUpdated = make(chan []int)
	c.roundFinished = make(chan struct{})

	go c.read()
	go c.watchPromises()
	if c.leader {
		go c.propose()
	}

	registryLock.Lock()
	registry[c.serverID] = c
	registryLock.Unlock()
}

func (c *client) read() {
	for c.scanner.Scan() {
		line := c.scanner.Text()
		if strings.Contains(line, "BAD COMMAND") {
			c.logError("<<< RECEIVED: %s", line)
		} else {
			c.log("<<< RECEIVED: %s", line)
		}
		if strings.Contains(line, "-> LEARN") {
			c.handleLearn(line)
		} else if strings.Contains(line, "-> PROMISE ") {
			c.promises <- parsePromise(line)
		} else if strings.Contains(line, ": PROMISE") && strings.Contains(line, "-> 9") {
			c.log("GOT PROMISE FOR 9")
			c.promises <- parseOurPromise(line)
		}
	}
}

func (c *client) quorum() int {
	c.Lock()
	defer c.Unlock()
	return len(c.serversList)/2 + 1
}

func (c *client) watchPromises() {
	pendingToAccept := make(chan promise, 100)
	for {
		select {
		case p := <-c.promises:
			pendingToAccept <- p

			q := c.quorum()
			if len(pendingToAccept) > q {
				c.log("Quorum on promises, need %d, got %d", q, len(pendingToAccept))
				go c.accept(pendingToAccept)
			} else if !c.leader {
				c.log("Not leader, just accepting")
				go c.accept(pendingToAccept)
			}
		case <-c.roundFinished:
			c.log("Resetting promises")
			if pendingToAccept != nil {
				close(pendingToAccept)
				pendingToAccept = make(chan promise, 100)
			}
		}
	}
}

func (c *client) acceptedServerList() []int {
	hash := c.serverHash()
	added := false
	_ = added
	registryLock.RLock()
	for srv := range registry {
		if _, ok := hash[srv]; !ok {
			hash[srv] = true
			added = true
			//BAD COMMAND IGNORED: Cluster membership must be modified one by one
			break
		}
	}

	c.Lock()
	if !added {
		for _, srv := range c.serversList {
			if _, ok := registry[srv]; !ok {
				delete(hash, srv)
				c.log("Removing %d from accepted server list", srv)
				break
			}
		}
	}
	c.Unlock()

	registryLock.RUnlock()
	var list []int
	for srv := range hash {
		list = append(list, srv)
	}
	sort.Ints(list)
	return list
}

func (c *client) accept(promises chan promise) {
	c.Lock()
	secretOwner := c.secretOwner
	c.Unlock()
	srvListStrings := isToAs(c.acceptedServerList())
	if len(c.ourActiveServers()) >= c.quorum() {
		c.logError("We have the quorum, let's try to hijack")
		secretOwner = 9
	}
	for p := range promises {
		msg := fmt.Sprintf("ACCEPT {id: {%d,%d}, value: {servers: [%s], secret_owner: %d}} -> %d",
			p.n, p.v, strings.Join(srvListStrings, ","), secretOwner, p.from)
		_, err := c.conn.Write([]byte(msg + "\n"))
		c.log(">>> SENDING: %s (for round %d)", msg, p.round)
		assertNoError(err)
	}
	c.log("Finished accepting")
}

func (c *client) propose() {
	proposal := 1
	for srvList := range c.serversUpdated {
		proposal++
		c.log("Will propose to %v", srvList)
		for _, dest := range srvList {
			msg := fmt.Sprintf("PREPARE {%d,%d} -> %d", proposal, c.serverID, dest)
			c.log(">>> SENDING: %s", msg)
			_, err := c.conn.Write([]byte(msg + "\n"))
			assertNoError(err)
		}
		time.Sleep(time.Second)
	}
}

func (c *client) handleLearn(line string) {
	srvList, secret, ok := parseLearn(line)
	if !ok {
		panic(fmt.Errorf("cant parse learn: %s", line))
	}
	c.Lock()
	c.serversList = srvList
	c.secretOwner = secret
	c.Unlock()
	c.log("Updated servers to %v secret owner to %d", srvList, secret)
	c.log("Our servers are: %v, active=%d, len=%d, quorum=%d", ourServerList(), c.ourActiveServers(), len(c.ourActiveServers()), c.quorum())
	if len(c.ourActiveServers()) >= c.quorum() {
		c.logError("WE HAVE THE QUORUM")
	}
	select {
	case c.serversUpdated <- srvList:
	default:
		c.log("Can't notify updated servers")
	}
	c.roundFinished <- struct{}{}
	if secret == 9 {
		for j := 0; j <= 10; j++ {
			c.logError("WE HAVE LEARNED THE SECRET ==============================================================")
		}
	}
}

func ourServerList() []int {
	registryLock.RLock()
	defer registryLock.RUnlock()
	var list []int
	for srv := range registry {
		list = append(list, srv)
	}
	return list
}

func (c *client) ourActiveServers() []int {
	our := ourServerList()
	hash := c.serverHash()

	var list []int
	for _, srv := range our {
		if _, ok := hash[srv]; ok {
			list = append(list, srv)
		}
	}
	return list
}

func (c *client) serverHash() map[int]bool {
	c.Lock()
	defer c.Unlock()
	hash := map[int]bool{}
	for _, srv := range c.serversList {
		hash[srv] = true
	}
	return hash
}

func (c *client) log(format string, args ...interface{}) {
	fmt.Println(c.color.Sprintf("#%2d", c.serverID), c.color.Sprintf(format, args...))
}
func (c *client) logError(format string, args ...interface{}) {
	fmt.Println(color.HiRedString("#%2d", c.serverID), color.HiRedString(format, args...))
}

func assertNoError(err error) {
	if err != nil {
		panic(err)
	}
}

func isToAs(ints []int) []string {
	ss := make([]string, len(ints))
	for i, n := range ints {
		ss[i] = strconv.Itoa(n)
	}
	return ss
}

//ROUND 960: PROMISE {29,9} no_proposal -> 9
func parseOurPromise(line string) promise {
	regex := regexp.MustCompile(`ROUND (\d+): PROMISE \{(\d+),(\d+)\} [\w_]+ -> (\d+)`)
	const (
		idxRound = iota + 1
		idxN
		idxV
		idxFrom
	)
	regexMatches := regex.FindStringSubmatch(line)
	round, err := strconv.Atoi(regexMatches[idxRound])
	assertNoError(err)

	from, err := strconv.Atoi(regexMatches[idxFrom])
	assertNoError(err)

	n, err := strconv.Atoi(regexMatches[idxN])
	assertNoError(err)

	v, err := strconv.Atoi(regexMatches[idxV])
	assertNoError(err)
	return promise{round, from, n, v}
}

func parsePromise(line string) promise {
	regex := regexp.MustCompile(`ROUND (\d+): (\d+) -> PROMISE \{(\d+),(\d+)\}`)
	const (
		idxRound = iota + 1
		idxFrom
		idxN
		idxV
	)
	regexMatches := regex.FindStringSubmatch(line)
	round, err := strconv.Atoi(regexMatches[idxRound])
	assertNoError(err)

	from, err := strconv.Atoi(regexMatches[idxFrom])
	assertNoError(err)

	n, err := strconv.Atoi(regexMatches[idxN])
	assertNoError(err)

	v, err := strconv.Atoi(regexMatches[idxV])
	assertNoError(err)
	return promise{round, from, n, v}
}

func parseLearn(line string) ([]int, int, bool) {
	all := learnRegex.FindStringSubmatch(line)
	if len(all) < 3 {
		return nil, 0, false
	}
	ss := strings.Split(all[1], ",")
	ints := make([]int, len(ss))
	var err error
	for i, s := range ss {
		ints[i], err = strconv.Atoi(s)
		assertNoError(err)
	}
	n, err := strconv.Atoi(all[2])
	assertNoError(err)
	return ints, n, true
}
