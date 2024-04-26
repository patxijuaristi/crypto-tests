package main

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/test/scripts/ecdsa"
	"crypto/test/scripts/sphincs"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

func main() {

	var choice int
	for {
		fmt.Println("\n======= CRYPTOGRAPHY METHODS COMPARISON =======")
		fmt.Println("1 - Generate Key")
		fmt.Println("2 - Generate Signature")
		fmt.Println("3 - Verify Signature")
		fmt.Println("0 - Exit")
		fmt.Println("===============================================")
		fmt.Print("Enter your choice: ")
		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		switch choice {
		case 1:
			generateKeyTest()
		case 2:
			generateSignatureTest()
		case 3:
			verifySignatureTest()
		case 0:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
		fmt.Print("Press Enter to continue...")
		fmt.Scanln() // Wait for user to press Enter before clearing the screen
	}
}

// CPU profiling
func cpu_profiling() {
	cpuProfile, err := os.Create("cpu_profile.prof")
	if err != nil {
		fmt.Println("Error creating CPU profile:", err)
		return
	}
	defer cpuProfile.Close()
	if err := pprof.StartCPUProfile(cpuProfile); err != nil {
		fmt.Println("Error starting CPU profile:", err)
		return
	}
	defer pprof.StopCPUProfile()

	// Trigger garbage collection to reset memory usage metrics
	runtime.GC()
}

// Memory usage
func memory_usage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf(" - Alloc -------> %v KB\n", m.Alloc/1024)
	fmt.Printf(" - TotalAlloc --> %v KB\n", m.TotalAlloc/1024)
	fmt.Printf(" - Sys ---------> %v KB\n", m.Sys/1024)
	fmt.Printf(" - NumGC -------> %v\n", m.NumGC)
	fmt.Printf("------------------------------------------------\n")
}

func measureExecutionTime(fn func(), name string, iterations int) {
	totalTime := time.Duration(0)

	for i := 0; i < iterations; i++ {
		start := time.Now()
		fn()
		totalTime += time.Since(start)
	}

	averageTime := totalTime / time.Duration(iterations)
	fmt.Printf(" - %s execution time (average of %d iterations): %v\n", name, iterations, averageTime)
}

func generateKeyTest() {
	// ECDSA algorithm key generation
	cpu_profiling()
	measureExecutionTime(ecdsa.GenerateKeyECDSAWrapper, "GenerateKey - ECDSA", 10)
	memory_usage()

	// SPHINCS+ algorithm key generation
	cpu_profiling()
	measureExecutionTime(sphincs.GenerateKeySPHINCSWrapper, "GenerateKey - SPHINCS", 10)
	memory_usage()
}

func generateSignatureTest() {
	// ECDSA algorithm signature generation
	cpu_profiling()
	key, _ := ecdsa.GenerateKeyECDSA()
	wrappedSignECDSA := func() { ecdsa.SignECDSA(generateRandomHash(), key) }
	measureExecutionTime(wrappedSignECDSA, "GenerateSignature - ECDSA", 100)
	memory_usage()

	// SPHINCS+ algorithm signature generation
	cpu_profiling()
	sk, _ := sphincs.GenerateKeySPHINCS()
	wrappedSignSPHINCS := func() { sphincs.SignSPHINCS(generateRandomHash(), sk) }
	measureExecutionTime(wrappedSignSPHINCS, "GenerateSignature - SPHINCS", 100)
	memory_usage()
}

func verifySignatureTest() {
	// ECDSA algorithm signature generation
	cpu_profiling()
	key, _ := ecdsa.GenerateKeyECDSA()
	pubkey := ecdsa.FromECDSA(key)
	hash := generateRandomHash()
	signature, _ := ecdsa.SignECDSA(hash, key)
	fmt.Println(hash)
	//wrappedSignECDSA := func() { ecdsa.VerifySignatureECDSA(pubkey, hash, signature) }
	//measureExecutionTime(wrappedSignECDSA, "VerifySignature - ECDSA", 1000)
	fmt.Println(ecdsa.VerifySignatureECDSA(pubkey, hash, signature))
	memory_usage()

	// SPHINCS+ algorithm signature generation
	cpu_profiling()
	sk, pk := sphincs.GenerateKeySPHINCS()
	hash2 := generateRandomHash()
	signature2 := sphincs.SignSPHINCS(hash2, sk)
	wrappedSignSPHINCS := func() { sphincs.VerifySignatureSPHINCS(hash2, signature2, pk) }
	measureExecutionTime(wrappedSignSPHINCS, "VerifySignature - SPHINCS", 100)
	memory_usage()
}

// ECDSA Cryptography testing
func ecdsa_test() {
	fmt.Printf("\n---------------- GETH ECDSA SIGNATURE TESTING ----------------\n\n")

	// Measure execution time for each function
	measureExecutionTime(ecdsa.TestEcrecover, "TestEcrecover", 100)
	measureExecutionTime(ecdsa.TestVerifySignature, "TestVerifySignature", 100)
	measureExecutionTime(ecdsa.TestVerifySignatureMalleable, "TestVerifySignatureMalleable", 100)
	measureExecutionTime(ecdsa.TestDecompressPubkey, "TestDecompressPubkey", 100)
	measureExecutionTime(ecdsa.TestCompressPubkey, "TestCompressPubkey", 100)
	measureExecutionTime(ecdsa.TestPubkeyRandom, "TestPubkeyRandom", 100)

	// Run benchmarks
	ecdsa.BenchmarkEcrecoverSignature()
	ecdsa.BenchmarkVerifySignature()
	ecdsa.BenchmarkDecompressPubkey()
}

// SPHINCS+ Cryptography testing
func sphincs_test() {
	fmt.Printf("\n----------------- SPHINCS+ SIGNATURE TESTING -----------------\n\n")

	// Measure execution time for each function
	measureExecutionTime(sphincs.TestSphincs, "TestSphincs", 100)
}

// Generate 32 random bytes
func generateRandomHash() []byte {
	// Generate 32 random bytes
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil
	}

	// Hash the random bytes using SHA-256
	hash := sha256.Sum256(randomBytes)

	return hash[:]
}

func testing() {
	ecdsa_privatekey, err := ecdsa.GenerateKeyECDSA()
	if err == nil {
		fmt.Println(ecdsa_privatekey.PublicKey.X) // The x-coordinate of the public key point on the elliptic curve.
		fmt.Println(ecdsa_privatekey.PublicKey.Y) // The y-coordinate of the public key point on the elliptic curve.
		fmt.Println(ecdsa_privatekey.D)           // The private key scalar
	} else {
		fmt.Printf("Error")
	}

	fmt.Printf("\n------------------- BIN KEY ------------------------------\n")
	bin_key := ecdsa.FromECDSA(ecdsa_privatekey)
	fmt.Println(bin_key)
	fmt.Printf("\nBinary key size (bytes): %d\n", len(bin_key))

	fmt.Printf("\n---------------- KEY TO ADDRESS --------------------------\n")
	address := ecdsa.PubkeyToAddress(ecdsa_privatekey.PublicKey)
	fmt.Println(address)

	fmt.Printf("\n--------------- SAVE KEY TO FILE ----------------------\n")
	ecdsa.SaveECDSA("./data/ecdsa-private-key.txt", ecdsa_privatekey)

	fmt.Printf("\n--------------- LOAD KEY FROM FILE --------------------\n")
	key, err4 := ecdsa.LoadECDSA("./data/ecdsa-private-key.txt")
	if err4 == nil {
		fmt.Println(key.PublicKey.X) // The x-coordinate of the public key point on the elliptic curve.
		fmt.Println(key.PublicKey.Y) // The y-coordinate of the public key point on the elliptic curve.
		fmt.Println(key.D)           // The private key scalar
	} else {
		fmt.Println(err4)
	}

	fmt.Printf("\n--------------------------------------------------------\n")

	sk, pk := sphincs.GenerateKeySPHINCS()
	fmt.Printf("\nPrivate key (seed): %x\n", sk.SKseed)
	fmt.Printf("Private key: (prf) %x\n", sk.SKprf)
	fmt.Printf("Public key (seed): %x\n", pk.PKseed)
	fmt.Printf("Public key (root): %x\n", pk.PKroot)
}

// Function used for testing
func consumeMemoryAndTime() {
	// Allocate memory
	const size = 1024 * 1024 // 1 MB
	var data []byte
	for i := 0; i < 10*size; i++ {
		data = append(data, 'a')
		_ = data // Discard the result to satisfy static analysis
	}

	// Simulate processing time
	time.Sleep(10 * time.Second)

	// Release memory
	data = nil
}
