package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// use godot package to load/read the .env file and
// return the value of the key
func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load("../../.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func main() {
	port := goDotEnvVariable("PORT")
	router := NewRouter()

	log.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(port, router))
}
