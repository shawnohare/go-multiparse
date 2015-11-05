package multiparse

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseBool(t *testing.T) {
	tests := []struct {
		in  string
		out bool
		err bool
	}{
		{"yes", true, false},
		{"true", true, false},
		{"1", true, false},
		{"no", false, false},
		{"false", false, false},
		{"0", false, false},
		{"abc", false, true},
		{"123", false, true},
		{"123.4", false, true},
		{"$123.4", false, true},
	}

	for _, tt := range tests {
		parser := NewBooleanParser()
		b, err := ParseBool(tt.in)
		b2, err2 := parser.ParseBool(tt.in)
		assert.Equal(t, tt.out, b)
		assert.Equal(t, tt.out, b2)
		if !tt.err {
			assert.NoError(t, err)
			assert.NoError(t, err2)
		} else {
			assert.Error(t, err)
			assert.Error(t, err2)
		}
	}

}
