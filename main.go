package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	password []byte
	First    string
}

type customClaims struct {
	jwt.StandardClaims
	SID string
}

// key is email, value is user
var db = map[string]user{}
var session = map[string]string{}

var key = []byte("my secret key 007 james bond rule the world")

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("sessionID")
	if err != nil {
		c = &http.Cookie{
			Name:  "sessionID",
			Value: "",
		}
	}

	sID, err := parseToken(c.Value)
	if err != nil {
		log.Println("Index parseToken", err)
	}

	var e string
	if sID != "" {
		e = session[sID]
	}

	var f string
	if user, ok := db[e]; ok {
		f = user.First
	}

	errMsg := r.FormValue("msg")
	fmt.Fprint(w, `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Document</title>
	</head>
	<body>
		<h1>IF YOU HAVE A SESSION, HERE IS YOUR EMAIL:`, e, `</h1>
		<h1>IF YOU HAVE A SESSION, HERE IS YOUR NAME:`, f, `</h1>
		<h1>IF THERE WAS ANY MESSAGE, HERE IT IS:`, errMsg, `</h1>
        <h2>REGISTER</h2>
		<form action="/register" method="POST">
			<label for = "first">First</label>
			<input type="text" name="first" placeholder="First" id ="first"><br>
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
	f := r.FormValue("first")
	if f == "" {
		msg := url.QueryEscape("Your first name was empty. It must not be empty.")
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
	db[e] = user{
		password: bsp,
		First:    f,
	}

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

	err := bcrypt.CompareHashAndPassword(db[e].password, []byte(p))
	if err != nil {
		msg := url.QueryEscape("Your email or password didn't match.")
		http.Redirect(w, r, "/msg="+msg, http.StatusSeeOther)
		return
	}

	sUUID := uuid.New().String()
	session[sUUID] = e
	token, err := createToken(sUUID)
	if err != nil {
		log.Println("Couldn't createToken in login", err)
		msg := url.QueryEscape("Our server didnt get enough. Try Again.")
		http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
		return
	}

	c := http.Cookie{
		Name:  "sessionID",
		Value: token,
	}

	http.SetCookie(w, &c)

	msg := url.QueryEscape("You logged in " + e)
	http.Redirect(w, r, "/msg="+msg, http.StatusSeeOther)
}

func createToken(sid string) (string, error) {
	cc := customClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		},
		SID: sid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cc)
	st, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("coudln't sign token in createToken %w", err)
	}
	return st, nil

	// mac := hmac.New(sha256.New, key)
	// mac.Write([]byte(sid))
	//to hex
	//signedMac := fmt.Sprintf("%x", mac.Sum(nil))

	//to base64
	// signedMac := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	// //signedSessionID as base64 | created from sid
	// return signedMac + "|" + sid
}

func parseToken(st string) (string, error) {
	token, err := jwt.ParseWithClaims(st, &customClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("parseWithClaims different algorithms used.")
		}
		return key, nil
	})

	if err != nil {
		return "", fmt.Errorf("couldn't parseWithClaims in parseToken %w", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("couldn't parseWithClaims in parseToken")
	}

	return token.Claims.(*customClaims).SID, nil

	// xs := strings.SplitN(st, "|", 2)
	// if len(xs) != 2 {
	// 	return "", fmt.Errorf("Wrong number of items in string parseToken")
	// }

	// //SIGNEDSESSIONID AS BASE64 | Created from sId
	// b64 := xs[0]
	// xb, err := base64.StdEncoding.DecodeString(b64)
	// if err != nil {
	// 	return "", fmt.Errorf("Couldnt parseToken decodestring: %w", err)
	// }

	// // //signedSessionID as base64 | CREATED FROM SID
	// mac := hmac.New(sha256.New, key)
	// mac.Write([]byte(xs[1]))

	// ok := hmac.Equal(xb, mac.Sum(nil))
	// if !ok {
	// 	return "", fmt.Errorf("Couldnt parseToken not equal signed sid")
	// }
	// return xs[1], nil
}
