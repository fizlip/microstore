// package models contains definitions of common types and models used accross the API.
package models 

// ToolName is a struct that defines the TNS table.
// The TNS or Tool Name Service is a mapping between Unique Tool ID and IPs that
// host the tools.
type ToolName struct {
    UTI     string      `dynamodbav:"uti" json:"uti"`
    IPs     []string    `dynamodbav:"ips" json:"ips"`
}
