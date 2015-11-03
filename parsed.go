package multiparse

// Type reporting for types with multiple subtypes.
type Type interface {
	Type() string
	Value() interface{}
}

// Parsed is the most general type description of a string.
type Parsed struct {
	original  string
	isNumeric bool
	isTime    bool
	numeric   *Numeric
	time      *Time
}

func NewParsed() *Parsed {
	return &Parsed{
		numeric: new(Numeric),
		time:    new(Time),
	}
}

func (p Parsed) String() string {
	return p.original
}

func (p Parsed) Numeric() (*Numeric, bool) {
	if !p.isNumeric {
		return nil, false
	}

	return p.numeric, true
}

func (p Parsed) Time() (*Time, bool) {
	if !p.isTime {
		return nil, false
	}

	return p.time, true
}

func (p Parsed) Type() string {
	if p.isNumeric {
		return p.numeric.Type()
	}

	if p.isTime {
		return p.time.Type()
	}

	return "None"
}
