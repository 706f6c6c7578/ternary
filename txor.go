package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
)

func main() {
    decrypt := flag.Bool("d", false, "Decrypt the input")
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

func txorOperation(a, b byte) byte {
    txorTable := [3][3]byte{
        {0, 1, 2},
        {1, 2, 0},
        {2, 0, 1},
    }
    return txorTable[a][b]
}

func txorEncrypt(data string, keys []string) string {
    key := selectKey(keys)
    result := make([]byte, len(data))
    for i := range data {
        a := data[i] - '0'
        b := key[i%len(key)] - '0'
        result[i] = txorOperation(a, b) + '0'
    }
    return string(result)
}

func txorDecrypt(data string, keys []string) string {
    key := selectKey(keys)
    result := make([]byte, len(data))
    for i := range data {
        a := data[i] - '0'
        b := key[i%len(key)] - '0'
        // Find inverse operation in TXOR table
        for x := byte(0); x < 3; x++ {
            if txorOperation(x, b) == a {
                result[i] = x + '0'
                break
            }
        }
    }
    return string(result)
}

func selectKey(keys []string) string {
    return keys[0]
}

