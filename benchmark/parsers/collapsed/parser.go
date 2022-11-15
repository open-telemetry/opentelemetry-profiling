package collapsed

import (
	"bufio"
	"io"
	"otelprofiling/parsers"
	"strconv"
	"strings"

	"github.com/jzelinskie/must"
)

type Parser struct{}

// TODO(petethepig): implement support for timestamps and labels, e.g:
// slow;code 100
// fast;code 20
// foo;bar 20 region=us-east-1,container=abc123
// foo;bar 20 1662426304123
// foo;bar 20 region=us-east-1,container=abc123 1662426304123

// more info on that here: https://docs.google.com/document/d/1Ba1F1PjRtW4AZ6qurUrBjw6bz5NiKqIs2WQQ2N2FR6A/edit#heading=h.nhcufnxt0tol
func (_ *Parser) Parse(r io.Reader) ([]parsers.Sample, error) {
	res := make([]parsers.Sample, 0)

	scanner := bufio.NewScanner(r)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		str := scanner.Text()
		lastSpace := strings.LastIndex(str, " ")
		stacktrace := strings.Split(str[:lastSpace], ";")
		value := must.NotError(strconv.Atoi(str[lastSpace+1:]))
		res = append(res, parsers.Sample{Stacktrace: stacktrace, Value: value})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
