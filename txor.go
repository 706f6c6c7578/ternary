package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Function for encryption: (x + y) mod 3
func encrypt(x, y int) int {
	return (x + y) % 3
}

// Function for decryption: (z - y) mod 3
func decrypt(z, y int) int {
	return (z - y + 3) % 3 // +3 to avoid negative results
}

func main() {
	// Check if key filename was provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: txor key.txt < infile > outfile")
		os.Exit(1)
	}

	// Read key file
	keyFile := os.Args[1]
	keyBytes, err := os.ReadFile(keyFile)
	if err != nil {
		fmt.Printf("Error reading key file: %v\n", err)
		os.Exit(1)
	}
	key := strings.TrimSpace(string(keyBytes))

	// Read message from stdin
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	message := scanner.Text()

	// Check if message and key have the same length
	if len(message) != len(key) {
		fmt.Println("Error: Message and key must have the same length")
		os.Exit(1)
	}

	// Process each character
	var result strings.Builder
	for i := 0; i < len(message); i++ {
		// Convert characters to integers
		x, err1 := strconv.Atoi(string(message[i]))
		y, err2 := strconv.Atoi(string(key[i]))

		// Only process valid characters (0, 1, 2)
		if err1 != nil || err2 != nil || x < 0 || x > 2 || y < 0 || y > 2 {
			fmt.Println("Error: Invalid character in message or key")
			os.Exit(1)
		}

		// Encrypt or decrypt (here we decrypt as per the original script)
		decryptedValue := decrypt(x, y)
		result.WriteString(strconv.Itoa(decryptedValue))
	}

	// Output result
	fmt.Println(result.String())
}