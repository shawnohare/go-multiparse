package multiparse

import (
	"errors"
	"regexp"
	"time"
)

var commonTimeLayouts = []string{
	time.RFC3339,
	time.RFC3339Nano,
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
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
}

var commonDateLayouts = []string{
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

// Time structure that represents a string which parses as a datetime.
type Time struct {
	parsed string
	layout string
	time   time.Time
}

// Type is the string name for this any Time instance.
func (t Time) Type() string { return "time" }

// Value returns the instance as an interface.
func (t Time) Value() interface{} { return &t }

// String that can parse into a datetime.
func (t Time) String() string { return t.parsed }

// Time returns the parsed string as a time.Time instace .
func (t Time) Time() time.Time { return t.time }

// Layout detected for datetime conversion.
func (t Time) Layout() string { return t.layout }

// TimeParser instances are responsible for parsing a string to determine
// whether it is a datetime representation.  It is simply a container for
// a number of datetime and date layouts.  The parser iterates over
// these layouts and attempts to parse a string against them.
type TimeParser struct {
	timeLayouts []string
	dateLayouts []string
}

// MakeGeneralTimeParser returns a ready to use datetime parser that
// attempts to detect datetimes using a number of standard layouts.
func MakeGeneralTimeParser() *TimeParser {
	return MakeTimeParser(commonTimeLayouts, commonDateLayouts)
}

// MakeTimeParser produces a custom parser that will attempt to parse
// a string using the user input time and date layouts.
func MakeTimeParser(timeLayouts []string, dateLayouts []string) *TimeParser {
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
func (p TimeParser) ParseTime(s string) (*Time, error) {
	return p.parse(s)
}

// The main datetime parsing logic.
func (p TimeParser) parse(s string) (*Time, error) {
	// Determine whether s has a valid layout that includes time.
	for _, layout := range p.timeLayouts {
		if pt, err := time.Parse(layout, s); err == nil {
			t := &Time{
				parsed: s,
				layout: layout,
				time:   pt,
			}
			return t, nil
		}
	}

	// Detect if the input has a date-like substring and try to parse that.
	re := regexp.MustCompile("^\\d{1,4}([-/\\s]\\d{1,4}){2}")
	d := re.FindString(s)

	if d == "" {
		return nil, errors.New(ParseTimeError)
	}

	t := new(Time)
	for _, layout := range p.dateLayouts {
		if pt, err := time.Parse(layout, d); err == nil {
			t = &Time{
				parsed: d,
				layout: layout,
				time:   pt,
			}
			break
		}
	}
	return t, nil
}
