package main

import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
)

func main() {
    password := "password123"
    
    // Generate hash
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Password: %s\n", password)
    fmt.Printf("Hash: %s\n", string(hash))
    
    // Test verification
    err = bcrypt.CompareHashAndPassword(hash, []byte(password))
    if err != nil {
        fmt.Printf("Verification FAILED: %v\n", err)
    } else {
        fmt.Printf("Verification SUCCESS\n")
    }
}
