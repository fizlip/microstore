package f4rmerhandlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/internal/services"
	log "github.com/sirupsen/logrus"
)

// Constants for SSE source and configuration
const (
	sourceSSEServer = "http://localhost:8080/sse"
	baseSourceURL   = "http://localhost:8080"
)

// Client connection store for tracking active connections
type ClientStore struct {
	sync.RWMutex
	clients map[string]*http.Client
}

// Global store for client connections
var clientStore = ClientStore{
	clients: make(map[string]*http.Client),
}

// Define structs to match your JSON structure
type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Capabilities struct {
	Type       string   `json:"type"`
	InputValue []string `json:"input_value"`
}

type Params struct {
	ProtocolVersion string       `json:"protocolVersion"`
	Capabilities    Capabilities `json:"capabilities"`
	ClientInfo      ClientInfo   `json:"clientInfo"`
}

type JsonRpcRequest struct {
	JsonRpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
}

func SSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache")

	uti := r.URL.Query().Get("uti")

	log.Println("SSE UTI: ", uti)
	pings := 0

	// Flush writer
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
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
	result, err := (*tns).GetTool("ascii_gen")
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

	baseURL := "http://" + result[0].IPs[0] + ":8080"

	log.Println("BASEURL: ", baseURL)

	// Connect to upstream SSE server
	upstreamResp, err := http.Get(baseURL + "/sse")
	if err != nil {
		http.Error(w, "Failed to connect to upstream SSE", http.StatusBadGateway)
		return
	}
	defer upstreamResp.Body.Close()

	upstreamReader := bufio.NewReader(upstreamResp.Body)

	// Stream each line from upstream to client
	for pings < 2 {
		line, err := upstreamReader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", err.Error())
			flusher.Flush()
			break
		}
		if strings.HasPrefix(line, "data: /messages") {
			line = strings.ReplaceAll(line, "\r\n", "") + fmt.Sprintf("&uti=%s\r\n", uti)
		}
		// Return raw text to client
		fmt.Fprintf(w, line)
		flusher.Flush()

		// If too many pings exit endpoint
		if strings.Contains(line, ": ping") {
			pings++
		}
		flusher.Flush()
	}

	err = json.NewEncoder(w).Encode("Terminate connection")
	if err != nil {
		return
	}
}
