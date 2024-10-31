package services

import (
    "github.com/fizlip/microstore/internal/models"

    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
    "github.com/google/uuid"
)

type ReviewService struct{
    reviewStorage models.DatabaseReadWriter[models.Review] 
}

func (r *ReviewService) GetReviews(id string) ([]models.Review, error) {
    return r.reviewStorage.Read(id) 
}

func (r *ReviewService) Create(review models.Review) error {
    return r.reviewStorage.Create(review, reviewToAttributeMapConverter)
}
    
func (r *ReviewService) Delete(id string) error {
    return r.reviewStorage.Delete(id)
}

func NewReviewService() (*ReviewService, error) {
    //reviewT, err := models.NewMockDynamoTable[models.Review]("f4-tools", "user", "uid")
    //if err != nil {
    //    return &ReviewService{}, err 
    //}
    reviewT, err := models.NewDynamoTable[models.Review]("f4-reviews", "tool_id", "uid")
    if err != nil {
        return &ReviewService{}, err 
    }
    vendor := ReviewService{reviewStorage: reviewT}
    return &vendor, nil
}

// Utils
func reviewToAttributeMapConverter(review models.Review) map[string]types.AttributeValue {
    return map[string]types.AttributeValue{
        "tool_id": &types.AttributeValueMemberS{Value: review.PID},
        "uid": &types.AttributeValueMemberS{Value: uuid.New().String()},
        "title": &types.AttributeValueMemberS{Value: review.Title},
        "value": &types.AttributeValueMemberS{Value: review.Value},
        "content": &types.AttributeValueMemberS{Value: review.Content},
        "creator": &types.AttributeValueMemberS{Value: review.Creator},
        "date": &types.AttributeValueMemberS{Value: review.Date},
    }
}
