package api

import (
  log "github.com/sirupsen/logrus"
)

type IProductAnalyticsTable interface {
  Insert(r ProductAnalytics) (bool, error)
  Read(pid string) ([]ProductAnalytics, error)
  Update(id string) (bool, error)
  Delete(id string) (bool, error)
}

type DynamoDBProductAnalyticsTable struct {
  TableName string
}

func (r *DynamoDBProductAnalyticsTable) Insert(p ProductAnalytics) (bool, error){
  log.Println("Inserting value into table")
  return true,nil
} 

func (r *DynamoDBProductAnalyticsTable) Read(pid string) ([]ProductAnalytics, error){
  log.Println("Reading reviews")
  return []ProductAnalytics{},nil
}

func (r *DynamoDBProductAnalyticsTable) Update(pid string) (bool, error){
  log.Println("Updating review")
  return true,nil
}

func (r *DynamoDBProductAnalyticsTable) Delete(pid string) (bool, error){
  log.Println("Deleting review")
  return true,nil
}

func NewProductAnalyticsTable(tableName string) (*IProductAnalyticsTable, error) {
  var analyticsTable IProductAnalyticsTable = &DynamoDBProductAnalyticsTable{tableName}
  return &analyticsTable, nil
}
