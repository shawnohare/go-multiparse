package multiparse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMoneyString(t *testing.T) {
	m := Money{original: "$123.40"}
	assert.Equal(t, "$123.40", m.String())
}

func TestMoneyParsedString(t *testing.T) {
	m := Money{original: "$123.40", parsed: "123.40"}
	assert.Equal(t, "123.40", m.ParsedString())
}

func TestMoneyBigFloat(t *testing.T) {
	p := NewStandardMoneyParser()
	m, _ := p.ParseMoney("123")
	_, err := m.BigFloat()
	assert.NoError(t, err)
}

func TestMoneyParserParse(t *testing.T) {
	passes := []string{
		"123",
	}
	fails := []string{
		"abc",
	}

	for _, tt := range passes {
		p := MakeMoneyParser("", "", "")
		_, err := p.Parse(tt)
		assert.NoError(t, err)
	}
	for _, tt := range fails {
		p := MakeMoneyParser("", "", "")
		_, err := p.Parse(tt)
		assert.Error(t, err)
	}
}

func TestMakeMoneyParserCustom(t *testing.T) {
	p := NewStandardMoneyParser()
	c := MakeMoneyParser("^[\\$]", "[,]", "[\\.]")
	d := MakeMoneyParser("$", ",", ".")
	assert.Equal(t, "^[\\$]", p.currencyReStr)
	assert.Equal(t, "^[\\$]", c.currencyReStr)
	assert.Equal(t, "^[\\$]", d.currencyReStr)
	assert.Equal(t, "[,]", p.digitReStr)
	assert.Equal(t, "[,]", c.digitReStr)
	assert.Equal(t, "[,]", d.digitReStr)
	assert.Equal(t, "[\\.]", p.decimalReStr)
	assert.Equal(t, "[\\.]", c.decimalReStr)
	assert.Equal(t, "[\\.]", d.decimalReStr)
}

func TestStandardMoneyParserSanitize(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"$123,456", "123456"},
		{"123,456", "123456"},
	}

	s := NewStandardMoneyParser()
	for _, tt := range tests {
		sanitized, err := s.sanitize(tt.in)
		assert.Equal(t, tt.out, sanitized)
		if tt.out != "" {
			assert.NoError(t, err)
		}
	}
}

func TestMoneyParserSanitize(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"$123,456", "123.456"},
		{"USD123,456", "123.456"},
		{"CURRENCY 123,456", "123.456"},
		{"123.456", "123456"},
	}

	s := MakeMoneyParser("", ".", ",")
	for _, tt := range tests {
		sanitized, err := s.sanitize(tt.in)
		assert.Equal(t, tt.out, sanitized)
		if tt.out != "" {
			assert.NoError(t, err)
		}
	}

}

func TestStandardMoneyParserParse(t *testing.T) {
	passes := []string{
		"$123,234.00",
		"123,234.00",
		"123234.00",
		"12323400",
		"$123,234",
		"$123,234.00",
	}

	fails := []string{
		"USD 1234",
		"123.123.00",
		"123,1.",
		"$",
		"..",
		",,.",
		",.",
		"",
	}

	parser := MakeMoneyParser("$", ",", ".")

	for _, tt := range passes {
		// t.Log(tt)
		_, err := parser.Parse(tt)
		assert.NoError(t, err)
	}
	for _, tt := range fails {
		// t.Log(tt)
		_, err := parser.Parse(tt)
		assert.Error(t, err)
	}
}

func TestParseMonetaryString(t *testing.T) {
	m1 := 123456789.12
	m2 := 123456789.0
	m3 := 123456.0
	m4 := 1234.0
	tests := []struct {
		in  string
		out float64
	}{
		{"$123,456,789.12", m1},
		{"-123,456,789.12", -1 * m1},
		{"+123.456.789,12", m1},
		{"$123.456.789,12", m1},
		{"CURRENCY 123,456,789.12", m1},
		{"CURRENCY 123.456.789,12", m1},
		{"$123,456,789", m2},
		{"$123.456.789", m2},
		{"CURRENCY 123,456,789", m2},
		{"123,456,789", m2},
		{"123.456.789", m2},
		{"$123,456.00", m3},
		{"$123,456", m3},
		{"CURRENCY 123,456", m3},
		{"CURRENCY -123,456", -1 * m3},
		{"123,456", m3},
		{"123456", m3},
		{"123456.00", m3},
		{"123456.", m3},
		{"123456,", m3},
		{"$1,234", m4},
		{"1,234", m4},
		{"1.234", m4},
		{"1234", m4},
		{"1234.567", 1234.567},
		{"CURRENCY 1.234", m4},
		{"1234.5678", 1234.5678},
		{"1", 1.0},
		{"1.00", 1.0},
		{".12", .12},
		{"12.", 12.0},
	}

	errorTests := []string{
		"",
		" ",
		"!",
		"$",
		"USD .",
		"123.234,234,00",
		"123,23.234",
		"12345,23.234",
		"abc",
		"abc++",
		"??#$*(@)",
		".",
		"..",
		"-.",
		"--123",
		"$12345.3D",
		"12abc345.3",
	}

	for _, tt := range tests {
		parser := NewMoneyParser()
		m, err := ParseMonetaryString(tt.in)
		m2, err2 := parser.ParseMoney(tt.in)
		assert.NoError(t, err)
		assert.NoError(t, err2)
		assert.NotNil(t, m)
		assert.NotNil(t, m2)
		// t.Log(tt.in)
		if err == nil {
			ac, _ := m.Float64()
			ac2, _ := m2.Float64()
			assert.Equal(t, tt.out, ac)
			assert.Equal(t, tt.out, ac2)
		}
	}

	for _, tt := range errorTests {
		m, err := ParseMonetaryString(tt)
		assert.Nil(t, m)
		assert.Error(t, err)
	}
}
