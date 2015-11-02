package multiparse

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
				parsed:  "123",
				isInt:   true,
				isFloat: true,
				isMoney: true,
			},
		},
		// Float
		{
			"123.4",
			&Numeric{
				parsed:  "123.4",
				isInt:   false,
				isFloat: true,
				isMoney: true,
			},
		},
		// Only money
		{
			"$123.45",
			&Numeric{
				parsed:  "123.45",
				isInt:   false,
				isFloat: false,
				isMoney: true,
			},
		},
		// Another money
		{
			"€123,45",
			&Numeric{
				parsed:  "123.45",
				isInt:   false,
				isFloat: false,
				isMoney: true,
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
			assert.Equal(t, *m, *n)
			assert.NoError(t, err)
			assert.Equal(t, *tt.out, *n)
		} else {
			assert.Error(t, err)
			assert.Nil(t, nI)
		}
	}

}

func TestStandardNumericParserParse(t *testing.T) {
	tests := []struct {
		in  string
		out *Numeric
	}{
		// Int
		{
			"123",
			&Numeric{
				parsed:  "123",
				isInt:   true,
				isFloat: true,
				isMoney: true,
			},
		},
		// Float
		{
			"123.4",
			&Numeric{
				parsed:  "123.4",
				isInt:   false,
				isFloat: true,
				isMoney: true,
			},
		},
		// Only money
		{
			"$123.45",
			&Numeric{
				parsed:  "123.45",
				isInt:   false,
				isFloat: false,
				isMoney: true,
			},
		},
		// Another money
		{
			"$123,456",
			&Numeric{
				parsed:  "123456",
				isInt:   false,
				isFloat: false,
				isMoney: true,
			},
		},
		// Fail case
		{
			"abc",
			nil,
		},
	}

	p := NewStandardNumericParser()
	for _, tt := range tests {
		nI, err := p.Parse(tt.in)
		m, _ := p.ParseNumeric(tt.in)
		if tt.out != nil {
			n := nI.(*Numeric)
			assert.Equal(t, *m, *n)
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
		outInt   int
		outFloat float64
		outMoney *Money
	}{
		// Int
		{
			&Numeric{
				parsed:  "123",
				isInt:   true,
				isFloat: true,
				isMoney: true,
			},
			123,
			123.0,
			&Money{
				original: "123",
				parsed:   "123",
			},
		},
		// Float
		{
			&Numeric{
				parsed:  "123.4",
				isInt:   false,
				isFloat: true,
				isMoney: true,
			},
			0,
			123.4,
			&Money{
				original: "123.4",
				parsed:   "123.4",
			},
		},
		// Only money
		{
			&Numeric{
				parsed:  "123.45",
				isInt:   false,
				isFloat: false,
				isMoney: true,
			},
			0,
			0.0,
			&Money{
				original: "123.45",
				parsed:   "123.45",
			},
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.in.parsed, tt.in.String())

		x, in := tt.in.Int()
		assert.Equal(t, tt.in.isInt, in)
		assert.Equal(t, tt.outInt, x)

		y, in := tt.in.Float()
		assert.Equal(t, tt.in.isFloat, in)
		assert.Equal(t, tt.outFloat, y)

		z, in := tt.in.Money()
		assert.Equal(t, tt.in.isMoney, in)
		assert.Equal(t, *tt.outMoney, *z)
	}

}
