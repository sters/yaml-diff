package yamldiff

import (
	"sort"
	"strings"

	"github.com/google/go-cmp/cmp"
)

type DiffStatus int

const (
	DiffStatusExists   DiffStatus = 1
	DiffStatusSame     DiffStatus = 2
	DiffStatus1Missing DiffStatus = 3
	DiffStatus2Missing DiffStatus = 4
)

type Diff struct {
	Diff      string
	difflines int

	Yaml1Struct *RawYaml
	Yaml2Struct *RawYaml

	Status DiffStatus
}

type Diffs []*Diff

func Do(list1 RawYamlList, list2 RawYamlList) Diffs {
	var result Diffs

	checked := map[string]struct{}{} // RawYaml.id => struct{}

	matchFuncs := []func([]*Diff) *Diff{
		func(diffs []*Diff) *Diff {
			for _, d := range diffs {
				if d.Status == DiffStatusSame {
					return d
				}
			}
			return nil
		},
		func(diffs []*Diff) *Diff {
			sort.Slice(diffs, func(i, j int) bool {
				return diffs[i].difflines < diffs[j].difflines
			})

			return diffs[0]
		},
	}

	for _, matchFunc := range matchFuncs {
		for _, yaml1 := range list1 {
			if _, ok := checked[yaml1.id]; ok {
				continue
			}

			diffs := make([]*Diff, 0, len(list2))

			for _, yaml2 := range list2 {
				if _, ok := checked[yaml2.id]; ok {
					continue
				}

				s := &Diff{
					Diff:        cmp.Diff(yaml1.Raw, yaml2.Raw),
					Yaml1Struct: yaml1,
					Yaml2Struct: yaml2,
					Status:      DiffStatusExists,
				}

				if len(strings.TrimSpace(s.Diff)) < 1 {
					s.Status = DiffStatusSame
					s.Diff = createSameFormat(yaml1, s.Status)
				} else {
					for _, str := range strings.Split(s.Diff, "\n") {
						trimmedstr := strings.TrimSpace(str)
						if strings.HasPrefix(trimmedstr, "+") || strings.HasPrefix(str, "-") {
							s.difflines++
						}
					}
				}

				diffs = append(diffs, s)
			}

			if len(diffs) == 0 {
				continue
			}

			d := matchFunc(diffs)
			if d == nil {
				continue
			}

			result = append(result, d)
			checked[d.Yaml1Struct.id] = struct{}{}
			checked[d.Yaml2Struct.id] = struct{}{}
		}
	}

	// check the unmarked items in list1
	for _, Yaml1Struct := range list1 {
		if _, ok := checked[Yaml1Struct.id]; ok {
			continue
		}

		result = append(
			result,
			&Diff{
				Diff:        "",
				Yaml1Struct: Yaml1Struct,
				Status:      DiffStatus2Missing,
			},
		)
	}

	for _, yaml2 := range list2 {
		if _, ok := checked[yaml2.id]; ok {
			continue
		}

		result = append(
			result,
			&Diff{
				Diff:        "",
				Yaml2Struct: yaml2,
				Status:      DiffStatus1Missing,
			},
		)
	}

	return result
}

func createSameFormat(y *RawYaml, status DiffStatus) string {
	result := strings.Builder{}

	prefix := ""
	switch status {
	case DiffStatusSame:
		prefix = "  "
	case DiffStatus1Missing:
		prefix = "+ "
	case DiffStatus2Missing:
		prefix = "- "
	}

	diff := cmp.Diff(y.Raw, interface{}(1))

	for _, str := range strings.Split(diff, "\n") {
		if !strings.HasPrefix(str, "-") {
			continue
		}

		str = strings.TrimSpace(str)
		str = strings.Replace(str, "- \t", "", 1)

		result.WriteString(prefix)
		result.WriteString(str)
		result.WriteRune('\n')
	}

	return result.String()
}
