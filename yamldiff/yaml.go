package yamldiff

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
)

type RawYaml struct {
	raw interface{}
	id  string
}

type RawYamlList []*RawYaml

type Diffs []*diff

func newRawYaml(raw interface{}) *RawYaml {
	return &RawYaml{
		raw: raw,
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
			return nil, fmt.Errorf("yamldiff: failed to unmarshal yaml: %w", err)
		}

		results = append(
			results,
			newRawYaml(out),
		)
	}

	return results, nil
}

type YamlDiff struct {
	d   *diff
	idA string
	idB string
}

func (y *YamlDiff) Dump() string {
	return y.d.Dump()
}

func Do(rawA RawYamlList, rawB RawYamlList) []*YamlDiff {
	return sortResult(rawA, rawB, findMinimumDiffs(performAllDiff(rawA, rawB)))
}

func performAllDiff(rawA RawYamlList, rawB RawYamlList) []*YamlDiff {
	diffs := make([]*YamlDiff, 0, len(rawA)*len(rawB))
	for _, a := range rawA {
		for _, b := range rawB {
			diffs = append(diffs, &YamlDiff{
				d:   performDiff(a.raw, b.raw, 0),
				idA: a.id,
				idB: b.id,
			})
		}
	}

	// Make more diffs `A:nil`` and `nil:B`` to find missing entry
	for _, a := range rawA {
		diffs = append(diffs, &YamlDiff{
			d:   performDiff(a.raw, nil, 0),
			idA: a.id,
			idB: fmt.Sprintf("empty-%d-%d", time.Now().UnixNano(), randInt()),
		})
	}

	for _, b := range rawB {
		diffs = append(diffs, &YamlDiff{
			d:   performDiff(nil, b.raw, 0),
			idA: fmt.Sprintf("empty-%d-%d", time.Now().UnixNano(), randInt()),
			idB: b.id,
		})
	}

	return diffs
}

func findMinimumDiffs(diffs []*YamlDiff) []*YamlDiff {
	sort.Slice(diffs, func(i, j int) bool {
		if diffs[i].d.status == diffStatusSame {
			return true
		}
		return diffs[i].d.diffCount < diffs[j].d.diffCount
	})

	result := []*YamlDiff{}
	checked := map[string]interface{}{}

	for _, d := range diffs {
		if _, ok := checked[d.idA]; ok {
			continue
		}
		if _, ok := checked[d.idB]; ok {
			continue
		}

		result = append(result, d)

		checked[d.idA] = struct{}{}
		checked[d.idB] = struct{}{}
	}

	// Even if missing entries in A or B, it should be covered by A:nil or nil:B case.

	return result
}

func sortResult(rawA RawYamlList, rawB RawYamlList, diffs []*YamlDiff) []*YamlDiff {
	result := []*YamlDiff{}
	checked := map[string]interface{}{}

	for _, b := range rawA {
		for _, d := range diffs {
			if b.id != d.idA {
				continue
			}
			if _, ok := checked[d.idA]; ok {
				continue
			}
			if _, ok := checked[d.idB]; ok {
				continue
			}

			result = append(result, d)
			checked[d.idA] = struct{}{}
			checked[d.idB] = struct{}{}
		}
	}

	for _, b := range rawB {
		for _, d := range diffs {
			if b.id != d.idB {
				continue
			}
			if _, ok := checked[d.idA]; ok {
				continue
			}
			if _, ok := checked[d.idB]; ok {
				continue
			}

			result = append(result, d)
			checked[d.idA] = struct{}{}
			checked[d.idB] = struct{}{}
		}
	}

	return result
}
