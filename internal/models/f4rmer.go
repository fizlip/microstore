// package models contains definitions of common types and models used accross the API.
package models

// F4rmer is an autonomous LLM controlled agent. It is described by this type.
type F4rmer struct {
    Creator         string   `dynamodbav:"creator" json:"creator"`
    UID             string   `dynamodbav:"uid" json:"uid"`
    Title           string   `dynamodbav:"title" json:"title"`
    JobDescription  string   `dynamodbav:"job_description" json:"jobDescription"`
    Toolbox         []Tool   `dynamodbav:"tool_box" json:"toolbox"`
    Created         string   `dynamodbav:"created" json:"created"`
}
