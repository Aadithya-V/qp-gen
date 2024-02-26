package main

import (
	"math/rand"
	"time"
)

func Qid() int {
	// Seed the random number generator with current time
	rand.Seed(time.Now().UnixNano())

	// Generate a random number between 0 and 100
	return rand.Intn(10001) // Generates a random number in [0, 101)
}
