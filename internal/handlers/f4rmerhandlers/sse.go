package f4rmerhandlers

import (
	"bufio"
	"fmt"
	"net/http"
	"sync"
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

	//clientIntilizationRequest := JsonRpcRequest{
	//	JsonRpc: "2.0",
	//	ID:      1,
	//	Method:  "initialize",
	//	Params: Params{
	//		ProtocolVersion: "2025-03-26",
	//		Capabilities: Capabilities{
	//			Type:       "claude-3-5-haiku-20241022",
	//			InputValue: []string{},
	//		},
	//		ClientInfo: ClientInfo{
	//			Name:    "f4rmhouse-client",
	//			Version: "1.0.0",
	//		},
	//	},
	//}
	// Convert to JSON
	//jsonInitilizationRequest, err := json.Marshal(clientIntilizationRequest)
	//if err != nil {
	//	fmt.Printf("Error marshaling JSON: %v\n", err)
	//	return
	//}

	// Flush writer
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Connect to upstream SSE server
	upstreamResp, err := http.Get("http://localhost:8080/sse")
	if err != nil {
		http.Error(w, "Failed to connect to upstream SSE", http.StatusBadGateway)
		return
	}
	defer upstreamResp.Body.Close()

	upstreamReader := bufio.NewReader(upstreamResp.Body)

	// Stream each line from upstream to client
	//
	for {
		line, err := upstreamReader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", err.Error())
			flusher.Flush()
			break
		}

		fmt.Fprintf(w, line)
		fmt.Println("line: ", line)
		flusher.Flush()

		//err = json.NewEncoder(w).Encode(line)
		//if err != nil {
		//	log.Error(err)
		//	api.InternalErrorHandler(w)
		//}

		//if strings.HasPrefix(line, "data: {") {
		//	fmt.Println("Return init from server: ", line)
		//	fmt.Fprintf(w, strings.ReplaceAll(line, "\r\n", ""))
		//	flusher.Flush()
		//}

		//if strings.HasPrefix(line, "data: /messages") {
		//	url := "http://localhost:8080" + strings.ReplaceAll(strings.ReplaceAll(line, "data: ", ""), "\r\n", "")
		//	fmt.Println(url)
		//	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonInitilizationRequest))
		//	if err != nil {
		//		fmt.Printf("Error creating request: %v\n", err)
		//		return
		//	}

		//	req.Header.Set("Content-Type", "application/json")

		//	// Create HTTP client and send the request
		//	client := &http.Client{}
		//	resp, err := client.Do(req)
		//	if err != nil {
		//		fmt.Printf("Error sending request: %v\n", err)
		//		return
		//	}
		//	defer resp.Body.Close()

		//	body, err := ioutil.ReadAll(resp.Body)
		//	if err != nil {
		//		fmt.Printf("Error reading response: %v\n", err)
		//		return
		//	}

		//	fmt.Println("Response Status:", resp.Status)
		//	fmt.Println("Response Body:", string(body))

		//	flusher.Flush()
		//} else {
		//	fmt.Fprint(w, line)
		//	flusher.Flush()
		//}
	}
}
