package multiparse

import "errors"

type MultiParse interface {
	Parse(string) (interface{}, error)
}

type Parser struct {
	numeric MultiParse
	time    MultiParse
}

func MakeGeneralParser() *Parser {
	return MakeParser(MakeGeneralNumericParser(), MakeGeneralTimeParser())
}

func MakeParser(numericParser MultiParse, timeParser MultiParse) *Parser {
	return &Parser{
		numeric: numericParser,
		time:    timeParser,
	}
}

// Parse a string to determine if it is a numeric or monetary value.
// If so, return the value as a *Parsed instance.
func (p Parser) Parse(s string) (interface{}, error) {
	return p.parse(s)
}

func (p Parser) ParseType(s string) (*Parsed, error) {
	return p.parse(s)
}

// parse a string to determine if it is a valid numeric or time value
// Error when either the underlying parsers return values that cannot
// convert to the appropriate types or when the string does not
// parse into either a numeric or time type.
func (p Parser) parse(s string) (*Parsed, error) {
	var numericAssertError error
	var timeAssertError error

	parsed := NewParsed()
	parsed.original = s

	x, numericError := p.numeric.Parse(s)
	if numericError == nil {
		switch t := x.(type) {
		case *Numeric:
			parsed.isNumeric = true
			parsed.numeric = t
		default:
			numericAssertError = errors.New(ParseTypeAssertError)
		}
	}

	ti, timeError := p.time.Parse(s)
	if timeError == nil {
		switch t := ti.(type) {
		case *Time:
			parsed.isTime = true
			parsed.time = t
		default:
			timeAssertError = errors.New(ParseTypeAssertError)
		}
	}

	if numericAssertError != nil || timeAssertError != nil {
		err := errors.New(numericAssertError.Error() + timeAssertError.Error())
		return nil, err
	}

	if numericError != nil && timeError != nil {
		return nil, errors.New(ParseError)
	}

	return parsed, nil
}

func Parse(s string) (*Parsed, error) {
	p := MakeGeneralParser()
	return p.parse(s)
}
