package services 

import (
	log "github.com/sirupsen/logrus"
    "github.com/fizlip/microstore/internal/models"
)

type ITable interface {
	Insert(r models.Review) (bool, error)
	Read(pid string) ([]models.Review, error)
	Update(id string) (bool, error)
	Delete(id string) (bool, error)
}

type DynamoDBReviewTable struct {
	TableName string
}

func (r *DynamoDBReviewTable) Insert(review models.Review) (bool, error) {
	log.Println("Inserting value into table")
	return true, nil
}

func (r *DynamoDBReviewTable) Read(pid string) ([]models.Review, error) {
	log.Println("Reading reviews")
	return []models.Review{}, nil
}

func (r *DynamoDBReviewTable) Update(pid string) (bool, error) {
	log.Println("Updating review")
	return true, nil
}

func (r *DynamoDBReviewTable) Delete(pid string) (bool, error) {
	log.Println("Deleting review")
	return true, nil
}

func NewReviewTable(tableName string) (*ITable, error) {
	var reviewTable ITable = &DynamoDBReviewTable{tableName}
	return &reviewTable, nil
}
