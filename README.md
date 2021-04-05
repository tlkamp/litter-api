# Litter API  [![Go Report Badge]][Go Report] [![GoDocBadge]][GoDocLink]

Litter API is an interface to the [Litter Robot](https://www.litter-robot.com/) API.

**This is an experimental API.** The upstream Litter Robot API is not publicly documented and may cause breaking
changes with no notice. Breaking changes will be handled as soon as possible.

## Examples

```go
package main

import (
	"fmt"
	. "github.com/tlkamp/litter-api"
	"log"
)

func main() {
	lc, err := NewClient(Config{
		Email:    "your-email@example.com",
		Password: "your-password-here",
		APIKey:   "your-api-key",
	})
	if err != nil {
		log.Fatalln(err)
	}

	states, err := lc.States()
	if err != nil {
		log.Fatalln(err)
	}

	for _, state := range states {
		log.Println(fmt.Sprintf("%+v", state))
	}
}
```

## Logging
[Logrus] is the logger used in this project. The log level and format can be altered
accordingly.

```go
log.SetLevel(log.DebugLevel)
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

[Logrus]: https://github.com/sirupsen/logrus
[Go Report Badge]: https://goreportcard.com/badge/github.com/tlkamp/litter-api
[Go Report]: https://goreportcard.com/report/github.com/tlkamp/litter-api
[GoDocBadge]: https://godoc.org/github.com/tlkamp/litter-api?status.svg
[GoDocLink]: https://godoc.org/github.com/tlkamp/litter-api
