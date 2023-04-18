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

	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatalf("couldn't load .env file")
	// }

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

	//data ingest bullshit again
	// ctx := context.Background()

	// perf_data, err := saul.IngestPerformance("./data/toy_data.csv")

	// if err != nil {
	// 	log.Fatalf("Couldn't read in performance data: %v", err)
	// }

	// tst_data, err := saul.IngestTests("./data/tsts.csv")

	// if err != nil {
	// 	log.Fatalf("Couldn't read in test data: %v", err)
	// }

	// err = ps.CreatePerformances(ctx, perf_data)

	// if err != nil {
	// 	log.Fatalf("Couldn't write performance data to firestore: %v", err)
	// }

	// err = ts.CreateTests(ctx, tst_data)

	// if err != nil {
	// 	log.Fatalf("Couldn't write test data to firestore: %v", err)
	// }

	s.Run()
}
