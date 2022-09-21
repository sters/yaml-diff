package yamldiff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_diff(t *testing.T) {
	tests := map[string]map[string]struct {
		a    rawRaw
		b    rawRaw
		want *rawDiff
	}{
		"primitive": {
			"int ok": {
				a: 1,
				b: 1,
				want: &rawDiff{
					status: DiffStatusSame,
				},
			},
			"int diff": {
				a: 1,
				b: 0,
				want: &rawDiff{
					status:    DiffStatusDiff,
					diffCount: 1,
				},
			},
			"int missing a": {
				b: 1,
				want: &rawDiff{
					status:    DiffStatus1Missing,
					diffCount: 1,
				},
			},
			"int missing b": {
				a: 11,
				want: &rawDiff{
					status:    DiffStatus2Missing,
					diffCount: 2,
				},
			},
			"string ok": {
				a: "1",
				b: "1",
				want: &rawDiff{
					status: DiffStatusSame,
				},
			},
			"string diff": {
				a: "1",
				b: "0",
				want: &rawDiff{
					status:    DiffStatusDiff,
					diffCount: 1,
				},
			},
			"string missing a": {
				b: "1",
				want: &rawDiff{
					status:    DiffStatus1Missing,
					diffCount: 1,
				},
			},
			"string missing b": {
				a: "11",
				want: &rawDiff{
					status:    DiffStatus2Missing,
					diffCount: 2,
				},
			},
			"int vs string": {
				a: "1",
				b: 1,
				want: &rawDiff{
					status:    DiffStatusDiff,
					diffCount: 0, // because it's only diff on type
				},
			},
			"int vs string 2": {
				a: "1",
				b: 0,
				want: &rawDiff{
					status:    DiffStatusDiff,
					diffCount: 1,
				},
			},
			"int vs float": {
				a: 1,
				b: 0.5,
				want: &rawDiff{
					status:    DiffStatusDiff,
					diffCount: 3,
				},
			},
			"float vs string": {
				a: "1",
				b: 0.5,
				want: &rawDiff{
					status:    DiffStatusDiff,
					diffCount: 3,
				},
			},
			"int vs bool": {
				a: 1,
				b: false,
				want: &rawDiff{
					status:    DiffStatusDiff,
					diffCount: 5,
				},
			},
		},
		"array": {
			"simple same": {
				a: rawArray{1, 2, 3},
				b: rawArray{1, 2, 3},
				want: &rawDiff{
					child: &diffChildren{
						a: diffChildrenArray{
							{a: 1, b: 1, status: DiffStatusSame},
							{a: 2, b: 2, status: DiffStatusSame},
							{a: 3, b: 3, status: DiffStatusSame},
						},
					},
					status: DiffStatusSame,
				},
			},
			"simple same different order": {
				a: rawArray{1, 2, 3},
				b: rawArray{3, 1, 2},
				want: &rawDiff{
					child: &diffChildren{
						a: diffChildrenArray{
							{a: 1, b: 1, status: DiffStatusSame},
							{a: 2, b: 2, status: DiffStatusSame},
							{a: 3, b: 3, status: DiffStatusSame},
						},
					},
					status: DiffStatusSame,
				},
			},
			"missing in A": {
				a: rawArray{1, 2},
				b: rawArray{1, 2, 3},
				want: &rawDiff{
					child: &diffChildren{
						a: diffChildrenArray{
							{a: 1, b: 1, status: DiffStatusSame},
							{a: 2, b: 2, status: DiffStatusSame},
							{b: 3, status: DiffStatus1Missing, diffCount: 1},
						},
					},
					diffCount: 1,
					status:    DiffStatusDiff,
				},
			},
			"missing in B": {
				a: rawArray{1, 2, 3},
				b: rawArray{1, 3},
				want: &rawDiff{
					child: &diffChildren{
						a: diffChildrenArray{
							{a: 1, b: 1, status: DiffStatusSame},
							{a: 3, b: 3, status: DiffStatusSame},
							{a: 2, status: DiffStatus2Missing, diffCount: 1}, // missing is added by last
						},
					},
					diffCount: 1,
					status:    DiffStatusDiff,
				},
			},
			"missing in A and B": {
				a: rawArray{1, 2, 3},
				b: rawArray{1, 3, 4},
				want: &rawDiff{
					child: &diffChildren{
						a: diffChildrenArray{
							{a: 1, b: 1, status: DiffStatusSame},
							{a: 3, b: 3, status: DiffStatusSame},
							{a: 2, b: 4, status: DiffStatusDiff, diffCount: 1}, // because can't find missing, it's diff.
						},
					},
					diffCount: 1,
					status:    DiffStatusDiff,
				},
			},
			"complicated": {
				a: rawArray{1, rawArray{2, 3, 4}, 5, 6},
				b: rawArray{1, 5, rawArray{2}},
				want: &rawDiff{
					child: &diffChildren{
						a: diffChildrenArray{
							{a: 1, b: 1, status: DiffStatusSame},
							{a: 5, b: 5, status: DiffStatusSame},
							{
								a: rawArray{2, 3, 4},
								b: rawArray{2},
								child: &diffChildren{
									a: diffChildrenArray{
										{a: 2, b: 2, status: DiffStatusSame},
										{a: 3, status: DiffStatus2Missing, diffCount: 1},
										{a: 4, status: DiffStatus2Missing, diffCount: 1},
									},
								},
								diffCount: 2,
								status:    DiffStatusDiff,
							},
							{a: 6, status: DiffStatus2Missing, diffCount: 1},
						},
					},
					diffCount: 2 + 1,
					status:    DiffStatusDiff,
				},
			},
		},
		"map": {
			"simple same": {
				a: rawMap{
					"foo": "bar",
				},
				b: rawMap{
					"foo": "bar",
				},
				want: &rawDiff{
					child: &diffChildren{
						m: diffChildrenMap{
							"foo": {
								a:      "bar",
								b:      "bar",
								status: DiffStatusSame,
							},
						},
					},
					status: DiffStatusSame,
				},
			},
			"simple diff": {
				a: rawMap{
					"foo": "bar",
				},
				b: rawMap{
					"foo": "baz",
				},
				want: &rawDiff{
					child: &diffChildren{
						m: diffChildrenMap{
							"foo": {
								a:         "bar",
								b:         "baz",
								status:    DiffStatusDiff,
								diffCount: 1,
							},
						},
					},
					diffCount: 1,
					status:    DiffStatusDiff,
				},
			},
			"simple diff type": {
				a: rawMap{
					"foo": "bar",
				},
				b: rawMap{
					"foo": 1,
				},
				want: &rawDiff{
					child: &diffChildren{
						m: diffChildrenMap{
							"foo": {
								a:         "bar",
								b:         1,
								diffCount: 3,
								status:    DiffStatusDiff,
							},
						},
					},
					diffCount: 3,
					status:    DiffStatusDiff,
				},
			},
			"simple diff type map and primitive": {
				a: rawMap{
					"foo": "bar",
				},
				b: rawMap{
					"foo": rawMap{
						"bar": "baz",
					},
				},
				want: &rawDiff{
					child: &diffChildren{
						m: diffChildrenMap{
							"foo": {
								a: "bar",
								b: rawMap{
									"bar": "baz",
								},
								diffCount: len("map[bar:baz]"), // as primitive diff
								status:    DiffStatusDiff,
							},
						},
					},

					diffCount: len("map[bar:baz]"), // from child
					status:    DiffStatusDiff,
				},
			},
			"complicated": {
				a: rawMap{
					"foo": rawMap{
						"bar":  "baz",
						"baz":  1,
						"barr": false,
					},
					"bar": 1,
					"baz": "1",
					"zoo": 1,
				},
				b: rawMap{
					"foo": rawMap{
						"bar": "baz",
						"baz": rawMap{
							"a": "b",
						},
						"bazz": 1,
					},
					"bar": "1",
					"baz": 1,
					"boo": 1,
				},
				want: &rawDiff{
					child: &diffChildren{
						m: diffChildrenMap{
							"foo": {
								a: rawMap{
									"bar":  "baz",
									"baz":  1,
									"barr": false,
								},
								b: rawMap{
									"bar": "baz",
									"baz": rawMap{
										"a": "b",
									},
									"bazz": 1,
								},
								child: &diffChildren{
									m: diffChildrenMap{
										"bar": {
											a:      "baz",
											b:      "baz",
											status: DiffStatusSame,
										},
										"baz": {
											a: 1,
											b: rawMap{
												"a": "b",
											},
											status:    DiffStatusDiff,
											diffCount: len("map[a:b]"),
										},
										"barr": {
											a:         false,
											status:    DiffStatus2Missing,
											diffCount: 5,
										},
										"bazz": {
											b:         1,
											status:    DiffStatus1Missing,
											diffCount: 1,
										},
									},
								},
								diffCount: (len("map[a:b]")) + (5) + (1),
								status:    DiffStatusDiff,
							},
							"bar": {
								a:         1,
								b:         "1",
								diffCount: 0,
								status:    DiffStatusDiff,
							},
							"baz": {
								a:         "1",
								b:         1,
								diffCount: 0,
								status:    DiffStatusDiff,
							},
							"zoo": {
								a:         1,
								diffCount: 1,
								status:    DiffStatus2Missing,
							},
							"boo": {
								b:         1,
								diffCount: 1,
								status:    DiffStatus1Missing,
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

					got := diff(tc.a, tc.b)

					tc.want.a = tc.a
					tc.want.b = tc.b
					assert.Equal(t, tc.want, got)
				})
			}
		})
	}
}
