package yamldiff

import (
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
)

func Test_performDiff(t *testing.T) {
	tests := map[string]map[string]struct {
		a    rawType
		b    rawType
		want *diff
	}{
		"primitive": {
			"int ok": {
				a: 1,
				b: 1,
				want: &diff{
					status: diffStatusSame,
				},
			},
			"int diff": {
				a: 1,
				b: 0,
				want: &diff{
					status:    diffStatusDiff,
					diffCount: 1,
				},
			},
			"int missing a": {
				b: 1,
				want: &diff{
					status:    diffStatus1Missing,
					diffCount: 1,
				},
			},
			"int missing b": {
				a: 11,
				want: &diff{
					status:    diffStatus2Missing,
					diffCount: 2,
				},
			},
			"string ok": {
				a: "1",
				b: "1",
				want: &diff{
					status: diffStatusSame,
				},
			},
			"string diff": {
				a: "1",
				b: "0",
				want: &diff{
					status:    diffStatusDiff,
					diffCount: 1,
				},
			},
			"string missing a": {
				b: "1",
				want: &diff{
					status:    diffStatus1Missing,
					diffCount: 1,
				},
			},
			"string missing b": {
				a: "11",
				want: &diff{
					status:    diffStatus2Missing,
					diffCount: 2,
				},
			},
			"int vs string": {
				a: "1",
				b: 1,
				want: &diff{
					status:    diffStatusDiff,
					diffCount: 0, // because it's only diff on type
				},
			},
			"int vs string 2": {
				a: "1",
				b: 0,
				want: &diff{
					status:    diffStatusDiff,
					diffCount: 1,
				},
			},
			"int vs float": {
				a: 1,
				b: 0.5,
				want: &diff{
					status:    diffStatusDiff,
					diffCount: 3,
				},
			},
			"float vs string": {
				a: "1",
				b: 0.5,
				want: &diff{
					status:    diffStatusDiff,
					diffCount: 3,
				},
			},
			"int vs bool": {
				a: 1,
				b: false,
				want: &diff{
					status:    diffStatusDiff,
					diffCount: 5,
				},
			},
		},
		"array": {
			"simple same": {
				a: rawTypeArray{1, 2, 3},
				b: rawTypeArray{1, 2, 3},
				want: &diff{
					children: &diffChildren{
						a: diffChildrenArray{
							{a: 1, b: 1, status: diffStatusSame, treeLevel: 1},
							{a: 2, b: 2, status: diffStatusSame, treeLevel: 1},
							{a: 3, b: 3, status: diffStatusSame, treeLevel: 1},
						},
					},
					status: diffStatusSame,
				},
			},
			"simple same different order": {
				a: rawTypeArray{1, 2, 3},
				b: rawTypeArray{3, 1, 2},
				want: &diff{
					children: &diffChildren{
						a: diffChildrenArray{
							{a: 1, b: 1, status: diffStatusSame, treeLevel: 1},
							{a: 2, b: 2, status: diffStatusSame, treeLevel: 1},
							{a: 3, b: 3, status: diffStatusSame, treeLevel: 1},
						},
					},
					status: diffStatusSame,
				},
			},
			"missing in A": {
				a: rawTypeArray{1, 2},
				b: rawTypeArray{1, 2, 3},
				want: &diff{
					children: &diffChildren{
						a: diffChildrenArray{
							{a: 1, b: 1, status: diffStatusSame, treeLevel: 1},
							{a: 2, b: 2, status: diffStatusSame, treeLevel: 1},
							{b: 3, status: diffStatus1Missing, diffCount: 1, treeLevel: 1},
						},
					},
					diffCount: 1,
					status:    diffStatusDiff,
				},
			},
			"missing in B": {
				a: rawTypeArray{1, 2, 3},
				b: rawTypeArray{1, 3},
				want: &diff{
					children: &diffChildren{
						a: diffChildrenArray{
							{a: 1, b: 1, status: diffStatusSame, treeLevel: 1},
							{a: 3, b: 3, status: diffStatusSame, treeLevel: 1},
							{a: 2, status: diffStatus2Missing, diffCount: 1, treeLevel: 1}, // missing is added by last
						},
					},
					diffCount: 1,
					status:    diffStatusDiff,
				},
			},
			"missing in A and B": {
				a: rawTypeArray{1, 2, 3},
				b: rawTypeArray{1, 3, 4},
				want: &diff{
					children: &diffChildren{
						a: diffChildrenArray{
							{a: 1, b: 1, status: diffStatusSame, treeLevel: 1},
							{a: 3, b: 3, status: diffStatusSame, treeLevel: 1},
							{a: 2, b: 4, status: diffStatusDiff, diffCount: 1, treeLevel: 1}, // because can't find missing, it's diff.
						},
					},
					diffCount: 1,
					status:    diffStatusDiff,
				},
			},
			"complicated": {
				a: rawTypeArray{1, rawTypeArray{2, 3, 4}, 5, 6},
				b: rawTypeArray{1, 5, rawTypeArray{2}},
				want: &diff{
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
			},
		},
		"map": {
			"simple same": {
				a: rawTypeMap{
					yaml.MapItem{Key: "foo", Value: "bar"},
				},
				b: rawTypeMap{
					yaml.MapItem{Key: "foo", Value: "bar"},
				},
				want: &diff{
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
			},
			"simple diff": {
				a: rawTypeMap{
					yaml.MapItem{Key: "foo", Value: "bar"},
				},
				b: rawTypeMap{
					yaml.MapItem{Key: "foo", Value: "baz"},
				},
				want: &diff{
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
					diffCount: 1,
					status:    diffStatusDiff,
				},
			},
			"simple diff type": {
				a: rawTypeMap{
					yaml.MapItem{Key: "foo", Value: "bar"},
				},
				b: rawTypeMap{
					yaml.MapItem{Key: "foo", Value: 1},
				},
				want: &diff{
					children: &diffChildren{
						m: diffChildrenMap{
							"foo": {
								a:         "bar",
								b:         1,
								diffCount: 3,
								status:    diffStatusDiff,
								treeLevel: 1,
							},
						},
					},
					diffCount: 3,
					status:    diffStatusDiff,
				},
			},
			"simple diff type map and primitive": {
				a: rawTypeMap{
					yaml.MapItem{Key: "foo", Value: "bar"},
				},
				b: rawTypeMap{
					yaml.MapItem{
						Key: "foo",
						Value: yaml.MapItem{
							Key:   "bar",
							Value: "baz",
						},
					},
				},
				want: &diff{
					children: &diffChildren{
						m: diffChildrenMap{
							"foo": {
								a: "bar",
								b: yaml.MapItem{
									Key:   "bar",
									Value: "baz",
								},
								diffCount: len("{bar baz}"), // as primitive diff
								status:    diffStatusDiff,
								treeLevel: 1,
							},
						},
					},

					diffCount: len("{bar baz}"), // from child
					status:    diffStatusDiff,
				},
			},
			"complicated": {
				a: rawTypeMap{
					yaml.MapItem{
						Key: "foo",
						Value: rawTypeMap{
							yaml.MapItem{Key: "bar", Value: "baz"},
							yaml.MapItem{Key: "baz", Value: 1},
							yaml.MapItem{Key: "barr", Value: false},
						},
					},
					yaml.MapItem{Key: "bar", Value: 1},
					yaml.MapItem{Key: "baz", Value: "1"},
					yaml.MapItem{Key: "zoo", Value: 1},
				},
				b: rawTypeMap{
					yaml.MapItem{
						Key: "foo",
						Value: rawTypeMap{
							yaml.MapItem{Key: "bar", Value: "baz"},
							yaml.MapItem{Key: "baz", Value: rawTypeMap{
								yaml.MapItem{Key: "a", Value: "b"},
							}},
							yaml.MapItem{Key: "bazz", Value: 1},
						},
					},
					yaml.MapItem{Key: "bar", Value: "1"},
					yaml.MapItem{Key: "baz", Value: 1},
					yaml.MapItem{Key: "boo", Value: 1},
				},
				want: &diff{
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
											diffCount: len("[{a b}]"),
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

					got := performDiff(tc.a, tc.b, 0)

					tc.want.a = tc.a
					tc.want.b = tc.b
					assert.Equal(t, tc.want, got)
				})
			}
		})
	}
}
