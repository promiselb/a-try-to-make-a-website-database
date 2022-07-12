package main

import (
	"net/http"

	"github.com/promiselb/website"
)

func main() {
	website.Port = ":8080"
	server := website.NewServer()
	http.ListenAndServe(website.Port, server)
}
