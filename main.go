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
	http.HandleFunc("/login", login)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	errMsg := r.FormValue("msg")
	fmt.Fprint(w, `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Document</title>
	</head>
	<body>
		<h1>IF THERE WAS ANY MESSAGE, HERE IT IS:`, errMsg, `</h1>
        <h2>REGISTER</h2>
		<form action="/register" method="POST">
			<input type="email" name="email"><br>
			<input type="password" name="password"><br>
			<input type="submit">
		</form>
    <h1>LOG IN</h1>
    <form action="/login" method="POST">
        <input type="email" name="email" id=""><br>
        <input type="password" name="password" id=""><br>
        <input type="submit">
    </form>
	</body>
	</html>`)
}

func register(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		msg := url.QueryEscape("Your method was not post")
		http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
		return
		//RETURN must be explicitly used
		//since rediret just sets up the response
		//to redirect the client
		//only after return
		//client is usually redirected to the site
	}

	e := r.FormValue("email")
	if e == "" {
		msg := url.QueryEscape("Your email was empty. It must not be empty")
		http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
		return
	}
	p := r.FormValue("password")
	if p == "" {
		msg := url.QueryEscape("Your password was empty. It must not be empty.")
		http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
		return
	}

	bsp, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		msg := "There was an internal server error."
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	log.Println("password", e)
	log.Println("brcypted", bsp)
	db[e] = bsp

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		msg := url.QueryEscape("Your method was not post")
		http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
		return
		//RETURN must be explicitly used
		//since rediret just sets up the response
		//to redirect the client
		//only after return
		//client is usually redirected to the site
	}

	e := r.FormValue("email")
	if e == "" {
		msg := url.QueryEscape("Your email was empty. It must not be empty")
		http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
		return
	}
	p := r.FormValue("password")
	if p == "" {
		msg := url.QueryEscape("Your password was empty. It must not be empty.")
		http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
		return
	}

	if _, ok := db[e]; !ok {
		msg := url.QueryEscape("Your email and password didn't match.")
		http.Redirect(w, r, "/?msg"+msg, http.StatusSeeOther)
		return
	}

	err := bcrypt.CompareHashAndPassword(db[e], []byte(p))
	if err != nil {
		msg := url.QueryEscape("Your email or password didn't match.")
		http.Redirect(w, r, "/msg="+msg, http.StatusSeeOther)
		return
	}

	msg := url.QueryEscape("You logged in " + e)
	http.Redirect(w, r, "/msg="+msg, http.StatusSeeOther)

}
