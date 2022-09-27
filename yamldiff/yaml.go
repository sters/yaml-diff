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
	Raw interface{}
	id  string
}

type RawYamlList []*RawYaml

type Diffs []*diff

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
			return nil, fmt.Errorf("yamldiff: failed to unmarshal yaml: %w", err)
		}

		results = append(
			results,
			newRawYaml(out),
		)
	}

	return results, nil
}

type diffData struct {
	d   *diff
	idA string
	idB string
}

func Do(rawA RawYamlList, rawB RawYamlList) []*diff {
	return sortResult(rawA, findMinimumDiffs(performAllDiff(rawA, rawB)))
}

func performAllDiff(rawA RawYamlList, rawB RawYamlList) []*diffData {
	diffs := make([]*diffData, 0, len(rawA)*len(rawB))
	for _, a := range rawA {
		for _, b := range rawB {
			diffs = append(diffs, &diffData{
				d:   performDiff(a.Raw, b.Raw, 0),
				idA: a.id,
				idB: b.id,
			})
		}
	}

	// Make more diffs `A:nil`` and `nil:B`` to find missing entry
	for _, a := range rawA {
		diffs = append(diffs, &diffData{
			d:   performDiff(a.Raw, nil, 0),
			idA: a.id,
			idB: fmt.Sprintf("empty-%d-%d", time.Now().UnixNano(), randInt()),
		})
	}

	for _, b := range rawB {
		diffs = append(diffs, &diffData{
			d:   performDiff(nil, b.Raw, 0),
			idA: fmt.Sprintf("empty-%d-%d", time.Now().UnixNano(), randInt()),
			idB: b.id,
		})
	}

	return diffs
}

func findMinimumDiffs(diffs []*diffData) []*diffData {
	sort.Slice(diffs, func(i, j int) bool {
		if diffs[i].d.status == DiffStatusSame {
			return true
		}
		return diffs[i].d.diffCount < diffs[j].d.diffCount
	})

	result := []*diffData{}
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

	// // add missing diffs for A
	// for _, d := range diffs {
	// 	if _, ok := checked[d.idA]; ok {
	// 		continue
	// 	}

	// 	dd := *(d.d)

	// 	result = append(result, &diffData{
	// 		d:   fillDummyDiffStatus(DiffStatus2Missing, &dd),
	// 		idA: d.idA,
	// 	})
	// 	checked[d.idA] = struct{}{}
	// }

	// // add missing diffs for B
	// for _, d := range diffs {
	// 	if _, ok := checked[d.idB]; ok {
	// 		continue
	// 	}

	// 	dd := *(d.d)

	// 	result = append(result, &diffData{
	// 		d:   fillDummyDiffStatus(DiffStatus1Missing, &dd),
	// 		idB: d.idB,
	// 	})
	// 	checked[d.idB] = struct{}{}
	// }

	return result
}

func fillDummyDiffStatus(status DiffStatus, target *diff) *diff {
	if target.children != nil {
		if target.children.a != nil {
			for k, a := range target.children.a {
				target.children.a[k] = fillDummyDiffStatus(status, a)
			}
		}

		if target.children.m != nil {
			for k, a := range target.children.m {
				target.children.m[k] = fillDummyDiffStatus(status, a)
			}
		}
	}

	target.status = status

	return target
}

func sortResult(base RawYamlList, diffs []*diffData) []*diff {
	result := []*diff{}
	checked := map[string]interface{}{}

	for _, b := range base {
		for _, d := range diffs {
			if b.id != d.idA {
				continue
			}

			result = append(result, d.d)
			checked[d.idA] = struct{}{}
			checked[d.idB] = struct{}{}
		}
	}

	for _, d := range diffs {
		if _, ok := checked[d.idA]; ok {
			continue
		}
		if _, ok := checked[d.idB]; ok {
			continue
		}

		result = append(result, d.d)
	}

	return result
}
