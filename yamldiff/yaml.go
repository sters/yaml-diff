package yamldiff

import (
	"crypto/rand"
	"fmt"
	"log"
	"math"
	"math/big"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
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

func Load(s string) RawYamlList {
	yamls := strings.Split(s, "\n---\n")

	results := make(RawYamlList, 0, len(yamls))
	for _, y := range yamls {
		var out interface{}
		if err := yaml.Unmarshal([]byte(y), &out); err != nil {
			log.Fatalf("f1, %+v", err)
		}
		results = append(
			results,
			newRawYaml(out),
		)
	}

	return results
}
