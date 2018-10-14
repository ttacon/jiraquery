package main

import (
	"fmt"
	"time"

	jq "github.com/ttacon/jiraquery"
)

func main() {
	query := jq.AndBuilder().
		Project("FOO").
		IssueType("Bug").
		Wrapped(
			jq.OrBuilder().
				CreatedAfter(time.Now().AddDate(0, -30, 0)).
				NotEq(jq.Word("statusCategory"), jq.Word("Done")).
				Value(),
		).
		Value()

	fmt.Println(query.String())
	// project = "FOO" OR issueType = "Bug"
}
