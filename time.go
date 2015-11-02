package multiparse

import (
	"errors"
	"regexp"
	"time"
)

var commonTimeLayouts = []string{
	// Layouts
	time.RFC3339,
	time.RFC3339Nano,
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	// Handy time stamps.
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,
	// Custom
	"2006-01-02T15:04:05Z0700",
	"2006-01-02 15:04:05Z07:00",
	"2006-01-02 15:04:05Z0700",
}

var commonDateLayouts = []string{
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

// TODO write general regex that recognizes

type Time struct {
	parsed string
	layout string
	time   time.Time
}

func (t Time) Type() string       { return "time" }
func (t Time) Value() interface{} { return &t }
func (t Time) String() string     { return t.parsed }
func (t Time) Time() time.Time    { return t.time }
func (t Time) Layout() string     { return t.layout }

// ParseTime determines whether the input string parses for any of
// a number of common layouts.
func ParseTime(s string) (*Time, error) {
	// Determine whether s has a valid layout that includes time.
	for _, layout := range commonTimeLayouts {
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
	re := regexp.MustCompile("^\\d{2,4}([-/\\s]\\d{2,4}){2}")
	d := re.FindString(s)

	if d == "" {
		return nil, errors.New(ParseTimeError)
	} else {
		t := new(Time)
		for _, layout := range commonDateLayouts {
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
}
