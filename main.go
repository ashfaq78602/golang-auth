package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type myClaims struct {
	jwt.StandardClaims
	Email string
}

const myKey = "i love thursdays when it rains 8273 inches"

func main() {
	http.HandleFunc("/", foo)
	http.HandleFunc("/submit", bar)
	http.ListenAndServe(":8080", nil)
}

func foo(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session")
	if err != nil {
		c = &http.Cookie{}
	}

	ss := c.Value
	afterVerificationToken, err := jwt.ParseWithClaims(ss, &myClaims{}, func(beforeVerificationToken *jwt.Token) (interface{}, error) {
		return []byte(myKey), nil
	})

	// Standard claims has the Valid() error method
	// which means it implements the claims interface

	/* type Claims struct {
		Valid() error
	}
	*/
	// when you parseClaims as with "ParseWithClaims"
	// the Valid() method gets run
	// and if all is well, it returns no error and
	// type TOKEN which has field VALID will be true

	isEqual := afterVerificationToken.Valid && err == nil

	message := "Not logged in!!!"
	if isEqual {
		message = "Logged in"
		claims := afterVerificationToken.Claims.(*myClaims)
		fmt.Println(claims.Email)
		fmt.Println(claims.ExpiresAt)
	}

	html := `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>HMAC Example</title>
	</head>
	<body>
		<p> Cookie value:` + c.Value + `</p>
		<p> Message: ` + message + `</p>
		<form action="/submit" method="post">
			<input type="email" name="email"/>
			<input type="submit" />
		</form>
	</body>
	</html>`

	io.WriteString(w, html)
}

func getJWT(msg string) (string, error) {
	claims := myClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		},
		Email: msg,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	ss, err := token.SignedString([]byte(myKey))

	if err != nil {
		return "", fmt.Errorf("Couldnt get signed string in NewWithClaims %v", err)
	}
	return ss, nil

}

func bar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	if email == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	ss, err := getJWT(email)
	if err != nil {
		http.Error(w, "couldnt getJWT", http.StatusInternalServerError)
		return
	}

	c := http.Cookie{
		Name:  "session",
		Value: ss,
	}

	http.SetCookie(w, &c)
	//fmt.Printf("%v", c.Value)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
