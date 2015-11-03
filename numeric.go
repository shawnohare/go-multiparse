package multiparse

import (
	"errors"
	"strconv"
)

// Numeric instances are containers for the various valid numerical types
// that a string may be parsed into.
type Numeric struct {
	parsed  string
	isInt   bool
	isFloat bool
	isMoney bool
	money   *Money
}

type NumericParser struct {
	money MultiParse
}

func (x Numeric) Type() string {
	if x.isInt {
		return "int"
	} else if x.isFloat {
		return "float"
	} else if x.isMoney {
		return "money"
	} else {
		return "None"
	}
}

func MakeGeneralNumericParser() *NumericParser {
	mp := MakeGeneralMoneyParser()
	return MakeNumericParser(mp)
}

func MakeStandardNumericParser() *NumericParser {
	mp := MakeStandardMoneyParser()
	return MakeNumericParser(mp)
}

func MakeNumericParser(moneyParser MultiParse) *NumericParser {
	return &NumericParser{
		money: moneyParser,
	}
}

func (p NumericParser) Parse(s string) (interface{}, error) {
	return p.parse(s)
}

func (p NumericParser) ParseNumeric(s string) (*Numeric, error) {
	return p.parse(s)
}

func (p NumericParser) parse(s string) (*Numeric, error) {
	var n *Numeric
	var err error

	_, err = strconv.Atoi(s)
	if err == nil {
		n = &Numeric{
			parsed:  s,
			isInt:   true,
			isFloat: true,
			isMoney: true,
		}
		return n, nil
	}

	_, err = strconv.ParseFloat(s, 64)
	if err == nil {
		n = &Numeric{
			parsed:  s,
			isFloat: true,
			isMoney: true,
		}
		return n, nil
	}

	mI, err := p.money.Parse(s)
	if err == nil {
		m := mI.(*Money)
		n = &Numeric{
			parsed:  m.ParsedString(),
			isMoney: true,
			money:   m,
		}
		return n, nil
	}

	return nil, errors.New(ParseNumericError)
}

func (x Numeric) String() string {
	return x.parsed
}

// Int reports whether the Numeric instance can be an integer
// and returns its value.
func (x Numeric) Int() (int, bool) {
	var y int
	if x.isInt {
		y, _ = strconv.Atoi(x.parsed)
	}
	return y, x.isInt
}

// Float reports whether the Numeric instance can be a float
// and returns its value.
func (x Numeric) Float() (float64, bool) {
	var y float64
	if x.isFloat {
		y, _ = strconv.ParseFloat(x.parsed, 64)
	}
	return y, x.isFloat
}

func (x Numeric) Money() (*Money, bool) {
	var y *Money
	if x.isMoney && x.money != nil && x.money.original != "" {
		y = x.money
	} else {
		p := MakeStandardMoneyParser()
		y, _ = p.parse(x.parsed)
	}
	return y, x.isMoney
}

func ParseNumeric(s string) (*Numeric, error) {
	p := MakeGeneralNumericParser()
	return p.parse(s)
}
