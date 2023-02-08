package main

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"encoding/csv"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/jzelinskie/must"
	"github.com/open-telemetry/opentelemetry-profiling/stateful-benchmarks/pprof"
	"github.com/twmb/murmur3"
	"google.golang.org/protobuf/proto"
)

// TODO: provide a way to set flags
type config struct {
}

func parseConfig() *config {
	cfg := config{}

	flag.Parse()
	return &cfg
}

func readPprof(path string) (*pprof.Profile, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if len(b) > 2 && b[0] == 0x1f && b[1] == 0x8b {
		gzipr, err := gzip.NewReader(bytes.NewReader(b))
		if err != nil {
			return nil, err
		}
		defer gzipr.Close()
		b, err = ioutil.ReadAll(gzipr)
		if err != nil {
			return nil, err
		}
	}

	var p pprof.Profile
	if err = proto.Unmarshal(b, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func writePprof(p *pprof.Profile) ([]byte, error) {
	b := bytes.Buffer{}
	gwriter := gzip.NewWriter(&b)
	res, err := proto.Marshal(p)
	if err != nil {
		return nil, err
	}

	_, err = gwriter.Write(res)
	if err != nil {
		return nil, err
	}
	gwriter.Close()

	return b.Bytes(), nil
}

const murmurSeed = 6231912

func utf8CompatibleHash(str string) string {
	buf := make([]byte, 16)
	a, b := murmur3.SeedSum128(murmurSeed, murmurSeed, []byte(str))

	binary.LittleEndian.PutUint64(buf, a)
	binary.LittleEndian.PutUint64(buf[8:], b)

	return hex.EncodeToString(buf[:8])
}

func perc(a, b int) string {
	sign := ""
	if a > b {
		sign = "+"
	}
	return fmt.Sprintf("%s%.2f%%", sign, float64(a-b)/float64(b)*100)
}

func writePprofOrPanic(p *pprof.Profile) []byte {
	res, err := writePprof(p)
	if err != nil {
		panic(err)
	}
	return res
}

func main() {
	parseConfig()
	files := flag.Args()

	allSymbolsMap := make(map[string]bool)

	writer := csv.NewWriter(os.Stdout)
	writer.Comma = '\t'

	writer.Write([]string{
		"noSymbols",
		"symbolsAsStrings",
		"symbolsAsStrings percent difference",
		"hashedSymbols",
		"hashedSymbols percent difference",
	})

	for _, filepath := range files {
		prof, err := readPprof(filepath)
		if err != nil {
			continue
		}

		symbolsAsStrings := must.NotError(writePprof(prof))

		for i, s := range prof.StringTable {
			allSymbolsMap[s] = true
			prof.StringTable[i] = utf8CompatibleHash(s)
		}
		hashedSymbols := must.NotError(writePprof(prof))

		prof.StringTable = []string{}
		noSymbols := must.NotError(writePprof(prof))

		writer.Write([]string{
			strconv.Itoa(len(noSymbols)),
			strconv.Itoa(len(symbolsAsStrings)),
			perc(len(symbolsAsStrings), len(noSymbols)),
			strconv.Itoa(len(hashedSymbols)),
			perc(len(hashedSymbols), len(noSymbols)),
		})
	}
	writer.Flush()
	separateSymbols := pprof.Profile{}
	allSymbols := make([]string, 0, len(allSymbolsMap))
	for s := range allSymbolsMap {
		allSymbols = append(allSymbols, s)
	}
	separateSymbols.StringTable = allSymbols

	separateSymbolsBytes, err := writePprof(&separateSymbols)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nseparateSymbols: %d\n", len(separateSymbolsBytes))
}
