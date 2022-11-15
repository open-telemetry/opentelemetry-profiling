package pprofgzip

import (
	"compress/gzip"
	"io"
	"otelprofiling/encoders/pprof"
	"otelprofiling/parsers"
)

type Builder struct {
	Builder *pprof.Builder
}

func New() *Builder {
	return &Builder{
		Builder: pprof.New(),
	}
}

func (b *Builder) Name() string {
	return "pprof-gzip"
}

func (b *Builder) Append(s parsers.Sample) {
	b.Builder.Append(s)
}

func (b *Builder) Serialize(w io.Writer) error {
	gw := gzip.NewWriter(w)
	defer gw.Close()
	return b.Builder.Serialize(gw)
}
