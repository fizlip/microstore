package tools

import (
  "time"
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

func (d *mockDB) SetupDatabase() error {
  return nil
}

