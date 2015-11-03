package multiparse

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParsedString(t *testing.T) {
	p := &Parsed{
		original: "1234.4",
	}

	assert.Equal(t, "1234.4", p.String())
}

func TestNewParsed(t *testing.T) {
	parsed := NewParsed()
	_, ok := parsed.Numeric()
	assert.False(t, ok)
	_, ok = parsed.Time()
	assert.False(t, ok)
	assert.Equal(t, "None", parsed.Type())
}

func TestParsedType(t *testing.T) {
	p := &Parsed{
		original:  "1234.4",
		isNumeric: true,
		numeric:   &Numeric{isFloat: true},
	}
	assert.Equal(t, "float", p.Type())

	q := &Parsed{
		original: "2015-06-01",
		isTime:   true,
		time:     &Time{},
	}
	assert.Equal(t, "time", q.Type())

}
