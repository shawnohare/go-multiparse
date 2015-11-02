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
}

type NumericParser struct {
	money Parser
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

func NewNumericParser() *NumericParser {
	mp := NewMoneyParser()
	return MakeNumericParser(mp)
}

func NewStandardNumericParser() *NumericParser {
	mp := NewStandardMoneyParser()
	return MakeNumericParser(mp)
}

func MakeNumericParser(moneyParser Parser) *NumericParser {
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

func (x Numeric) Float() (float64, bool) {
	var y float64
	if x.isFloat {
		y, _ = strconv.ParseFloat(x.parsed, 64)
	}
	return y, x.isFloat
}

func (x Numeric) Money() (*Money, bool) {
	var y *Money
	if x.isMoney {
		// In this case, x.parsed is a standard string and can be
		// converted to a Money instance via the standard parser.
		p := NewStandardMoneyParser()
		y, _ = p.parse(x.parsed)
	}
	return y, x.isMoney
}

func ParseNumeric(s string) (*Numeric, error) {
	// TODO
	return nil, errors.New("bad")
}
