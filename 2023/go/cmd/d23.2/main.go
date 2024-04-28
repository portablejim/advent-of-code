package main

import (
	"container/heap"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	y int16
	x int16
}

func (p Point) move(dir string, num int16) Point {
	return move(dir, p, num)
}

func (p Point) isInvalid(testMap2d [][]MapTile) bool {
	output := false
	output = output || p.y < 0
	output = output || p.x < 0
	output = output || int(p.y) >= len(testMap2d)
	output = output || int(p.x) >= len(testMap2d[p.y])
	return output
}

func move(dir string, from_pos Point, num int16) Point {
	output := Point{-1, -1}
	if dir == "U" {
		output = Point{from_pos.y - num, from_pos.x}
	} else if dir == "D" {
		output = Point{from_pos.y + num, from_pos.x}
	} else if dir == "L" {
		output = Point{from_pos.y, from_pos.x - num}
	} else if dir == "R" {
		output = Point{from_pos.y, from_pos.x + num}
	}

	return output
}

func doPointsMatch(a Point, b Point) bool {
	return a.x == b.x && a.y == b.y
}

type MapTile struct {
	pos                Point
	isPath             bool
	tileChar           string
	visited            bool
	highest_steps      int
	highest_steps_path []Point
	isNode             bool
}

type GraphEdge struct {
	direction       string
	weight          int
	connectionPoint Point
	connectionIndex int
	connectionPath  []Point
}

type GraphNode struct {
	pos   Point
	edges []GraphEdge
}

type GraphNodeDiscovery struct {
	node      GraphNode
	direction string
}

type PendingWalk struct {
	pos     Point
	steps   int
	history []Point
}

type LongestWalk []PendingWalk

func (lw LongestWalk) Len() int { return len(lw) }

func (lw LongestWalk) Less(i, j int) bool {
	return lw[i].steps > lw[j].steps
}

func (lw LongestWalk) Swap(i, j int) {
	lw[i], lw[j] = lw[j], lw[i]
}

func (lw *LongestWalk) Push(x any) {
	item := x.(PendingWalk)
	*lw = append(*lw, item)
}

func (lw *LongestWalk) Pop() any {
	old := *lw
	n := len(old)
	item := old[n-1]
	*lw = old[0 : n-1]

	return item
}

func generatePositionKey(currentWalk PendingWalk) int64 {
	currentPosKey := int64(0)

	currentPosKey += int64(currentWalk.pos.y) << 16
	currentPosKey += int64(currentWalk.pos.x)

	return currentPosKey
}

func moveDirection(direction string, currentPos PendingWalk) PendingWalk {
	nextPos := currentPos.pos.move(direction, 1)
	nextHistory := []Point{}
	nextHistory = append(nextHistory, currentPos.history...)
	nextHistory = append(nextHistory, nextPos)
	//fmt.Printf("MoveDir %s | %v | %v\n", direction, currentPos.history, nextHistory)
	return PendingWalk{nextPos, currentPos.steps + 1, nextHistory}
}

func getGraphNodes(mapTiles2d [][]MapTile, startingNode GraphNode, endingNode GraphNode, ignoreTiles bool) []GraphNode {
	pendingNodes := []GraphNodeDiscovery{{startingNode, "D"}}

	allDirs := []string{"R", "D", "L", "U"}
	inverseDirs := map[string]string{"R": "L", "D": "U", "L": "R", "U": "D"}
	validDirs := map[string][]string{".": {"U", "R", "D", "L"}, "^": {"U"}, ">": {"R"}, "v": {"D"}, "<": {"L"}, "#": {}}

	foundNodes := []GraphNode{}
	foundNodes = append(foundNodes, startingNode)
	foundNodes = append(foundNodes, endingNode)

	for len(pendingNodes) > 0 {
		currentNodeDiscovery := pendingNodes[0]
		currentNode := currentNodeDiscovery.node
		pendingNodes = pendingNodes[1:]

		fmt.Printf("Node %d,%d from %s\n", currentNode.pos.y, currentNode.pos.x, currentNodeDiscovery.direction)

	initDirLoop:
		for _, initDir := range allDirs {
			nextSteps := 1
			nextPos := currentNode.pos.move(initDir, 1)
			nextDir := initDir
			nextPosNotNode := true
			history := []Point{}
			isOneWay := false

			if nextPos.y == endingNode.pos.y && nextPos.x == endingNode.pos.x {
				fmt.Printf("Ending node\n")
			}

			for nextPosNotNode {
				fmt.Printf("Next pos: %d,%d | %s | %s\n", nextPos.y, nextPos.x, nextDir, currentNodeDiscovery.direction)

				if nextPos.isInvalid(mapTiles2d) {
					fmt.Printf("Invalid\n")
					continue initDirLoop
				}

				if inverseDirs[currentNodeDiscovery.direction] == nextDir {
					// Can't go back
					fmt.Printf("Back %s\n", inverseDirs[currentNodeDiscovery.direction])
					//continue initDirLoop
				}

				nextTile := mapTiles2d[nextPos.y][nextPos.x]

				if !nextTile.isPath {
					// Can't go to a non-path.
					fmt.Printf("Non-path %v %v\n", nextPos, nextTile)
					continue initDirLoop
				}

				nextCandidateDirs, tileHasDirs := validDirs[nextTile.tileChar]
				if !tileHasDirs || ignoreTiles {
					nextCandidateDirs = allDirs
				}
				if len(nextCandidateDirs) == 1 {
					isOneWay = true
				}

				nextValidDirs := []string{}
				for _, nextTestDir := range nextCandidateDirs {
					if nextTestDir == inverseDirs[nextDir] {
						continue
					}

					candidatePos := nextPos.move(nextTestDir, 1)

					if candidatePos.isInvalid(mapTiles2d) {
						continue
					}

					candidateTile := mapTiles2d[candidatePos.y][candidatePos.x]

					if !candidateTile.isPath {
						continue
					}

					fmt.Printf("Appending valid dirs: np %v cp %v ct %v nextdir %s from %s\n", nextPos, candidatePos, candidateTile, nextTestDir, nextDir)
					nextValidDirs = append(nextValidDirs, nextTestDir)
				}

				numCandidateDirs := len(nextValidDirs)
				if numCandidateDirs == 1 {
					history = append(history, nextPos)
					nextSteps += 1
					nextPos = nextPos.move(nextValidDirs[0], 1)
					nextDir = nextValidDirs[0]
					fmt.Printf("nextpos %v nextDir %v | cands %v\n", nextPos, nextDir, nextValidDirs)
				} else if numCandidateDirs == 2 || numCandidateDirs == 3 || (numCandidateDirs == 0 && nextPos.y == endingNode.pos.y && nextPos.x == endingNode.pos.x) {
					fmt.Printf("Node nextvaliddirs: %v | %v\n", nextPos, nextValidDirs)
					// Found a node
					nextNodeNum := -1
					currentNodeNum := -1
					for ndeI, nde := range foundNodes {
						if nde.pos == nextPos {
							nextNodeNum = ndeI
						}
						if nde.pos == currentNode.pos {
							currentNodeNum = ndeI
						}
					}
					if currentNodeNum > -1 {
						if nextNodeNum > -1 {
							// Found an existing node.
							nextNodeEdge := GraphEdge{inverseDirs[nextDir], nextSteps, currentNode.pos, currentNodeNum, history}
							edgeExists := false
							for _, ege := range foundNodes[nextNodeNum].edges {
								if ege.direction == nextNodeEdge.direction {
									edgeExists = true
								}
							}
							if !edgeExists && !isOneWay {
								// Add edge back for revese direction.
								foundNodes[nextNodeNum].edges = append(foundNodes[nextNodeNum].edges, nextNodeEdge)
							}
							fmt.Printf("Connecting to existing node %d at %d,%d\n", nextNodeNum, foundNodes[nextNodeNum].pos.y, foundNodes[nextNodeNum].pos.x)

						} else {
							nextNodeNum = len(foundNodes)

							nextNodeEdges := []GraphEdge{{inverseDirs[nextDir], nextSteps, currentNode.pos, currentNodeNum, history}}
							if isOneWay {
								// If one way, there is no reverse edge to come back.
								nextNodeEdges = []GraphEdge{}
							}
							nextNodeDis := GraphNodeDiscovery{GraphNode{nextPos, nextNodeEdges}, nextDir}
							foundNodes = append(foundNodes, nextNodeDis.node)
							pendingNodes = append(pendingNodes, nextNodeDis)
							fmt.Printf("Adding node at %d,%d (%d)\n", nextPos.y, nextPos.x, len(foundNodes))
						}

						// Add edge to next node
						curNodeEdge := GraphEdge{initDir, nextSteps, nextPos, nextNodeNum, history}
						edgeExists := false
						for _, ege := range foundNodes[currentNodeNum].edges {
							if ege.direction == curNodeEdge.direction {
								edgeExists = true
							}
						}
						if !edgeExists {
							foundNodes[currentNodeNum].edges = append(foundNodes[currentNodeNum].edges, curNodeEdge)
						}
					}
					nextPosNotNode = false
				} else {
					fmt.Printf("Bad nextvaliddirs: %v | %v\n", nextPos, nextValidDirs)
					break
				}
			}
		}

	}

	fmt.Printf("Total found: %d\n", len(foundNodes))
	for _, foundN := range foundNodes {
		fmt.Printf("@ %d,%d | ", foundN.pos.y, foundN.pos.x)
		for _, foundE := range foundN.edges {
			fmt.Printf(" (%s %d %v)", foundE.direction, foundE.weight, foundE.connectionPoint)
		}
		fmt.Printf("\n")
	}

	return foundNodes
}

func printVisitedGraph(mapTiles2d [][]MapTile, visitedPoints []Point) {

	visitedData := [][]bool{}
	for _, graphLine := range mapTiles2d {
		vr := []bool{}
		for range graphLine {
			vr = append(vr, false)
		}
		visitedData = append(visitedData, vr)
	}
	for _, vp := range visitedPoints {
		visitedData[vp.y][vp.x] = true
	}

	for _, graph_line := range mapTiles2d {
		for _, graph_nde := range graph_line {
			if visitedData[graph_nde.pos.y][graph_nde.pos.x] {
				if graph_nde.isPath {
					if graph_nde.tileChar == "." {
						fmt.Printf("O")
					} else {
						fmt.Printf("%s", graph_nde.tileChar)
					}
				} else {
					fmt.Printf("!")
				}
			} else {
				//fmt.Printf(" ")
				fmt.Printf("\u001b[31m%s\u001b[0m", graph_nde.tileChar)
			}
		}
		fmt.Printf("\n")
	}
}

func stepsStringify(steps []Point) string {
	output := ""
	for _, step := range steps {
		output += fmt.Sprintf(";%d,%d", step.y, step.x)
	}
	return output
}

func stepsParse(stepsStr string) []Point {
    output := []Point{}

    for _,stepStr := range strings.Split(stepsStr, ";") {
        stepSplit := strings.Split(stepStr, ",")
        if len(stepSplit) == 2 {
            stepY, stepYParsed := strconv.ParseInt(stepSplit[1], 10, 16)
            stepX, stepXParsed := strconv.ParseInt(stepSplit[0], 10, 16)

            if stepYParsed != nil && stepXParsed != nil {
                output = append(output, Point{int16(stepY), int16(stepX)})
            }
        }
    }

    return output
}

func main() {
	var filename = flag.String("f", "../inputs/d23.sample1.txt", "file to use")
	var part2 = flag.Bool("part2", false, "do part 2")
	var visualise = flag.String("visualise", "", "check a string")
	flag.Parse()
	dat, err := os.ReadFile(*filename)
	if err != nil {
		log.Fatalf("unable to read file: %f", err)
	}

	mapTiles2d := [][]MapTile{}

	pointStart := Point{-1, -1}
	pointEnd := Point{-1, -1}

	// Parse the data.
	for l_num, f_line := range strings.Split(strings.Replace(string(dat), "\r\n", "\n", -1), "\n") {
		if len(f_line) == 0 {
			continue
		}
		f_line = strings.Trim(f_line, " \n")

		graph_node_row := []MapTile{}
		for c_num := 0; c_num < len(f_line); c_num += 1 {
			node_char := f_line[c_num]
			isPath := node_char != '#'
			if isPath {
				pointEnd = Point{int16(l_num), int16(c_num)}
				if l_num == 0 {
					pointStart = Point{int16(l_num), int16(c_num)}
				}
			}
			graph_node_row = append(graph_node_row, MapTile{Point{int16(l_num), int16(c_num)}, isPath, string(node_char), false, -1, []Point{}, false})
		}
		mapTiles2d = append(mapTiles2d, graph_node_row)
	}

	if len(mapTiles2d) == 0 || len(mapTiles2d[0]) == 0 {
		fmt.Fprintf(os.Stderr, "Error when reading file.\n")
		return
	}
	if pointStart.x < 0 || pointStart.y < 0 {
		fmt.Fprintf(os.Stderr, "Start point not found.\n")
		return
	}

	getGraphNodes(mapTiles2d, GraphNode{pointStart, []GraphEdge{}}, GraphNode{pointEnd, []GraphEdge{}}, *part2)
	numNodes2 := 0
	numNodes3 := 0
	for _, graphLine := range mapTiles2d {
		for _, graphColumn := range graphLine {
			dirs := []string{"U", "R", "D", "L"}
			countPaths := 0
			if !graphColumn.isPath {
				continue
			}
			for _, dir := range dirs {
				testPos := graphColumn.pos.move(dir, 1)
				if testPos.x < 0 || testPos.y < 0 || int(testPos.x) >= len(graphLine) || int(testPos.y) >= len(mapTiles2d) {
					continue
				}
				testItem := mapTiles2d[testPos.y][testPos.x]
				if testItem.isPath {
					countPaths += 1
				}
			}
			if countPaths > 2 {
				fmt.Printf("%d,%d: %d\n", graphColumn.pos.y, graphColumn.pos.x, countPaths)
				mapTiles2d[graphColumn.pos.y][graphColumn.pos.x].isNode = true
				numNodes2 += 1
			}
			if countPaths > 3 {
				numNodes3 += 1
			}
		}
	}

	fmt.Printf("Num nodes >2: %d, >3: %d\n", numNodes2, numNodes3)
	for _, graph_line := range mapTiles2d {
		for _, graph_nde := range graph_line {
			if graph_nde.isNode {
				fmt.Printf("X")
				//fmt.Printf("X", graph_nde.tileChar)
			} else {
				//fmt.Printf(" ")
				if graph_nde.isPath {
					fmt.Printf("\u001b[31m%s\u001b[0m", graph_nde.tileChar)
				} else {
					fmt.Printf(" ")
				}
			}
		}
		fmt.Printf("\n")
	}
	if true {
		return
	}

	fmt.Printf("Start point: %v\n", pointStart)
	fmt.Printf("End point: %v\n", pointEnd)
	mapTiles2d[pointStart.y][pointStart.x].highest_steps = 0

	max_y := len(mapTiles2d) - 1
	max_x := len(mapTiles2d[0]) - 1

	starting_node := PendingWalk{pointStart, 0, []Point{pointStart}}

	total := -1

	pending_positions := make(LongestWalk, 0)
	heap.Init(&pending_positions)
	heap.Push(&pending_positions, starting_node)

	valid_dirs := map[string][]string{".": {"U", "R", "D", "L"}, "^": {"U"}, ">": {"R"}, "v": {"D"}, "<": {"L"}}

	visualiseSteps := []Point{}

	loop_num := 0
	if len(*visualise) == 0 {
		for pending_positions.Len() > 0 {
			loop_num += 1

			current_position := heap.Pop(&pending_positions).(PendingWalk)

			current_node := mapTiles2d[current_position.pos.y][current_position.pos.x]
			//fmt.Printf("Cur: %v | %v\n", current_position, current_node)
			//fmt.Printf("Len: %d, %d\n", pending_positions.Len(), current_position.steps)
			//fmt.Printf("Len: %d, %d, %v\n", pending_positions.Len(), current_position.steps, nextPosition.history)

			if !current_node.isPath {
				continue
			}

			if current_node.pos.x == pointEnd.x && current_node.pos.y == pointEnd.y {
				fmt.Printf("Ending: %d | %v\n", current_node.highest_steps+15, current_node.pos)
				continue
			}

			try_dirs, is_valid_dir := valid_dirs[current_node.tileChar]
			if !is_valid_dir || *part2 {
				try_dirs = []string{"R", "D", "L", "U"}
			}
			for _, dir := range try_dirs {
				nextPosition := moveDirection(dir, current_position)

				if nextPosition.pos.x < 0 || nextPosition.pos.y < 0 || nextPosition.pos.x > int16(max_x) || nextPosition.pos.y > int16(max_y) {
					// Out of bounds
					continue
				}

				nextNode := mapTiles2d[nextPosition.pos.y][nextPosition.pos.x]

				if nextNode.visited {
					// Already visited.
					//continue
				}

				if !nextNode.isPath {
					// Not a path, can't traverse.
					continue
				}

				isVisited := false
				for _, prevPos := range current_position.history {
					if prevPos.x == nextNode.pos.x && prevPos.y == nextNode.pos.y {
						isVisited = true
						break
					}
				}
				if isVisited {
					continue
				}
				//fmt.Printf("Next3: %s | %v | %v\n", dir, nextPosition, nextNode)

				if nextPosition.steps > nextNode.highest_steps {
					nextNode.highest_steps = nextPosition.steps
					nextNode.highest_steps_path = nextPosition.history
					//fmt.Printf("New highest: %d,%d: %d %d\n", nextNode.pos.y, nextNode.pos.x, nextNode.highest_steps, nextNode.highest_steps)
					mapTiles2d[nextPosition.pos.y][nextPosition.pos.x] = nextNode

				} else {
					//fmt.Printf("Not highest: %d,%d: %d %d\n", nextNode.pos.y, nextNode.pos.x, nextNode.highest_steps, nextNode.highest_steps)
				}

				heap.Push(&pending_positions, nextPosition)

			}
		}

		visualiseSteps = mapTiles2d[pointEnd.y][pointEnd.x].highest_steps_path
		/*
			for i_y, graph_line := range graph_nodes {
				for i_x, graph_nde := range graph_line {
					if graph_nde.lowest_cost < 10_000 {
						fmt.Printf("%d:%3d|", graph_nde.cost, graph_nde.lowest_cost)
					} else {
						fmt.Printf("%d:   |", graph_nde.cost)
					}
					graph_nodes[i_y][i_x].visited_v = false
					graph_nodes[i_y][i_x].visited_h = false
				}
				fmt.Printf("\n")
			}
		*/

		total = mapTiles2d[pointEnd.y][pointEnd.x].highest_steps
		//*visualise = graph_nodes[ending_node.y][ending_node.x].highest_steps_path

		fmt.Printf("T: %d\n", total)

        fmt.Printf("Path: %s\n", stepsStringify(visualiseSteps))
	}

	if len(*visualise) > 0 {
        visualiseSteps = stepsParse(*visualise)
	}

    printVisitedGraph(mapTiles2d, visualiseSteps)
}
