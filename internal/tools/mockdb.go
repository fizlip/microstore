package tools

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/internal/models"

	log "github.com/sirupsen/logrus"
)

type mockDB struct{}

var mockLoginDetails = map[string]LoginDetails{
	"alex": {
		AuthToken: "123ABC",
		Username:  "alex",
	},
	"marie": {
		AuthToken: "456DEF",
		Username:  "marie",
	},
	"filip": {
		AuthToken: "789GHI",
		Username:  "filip",
	},
}

var mockServiceDetails = map[string]ServiceDetails{
	"alex": {
		Service:  100,
		Username: "alex",
	},
	"jason": {
		Service:  200,
		Username: "jason",
	},
	"filip": {
		Service:  300,
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

func (d *mockDB) GetProduct(uid string) models.Tool {

	return models.Tool{
		UID:          uid,
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
		Showcase:     []string{"https://example.com/showcase1.jpg", "https://example.com/showcase2.jpg"},
		Tags:         []string{"sample", "product", "example"},
	}
}

func (d *mockDB) GetMetadata(digests []types.ImageIdentifier) []models.Tool {
	result := make([]models.Tool, len(digests))

	for i, digest := range digests {
		result[i] = models.Tool{
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
			Showcase:     []string{"https://example.com/showcase1.jpg", "https://example.com/showcase2.jpg"},
			Tags:         []string{"sample", "product", "example"},
		}
	}

	return result
}

func (d *mockDB) CreateMetadata(product models.Tool) error {
	return nil
}

func (d *mockDB) UpdateMetadata(upd api.UpdateParams) error {
	log.Println(fmt.Sprintf("Updating '%s' setting '%s' = '%s'", upd.ID, upd.Field, upd.Value))
	return nil
}

func (d *mockDB) GetTrending() [9]models.Tool {

	var products [9]models.Tool

	names := []string{
		"Tsk Master",
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
		0.00, // No price provided for MindSpace, using 0.00
		2.99,
		14.99,
		0.99,
		5.49,
		24.99,
	}

	for i, _ := range names {
		product := models.Tool{}
		product.UID = strconv.Itoa(i)
		product.Title = names[i]
		product.Description = descriptions[i]
		product.Price = prices[i]
		products[i] = product
	}

	return products
}

func (d *mockDB) SearchProduct(query string) []models.Tool {
	products := []models.Tool{
		{
			UID:          "uuid1", // Replace with actual UUIDs
			Title:        "Super Productivity App",
			Description:  "Boost your productivity with our amazing app!",
			Rating:       2.7,
			Price:        19.99,
			Thumbnail:    "https://example.com/thumbnails/productivity-app.jpg",
			Overview:     "This app helps you manage tasks, track time, and collaborate with your team.",
			CommunityURL: "https://community.superproductivityapp.com",
			Reviews:      1250,
			Developer:    "ProductiveSoft Inc.",
			PricingType:  "Subscription",
			ReleaseType:  "Stable",
			Showcase: []string{
				"https://example.com/showcase/productivity-app-1.jpg",
				"https://example.com/showcase/productivity-app-2.jpg",
				"https://example.com/showcase/productivity-app-3.jpg",
			},
			Tags: []string{"productivity", "value"},
		},
		{
			UID:          "uuid2", // Replace with actual UUIDs
			Title:        "Awesome Game Engine",
			Description:  "Create stunning games with our powerful engine!",
			Rating:       4.9,
			Price:        299.99,
			Thumbnail:    "https://example.com/thumbnails/game-engine.jpg",
			Overview:     "A comprehensive game development platform for both 2D and 3D games.",
			CommunityURL: "https://community.awesomegameengine.com",
			Reviews:      3500,
			Developer:    "GameTech Solutions",
			PricingType:  "One-time purchase",
			ReleaseType:  "Beta",
			Showcase: []string{
				"https://example.com/showcase/game-engine-1.jpg",
				"https://example.com/showcase/game-engine-2.jpg",
				"https://example.com/showcase/game-engine-3.jpg",
			},
			Tags: []string{"gaming", "development"},
		},
		{
			UID:          "uuid3", // Replace with actual UUIDs
			Title:        "Smart Home Controller",
			Description:  "Control your entire home with a single app!",
			Rating:       4.5,
			Price:        49.99,
			Thumbnail:    "https://example.com/thumbnails/smart-home.jpg",
			Overview:     "Integrate all your smart home devices and control them from one central hub.",
			CommunityURL: "https://community.smarthomecontroller.com",
			Reviews:      875,
			Developer:    "SmartTech Innovations",
			PricingType:  "Freemium",
			ReleaseType:  "Stable",
			Showcase: []string{
				"https://example.com/showcase/smart-home-1.jpg",
				"https://example.com/showcase/smart-home-2.jpg",
				"https://example.com/showcase/smart-home-3.jpg",
			},
			Tags: []string{"iot", "smarthome", "fun"},
		},
		{
			UID:          "uuid4", // Replace with actual UUIDs
			Title:        "Photo Editor Pro",
			UTI:          "net-tool",
			Description:  "Edit your photos like a pro with our powerful tools!",
			Rating:       4.8,
			Price:        79.99,
			Thumbnail:    "https://example.com/thumbnails/photo-editor.jpg",
			Overview:     "Advanced photo editing software with AI-powered features and intuitive interface.",
			CommunityURL: "https://community.photoeditorpro.com",
			Reviews:      2100,
			Developer:    "CreativeSoft Labs",
			PricingType:  "Subscription",
			ReleaseType:  "Stable",
			Showcase: []string{
				"https://example.com/showcase/photo-editor-1.jpg",
				"https://example.com/showcase/photo-editor-2.jpg",
				"https://example.com/showcase/photo-editor-3.jpg",
			},
			Tags: []string{"photography", "creativity", "editing", "social"},
		},
		{
			UID:          "uuid5", // Replace with actual UUIDs
			Title:        "Fitness Tracker Deluxe",
			Description:  "Track your fitness journey with precision!",
			Rating:       1.6,
			Price:        0,
			Thumbnail:    "https://example.com/thumbnails/fitness-tracker.jpg",
			Overview:     "Comprehensive fitness tracking app with personalized workout plans and nutrition advice.",
			CommunityURL: "https://community.fitnesstrackerdeluxe.com",
			Reviews:      1800,
			Developer:    "HealthTech Solutions",
			PricingType:  "Free with in-app purchases",
			ReleaseType:  "Stable",
			Showcase: []string{
				"https://example.com/showcase/fitness-tracker-1.jpg",
				"https://example.com/showcase/fitness-tracker-2.jpg",
				"https://example.com/showcase/fitness-tracker-3.jpg",
			},
			Tags: []string{"fitness", "tracker", "deluxe"},
		},
	}
	return products
}

func (d *mockDB) SetupDatabase() error {
	return nil
}
