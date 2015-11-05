package multiparse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFail(t *testing.T) {
	fails := []string{
		"class 1",
		"$$",
		"blah",
		"abc",
	}

	for _, tt := range fails {
		t.Log(tt)
		p, err := ParseType(tt)
		if p != nil {
			t.Log(p.IsFloat())
			t.Log(p.IsMoney())
		}
		assert.Error(t, err)
		assert.Nil(t, p)
	}
}
