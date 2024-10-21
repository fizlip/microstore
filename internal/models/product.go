package models

type Tool struct {
	UID          string   `dynamodbav:"uid" json:"uid"`
	Title        string   `dynamodbav:"title" json:"title"`
	Description  string   `dynamodbav:"description" json:"description"`
	Rating       float64  `dynamodbav:"rating" json:"rating"`
	Price        float64  `dynamodbav:"price" json:"price"`
	Thumbnail    string   `dynamodbav:"thumbnail" json:"thumbnail"`
	Overview     string   `dynamodbav:"overview" json:"overview"`
	CommunityURL string   `dynamodbav:"communityURL" json:"communityURL"`
	Reviews      int      `dynamodbav:"reviews" json:"reviews"`
	Developer    string   `dynamodbav:"developer" json:"developer"`
	PricingType  string   `dynamodbav:"pricingType" json:"pricingType"`
	ReleaseType  string   `dynamodbav:"releaseType" json:"releaseType"`
	Version      string   `dynamodbav:"version" json:"version"`
	Showcase     []string `dynamodbav:"showcase" json:"showcase"`
	Tags         []string `dynamodbav:"tags" json:"tags"`
	RunEndpoint  string   `dynamodbav:"run_endpoints" json:"run_endpoint"`
	Pulls        int      `dynamodbav:"pulls" json:"pulls"`
    License      string   `dynamodbav:"license" json:"license"`
    Verified     bool     `dynamodbav:"verified" json:"verified"`
}

