package mocks 

type Product struct {
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

var MockProductCatalog = []Product{
	{
		UID:          "prod-001",
		Title:        "SaaS CRM",
		Description:  "A cloud-based CRM tool for businesses.",
		Rating:       4.7,
		Price:        49.99,
		Thumbnail:    "https://example.com/crm-thumbnail.png",
		Overview:     "Complete customer relationship management system.",
		CommunityURL: "https://community.example.com/crm",
		Reviews:      350,
		Developer:    "CRM Corp",
		PricingType:  "Subscription",
		ReleaseType:  "Stable",
		Version:      "v3.2.1",
		Showcase:     []string{"Feature 1", "Feature 2", "Feature 3"},
		Tags:         []string{"CRM", "Sales", "Cloud"},
		RunEndpoint:  "https://api.example.com/run/crm",
		Pulls:        12000,
	},
	{
		UID:          "prod-002",
		Title:        "Project Management Tool",
		Description:  "An advanced tool for managing projects and teams.",
		Rating:       4.9,
		Price:        29.99,
		Thumbnail:    "https://example.com/pm-thumbnail.png",
		Overview:     "Simplify team collaboration and project tracking.",
		CommunityURL: "https://community.example.com/pm",
		Reviews:      400,
		Developer:    "PM Tools Inc",
		PricingType:  "One-time",
		ReleaseType:  "Beta",
		Version:      "v2.5.3",
		Showcase:     []string{"Gantt Charts", "Task Automation", "Kanban"},
		Tags:         []string{"Project Management", "Collaboration"},
		RunEndpoint:  "https://api.example.com/run/pm",
		Pulls:        9500,
	},
	{
		UID:          "p_001",
		Title:        "Project Management Tool",
		Description:  "An advanced tool for managing projects and teams.",
		Rating:       4.9,
		Price:        29.99,
		Thumbnail:    "https://example.com/pm-thumbnail.png",
		Overview:     "Simplify team collaboration and project tracking.",
		CommunityURL: "https://community.example.com/pm",
		Reviews:      400,
		Developer:    "PM Tools Inc",
		PricingType:  "One-time",
		ReleaseType:  "Beta",
		Version:      "v2.5.3",
		Showcase:     []string{"Gantt Charts", "Task Automation", "Kanban"},
		Tags:         []string{"Project Management", "Collaboration"},
		RunEndpoint:  "https://api.example.com/run/pm",
		Pulls:        9500,
	},
}
