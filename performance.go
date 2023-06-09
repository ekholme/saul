package saul

import (
	"context"
	"encoding/csv"
	"errors"
	"log"
	"os"
	"sort"
	"strconv"

	firestore "cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

const perfColl = "performances"

type Performance struct {
	SchName        string  `json:"schName"`
	Test           string  `json:"test"`
	ItemDescriptor string  `json:"itemDescriptor"`
	BestPractice   string  `json:"bestPractice"`
	Q              float64 `json:"q"`
}

// helper to facilitate passing data to the html template
type PerformanceRequest struct {
	URL          string
	Performances []*Performance
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

// get performances by school and test
func (ps *PerformanceService) GetPerfBySchoolAndTest(ctx context.Context, sch string, tst string) ([]*Performance, error) {

	iter := ps.Client.Collection(perfColl).Where("SchName", "==", sch).Where("Test", "==", tst).Documents(ctx)

	defer iter.Stop()

	var perfs []*Performance

	for {
		doc, err := iter.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, err
		}

		var perf *Performance

		doc.DataTo(&perf)

		perfs = append(perfs, perf)
	}

	//sort of janky, but right now we're returning everything, sorting, then filtering only to the top 3
	//this is a limitation of firestore
	sort.SliceStable(perfs, func(i, j int) bool {
		return perfs[i].Q < perfs[j].Q
	})

	f := len(perfs) - 3

	p := perfs[f:]

	return p, nil
}

func (ps *PerformanceService) GetPerfBySchTestItem(ctx context.Context, sch string, tst string, item string) (*Performance, error) {

	iter := ps.Client.Collection(perfColl).Where("SchName", "==", sch).Where("Test", "==", tst).Where("ItemDescriptor", "==", item).Documents(ctx)

	defer iter.Stop()

	var perfs []*Performance

	for {
		doc, err := iter.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, err
		}

		var perf *Performance

		doc.DataTo(&perf)

		perfs = append(perfs, perf)
	}

	if len(perfs) < 1 {
		return nil, errors.New("couldn't retrieve performances")
	}

	p := perfs[0]

	return p, nil
}

// below this are utility functions to write stuff to firestore
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
// this is a utility function to help with ingesting data
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

		q, err := strconv.ParseFloat(v[4], 64)

		if err != nil {
			log.Fatal("couldn't parse Q as float")
		}

		p := &Performance{
			SchName:        v[0],
			Test:           v[1],
			ItemDescriptor: v[2],
			BestPractice:   v[3],
			Q:              q,
		}

		perfs[k] = p
	}

	return perfs, nil
}
