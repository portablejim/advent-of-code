package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

type Coord struct {
	y int64
	x int64
}

func move(dir string, from_pos Coord, num int64) Coord {
	output := Coord{-1, -1}
	if dir == "U" {
		output = Coord{from_pos.y - num, from_pos.x}
	} else if dir == "D" {
		output = Coord{from_pos.y + num, from_pos.x}
	} else if dir == "L" {
		output = Coord{from_pos.y, from_pos.x - num}
	} else if dir == "R" {
		output = Coord{from_pos.y, from_pos.x + num}
	}

	return output
}

func doCoordsMatch(a Coord, b Coord) bool {
	return a.x == b.x && a.y == b.y
}

type YardTile struct {
	pos         Coord
	isGarden    bool
	isStart     bool
	reachableIn int
	inSolution  bool
}

type CandidateTile struct {
	pos      Coord
	distance int
}

func main() {
	var filename = flag.String("f", "../inputs/d21.sample1.txt", "file to use")
    var steps = flag.Int("steps", 6, "steps for the day")
	var part2 = flag.Bool("part2", false, "do part 2")
	var visualise = flag.Bool("visualise", false, "check a string")
	flag.Parse()
	dat, err := os.ReadFile(*filename)
	if err != nil {
		log.Fatalf("unable to read file: %f", err)
	}
	if *part2 {
		fmt.Printf("Part 2")
	}

	yardPlots := [][]YardTile{}
	startingPoint := Coord{-1, -1}

	// Parse the data.
	for l_num, f_line := range strings.Split(strings.Replace(string(dat), "\r\n", "\n", -1), "\n") {
		if len(f_line) == 0 {
			continue
		}
		f_line = strings.Trim(f_line, " \n")

		yardTileRow := []YardTile{}
		for c_num := 0; c_num < len(f_line); c_num += 1 {
			yardTile := f_line[c_num]
			isGarden := yardTile == '.' || yardTile == 'S'
			isStart := yardTile == 'S'
			yardTileRow = append(yardTileRow, YardTile{Coord{int64(l_num), int64(c_num)}, isGarden, isStart, -1, false})
			if isStart {
				startingPoint = Coord{int64(l_num), int64(c_num)}
			}
		}
		yardPlots = append(yardPlots, yardTileRow)
	}

    // Check valid conditions.
	if len(yardPlots) == 0 || len(yardPlots[0]) == 0 {
		fmt.Fprintf(os.Stderr, "Error when reading file.\n")
		return
	}
	if startingPoint.y < 0 || startingPoint.x < 0 {
		fmt.Fprintf(os.Stderr, "Error finding starting point.\n")
		return
	}

    yardHeight := len(yardPlots)
    yardWidth := len(yardPlots[0])

    // Setup starting conditions
	yardPlots[startingPoint.y][startingPoint.x].reachableIn = 0
	pointQueue := []CandidateTile{{startingPoint, 0}}
	total := 0
	dirs := []string{"U", "D", "L", "R"}
	stepAllowance := *steps

    // Process queue.
	for len(pointQueue) > 0 {
		currentQueueItem := pointQueue[0]
		pointQueue = pointQueue[1:]

		for _, dir := range dirs {
			nextPos := move(dir, currentQueueItem.pos, 1)
            if nextPos.x < 0 || nextPos.y < 0 || nextPos.y >= int64(yardHeight) || nextPos.x >= int64(yardWidth) {
                continue
            }
            //fmt.Printf("D: %s %v %v\n", dir, currentQueueItem.pos, nextPos)
			nextTile := &yardPlots[nextPos.y][nextPos.x]
            nextReachable := currentQueueItem.distance + 1
			if nextTile.isGarden && (nextTile.reachableIn == -1 || nextTile.reachableIn > nextReachable) {
				nextTile.reachableIn = nextReachable
				pointQueue = append(pointQueue, CandidateTile{nextPos, nextReachable})
			}
		}

        slices.SortFunc(pointQueue, func(a CandidateTile, b CandidateTile) int {
            if a.distance == b.distance {
                return 0
            } else if a.distance < b.distance {
                return -1
            } else {
                return 1
            }
        })
        /*
        if len(pointQueue) > 1 {
            fmt.Printf("Q: %d %d\n", len(pointQueue), pointQueue[0].distance)
        }
        if len(pointQueue) > 100 {
            break
        }
        */
	}

	for i, plotLine := range yardPlots {
		for j, plotItem := range plotLine {
			if plotItem.isGarden && plotItem.reachableIn <= stepAllowance && (plotItem.reachableIn%2) == (stepAllowance%2) {
				total += 1
				currentPlotItem := yardPlots[i][j]
				currentPlotItem.inSolution = true
				yardPlots[i][j] = currentPlotItem
			}
		}
	}

	fmt.Printf("T: %d\n", total)

    if *visualise {
	for _, graph_line := range yardPlots {
		for _, graph_nde := range graph_line {
            printDistance := false
            if printDistance && graph_nde.reachableIn > 0 && graph_nde.reachableIn < 16 {
				fmt.Printf("%X", graph_nde.reachableIn)
            } else if graph_nde.inSolution {
				fmt.Printf("O")
			} else if graph_nde.isGarden {
				if graph_nde.isStart {
					fmt.Printf("S")
				} else if graph_nde.reachableIn > -1 {
					fmt.Printf(".")
				} else {
					fmt.Printf("!")
				}
			} else {
				fmt.Printf("#")
				//fmt.Printf("\u001b[31m%d\u001b[0m", graph_nde.cost)
			}
		}
		fmt.Printf("\n")
	}
    }
}
