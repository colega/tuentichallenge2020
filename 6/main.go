package main

import (
	"fmt"
	"net"

	"github.com/fatih/color"
)

type tile int

const (
	UNKNOWN tile = iota
	FREE
	WALL
	VISITED
	KNIGHT
	PRINCESS
)

var tileRepr = map[tile]string{
	UNKNOWN:  " ",
	FREE:     color.BlueString("."),
	WALL:     color.WhiteString("#"),
	VISITED:  color.GreenString("@"),
	KNIGHT:   color.New(color.Bold, color.BgHiRed, color.FgHiBlack).Sprintf("K"),
	PRINCESS: color.New(color.Bold, color.BgHiWhite, color.FgHiRed).Sprintf("P"),
}
var mapChar = map[byte]tile{
	'.': FREE,
	'#': WALL,
	'K': KNIGHT,
	'P': PRINCESS,
}
var buf = make([]byte, 1024)

func main() {
	conn, err := net.Dial("tcp", "52.49.91.111:2003")
	defer conn.Close()

	assertNoError(err)

	m := map[xy]tile{}
	k := xy{0, 0}
	update(conn, m, k)

	search(conn, m, k)
}

func search(conn net.Conn, m map[xy]tile, k xy) bool {
	for _, mov := range movements(m, k) {
		k = move(conn, m, k, mov)
		if m[k] == PRINCESS {
			n, err := conn.Read(buf)
			assertNoError(err)
			fmt.Printf("We should be at the princess now, received:\n%s", buf[:n])
			return true
		}

		update(conn, m, k)
		if found := search(conn, m, k); found {
			return found
		}

		k = move(conn, m, k, mov.reverse())
		update(conn, m, k)
	}
	return false
}

func move(conn net.Conn, m map[xy]tile, k xy, mov xy) xy {
	fmt.Printf("Moving %s\n", string(mov.movementCommand()))
	_, err := conn.Write(mov.movementCommand())
	assertNoError(err)

	m[k] = VISITED
	k = k.move(mov)
	return k
}

func update(conn net.Conn, m map[xy]tile, k xy) {
	n, err := conn.Read(buf)
	assertNoError(err)
	readMap(m, k, buf[:n])

	fmt.Printf("Received:\n%s\n", buf[:n])
	fmt.Println("Map now is:")
	drawMap(m)
}

func movements(m map[xy]tile, k xy) []xy {
	var free []xy
	for _, mov := range allMovements {
		pos := k.move(mov)
		switch m[pos] {
		case PRINCESS:
			return []xy{mov}
		case FREE:
			free = append(free, mov)
		}
	}
	return free
}

func readMap(m map[xy]tile, k xy, recv []byte) {
	rk, ok := findK(recv)
	if !ok {
		panic("can't find king")
	}
	offX := k.x - rk.x
	offY := k.y - rk.y

	var x, y int
	for i, c := range recv {
		if c == '\n' {
			y++
			x = 0
			continue
		}
		if c == '-' {
			break
		}
		t, ok := mapChar[c]
		if !ok {
			panic(fmt.Errorf("unexpected char %c at %d in:\n%s", c, i, string(recv)))
		}
		p := xy{x + offX, y + offY}
		if t == KNIGHT || m[p] != VISITED {
			m[p] = t
		}
		x++
	}
}

func findK(recv []byte) (xy, bool) {
	var x, y int
	for _, c := range recv {
		if c == '\n' {
			y++
			x = 0
			continue
		}
		t, ok := mapChar[c]
		if !ok {
			return xy{}, false
		}
		if t == KNIGHT {
			return xy{x, y}, true
		}
		x++
	}
	return xy{}, false
}

func drawMap(m map[xy]tile) {
	var x0, y0, xM, yM int
	for p := range m {
		x0 = min(x0, p.x)
		y0 = min(y0, p.y)
		xM = max(xM, p.x)
		yM = max(yM, p.y)
	}

	for y := y0; y <= yM; y++ {
		for x := x0; x <= xM; x++ {
			fmt.Print(tileRepr[m[xy{x, y}]])
		}
		fmt.Println()
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func assertNoError(err error) {
	if err != nil {
		panic(err)
	}
}

type xy struct{ x, y int }

func (p xy) move(to xy) xy {
	p.x += to.x
	p.y += to.y
	return p
}

func (p xy) reverse() xy {
	return xy{-p.x, -p.y}
}

func (p xy) movementCommand() []byte {
	mv := ""
	if p.y > 0 {
		mv += fmt.Sprintf("%dD", p.y)
	} else {
		mv += fmt.Sprintf("%dU", -p.y)
	}
	if p.x > 0 {
		mv += fmt.Sprintf("%dR", p.x)
	} else {
		mv += fmt.Sprintf("%dL", -p.x)
	}
	return []byte(mv)
}

var allMovements = []xy{
	{-1, -2},
	{1, -2},
	{-2, -1},
	{2, -1},
	{-2, 1},
	{2, 1},
	{-1, 2},
	{1, 2},
}
