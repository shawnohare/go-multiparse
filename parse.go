// Package multiparse provides tools to perform basic type detection and
// parsing on strings.
package multiparse

var context = struct {
	p *Parser
}{
	p: NewParser(),
}

// MultiParse is the general interface all multiparse parsers implement.
type MultiParse interface {
	Parse(string) (interface{}, error)
}

// Parse a string to determine whether it represents a numeric or time value.
// This is a convenience function that is equivalent to calling the
// ParseType method on a general purpose parser instance returned by
// NewGeneralParser.
func Parse(s string) (*Parsed, error) {
	return context.p.parse(s)
}

// ParseType determines whether a numeric or time representation
// according to the initialzed parsers.
func ParseType(s string) (*Parsed, error) {
	return context.p.ParseType(s)
}

// ParseInt reports whether the string parses to an integer according
// to the parser rules.
func ParseInt(s string) (int, error) {
	return context.p.ParseInt(s)
}

// ParseFloat reports whether the string parses to a float according
// to the parser rules.
func ParseFloat(s string) (float64, error) {
	return context.p.ParseFloat(s)
}

// ParseMoney reports whether the string parses to a moneytary value
// according to the parser rules.
func ParseMoney(s string) (*Money, error) {
	return context.p.ParseMoney(s)
}

// ParseTime determines whether the input string parses for any of
// a number of common layouts. Calling this function is equivalent to
// constructing a general time parser with NewGeneralTimeParser
// and invoking its ParseTime method.
func ParseTime(s string) (*Time, error) {
	return context.p.ParseTime(s)
}

// ParseNumeric determines whether the string represents a numeric type.
func ParseNumeric(s string) (*Numeric, error) {
	return context.p.ParseNumeric(s)
}

// ParseBool determines whether the string represents a boolean value.
// The strings "0" and "1" are interpreted as Boolean in this case.
func ParseBool(s string) (bool, error) {
	p := NewBooleanParser()
	return p.parse(s)
}
