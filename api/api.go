package api

import (
  "encoding/json"
  "net/http"
)

// HTTP parameter structs
type SearchParams struct {
  Query string
}

type UserServiceParams struct {
  Username string
}

type GetProductParams struct {
  UID string
}

type UpdateParams struct {
  ID    string `json:"id"`
  Field string `json:"field"`
  Value string `json:"value"`
}

// HTTP response structs
type UserServiceResponse struct {
  Code int
  Balance int64
}

type Success struct {
  Code int
  Message string
}

type Error struct {
  Code int
  Message string
}

// F4 object structs
type Product struct {
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
}

type Taskdef struct {
  ContainerName		string `json:"container_name"`
  Image			string `json:"image"`
  Cpu			int    `json:"cpu"` 
  Memory		int    `json:"memory"` 
  PortMappings		string `json:"port_mappings"` 
  Group			string `json:"group"` 
  Region		string `json:"region"` 
  StreamPrefix		string `json:"stream_prefix"` 
  Family		string `json:"family"` 
  ExecutionRoleARN 	string `json:"execution_role_arn"` 
  TaskRoleARN 		string `json:"task_role_arn"` 
}

type ClusterParams struct {
  TaskDefARN	string `json:"taskdef_arn"`
  Cluster 	string `json:"cluster"`
  LaunchType	string `json:"launch_type"`
  ContainerName	string `json:"container_name"`
  TaskCount	string `json:"task_count"`
}

func writeError(w http.ResponseWriter, message string, code int) {

  resp := Error {
    Code: code,
    Message: message,
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(code)

  json.NewEncoder(w).Encode(resp)

}

var (
  RequestErrorHandler = func(w http.ResponseWriter, err error) {
    writeError(w, err.Error(), http.StatusBadRequest)
  }
  InternalErrorHandler = func(w http.ResponseWriter) {
    writeError(w, "An unexpected Error occured", http.StatusInternalServerError)
  }
)
