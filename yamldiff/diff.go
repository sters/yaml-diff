package yamldiff

import (
	"fmt"
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
			s := Diff{n: n, difflines: 0}

			if _, ok := marker[n]; ok {
				continue
			}

			s.Diff = cmp.Diff(y1, y2)

			if len(strings.TrimSpace(s.Diff)) < 1 {
				s.Diff = fmt.Sprintf(
					"Same Content: %s..., %s...",
					fmt.Sprintf("%+v", y1)[0:100],
					fmt.Sprintf("%+v", y2)[0:100],
				)
			}

			for _, str := range strings.Split(s.Diff, "\n") {
				trimmedstr := strings.TrimSpace(str)
				if strings.HasPrefix(trimmedstr, "+") || strings.HasPrefix(str, "-") {
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
