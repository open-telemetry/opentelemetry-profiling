package main

import (
	"fmt"
	"os"
	"otelprofiling/anonymizer"
	"otelprofiling/parsers"
	"otelprofiling/parsers/collapsed"
	"otelprofiling/parsers/pprof"
	"path/filepath"
	"strings"

	"github.com/jzelinskie/must"
)

func main() {
	filepaths, err := filepath.Glob("./profiles/src/*")

	if err != nil {
		panic(err)
	}

	for _, path := range filepaths {
		var parser parsers.Parser
		ext := filepath.Ext(path)
		switch ext {
		case ".collapsed":
			parser = &collapsed.Parser{}
		case ".pprof":
			parser = &pprof.Parser{}
		default:
			panic("unknown file extension: " + ext)
		}

		inputFile := must.NotError(os.Open(path))
		outPath := strings.Replace(path, "src", "intermediary", 1) + ".intermediary"
		outputFile := must.NotError(os.Create(outPath))
		a := anonymizer.New()
		samples := must.NotError(parser.Parse(inputFile))
		for _, s := range samples {
			if strings.Contains(path, "sensitive") {
				for i := 0; i < len(s.Stacktrace); i++ {
					s.Stacktrace[i] = a.Anonymize(s.Stacktrace[i])
				}
			}
			stacktraceJoined := strings.Join(s.Stacktrace, ";")
			outputFile.Write([]byte(fmt.Sprintf("%s %v\n", stacktraceJoined, s.Value)))
		}
	}
}
