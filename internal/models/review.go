// package models contains definitions of common types and models used accross the API.
package models 

// Review are used to give users due diligence before commiting to 
// using a tool.
type Review struct {
	PID         string `dynamodbav:"tool_id" json:"pid"`
	UID         string `dynamodbav:"uid" json:"uid"`
	Title       string `dynamodbav:"title" json:"title"`
	Value       string `dynamodbav:"value" json:"value"`
	Content     string `dynamodbav:"content" json:"content"`
	Creator     string `dynamodbav:"creator" json:"creator"`
	CreatorUID  string `dynamodbav:"creator" json:"creator_uid"`
	Date        string `dynamodbav:"date" json:"date"`
	Rating      int    `dynamodbav:"rating" json:"rating"`
    Rated       bool   `json:"rated"`

}
