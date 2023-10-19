package main

import (
	"encoding/base64"
	"fmt"
	"log"
)

func main() {
	msg := "Hello. I am happy to learning things really hands-on and thank you for sharing this info with me. I am really grateful and this will go a long way in helping me understand the basics of encryption. Lets gets started."
	encoded := encoded(msg)
	fmt.Println("Encoded message: ", encoded)

	s, err := decode(encoded)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Decoded message: ", s)
}

func encoded(msg string) string {
	return base64.URLEncoding.EncodeToString([]byte(msg))
}

func decode(encoded string) (string, error) {
	s, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("couldnt decode string")
	}
	return string(s), nil
}
