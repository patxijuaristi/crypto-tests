package main

import (
	"crypto/test/scripts"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func test() {
	scripts.Testing()
}

func main() {
	// Create a new Gorilla Mux router
	router := mux.NewRouter()

	// Define your API endpoints
	router.HandleFunc("/generateKey", generateKeyHandler).Methods("GET")
	router.HandleFunc("/generateSignature", generateSignatureHandler).Methods("GET")
	router.HandleFunc("/verifySignature", verifySignatureHandler).Methods("GET")
	router.HandleFunc("/keySignatureSizes", keySignatureSizesHandler).Methods("GET")

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
