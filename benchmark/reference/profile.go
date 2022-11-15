package reference

import (
	"encoding/json"
	"io"
)

type line struct {
	Stacktrace string
	Value      int
}

type Profile struct {
	Lines []line
}

func New() *Profile {
	return &Profile{}
}

func (p *Profile) Append(stacktrace string, value int) {
	p.Lines = append(p.Lines, line{Stacktrace: stacktrace, Value: value})
}

func (p *Profile) Serialize(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}
