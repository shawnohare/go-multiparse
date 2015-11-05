package multiparse

import (
	"fmt"
	"testing"

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

type BadBoolParser struct{}

func (b BadBoolParser) Parse(s string) (interface{}, error) {
	return "not a bool", nil
}

func TestNewUSDParser(t *testing.T) {
	parser := NewUSDParser()
	parsed, err := parser.ParseType("$123,456")
	assert.NoError(t, err)
	assert.Equal(t, 123456.0, parsed.Float())
}

func TestParseIntCase(t *testing.T) {

	parsed, err := Parse("123")
	assert.NoError(t, err)
	assert.True(t, parsed.IsInt())
	assert.True(t, parsed.IsNumeric())
	f := parsed.Int()
	f2 := parsed.Int()
	assert.Equal(t, 123, f)
	assert.Equal(t, 123, f2)
}

func TestParseFloatCase(t *testing.T) {
	parsed, err := Parse("1234.5")
	assert.NoError(t, err)
	assert.True(t, parsed.IsFloat())
	assert.True(t, parsed.IsNumeric())
	f := parsed.Float()
	assert.Equal(t, 1234.5, f)
}

func TestParserMoneyCase(t *testing.T) {
	input := "$1234.5"
	parsed, err := Parse(input)
	assert.NoError(t, err)
	assert.True(t, parsed.IsMoney())
	assert.True(t, parsed.IsNumeric())
	assert.Equal(t, 1234.5, parsed.Float())
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
	b3 := new(BadBoolParser)
	parser := NewCustomParser(b1, b2, b3)
	_, err := parser.Parse("123")
	assert.Error(t, err)
	_, err = parser.ParseType("123")
}

func TestParserParseType(t *testing.T) {
	var err error
	p := NewParser()

	// Pass
	passes := []string{"123"}
	for _, tt := range passes {
		_, err = p.ParseType(tt)
		assert.NoError(t, err)
		_, err = ParseType(tt)
		assert.NoError(t, err)
	}

	// Fail
	failures := []string{
		"abc",
	}

	for _, f := range failures {
		_, err = p.ParseType(f)
		assert.Error(t, err)
		_, err = ParseType(f)
		assert.Error(t, err)
	}

}

func TestParserParseInt(t *testing.T) {
	var err error
	p := NewParser()

	// Pass
	_, err = p.ParseInt("123")
	assert.NoError(t, err)
	_, err = ParseInt("123")
	assert.NoError(t, err)

	// Fail
	failures := []string{
		"abc",
		"123.4",
		"$123.4",
	}

	for _, f := range failures {
		_, err = p.ParseInt(f)
		assert.Error(t, err)
		_, err = ParseInt(f)
		assert.Error(t, err)
	}
}

func TestParserParseFloat(t *testing.T) {
	var err error
	p := NewParser()

	// Pass
	passes := []string{
		"123",
		"123.4",
		"123,456",
		"$123.4",
	}

	for _, tt := range passes {
		_, err = p.ParseFloat(tt)
		assert.NoError(t, err)
		_, err = ParseFloat(tt)
		assert.NoError(t, err)
	}
	_, err = p.ParseFloat("123")
	_, err = ParseFloat("123")
	assert.NoError(t, err)

	// Fail
	failures := []string{
		"abc",
		"",
		"..",
		"3i",
	}

	for _, f := range failures {
		_, err = p.ParseFloat(f)
		assert.Error(t, err)
		_, err = ParseFloat(f)
		assert.Error(t, err)
	}
}

func TestParserParseTime(t *testing.T) {
	var err error
	p := NewParser()

	// Pass
	_, err = p.ParseTime("2015-06-15")
	assert.NoError(t, err)
	_, err = ParseTime("2015-06-15")
	assert.NoError(t, err)

	// Fail
	failures := []string{
		"123",
		"abc",
		"$123.4",
	}

	for _, f := range failures {
		_, err = p.ParseTime(f)
		assert.Error(t, err)
		_, err = ParseTime(f)
		assert.Error(t, err)
	}
}

func ExampleParseType() {
	var p *Parsed

	p, _ = Parse("$12,345")
	fmt.Println(p.Float())
	// output: 12345
}
