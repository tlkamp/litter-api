# Litter API  [![Go Report Badge]][Go Report] [![GoDocBadge]][GoDocLink]

Litter API is an interface to the [Litter Robot](https://www.litter-robot.com/) API.

**This is an experimental API.** The upstream Litter Robot API is not publicly documented and may cause breaking
changes with no notice. Breaking changes will be handled as soon as possible.

## Examples

```go
package main

import (
	"fmt"
	lr "github.com/tlkamp/litter-api/v2/pkg/client"
)

func main() {
	api := lr.New("your-email", "your-password")

	ctx := context.Background()

	err := api.Login(ctx)
	if err != nil {
		fmt.Println("Error logging in - ", err)
	}

	if err := api.FetchRobots(ctx); err != nil {
		fmt.Println("Error fetching robots - ", err)
	}

	for _, robot := range api.Robots() {
		fmt.Println(robot.LitterRobotID, "-", robot.Name)
		api.Cycle(ctx, robot.LitterRobotID)
	}
}
```

## Unit Status
The unit status is represented by a non-negative integer.

| **String** | **Int** | **Description**                      |
|------------|---------|--------------------------------------|
| RDY        | 0       | Ready                                |
| CCP        | 1       | Clean Cycle in Progress              |
| CCC        | 2       | Clean Cycle Complete                 |
| CSF        | 3       | Cat Sensor Fault                     |
| DF1        | 4       | Drawer full - will still cycle       |
| DF2        | 5       | Drawer full - will still cycle       |
| CST        | 6       | Cat Sensor Timing                    |
| CSI        | 7       | Cat Sensor Interrupt                 |
| BR         | 8       | Bonnet Removed                       |
| P          | 9       | Paused                               |
| OFF        | 10      | Off                                  |
| SDF        | 11      | Drawer full - will not cycle         |
| DFS        | 12      | Drawer full - will not cycle         |

[Go Report Badge]: https://goreportcard.com/badge/github.com/tlkamp/litter-api
[Go Report]: https://goreportcard.com/report/github.com/tlkamp/litter-api
[GoDocBadge]: https://godoc.org/github.com/tlkamp/litter-api?status.svg
[GoDocLink]: https://godoc.org/github.com/tlkamp/litter-api
