package scripts

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func SimulateExecutions(database *sql.DB) {
	currentTime := time.Now()
	for i := 0; i < 50; i++ {
		simulateGenerateKey(database, currentTime)
		simulateGenerateSignature(database, currentTime)
		simulateVerifySignature(database, currentTime)

		currentTime = currentTime.Add(time.Minute)
		fmt.Println("* Test nº = ", currentTime)
	}
}

func simulateGenerateKey(database *sql.DB, currentTime time.Time) {
	keyTestData := GenerateKeyTest()
	for i := 0; i < len(keyTestData); i++ {
		result := keyTestData[i]
		InsertDB(database, "GenerateKey", currentTime, result)
	}
}

func simulateGenerateSignature(database *sql.DB, currentTime time.Time) {
	hash := generateRandomHash()
	hexHash := convertToHexadecimal(hash)

	signatureTestData := GenerateSignatureTest(hash)
	for i := 0; i < len(signatureTestData); i++ {
		result := signatureTestData[i]
		InsertDBHash(database, "GenerateSignature", currentTime, result, hexHash)
	}
}

func simulateVerifySignature(database *sql.DB, currentTime time.Time) {
	hash := generateRandomHash()
	hexHash := convertToHexadecimal(hash)

	verifySignatureTestData := VerifySignatureTest(hash)
	for i := 0; i < len(verifySignatureTestData); i++ {
		result := verifySignatureTestData[i]
		InsertDBHash(database, "VerifySignature", currentTime, result, hexHash)
	}
}

// Generate 32 random bytes
func generateRandomHash() []byte {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil
	}

	// Hash the random bytes using SHA-256
	hash := sha256.Sum256(randomBytes)

	return hash[:]
}

func convertToHexadecimal(hash []byte) string {
	hexString := "0x" + hex.EncodeToString(hash[:])

	return hexString
}
