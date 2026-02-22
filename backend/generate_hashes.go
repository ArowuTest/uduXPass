package main

import (
"fmt"
"golang.org/x/crypto/bcrypt"
)

func main() {
passwords := map[string]string{
"Admin@123456": "admin",
"Scanner@123": "scanner",
"User@123": "user",
}

for password, label := range passwords {
hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
if err != nil {
fmt.Printf("Error hashing %s: %v\n", label, err)
continue
}
fmt.Printf("%s (%s): %s\n", label, password, string(hash))
}
}
