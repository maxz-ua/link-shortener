package random

import (
	"math/rand"
	"time"
)

// Global rand instance with a custom source initialized once at the start.
var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

// NewRandomString generates a random string of the given size.
func NewRandomString(size int) string {
	// Characters to use in the random string.
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	// Create a rune slice of the given size.
	b := make([]rune, size)
	// Randomly select characters to form the string.
	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))]
	}
	return string(b)
}
