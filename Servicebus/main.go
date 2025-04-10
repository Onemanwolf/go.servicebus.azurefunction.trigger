package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type InvokeRequest struct {
	Data     map[string]json.RawMessage
	Metadata map[string]interface{}
}

func queueHandler(w http.ResponseWriter, r *http.Request) {
    var invokeRequest InvokeRequest

    // Decode the incoming message
    d := json.NewDecoder(r.Body)
    err := d.Decode(&invokeRequest)
    if err != nil {
        // Return 500 to indicate a failure (message will be retried)
        http.Error(w, `{"status":"error","message":"Failed to decode message"}`, http.StatusInternalServerError)
        return
    }

    var parsedMessage string
    err = json.Unmarshal(invokeRequest.Data["queueItem"], &parsedMessage)
    if err != nil {
        // Return 500 to indicate a failure (message will be retried)
        http.Error(w, `{"status":"error","message":"Failed to parse message"}`, http.StatusInternalServerError)
        return
    }

    // Simulate message processing
    fmt.Println("Processing message:", parsedMessage)

    // Example: Simulate a failure for specific messages
    if parsedMessage == "fail" {
        // Return 500 to indicate a failure (message will be retried)
        http.Error(w, `{"status":"error","message":"Processing failed"}`, http.StatusInternalServerError)
        return
    }

    // If processing succeeds, return 200 with a JSON response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status":  "success",
        "message": "Message processed successfully",
    })
}

func main() {
	customHandlerPort, exists := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if !exists {
		customHandlerPort = "8080"
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", queueHandler)
	fmt.Println("Go server Listening on: ", customHandlerPort)
	log.Fatal(http.ListenAndServe(":"+customHandlerPort, mux))
	//MessageProcessorFunction

}