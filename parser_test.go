package multiparse

import (
	"fmt"
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

func TestMakeUSDParser(t *testing.T) {
	parser := MakeUSDParser()
	parsed, _ := parser.ParseType("$123,456")
	m, _ := parsed.Money()
	f, err := m.Float64()
	assert.NoError(t, err)
	assert.Equal(t, 123456.0, f)
}

func TestParseIntCase(t *testing.T) {

	parsed, err := Parse("123")
	assert.NoError(t, err)
	assert.Equal(t, "123", parsed.String())

	assert.True(t, parsed.IsInt())
	assert.True(t, parsed.IsNumeric())
	y, _ := parsed.Numeric()
	f, _ := y.Int()
	f2, _ := parsed.Int()
	assert.Equal(t, 123, f)
	assert.Equal(t, 123, f2)
}

func TestParseFloatCase(t *testing.T) {

	parsed, err := Parse("1234.5")
	assert.NoError(t, err)
	assert.Equal(t, "1234.5", parsed.String())

	assert.True(t, parsed.IsFloat())
	assert.True(t, parsed.IsNumeric())
	y, _ := parsed.Numeric()
	f, _ := y.Float()
	f2, _ := parsed.Float()
	assert.Equal(t, 1234.5, f)
	assert.Equal(t, 1234.5, f2)
}

func TestParserMoneyCase(t *testing.T) {

	input := "$1234.5"
	parsed, err := Parse(input)
	assert.NoError(t, err)
	assert.Equal(t, "$1234.5", parsed.String())

	assert.True(t, parsed.IsMoney())
	assert.True(t, parsed.IsNumeric())
	y, _ := parsed.Numeric()
	f, _ := y.Money()
	f2, _ := ParseMoney(input)
	f3, _ := parsed.Money()
	assert.Equal(t, *f2, *f)
	assert.Equal(t, *f2, *f3)
}

func TestParseTimeCase(t *testing.T) {
	input := "2015-01-02"
	parsed, err := Parse(input)
	assert.NoError(t, err)
	assert.Equal(t, input, parsed.String())

	assert.True(t, parsed.IsTime())
	y, _ := parsed.Time()
	z, _ := time.Parse("2006-01-02", input)
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

func ExampleParseType() {
	var p *Parsed

	p, _ = Parse("$12,345")
	fmt.Println(p.Type())
	// output: money
}
