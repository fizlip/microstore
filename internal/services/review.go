// package services contains implementaions of models and interfaces
// and provides the business logic of the API
package services

import (
    "fmt"
    "errors"

    "github.com/fizlip/microstore/internal/models"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
    "github.com/google/uuid"
    
    log "github.com/sirupsen/logrus"

)

// ReviewService handles all logic related to reviews
type ReviewService struct{
    reviewStorage models.DatabaseReadWriter[models.Review] 
    interactionStorage models.DatabaseReadWriter[models.Interaction] 
}

// GetInteractions returns likes that a review has received
func (r *ReviewService) GetInteractions(uids []string, username string) ([]models.Interaction, error) {
    if username == "" {
        return []models.Interaction{}, nil
    }

    unames := make([]string, len(uids))
    for i := range unames {
        unames[i] = username 
    }
    return r.interactionStorage.BatchRead(uids, unames)   
}

func (r *ReviewService) DeleteInteractions(username string) (error) {
    userInteractions, err := r.interactionStorage.ReadCol("username", username)
    if err != nil {
        return errors.New("Could not read user interactions")
    }
    for _, interaction := range userInteractions {
        err = r.interactionStorage.DeleteDetailed(interaction.UID, interaction.Username) 
        if err != nil {
            return errors.New(fmt.Sprintf("could not update a reviews: ", err))
        }
    }
    return nil
}

// GetReviews returns all reviews associated with a product
func (r *ReviewService) Anonymize(username string) (error) {
    userReviews, err := r.reviewStorage.ReadCol("creator", username)
    if err != nil {
        return errors.New("Could not read reviews")
    }
    for _, review := range userReviews {
        err = r.reviewStorage.Update(review.PID, review.UID, "S,creator", "anon") 
        if err != nil {
            return errors.New(fmt.Sprintf("could not update a reviews: ", err))
        }
    }
    return nil
}

// GetReviews returns all reviews associated with a product
func (r *ReviewService) GetReviews(id string) ([]models.Review, error) {
    return r.reviewStorage.Read(id) 
}

// GetReview returns a specific review
func (r *ReviewService) GetReview(sk, pk string) (models.Review, error) {
    return r.reviewStorage.ReadSpecific(sk, pk) 
}

// Create will create a review
func (r *ReviewService) Create(review models.Review) error {
    return r.reviewStorage.Create(review, reviewToAttributeMapConverter)
}
    
// Delete will delete a review
func (r *ReviewService) Delete(id string) error {
    return r.reviewStorage.Delete(id)
}

// UpdateRating will increment the rating on a review or decrement it if
// the user has already liked it
func (r *ReviewService) UpdateRating(sk, pk, username string) error {
    var err error
    var add = 1
    
    interaction, err := r.interactionStorage.ReadSpecific(pk, username)

    newInteraction := models.Interaction{UID: pk, Username: username}

    if interaction.UID != "" {
        err := r.interactionStorage.DeleteDetailed(pk, username)
        if err != nil {
            return errors.New(fmt.Sprintf("Could not delete interaction: %s, %s"))
        }
        add = -1
    } else {
        err = r.interactionStorage.Create(newInteraction, interactionToAttributeMapConverter)
        if err != nil {
            return err
        }
        add = 1
    }

    review, err := r.reviewStorage.ReadSpecific(sk, pk)
    if err != nil {
        return err
    } 
    log.Println("Rating review")
    return r.reviewStorage.Update(sk, pk, "N,rating", fmt.Sprintf("%d", review.Rating + add))
}

// NewReviewService returns a new ReviewService
func NewReviewService() (*ReviewService, error) {
    reviewT, err := models.NewDynamoTable[models.Review]("f4-reviews", "tool_id", "uid")
    if err != nil {
        return &ReviewService{}, err 
    }

    interactionT, err := models.NewDynamoTable[models.Interaction]("f4-interactions", "uid", "username")
    if err != nil {
        return &ReviewService{}, err 
    }
    vendor := ReviewService{reviewStorage: reviewT, interactionStorage: interactionT}
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
        "rating": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", 0)},
    }
}

// Utils
func interactionToAttributeMapConverter(inter models.Interaction) map[string]types.AttributeValue {
    return map[string]types.AttributeValue{
        "uid": &types.AttributeValueMemberS{Value: inter.UID},
        "username": &types.AttributeValueMemberS{Value: inter.Username},
    }
}
