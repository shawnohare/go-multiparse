package multiparse

import (
	"errors"
	"time"
)

// Parser instances determine whether a string is a numeric or
// time representation.  Each Parser instance implements
// the Interface interface.  Moreover, it is a wrapper for two more
// Interface interfaces: one for parsing numeric strings and the other
// for parsing datetime strings.
type Parser struct {
	numeric Interface
	time    Interface
	b       Interface
}

// NewGeneralParser constructs a general purpose top-level Parser instance.
// It is initialized with the  general numeric and time parsers provided
// by NewGeneralNumericParser and NewGeneralTimeParser.
func NewParser() *Parser {
	n := NewNumericParser()
	t := NewTimeParser()
	b := NewBooleanParser()
	return NewCustomParser(n, t, b)
}

// NewUSDParser constructs a top-level Parser instance that can more
// accurately detect USD monetary strings. For example, it will
// parse "$123,456" as a monetary integer.
func NewUSDParser() *Parser {
	n := NewUSDNumericParser()
	t := NewTimeParser()
	b := NewBooleanParser()
	return NewCustomParser(n, t, b)
}

// NewParser is a general purpose parser that uses the passed in
// Interface interfaces to determine whether a string is a numeric or
// time representation.  The provided parsers should return *Numeric and
// *Time instances, respectively.
func NewCustomParser(numeric, time, boolean Interface) *Parser {
	return &Parser{
		numeric: numeric,
		time:    time,
		b:       boolean,
	}
}

// Parse a string to determine if it is a numeric or monetary value.
// This method is defined primarily so that the Parser struct satifies
// the Interface interface.
func (p Parser) Parse(s string) (interface{}, error) {
	return p.parse(s)
}

// ParseType determines whether a numeric or time representation
// according to the initialzed parsers.
func (p Parser) ParseType(s string) (*Parsed, error) {
	return p.parse(s)
}

// ParseTime reports whether the string parses to a datetime
// according to the parser rules.
func (p Parser) ParseTime(s string) (time.Time, error) {
	parsed, err := p.parse(s)
	if err != nil || !parsed.isTime {
		var t time.Time
		return t, errors.New(ParseTimeError)
	}
	return parsed.time, nil
}

func (p Parser) ParseNumeric(s string) (*Numeric, error) {
	parsed, err := p.parse(s)
	if err != nil || !parsed.isNumeric {
		return nil, errors.New(ParseNumericError)
	}
	return parsed.Numeric, nil
}

// ParseInt reports whether the string parses to an integer according
// to the parser rules.
func (p Parser) ParseInt(s string) (int, error) {
	parsed, err := p.parse(s)
	if err != nil || !parsed.isInt {
		return 0, errors.New(ParseIntError)
	}
	return parsed.Int(), nil
}

// ParseFloat reports whether the string parses to a float according
// to the parser rules.
func (p Parser) ParseFloat(s string) (float64, error) {
	parsed, err := p.parse(s)
	if err != nil || !parsed.isFloat {
		return 0.0, errors.New(ParseFloatError)
	}
	return parsed.Float(), nil
}

func (p Parser) ParseBool(s string) (bool, error) {
	parsed, err := p.parse(s)
	if err != nil || !parsed.isBool {
		return false, errors.New(ParseBoolError)
	}
	return parsed.b, nil
}

// parse a string to determine if it is a valid numeric or time value
// Error when either the underlying parsers return values that cannot
// convert to the appropriate types or when the string does not
// parse into either a numeric or time type.
func (p Parser) parse(s string) (*Parsed, error) {
	var numericAssertErr error
	var timeAssertErr error
	var boolAssertErr error

	parsed := NewParsed()

	x, numericError := p.numeric.Parse(s)
	if numericError == nil {
		switch t := x.(type) {
		case *Numeric:
			parsed.isNumeric = true
			parsed.Numeric = t
		default:
			numericAssertErr = errors.New(ParseTypeAssertError)
		}
	}

	ti, timeError := p.time.Parse(s)
	if timeError == nil {
		switch t := ti.(type) {
		case time.Time:
			parsed.isTime = true
			parsed.time = t
		default:
			timeAssertErr = errors.New(ParseTypeAssertError)
		}
	}

	b, boolError := p.b.Parse(s)
	if boolError == nil {
		switch t := b.(type) {
		case bool:
			parsed.isBool = true
			parsed.b = t
		default:
			boolAssertErr = errors.New(ParseBoolError)
		}
	}

	if numericAssertErr != nil || timeAssertErr != nil || boolAssertErr != nil {
		err := numericAssertErr.Error()
		err += timeAssertErr.Error()
		err += boolAssertErr.Error()
		return nil, errors.New(err)
	}

	if numericError != nil && timeError != nil && boolError != nil {
		return nil, errors.New(ParseError)
	}

	return parsed, nil
}
