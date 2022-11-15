package json

import (
	"encoding/json"
	"io"
	"otelprofiling/parsers"
	"strings"
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

func (p *Profile) Name() string {
	return "json"
}

func (p *Profile) Append(s parsers.Sample) {
	p.Lines = append(p.Lines, line{Stacktrace: strings.Join(s.Stacktrace, ";"), Value: s.Value})
}

func (p *Profile) Serialize(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}
