package main

import (
	"crypto/test/scripts"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

func main() {

	// Start CPU profiling
	cpu_profiling()
	// Perform ECDSA tests
	ecdsa_test()
	// Print memory usage after test
	memory_usage()

	// Start CPU profiling
	cpu_profiling()
	// Perform SPHINCS+ tests
	sphincs_test()
	// Print memory usage after test
	memory_usage()
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
	fmt.Printf("\n - Alloc -------> %v KB\n", m.Alloc/1024)
	fmt.Printf(" - TotalAlloc --> %v KB\n", m.TotalAlloc/1024)
	fmt.Printf(" - Sys ---------> %v KB\n", m.Sys/1024)
	fmt.Printf(" - NumGC -------> %v\n", m.NumGC)
}

// ECDSA Cryptography testing
func ecdsa_test() {
	fmt.Printf("\n---------------- GETH ECDSA SIGNATURE TESTING ----------------\n\n")

	// Measure execution time for each function
	start := time.Now()
	scripts.TestEcrecover()
	fmt.Println(" - TestEcrecover execution time -----------------> ", time.Since(start))

	start = time.Now()
	scripts.TestVerifySignature()
	fmt.Println(" - TestVerifySignature execution time -----------> ", time.Since(start))

	start = time.Now()
	scripts.TestVerifySignatureMalleable()
	fmt.Println(" - TestVerifySignatureMalleable execution time --> ", time.Since(start))

	start = time.Now()
	scripts.TestDecompressPubkey()
	fmt.Println(" - TestDecompressPubkey execution time ----------> ", time.Since(start))

	start = time.Now()
	scripts.TestCompressPubkey()
	fmt.Println(" - TestCompressPubkey execution time ------------> ", time.Since(start))

	start = time.Now()
	scripts.TestPubkeyRandom()
	fmt.Println(" - TestPubkeyRandom execution time --------------> ", time.Since(start))

	// Run benchmarks
	scripts.BenchmarkEcrecoverSignature()
	scripts.BenchmarkVerifySignature()
	scripts.BenchmarkDecompressPubkey()
}

// SPHINCS+ Cryptography testing
func sphincs_test() {
	fmt.Printf("\n----------------- SPHINCS+ SIGNATURE TESTING -----------------\n\n")

	// Measure execution time for each function
	start := time.Now()
	scripts.TestSphincs()
	fmt.Println(" - TestSphincs execution time -----------------> ", time.Since(start))
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
