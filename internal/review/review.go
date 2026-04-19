package review

import (
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Review struct {
	MutationID          uuid.UUID `gorm:"primaryKey"`
	FrameworkMutationID string    // improve compatibility where frameworks issue mutant ids
	Framework           string
	Review              string
}

type Repository struct {
	db *gorm.DB
}

func NewRepository() (*Repository, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&Review{}); err != nil {
		return nil, err
	}
	return &Repository{db: db}, err
}

func (r *Repository) GetReviewByMutationID(id uuid.UUID) (*Review, error) {
	var review Review
	if err := r.db.First(&review, "mutation_id = ?", id).Error; err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *Repository) SaveReview(review *Review) error {
	return r.db.FirstOrCreate(review).Error
}

func (r *Repository) GetReviewsForFramework(framework string) ([]Review, error) {
	var reviews []Review
	err := r.db.Where("framework = ? AND review IS NOT NULL AND review != ''", framework).Find(&reviews).Error
	if err != nil {
		return nil, err
	}
	return reviews, nil
}
