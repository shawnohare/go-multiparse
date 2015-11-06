package multiparse

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumericParserParse(t *testing.T) {
	tests := []struct {
		in  string
		out *Numeric
	}{
		// Int
		{
			"123",
			&Numeric{
				isInt:   true,
				isFloat: true,
				isMoney: false,
				f:       123.0,
			},
		},
		// Int
		{
			"123,456",
			&Numeric{
				isInt:   true,
				isFloat: true,
				isMoney: false,
				f:       123456,
			},
		},
		// Float
		{
			"123.4",
			&Numeric{
				isInt:   false,
				isFloat: true,
				isMoney: false,
				f:       123.4,
			},
		},
		// Float
		{
			"12,345.67",
			&Numeric{
				isInt:   false,
				isFloat: true,
				isMoney: false,
				f:       12345.67,
			},
		},
		// Only money
		{
			"$123.45",
			&Numeric{
				isInt:   false,
				isFloat: true,
				isMoney: true,
				f:       123.45,
			},
		},
		// Another money
		{
			"€123.45",
			&Numeric{
				isInt:   false,
				isFloat: true,
				isMoney: true,
				f:       123.45,
			},
		},
		// Fail case
		{
			"abc",
			nil,
		},
	}

	p := NewNumericParser()
	for _, tt := range tests {
		nI, err := p.Parse(tt.in)
		m, _ := p.ParseNumeric(tt.in)
		if tt.out != nil {
			n := nI.(*Numeric)
			assert.True(t, reflect.DeepEqual(m, n))
			assert.NoError(t, err)
			assert.True(t, reflect.DeepEqual(tt.out, n))
		} else {
			assert.Error(t, err)
			assert.Nil(t, nI)
		}
	}

}

func TestCustomNumericParserParse(t *testing.T) {
	in := "€123,45"
	expected := &Numeric{
		isInt:   false,
		isFloat: true,
		isMoney: true,
		f:       123.45,
	}
	p := NewCustomNumericParser("", "", "")
	actual, err := p.parse(in)
	assert.NoError(t, err)
	assert.Equal(t, *expected, *actual)

}

func TestUSDNumericParserParse(t *testing.T) {
	tests := []struct {
		in  string
		out *Numeric
	}{
		// Int
		{
			"123",
			&Numeric{
				isInt:   true,
				isFloat: true,
				isMoney: false,
				f:       123.0,
			},
		},
		// Float
		{
			"123.4",
			&Numeric{
				isInt:   false,
				isFloat: true,
				isMoney: false,
				f:       123.4,
			},
		},
		// Only money
		{
			"$123.45",
			&Numeric{
				isInt:   false,
				isFloat: true,
				isMoney: true,
				f:       123.45,
			},
		},
		// Another money
		{
			"$123,456",
			&Numeric{
				isInt:   true,
				isFloat: true,
				isMoney: true,
				f:       123456.0,
			},
		},
		// Fail case
		{
			"abc",
			nil,
		},
	}

	p := NewUSDNumericParser()
	for _, tt := range tests {
		nI, err := p.Parse(tt.in)
		m, _ := p.ParseNumeric(tt.in)
		if tt.out != nil {
			n := nI.(*Numeric)
			assert.Equal(t, m.isInt, n.isInt)
			assert.Equal(t, m.isFloat, n.isFloat)
			assert.Equal(t, m.isMoney, n.isMoney)
			assert.Equal(t, tt.out.isInt, m.isInt)
			assert.Equal(t, tt.out.isFloat, m.isFloat)
			assert.Equal(t, tt.out.isMoney, m.isMoney)
			assert.NoError(t, err)
			assert.Equal(t, *tt.out, *n)
		} else {
			assert.Error(t, err)
			assert.Nil(t, nI)
		}
	}

}

func TestNumericMethods(t *testing.T) {
	tests := []struct {
		in       *Numeric
		outFloat float64
	}{
		// Int
		{
			&Numeric{
				isInt:   true,
				isFloat: true,
				isMoney: false,
				f:       123.0,
			},
			123.0,
		},
		// Float
		{
			&Numeric{
				isInt:   false,
				isFloat: true,
				isMoney: false,
				f:       123.4,
			},
			123.4,
		},
	}

	for _, tt := range tests {

		x := tt.in.Int()
		assert.Equal(t, int(tt.outFloat), x)

		y := tt.in.Float()
		assert.Equal(t, tt.outFloat, y)
	}

}

func TestParseNumeric(t *testing.T) {
	tests := []string{
		"123",
		"123.5",
		"$123.4",
		"abc",
	}

	p := NewNumericParser()
	for _, tt := range tests {
		x, err := ParseNumeric(tt)
		y, err2 := p.parse(tt)
		assert.Equal(t, err, err2)
		if err == nil {
			assert.Equal(t, *y, *x)
		}
	}
}
