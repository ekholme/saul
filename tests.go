package saul

import (
	"context"
	"encoding/csv"
	"os"

	"cloud.google.com/go/firestore"
)

const testColl = "tests"

type Test struct {
	SchName string `json:"schName"`
	Test    string `json:"test"`
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

// method to creat a bunch of tests
func (ts *TestService) CreateTests(ctx context.Context, tsts []*Test) error {

	for _, v := range tsts {
		_, _, err := ts.Client.Collection(testColl).Add(ctx, v)

		if err != nil {
			return err
		}
	}

	return nil
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
