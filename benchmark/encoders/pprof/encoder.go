package pprof

import (
	"io"
	"otelprofiling/parsers"
	"time"

	"google.golang.org/protobuf/proto"
)

type Encoder struct {
	locations map[string]uint64
	functions map[string]uint64
	strings   map[string]int64
	profile   *Profile
}

func New() *Encoder {
	return &Encoder{
		locations: make(map[string]uint64),
		functions: make(map[string]uint64),
		strings:   make(map[string]int64),
		profile: &Profile{
			StringTable: []string{""},
		},
	}
}

func (b *Encoder) Name() string {
	return "pprof"
}

func (b *Encoder) Append(s parsers.Sample) {
	valueSlice := []int64{int64(s.Value)}
	loc := []uint64{}
	for _, l := range s.Stacktrace {
		loc = append(loc, b.newLocation(l))
	}
	sample := &Sample{LocationId: loc, Value: valueSlice}
	b.profile.Sample = append(b.profile.Sample, sample)
}

func (b *Encoder) Serialize(w io.Writer) error {
	b.profile.SampleType = []*ValueType{{Type: b.newString("cpu"), Unit: b.newString("samples")}}
	b.profile.TimeNanos = time.Now().UnixNano()
	b.profile.DurationNanos = (10 * time.Second).Nanoseconds()

	res, err := proto.Marshal(b.profile)
	if err != nil {
		return err
	}

	_, err = w.Write(res)
	if err != nil {
		return err
	}

	return nil
}

func (b *Encoder) newString(value string) int64 {
	id, ok := b.strings[value]
	if !ok {
		id = int64(len(b.profile.StringTable))
		b.profile.StringTable = append(b.profile.StringTable, value)
		b.strings[value] = id
	}
	return id
}

func (b *Encoder) newLocation(location string) uint64 {
	id, ok := b.locations[location]
	if !ok {
		id = uint64(len(b.profile.Location) + 1)
		newLoc := &Location{
			Id:   id,
			Line: []*Line{{FunctionId: b.newFunction(location)}},
		}
		b.profile.Location = append(b.profile.Location, newLoc)
		b.locations[location] = newLoc.Id
	}
	return id
}

func (b *Encoder) newFunction(function string) uint64 {
	id, ok := b.functions[function]
	if !ok {
		id = uint64(len(b.profile.Function) + 1)
		name := b.newString(function)
		newFn := &Function{
			Id:         id,
			Name:       name,
			SystemName: name,
		}
		b.functions[function] = id
		b.profile.Function = append(b.profile.Function, newFn)
	}
	return id
}
