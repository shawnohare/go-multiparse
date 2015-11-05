package multiparse

// Parsed is the most general type description of a string.
type Parsed struct {
	original  string
	isNumeric bool
	isTime    bool
	isBool    bool
	numeric   *Numeric
	time      *Time
	b         bool
}

// NewParsed returns a Parsed instance with zero values.
func NewParsed() *Parsed {
	return &Parsed{
		numeric: &Numeric{money: new(Money)},
		time:    new(Time),
	}
}

func (p Parsed) String() string {
	return p.original
}

// Numeric instance of the string if it parses as such, or
// the default value if it does not.
func (p Parsed) Numeric() *Numeric {
	if !p.isNumeric {
		return nil
	}
	return p.numeric
}

// Time instance of the string if it parses as such, or
// the default value if it does not.
func (p Parsed) Time() *Time {
	if !p.isTime {
		return nil
	}
	return p.time
}

// Bool instance of the string if it parses as such, or
// the default value if it does not.
func (p Parsed) Money() *Money {
	if !p.IsMoney() {
		return nil
	}
	return p.numeric.Money()
}

// Bool instance of the string if it parses as such, or
// the default value if it does not.
func (p Parsed) Bool() bool {
	if !p.isBool {
		return false
	}
	return true
}

// Type that the parsed string most specifically represents.
func (p Parsed) Type() string {
	if p.isNumeric {
		return p.numeric.Type()
	}

	if p.isTime {
		return p.time.Type()
	}

	return "None"
}

// IsTime reports if the parsed string represents a datetime.
func (p Parsed) IsTime() bool {
	return p.isTime
}

// IsNumeric reports if the parsed string represents a numeric value.
func (p Parsed) IsNumeric() bool {
	return p.isNumeric
}

// IsInt reports if the parsed string represents an integer.
func (p Parsed) IsInt() bool {
	return p.numeric.IsInt()
}

// IsFloat reports if the parsed string represents a floating point number.
func (p Parsed) IsFloat() bool {
	return p.numeric.IsFloat()
}

// IsMoney reports if the parsed string represents a monetary value.
func (p Parsed) IsMoney() bool {
	return p.numeric.IsMoney()
}

// Int reports if the parsed string is an integer and returns the integer.
func (p Parsed) Int() (int, bool) {
	return p.numeric.Int()
}

// Float reports if the parsed string is a float and returns the float.
func (p Parsed) Float() (float64, bool) {
	return p.numeric.Float()
}
