package random

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRandomString(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{
			name: "size = 0",
			size: 0,
		},
		{
			name: "size = 1",
			size: 1,
		},
		{
			name: "size = 5",
			size: 5,
		},
		{
			name: "size = 10",
			size: 10,
		},
		{
			name: "size = 20",
			size: 20,
		},
		{
			name: "size = 30",
			size: 30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str1 := NewRandomString(tt.size)
			str2 := NewRandomString(tt.size)
			// Check the length of the generated strings
			assert.Len(t, str1, tt.size)
			assert.Len(t, str2, tt.size)
			// If size is 0, ensure the string is empty
			if tt.size == 0 {
				assert.Empty(t, str1)
				assert.Empty(t, str2)
			} else {
				// Check that two generated strings are different (a heuristic for randomness)
				assert.NotEqual(t, str1, str2)
			}
		})
	}
}
