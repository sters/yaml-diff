package yamldiff

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
)

func Test_diff_Dump(t *testing.T) {
	tests := map[string]map[string]struct {
		d    *diff
		want string
	}{
		"primitive": {
			"simple": {
				d: &diff{
					a:      1,
					b:      1,
					status: diffStatusSame,
				},
				want: "  1\n",
			},
			"diff": {
				d: &diff{
					a:      1,
					b:      2,
					status: diffStatusDiff,
				},
				want: "- 1\n+ 2\n",
			},
			"missing A": {
				d: &diff{
					b:      2,
					status: diffStatus1Missing,
				},
				want: "+ 2\n",
			},
			"missing B": {
				d: &diff{
					a:      1,
					status: diffStatus2Missing,
				},
				want: "- 1\n",
			},
			"diff but string": {
				d: &diff{
					a:      "1",
					b:      "2",
					status: diffStatusDiff,
				},
				want: "- \"1\"\n+ \"2\"\n",
			},
			"diff but bool": {
				d: &diff{
					a:      false,
					b:      true,
					status: diffStatusDiff,
				},
				want: "- false\n+ true\n",
			},
		},
		"map": {
			"simple": {
				d: &diff{
					children: &diffChildren{
						m: diffChildrenMap{
							"foo": {
								a:         "bar",
								b:         "bar",
								status:    diffStatusSame,
								treeLevel: 1,
							},
						},
					},
					status: diffStatusSame,
				},
				want: "  foo: \"bar\"\n",
			},
			"diff": {
				d: &diff{
					children: &diffChildren{
						m: diffChildrenMap{
							"foo": {
								a:         "bar",
								b:         "baz",
								status:    diffStatusDiff,
								diffCount: 1,
								treeLevel: 1,
							},
						},
					},
					status: diffStatusDiff,
				},
				want: "- foo: \"bar\"\n+ foo: \"baz\"\n",
			},
			"missing A": {
				d: &diff{
					children: &diffChildren{
						m: diffChildrenMap{
							"foo": {
								b:         "baz",
								status:    diffStatus1Missing,
								diffCount: 3,
								treeLevel: 1,
							},
						},
					},
					status: diffStatusDiff,
				},
				want: "+ foo: \"baz\"\n",
			},
			"missing B": {
				d: &diff{
					children: &diffChildren{
						m: diffChildrenMap{
							"foo": {
								a:         "bar",
								status:    diffStatus2Missing,
								diffCount: 3,
								treeLevel: 1,
							},
						},
					},
					status: diffStatusDiff,
				},
				want: "- foo: \"bar\"\n",
			},
			"complicated": {
				d: &diff{
					children: &diffChildren{
						m: diffChildrenMap{
							"foo": {
								a: rawTypeMap{
									yaml.MapItem{Key: "bar", Value: "baz"},
									yaml.MapItem{Key: "baz", Value: 1},
									yaml.MapItem{Key: "barr", Value: false},
								},
								b: rawTypeMap{
									yaml.MapItem{Key: "bar", Value: "baz"},
									yaml.MapItem{Key: "baz", Value: rawTypeMap{
										yaml.MapItem{Key: "a", Value: "b"},
									}},
									yaml.MapItem{Key: "bazz", Value: 1},
								},
								children: &diffChildren{
									m: diffChildrenMap{
										"bar": {
											a:         "baz",
											b:         "baz",
											status:    diffStatusSame,
											treeLevel: 2,
										},
										"baz": {
											a: 1,
											b: rawTypeMap{
												yaml.MapItem{Key: "a", Value: "b"},
											},
											status:    diffStatusDiff,
											diffCount: len("map[a:b]"),
											treeLevel: 2,
										},
										"barr": {
											a:         false,
											status:    diffStatus2Missing,
											diffCount: 5,
											treeLevel: 2,
										},
										"bazz": {
											b:         1,
											status:    diffStatus1Missing,
											diffCount: 1,
											treeLevel: 2,
										},
									},
								},
								diffCount: (len("[{a b}]")) + (5) + (1),
								status:    diffStatusDiff,
								treeLevel: 1,
							},
							"bar": {
								a:         1,
								b:         "1",
								diffCount: 0,
								status:    diffStatusDiff,
								treeLevel: 1,
							},
							"baz": {
								a:         "1",
								b:         1,
								diffCount: 0,
								status:    diffStatusDiff,
								treeLevel: 1,
							},
							"zoo": {
								a:         1,
								diffCount: 1,
								status:    diffStatus2Missing,
								treeLevel: 1,
							},
							"boo": {
								b:         1,
								diffCount: 1,
								status:    diffStatus1Missing,
								treeLevel: 1,
							},
						},
					},
					diffCount: ((len("[{a b}]")) + (5) + (1)) + (0) + (0) + (1) + (1),
					status:    diffStatusDiff,
				},
				want: `
  foo:
    bar: "baz"
-   baz: 1
+   baz:
+     a: "b"
-   barr: false
+   bazz: 1
- bar: 1
+ bar: "1"
- baz: "1"
+ baz: 1
- zoo: 1
+ boo: 1`,
			},
		},
		"array": {
			"simple": {
				d: &diff{
					children: &diffChildren{
						a: diffChildrenArray{
							{
								a:         "bar",
								b:         "bar",
								status:    diffStatusSame,
								treeLevel: 1,
							},
						},
					},
					status: diffStatusSame,
				},
				want: "  - \"bar\"\n",
			},
			"diff": {
				d: &diff{
					children: &diffChildren{
						a: diffChildrenArray{
							{
								a:         "bar",
								b:         "baz",
								status:    diffStatusDiff,
								treeLevel: 1,
							},
						},
					},
					status: diffStatusDiff,
				},
				want: "- - \"bar\"\n+ - \"baz\"\n",
			},
			"missing A": {
				d: &diff{
					children: &diffChildren{
						a: diffChildrenArray{
							{
								b:         "baz",
								status:    diffStatus1Missing,
								treeLevel: 1,
							},
						},
					},
					status: diffStatusDiff,
				},
				want: "+ - \"baz\"\n",
			},
			"missing B": {
				d: &diff{
					children: &diffChildren{
						a: diffChildrenArray{
							{
								a:         "bar",
								status:    diffStatus2Missing,
								treeLevel: 1,
							},
						},
					},
					status: diffStatusDiff,
				},
				want: "- - \"bar\"\n",
			},
			"complicated": {
				d: &diff{
					children: &diffChildren{
						a: diffChildrenArray{
							{a: 1, b: 1, status: diffStatusSame, treeLevel: 1},
							{a: 5, b: 5, status: diffStatusSame, treeLevel: 1},
							{
								a: rawTypeArray{2, 3, 4},
								b: rawTypeArray{2},
								children: &diffChildren{
									a: diffChildrenArray{
										{a: 2, b: 2, status: diffStatusSame, treeLevel: 2},
										{a: 3, status: diffStatus2Missing, diffCount: 1, treeLevel: 2},
										{a: 4, status: diffStatus2Missing, diffCount: 1, treeLevel: 2},
									},
								},
								diffCount: 2,
								status:    diffStatusDiff,
								treeLevel: 1,
							},
							{a: 6, status: diffStatus2Missing, diffCount: 1, treeLevel: 1},
						},
					},
					diffCount: 2 + 1,
					status:    diffStatusDiff,
				},
				want: fmt.Sprintf(`
  - 1
  - 5
  %s
    - 2
-   - 3
-   - 4
- - 6`, "- "),
			},
		},
	}
	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			// t.Parallel()

			for n, tc := range tt {
				tc := tc
				t.Run(n, func(t *testing.T) {
					// t.Parallel()

					gotSorted := strings.Split(tc.d.Dump(), "\n")
					sort.SliceStable(gotSorted, func(i, j int) bool { return gotSorted[i] < gotSorted[j] })

					wantSorted := strings.Split(tc.want, "\n")
					sort.SliceStable(wantSorted, func(i, j int) bool { return wantSorted[i] < wantSorted[j] })

					assert.Equal(t, wantSorted, gotSorted)
				})
			}
		})
	}
}
