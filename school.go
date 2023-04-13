package saul

import (
	"context"

	firestore "cloud.google.com/go/firestore"
)

const schColl = "schools"

type School struct {
	SchName string `json:"schName"`
}

type SchoolService struct {
	Client *firestore.Client
}

// constructor
func NewSchoolService(client *firestore.Client) *SchoolService {
	return &SchoolService{
		Client: client,
	}
}

// method to get all schools
func (ss *SchoolService) GetAllSchools(ctx context.Context) ([]*School, error) {
	docs, err := ss.Client.Collection(schColl).Documents(ctx).GetAll()

	if err != nil {
		return nil, err
	}

	var schs []*School

	for _, doc := range docs {

		var s *School

		doc.DataTo(&s)

		schs = append(schs, s)
	}

	return schs, nil
}
