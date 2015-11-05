package multiparse

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumericParserParseType(t *testing.T) {
	p := NewNumericParser()
	s := "123"
	n, _ := p.ParseType(s)
	m, _ := p.ParseNumeric(s)
	assert.True(t, reflect.DeepEqual(n, m))
}

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
				money:   &Money{"123", "123"},
			},
		},
		// Int
		{
			"123,456",
			&Numeric{
				parsed:  "123456",
				isInt:   true,
				isFloat: true,
				isMoney: true,
				money:   &Money{"123,456", "123456"},
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
				money:   &Money{"123.4", "123.4"},
			},
		},
		// Float
		{
			"12,345.67",
			&Numeric{
				parsed:  "12345.67",
				isInt:   false,
				isFloat: true,
				isMoney: true,
				money:   &Money{"12,345.67", "12345.67"},
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
				money:   &Money{"$123.45", "123.45"},
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
				money:   &Money{"€123,45", "123.45"},
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
			assert.Equal(t, *(tt.out.money), *(n.money))
			assert.True(t, reflect.DeepEqual(tt.out, n))
		} else {
			assert.Error(t, err)
			assert.Nil(t, nI)
		}
	}

}

func TestParseNumericParseMoney(t *testing.T) {
	var err error
	good := "$123"
	bad := "abc"
	p := NewNumericParser()
	_, err = p.ParseMoney(good)
	assert.NoError(t, err)
	_, err = p.ParseMoney(bad)
	assert.Error(t, err)
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
				money:   &Money{original: "123", parsed: "123"},
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
				money:   &Money{original: "123.4", parsed: "123.4"},
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
				money:   &Money{"$123.45", "123.45"},
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
				money:   &Money{"$123,456", "123456"},
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
			if m.money != nil {
				assert.Equal(t, *m.money, *n.money)
				assert.Equal(t, *tt.out.money, *m.money)
			}
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
				money: &Money{
					original: "123",
					parsed:   "123",
				},
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
				money: &Money{
					original: "123.4",
					parsed:   "123.4",
				},
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
				money: &Money{
					original: "123.45",
					parsed:   "123.45",
				},
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

func TestNumericType(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"123", "int"},
		{"123.5", "float"},
		{"$123.5", "money"},
	}

	p := NewNumericParser()
	for _, tt := range tests {
		x, _ := p.parse(tt.in)
		assert.Equal(t, tt.out, x.Type())
	}

	x := new(Numeric)
	assert.Equal(t, "None", x.Type())
}
