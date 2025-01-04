package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"strings"
)

// Encode binary data to ternary numbers as string with line length of 64
func encodeBinaryToTernaryNumbers(data []byte) string {
	binaryInt := new(big.Int).SetBytes(data)
	ternaryStr := binaryInt.Text(3)

	// Format the output to have lines of 64 characters
	var result strings.Builder
	for i, c := range ternaryStr {
		result.WriteRune(c)
		if (i+1)%64 == 0 {
			result.WriteRune('\n')
		}
	}
	// Add a final newline for cleaner output
	result.WriteRune('\n')
	return result.String()
}

// Decode ternary numbers string back to binary data
func decodeTernaryNumbersToBinary(data string) ([]byte, error) {
	// Remove newlines for decoding
	data = strings.ReplaceAll(data, "\n", "")

	// Validate input (ensure only 0, 1, 2 are present)
	for _, c := range data {
		if c < '0' || c > '2' {
			return nil, fmt.Errorf("invalid character '%c' in input", c)
		}
	}

	ternaryInt, ok := new(big.Int).SetString(data, 3)
	if !ok {
		return nil, fmt.Errorf("invalid ternary data")
	}

	return ternaryInt.Bytes(), nil
}

func main() {
	decodeFlag := flag.Bool("d", false, "decode mode")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [-d]\n", os.Args[0])
		fmt.Fprintln(flag.CommandLine.Output(), `
Description:
  ternary encodes binary data to a ternary (base-3) number system or decodes it back.
  By default, it reads from stdin and writes to stdout.

Options:
  -d    Decode mode: convert ternary back to binary.
`)
	}
	flag.Parse()

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Error reading input: %v\n", err)
	}

	if len(input) == 0 {
		fmt.Fprintln(os.Stderr, "Error: No input provided.")
		flag.Usage()
		os.Exit(1)
	}

	if *decodeFlag {
		result, err := decodeTernaryNumbersToBinary(string(input))
		if err != nil {
			log.Fatalf("Decoding error: %v\n", err)
		}
		_, writeErr := os.Stdout.Write(result)
		if writeErr != nil {
			log.Fatalf("Error writing output: %v\n", writeErr)
		}
	} else {
		result := encodeBinaryToTernaryNumbers(input)
		_, writeErr := os.Stdout.Write([]byte(result))
		if writeErr != nil {
			log.Fatalf("Error writing output: %v\n", writeErr)
		}
	}
}
