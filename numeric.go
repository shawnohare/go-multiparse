package multiparse

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// Numeric instances are containers for the various valid numerical types
// that a string may be parsed into.
type Numeric struct {
	isInt   bool
	isFloat bool
	isMoney bool
	f       float64
}

// A NumericParser ingests a string and determines whether it is
// numeric or monetary value by using
// its configuration dictionary.  This dictionary consists
// of a currency symbol, a digit separator and a decimal separator.
type NumericParser struct {
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

// NewGeneralNumericParser with the most general dictionary.  The dictionary
// values are all empty strings, and so this parser is agnostic to
// whether "," indicates a digit or decimal separator.  Moreover, it
// considers anything in the compliment of { +, -, 0, 1, ..., 9, \s }
// to be a currency symbol.
//
// This parser interprets "123.456" and "123,456" as integer values.
func NewNumericParser() *NumericParser {
	return NewCustomNumericParser("", "", "")
}

// NewUSDNumericParser is configured with a currency symbol "$",
// digit separator ",", and decimal separator ".".
// It properly detects that "123.456" is a real number, but not an integer.
func NewUSDNumericParser() *NumericParser {
	return NewCustomNumericParser("$", ",", ".")
}

// NewNumericParser with the dictionary defined by the inititialization
// parameters.
//
// Valid inputs for the currency symbol are: "", "$", or any
// regular expression.
//
// Valid inputs for the separators are: "", ".", ",", or any
// regular expression.
func NewCustomNumericParser(currencySym, digitSep, decimalSep string) *NumericParser {
	p := &NumericParser{
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

	// Construct the predefined currency regex (dcre) in a series of steps.
	tmp := []string{
		"^\\p{Sc}",                    // Perl currency symbol
		"^Дин\\.",                     // Serbian Dinar
		"^p\\.",                       // Belarus Ruble
		"^[SB]/\\.",                   // Peru Nuevo Sol and Panama Balboa
		"^[^0-9\\.\\+\\s-_]{1,3}\\s?", // Match things like Lek, HK$ or "USD "
	}
	cre := strings.Join(tmp, "|")

	currencyMap := map[string]string{
		"":  cre,
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

// Parse a string to determine if it represents a numeric type.
func (p NumericParser) Parse(s string) (interface{}, error) {
	return p.parse(s)
}

// ParseNumeric a string to determine if it represents a numeric type.
func (p NumericParser) ParseNumeric(s string) (*Numeric, error) {
	return p.parse(s)
}

func (p NumericParser) removeCurrencySymbol(s string) string {
	loc := p.currencyRegex.FindStringIndex(s)
	if len(loc) == 2 {
		return s[loc[1]:]
	}
	return s
}

func (p NumericParser) removeDigitSeparators(s string) (string, error) {
	if p.digitReStr == p.decimalReStr {
		return "", errors.New(ParseMoneySeparatorError)
	}

	cleaned := p.digitRegex.ReplaceAllString(s, "")
	return cleaned, nil
}

// replace the last occurence of a decimal separator with "."
func (p NumericParser) replaceDecimalSeparator(s string) (string, error) {
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

func (p NumericParser) sanitize(s string) (string, error) {

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

// Parse a string representation of a money value which has one "." or ",".
// Any string passed in should not begin or end with a delimiter.
func (p NumericParser) parseOneUnknownSeparator(m string, i int) (string, error) {
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
func (p NumericParser) parseManyUnknownSeparators(m string, locs [][]int) (string, error) {
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

func (p NumericParser) parse(s string) (*Numeric, error) {
	var (
		n        *Numeric
		err      error
		sign     string
		reStr    string
		re       *regexp.Regexp
		parseErr = errors.New(ParseNumericError)
	)

	// Record whether the input string has a currency symbol.
	// If so, it can only be a monetary value.
	hasCurrency := p.currencyRegex.MatchString(s)
	if hasCurrency {
		s = p.removeCurrencySymbol(s)
	}

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

	// Since currency and sign symbols have been stripped, we now check that the
	// expression begins with a decimal separator (possibly) and digit.
	// Valid strings thus look like either: .x* or x*.
	reStr = "^" + p.decimalReStr + "?" + "[0-9]"
	re = regexp.MustCompile(reStr)
	if !re.MatchString(s) {
		return nil, parseErr
	}

	// Prepend a 0 if the string begins with a decimal separator.
	reStr = "^" + p.decimalReStr
	re = regexp.MustCompile(reStr)
	if re.MatchString(s) {
		s = "0" + s
	}

	// If the input ends with the decimal separator, remove it.
	re = regexp.MustCompile(p.decimalReStr + "$")
	if re.MatchString(s) {
		s = re.ReplaceAllString(s, "")
	}

	// Create the main validating regex.
	reStr = "^\\d+" + "(" + p.digitReStr + "\\d{3})*" + p.decimalReStr + "?\\d*$"
	re = regexp.MustCompile(reStr)
	if !re.MatchString(s) {
		return nil, parseErr
	}

	// We can now assume that the string is valid except for
	// intermediate delimiters.
	// Before attempting to parse the string further, we (possibly) perform
	// some basic sanitization.
	var parsed string
	tmp, err := p.sanitize(s)
	if err == nil {
		parsed = tmp
	} else {
		// Probably the parser cannot distinguish between decimal and digit
		// separators.  So we handle this case separately.
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

	parsed = sign + parsed
	f, ferr := strconv.ParseFloat(parsed, 64)
	if err != nil || ferr != nil {
		return nil, err
	}

	// We now know that the parsed string correctly parses as a float.
	n = &Numeric{
		isFloat: true,
		f:       f,
	}
	if hasCurrency {
		n.isMoney = true
	}
	_, err = strconv.Atoi(parsed)
	if err == nil {
		n.isInt = true
	}

	return n, nil
}

// Int reports whether the Numeric instance can be an integer
// and returns its value.
func (x Numeric) Int() int {
	return int(x.f)
}

// Float reports whether the Numeric instance can be a float
// and returns its value.
func (x Numeric) Float() float64 {
	return x.f
}

// IsInt reports if the instance can represent an integer.
func (x Numeric) IsInt() bool {
	return x.isInt
}

// IsFloat reports if the instance can represent a floating point number.
func (x Numeric) IsFloat() bool {
	return x.isFloat
}

// IsMoney reports if the original string has currency symbols
// and correctly prases to a numeric value.
func (x Numeric) IsMoney() bool {
	return x.isMoney
}
