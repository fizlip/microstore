package models 

type Review struct {
	PID     string `dynamodbav:"tool_id" json:"pid"`
	UID     string `dynamodbav:"uid" json:"uid"`
	Title   string `dynamodbav:"title" json:"title"`
	Value   string `dynamodbav:"value" json:"value"`
	Content string `dynamodbav:"content" json:"content"`
	Creator string `dynamodbav:"creator" json:"creator"`
	Date    string `dynamodbav:"date" json:"date"`
	Rating  int    `dynamodbav:"rating" json:"rating"`
    Rated   bool   `json:"rated"`

}
