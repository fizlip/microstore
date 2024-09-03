package tools

import (
  "time"
  "github.com/aws/aws-sdk-go-v2/service/ecr/types"
)

type mockDB struct{}

var mockLoginDetails = map[string]LoginDetails{
  "alex": {
    AuthToken: "123ABC",
    Username: "alex",
  },
  "marie": {
    AuthToken: "456DEF",
    Username: "marie",
  },
  "filip": {
    AuthToken: "789GHI",
    Username: "filip",
  },
}

var mockServiceDetails = map[string]ServiceDetails{
  "alex": {
    Service: 100,
    Username: "alex",
  },
  "jason": {
    Service: 200,
    Username: "jason",
  },
  "filip": {
    Service: 300,
    Username: "filip",
  },
}

func (d *mockDB) GetUserLoginDetails(username string) *LoginDetails {
  time.Sleep(time.Second * 1)

  var clientData = LoginDetails{}
  clientData, ok := mockLoginDetails[username]

  if !ok {
    return nil
  }

  return &clientData
}

func (d *mockDB) GetUserServices(username string) *ServiceDetails {
  time.Sleep(time.Second * 1)

  var clientData = ServiceDetails{}
  clientData, ok := mockServiceDetails[username]

  if !ok {
    return nil
  }

  return &clientData
}

func (d *mockDB) GetMetadata(digests []types.ImageIdentifier) []Metadata{
  result := make([]Metadata, len(digests))

  for i, digest := range digests {
    result[i] = Metadata{
      Digest: *digest.ImageDigest,
      Created: 1234,
      Title: "title",
      Description: "Description",
      Reviews: 102,
      Tag: *digest.ImageTag,
      Creator: "filip",
      Requests: 123232,
      Pulls: 12,
      Price: 0.001,      
    }
  }
  
  return result
}

func (d *mockDB) SetupDatabase() error {
  return nil
}

