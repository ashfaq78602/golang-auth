package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Person struct {
	First string `json:"first"`
}

func main() {
	http.HandleFunc("/encode", foo)
	http.HandleFunc("/decode", bar)
	http.ListenAndServe(":8080", nil)
}

func foo(w http.ResponseWriter, r *http.Request) {
	p1 := Person{
		First: "Ashfaq",
	}
	p2 := Person{
		First: "Reynis",
	}
	people := []Person{p1, p2}

	err := json.NewEncoder(w).Encode(people)
	if err != nil {
		log.Println("Bad data is encoded", err)
	}
}

func bar(w http.ResponseWriter, r *http.Request) {
	people := []Person{}
	err := json.NewDecoder(r.Body).Decode(&people)
	if err != nil {
		log.Println("Bad decoded data", err)
	}
	log.Println(people)
}
