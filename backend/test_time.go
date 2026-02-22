package main

import (
"fmt"
"time"
)

func main() {
now := time.Now()
fmt.Printf("Now (local): %v\n", now)
fmt.Printf("Now (UTC): %v\n", now.UTC())

expires1 := now.Add(15 * time.Minute).UTC()
fmt.Printf("\nMethod 1 - Add then UTC: %v\n", expires1)

expires2 := now.UTC().Add(15 * time.Minute)
fmt.Printf("Method 2 - UTC then Add: %v\n", expires2)

fmt.Printf("\nAre they equal? %v\n", expires1.Equal(expires2))
}
