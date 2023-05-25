package main

import (
	"html/template"
	"log"
	"os"

	"github.com/ekholme/saul"
	"github.com/gorilla/mux"
	openai "github.com/sashabaranov/go-openai"
)

func main() {

	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("couldn't parse templates: %v", err)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")

	//openai client
	client := openai.NewClient(apiKey)

	//firestore client & performance service
	fsClient := saul.NewFirestoreClient()

	//performance service
	ps := saul.NewPerformanceService(fsClient)

	//test service
	ts := saul.NewTestService(fsClient)

	//school service
	ss := saul.NewSchoolService(fsClient)

	//create router
	r := mux.NewRouter()

	//create server
	s := saul.NewServer(r, client, tmpl, ps, ts, ss)

	s.Run()
}
