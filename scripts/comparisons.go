package scripts

import (
	"crypto/test/scripts/dilithium"
	"crypto/test/scripts/ecdsa"
	"crypto/test/scripts/falcon"
	"crypto/test/scripts/sphincs"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"sort"
	"time"
)

const (
	nSphincs     = 6
	nDilithium   = 5
	printResults = false
)

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
func memory_usage(print bool) map[string]interface{} {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	if print {
		fmt.Printf(" - Alloc -------> %v KB\n", m.Alloc/1024)
		fmt.Printf(" - TotalAlloc --> %v KB\n", m.TotalAlloc/1024)
		fmt.Printf(" - Sys ---------> %v KB\n", m.Sys/1024)
		fmt.Printf(" - NumGC -------> %v\n", m.NumGC)
		fmt.Printf("------------------------------------------------\n")
	}

	memStats := map[string]interface{}{
		"alloc_kb":       int(m.Alloc / 1024),
		"total_alloc_kb": int(m.TotalAlloc / 1024),
		"sys_kb":         int(m.Sys / 1024),
		"num_gc":         int(m.NumGC),
	}
	return memStats
}

func measureExecutionTime(print bool, fn func(), name string, iterations int) int {
	totalTime := time.Duration(0)

	for i := 0; i < iterations; i++ {
		start := time.Now()
		fn()
		totalTime += time.Since(start)
	}

	averageTime := totalTime / time.Duration(iterations)
	averageMicroseconds := int(averageTime / time.Microsecond)

	if print {
		fmt.Printf(" - %s execution time (average of %d iterations): %d microseconds\n", name, iterations, averageMicroseconds)
	}

	return averageMicroseconds
}

func sortByAlgorithm(result []map[string]interface{}) {
	sort.Slice(result, func(i, j int) bool {
		return result[i]["algorithm"].(string) < result[j]["algorithm"].(string)
	})
}

func GenerateKeyTest(nIterations int) []map[string]interface{} {
	var result []map[string]interface{}

	// ECDSA algorithm key generation
	cpu_profiling()
	ecdsaTime := measureExecutionTime(printResults, ecdsa.GenerateKeyECDSAWrapper, "GenerateKey - ECDSA", nIterations)
	ecdsaStats := memory_usage(printResults)
	ecdsaStats["algorithm"] = "ECDSA"
	ecdsaStats["execution_time"] = ecdsaTime
	result = append(result, ecdsaStats)

	// SPHINCS+ algorithm key generation
	for i := 0; i < nSphincs; i++ {
		mode := sphincs.GetCurrentSphincsMode()
		sphincsTime := measureExecutionTime(printResults, sphincs.GenerateKeySPHINCSWrapper, "GenerateKey - "+mode, nIterations)
		sphincsStats := memory_usage(printResults)
		sphincsStats["algorithm"] = mode
		sphincsStats["execution_time"] = sphincsTime
		result = append(result, sphincsStats)
		sphincs.ChangeSphincsMode()
	}

	// Dilithium algorithm key generation
	for i := 0; i < nDilithium; i++ {
		mode := dilithium.GetCurrentDilithiumMode()
		dilithiumTime := measureExecutionTime(printResults, dilithium.GenerateKeyDilithiumWrapper, "GenerateKey - "+mode, nIterations)
		dilithiumStats := memory_usage(printResults)
		dilithiumStats["algorithm"] = mode
		dilithiumStats["execution_time"] = dilithiumTime
		result = append(result, dilithiumStats)
		dilithium.ChangeDilithiumMode()
	}

	// Falcon algorithm key generation
	seed := make([]byte, 64)
	cpu_profiling()
	wrappedFalconGenKey := func() { falcon.GenerateKey(seed) }
	falconTime := measureExecutionTime(printResults, wrappedFalconGenKey, "GenerateKey - Falcon 1024", nIterations)
	falconStats := memory_usage(printResults)
	falconStats["algorithm"] = "Falcon 1024"
	falconStats["execution_time"] = falconTime
	result = append(result, falconStats)

	sortByAlgorithm(result)
	return result
}

func GenerateSignatureTest(hash []byte, nIterations int) []map[string]interface{} {
	var result []map[string]interface{}

	// ECDSA algorithm signature generation
	cpu_profiling()
	key, _ := ecdsa.GenerateKeyECDSA()
	wrappedSignECDSA := func() { ecdsa.SignECDSA(hash, key) }
	ecdsaTime := measureExecutionTime(printResults, wrappedSignECDSA, "GenerateSignature - ECDSA", nIterations)
	ecdsaStats := memory_usage(printResults)
	ecdsaStats["algorithm"] = "ECDSA"
	ecdsaStats["execution_time"] = ecdsaTime
	result = append(result, ecdsaStats)

	// SPHINCS+ algorithm signature generation
	for i := 0; i < nSphincs; i++ {
		mode := sphincs.GetCurrentSphincsMode()
		cpu_profiling()
		sk, _ := sphincs.GenerateKeySPHINCS()
		wrappedSignSPHINCS := func() { sphincs.SignSPHINCS(hash, sk) }
		sphincsTime := measureExecutionTime(printResults, wrappedSignSPHINCS, "GenerateSignature - "+mode, nIterations)
		sphincsStats := memory_usage(printResults)
		sphincsStats["algorithm"] = mode
		sphincsStats["execution_time"] = sphincsTime
		result = append(result, sphincsStats)
		sphincs.ChangeSphincsMode()
	}

	// Dilithium algorithm signature generation
	for i := 0; i < nDilithium; i++ {
		mode := dilithium.GetCurrentDilithiumMode()
		cpu_profiling()
		_, sk, _ := dilithium.GenerateKeyDilithium()
		wrappedSignDilithium := func() { dilithium.SignDilithium(sk, hash) }
		dilithiumTime := measureExecutionTime(printResults, wrappedSignDilithium, "GenerateSignature - "+dilithium.GetCurrentDilithiumMode(), nIterations)
		dilithiumStats := memory_usage(printResults)
		dilithiumStats["algorithm"] = mode
		dilithiumStats["execution_time"] = dilithiumTime
		result = append(result, dilithiumStats)
		dilithium.ChangeDilithiumMode()
	}

	// Falcon algorithm signature generation
	seed := make([]byte, 64)
	cpu_profiling()
	_, sk, _ := falcon.GenerateKey(seed)
	wrappedSignFalcon := func() { sk.SignCompressed(hash) }
	falconTime := measureExecutionTime(printResults, wrappedSignFalcon, "GenerateSignature - Falcon 1024", nIterations)
	falconStats := memory_usage(printResults)
	falconStats["algorithm"] = "Falcon 1024"
	falconStats["execution_time"] = falconTime
	result = append(result, falconStats)

	sortByAlgorithm(result)
	return result
}

func VerifySignatureTest(hash []byte, nIterations int) []map[string]interface{} {
	var result []map[string]interface{}

	// ECDSA algorithm signature verification
	cpu_profiling()
	key, _ := ecdsa.GenerateKeyECDSA()
	pubkey := ecdsa.FromECDSAPub(&key.PublicKey)
	signature, _ := ecdsa.SignECDSA(hash, key)
	signature = signature[:len(signature)-1] // remove recovery id
	wrappedVerifySignECDSA := func() { ecdsa.VerifySignatureECDSA(pubkey, hash, signature) }
	ecdsaTime := measureExecutionTime(printResults, wrappedVerifySignECDSA, "VerifySignature - ECDSA", nIterations)
	ecdsaStats := memory_usage(printResults)
	ecdsaStats["algorithm"] = "ECDSA"
	ecdsaStats["execution_time"] = ecdsaTime
	result = append(result, ecdsaStats)

	// SPHINCS+ algorithm signature verification
	for i := 0; i < nSphincs; i++ {
		mode := sphincs.GetCurrentSphincsMode()
		cpu_profiling()
		sk, pk := sphincs.GenerateKeySPHINCS()
		signature2 := sphincs.SignSPHINCS(hash, sk)
		wrappedVerifySignSPHINCS := func() { sphincs.VerifySignatureSPHINCS(hash, signature2, pk) }
		sphincsTime := measureExecutionTime(printResults, wrappedVerifySignSPHINCS, "VerifySignature - "+mode, nIterations)
		sphincsStats := memory_usage(printResults)
		sphincsStats["algorithm"] = mode
		sphincsStats["execution_time"] = sphincsTime
		result = append(result, sphincsStats)
		sphincs.ChangeSphincsMode()
	}

	// Dilithium algorithm signature generation
	for i := 0; i < nDilithium; i++ {
		mode := dilithium.GetCurrentDilithiumMode()
		cpu_profiling()
		pk, sk, _ := dilithium.GenerateKeyDilithium()
		signature3 := dilithium.SignDilithium(sk, hash)
		wrappedVerifySignDilithium := func() { dilithium.VerifySignatureDilithium(pk, hash, signature3) }
		dilithiumTime := measureExecutionTime(printResults, wrappedVerifySignDilithium, "GenerateSignature - "+mode, nIterations)
		dilithiumStats := memory_usage(printResults)
		dilithiumStats["algorithm"] = mode
		dilithiumStats["execution_time"] = dilithiumTime
		result = append(result, dilithiumStats)
		dilithium.ChangeDilithiumMode()
	}

	// Falcon algorithm signature verification
	seed := make([]byte, 64)
	cpu_profiling()
	pk, sk, _ := falcon.GenerateKey(seed)
	signatureFalcon, _ := sk.SignCompressed(hash)
	wrappedVerifySignFalcon := func() { pk.Verify(signatureFalcon, hash) }
	falconTime := measureExecutionTime(printResults, wrappedVerifySignFalcon, "VerifySignature - Falcon 1024", nIterations)
	falconStats := memory_usage(printResults)
	falconStats["algorithm"] = "Falcon 1024"
	falconStats["execution_time"] = falconTime
	result = append(result, falconStats)

	sortByAlgorithm(result)
	return result
}

func KeySignatureSizes(hash []byte) []map[string]interface{} {
	var result []map[string]interface{}

	// ECDSA
	key, _ := ecdsa.GenerateKeyECDSA()
	pubkey := ecdsa.FromECDSAPub(&key.PublicKey)
	pubkey = pubkey[:len(pubkey)-1] // remove recovery id
	signature, _ := ecdsa.SignECDSA(hash, key)
	signature = signature[:len(signature)-1] // remove recovery id
	ecdsaStats := getKeySignatureSizes("ECDSA", len(pubkey), len(ecdsa.FromECDSA(key)), len(signature), printResults)
	ecdsaStats["algorithm"] = "ECDSA"
	result = append(result, ecdsaStats)

	//SPHINCS
	for i := 0; i < nSphincs; i++ {
		mode := sphincs.GetCurrentSphincsMode()
		sk, pk := sphincs.GenerateKeySPHINCS()
		skBytes, pkBytes := sphincs.KeysToBytes(sk, pk)
		signature2 := sphincs.SignSPHINCS(hash, sk)
		sigBytes := sphincs.SignatureToBytes(signature2)
		sphincsStats := getKeySignatureSizes(mode, len(pkBytes), len(skBytes), len(sigBytes), printResults)
		sphincsStats["algorithm"] = mode
		result = append(result, sphincsStats)
		sphincs.ChangeSphincsMode()
	}

	//DILITHIUM
	for i := 0; i < nDilithium; i++ {
		mode := dilithium.GetCurrentDilithiumMode()
		pk, sk, _ := dilithium.GenerateKeyDilithium()
		signature3 := dilithium.SignDilithium(sk, hash)
		dilithiumStats := getKeySignatureSizes(mode, len(pk.Bytes()), len(sk.Bytes()), len(signature3), printResults)
		dilithiumStats["algorithm"] = mode
		result = append(result, dilithiumStats)
		dilithium.ChangeDilithiumMode()
	}

	// Falcon 1024
	seed := make([]byte, 64)
	pk, sk, _ := falcon.GenerateKey(seed)
	signatureFalcon, _ := sk.SignCompressed(hash)
	ctSignatureFalcon, _ := signatureFalcon.ConvertToCT()
	falconStats := getKeySignatureSizes("Falcon 1024", len(pk), len(sk), len(signature), printResults)
	falconStats["algorithm"] = "Falcon 1024"
	result = append(result, falconStats)
	falconCTStats := getKeySignatureSizes("Falcon 1024 (CT-format)", len(pk), len(sk), len(ctSignatureFalcon), printResults)
	falconCTStats["algorithm"] = "Falcon 1024 (CT-format)"
	result = append(result, falconCTStats)

	sortByAlgorithm(result)
	return result
}

func getKeySignatureSizes(algorithm string, pkBytes int, skBytes int, signatureBytes int, print bool) map[string]interface{} {
	if print {
		fmt.Printf("\n** %s", algorithm)
		fmt.Printf("\n - Private key --> %d", skBytes)
		fmt.Printf("\n - Public key ---> %d", pkBytes)
		fmt.Printf("\n - Signature ----> %d", signatureBytes)
		fmt.Printf("\n")
	}

	sizes := map[string]interface{}{
		"private_key": skBytes,
		"public_key":  pkBytes,
		"signature":   signatureBytes,
	}
	return sizes
}

func Testing() {
	fmt.Printf("\n========================\n")
}
