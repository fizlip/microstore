// package services contains implementaions of models and interfaces
// and provides the business logic of the API
package services

import (
    "errors"

    "github.com/fizlip/microstore/internal/models"

    log "github.com/sirupsen/logrus"
)

type DBItem struct {
    UID string
}

type Store struct {
    StoreStorage            models.DatabaseReadWriter[models.Tool]
    TrendingBusinessStorage models.DatabaseReadWriter[DBItem]
}


func (s *Store) GetTrendingBusiness() ([]models.Tool, error){
    log.Println("Get trending")
    var ids []string

    businessTrending, err := s.TrendingBusinessStorage.ReadAll()
    if err != nil {
        return []models.Tool{}, errors.New("Could not read trending")
    }

    for _, item := range businessTrending {
        ids = append(ids, item.UID)
    }

    ts, err := s.StoreStorage.BatchReadPks(ids)
    if err != nil {
        return []models.Tool{}, errors.New("Could not read trending")
    }

    return ts, nil
}

// NewVendorController is a factory method that initialised a data store and
// returns a Vendor struct. If you need to change the data store for some reason
// you only need to update this method.
func NewStoreController() (*Store, error){
    // Init tools data store connection
    storeT, err := models.NewDynamoTable[models.Tool]("f4-tools-prod", "uid", "")
    if err != nil {
        return &Store{}, err 
    }
    trendingBT, err := models.NewDynamoTable[DBItem]("f4-trending-business", "uid", "")
    if err != nil {
        return &Store{}, err 
    }
    store := Store{
        StoreStorage: storeT,
        TrendingBusinessStorage: trendingBT,
    }
    return &store, nil
}
