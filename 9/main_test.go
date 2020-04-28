package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	key := "123456"
	msg := "asdfgh"
	encrypted := "677660656569"

	assert.Equal(t, encrypted, encrypt(key, msg))
}

func TestDecrypt(t *testing.T) {
	key := "123456"
	msg := "asdfgh"
	encrypted := "677660656569"

	assert.Equal(t, msg, decrypt(key, encrypted))
}

func TestFindKey(t *testing.T) {
	key := "123456"
	msg := "asdfgh"
	encrypted := "677660656569"

	assert.Equal(t, key, findKey(msg, encrypted))
}
