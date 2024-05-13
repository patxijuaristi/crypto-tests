package main

import (
	"crypto/test/scripts"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var database = scripts.CreateDB()

func main() {
	scripts.SimulateExecutions(database)
	/*// Create a new Gorilla Mux router
	router := mux.NewRouter()

	// Define your API endpoints
	router.HandleFunc("/generateKey", generateKeyHandler).Methods("GET")
	router.HandleFunc("/generateSignature", generateSignatureHandler).Methods("GET")
	router.HandleFunc("/verifySignature", verifySignatureHandler).Methods("GET")
	router.HandleFunc("/keySignatureSizes", keySignatureSizesHandler).Methods("GET")
	router.HandleFunc("/test", testApi).Methods("POST")

	// Start the HTTP server
	fmt.Println("Server listening on port 8080...")
	http.ListenAndServe(":8080", router)*/
}

// Handler functions for each endpoint

func generateKeyHandler(w http.ResponseWriter, r *http.Request) {
	// Generate key test data
	keyTestData := scripts.GenerateKeyTest()

	returnData(w, r, keyTestData)
}

func generateSignatureHandler(w http.ResponseWriter, r *http.Request) {
	// Generate signature test data
	signatureTestData := scripts.GenerateSignatureTest(nil)

	returnData(w, r, signatureTestData)
}

func verifySignatureHandler(w http.ResponseWriter, r *http.Request) {
	// Verify signature test data
	verifySignatureTestData := scripts.VerifySignatureTest(nil)

	returnData(w, r, verifySignatureTestData)
}

func keySignatureSizesHandler(w http.ResponseWriter, r *http.Request) {
	// Keya and signature size test data
	keySignatureSizesTestData := scripts.KeySignatureSizes(nil)

	returnData(w, r, keySignatureSizesTestData)
}

func returnData(w http.ResponseWriter, r *http.Request, data []map[string]interface{}) {
	// Convert signature test data to JSON
	jsonResponse, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set content type and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func testApi(w http.ResponseWriter, r *http.Request) {
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

	fmt.Println("----------------------------")
	perform_tests(hash, functionName)
	fmt.Println("----------------------------")
}

func perform_tests(hash []byte, functionName string) {
	fmt.Println("Performing tests", functionName)

	if functionName == "BlockValidator" {
		currentTime := time.Now()
		fmt.Println("** GENERATE SIGNATURE TEST **")
		resultGenerateSig := scripts.GenerateSignatureTest(hash)
		fmt.Println("** VERIFY SIGNATURE TEST **")
		resultVerifySig := scripts.VerifySignatureTest(hash)

		hexHash := convertToHexadecimal(hash)
		for i := 0; i < len(resultGenerateSig); i++ {
			result := resultGenerateSig[i]
			scripts.InsertDBHash(database, "GenerateSignature", currentTime, result, hexHash)
		}
		for i := 0; i < len(resultVerifySig); i++ {
			result := resultVerifySig[i]
			scripts.InsertDBHash(database, "GenerateSignature", currentTime, result, hexHash)
		}
	}

	// This is how it should be
	/*
		switch functionName {
		case "Signature":
			resultGenerateSig := scripts.GenerateSignatureTest(hash)
		case "VerifySignature", "BlockValidator":
			resultVerifySig := scripts.VerifySignatureTest(hash)
		default: //Generate key
			resultGenerateKey := scripts.GenerateKeyTest()
		}
	*/
}

func test() {
	scripts.Testing()
}

func convertToHexadecimal(hash []byte) string {
	hexString := "0x" + hex.EncodeToString(hash[:])

	return hexString
}
