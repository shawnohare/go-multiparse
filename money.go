package multiparse

import (
	"errors"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

// A Money instance is a simple representation of a monetary value.  It
// has access to the original pre-parsed string as well as a few numeric
// types.
type Money struct {
	original string
	parsed   string
}

// A MoneyParser ingests a string and determines whether it is a monetary
// value by using its configuration dictionary.  This dictionary consists
// of a currency symbol, a digit separator and a decimal separator.
type MoneyParser struct {
	CurrencySymbol   string
	DigitSeparator   string
	DecimalSeparator string
	// Unexported fields.
	digitReStr    string
	decimalReStr  string
	currencyReStr string
	digitRegex    *regexp.Regexp
	decimalRegex  *regexp.Regexp
	currencyRegex *regexp.Regexp
}

// NewMoneyParser with the most general dictionary.  The dictionary
// values are all empty strings, and so this parser is agnostic to
// whether "," indicates a digit or decimal separator.  Moreover, it
// considers anything in the compliment of { +, -, 0, 1, ..., 9, \s }
// to be a currency symbol.
//
// This parser interprets "123.456" and "123,456" as integer values.
func NewMoneyParser() *MoneyParser {
	return MakeMoneyParser("", "", "")
}

// NewStandardMoneyParser is configured with a currency symbol "$",
// digit separator ",", and decimal separator ".".
func NewStandardMoneyParser() *MoneyParser {
	return MakeMoneyParser("$", ",", ".")
}

// MakeMoneyParser with the dictionary defined by the inititialization
// parameters.
//
// Valid inputs for the currency symbol are: "", "$", or any
// regular expression.
//
// Valid inputs for the separators are: "", ".", ",", or any
// regular expression.
func MakeMoneyParser(currencySym, digitSep, decimalSep string) *MoneyParser {
	p := &MoneyParser{
		CurrencySymbol:   currencySym,
		DigitSeparator:   digitSep,
		DecimalSeparator: decimalSep,
	}

	// Define the regular expression maps to convert string inputs into valid
	// regular expressions.
	sepMap := map[string]string{
		"":  "[\\.,]",
		".": "[\\.]",
		",": "[,]",
	}

	currencyMap := map[string]string{
		"":  "^[^0-9-\\+\\.]+",
		"$": "^[\\$]",
	}

	// Input -> regex string
	f := func(t string, m map[string]string) string {
		if restr, prs := m[t]; prs {
			return restr
		}

		return t
	}

	p.digitReStr = f(p.DigitSeparator, sepMap)
	p.decimalReStr = f(p.DecimalSeparator, sepMap)
	p.currencyReStr = f(p.CurrencySymbol, currencyMap)
	p.digitRegex = regexp.MustCompile(p.digitReStr)
	p.decimalRegex = regexp.MustCompile(p.decimalReStr)
	p.currencyRegex = regexp.MustCompile(p.currencyReStr)

	return p
}

// Parse an input string and return the *Money result as an interface.
func (p MoneyParser) Parse(s string) (interface{}, error) {
	return p.parse(s)
}

// ParseMoney parses the input string and return the result
// as a *Money instance.
func (p MoneyParser) ParseMoney(s string) (*Money, error) {
	return p.parse(s)
}

func (p MoneyParser) removeCurrencySymbol(s string) string {
	loc := p.currencyRegex.FindStringIndex(s)
	if len(loc) == 2 {
		return s[loc[1]:]
	}
	return s
}

func (p MoneyParser) removeDigitSeparators(s string) (string, error) {
	if p.digitReStr == p.decimalReStr {
		return "", errors.New(ParseMoneySeparatorError)
	}

	cleaned := p.digitRegex.ReplaceAllString(s, "")
	return cleaned, nil
}

// replace the last occurence of a decimal separator with "."
func (p MoneyParser) replaceDecimalSeparator(s string) (string, error) {
	if p.digitReStr == p.decimalReStr {
		return "", errors.New(ParseMoneySeparatorError)
	}

	// Do not need to do anything if the decimal separator is already "."
	if p.decimalReStr == "[\\.]" {
		return s, nil
	}

	locs := p.decimalRegex.FindAllStringIndex(s, -1)
	if len(locs) == 0 {
		return s, nil
	}

	loc := locs[len(locs)-1]
	cleaned := s[:loc[0]] + "." + s[loc[1]:]
	return cleaned, nil
}

func (p MoneyParser) sanitize(s string) (string, error) {

	s = p.removeCurrencySymbol(s)
	tmp, err1 := p.removeDigitSeparators(s)
	if err1 == nil {
		s = tmp
	}
	tmp, err2 := p.replaceDecimalSeparator(s)
	if err2 == nil {
		s = tmp
	}

	var err error
	if err1 != nil || err2 != nil {
		err = errors.New(err1.Error() + err2.Error())
	}
	return s, err
}

func (p MoneyParser) parse(s string) (*Money, error) {
	var (
		sign     string
		reStr    string
		re       *regexp.Regexp
		err      error
		parseErr = errors.New(ParseMonetaryStringError)
		original = s
	)

	// Ensure the input has at least one digit.
	re = regexp.MustCompile(".*[0-9].*")
	if !re.MatchString(s) {
		return nil, parseErr
	}

	// Remove the first currency symbols that appear.
	s = p.removeCurrencySymbol(s)

	// Now determine whether the string's initial character is a + or -.
	// If so, strip it away and record the sign.
	sign = ""
	re = regexp.MustCompile("^[\\+-]")
	if re.MatchString(s) {
		if re.FindString(s) == "-" {
			sign = "-"
		}
		s = s[1:]
	}

	// A valid string now either begins with digits or a decimal separator.
	// If it begns with the later, prepend a 0.
	reStr = "^" + p.decimalReStr
	re = regexp.MustCompile(reStr)
	if re.MatchString(s) {
		s = "0" + s
	}

	// Append decimal zeros if necessary.  E.g., 123. -> 123.00
	reStr = "^.*" + p.decimalReStr + "$"
	re = regexp.MustCompile(reStr)
	if re.MatchString(s) {
		s = s + "00"
	}

	// Create the main validating regex.
	reStr = "^\\d+" + "(" + p.digitReStr + "\\d{3})*" + p.decimalReStr + "?\\d*$"
	re = regexp.MustCompile(reStr)
	if !re.MatchString(s) {
		return nil, parseErr
	}

	// We can now assume that the string is valid except for extra delimiters.
	// Before attempting to parse the string further, we (possibly) perform
	// some basic sanitization.
	var parsed string
	tmp, err := p.sanitize(s)
	if err == nil {
		parsed = tmp
	} else {
		re = regexp.MustCompile(p.digitReStr + "|" + p.decimalReStr)
		locs := re.FindAllStringSubmatchIndex(s, -1)
		switch len(locs) {
		case 0: // The number is an integer.  No additional parsing needed.
			parsed = s
			err = nil
		case 1: // Need to deal with 1,234 vs 123,456 vs 12.345, etc.
			parsed, err = p.parseOneUnknownSeparator(s, locs[0][0])
		default: // Try to find the last separator and determine its type.
			parsed, err = p.parseManyUnknownSeparators(s, locs)
		}

	}

	if err != nil {
		return nil, err
	}

	m := &Money{
		original: original,
		parsed:   sign + parsed,
	}
	return m, nil

}

// Parse a string representation of a money value which has one "." or ",".
// Any string passed in should not begin or end with a delimiter.
func (p MoneyParser) parseOneUnknownSeparator(m string, i int) (string, error) {
	// Split the string at the delimiter.
	before := m[:i]
	after := m[i+1:]

	// Initially, assume the after var contins the non- integral part of m.
	decimalSep := "."
	if len(before) <= 3 && len(after) == 3 {
		// Then m looks like "1,238" or "123.456" and we assume m is integral.
		decimalSep = ""
	}

	fs := strings.Join([]string{before, after}, decimalSep)
	return fs, nil
}

// Parse a string representation of a money value which has many "." or ",".
func (p MoneyParser) parseManyUnknownSeparators(m string, locs [][]int) (string, error) {
	i := locs[len(locs)-1][0]
	j := locs[len(locs)-2][0]
	ci := string(m[i]) // last delimiter
	cj := string(m[j]) // second to last delimiter

	// Split the string at the last delimiter we encounter.
	before := m[:i]
	after := m[i+1:]

	decimalSep := "."
	if ci == cj {
		// If the last two delimiters are equal, check whether all delimiters
		// are equal.  If not, error, otherwise, we assume the string
		// represents an integer.
		for i, loc := range locs {
			if i == len(locs)-1 {
				break
			}
			if m[loc[0]] != m[locs[i+1][0]] {
				return "", errors.New(ParseMonetaryStringError)
			}
		}
		decimalSep = ""
	}

	// Remove all non-decimal indicator delimiters.
	re := regexp.MustCompile(p.digitReStr)
	before = re.ReplaceAllString(before, "")

	fs := strings.Join([]string{before, after}, decimalSep)
	return fs, nil
}

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
	return strconv.ParseFloat(m.parsed, 64)
}

// BigFloat returns a big.Float representation of the monetary value.
// This type is more appropriate for accounting.
func (m Money) BigFloat() (*big.Float, error) {
	tmp := new(big.Float)
	bf, _, err := tmp.Parse(m.parsed, 10)
	return bf, err
}

// ParseMonetaryString parses an input string representing a monetary value
// and returns the *Money result.  This convenience function utilizes the
// generic MoneyParser returned by NewMoneyParser.
func ParseMonetaryString(s string) (*Money, error) {
	// Make a generic money parser with no opinion
	parser := MakeMoneyParser("", "", "")
	return parser.ParseMoney(s)
}
