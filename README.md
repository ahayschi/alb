# Assembly Line Balancing
ALB is a library for constructing, manipulating, and balancing abstract assembly lines. In addition to the library, there is a supporting main command to construct a line from an in2 file for balancing and stat logging.

#### Examples

Constructing a new line with tasks and stations:
```go
package main

import (
	"github.com/parallelworks/alb"
)

func main() {
        line := alb.NewLine("Door Assembly North Region")
        tasks := []*alb.Task{
                alb.NewTask(1, 7.0),
        }
        line.AddTasks(tasks)

        stations := []*alb.Station{
                alb.NewStation(1),
        }
        line.AddStations(stations)
}
```

Adding task-station assignment constraints to a line:
```go
constraints := []alb.Constraint{
        &alb.SingleTaskAssignment{},
}
line.AddConstraints(constraints)
```

Adding predecessors to a task:
```go
task1 := alb.NewTask(1, 7.0)
task2 := alb.NewTask(2, 10.0)

task1.AddPred(task2)
```

#### Constraints
TODO

#### Heuristics
TODO

#### Balancing

ALB currently has two balancing methods: the first balances a line by iterating over a line's stations in order by Id, making task assignments. The second balances a line by making task assignments to the station with the shortest time at assignment.

Whichever balance method you choose, you need to provide a heuristic for picking the task to assign from a set of valid tasks. I recommend you use either ```ShortestTaskTime``` or ```LongestTaskTime```, as they are the simplest to verify and test. LTT has been shown to produce better results than STT.

## Development
Since this is currently a private repository, you will need to manually put it in the right place in your ```GOPATH```.

```bash
cd $GOPATH/src/github.com/parallelworks/alb
git clone git@github.com:parallelworks/alb.git
```

The Makefile currently supports two build targets:

```bash
make mac
make linux
```

### Testing the Library
Much work needs to be done just on the current codebase to ensure correctnes. Stubs of some unit tests exist for stations and tasks and they can be run by:
```bash
make test
```

### Running Balance Command
After building, run with the in2 file and the cycle time:

```bash
./bin/balance -file=specs/buxey.in2 -cycle=37
```
