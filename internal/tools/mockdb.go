package tools

import (
  "fmt"
  "time"
  "github.com/aws/aws-sdk-go-v2/service/ecr/types"
  "github.com/fizlip/microstore/api"

  log "github.com/sirupsen/logrus"
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

func (d *mockDB) GetMetadata(digests []types.ImageIdentifier) []api.Product{
  result := make([]api.Product, len(digests))

  for i, digest := range digests {
    result[i] = api.Product{
	UID:          *digest.ImageDigest,
	Title:        "Sample Product",
	Description:  "This is a sample product description.",
	Rating:       4.5,
	Price:        99.99,
	Thumbnail:    "https://example.com/sample-thumbnail.jpg",
	Overview:     "A brief overview of the sample product.",
	CommunityURL: "https://community.example.com/sample-product",
	Reviews:      100,
	Developer:    "Sample Developer",
	PricingType:  "One-time purchase",
	ReleaseType:  "Stable",
	Version:      "1.0.0",
	Showcase:     []string{"https://example.com/showcase1.jpg", "https://example.com/showcase2.jpg"},
	Tags:         []string{"sample", "product", "example"},
    }
  }
  
  return result
}

func (d *mockDB) CreateMetadata(product api.Product) error{
  log.Println(fmt.Sprintf("Created OCI with title: '%s' and endpoint: '%s'", product.Title, product.RunEndpoint))  
  return nil 
}

func (d *mockDB) UpdateMetadata(upd api.UpdateParams) error{
  log.Println(fmt.Sprintf("Updating '%s' setting '%s' = '%s'", upd.ID, upd.Field, upd.Value))  
  return nil 
}

func (d *mockDB) GetTrending() [9]api.Product{

  var products [9]api.Product;

  names := []string{
	"Task Master",
	"FitTrack",
	"PhotoStudio",
	"MindSpace",
	"FinanceHub",
	"LanguagePro",
	"GameZilla",
	"Recipe Wizard",
	"Home Design 3D",
  }

  descriptions := []string{
	"Manage and track tasks",
	"Monitor workouts progress with personalized plans.",
	"Edit photos pro with filters and editing tools.",
	"Meditate, relax, and sleep better with guided sessions.",
	"Track expenses, set budgets",
	"Language learning app",
	"Hundreds of casual and competitive games.",
	"Create delicious recipes tailored to your diet",
	"Interior designs with easy-to-use 3D modeling",
  }

  prices := []float64{
	9.99,
	4.99,
	19.99,
	0.00,  // No price provided for MindSpace, using 0.00
	2.99,
	14.99,
	0.99,
	5.49,
	24.99,
  }

  for i, _ := range names {
    product := api.Product{};
    product.Title = names[i];
    product.Description = descriptions[i];
    product.Price = prices[i];
    products[i] = product;
  }

  return products;
}

func (d *mockDB) SetupDatabase() error {
  return nil
}

