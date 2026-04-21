package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Restartor/web-streaming/routes"
)

func main() {
	router := routes.NewRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, router))
}
