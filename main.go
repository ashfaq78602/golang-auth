package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/crypto/bcrypt"
)

var db = map[string][]byte{}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/register", register)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	errMsg := r.FormValue("errormsg")
	fmt.Fprintf(w, `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Document</title>
	</head>
	<body>
		<h1>IF THERE WAS ANY ERROR, HERE IT IS: %s</h1>
		<form action="/register" method="post">
			Email:<input type="email" name="email" id="email"><br>
			Password:<input type="password" name="password" id="password"><br>
			<input type="submit">
		</form>
	</body>
	</html>`, errMsg)
}

func register(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		errorMsg := url.QueryEscape("Your method was not post")
		http.Redirect(w, r, "/?errormsg"+errorMsg, http.StatusSeeOther)
		return
		//RETURN must be explicitly used
		//since rediret just sets up the response
		//to redirect the client
		//only after return
		//client is usually redirected to the site
	}

	e := r.FormValue("email")
	if e == "" {
		errorMsg := url.QueryEscape("Your email was empty. It must not be empty")
		http.Redirect(w, r, "/?errormsg"+errorMsg, http.StatusSeeOther)
		return
	}
	p := r.FormValue("password")
	if p == "" {
		errorMsg := url.QueryEscape("Your password was empty. It must not be empty.")
		http.Redirect(w, r, "/?errormsg"+errorMsg, http.StatusSeeOther)
		return
	}

	bsp, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		errorMsg := "There was an internal server error."
		http.Error(w, errorMsg, http.StatusInternalServerError)
		return
	}
	log.Println("password", e)
	log.Println("brcypted", bsp)
	db[e] = bsp

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
