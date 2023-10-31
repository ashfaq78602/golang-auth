package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/oauth/github", startGithubOauth)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Document</title>
	</head>
	<body>
        <h2>REGISTER</h2>
		<form action="/oauth/github" method="post">
			<input type="submit" value="Login with GitHub">
	</body>
	</html>`)
}

func startGithubOauth(w http.ResponseWriter, r *http.Request) {

}
