package f4rmerhandlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fizlip/microstore/api"
	log "github.com/sirupsen/logrus"
)

func HandleUpdate(w http.ResponseWriter, r *http.Request) {
	log.Println("UPDATE")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache")

	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		http.Error(w, "missing session_id", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "{'jsonrpc': '2.0', 'result': ':)', 'id': 2}")
	url := "http://localhost:8080/messages/?session_id=" + sessionID
	fmt.Println(url)
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
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	fmt.Println(body)

	fmt.Println("Response Status:", resp.Status)
	err = json.NewEncoder(w).Encode("Accepted")
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}
