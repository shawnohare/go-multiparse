package multiparse

type Parser interface {
	Parse(string) (interface{}, error)
}
