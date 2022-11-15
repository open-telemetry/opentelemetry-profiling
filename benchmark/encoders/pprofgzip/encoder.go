package pprofgzip

import (
	"compress/gzip"
	"io"
	"otelprofiling/encoders/pprof"
	"otelprofiling/parsers"
)

type Encoder struct {
	pprofEncoder *pprof.Encoder
}

func New() *Encoder {
	return &Encoder{
		pprofEncoder: pprof.New(),
	}
}

func (b *Encoder) Name() string {
	return "pprof-gzip"
}

func (b *Encoder) Append(s parsers.Sample) {
	b.pprofEncoder.Append(s)
}

func (b *Encoder) Serialize(w io.Writer) error {
	gw := gzip.NewWriter(w)
	defer gw.Close()
	return b.pprofEncoder.Serialize(gw)
}
