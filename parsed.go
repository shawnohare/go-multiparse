package multiparse

// Parsed is the most general type description of a string.
type Parsed struct {
	original  string
	isNumeric bool
	isTime    bool
	numeric   *Numeric
	time      *Time
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

// Numeric reports whether the parsed string is a Numeric type and
// returns its value.
func (p Parsed) Numeric() (*Numeric, bool) {
	if !p.isNumeric {
		return nil, false
	}

	return p.numeric, true
}

// Time reports whether the parsed string represents a datetime and the
// the *Time that has been parsed.
func (p Parsed) Time() (*Time, bool) {
	if !p.isTime {
		return nil, false
	}

	return p.time, true
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

// Money reports if the parsed string is a monetary value and returns the
// corresponding Money instance.
func (p Parsed) Money() (*Money, bool) {
	return p.numeric.Money()
}
