package f4rmerhandlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/internal/services"
	log "github.com/sirupsen/logrus"
)

func HandleUpdate(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	UTI := r.URL.Query().Get("uti")
	log.Println("UTI: ", UTI)
	if sessionID == "" {
		http.Error(w, "missing session_id", http.StatusBadRequest)
		return
	}

	// the tns is used to get the tools IP (like a DNS in classical web)
	var tns *services.ToolNameService
	tns, err := services.NewToolNameService()
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
	result, err := (*tns).GetTool(UTI)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

	baseURL := "http://" + result[0].IPs[0] + ":8080"

	w.WriteHeader(http.StatusOK)
	url := baseURL + "/messages/?session_id=" + sessionID
	req, err := http.NewRequest("POST", url, r.Body)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	// Create HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Error sending request: ", err)
		return
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error reading response: ", err)
		return
	}

	log.Println("Initializing MCP server connection for: ", sessionID)
	err = json.NewEncoder(w).Encode("Accepted")
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}
