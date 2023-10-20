package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Open("sample.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	h := sha256.New()

	_, err = io.Copy(h, f)
	if err != nil {
		log.Fatalln("couldnt io.copy: ", err)
	}

	fmt.Printf("Here's the type before Sum: %T\n", h)
	fmt.Printf("%v\n", h)
	sb := h.Sum(nil)
	fmt.Printf("Here's the type after Sum: %T\n", sb)
	fmt.Printf("%x\n", sb)

	sb = h.Sum(nil)
	fmt.Printf("Here's the type after second Sum: %T\n", sb)
	fmt.Printf("%x\n", sb)

	sb = h.Sum(sb)
	fmt.Printf("Here's the type after third Sum and passing in sb: %T\n", sb)
	fmt.Printf("%x\n", sb)
}
