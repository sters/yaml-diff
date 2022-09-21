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
					status: DiffStatusDiff,
				},
			},
			"int missing a": {
				b: 1,
				want: &rawDiff{
					status: DiffStatus1Missing,
				},
			},
			"int missing b": {
				a: 1,
				want: &rawDiff{
					status: DiffStatus2Missing,
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
					status: DiffStatusDiff,
				},
			},
			"string missing a": {
				b: "1",
				want: &rawDiff{
					status: DiffStatus1Missing,
				},
			},
			"string missing b": {
				a: "1",
				want: &rawDiff{
					status: DiffStatus2Missing,
				},
			},
			"int vs string": {
				a: "1",
				b: 0,
				want: &rawDiff{
					status: DiffStatusDiff,
				},
			},
			"int vs float": {
				a: 1,
				b: 0.5,
				want: &rawDiff{
					status: DiffStatusDiff,
				},
			},
			"float vs string": {
				a: "1",
				b: 0.5,
				want: &rawDiff{
					status: DiffStatusDiff,
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
						m: map[string]*rawDiff{
							"foo": {
								a:      "bar",
								b:      "bar",
								status: DiffStatusSame,
							},
						},
					},
					status: DiffStatusUnknown,
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
						m: map[string]*rawDiff{
							"foo": {
								a:      "bar",
								b:      "baz",
								status: DiffStatusDiff,
							},
						},
					},
					status: DiffStatusUnknown,
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
						m: map[string]*rawDiff{
							"foo": {
								a:      "bar",
								b:      1,
								status: DiffStatusDiff,
							},
						},
					},
					status: DiffStatusUnknown,
				},
			},
			"simple diff type 2": {
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
						m: map[string]*rawDiff{
							"foo": {
								a: "bar",
								b: rawMap{
									"bar": "baz",
								},
								status: DiffStatusDiff,
							},
						},
					},
					status: DiffStatusUnknown,
				},
			},
			"complicated": {
				a: rawMap{
					"foo": rawMap{
						"bar":  "baz",
						"barr": false,
					},
					"bar": 1,
					"baz": "1",
					"zoo": 1,
				},
				b: rawMap{
					"foo": rawMap{
						"bar":  "baz",
						"bazz": 1,
					},
					"bar": "1",
					"baz": 1,
					"boo": 1,
				},
				want: &rawDiff{
					child: &diffChildren{
						m: map[string]*rawDiff{
							"foo": {
								a: rawMap{
									"bar":  "baz",
									"barr": false,
								},
								b: rawMap{
									"bar":  "baz",
									"bazz": 1,
								},
								child: &diffChildren{
									m: map[string]*rawDiff{
										"bar": {
											a:      "baz",
											b:      "baz",
											status: DiffStatusSame,
										},
										"barr": {
											a:      false,
											status: DiffStatus2Missing,
										},
										"bazz": {
											b:      1,
											status: DiffStatus1Missing,
										},
									},
								},
								status: DiffStatusUnknown,
							},
							"bar": {
								a:      1,
								b:      "1",
								status: DiffStatusDiff,
							},
							"baz": {
								a:      "1",
								b:      1,
								status: DiffStatusDiff,
							},
							"zoo": {
								a:      1,
								status: DiffStatus2Missing,
							},
							"boo": {
								b:      1,
								status: DiffStatus1Missing,
							},
						},
					},
					status: DiffStatusUnknown,
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
