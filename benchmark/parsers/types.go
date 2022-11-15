package parsers

import "io"

type Parser interface {
	Parse(r io.Reader) ([]Sample, error)
}

type Sample struct {
	Stacktrace []string
	Value      int
}
