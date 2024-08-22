package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/frontend"

	"github.com/PolyhedraZK/ExpanderCompilerCollection"
	"github.com/PolyhedraZK/ExpanderCompilerCollection/test"
)

type Circuit struct {
	X      frontend.Variable `gnark:",public"`
	Y      frontend.Variable `gnark:",public"`
	Result frontend.Variable `gnark:",public"`
}

// / expander compiler optimized this circuit, it is not the same as my fibonacci
// / we need to write eache intermediate value as frontend variable if we want to genetate the same circuit
// as expander compiler reuses the same gate to represent different rounds, but modify the coef of the gate
func (circuit *Circuit) Define(api frontend.API) error {
	var rounds int = 2

	var FibonacciResult frontend.Variable = api.Add(circuit.X, circuit.Y)

	for i := 0; i < rounds; i++ {
		fmt.Println("round")
		circuit.X = circuit.Y
		circuit.Y = FibonacciResult
		FibonacciResult = api.Add(circuit.Y, circuit.X)
	}
	api.AssertIsEqual(circuit.Result, FibonacciResult)
	return nil
}

func main() {
	assignment := &Circuit{X: 1, Y: 1, Result: 5}

	circuit, _ := ExpanderCompilerCollection.Compile(ecc.BN254.ScalarField(), &Circuit{})
	//fmt.Println("circuit is ", circuit)
	c := circuit.GetLayeredCircuit()
	//fmt.Println("layered circuit is ", c)

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Fatalf("Failed to serialize data: %v", err)
	}

	err = os.WriteFile("circuitjson.txt", data, 0o644)
	if err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}

	os.WriteFile("circuit.txt", c.Serialize(), 0o644)
	inputSolver := circuit.GetInputSolver()
	witness, _ := inputSolver.SolveInputAuto(assignment)

	witnessdata, err := json.MarshalIndent(witness, "", "  ")
	if err != nil {
		log.Fatalf("Failed to serialize data: %v", err)
	}

	err = os.WriteFile("witnessjson.txt", witnessdata, 0o644)
	if err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}

	os.WriteFile("witness.txt", witness.Serialize(), 0o644)

	if !test.CheckCircuit(c, witness) {
		panic("verification failed")
	}

}
