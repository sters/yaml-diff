package yamldiff

import (
	"sort"

	"github.com/google/go-cmp/cmp"
)

type Diff struct {
	n    int
	Diff string
}

type Diffs []Diff

func Do(yamls1 []interface{}, yamls2 []interface{}) []Diffs {
	var diffs []Diffs

	marker := map[int]bool{}
	for _, y1 := range yamls1 {

		d := make([]Diff, 0, len(yamls2))
		for n, y2 := range yamls2 {
			s := Diff{n: n}

			if _, ok := marker[n]; ok {
				s.Diff = cmp.Diff(y1, "fake")
				d = append(d, s)
				continue
			}

			s.Diff = cmp.Diff(y1, y2)
			d = append(d, s)
		}

		sort.Slice(d, func(i, j int) bool {
			return len(d[i].Diff) < len(d[j].Diff)
		})

		diffs = append(diffs, d)
		marker[d[0].n] = true
	}

	return diffs
}
