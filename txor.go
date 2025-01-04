package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "time"
)

func main() {
    decrypt := flag.Bool("d", false, "Decrypt the input (default is encrypt)")
    keyfile := flag.String("k", "", "Key file for encryption/decryption")
    flag.Parse()

    if *keyfile == "" {
        fmt.Println("Key file is required for encryption/decryption")
        return
    }

    keys, err := readKeysFromFile(*keyfile)
    if err != nil {
        fmt.Println("Error reading keys:", err)
        return
    }

    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        line := scanner.Text()
        var result string
        if *decrypt {
            result = txorDecrypt(line, keys)
        } else {
            result = txorEncrypt(line, keys)
        }
        fmt.Println(result)
    }

    if err := scanner.Err(); err != nil {
        fmt.Println("Error reading from stdin:", err)
    }
}

func readKeysFromFile(filename string) ([]string, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var keys []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        keys = append(keys, scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        return nil, err
    }

    return keys, nil
}

func txorEncrypt(data string, keys []string) string {
    key := selectKey(keys)
    result := make([]byte, len(data))
    for i := range data {
        a := data[i] - '0'
        b := key[i%len(key)] - '0'
        result[i] = ((a + b) % 3) + '0'
    }
    return string(result)
}

func txorDecrypt(data string, keys []string) string {
    key := selectKey(keys)
    result := make([]byte, len(data))
    for i := range data {
        a := data[i] - '0'
        b := key[i%len(key)] - '0'
        result[i] = ((a - b + 3) % 3) + '0' // Add 3 to handle negative values
    }
    return string(result)
}

func selectKey(keys []string) string {
    day := time.Now().UTC().Day() // Change this logic if the key depends on a specific day
    return keys[day%len(keys)]
}
