package multiparse

import (
	"errors"
	"log"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

const ParseMonetaryStringError = "Cannot parse string as a monetary value."

type Money struct {
	original string
	parsed   string
}

type MoneyParser struct {
	CurrencySymbol   string
	DigitSeparator   string
	DecimalSeparator string
	// Unexported fields.
	digitReStr   string
	decimalReStr string
	digitRegex   *regexp.Regexp
	decimalRegex *regexp.Regexp
}

func MakeMoneyParser(currencySym, digitSep, decimalSep string) *MoneyParser {
	p := &MoneyParser{
		CurrencySymbol:   currencySym,
		DigitSeparator:   digitSep,
		DecimalSeparator: decimalSep,
	}

	// Regex (string) delimiter creator.
	f := func(t string) string {
		var r string
		switch t {
		case "":
			r = "[\\.,]"
		case ".":
			r = "[\\.]"
		case ",":
			r = "[,]"
		}
		return r
	}

	p.digitReStr = f(p.DigitSeparator)
	p.decimalReStr = f(p.DecimalSeparator)
	p.digitRegex = regexp.MustCompile(p.digitReStr)
	p.decimalRegex = regexp.MustCompile(p.decimalReStr)

	return p
}

func MakeStandardMoneyParser() *MoneyParser {
	return MakeMoneyParser("$", ",", ".")
}

func (p MoneyParser) removeCurrencySymbol(s string) string {
	var cleaned string
	if p.CurrencySymbol != "" {
		cleaned = strings.Replace(s, p.CurrencySymbol, "", 1)
	} else {
		re := regexp.MustCompile("^[^0-9-\\+\\.]+")
		loc := re.FindStringIndex(s)
		if len(loc) == 2 {
			cleaned = s[loc[1]:]
		} else {
			cleaned = s
		}
	}
	return cleaned
}

func (p MoneyParser) Parse(s string) (*Money, error) {
	parsed, err := p.parse(s)
	if err != nil {
		return nil, err
	}

	m := &Money{
		original: s,
		parsed:   parsed,
	}
	return m, nil

}

// parse a string and determine whether it passes some basic validation tests.
func (p MoneyParser) parse(s string) (string, error) {
	var (
		sign     string
		reStr    string
		re       *regexp.Regexp
		err      error
		parseErr error = errors.New(ParseMonetaryStringError)
	)

	log.Println("Initial input:", s)

	// Ensure the input has at least one digit.
	re = regexp.MustCompile(".*[0-9].*")
	if !re.MatchString(s) {
		return "", parseErr
	}

	// Remove the first currency symbols that appear.
	s = p.removeCurrencySymbol(s)
	log.Println("Removed currency symbols:", s)

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
	log.Println("main validation re:", re.String())
	if !re.MatchString(s) {
		// log.Println("Didn't pass main validating regexp:", s)
		return "", parseErr
	}
	log.Println("Input passed main regex test:", s)

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

	log.Println("Parsed string:", parsed)

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

// Parse a string representing a monetary value and covert it to a
// *Money instance.  The separator parameters denote any optional formatting
// and decimal separators (e.g., "," and "." in "123,456.05", resp.).
// If the separators are "" then the function tries to automatically detect
// formatting and decimal separators.
func ParseMonetaryString(s string) (*Money, error) {
	// Make a generic money parser with no opinion
	parser := MakeMoneyParser("", "", "")
	return parser.Parse(s)
}
