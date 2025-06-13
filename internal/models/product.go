// package models contains definitions of common types and models used accross the API.
package models

// Tool is a type that describes a tool available to a f4rmer.
type Tool struct {
	UID                  string     `dynamodbav:"uid" json:"uid"`
	Title                string     `dynamodbav:"title" json:"title"`
	UTI                  string     `dynamodbav:"uti" json:"uti"`
	Description          string     `dynamodbav:"description" json:"description"`
	Rating               float64    `dynamodbav:"rating" json:"rating"`
	Price                float64    `dynamodbav:"price" json:"price"`
	Thumbnail            string     `dynamodbav:"thumbnail" json:"thumbnail"`
	Overview             string     `dynamodbav:"overview" json:"overview"`
	CommunityURL         string     `dynamodbav:"communityURL" json:"communityURL"`
	Reviews              int        `dynamodbav:"reviews" json:"reviews"`
	Developer            string     `dynamodbav:"developer" json:"developer"`
	PricingType          string     `dynamodbav:"pricingType" json:"pricingType"`
	ReleaseType          string     `dynamodbav:"releaseType" json:"releaseType"`
	Version              string     `dynamodbav:"version" json:"version"`
	Showcase             []string   `dynamodbav:"showcase" json:"showcase"`
	Tags                 []string   `dynamodbav:"tags" json:"tags"`
	RunEndpoint          string     `dynamodbav:"run_endpoints" json:"run_endpoint"`
	Pulls                int        `dynamodbav:"pulls" json:"pulls"`
	License              string     `dynamodbav:"license" json:"license"`
	Verified             bool       `dynamodbav:"verified" json:"verified"`
	EndpointNames        []string   `dynamodbav:"endpointNames" json:"endpoint_names"`
	EndpointDescriptions []string   `dynamodbav:"endpointDescriptions" json:"endpoint_descriptions"`
	Endpoints            []Endpoint `dynamodbav:"endpoints" json:"endpoints"`
	Created              int        `dynamodbav:"created" json:"created"`
	Deployed             bool       `dynamodbav:"deployed" json:"deployed"`
	DeploymentType       string     `dynamodbav:"deployment_type" json:"deployment_type"`
	Public               bool       `dynamodbav:"public" json:"public"`
}

type ToolData struct {
	ProductID   string
	ProductName string
	Price       float64
	NumRequests int
}

type Endpoint struct {
	Name          string             `dynamodbav:"name" json:"name"`
	Parameters    []ParameterWrapper `dynamodbav:"parameters" json:"parameters"`
	Description   string             `dynamodbav:"description" json:"description"`
	Authorization Authorization      `dynamodbav:"authorization" json:"authorization"`
}

type ParameterWrapper struct {
	Parameter Parameter `dynamodbav:"parameter" json:"parameter"`
}

type Parameter struct {
	Name        string `dynamodbav:"name" json:"name"`
	Type        string `dynamodbav:"type" json:"type"`
	Description string `dynamodbav:"description" json:"description"`
}

type Authorization struct {
	Provider string `dynamodbav:"provider" json:"provider"`
	Type     string `dynamodbav:"type" json:"type"`
	Scope    string `dynamodbav:"scope" json:"scope"`
}
