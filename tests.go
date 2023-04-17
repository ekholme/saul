package saul

import (
	"context"
	"encoding/csv"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

const testColl = "tests"

type Test struct {
	SchName string `json:"schName"`
	Test    string `json:"test"`
}

// small helper to faciltitate template rendering
type TestRequest struct {
	URL   string `json:"url"`
	Tests []*Test
}

type TestService struct {
	Client *firestore.Client
}

// constructor for test service
func NewTestService(client *firestore.Client) *TestService {
	return &TestService{
		Client: client,
	}
}

// method to create a bunch of tests
func (ts *TestService) CreateTests(ctx context.Context, tsts []*Test) error {

	for _, v := range tsts {
		_, _, err := ts.Client.Collection(testColl).Add(ctx, v)

		if err != nil {
			return err
		}
	}

	return nil
}

// method to get tests by school
func (ts *TestService) GetTestBySchool(ctx context.Context, sch string) ([]*Test, error) {

	iter := ts.Client.Collection(testColl).Where("SchName", "==", sch).Documents(ctx)

	defer iter.Stop()

	var tsts []*Test

	for {
		doc, err := iter.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, err
		}

		var tst *Test

		doc.DataTo(&tst)

		tsts = append(tsts, tst)
	}

	return tsts, nil
}

// helper to read in tests
func IngestTests(path string) ([]*Test, error) {
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

	tsts := make([]*Test, len(lines))

	for k, v := range lines {

		t := &Test{
			SchName: v[0],
			Test:    v[1],
		}

		tsts[k] = t
	}

	return tsts, nil
}
