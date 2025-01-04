package main

import (
    "flag"
    "fmt"
    "github.com/google/go-tpm/legacy/tpm2"
    "os"
    "strings"
    "io"
)

func main() {
    length := flag.Int("l", 180, "Length of each key")
    numKeys := flag.Int("n", 1, "Number of keys to generate")
    flag.Parse()

    rwc, err := tpm2.OpenTPM()
    if err != nil {
        fmt.Printf("Failed to open TPM: %v\n", err)
        return
    }
    defer rwc.Close()

    for i := 1; i <= *numKeys; i++ {
        key := generateTernaryKey(rwc, *length)
        filename := fmt.Sprintf("k-%d.txt", i)
        
        err := os.WriteFile(filename, []byte(key), 0644)
        if err != nil {
            fmt.Printf("Error writing file %s: %v\n", filename, err)
            continue
        }
        fmt.Printf("Key %d written to %s\n", i, filename)
    }
}

func generateTernaryKey(rwc io.ReadWriteCloser, length int) string {
    var builder strings.Builder

    for builder.Len() < length {
        random, err := tpm2.GetRandom(rwc, 1)
        if err != nil {
            fmt.Printf("Failed to generate random number: %v\n", err)
            return ""
        }

        // Accept only values that fall within the largest multiple of 3
        if random[0] < 252 { // 252 is the largest multiple of 3 less than 256
            digit := random[0] % 3
            builder.WriteByte('0' + byte(digit))
        }
    }

    return builder.String()
}
