package jiraquery

import (
	"testing"
	"time"
)

var (
	TIME_OCT_12_2018, _ = time.Parse("2006-1-2", "2018-10-12")
	TIME_OCT_15_2018, _ = time.Parse("2006-1-2", "2018-10-15")
)

func Test_Queries(t *testing.T) {
	tests := []struct {
		query    Query
		expected string
	}{
		{
			Query(
				And(
					Word("foo"),
					Word("bar"),
				),
			),
			"foo AND bar",
		},
		{
			Query(
				Or(
					Word("foo"),
					Word("bar"),
				),
			),
			"foo OR bar",
		},
		{
			Query(
				Not(
					Word("foo"),
				),
			),
			"NOT foo",
		},
		{
			Query(
				And(
					Word("foo"),
					Wrapped(
						Or(
							Word("bar"),
							Word("baz"),
						),
					),
				),
			),
			"foo AND ( bar OR baz )",
		},
		{
			Query(
				In(
					Word("foo"),
					List("yolo", "solo"),
				),
			),
			"foo IN [ yolo, solo ]",
		},
		{
			Query(
				Project("JIRA"),
			),
			"project = \"JIRA\"",
		},
		{
			Query(
				IssueType("Bug"),
			),
			"issueType = \"Bug\"",
		},
		{
			Query(
				CreatedBefore(TIME_OCT_12_2018),
			),
			"created < 2018-10-12 00:00",
		},
		{
			Query(
				CreatedAfter(TIME_OCT_12_2018),
			),
			"created > 2018-10-12 00:00",
		},
		{
			Query(
				NotEq(Word("yolo"), Word("solo")),
			),
			"yolo != solo",
		},
		{
			Query(
				GreaterThan(Word("field"), Word("5")),
			),
			"field > 5",
		},
		{
			Query(
				LessThan(Word("field"), Word("5")),
			),
			"field < 5",
		},
		{
			Query(
				MultiAnd(
					GreaterThan(Word("field"), Word("5")),
					LessThan(Word("field"), Word("10")),
				),
			),
			"field > 5 AND field < 10",
		},
		{
			Query(
				MultiOr(
					GreaterThan(Word("field"), Word("5")),
					LessThan(Word("field"), Word("10")),
				),
			),
			"field > 5 OR field < 10",
		},
	}

	for i, test := range tests {
		query, expected := test.query, test.expected

		got := query.String()
		if got != expected {
			t.Errorf("[test %d] got %q, expected %q\n", i, got, expected)
		}
	}
}

func Test_QueryBuilder(t *testing.T) {
	tests := []struct {
		query    Query
		expected string
	}{
		{
			AndBuilder().
				Project("JIRA").
				IssueType("Bug").
				Value(),
			"project = \"JIRA\" AND issueType = \"Bug\"",
		}, {
			AndBuilder().
				CreatedAfter(TIME_OCT_12_2018).
				CreatedBefore(TIME_OCT_15_2018).
				Value(),
			"created > 2018-10-12 00:00 AND created < 2018-10-15 00:00",
		}, {
			OrBuilder().
				Not(Wrapped(Project("JIRA"))).
				Eq(Word("foo"), Word("bar")).
				NotEq(Word("foo"), Word("baz")).
				Value(),
			"NOT ( project = \"JIRA\" ) OR foo = bar OR foo != baz",
		}, {
			AndBuilder().
				GreaterThan(Word("foo"), Word("5")).
				LessThan(Word("foo"), Word("10")).
				In(Word("bar"), List("baz", "bolo")).
				Value(),
			"foo > 5 AND foo < 10 AND bar IN [ baz, bolo ]",
		}, {
			AndBuilder().
				Wrapped(
					OrBuilder().
						Eq(Word("foo"), Word("5")).
						NotEq(Word("foo"), Word("10")).
						Value(),
				).
				In(Word("bar"), List("baz", "bolo")).
				Value(),
			"( foo = 5 OR foo != 10 ) AND bar IN [ baz, bolo ]",
		},
	}

	for i, test := range tests {
		query, expected := test.query, test.expected

		got := query.String()
		if got != expected {
			t.Errorf("[test %d] got %q, expected %q\n", i, got, expected)
		}
	}
}
