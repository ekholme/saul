package main

import (
	"log"
	"os"

	"github.com/ekholme/saul"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("couldn't load .env file")
	}

	apiKey := os.Getenv("OPENAI_API_KEY")

	client := openai.NewClient(apiKey)

	r := mux.NewRouter()

	s := saul.NewServer(r, client)

	s.Run()
}
