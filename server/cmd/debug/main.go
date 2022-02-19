package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
)

func main() {
	fmt.Println(uuid.New().String())
	fmt.Println(time.Now().UnixNano())
	// rand
	fmt.Println(rand.Int(rand.Reader, big.NewInt(1000)))
}


