/**
api package contains types that are used in server-client communication
**/
package api

import (
	"encoding/json"
	"net/http"
)

type UsernameParams struct {
	Username string
}

type UpdateReviewParams struct {
    TID     string `json:"tid"`
    UID     string `json:"uid"`
    Col     string `json:"col"`
    Val     string `json:"val"`
}

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
	Code    int
	Balance int64
}

type Success struct {
	Code    int
	Message string
}

type Error struct {
	Code    int
	Message string
}

type Taskdef struct {
	ContainerName    string `json:"container_name"`
	Image            string `json:"image"`
	Cpu              int    `json:"cpu"`
	Memory           int    `json:"memory"`
	PortMappings     string `json:"port_mappings"`
	Group            string `json:"group"`
	Region           string `json:"region"`
	StreamPrefix     string `json:"stream_prefix"`
	Family           string `json:"family"`
	ExecutionRoleARN string `json:"execution_role_arn"`
	TaskRoleARN      string `json:"task_role_arn"`
}

type ClusterParams struct {
	TaskDefARN    string `json:"taskdef_arn"`
	Cluster       string `json:"cluster"`
	LaunchType    string `json:"launch_type"`
	ContainerName string `json:"container_name"`
	TaskCount     string `json:"task_count"`
}

func writeError(w http.ResponseWriter, message string, code int) {

	resp := Error{
		Code:    code,
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
    AuthenticationErrorHandler = func(w http.ResponseWriter, err error) {
		writeError(w, err.Error(), http.StatusUnauthorized)
	}
    AuthorizationErrorHandler = func(w http.ResponseWriter, err error) {
		writeError(w, err.Error(), http.StatusForbidden)
	}
	InternalErrorHandler = func(w http.ResponseWriter) {
		writeError(w, "An unexpected error occured", http.StatusInternalServerError)
	}
    PayloadTooLargeHandler = func(w http.ResponseWriter, err error) {
		writeError(w, err.Error(), http.StatusRequestEntityTooLarge)
	}
    CustomInternalErrorHandler = func(w http.ResponseWriter, err error) {
		writeError(w, err.Error(), http.StatusInternalServerError)
	}
)
