package main

import (
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload" // Loads the .env
)

func main() {
	// Setup the web server
	mainHandler := MainHandler()
	log.Fatal(http.ListenAndServe(":8080", mainHandler))
}
