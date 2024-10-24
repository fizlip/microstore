package services

import (
    "github.com/fizlip/microstore/internal/models"
)

// VendorController is an interface that all Vendors implement
type ReviewController interface {
    GetReviews(toolID string) ([]models.Tool, error)
}

type ReviewService struct{
    reviewStorage models.DatabaseReadWriter[models.Review] 
}

func (r *ReviewService) GetReviews(id string) ([]models.Review, error) {
    return r.reviewStorage.Read(id) 
}

func NewReviewService() (*ReviewService, error) {
    reviewT, err := models.NewMockDynamoTable[models.Review]("f4-tools", "user", "uid")
    if err != nil {
        return &ReviewService{}, err 
    }
    reviewT, err = models.NewMockDynamoTable[models.Review]("f4-reviews", "user", "uid")
    if err != nil {
        return &ReviewService{}, err 
    }
    vendor := ReviewService{reviewT}
    return &vendor, nil
}
