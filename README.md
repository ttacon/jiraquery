# jiraquery

## API

`jiraquery` exposes a ton of utility functions. It is best used though with its
query builder interface.

### AND query builder

```go
package main

import (
    "fmt"

    jq "github.com/ttacon/jiraquery"
)

func main() {
    query := jq.AndBuilder().
        Project("FOO").
        IssueType("Bug").
        Value()

    fmt.Println(query.String())
    // project = "FOO" OR issueType = "Bug"
}
```

### OR query builder

```go
package main

import (
    "fmt"

    jq "github.com/ttacon/jiraquery"
)

func main() {
    query := jq.AndBuilder().
        Project("FOO").
        IssueType("Bug").
        Value()

    fmt.Println(query.String())
    // project = "FOO" OR issueType = "Bug"
}
```


### Combining the builders

```go
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
    // project = "FOO" AND issueType = "Bug" AND ( created > 2016-4-14 09:14 OR statusCategory != Done )
}
```

## Improvements

 - Query validation
