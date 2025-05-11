// package f4rmer_handlers provides functions used to interact with f4rmers.
package f4rmerhandlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/internal/services"

	log "github.com/sirupsen/logrus"
)

// Run will execute the UTI and increment all necessary analytics metrics
func Run(w http.ResponseWriter, r *http.Request) {
	var err error

	var user *services.Customer
	var instigator = r.Header.Get("X-Username")

	uti := r.URL.Query().Get("uti")
	endpoint := r.URL.Query().Get("endpoint")
	params := r.URL.Query().Get("params")

	// Read body, this will contain variables given to the function that's executed
	// unmarshal body to get JSON

	user, err = services.NewCustomer(instigator)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
	// the tns is used to get the tools IP (like a DNS in classical web)
	var tns *services.ToolNameService
	tns, err = services.NewToolNameService()
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
	result, err := (*tns).GetTool(uti)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

	// Execute the action
	var res string
	if result[0].DeploymentType == "mcp_sse" {
		proxySSEHandler(w, r)
		return
		//res, err = user.RunSSE(result[0].IPs[0], runParams.Endpoint, runParams.Params, result[0].UTI)
		//if err != nil {
		//	log.Error(err)
		//	api.InternalErrorHandler(w)
		//	return
		//}
	} else {
		res, err = user.Run(result[0].IPs[0], endpoint, params, result[0].UTI)
		if err != nil {
			log.Error(err)
			api.InternalErrorHandler(w)
			return
		}
	}

	log.Println(fmt.Sprintf("%s was executed sucessfully", result[0].UTI))
	err = json.NewEncoder(w).Encode(api.RunSuccess{Code: 200, Message: res})
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

func proxySSEHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Check if client supports flushing
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}
	flusher.Flush()

	// Create a request to the source SSE server
	req, err := http.NewRequest("GET", "http://localhost:8080/sse", nil)
	if err != nil {
		log.Printf("Error creating request to source SSE: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Copy relevant headers from the client request
	// For example, authorization headers if needed
	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	// If the source SSE server expects any specific headers
	// req.Header.Set("X-Custom-Header", "value")

	// Create HTTP client with extended timeout for long-lived connections
	client := &http.Client{
		Timeout: 0, // No timeout for SSE connections
	}

	// Connect to the source SSE server
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error connecting to source SSE: %v\n", err)
		http.Error(w, "Could not connect to event source", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Check if the source responded with a successful status
	if resp.StatusCode != http.StatusOK {
		log.Printf("Source SSE server responded with status: %d\n", resp.StatusCode)
		http.Error(w, fmt.Sprintf("Event source responded with status: %d", resp.StatusCode), resp.StatusCode)
		return
	}

	// Channel to notify of client disconnect
	clientGone := r.Context().Done()

	// Create a done channel to signal when goroutines should exit
	done := make(chan struct{})
	defer close(done)

	// Start a goroutine to check for client disconnection
	go func() {
		select {
		case <-clientGone:
			log.Println("Client disconnected")
			// Nothing to do here, the main handler will exit
		case <-done:
			// Handler is done
		}
	}()

	// Read from the source SSE and forward to client
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("Source SSE connection closed")
			} else {
				log.Printf("Error reading from source SSE: %v\n", err)
			}
			return
		}

		// Write the line to the client
		_, err = w.Write(line)
		if err != nil {
			log.Printf("Error writing to client: %v\n", err)
			return
		}

		// If this is the end of an event (empty line), flush the data
		if len(line) <= 2 { // "\n" or "\r\n"
			flusher.Flush()
		}

		// Check if client disconnected
		select {
		case <-clientGone:
			log.Println("Client disconnected while streaming")
			return
		default:
			// Continue streaming
		}
	}
}
