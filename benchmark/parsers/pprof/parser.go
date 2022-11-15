package pprof

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"

	"otelprofiling/encoders/pprof"
	"otelprofiling/parsers"

	"github.com/jzelinskie/must"
	"google.golang.org/protobuf/proto"
)

type Parser struct{}

func (_ *Parser) Parse(r io.Reader) ([]parsers.Sample, error) {
	res := make([]parsers.Sample, 0)

	p := must.NotError(unmarshalProtobuf(r))

	locMap := make(map[uint64]*pprof.Location)
	funcMap := make(map[uint64]*pprof.Function)

	for _, l := range p.Location {
		locMap[l.Id] = l
	}
	for _, f := range p.Function {
		funcMap[f.Id] = f
	}
	for _, s := range p.Sample {
		stack := []string{}
		for _, locid := range s.LocationId {
			loc := locMap[locid]
			for _, l := range loc.Line {
				f := funcMap[l.FunctionId]
				stack = append(stack, p.StringTable[f.Name])
			}
		}
		res = append(res, parsers.Sample{Stacktrace: stack, Value: int(s.Value[0])})
	}

	return res, nil
}

var gzipMagicBytes = []byte{0x1f, 0x8b}

func unmarshalProtobuf(r io.Reader) (*pprof.Profile, error) {
	// this allows us to support both gzipped and not gzipped pprof
	// TODO: this might be allocating too much extra memory, maybe optimize later
	bufioReader := bufio.NewReader(r)
	header, err := bufioReader.Peek(2)
	if err != nil {
		return nil, fmt.Errorf("unable to read profile file header: %w", err)
	}

	if header[0] == gzipMagicBytes[0] && header[1] == gzipMagicBytes[1] {
		r, err = gzip.NewReader(bufioReader)
		if err != nil {
			return nil, err
		}
	} else {
		r = bufioReader
	}

	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	profile := &pprof.Profile{}
	if err := proto.Unmarshal(b, profile); err != nil {
		return nil, err
	}

	return profile, nil
}
