package yamldiff

import (
	"testing"

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
					status: DiffStatusSame,
				},
			},
			"int diff": {
				a: 1,
				b: 0,
				want: &diff{
					status:    DiffStatusDiff,
					diffCount: 1,
				},
			},
			"int missing a": {
				b: 1,
				want: &diff{
					status:    DiffStatus1Missing,
					diffCount: 1,
				},
			},
			"int missing b": {
				a: 11,
				want: &diff{
					status:    DiffStatus2Missing,
					diffCount: 2,
				},
			},
			"string ok": {
				a: "1",
				b: "1",
				want: &diff{
					status: DiffStatusSame,
				},
			},
			"string diff": {
				a: "1",
				b: "0",
				want: &diff{
					status:    DiffStatusDiff,
					diffCount: 1,
				},
			},
			"string missing a": {
				b: "1",
				want: &diff{
					status:    DiffStatus1Missing,
					diffCount: 1,
				},
			},
			"string missing b": {
				a: "11",
				want: &diff{
					status:    DiffStatus2Missing,
					diffCount: 2,
				},
			},
			"int vs string": {
				a: "1",
				b: 1,
				want: &diff{
					status:    DiffStatusDiff,
					diffCount: 0, // because it's only diff on type
				},
			},
			"int vs string 2": {
				a: "1",
				b: 0,
				want: &diff{
					status:    DiffStatusDiff,
					diffCount: 1,
				},
			},
			"int vs float": {
				a: 1,
				b: 0.5,
				want: &diff{
					status:    DiffStatusDiff,
					diffCount: 3,
				},
			},
			"float vs string": {
				a: "1",
				b: 0.5,
				want: &diff{
					status:    DiffStatusDiff,
					diffCount: 3,
				},
			},
			"int vs bool": {
				a: 1,
				b: false,
				want: &diff{
					status:    DiffStatusDiff,
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
							{a: 1, b: 1, status: DiffStatusSame, treeLevel: 1},
							{a: 2, b: 2, status: DiffStatusSame, treeLevel: 1},
							{a: 3, b: 3, status: DiffStatusSame, treeLevel: 1},
						},
					},
					status: DiffStatusSame,
				},
			},
			"simple same different order": {
				a: rawTypeArray{1, 2, 3},
				b: rawTypeArray{3, 1, 2},
				want: &diff{
					children: &diffChildren{
						a: diffChildrenArray{
							{a: 1, b: 1, status: DiffStatusSame, treeLevel: 1},
							{a: 2, b: 2, status: DiffStatusSame, treeLevel: 1},
							{a: 3, b: 3, status: DiffStatusSame, treeLevel: 1},
						},
					},
					status: DiffStatusSame,
				},
			},
			"missing in A": {
				a: rawTypeArray{1, 2},
				b: rawTypeArray{1, 2, 3},
				want: &diff{
					children: &diffChildren{
						a: diffChildrenArray{
							{a: 1, b: 1, status: DiffStatusSame, treeLevel: 1},
							{a: 2, b: 2, status: DiffStatusSame, treeLevel: 1},
							{b: 3, status: DiffStatus1Missing, diffCount: 1, treeLevel: 1},
						},
					},
					diffCount: 1,
					status:    DiffStatusDiff,
				},
			},
			"missing in B": {
				a: rawTypeArray{1, 2, 3},
				b: rawTypeArray{1, 3},
				want: &diff{
					children: &diffChildren{
						a: diffChildrenArray{
							{a: 1, b: 1, status: DiffStatusSame, treeLevel: 1},
							{a: 3, b: 3, status: DiffStatusSame, treeLevel: 1},
							{a: 2, status: DiffStatus2Missing, diffCount: 1, treeLevel: 1}, // missing is added by last
						},
					},
					diffCount: 1,
					status:    DiffStatusDiff,
				},
			},
			"missing in A and B": {
				a: rawTypeArray{1, 2, 3},
				b: rawTypeArray{1, 3, 4},
				want: &diff{
					children: &diffChildren{
						a: diffChildrenArray{
							{a: 1, b: 1, status: DiffStatusSame, treeLevel: 1},
							{a: 3, b: 3, status: DiffStatusSame, treeLevel: 1},
							{a: 2, b: 4, status: DiffStatusDiff, diffCount: 1, treeLevel: 1}, // because can't find missing, it's diff.
						},
					},
					diffCount: 1,
					status:    DiffStatusDiff,
				},
			},
			"complicated": {
				a: rawTypeArray{1, rawTypeArray{2, 3, 4}, 5, 6},
				b: rawTypeArray{1, 5, rawTypeArray{2}},
				want: &diff{
					children: &diffChildren{
						a: diffChildrenArray{
							{a: 1, b: 1, status: DiffStatusSame, treeLevel: 1},
							{a: 5, b: 5, status: DiffStatusSame, treeLevel: 1},
							{
								a: rawTypeArray{2, 3, 4},
								b: rawTypeArray{2},
								children: &diffChildren{
									a: diffChildrenArray{
										{a: 2, b: 2, status: DiffStatusSame, treeLevel: 2},
										{a: 3, status: DiffStatus2Missing, diffCount: 1, treeLevel: 2},
										{a: 4, status: DiffStatus2Missing, diffCount: 1, treeLevel: 2},
									},
								},
								diffCount: 2,
								status:    DiffStatusDiff,
								treeLevel: 1,
							},
							{a: 6, status: DiffStatus2Missing, diffCount: 1, treeLevel: 1},
						},
					},
					diffCount: 2 + 1,
					status:    DiffStatusDiff,
				},
			},
		},
		"map": {
			"simple same": {
				a: rawTypeMap{
					"foo": "bar",
				},
				b: rawTypeMap{
					"foo": "bar",
				},
				want: &diff{
					children: &diffChildren{
						m: diffChildrenMap{
							"foo": {
								a:         "bar",
								b:         "bar",
								status:    DiffStatusSame,
								treeLevel: 1,
							},
						},
					},
					status: DiffStatusSame,
				},
			},
			"simple diff": {
				a: rawTypeMap{
					"foo": "bar",
				},
				b: rawTypeMap{
					"foo": "baz",
				},
				want: &diff{
					children: &diffChildren{
						m: diffChildrenMap{
							"foo": {
								a:         "bar",
								b:         "baz",
								status:    DiffStatusDiff,
								diffCount: 1,
								treeLevel: 1,
							},
						},
					},
					diffCount: 1,
					status:    DiffStatusDiff,
				},
			},
			"simple diff type": {
				a: rawTypeMap{
					"foo": "bar",
				},
				b: rawTypeMap{
					"foo": 1,
				},
				want: &diff{
					children: &diffChildren{
						m: diffChildrenMap{
							"foo": {
								a:         "bar",
								b:         1,
								diffCount: 3,
								status:    DiffStatusDiff,
								treeLevel: 1,
							},
						},
					},
					diffCount: 3,
					status:    DiffStatusDiff,
				},
			},
			"simple diff type map and primitive": {
				a: rawTypeMap{
					"foo": "bar",
				},
				b: rawTypeMap{
					"foo": rawTypeMap{
						"bar": "baz",
					},
				},
				want: &diff{
					children: &diffChildren{
						m: diffChildrenMap{
							"foo": {
								a: "bar",
								b: rawTypeMap{
									"bar": "baz",
								},
								diffCount: len("map[bar:baz]"), // as primitive diff
								status:    DiffStatusDiff,
								treeLevel: 1,
							},
						},
					},

					diffCount: len("map[bar:baz]"), // from child
					status:    DiffStatusDiff,
				},
			},
			"complicated": {
				a: rawTypeMap{
					"foo": rawTypeMap{
						"bar":  "baz",
						"baz":  1,
						"barr": false,
					},
					"bar": 1,
					"baz": "1",
					"zoo": 1,
				},
				b: rawTypeMap{
					"foo": rawTypeMap{
						"bar": "baz",
						"baz": rawTypeMap{
							"a": "b",
						},
						"bazz": 1,
					},
					"bar": "1",
					"baz": 1,
					"boo": 1,
				},
				want: &diff{
					children: &diffChildren{
						m: diffChildrenMap{
							"foo": {
								a: rawTypeMap{
									"bar":  "baz",
									"baz":  1,
									"barr": false,
								},
								b: rawTypeMap{
									"bar": "baz",
									"baz": rawTypeMap{
										"a": "b",
									},
									"bazz": 1,
								},
								children: &diffChildren{
									m: diffChildrenMap{
										"bar": {
											a:         "baz",
											b:         "baz",
											status:    DiffStatusSame,
											treeLevel: 2,
										},
										"baz": {
											a: 1,
											b: rawTypeMap{
												"a": "b",
											},
											status:    DiffStatusDiff,
											diffCount: len("map[a:b]"),
											treeLevel: 2,
										},
										"barr": {
											a:         false,
											status:    DiffStatus2Missing,
											diffCount: 5,
											treeLevel: 2,
										},
										"bazz": {
											b:         1,
											status:    DiffStatus1Missing,
											diffCount: 1,
											treeLevel: 2,
										},
									},
								},
								diffCount: (len("map[a:b]")) + (5) + (1),
								status:    DiffStatusDiff,
								treeLevel: 1,
							},
							"bar": {
								a:         1,
								b:         "1",
								diffCount: 0,
								status:    DiffStatusDiff,
								treeLevel: 1,
							},
							"baz": {
								a:         "1",
								b:         1,
								diffCount: 0,
								status:    DiffStatusDiff,
								treeLevel: 1,
							},
							"zoo": {
								a:         1,
								diffCount: 1,
								status:    DiffStatus2Missing,
								treeLevel: 1,
							},
							"boo": {
								b:         1,
								diffCount: 1,
								status:    DiffStatus1Missing,
								treeLevel: 1,
							},
						},
					},
					diffCount: ((len("map[a:b]")) + (5) + (1)) + (0) + (0) + (1) + (1),
					status:    DiffStatusDiff,
				},
			},
		},
	}
	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			for n, tc := range tt {
				tc := tc
				t.Run(n, func(t *testing.T) {
					t.Parallel()

					got := performDiff(tc.a, tc.b, 0)

					tc.want.a = tc.a
					tc.want.b = tc.b
					assert.Equal(t, tc.want, got)
				})
			}
		})
	}
}
