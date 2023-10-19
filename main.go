package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	msg := "Hello. I am happy to learning things really hands-on and thank you for sharing this info with me. I am really grateful and this will go a long way in helping me understand the basics of encryption. Lets gets started."

	password := "ilovecats"
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Panic("couldn't decrypt password")
	}
	bs := b[:16]
	rslt, err := enDecode(bs, msg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(rslt)

	rslt2, err := enDecode(bs, string(rslt))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(rslt2))
}

func enDecode(key []byte, input string) ([]byte, error) {
	b, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("couldn't newciper %w", err)
	}

	//initialization vector
	iv := make([]byte, aes.BlockSize)

	//create a cipher
	s := cipher.NewCTR(b, iv)

	buff := &bytes.Buffer{}
	sw := cipher.StreamWriter{
		S: s,
		W: buff,
	}
	_, err = sw.Write([]byte(input))
	if err != nil {
		return nil, fmt.Errorf("couldn't sw.Write to streamwriter %w", err)
	}

	return buff.Bytes(), nil
}
