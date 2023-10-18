package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Person struct {
	First string
}

func main() {

	// p2 := Person{
	// 	First: "Ashfaq",
	// }

	// xp := []Person{p1, p2}

	// bs, err := json.Marshal(xp)
	// if err != nil {
	// 	log.Panic(err)
	// }
	// fmt.Println("Json type of Go Data: ", string(bs))

	// xp2 := []Person{}

	// err = json.Unmarshal(bs, &xp2)
	// if err != nil {
	// 	log.Panic(err)
	// }

	// fmt.Println("Back to GO Data: ", xp2)

	http.HandleFunc("/encode", foo)
	http.HandleFunc("/decode", bar)
	http.ListenAndServe(":8080", nil)
}

func foo(w http.ResponseWriter, r *http.Request) {
	p1 := Person{
		First: "Reynis",
	}
	err := json.NewEncoder(w).Encode(p1)
	if err != nil {
		log.Println("Encoded bad data!", err)
	}
}

func bar(w http.ResponseWriter, r *http.Request) {

}
