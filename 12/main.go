package main

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"math/big"
)

func main() {
	m1b, err := ioutil.ReadFile("plaintexts/test1.txt")
	assertNoError(err)
	m2b, err := ioutil.ReadFile("plaintexts/test2.txt")
	assertNoError(err)
	c1b, err := ioutil.ReadFile("ciphered/test1.txt")
	assertNoError(err)
	c2b, err := ioutil.ReadFile("ciphered/test2.txt")
	assertNoError(err)

	m1 := new(big.Int).SetBytes(m1b)
	m2 := new(big.Int).SetBytes(m2b)

	c1 := new(big.Int).SetBytes(c1b)
	c2 := new(big.Int).SetBytes(c2b)

	// exps:
	// https://security.stackexchange.com/questions/2335/should-rsa-public-exponent-be-only-in-3-5-17-257-or-65537-due-to-security-c
	for _, exp := range []int64{3, 5, 17, 257, 65537} {
		e := big.NewInt(exp)

		m1e := new(big.Int).Exp(m1, e, nil)
		m2e := new(big.Int).Exp(m2, e, nil)

		d1 := new(big.Int).Sub(c1, m1e)
		d2 := new(big.Int).Sub(c2, m2e)

		// https://crypto.stackexchange.com/questions/43583/deduce-modulus-n-from-public-exponent-and-encrypted-data
		modulus := new(big.Int).GCD(nil, nil, d1, d2)

		key := &rsa.PublicKey{
			N: modulus,
			E: int(e.Int64()),
		}
		again1 := encrypt_RSA(key, m1b)
		equals := string(c1b) == string(again1)
		fmt.Println("modulus size", modulus.BitLen(), " for ", exp)
		fmt.Println(modulus)
		if equals {
			fmt.Println(modulus)
		}
	}
}

func encrypt_RSA(pub *rsa.PublicKey, data []byte) []byte {
	encrypted := new(big.Int)
	e := big.NewInt(int64(pub.E))
	payload := new(big.Int).SetBytes(data)
	encrypted.Exp(payload, e, pub.N)
	return encrypted.Bytes()
}

func assertNoError(err error) {
	if err != nil {
		panic(err)
	}
}
