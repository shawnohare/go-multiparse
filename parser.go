// Package multiparse provides tools to perform basic type detection and
// parsing on strings.
package multiparse

import "errors"

// MultiParse is the general interface all multiparse parsers implement.
type MultiParse interface {
	Parse(string) (interface{}, error)
}

// Parser instances determine whether a string is a numeric or
// time representation.  Each Parser instance implements
// the MultiParse interface.  Moreover, it is a wrapper for two more
// MultiParse interfaces: one for parsing numeric strings and the other
// for parsing datetime strings.
type Parser struct {
	numeric MultiParse
	time    MultiParse
}

// MakeGeneralParser constructs a general purpose top-level Parser instance.
// It is initialized with the  general numeric and time parsers provided
// by MakeGeneralNumericParser and MakeGeneralTimeParser.
func MakeGeneralParser() *Parser {
	return MakeParser(MakeGeneralNumericParser(), MakeGeneralTimeParser())
}

// MakeUSDParser constructs a top-level Parser instance that can more
// accurately detect USD monetary strings. For example, it will
// parse "$123,456" as a monetary integer.
func MakeUSDParser() *Parser {
	return MakeParser(MakeUSDNumericParser(), MakeGeneralTimeParser())
}

// MakeParser is a general purpose parser that uses the passed in
// MultiParse interfaces to determine whether a string is a numeric or
// time representation.  The provided parsers should return *Numeric and
// *Time instances, respectively.
func MakeParser(numericParser MultiParse, timeParser MultiParse) *Parser {
	return &Parser{
		numeric: numericParser,
		time:    timeParser,
	}
}

// Parse a string to determine if it is a numeric or monetary value.
// This method is defined primarily so that the Parser struct satifies
// the MultiParse interface.
func (p Parser) Parse(s string) (interface{}, error) {
	return p.parse(s)
}

// ParseType determines whether a numeric or time representation
// according to the initialzed parsers.
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

// Parse a string to determine whether it represents a numeric or time value.
// This is a convenience function that is equivalent to calling the
// ParseType method on a general purpose parser instance returned by
// MakeGeneralParser.
func Parse(s string) (*Parsed, error) {
	p := MakeGeneralParser()
	return p.parse(s)
}
