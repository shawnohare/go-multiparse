package multiparse

import (
	"errors"
	"regexp"
	"time"
)

// TODO
// - [ ] try to split date from time.
// - [ ] Sanitize both.
// - [ ] see if the input parses in multiple classes.  If so, it's ambiguous.

// expression:
// timeEquivalenceClass is an equivalence class of datetime formats
// with respect to the relation ~ defined by f ~ g if and only if
// T(time.Parse(f, f)) = T(time.Parse(g, g)), where T is the time truncation
// function.
type timeEquivalenceClass struct {
	dateRep string

	dates []string
	times []string
}

var commonTimeLayouts = []string{
	time.RFC3339,
	time.RFC3339Nano,
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	"2006-01-02T15:04:05+07:00", // ISO 8601 in UTC
	"2006-01-02T15:04:05Z0700",
	"2006-01-02 15:04:05Z07:00",
	"2006-01-02 15:04:05Z0700",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,
	"1/2/06 15:04",
	"2/1/06 15:04",
	"Jan. 2 2006 15:04:05",
}

var commonDateLayouts = []string{
	"1/2/06",
	"2/1/06",
	"01-02-06",
	"02-01-06", // This may not ever match
	"1/2/2006",
	"1-2-2006",
	"01-02-2006",
	"02-01-2006",
	"2006/01/02",
	"01/02/2006",
	"02/01/2006",
	"Jan 02 2006",
	"2006-01-02",
	"01-02-2006",
	"02-01-2006",
	"2006/01/02",
	"01/02/2006",
	"Jan 02 2006",
	"Jan. 2 2006",
}

// TimeParser instances are responsible for parsing a string to determine
// whether it is a datetime representation.  It is simply a container for
// a number of datetime and date layouts.  The parser iterates over
// these layouts and attempts to parse a string against them.
type TimeParser struct {
	timeLayouts []string
	dateLayouts []string
}

// NewGeneralTimeParser returns a ready to use datetime parser that
// attempts to detect datetimes using a number of standard layouts.
func NewTimeParser() *TimeParser {
	return NewCustomTimeParser(commonTimeLayouts, commonDateLayouts)
}

// NewTimeParser produces a custom parser that will attempt to parse
// a string using the user input time and date layouts.
func NewCustomTimeParser(timeLayouts []string, dateLayouts []string) *TimeParser {
	return &TimeParser{
		timeLayouts: timeLayouts,
		dateLayouts: dateLayouts,
	}
}

// Parse a string to determine if it represents a datetime.
func (p TimeParser) Parse(s string) (interface{}, error) {
	return p.parse(s)
}

// ParseTime is the same as Parse but returns a *Time instance.
func (p TimeParser) ParseTime(s string) (time.Time, error) {
	return p.parse(s)
}

// The main datetime parsing logic.
func (p TimeParser) parse(s string) (time.Time, error) {
	// Determine whether s has a valid layout that includes time.
	for _, layout := range p.timeLayouts {
		if pt, err := time.Parse(layout, s); err == nil {
			return pt, nil
		}
	}

	// Detect if the input has a date-like substring and try to parse that.
	re := regexp.MustCompile("^\\d{1,4}([-/\\s]\\d{1,4}){2}")
	d := re.FindString(s)

	if d == "" {
		var t time.Time
		return t, errors.New(ParseTimeError)
	}

	for _, layout := range p.dateLayouts {
		if pt, err := time.Parse(layout, d); err == nil {
			return pt, nil
		}
	}

	var t time.Time
	return t, nil
}
