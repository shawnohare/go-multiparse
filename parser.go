package multiparse

import "errors"

// Parser instances determine whether a string is a numeric or
// time representation.  Each Parser instance implements
// the MultiParse interface.  Moreover, it is a wrapper for two more
// MultiParse interfaces: one for parsing numeric strings and the other
// for parsing datetime strings.
type Parser struct {
	numeric MultiParse
	time    MultiParse
}

// NewGeneralParser constructs a general purpose top-level Parser instance.
// It is initialized with the  general numeric and time parsers provided
// by NewGeneralNumericParser and NewGeneralTimeParser.
func NewGeneralParser() *Parser {
	return NewParser(NewGeneralNumericParser(), NewGeneralTimeParser())
}

// NewUSDParser constructs a top-level Parser instance that can more
// accurately detect USD monetary strings. For example, it will
// parse "$123,456" as a monetary integer.
func NewUSDParser() *Parser {
	return NewParser(NewUSDNumericParser(), NewGeneralTimeParser())
}

// NewParser is a general purpose parser that uses the passed in
// MultiParse interfaces to determine whether a string is a numeric or
// time representation.  The provided parsers should return *Numeric and
// *Time instances, respectively.
func NewParser(numericParser MultiParse, timeParser MultiParse) *Parser {
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

// ParseInt reports whether the string parses to an integer according
// to the parser rules.
func (p Parser) ParseInt(s string) (int, error) {
	parsed, err := p.parse(s)
	if err != nil {
		return 0, err
	}

	if f, ok := parsed.Int(); ok {
		return f, nil
	}

	return 0, errors.New(ParseIntError)
}

// ParseFloat reports whether the string parses to a float according
// to the parser rules.
func (p Parser) ParseFloat(s string) (float64, error) {
	parsed, err := p.parse(s)
	if err != nil {
		return 0.0, err
	}

	if f, ok := parsed.Float(); ok {
		return f, nil
	}

	return 0.0, errors.New(ParseFloatError)
}

// ParseMoney reports whether the string parses to a moneytary value
// according to the parser rules.
func (p Parser) ParseMoney(s string) (*Money, error) {
	parsed, err := p.parse(s)
	if err != nil {
		return nil, err
	}

	if f, ok := parsed.Money(); ok {
		return f, nil
	}

	return nil, errors.New(ParseMoneyError)
}

// ParseTime reports whether the string parses to a datetime
// according to the parser rules.
func (p Parser) ParseTime(s string) (*Time, error) {
	parsed, err := p.parse(s)
	if err != nil {
		return nil, err
	}

	if f, ok := parsed.Time(); ok {
		return f, nil
	}

	return nil, errors.New(ParseTimeError)
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
