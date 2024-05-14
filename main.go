package main

import (
	"crypto/test/scripts"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var database = scripts.CreateDB()

const (
	nIterationsWeb = 5
	nIterationsApi = 20
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run . [mode]")
		fmt.Println("[mode] 1: Execute API endpoints")
		fmt.Println("[mode] 2: Simulate executions")
		return
	}

	mode := os.Args[1]

	if mode == "1" {
		// Execute the code for API endpoints
		router := mux.NewRouter()
		router.HandleFunc("/generateKey", generateKeyHandler).Methods("GET")
		router.HandleFunc("/generateSignature", generateSignatureHandler).Methods("GET")
		router.HandleFunc("/verifySignature", verifySignatureHandler).Methods("GET")
		router.HandleFunc("/keySignatureSizes", keySignatureSizesHandler).Methods("GET")
		router.HandleFunc("/test", testApi).Methods("POST")
		fmt.Println("Server listening on port 8080...")
		http.ListenAndServe(":8080", router)
	} else if mode == "2" {
		// Execute the code for simulating executions
		scripts.SimulateExecutions(database)
	} else {
		fmt.Println("Invalid mode. Please specify either '1' or '2'.")
	}
}

// Handler functions for each endpoint

// Test generate key modes
func generateKeyHandler(w http.ResponseWriter, r *http.Request) {
	keyTestData := scripts.GenerateKeyTest(nIterationsWeb)
	returnData(w, r, keyTestData)
}

// Test generate signature modes
func generateSignatureHandler(w http.ResponseWriter, r *http.Request) {
	signatureTestData := scripts.GenerateSignatureTest(scripts.GenerateRandomHash(), nIterationsWeb)
	returnData(w, r, signatureTestData)
}

// Test verify signature modes
func verifySignatureHandler(w http.ResponseWriter, r *http.Request) {
	verifySignatureTestData := scripts.VerifySignatureTest(scripts.GenerateRandomHash(), nIterationsWeb)
	returnData(w, r, verifySignatureTestData)
}

// Key and signature sizes test
func keySignatureSizesHandler(w http.ResponseWriter, r *http.Request) {
	keySignatureSizesTestData := scripts.KeySignatureSizes(scripts.GenerateRandomHash())
	returnData(w, r, keySignatureSizesTestData)
}

// Return data to the client
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

// Endpoint that receives the new block validations from the Ethereum network
// and performs the tests with the corresponding block hash
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

	perform_tests(hash, functionName)
}

// Perform the tests with the block hash. Generate signature and verify signature are tested.
// Tests are only performed when a block is validated. The API call is done when a signature is done and verified too,
// but we don't make the tests in these cases to avoid system overloading.
func perform_tests(hash []byte, functionName string) {
	fmt.Println("Performing tests", functionName)

	if functionName == "BlockValidator" {
		currentTime := time.Now()
		fmt.Println("** GENERATE SIGNATURE TEST ", currentTime)
		resultGenerateSig := scripts.GenerateSignatureTest(hash, nIterationsApi)
		fmt.Println("** VERIFY SIGNATURE TEST ", currentTime)
		resultVerifySig := scripts.VerifySignatureTest(hash, nIterationsApi)

		hexHash := convertToHexadecimal(hash)
		for i := 0; i < len(resultGenerateSig); i++ {
			result := resultGenerateSig[i]
			scripts.InsertDBHash(database, "GenerateSignature", currentTime, result, hexHash)
		}
		for i := 0; i < len(resultVerifySig); i++ {
			result := resultVerifySig[i]
			scripts.InsertDBHash(database, "VerifySignature", currentTime, result, hexHash)
		}
	}

	// This is how it should be with all API calls
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

// Convert the hash to hexadecimal to store it in the DB
func convertToHexadecimal(hash []byte) string {
	hexString := "0x" + hex.EncodeToString(hash[:])

	return hexString
}
