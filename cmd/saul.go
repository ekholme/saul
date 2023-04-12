package main

import (
	"html/template"
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
		log.Fatalf("couldn't load .env file")
	}

	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("couldn't parse templates: %v", err)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")

	client := openai.NewClient(apiKey)

	r := mux.NewRouter()

	s := saul.NewServer(r, client, tmpl)

	s.Run()
}
