package saul

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"strconv"

	firestore "cloud.google.com/go/firestore"
)

const perfColl = "performances"

const projectID = "saul-lesson-planner"

type Performance struct {
	SchName        string  `json:"schName"`
	Test           string  `json:"test"`
	ItemDescriptor string  `json:"itemDescriptor"`
	Q              float64 `json:"q'`
}

type PerformanceService struct {
	Client *firestore.Client
}

// constructor for performance service
func NewPerformanceService(client *firestore.Client) *PerformanceService {
	return &PerformanceService{
		Client: client,
	}
}

func NewFirestoreClient() *firestore.Client {
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, projectID)

	if err != nil {
		log.Fatalf("Couldn't create firestore client: %v", err)
	}

	return client
}

// method to write to firestore
func (ps *PerformanceService) CreatePerformance(ctx context.Context, p *Performance) error {
	_, _, err := ps.Client.Collection(perfColl).Add(ctx, p)

	if err != nil {
		return err
	}

	return nil
}

// method to write a bunch of performances
func (ps *PerformanceService) CreatePerformances(ctx context.Context, perfs []*Performance) error {

	for _, v := range perfs {
		_, _, err := ps.Client.Collection(perfColl).Add(ctx, v)

		if err != nil {
			return err
		}
	}
	return nil
}

// count the number of performances in the collection
func (ps *PerformanceService) CountPerformances(ctx context.Context) (map[string]interface{}, error) {
	query := ps.Client.Collection(perfColl).NewAggregationQuery().WithCount("count")

	res, err := query.Get(ctx)

	if err != nil {
		return nil, err
	}

	return res, nil

}

// helper to ingest csv with performances
func IngestPerformance(path string) ([]*Performance, error) {
	f, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	csvReader := csv.NewReader(f)

	lines, err := csvReader.ReadAll()

	if err != nil {
		return nil, err
	}

	perfs := make([]*Performance, len(lines))

	for k, v := range lines {

		q, err := strconv.ParseFloat(v[3], 64)

		if err != nil {
			log.Fatal("couldn't parse Q as float")
		}

		p := &Performance{
			SchName:        v[0],
			Test:           v[1],
			ItemDescriptor: v[2],
			Q:              q,
		}

		perfs[k] = p
	}

	return perfs, nil
}
