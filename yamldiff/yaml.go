package yamldiff

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

type RawYaml struct {
	raw interface{}
	id  string
}

type RawYamlList []*RawYaml

func newRawYaml(raw interface{}) *RawYaml {
	return &RawYaml{
		raw: raw,
		id:  fmt.Sprintf("%d-%d", time.Now().UnixNano(), randInt()),
	}
}

func randInt() int64 {
	n, _ := rand.Int(rand.Reader, big.NewInt(9223372036854775807))
	return n.Int64()
}
