package yamldiff

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"time"
)

type RawYaml struct {
	Raw interface{}
	id  string
}

type RawYamlList []*RawYaml

func newRawYaml(raw interface{}) *RawYaml {
	return &RawYaml{
		Raw: raw,
		id:  fmt.Sprintf("%d-%d", time.Now().UnixNano(), randInt()),
	}
}

func randInt() int64 {
	n, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	return n.Int64()
}
