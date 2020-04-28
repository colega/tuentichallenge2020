package main

import (
	"encoding/hex"
	"fmt"
)

func main() {
	ciphered1 := "3633363A33353B393038383C363236333635313A353336"
	msg1 := "514;248;980;347;145;332"
	key := findKey(msg1, ciphered1)

	ciphered2 := "3A3A333A333137393D39313C3C3634333431353A37363D"
	fmt.Println(decrypt(key, ciphered2))
}

func encrypt(key, msg string) string {
	encrypted := ""
	for i, char := range msg {
		keyPos := len(key) - 1 - i
		keyChar := key[keyPos] - '0'
		crptChar := uint8(char) ^ keyChar
		encrypted += fmt.Sprintf("%02x", crptChar)
		fmt.Printf("char=%c, ascii=%d, keyPos=%d, keyChar=%c=%d, crptChar=%d\n", char, char, keyPos, keyChar, keyChar, crptChar)
	}
	return encrypted
}

func decrypt(key, encrypted string) string {
	if len(encrypted)%2 != 0 {
		panic("not even length encrypted")
	}
	msg := ""
	encBytes := make([]byte, len(encrypted)/2)
	_, err := hex.Decode(encBytes, []byte(encrypted))
	if err != nil {
		panic(err)
	}
	for i, b := range encBytes {
		msgChar := b ^ (key[len(key)-1-i] - '0')
		msg += string(msgChar)
	}
	return msg
}

func findKey(msg, encrypted string) string {
	if len(encrypted)%2 != 0 {
		panic("not even length encrypted")
	}
	key := ""
	encBytes := make([]byte, len(encrypted)/2)
	_, err := hex.Decode(encBytes, []byte(encrypted))
	if err != nil {
		panic(err)
	}
	for i, b := range encBytes {
		keyChar := b ^ msg[i]
		key = string(keyChar+'0') + key
	}
	return key
}
