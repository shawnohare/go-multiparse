package multiparse

import (
	"errors"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

// A ParseMonetaryStringError is a generic indication that a string
// could not be parsed as a monetary value.
const ParseMonetaryStringError = "Cannot parse string as a monetary value."

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
		"$": "[\\$]",
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

func (p MoneyParser) removeCurrencySymbol(s string) string {
	var cleaned string
	loc := p.currencyRegex.FindStringIndex(s)
	if len(loc) == 2 {
		cleaned = s[loc[1]:]
	} else {
		cleaned = s
	}
	return cleaned
}

// Parse an input string and return the *Money result as an interface.
func (p MoneyParser) Parse(s string) (interface{}, error) {
	return p.ParseMoney(s)
}

// ParseMoney parses the input string and return the result
// as a *Money instance.
func (p MoneyParser) ParseMoney(s string) (*Money, error) {
	parsed, err := p.parseString(s)
	if err != nil {
		return nil, err
	}

	m := &Money{
		original: s,
		parsed:   parsed,
	}
	return m, nil

}

func (p MoneyParser) parseString(s string) (string, error) {
	var (
		sign     string
		reStr    string
		re       *regexp.Regexp
		err      error
		parseErr = errors.New(ParseMonetaryStringError)
	)

	// log.Println("Initial input:", s)

	// Ensure the input has at least one digit.
	re = regexp.MustCompile(".*[0-9].*")
	if !re.MatchString(s) {
		return "", parseErr
	}

	// Remove the first currency symbols that appear.
	s = p.removeCurrencySymbol(s)
	// log.Println("Removed currency symbols:", s)

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
	// log.Println("Determined the sign to be:", sign, "for ", s)

	// A valid string now either begins with digits or a decimal separator.
	// If it begns with the later, prepend a 0.
	reStr = "^" + p.decimalReStr
	re = regexp.MustCompile(reStr)
	if re.MatchString(s) {
		s = "0" + s
	}

	// A valid string could terminate with a decimal separator.  If so,
	// add some zeros.
	reStr = "^.*" + p.decimalReStr + "$"
	re = regexp.MustCompile(reStr)
	if re.MatchString(s) {
		s = s + "00"
	}

	// Create the main validating regex.
	reStr = "^\\d+" + "(" + p.digitReStr + "\\d{3})*" + p.decimalReStr + "?\\d*$"
	re = regexp.MustCompile(reStr)
	// log.Println("main validation re:", re.String())
	if !re.MatchString(s) {
		// log.Println("Didn't pass main validating regexp:", s)
		return "", parseErr
	}
	// log.Println("Input passed main regex test:", s)

	// We can now assume that the string is valid except for extra delimiters.
	var parsed string
	if p.DigitSeparator != "" && p.DigitSeparator != "" {
		// Remove extraneous digit separators.
		s = strings.Replace(s, p.DigitSeparator, "", -1)
		// Replace any custom decimal separator with one that will parse later.
		s = strings.Replace(s, p.DecimalSeparator, ".", 1)
		// The input string should now be properly sanitized and ready for conversion.
		parsed = s
	} else {
		re = regexp.MustCompile(p.digitReStr + "|" + p.decimalReStr)
		// FIXME
		locs := re.FindAllStringSubmatchIndex(s, -1)
		switch len(locs) {
		case 0:
			// The number is an integer.  No additional parsing needed.
			parsed = s
			err = nil
		case 1:
			parsed, err = p.parseOneUnknownSeparator(s, locs[0][0])
		default:
			parsed, err = p.parseManyUnknownSeparators(s, locs)
		}
	}

	// log.Println("Parsed string:", parsed)

	if err != nil {
		return "", err
	}

	// Add back the sign.
	parsed = sign + parsed

	return parsed, nil

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
