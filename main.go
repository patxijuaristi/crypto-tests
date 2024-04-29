package main

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/test/scripts/dilithium"
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
		fmt.Println("4 - Key and Signature Sizes")
		fmt.Println("5 - Testing")
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
		case 4:
			keySignatureSizes()
		case 5:
			testing()
		case 0:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
			continue // this should be removed if the following two lines are active
		}
		//fmt.Print("Press Enter to continue...")
		//fmt.Scanln() // Wait for user to press Enter before clearing the screen
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
	measureExecutionTime(wrappedSignECDSA, "GenerateSignature - ECDSA", 10)
	memory_usage()

	// SPHINCS+ algorithm signature generation
	cpu_profiling()
	sk, _ := sphincs.GenerateKeySPHINCS()
	wrappedSignSPHINCS := func() { sphincs.SignSPHINCS(generateRandomHash(), sk) }
	measureExecutionTime(wrappedSignSPHINCS, "GenerateSignature - SPHINCS", 10)
	memory_usage()
}

func verifySignatureTest() {
	hash := generateRandomHash()
	// ECDSA algorithm signature verification
	cpu_profiling()
	key, _ := ecdsa.GenerateKeyECDSA()
	pubkey := ecdsa.FromECDSAPub(&key.PublicKey)
	signature, _ := ecdsa.SignECDSA(hash, key)
	signature = signature[:len(signature)-1] // remove recovery id
	wrappedSignECDSA := func() { ecdsa.VerifySignatureECDSA(pubkey, hash, signature) }
	measureExecutionTime(wrappedSignECDSA, "VerifySignature - ECDSA", 100)
	memory_usage()

	// SPHINCS+ algorithm signature verification
	cpu_profiling()
	sk, pk := sphincs.GenerateKeySPHINCS()
	signature2 := sphincs.SignSPHINCS(hash, sk)
	wrappedSignSPHINCS := func() { sphincs.VerifySignatureSPHINCS(hash, signature2, pk) }
	measureExecutionTime(wrappedSignSPHINCS, "VerifySignature - SPHINCS", 100)
	memory_usage()
}

func keySignatureSizes() {
	hash := generateRandomHash()
	// ECDSA
	key, _ := ecdsa.GenerateKeyECDSA()
	pubkey := ecdsa.FromECDSAPub(&key.PublicKey)
	pubkey = pubkey[:len(pubkey)-1] // remove recovery id
	signature, _ := ecdsa.SignECDSA(hash, key)
	signature = signature[:len(signature)-1] // remove recovery id
	printKeySignatureSizes("ECDSA", len(ecdsa.FromECDSA(key)), len(pubkey), len(signature))

	//SPHINCS
	sk, pk := sphincs.GenerateKeySPHINCS()
	skBytes, pkBytes := sphincs.KeysToBytes(sk, pk)
	signature2 := sphincs.SignSPHINCS(hash, sk)
	sigBytes := sphincs.SignatureToBytes(signature2)
	printKeySignatureSizes("SPHINCA", len(skBytes), len(pkBytes), len(sigBytes))

	//DILITHIUM
	for i := 0; i < 5; i++ {
		mode := dilithium.GetCurrentDilithiumMode()
		pk, sk, _ := dilithium.GenerateKeyDilithium()
		printKeySignatureSizes(mode, len(sk.Bytes()), len(pk.Bytes()), len("xxxxxxxxxxxxxx"))
		dilithium.ChangeDilithiumMode()
	}
}

func printKeySignatureSizes(algorithm string, pkBytes int, skBytes int, signatureBytes int) {
	fmt.Printf("\n** %s", algorithm)
	fmt.Printf("\n - Private key --> %d", skBytes)
	fmt.Printf("\n - Public key ---> %d", pkBytes)
	fmt.Printf("\n - Signature ----> %d", signatureBytes)
	fmt.Printf("\n")
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

func testing() {

	for i := 0; i < 5; i++ {
		fmt.Println(dilithium.GetCurrentDilithiumMode())
		pk, sk, _ := dilithium.GenerateKeyDilithium()

		fmt.Printf("Public key size %d ", len(pk.Bytes()))
		fmt.Printf("\nPrivate key size %d ", len(sk.Bytes()))

		dilithium.ChangeDilithiumMode()
		fmt.Printf("\n")
	}
}
