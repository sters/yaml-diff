package yamldiff

import (
	"fmt"
	"sort"
	"strings"

	"github.com/google/go-cmp/cmp"
)

type Diff struct {
	Diff      string
	difflines int

	yaml1 *RawYaml
	yaml2 *RawYaml
}

type Diffs []*Diff

func Do(list1 RawYamlList, list2 RawYamlList) Diffs {
	var result Diffs

	checked := map[string]struct{}{} // RawYaml.id => struct{}

	for _, yaml1 := range list1 {

		diffs := make([]*Diff, 0, len(list2))

		for _, yaml2 := range list2 {
			if _, ok := checked[yaml2.id]; ok {
				continue
			}

			s := &Diff{
				Diff:  cmp.Diff(yaml1.raw, yaml2.raw),
				yaml1: yaml1,
				yaml2: yaml2,
			}

			if len(strings.TrimSpace(s.Diff)) < 1 {
				content1 := fmt.Sprintf("%+v", yaml1.raw)
				content1trimlen := 100
				if len(content1) < content1trimlen {
					content1trimlen = len(content1)
				}

				content2 := fmt.Sprintf("%+v", yaml2.raw)
				content2trimlen := 100
				if len(content2) < content2trimlen {
					content2trimlen = len(content2)
				}

				s.Diff = fmt.Sprintf(
					"Same Content: %s..., %s...",
					content1[:content1trimlen],
					content2[:content2trimlen],
				)
			}

			for _, str := range strings.Split(s.Diff, "\n") {
				trimmedstr := strings.TrimSpace(str)
				if strings.HasPrefix(trimmedstr, "+") || strings.HasPrefix(str, "-") {
					s.difflines++
				}
			}

			diffs = append(diffs, s)
		}

		if len(diffs) == 0 {
			continue
		}

		sort.Slice(diffs, func(i, j int) bool {
			return diffs[i].difflines < diffs[j].difflines
		})

		result = append(result, diffs[0])
		checked[diffs[0].yaml1.id] = struct{}{}
		checked[diffs[0].yaml2.id] = struct{}{}
	}

	// check the unmarked items in list1
	for _, yaml1 := range list1 {
		if _, ok := checked[yaml1.id]; ok {
			continue
		}

		checked[yaml1.id] = struct{}{}

		content := fmt.Sprintf("%+v", yaml1.raw)
		trimlen := 100
		if len(content) < trimlen {
			trimlen = len(content)
		}

		result = append(
			result,
			&Diff{
				Diff: fmt.Sprintf(
					"Not found on another one: %s...",
					content[:trimlen],
				),
				yaml1: yaml1,
			},
		)
	}
	for _, yaml2 := range list2 {
		if _, ok := checked[yaml2.id]; ok {
			continue
		}

		checked[yaml2.id] = struct{}{}

		content := fmt.Sprintf("%+v", yaml2.raw)
		trimlen := 100
		if len(content) < trimlen {
			trimlen = len(content)
		}

		result = append(
			result,
			&Diff{
				Diff: fmt.Sprintf(
					"Not found on another one: %s...",
					content[:trimlen],
				),
				yaml2: yaml2,
			},
		)
	}

	return result
}
