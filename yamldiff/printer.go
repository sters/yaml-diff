package yamldiff

import (
	"fmt"
	"io"
	"strings"
)

const (
	indentString = "  "
)

type sortedChildItem struct {
	k string
	v *diff
}

func indent(level int) string {
	return strings.Repeat(indentString, level)
}

func dumpData(b io.Writer, diffPrefix string, level int, v rawType) {
	if t, ok := tryMap(v); ok {
		dumpMap(b, diffPrefix, level, t)

		return
	}

	if t, ok := tryArray(v); ok {
		dumpArray(b, diffPrefix, level, t)

		return
	}

	dumpPrimitive(b, diffPrefix, level, "", v)
}

func dumpMap(b io.Writer, diffPrefix string, level int, m rawTypeMap) {
	for _, v := range m {
		k, ok := v.Key.(string)
		if !ok {
			k = ""
		}

		dumpMapItem(b, diffPrefix, level, k, v)
	}
}

func dumpArray(b io.Writer, diffPrefix string, level int, m rawTypeArray) {
	for _, v := range m {
		dumpArrayItem(b, diffPrefix, level, v)
	}
}

func dumpArrayItem(b io.Writer, diffPrefix string, level int, v rawType) {
	if t, ok := tryMap(v); ok {
		fmt.Fprintf(b, "%s %s-\n", diffPrefix, indent(level))
		dumpData(b, diffPrefix, level+1, t)

		return
	}

	if t, ok := tryArray(v); ok {
		fmt.Fprintf(b, "%s %s-\n", diffPrefix, indent(level))
		dumpData(b, diffPrefix, level+1, t)

		return
	}

	dumpPrimitive(b, diffPrefix, level, "- ", v)
}

func dumpMapItem(b io.Writer, diffPrefix string, level int, k string, v rawType) {
	if t, ok := tryMap(v); ok {
		fmt.Fprintf(b, "%s %s%s:\n", diffPrefix, indent(level), k)
		dumpData(b, diffPrefix, level+1, t)

		return
	}

	if t, ok := tryArray(v); ok {
		fmt.Fprintf(b, "%s %s%s:\n", diffPrefix, indent(level), k)
		dumpData(b, diffPrefix, level+1, t)

		return
	}

	if t, ok := tryMapItem(v); ok {
		dumpMapItem(b, diffPrefix, level, k, t.Value)

		return
	}

	if v == nil {
		dumpPrimitive(b, diffPrefix, level, fmt.Sprintf("%s:", k), v)

		return
	}

	dumpPrimitive(b, diffPrefix, level, fmt.Sprintf("%s: ", k), v)
}

func dumpPrimitive(b io.Writer, diffPrefix string, level int, somethingPrefix string, v rawType) {
	switch v.(type) {
	case nil, _missingKey:
		fmt.Fprintf(b, "%s %s%s\n", diffPrefix, indent(level), somethingPrefix)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		fmt.Fprintf(b, "%s %s%s%d\n", diffPrefix, indent(level), somethingPrefix, v)
	case float32, float64:
		fmt.Fprintf(b, "%s %s%s%f\n", diffPrefix, indent(level), somethingPrefix, v)
	case string:
		// try escape special characters
		fmt.Fprintf(b, "%s %s%s%#v\n", diffPrefix, indent(level), somethingPrefix, v)
	default:
		fmt.Fprintf(b, "%s %s%s%#v\n", diffPrefix, indent(level), somethingPrefix, v)
	}
}

func (d *diff) dump(b io.Writer, level int) {
	if d.children != nil {
		d.dumpTryArray(b, level)
		d.dumpTryMap(b, level)

		return
	}

	switch d.status {
	case DiffStatusSame:
		dumpData(b, " ", level, d.a)
	case DiffStatusDiff:
		if d.a != nil {
			dumpData(b, "-", level, d.a)
		}
		if d.b != nil {
			dumpData(b, "+", level, d.b)
		}
	case DiffStatus1Missing:
		dumpData(b, "+", level, d.b)
	case DiffStatus2Missing:
		dumpData(b, "-", level, d.a)
	}
}

func (d *diff) dumpTryArray(b io.Writer, level int) {
	if d.children.a == nil {
		return
	}

	for _, v := range d.children.a {
		if v.children != nil && (v.children.a != nil || v.children.m != nil) {
			fmt.Fprintf(b, "  %s-\n", indent(level))
			v.dump(b, level+1)

			continue
		}

		switch v.status {
		case DiffStatusSame:
			dumpArrayItem(b, " ", level, v.a)
		case DiffStatusDiff:
			dumpArrayItem(b, "-", level, v.a)
			dumpArrayItem(b, "+", level, v.b)
		case DiffStatus1Missing:
			dumpArrayItem(b, "+", level, v.b)
		case DiffStatus2Missing:
			dumpArrayItem(b, "-", level, v.a)
		}
	}
}

func (d *diff) dumpTryMap(b io.Writer, level int) {
	if d.children.m == nil {
		return
	}

	sortedChildren := []*sortedChildItem{}
	checked := map[string]struct{}{}

	appendSorted := func(r interface{}) {
		m, ok := tryMap(r)
		if !ok {
			return
		}

		for _, r := range m {
			for k, v := range d.children.m {
				if _, ok := checked[k]; ok {
					continue
				}
				if r.Key != k {
					continue
				}

				sortedChildren = append(sortedChildren, &sortedChildItem{
					k: k,
					v: v,
				})
				checked[k] = struct{}{}
			}
		}
	}

	appendSorted(d.a)
	appendSorted(d.b)

	for k, v := range d.children.m {
		if _, ok := checked[k]; ok {
			continue
		}

		sortedChildren = append(sortedChildren, &sortedChildItem{
			k: k,
			v: v,
		})
	}

	for _, r := range sortedChildren {
		if r.v.children != nil && (r.v.children.a != nil || r.v.children.m != nil) {
			fmt.Fprintf(b, "  %s%s:\n", indent(level), r.k)
			r.v.dump(b, level+1)

			continue
		}

		switch r.v.status {
		case DiffStatusSame:
			dumpMapItem(b, " ", level, r.k, r.v.a)
		case DiffStatusDiff:
			dumpMapItem(b, "-", level, r.k, r.v.a)
			dumpMapItem(b, "+", level, r.k, r.v.b)
		case DiffStatus1Missing:
			dumpMapItem(b, "+", level, r.k, r.v.b)
		case DiffStatus2Missing:
			dumpMapItem(b, "-", level, r.k, r.v.a)
		}
	}
}

func (d *diff) Dump() string {
	var b strings.Builder

	d.dump(&b, d.treeLevel)

	return b.String()
}
