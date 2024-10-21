package models

type Tool struct {
	UID          string   `json:"uid"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Rating       float64  `json:"rating"`
	Price        float64  `json:"price"`
	Thumbnail    string   `json:"thumbnail"`
	Overview     string   `json:"overview"`
	CommunityURL string   `json:"communityURL"`
	Reviews      int      `json:"reviews"`
	Developer    string   `json:"developer"`
	PricingType  string   `json:"pricingType"`
	ReleaseType  string   `json:"releaseType"`
	Version      string   `json:"version"`
	Showcase     []string `json:"showcase"`
	Tags         []string `json:"tags"`
	RunEndpoint  string   `json:"run_endpoint"`
	Pulls        int      `json:"pulls"`
    License      string   `json:"license"`
}

