package snark

import (
	"os"
	"bytes"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"
	"github.com/arnaucube/go-snark/circuitcompiler"
	"github.com/arnaucube/go-snark/groth16"
	"github.com/arnaucube/go-snark/r1csqap"
	"github.com/stretchr/testify/assert"
	"runtime/pprof"
	"runtime"
	
)

// Measure CPU utilization
func measureCPUUsage() float64 {
	var cpuUsage float64
	numCPU := runtime.NumCPU()
	startTime := time.Now()
	done := make(chan struct{})
	go func() {
		start := runtime.NumGoroutine()
		time.Sleep(1 * time.Second) // Measure over 1 sec
		end := runtime.NumGoroutine()
		elapsed := time.Since(startTime).Seconds()
		cpuUsage = float64(end-start) / elapsed * 100 / float64(numCPU)
		close(done)
	}()

	<-done
	return cpuUsage
}

// Measure memory usage
func measureMemoryUsage() (uint64, uint64) {
    var memStats runtime.MemStats
    runtime.ReadMemStats(&memStats)
    return memStats.Alloc / 1024, memStats.Sys / 1024 // Convert to KB
}

func TestGroth16MinimalFlow(t *testing.T) {
	fmt.Println("testing Groth16 minimal flow")
	// circuit function
	// y = x^3 + x + 5
	code := `
	func main(private s0, public s1):
		s2 = s0 * s0
		s3 = s2 * s0
		s4 = s3 + s0
		s5 = s4 + 5
		equals(s1, s5)
		out = 1 * 1
	`
	fmt.Print("\ncode of the circuit:")

	// parse the code
	parser := circuitcompiler.NewParser(strings.NewReader(code))
	circuit, err := parser.Parse()
	assert.Nil(t, err)

	b3 := big.NewInt(int64(3))
	privateInputs := []*big.Int{b3}
	b35 := big.NewInt(int64(35))
	publicSignals := []*big.Int{b35}

	// wittness
	w, err := circuit.CalculateWitness(privateInputs, publicSignals)
	assert.Nil(t, err)

	// code to R1CS
	fmt.Println("\ngenerating R1CS from code")
	a, b, c := circuit.GenerateR1CS()
	fmt.Println("\nR1CS:")
	fmt.Println("a:", a)
	fmt.Println("b:", b)
	fmt.Println("c:", c)

	// R1CS to QAP
	// TODO zxQAP is not used and is an old impl, TODO remove
	alphas, betas, gammas, _ := Utils.PF.R1CSToQAP(a, b, c)
	fmt.Println("qap")
	assert.Equal(t, 8, len(alphas))
	assert.Equal(t, 8, len(alphas))
	assert.Equal(t, 8, len(alphas))
	assert.True(t, !bytes.Equal(alphas[1][1].Bytes(), big.NewInt(int64(0)).Bytes()))

	ax, bx, cx, px := Utils.PF.CombinePolynomials(w, alphas, betas, gammas)
	assert.Equal(t, 7, len(ax))
	assert.Equal(t, 7, len(bx))
	assert.Equal(t, 7, len(cx))
	assert.Equal(t, 13, len(px))

	// ---
	// from here is the GROTH16
	// ---
	// calculate trusted setup
	fmt.Println("groth")
	setup, err := groth16.GenerateTrustedSetup(len(w), *circuit, alphas, betas, gammas)
	assert.Nil(t, err)
	fmt.Println("\nt:", setup.Toxic.T)

	hx := Utils.PF.DivisorPolynomial(px, setup.Pk.Z)
	div, rem := Utils.PF.Div(px, setup.Pk.Z)
	assert.Equal(t, hx, div)
	assert.Equal(t, rem, r1csqap.ArrayOfBigZeros(6))

	// hx==px/zx so px==hx*zx
	assert.Equal(t, px, Utils.PF.Mul(hx, setup.Pk.Z))

	// check length of polynomials H(x) and Z(x)
	assert.Equal(t, len(hx), len(px)-len(setup.Pk.Z)+1)

	start := time.Now()

	// Create a file to store CPU profiling data
	cpuProfile, err := os.Create("cpu_profile.prof")
	if err != nil {
		t.Fatal("Could not create CPU profile:", err)
	}
	defer cpuProfile.Close()

	// Start CPU profiling
	pprof.StartCPUProfile(cpuProfile)
	defer pprof.StopCPUProfile()

	// Measure proof generation time
	startTime := time.Now()

	memBefore, sysBefore := measureMemoryUsage()

	proof, err := groth16.GenerateProofs(*circuit, setup.Pk, w, px)
	if err != nil {
		t.Fatalf("Error generating proof: %v", err)
	}

	var rawProofBuffer bytes.Buffer

	// Write A, B, C directly in binary without using any encoder
	rawProofBuffer.Write(proof.PiA[0].Bytes())
	rawProofBuffer.Write(proof.PiA[1].Bytes())
	rawProofBuffer.Write(proof.PiA[2].Bytes())

	rawProofBuffer.Write(proof.PiB[0][0].Bytes())
	rawProofBuffer.Write(proof.PiB[0][1].Bytes())
	rawProofBuffer.Write(proof.PiB[1][0].Bytes())
	rawProofBuffer.Write(proof.PiB[1][1].Bytes())
	rawProofBuffer.Write(proof.PiB[2][0].Bytes())
	rawProofBuffer.Write(proof.PiB[2][1].Bytes())

	rawProofBuffer.Write(proof.PiC[0].Bytes())
	rawProofBuffer.Write(proof.PiC[1].Bytes())
	rawProofBuffer.Write(proof.PiC[2].Bytes())

	fmt.Printf("zk-SNARK Proof Size (Minimal Encoding): %d bytes\n", rawProofBuffer.Len())

	// Stop timer after proof generation
	elapsedTime := time.Since(startTime)

	// Print profiling results
	fmt.Printf("=== zk-SNARK CPU Profiling Done ===\n")
	fmt.Printf("Proof Generation Time: %v\n", elapsedTime)

	// Log results in Go test output
	t.Logf("Proof generated successfully in %v", elapsedTime)
	fmt.Println("CPU profile saved to cpu_profile.prof")

	elapsed := time.Since(start)

	// Measure memory usage after proof generation
	memAfter, sysAfter := measureMemoryUsage()

	// Print memory results
	fmt.Printf("Memory Usage Before: %d KB, After: %d KB\n", memBefore, memAfter)
	fmt.Printf("Total System Memory Before: %d KB, After: %d KB\n", sysBefore, sysAfter)

	// Log results
	t.Logf("Memory Usage Before: %d KB, After: %d KB", memBefore, memAfter)
	t.Logf("Total System Memory Before: %d KB, After: %d KB", sysBefore, sysAfter)

	fmt.Printf("Time taken to generate proof: %v\n", elapsed)

	fmt.Println("\n proofs:")
	fmt.Println(proof)

	// fmt.Println("public signals:", proof.PublicSignals)
	fmt.Println("\nsignals:", circuit.Signals)
	fmt.Println("witness:", w)
	b35Verif := big.NewInt(int64(35))
	publicSignalsVerif := []*big.Int{b35Verif}
	before := time.Now()
	assert.True(t, groth16.VerifyProof(setup.Vk, proof, publicSignalsVerif, true))
	fmt.Println("verify proof time elapsed:", time.Since(before))

	// check that with another public input the verification returns false
	bOtherWrongPublic := big.NewInt(int64(34))
	wrongPublicSignalsVerif := []*big.Int{bOtherWrongPublic}
	assert.True(t, !groth16.VerifyProof(setup.Vk, proof, wrongPublicSignalsVerif, false))
}
