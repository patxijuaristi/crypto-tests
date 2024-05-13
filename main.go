package main

import (
	"bytes"
	"crypto/test/scripts"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Create a new Gorilla Mux router
	router := mux.NewRouter()

	// Define your API endpoints
	router.HandleFunc("/generateKey", generateKeyHandler).Methods("GET")
	router.HandleFunc("/generateSignature", generateSignatureHandler).Methods("GET")
	router.HandleFunc("/verifySignature", verifySignatureHandler).Methods("GET")
	router.HandleFunc("/keySignatureSizes", keySignatureSizesHandler).Methods("GET")
	router.HandleFunc("/test", testApi).Methods("POST")

	// Start the HTTP server
	fmt.Println("Server listening on port 8080...")
	http.ListenAndServe(":8080", router)
}

// Handler functions for each endpoint

func generateKeyHandler(w http.ResponseWriter, r *http.Request) {
	// Generate key test data
	keyTestData := scripts.GenerateKeyTest()

	// Convert key test data to JSON
	jsonResponse, err := json.Marshal(keyTestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set content type and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func generateSignatureHandler(w http.ResponseWriter, r *http.Request) {
	// Generate signature test data
	signatureTestData := scripts.GenerateSignatureTest()

	// Convert signature test data to JSON
	jsonResponse, err := json.Marshal(signatureTestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set content type and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func verifySignatureHandler(w http.ResponseWriter, r *http.Request) {
	// Verify signature test data
	verifySignatureTestData := scripts.VerifySignatureTest()

	// Convert signature test data to JSON
	jsonResponse, err := json.Marshal(verifySignatureTestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set content type and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func keySignatureSizesHandler(w http.ResponseWriter, r *http.Request) {
	// Keya and signature size test data
	keySignatureSizesTestData := scripts.KeySignatureSizes()

	// Convert sizes test data to JSON
	jsonResponse, err := json.Marshal(keySignatureSizesTestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set content type and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func testApi(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Define a struct to unmarshal the request body into
	var requestData struct {
		Hash         []byte `json:"hash"`
		FunctionName string `json:"functionName"`
	}

	// Unmarshal request body into the struct
	if err := json.Unmarshal(body, &requestData); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	hash := requestData.Hash
	functionName := requestData.FunctionName

	fmt.Println("Entrado al endpoint")

	fmt.Println(hash)
	fmt.Println(functionName)

	// Respond with a success message
	response := map[string]string{"message": "Success"}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	// Set Content-Type header and send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := io.Copy(w, bytes.NewReader(jsonResponse)); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func test() {
	scripts.Testing()
}
