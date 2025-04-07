// package services contains implementaions of models and interfaces
// and provides the business logic of the API
package services

import (
    "errors"
    "sort"

    "github.com/fizlip/microstore/internal/models"

    // log "github.com/sirupsen/logrus"
)

type DBItem struct {
    UID string
}

type Store struct {
    StoreStorage            models.DatabaseReadWriter[models.Tool]
    TrendingBusinessStorage models.DatabaseReadWriter[DBItem]
}

func (s *Store) GetProduct(uuid string) (models.Tool, error) {
    product, err := s.StoreStorage.Read(uuid)
    if len(product) == 0 {
        return models.Tool{}, err
    }
    return product[0], err
}

func (s *Store) GetUTIProduct(uti string) (models.Tool, error) {
    product, err := s.StoreStorage.ReadCol("uti", uti)
    if len(product) == 0 {
        return models.Tool{}, err
    }
    return product[0], err
}


func (s *Store) GetTrendingBusiness() ([]models.Tool, error){
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

// GetRecent will get all actions in the tools table and sort it based on timestamp created from 
// highest to lowest
func (s *Store) GetRecent() ([]models.Tool, error){
    var filteredResult []models.Tool

    items, err := s.StoreStorage.ReadAll()
    if err != nil {
        return []models.Tool{}, errors.New("Could not read trending")
    }

    for _, item := range items {
        if item.Verified {
            filteredResult = append(filteredResult, item)
        }
    }

    sort.Slice(filteredResult, func(i, j int) bool {
        return filteredResult[i].Created > filteredResult[j].Created
    })

    return filteredResult, nil
}

func (s *Store) VerifyProduct(uti string) (error) {

    item, err := s.StoreStorage.GetGSIItem("uti-index", "uti", uti)
    if err != nil {
        return errors.New("Could not read product")
    }

    err = s.StoreStorage.UpdateAll(item.UID, "B,verified", "true")

    return err
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
