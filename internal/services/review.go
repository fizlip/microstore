package services

import (
    "fmt"
    "errors"

    "github.com/fizlip/microstore/internal/models"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
    "github.com/google/uuid"
    
    log "github.com/sirupsen/logrus"

)

type ReviewService struct{
    reviewStorage models.DatabaseReadWriter[models.Review] 
    interactionStorage models.DatabaseReadWriter[models.Interaction] 
}

func (r *ReviewService) GetInteractions(uids []string, username string) ([]models.Interaction, error) {
    if username == "" {
        return []models.Interaction{}, nil
    }

    unames := make([]string, len(uids))
    for i := range unames {
        unames[i] = username 
    }
    log.Println("unames: ", unames)
    log.Println("uids: ", uids)
    return r.interactionStorage.BatchRead(uids, unames)   
}

func (r *ReviewService) GetReviews(id string) ([]models.Review, error) {
    return r.reviewStorage.Read(id) 
}

func (r *ReviewService) GetReview(sk, pk string) (models.Review, error) {
    return r.reviewStorage.ReadSpecific(sk, pk) 
}

func (r *ReviewService) Create(review models.Review) error {
    return r.reviewStorage.Create(review, reviewToAttributeMapConverter)
}
    
func (r *ReviewService) Delete(id string) error {
    return r.reviewStorage.Delete(id)
}

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
