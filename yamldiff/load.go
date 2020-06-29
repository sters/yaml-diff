package yamldiff

import (
	"log"
	"strings"

	"gopkg.in/yaml.v3"
)

func Load(s string) []interface{} {
	yamls := strings.Split(s, "\n---\n")

	results := make([]interface{}, 0, len(yamls))
	for _, y := range yamls {
		if len(y) < 100 {
			continue
		}

		var out interface{}
		if err := yaml.Unmarshal([]byte(y), &out); err != nil {
			log.Fatalf("f1, %+v", err)
		}
		results = append(results, out)
	}

	return nil
}
