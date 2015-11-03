package multiparse

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type BadNumericParser struct{}

func (b BadNumericParser) Parse(s string) (interface{}, error) {
	return "not a number", nil
}

type BadTimeParser struct{}

func (b BadTimeParser) Parse(s string) (interface{}, error) {
	return "not a time", nil
}

func TestParseFloatCase(t *testing.T) {

	parsed, err := Parse("1234.5")
	assert.NoError(t, err)
	assert.Equal(t, "1234.5", parsed.String())

	y, ok := parsed.Numeric()
	f, _ := y.Float()
	assert.True(t, ok)
	assert.Equal(t, 1234.5, f)
}

func TestParserMoneyCase(t *testing.T) {

	input := "$1234.5"
	parsed, err := Parse(input)
	assert.NoError(t, err)
	assert.Equal(t, "$1234.5", parsed.String())

	y, ok := parsed.Numeric()
	f, _ := y.Money()
	f2, _ := ParseMoney(input)
	assert.True(t, ok)
	assert.Equal(t, *f2, *f)
}

func TestParseTimeCase(t *testing.T) {
	input := "2015-01-02"
	parsed, err := Parse(input)
	assert.NoError(t, err)
	assert.Equal(t, input, parsed.String())

	y, ok := parsed.Time()
	z, _ := time.Parse("2006-01-02", input)
	assert.True(t, ok)
	assert.Equal(t, z, y.Time())
}

func TestParseInvalidCase(t *testing.T) {
	failures := []string{
		"abc",
		"",
		"$",
		"123abc840",
	}

	for _, tt := range failures {
		_, err := Parse(tt)
		assert.Error(t, err)
	}
}

func TestBadParsers(t *testing.T) {
	b1 := new(BadNumericParser)
	b2 := new(BadTimeParser)
	parser := MakeParser(b1, b2)
	_, err := parser.Parse("123")
	assert.Error(t, err)
	_, err = parser.ParseType("123")
}
