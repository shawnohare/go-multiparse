package multiparse

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseTime(t *testing.T) {
	timepasses := []string{
		"2009-01-02T15:04:05Z",
		"2009-01-02T15:04:05-07:00",
		"2009-01-02T15:04:05-0700",
		"2009-01-02 15:04:05-0700",
	}

	datepasses := []string{
		"2009-01-02",
		"2009/01/02",
		"01/02/2009",
		"02/01/2009",
		"02/01/2009Tflaksdfj",
	}

	failures := []string{
		"abc",
		"",
		"134.00",
		"$134.00",
	}

	for _, st := range timepasses {
		tt, err := ParseTime(st)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, tt.Hour())
		// t.Log(tt.time)
	}

	for _, st := range datepasses {
		tt, err := ParseTime(st)
		assert.NoError(t, err)
		assert.NotEmpty(t, st, tt.String())
		assert.Equal(t, 0, tt.Hour())
		// t.Log(tt.time)
	}

	for _, st := range failures {
		_, err := ParseTime(st)
		assert.Error(t, err)
	}

}

func TestTimeParserParseType(t *testing.T) {
	p := NewTimeParser()
	expected, _ := time.Parse("2006-01-02", "2015-01-02")
	actual, err := p.ParseTime("2015-01-02")
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
