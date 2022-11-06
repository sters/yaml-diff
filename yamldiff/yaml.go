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

func (y *YamlDiff) Status() DiffStatus {
	return y.d.status
}

func (y *YamlDiff) Dump() string {
	return y.d.Dump()
}

type doOptions struct {
	emptyAsNull bool
}

type DoOptionFunc func(o *doOptions)

func EmptyAsNull() DoOptionFunc {
	return func(o *doOptions) {
		o.emptyAsNull = true
	}
}

func Do(rawA RawYamlList, rawB RawYamlList, options ...DoOptionFunc) []*YamlDiff {
	opts := &doOptions{}
	for _, o := range options {
		o(opts)
	}

	r := &runner{
		option: *opts,
		rawA:   rawA,
		rawB:   rawB,
	}

	r.performAllDiff()
	r.findMinimumDiffs()
	r.sortResult()

	return r.diffs
}

type runner struct {
	option doOptions
	rawA   RawYamlList
	rawB   RawYamlList
	diffs  []*YamlDiff
}

func (r *runner) performAllDiff() {
	diffs := make([]*YamlDiff, 0, len(r.rawA)*len(r.rawB))
	for _, a := range r.rawA {
		for _, b := range r.rawB {
			diffs = append(diffs, &YamlDiff{
				d:   r.performDiff(a.raw, b.raw, 0),
				idA: a.id,
				idB: b.id,
			})
		}
	}

	// Make more diffs `A:nil`` and `nil:B`` to find missing entry
	for _, a := range r.rawA {
		diffs = append(diffs, &YamlDiff{
			d:   r.performDiff(a.raw, nil, 0),
			idA: a.id,
			idB: fmt.Sprintf("empty-%d-%d", time.Now().UnixNano(), randInt()),
		})
	}

	for _, b := range r.rawB {
		diffs = append(diffs, &YamlDiff{
			d:   r.performDiff(nil, b.raw, 0),
			idA: fmt.Sprintf("empty-%d-%d", time.Now().UnixNano(), randInt()),
			idB: b.id,
		})
	}

	r.diffs = diffs
}

func (r *runner) findMinimumDiffs() {
	sort.Slice(r.diffs, func(i, j int) bool {
		if r.diffs[i].d.status == DiffStatusSame {
			return true
		}
		return r.diffs[i].d.diffCount < r.diffs[j].d.diffCount
	})

	result := []*YamlDiff{}
	checked := map[string]interface{}{}

	for _, d := range r.diffs {
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

	r.diffs = result
}

func (r *runner) sortResult() {
	result := []*YamlDiff{}
	checked := map[string]interface{}{}

	for _, b := range r.rawA {
		for _, d := range r.diffs {
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

	for _, b := range r.rawB {
		for _, d := range r.diffs {
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

	r.diffs = result
}
