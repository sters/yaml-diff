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
			"int null a": {
				b: 1,
				want: &diff{
					status:    DiffStatusDiff,
					diffCount: 5,
				},
			},
			"int null b": {
				a: 11,
				want: &diff{
					status:    DiffStatusDiff,
					diffCount: 5,
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
			"string null a": {
				b: "1",
				want: &diff{
					status:    DiffStatusDiff,
					diffCount: 5,
				},
			},
			"string null b": {
				a: "11",
				want: &diff{
					status:    DiffStatusDiff,
					diffCount: 5,
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
							{b: 3, status: DiffStatusDiff, diffCount: 5, treeLevel: 1},
						},
					},
					diffCount: 5,
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
							{a: 2, status: DiffStatusDiff, diffCount: 5, treeLevel: 1}, // missing is added by last
						},
					},
					diffCount: 5,
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
								a:         6,
								b:         rawTypeArray{2},
								diffCount: 3,
								status:    DiffStatusDiff,
								treeLevel: 1,
							},
							{
								a:         rawTypeArray{2, 3, 4},
								status:    DiffStatusDiff,
								diffCount: 7,
								treeLevel: 1,
							},
						},
					},
					diffCount: 3 + 7,
					status:    DiffStatusDiff,
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
								status:    DiffStatusDiff,
								treeLevel: 1,
							},
						},
					},

					diffCount: len("{bar baz}"), // from child
					status:    DiffStatusDiff,
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
											status:    DiffStatusSame,
											treeLevel: 2,
										},
										"baz": {
											a: 1,
											b: rawTypeMap{
												yaml.MapItem{Key: "a", Value: "b"},
											},
											status:    DiffStatusDiff,
											diffCount: len("[{a b}]"),
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
								diffCount: (len("[{a b}]")) + (5) + (1),
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
					diffCount: ((len("[{a b}]")) + (5) + (1)) + (0) + (0) + (1) + (1),
					status:    DiffStatusDiff,
				},
			},
		},
		"#27": {
			"creationTimestamp: null": {
				a: rawTypeMap{
					yaml.MapItem{Key: "metadata", Value: rawTypeMap{
						yaml.MapItem{Key: "creationTimestamp", Value: nil},
						yaml.MapItem{Key: "name", Value: "nginx"},
						yaml.MapItem{Key: "namespace", Value: "default"},
					}},
					yaml.MapItem{Key: "status", Value: nil},
				},
				b: rawTypeMap{
					yaml.MapItem{Key: "metadata", Value: rawTypeMap{
						yaml.MapItem{Key: "name", Value: "nginx"},
						yaml.MapItem{Key: "namespace", Value: "default"},
					}},
				},
				want: &diff{
					children: &diffChildren{
						m: diffChildrenMap{
							"metadata": {
								a: rawTypeMap{
									yaml.MapItem{Key: "creationTimestamp", Value: nil},
									yaml.MapItem{Key: "name", Value: "nginx"},
									yaml.MapItem{Key: "namespace", Value: "default"},
								},
								b: rawTypeMap{
									yaml.MapItem{Key: "name", Value: "nginx"},
									yaml.MapItem{Key: "namespace", Value: "default"},
								},
								children: &diffChildren{
									m: diffChildrenMap{
										"creationTimestamp": &diff{
											status:    DiffStatus2Missing,
											a:         nil,
											treeLevel: 2,
											diffCount: 5,
										},
										"name": &diff{
											status:    DiffStatusSame,
											a:         "nginx",
											b:         "nginx",
											treeLevel: 2,
										},
										"namespace": &diff{
											status:    DiffStatusSame,
											a:         "default",
											b:         "default",
											treeLevel: 2,
										},
									},
								},
								diffCount: 5,
								status:    DiffStatusDiff,
								treeLevel: 1,
							},
							"status": &diff{
								status:    DiffStatus2Missing,
								a:         nil,
								diffCount: 5,
								treeLevel: 1,
							},
						},
					},
					diffCount: 5 + 5,
					status:    DiffStatusDiff,
				},
			},
		},
		"#30": {
			"zero string": {
				a: rawTypeMap{
					yaml.MapItem{Key: "strA", Value: "foo"},
					yaml.MapItem{Key: "strB", Value: ""},
					yaml.MapItem{Key: "strC", Value: ""},
				},
				b: rawTypeMap{
					yaml.MapItem{Key: "strA", Value: "foo"},
					yaml.MapItem{Key: "strB", Value: ""},
				},
				want: &diff{
					children: &diffChildren{
						m: diffChildrenMap{
							"strA": &diff{
								status:    DiffStatusSame,
								a:         "foo",
								b:         "foo",
								treeLevel: 1,
							},
							"strB": &diff{
								status:    DiffStatusSame,
								a:         "",
								b:         "",
								treeLevel: 1,
							},
							"strC": &diff{
								status:    DiffStatus2Missing,
								a:         "",
								treeLevel: 1,
							},
						},
					},
					status: DiffStatusDiff,
				},
			},
			"zero int": {
				a: rawTypeMap{
					yaml.MapItem{Key: "intA", Value: 5},
					yaml.MapItem{Key: "intB", Value: 0},
					yaml.MapItem{Key: "intC", Value: 0},
				},
				b: rawTypeMap{
					yaml.MapItem{Key: "intA", Value: 5},
					yaml.MapItem{Key: "intB", Value: 0},
				},
				want: &diff{
					children: &diffChildren{
						m: diffChildrenMap{
							"intA": &diff{
								status:    DiffStatusSame,
								a:         5,
								b:         5,
								treeLevel: 1,
							},
							"intB": &diff{
								status:    DiffStatusSame,
								a:         0,
								b:         0,
								treeLevel: 1,
							},
							"intC": &diff{
								status:    DiffStatus2Missing,
								a:         0,
								diffCount: 1,
								treeLevel: 1,
							},
						},
					},
					diffCount: 1,
					status:    DiffStatusDiff,
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

					got := (&runner{}).performDiff(tc.a, tc.b, 0)

					tc.want.a = tc.a
					tc.want.b = tc.b
					assert.Equal(t, tc.want, got)
				})
			}
		})
	}
}

func Test_performDiff_emptyAsNull(t *testing.T) {
	tests := map[string]map[string]struct {
		a    rawType
		b    rawType
		want *diff
	}{
		"#27": {
			"creationTimestamp: null": {
				a: rawTypeMap{
					yaml.MapItem{Key: "metadata", Value: rawTypeMap{
						yaml.MapItem{Key: "creationTimestamp", Value: nil},
						yaml.MapItem{Key: "name", Value: "nginx"},
						yaml.MapItem{Key: "namespace", Value: "default"},
					}},
					yaml.MapItem{Key: "status", Value: nil},
				},
				b: rawTypeMap{
					yaml.MapItem{Key: "metadata", Value: rawTypeMap{
						yaml.MapItem{Key: "name", Value: "nginx"},
						yaml.MapItem{Key: "namespace", Value: "default"},
					}},
				},
				want: &diff{
					children: &diffChildren{
						m: diffChildrenMap{
							"metadata": {
								a: rawTypeMap{
									yaml.MapItem{Key: "creationTimestamp", Value: nil},
									yaml.MapItem{Key: "name", Value: "nginx"},
									yaml.MapItem{Key: "namespace", Value: "default"},
								},
								b: rawTypeMap{
									yaml.MapItem{Key: "name", Value: "nginx"},
									yaml.MapItem{Key: "namespace", Value: "default"},
								},
								children: &diffChildren{
									m: diffChildrenMap{
										"creationTimestamp": &diff{
											status:    DiffStatusSame,
											a:         nil,
											b:         missingKey,
											treeLevel: 2,
										},
										"name": &diff{
											status:    DiffStatusSame,
											a:         "nginx",
											b:         "nginx",
											treeLevel: 2,
										},
										"namespace": &diff{
											status:    DiffStatusSame,
											a:         "default",
											b:         "default",
											treeLevel: 2,
										},
									},
								},
								status:    DiffStatusSame,
								treeLevel: 1,
							},
							"status": &diff{
								status:    DiffStatusSame,
								a:         nil,
								b:         missingKey,
								treeLevel: 1,
							},
						},
					},
					status: DiffStatusSame,
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

					got := (&runner{option: doOptions{emptyAsNull: true}}).performDiff(tc.a, tc.b, 0)

					tc.want.a = tc.a
					tc.want.b = tc.b
					assert.Equal(t, tc.want, got)
				})
			}
		})
	}
}

func Test_performDiff_zeroAsNull(t *testing.T) {
	tests := map[string]map[string]struct {
		a    rawType
		b    rawType
		want *diff
	}{
		"#30": {
			"zero string": {
				a: rawTypeMap{
					yaml.MapItem{Key: "strA", Value: "foo"},
					yaml.MapItem{Key: "strB", Value: ""},
					yaml.MapItem{Key: "strC", Value: ""},
				},
				b: rawTypeMap{
					yaml.MapItem{Key: "strA", Value: "foo"},
					yaml.MapItem{Key: "strB", Value: ""},
				},
				want: &diff{
					children: &diffChildren{
						m: diffChildrenMap{
							"strA": &diff{
								status:    DiffStatusSame,
								a:         "foo",
								b:         "foo",
								treeLevel: 1,
							},
							"strB": &diff{
								status:    DiffStatusSame,
								a:         "",
								b:         "",
								treeLevel: 1,
							},
							"strC": &diff{
								status:    DiffStatusSame,
								a:         "",
								b:         missingKey,
								treeLevel: 1,
							},
						},
					},
					status: DiffStatusSame,
				},
			},
			"zero int": {
				a: rawTypeMap{
					yaml.MapItem{Key: "intA", Value: 5},
					yaml.MapItem{Key: "intB", Value: 0},
					yaml.MapItem{Key: "intC", Value: 0},
				},
				b: rawTypeMap{
					yaml.MapItem{Key: "intA", Value: 5},
					yaml.MapItem{Key: "intB", Value: 0},
				},
				want: &diff{
					children: &diffChildren{
						m: diffChildrenMap{
							"intA": &diff{
								status:    DiffStatusSame,
								a:         5,
								b:         5,
								treeLevel: 1,
							},
							"intB": &diff{
								status:    DiffStatusSame,
								a:         0,
								b:         0,
								treeLevel: 1,
							},
							"intC": &diff{
								status:    DiffStatusSame,
								a:         0,
								b:         missingKey,
								treeLevel: 1,
							},
						},
					},
					status: DiffStatusSame,
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

					got := (&runner{option: doOptions{zeroAsNull: true}}).performDiff(tc.a, tc.b, 0)

					tc.want.a = tc.a
					tc.want.b = tc.b
					assert.Equal(t, tc.want, got)
				})
			}
		})
	}
}
