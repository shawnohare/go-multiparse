package multiparse

import "errors"

type BooleanParser struct {
	m map[string]bool
}

func NewBooleanParser() *BooleanParser {
	m := map[string]bool{
		"1":     true,
		"yes":   true,
		"true":  true,
		"0":     false,
		"no":    false,
		"false": false,
	}
	return NewCustomBooleanParser(m)
}

func NewCustomBooleanParser(m map[string]bool) *BooleanParser {
	return &BooleanParser{m}
}

func (p BooleanParser) Parse(s string) (interface{}, error) {
	return p.parse(s)
}

func (p BooleanParser) ParseBool(s string) (bool, error) {
	return p.parse(s)
}

func (p BooleanParser) parse(s string) (bool, error) {
	b, prs := p.m[s]
	if !prs {
		return false, errors.New(ParseBoolError)
	}
	return b, nil
}
