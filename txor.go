package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "strconv"
    "strings"
)

func encrypt(x, y int) int {
    return (x + y) % 3
}

func decrypt(z, y int) int {
    return (z - y + 3) % 3
}

func main() {
    decrypt_mode := flag.Bool("d", false, "decrypt mode")
    flag.Parse()

    args := flag.Args()
    if len(args) < 1 {
        fmt.Println("Usage: txor [-d] key.txt < infile > outfile")
        os.Exit(1)
    }

    // Read and trim key file
    keyBytes, err := os.ReadFile(args[0])
    if err != nil {
        fmt.Printf("Error reading key file: %v\n", err)
        os.Exit(1)
    }
    key := strings.TrimSpace(string(keyBytes))

    // Read and trim message from stdin
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Scan()
    message := strings.TrimSpace(scanner.Text())

    if len(message) != len(key) {
        fmt.Println("Error: Message and key must have the same length")
        os.Exit(1)
    }

    var result strings.Builder
    for i := 0; i < len(message); i++ {
        x, err1 := strconv.Atoi(string(message[i]))
        y, err2 := strconv.Atoi(string(key[i]))

        if err1 != nil || err2 != nil || x < 0 || x > 2 || y < 0 || y > 2 {
            fmt.Println("Error: Invalid character in message or key")
            os.Exit(1)
        }

        var value int
        if *decrypt_mode {
            value = decrypt(x, y)
        } else {
            value = encrypt(x, y)
        }
        result.WriteString(strconv.Itoa(value))
    }

    fmt.Println(result.String())
}
