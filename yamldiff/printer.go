package yamldiff

import (
	"fmt"
	"io"
	"strings"
)

const (
	indentString = "  "
)

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
		fmt.Fprintf(b, "%s %s- \n", diffPrefix, indent(level))
		dumpData(b, diffPrefix, level+1, t)

		return
	}

	if t, ok := tryArray(v); ok {
		fmt.Fprintf(b, "%s %s- \n", diffPrefix, indent(level))
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

	dumpPrimitive(b, diffPrefix, level, fmt.Sprintf("%s: ", k), v)
}

func dumpPrimitive(b io.Writer, diffPrefix string, level int, somethingPrefix string, v rawType) {
	switch v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		fmt.Fprintf(b, "%s %s%s%d\n", diffPrefix, indent(level), somethingPrefix, v)
	case float32, float64:
		fmt.Fprintf(b, "%s %s%s%f\n", diffPrefix, indent(level), somethingPrefix, v)
	case string:
		fmt.Fprintf(b, "%s %s%s\"%s\"\n", diffPrefix, indent(level), somethingPrefix, v)
	default:
		fmt.Fprintf(b, "%s %s%s%#v\n", diffPrefix, indent(level), somethingPrefix, v)
	}
}

func (d *diff) dump(b io.Writer, level int) {
	if d.children != nil {
		if d.children.a != nil {
			for _, v := range d.children.a {
				if v.children != nil && (v.children.a != nil || v.children.m != nil) {
					fmt.Fprintf(b, "  %s- \n", indent(level))
					v.dump(b, level+1)

					continue
				}

				switch v.status {
				case diffStatusSame:
					dumpArrayItem(b, " ", level, v.a)
				case diffStatusDiff:
					dumpArrayItem(b, "-", level, v.a)
					dumpArrayItem(b, "+", level, v.b)
				case diffStatus1Missing:
					dumpArrayItem(b, "+", level, v.b)
				case diffStatus2Missing:
					dumpArrayItem(b, "-", level, v.a)
				}
			}
		}

		if d.children.m != nil {
			for k, v := range d.children.m {
				if v.children != nil && (v.children.a != nil || v.children.m != nil) {
					fmt.Fprintf(b, "  %s%s:\n", indent(level), k)
					v.dump(b, level+1)

					continue
				}

				switch v.status {
				case diffStatusSame:
					dumpMapItem(b, " ", level, k, v.a)
				case diffStatusDiff:
					dumpMapItem(b, "-", level, k, v.a)
					dumpMapItem(b, "+", level, k, v.b)
				case diffStatus1Missing:
					dumpMapItem(b, "+", level, k, v.b)
				case diffStatus2Missing:
					dumpMapItem(b, "-", level, k, v.a)
				}
			}
		}

		return
	}

	switch d.status {
	case diffStatusSame:
		dumpData(b, " ", level, d.a)
	case diffStatusDiff:
		dumpData(b, "-", level, d.a)
		dumpData(b, "+", level, d.b)
	case diffStatus1Missing:
		dumpData(b, "+", level, d.b)
	case diffStatus2Missing:
		dumpData(b, "-", level, d.a)
	}
}

func (d *diff) Dump() string {
	var b strings.Builder

	d.dump(&b, d.treeLevel)

	return b.String()
}
