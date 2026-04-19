package review

import (
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// only one db object because Marv should only ever require one review database.
var db *gorm.DB

func init() {
	initDB()
}

func initDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.AutoMigrate(&Review{})
	if err != nil {
		panic("failed to create in memory database")
	}
}

type Review struct {
	MutationID uuid.UUID `gorm:"primaryKey"`
	Framework  string
	Review     string
}

func GetReviewByMutationID(id uuid.UUID) (*Review, error) {
	var review Review
	if err := db.First(&review, "mutation_id = ?", id).Error; err != nil {
		return nil, err
	}
	return &review, nil
}

func SaveReview(review *Review) error {
	return db.FirstOrCreate(review).Error
}

func GetReviewsForFramework(framework string) ([]Review, error) {
	var reviews []Review
	if err := db.Where("framework = ? AND review IS NOT NULL AND review != ''", framework).Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}
