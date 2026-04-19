package review

import (
	"reflect"
	"testing"

	"github.com/SecretSheppy/marv/fws/mockfw"
	"github.com/google/uuid"
)

var (
	fw = mockfw.MockFW{}

	review = &Review{
		MutationID: uuid.New(),
		Framework:  fw.Meta().Name,
		Review:     "hello there",
	}
	review2 = &Review{
		MutationID: uuid.New(),
		Framework:  fw.Meta().Name,
		Review:     "review number 2",
	}
	review3 = &Review{
		MutationID: uuid.New(),
		Framework:  fw.Meta().Name,
		Review:     "review number 3",
	}
	review4 = &Review{
		MutationID: uuid.New(),
		Framework:  "not-mock-fw",
		Review:     "review number 4",
	}
	review5 = &Review{
		MutationID: uuid.New(),
		Framework:  fw.Meta().Name,
	}

	reviews = []*Review{review, review2, review3, review4, review5}

	reviewsValidForRetrievalByMockFW = 3
)

func saveReview(t *testing.T, db *Repository) *Review {
	if err := db.SaveReview(review); err != nil {
		t.Fatal(err)
	}
	fetched, err := db.GetReviewByMutationID(review.MutationID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(review, fetched) {
		t.Fatal("fetched review was not equal to created review")
	}
	return review
}

func TestSaveReview(t *testing.T) {
	db, err := NewRepository()
	if err != nil {
		t.Fatal(err)
	}
	saveReview(t, db)
}

func TestEditReview(t *testing.T) {
	db, err := NewRepository()
	if err != nil {
		t.Fatal(err)
	}
	fetched := saveReview(t, db)
	fetched.Review = "new review content"
	if err := db.SaveReview(fetched); err != nil {
		t.Fatal(err)
	}
	saved, err := db.GetReviewByMutationID(review.MutationID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(fetched, saved) {
		t.Fatal("saved review was not equal to updated review")
	}
}

func TestCreateMultipleReviews(t *testing.T) {
	db, err := NewRepository()
	if err != nil {
		t.Fatal(err)
	}
	for _, rev := range reviews {
		if err := db.SaveReview(rev); err != nil {
			t.Fatal(err)
		}
	}
	var count int64
	if err := db.db.Model(&Review{}).Count(&count).Error; err != nil {
		t.Fatal(err)
	}
	if int(count) != len(reviews) {
		t.Fatalf("expected %d records to be created by got %d", len(reviews), count)
	}
}

func TestGetReviewsForFramework(t *testing.T) {
	db, err := NewRepository()
	if err != nil {
		t.Fatal(err)
	}
	for _, rev := range reviews {
		if err := db.SaveReview(rev); err != nil {
			t.Fatal(err)
		}
	}
	revs, err := db.GetReviewsForFramework(fw.Meta().Name)
	if err != nil {
		t.Fatal(err)
	}
	if len(revs) != reviewsValidForRetrievalByMockFW {
		t.Fatalf("expected %d reviews but got %d", reviewsValidForRetrievalByMockFW, len(revs))
	}
}
