package yamldiff

import (
	"log"
	"strings"

	"github.com/goccy/go-yaml"
)

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
