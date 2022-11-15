package anonymizer

import (
	"encoding/base64"
	"math/rand"

	"github.com/jzelinskie/must"
)

type Anonymizer struct {
	r *rand.Rand
}

func New() *Anonymizer {
	return &Anonymizer{
		r: rand.New(rand.NewSource(23061912)),
	}
}

func (a *Anonymizer) Anonymize(str string) string {
	src := make([]byte, len(str))
	must.NotError(a.r.Read(src))

	return base64.StdEncoding.EncodeToString(src)[:len(str)]
}
