package multiparse

import (
	"errors"
	"math/big"
	"strconv"
)

// A Money instance is a simple representation of a monetary value.  It
// has access to the original pre-parsed string as well as a few numeric
// types.
type Money struct {
	original string
	parsed   string
}

func (m Money) Type() string       { return "money" }
func (m Money) Value() interface{} { return &m }

// String representation of the monetary value with any original currency
// symbols and formatting included.
func (m Money) String() string {
	return m.original
}

// ParsedString returns a cleaner version of the original monetary string.
func (m Money) ParsedString() string {
	return m.parsed
}

// Float64 representation of the monetary value.  It is not
// recommended that this type be used for accounting.
func (m Money) Float64() (float64, error) {
	f, err := strconv.ParseFloat(m.parsed, 64)
	if err != nil {
		err = errors.New(MoneyFloatError + err.Error())
	}
	return f, err
}

// BigFloat returns a big.Float representation of the monetary value.
// This type is more appropriate for accounting.
func (m Money) BigFloat() (*big.Float, error) {
	tmp := new(big.Float)
	bf, _, err := tmp.Parse(m.parsed, 10)
	if err != nil {
		err = errors.New(MoneyFloatError + err.Error())
	}
	return bf, err
}
