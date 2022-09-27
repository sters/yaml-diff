package yamldiff

import (
	"crypto/rand"
	"fmt"
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

func Load(s string) (RawYamlList, error) {
	yamls := strings.Split(s, "\n---\n")

	results := make(RawYamlList, 0, len(yamls))
	for _, y := range yamls {
		var out interface{}
		if err := yaml.UnmarshalWithOptions([]byte(y), &out, yaml.UseOrderedMap()); err != nil {
			return nil, fmt.Errorf("yamldiff: failed to unmarshal yaml: %+s", err)
		}
		results = append(
			results,
			newRawYaml(out),
		)
	}

	return results, nil
}
