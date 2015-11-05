package multiparse

import "time"

// Parsed is the most general type description of a string.
type Parsed struct {
	*Numeric
	isNumeric bool
	isTime    bool
	isBool    bool
	time      time.Time
	b         bool
}

// NewParsed returns a Parsed instance with zero values.
func NewParsed() *Parsed {
	return &Parsed{
		Numeric: new(Numeric),
	}
}

// IsTime reports if the parsed string represents a datetime.
func (p Parsed) IsTime() bool {
	return p.isTime
}

// IsBoolreports if the parsed string represents a boolean.
func (p Parsed) IsBool() bool {
	return p.isBool
}

// IsNumeric reports if the parsed string represents a numeric value.
func (p Parsed) IsNumeric() bool {
	return p.isNumeric
}

// Time instance of the string if it parses as such, or
// the default value if it does not.
func (p Parsed) Time() time.Time {
	var t time.Time
	if !p.isTime {
		return t
	}
	return p.time
}

// Bool instance of the string if it parses as such, or
// the default value if it does not.
func (p Parsed) Bool() bool {
	if !p.isBool {
		return false
	}
	return true
}
