# multiparse
Basic string parser written in Go.

Currently the parser focuses on determining whether a string can be an 
integer, float, monetary or time value.

### Parse

For single use applications the package provides the 
general-purpose `Parse(string) (*Parsed, error)`
function that will parse an input string and attempt to detect whether
it represents an integer, floating point number, moneytary value,
or a datetime.  It will interpret `"123.456"` and `"123,456"` as
integers, however, as this function does not take an opinion on
digit separators.


### Parsers

To parse multiple strings, we recommend using a MultiParse interface.
The `NewGeneralParser` function returns a general-purpose parser that
is read to use.  The `Parser.ParseType` method for this general-purpose
parser is equivalent to the package's `Parse`, except it avoids
uncessary struct initialization.

## USD 

The `NewUSDParser` function returns a ready to use parser that
recognizes `"$"` as the only currency symbol and assumes that the
digit and decimal separators are `","` and `"."`, respectively.


## Basic Usage 

```go
package main

import (
	mp "github.com/shawnohare/go-multiparse"
	"log"
	"reflect"
)

func main() {
	// Parsing a monetary string.
	s := "$123,456"
	parsed, err := mp.Parse(s)
	log.Println(err == nil)              // true
	log.Println(parsed.Type())           // money
	log.Println(parsed.IsNumeric())      // true
	log.Println(parsed.IsMoney())        // true
	log.Println(parsed.IsInt())          // false
	if money, ok := parsed.Money(); ok { // ok is true
		log.Println(money.Float64()) // 123456, nil (no error)
	}

	// Parsing a date string.
	d := "2015-01-02"
	parsed, err = mp.Parse(d)
	log.Println(err == nil)         // true
	log.Println(parsed.IsTime())    // true
	log.Println(parsed.IsNumeric()) // false
	// parsed.Time() returns a small wrapper for a time.Time instance.
	if t, ok := parsed.Time(); ok {
		log.Println(t.String()) // 2015-01-02
		log.Println(t.Layout()) // 2006-01-02
		// We can obtain the underlying time.Time instance via:
		log.Println(t.Time()) // 2015-01-02 00:00:00 +0000 UTC
	}

	// A general purpose parser that will produce the same results as above
	// is constructed by the NewGeneralParser method.
	parser := mp.NewGeneralParser()
	// The ParseType method for this parser is equivalent to the package's
	// Parse function.
	parsed1, _ := mp.Parse("123.4")
	parsed2, _ := parser.ParseType("123.4")
	log.Println(reflect.DeepEqual(parsed1, parsed2)) // true
	// The parser's Parse method returns an interface and error.
	I, err := parser.Parse("123.4")
	parsed3 := I.(*mp.Parsed)
	log.Println(reflect.DeepEqual(parsed1, parsed3)) // true

	// A less general but more accurate parser that can handle USD money
	// strings can be constructed via NewUSDParser()
}
```
