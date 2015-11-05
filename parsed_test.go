package multiparse

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewParsed(t *testing.T) {
	p := NewParsed()
	assert.False(t, p.IsBool())
	assert.False(t, p.IsNumeric())
	assert.False(t, p.IsTime())
	assert.False(t, p.IsInt())
	assert.False(t, p.IsFloat())

	assert.Equal(t, p.isBool, p.IsBool())
	assert.Equal(t, p.isNumeric, p.IsNumeric())
	assert.Equal(t, p.isTime, p.IsTime())
}

func TestParsedTime(t *testing.T) {
	var tt time.Time
	p := NewParsed()
	assert.Equal(t, tt, p.Time())
}

func TestParsedBool(t *testing.T) {
	b, _ := ParseType("yes")
	b2, _ := ParseType("123")
	assert.True(t, b.Bool())
	assert.False(t, b2.Bool())
}
