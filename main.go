package main

import (
	"log"
	"net/http"

	"github.com/Restartor/web-streaming/routes"
)

func main() {
	router := routes.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
