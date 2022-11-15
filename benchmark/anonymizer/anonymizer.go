package anonymizer

import (
	"encoding/base64"
	"math/rand"
	"strings"

	"github.com/jzelinskie/must"
)

type Anonymizer struct {
	r     *rand.Rand
	cache map[string]string
}

func New() *Anonymizer {
	return &Anonymizer{
		r:     rand.New(rand.NewSource(23061912)),
		cache: map[string]string{},
	}
}

func (a *Anonymizer) Anonymize(str string) string {
	if _, ok := a.cache[str]; !ok {
		src := make([]byte, len(str))
		must.NotError(a.r.Read(src))
		a.cache[str] = strings.ReplaceAll(base64.RawStdEncoding.EncodeToString(src)[:len(str)], " ", "")
	}

	return a.cache[str]
}
