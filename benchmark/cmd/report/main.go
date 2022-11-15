package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"otelprofiling/encoders/json"
	"otelprofiling/encoders/pprof"
	"otelprofiling/encoders/pprofgzip"
	"otelprofiling/parsers"
	pcollapsed "otelprofiling/parsers/collapsed"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jzelinskie/must"
)

type Encoder interface {
	Name() string
	Append(s parsers.Sample)
	Serialize(w io.Writer) error
}

// You can add your encoders here
var encoders = []func() Encoder{
	// first one is always "baseline"
	func() Encoder { return pprofgzip.New() },
	func() Encoder { return pprof.New() },
	func() Encoder { return json.New() },
}

func main() {
	mdReport := ""
	filepaths := must.NotError(filepath.Glob("./profiles/intermediary/*"))
	for _, path := range filepaths {
		basename := filepath.Base(path)
		fmt.Printf("Filename: %s\n", basename)
		mdReport += fmt.Sprintf("### Filename: %s\n", basename)
		fileBytes := must.NotError(ioutil.ReadFile(path))

		// this parser is always "collapsed" because that is the intermediary format
		parser := &pcollapsed.Parser{}
		samples := must.NotError(parser.Parse(bytes.NewReader(fileBytes)))

		t := table.NewWriter()
		t.AppendHeader(table.Row{"Encoder", "Duration", "% diff", "Byte Size", "% diff"})

		var baselineDur time.Duration
		var baselineSize int

		for i, encoder := range encoders {
			buf := &bytes.Buffer{}
			noopWriter := &noopWriter{}

			// step 1: measure byte size
			e := encoder()
			for _, sample := range samples {
				e.Append(sample)
			}
			e.Serialize(buf)
			res := buf.Bytes()
			size := len(res)

			// step 2: save file output for debugging purposes
			dir := filepath.Join("profiles", "out", e.Name())
			os.MkdirAll(dir, 0755)
			ioutil.WriteFile(filepath.Join(dir, filepath.Base(path)), res, 0644)

			// step 3: measure latency
			// we need to make sure this closure does as little extra work as possible
			// so that we're measuring the encoder's performance, not the parser's
			dur := measureEncoderLatency(func() {
				e := encoder()
				for _, sample := range samples {
					e.Append(sample)
				}
				e.Serialize(noopWriter)
			})

			if i == 0 {
				baselineDur = dur
				baselineSize = size
			}

			t.AppendRow([]interface{}{e.Name(), formatDuration(dur), percDiff(int(dur), int(baselineDur)), formatBytes(size), percDiff(size, baselineSize)})
		}
		fmt.Println(t.Render())
		fmt.Println("")
		mdReport += t.RenderMarkdown() + "\n\n"
	}

	ioutil.WriteFile("report.md", []byte(mdReport), 0644)
}

const N = 10

func measureEncoderLatency(cb func()) time.Duration {
	// first call is to warm up various runtime caches
	// if we don't do this it might skew the results
	cb()
	startTime := time.Now()
	for i := 0; i < N; i++ {
		cb()
	}
	return time.Since(startTime) / N
}

type noopWriter struct{}

func (n *noopWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func percDiff(a, b int) string {
	if a == b {
		return ""
	}
	sign := ""
	if a > b {
		sign = "+"
	}
	return fmt.Sprintf("%s%.2f%%", sign, 100*float64(a-b)/float64(b))
}

func formatBytes(b int) string {
	if b < 1024 {
		return fmt.Sprintf("%d bytes", b)
	}
	if b < 1024*1024 {
		return fmt.Sprintf("%.2f KiB", float64(b)/1024)
	}
	return fmt.Sprintf("%.2f MiB", float64(b)/(1024*1024))
}

func formatDuration(b time.Duration) string {
	if b < time.Microsecond {
		b = b.Round(time.Nanosecond)
	} else if b < time.Millisecond {
		b = b.Round(100 * time.Nanosecond)
	} else {
		b = b.Round(100 * time.Microsecond)
	}

	return b.String()
}
