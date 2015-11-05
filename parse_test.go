package multiparse

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseFail(t *testing.T) {
	fails := []string{
		"class 1",
		"$$",
		"blah",
		"abc",
	}

	for _, tt := range fails {
		p, err := ParseType(tt)
		assert.Error(t, err)
		assert.Nil(t, p)
	}
}

func TestParseTypeTimeCase(t *testing.T) {
	p, err := ParseType("2006/1/2")
	assert.NoError(t, err)
	assert.True(t, p.IsTime())
	assert.False(t, p.IsNumeric())
	assert.False(t, p.IsInt())
	assert.False(t, p.IsFloat())
	assert.False(t, p.IsMoney())
}

func TestParseTimeCase(t *testing.T) {
	input := "2015-01-02"
	parsed, err := Parse(input)
	assert.NoError(t, err)

	assert.True(t, parsed.IsTime())
	y := parsed.Time()
	z, _ := time.Parse("2006-01-02", input)
	assert.Equal(t, z, y)
}
