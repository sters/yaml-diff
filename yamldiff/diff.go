package yamldiff

import (
	"sort"
	"strings"

	"github.com/google/go-cmp/cmp"
)

type Diff struct {
	n         int
	Diff      string
	difflines int
}

type Diffs []Diff

func Do(yamls1 []interface{}, yamls2 []interface{}) Diffs {
	var diffs Diffs

	marker := map[int]bool{}
	for _, y1 := range yamls1 {

		d := make([]Diff, 0, len(yamls2))
		for n, y2 := range yamls2 {
			s := Diff{n: n}

			if _, ok := marker[n]; ok {
				continue
			}

			s.Diff = cmp.Diff(y1, y2)
			for _, str := range strings.Split(s.Diff, "\n") {
				if strings.HasPrefix("+", str) || strings.HasPrefix("-", str) {
					s.difflines++
				}
			}

			d = append(d, s)
		}

		sort.Slice(d, func(i, j int) bool {
			return d[i].difflines < d[j].difflines
		})

		diffs = append(diffs, d[0])
		marker[d[0].n] = true
	}

	return diffs
}
