package multiparse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCurrencySymbolMatching(t *testing.T) {
	symbols := []string{
		"Lek",
		"$",
		"₡",
		"؋",
		"﷼",
		"ƒ",
		"p.",
		"RD$",
		"₩",
		"S/.",
		"p.",
		"Дин.",
		"₪",
		"₹",
		"₱",
		"₺",
	}

	p := NewCustomNumericParser("", "", "")
	// New sure we accurately detect currency symbols.
	for _, s := range symbols {
		// t.Log(s)
		assert.True(t, p.currencyRegex.MatchString(s))
		// New sure we remove the entire currency symbol.
		assert.Equal(t, "", p.removeCurrencySymbol(s))
	}
}

func TestMoneyParserParse(t *testing.T) {
	passes := []string{
		"123",
	}
	fails := []string{
		"abc",
	}

	for _, tt := range passes {
		p := NewNumericParser()
		_, err := p.Parse(tt)
		assert.NoError(t, err)
	}
	for _, tt := range fails {
		p := NewNumericParser()
		_, err := p.Parse(tt)
		assert.Error(t, err)
	}
}

func TestNewNumericParserCustom(t *testing.T) {
	p := NewUSDNumericParser()
	c := NewCustomNumericParser("^[\\$]", "[,]", "[\\.]")
	d := NewCustomNumericParser("$", ",", ".")
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

func TestUSDMoneyParserSanitize(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"$123,456", "123456"},
		{"123,456", "123456"},
	}

	s := NewUSDNumericParser()
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
		{"USD 123,456", "123.456"},
		{"123.456", "123456"},
	}

	s := NewCustomNumericParser("", ".", ",")
	for _, tt := range tests {
		sanitized, err := s.sanitize(tt.in)
		assert.Equal(t, tt.out, sanitized)
		if tt.out != "" {
			assert.NoError(t, err)
		}
	}

}

func TestUSDMoneyParserParse(t *testing.T) {
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

	parser := NewCustomNumericParser("$", ",", ".")

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

func TestParseMoney(t *testing.T) {
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
		{"$123,456,789", m2},
		{"$123.456.789", m2},
		{"123,456,789", m2},
		{"123.456.789", m2},
		{"$123,456.00", m3},
		{"$123,456", m3},
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
		parser := NewCustomNumericParser("", "", "")
		m, err := parser.ParseNumeric(tt.in)
		assert.NoError(t, err)
		assert.NotNil(t, m)
		// t.Log(tt.in)
		if err == nil {
			ac := m.Float()
			assert.NotEqual(t, 0.0, ac)
			assert.Equal(t, tt.out, ac)
		}
	}

	for _, tt := range errorTests {
		parser := NewCustomNumericParser("", "", "")
		m, err := parser.ParseNumeric(tt)
		assert.Nil(t, m)
		assert.Error(t, err)
	}
}
