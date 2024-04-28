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
	nodeIndex          int
}

type GraphEdge struct {
	direction       string
	weight          int
	connectionPoint Point
	connectionIndex int
	connectionPath  []Point
}

type GraphNode struct {
	pos              Point
	edges            []GraphEdge
	highestSteps     int
	highestStepsPath []Point
}

type GraphNodeDiscovery struct {
	node      GraphNode
	direction string
}

type PendingWalk struct {
	pos          Point
	steps        int
	history      []Point
	historyNodes []int
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
	return PendingWalk{nextPos, currentPos.steps + 1, nextHistory, currentPos.historyNodes}
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

		//fmt.Printf("Node %d,%d from %s\n", currentNode.pos.y, currentNode.pos.x, currentNodeDiscovery.direction)

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
				//fmt.Printf("Next pos: %d,%d | %s | %s\n", nextPos.y, nextPos.x, nextDir, currentNodeDiscovery.direction)

				if nextPos.isInvalid(mapTiles2d) {
					continue initDirLoop
				}

				if inverseDirs[currentNodeDiscovery.direction] == nextDir {
					// Can't go back
					//continue initDirLoop
				}

				nextTile := mapTiles2d[nextPos.y][nextPos.x]

				if !nextTile.isPath {
					// Can't go to a non-path.
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
					nextValidDirs = append(nextValidDirs, nextTestDir)
				}

				numCandidateDirs := len(nextValidDirs)
				if numCandidateDirs == 1 {
					history = append(history, nextPos)
					nextSteps += 1
					nextPos = nextPos.move(nextValidDirs[0], 1)
					nextDir = nextValidDirs[0]
				} else if numCandidateDirs == 2 || numCandidateDirs == 3 || (numCandidateDirs == 0 && nextPos.y == endingNode.pos.y && nextPos.x == endingNode.pos.x) {
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

						} else {
							nextNodeNum = len(foundNodes)

							nextNodeEdges := []GraphEdge{{inverseDirs[nextDir], nextSteps, currentNode.pos, currentNodeNum, history}}
							if isOneWay {
								// If one way, there is no reverse edge to come back.
								nextNodeEdges = []GraphEdge{}
							}
							nextNodeDis := GraphNodeDiscovery{GraphNode{nextPos, nextNodeEdges, -1, []Point{}}, nextDir}
							foundNodes = append(foundNodes, nextNodeDis.node)
							pendingNodes = append(pendingNodes, nextNodeDis)
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
					break
				}
			}
		}

	}

	fmt.Printf("Total found: %d\n", len(foundNodes))
	for foundI, foundN := range foundNodes {
		fmt.Printf("%d @ %d,%d | ", foundI, foundN.pos.y, foundN.pos.x)
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
		output += fmt.Sprintf("|%d,%d", step.y, step.x)
	}
    return output[1:]
}

func stepsParse(stepsStr string) []Point {
	output := []Point{}

	for _, stepStr := range strings.Split(stepsStr, "|") {
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

func countNodesManually(mapTiles2d [][]MapTile) {
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
				mapTiles2d[graphColumn.pos.y][graphColumn.pos.x].nodeIndex = numNodes2
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
			if graph_nde.nodeIndex > -1 {
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
			graph_node_row = append(graph_node_row, MapTile{Point{int16(l_num), int16(c_num)}, isPath, string(node_char), false, -1, []Point{}, -1})
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

	graphNodes := getGraphNodes(mapTiles2d, GraphNode{pointStart, []GraphEdge{}, 0, []Point{}}, GraphNode{pointEnd, []GraphEdge{}, -1, []Point{}}, *part2)
	for i, nde := range graphNodes {
		mapTiles2d[nde.pos.y][nde.pos.x].nodeIndex = i
	}

	fmt.Printf("Start point: %v\n", pointStart)
	fmt.Printf("End point: %v\n", pointEnd)
	mapTiles2d[pointStart.y][pointStart.x].highest_steps = 0

	//max_y := len(mapTiles2d) - 1
	//max_x := len(mapTiles2d[0]) - 1

	starting_node := PendingWalk{pointStart, 0, []Point{pointStart}, []int{0}}

	total := -1

	pending_positions := make(LongestWalk, 0)
	heap.Init(&pending_positions)
	heap.Push(&pending_positions, starting_node)

	visualiseSteps := []Point{}

	loop_num := 0
	if len(*visualise) == 0 {
		for pending_positions.Len() > 0 {
			loop_num += 1

			currentWalk := heap.Pop(&pending_positions).(PendingWalk)

			currentTile := mapTiles2d[currentWalk.pos.y][currentWalk.pos.x]
			currentNode := graphNodes[currentTile.nodeIndex]

			//fmt.Printf("Current: w %v t %v n %v\n", currentWalk, currentTile, currentNode)

			if currentNode.pos.x == pointEnd.x && currentNode.pos.y == pointEnd.y {
				continue
			}

			for _, edge := range currentNode.edges {
                //fmt.Printf("Edge: %d -> %d\n", currentTile.nodeIndex, edge.connectionIndex)
				nextNode := graphNodes[edge.connectionIndex]
				nextSteps := currentWalk.steps + edge.weight
                nextHistoryPath := []Point{}
				nextHistoryPath = append(nextHistoryPath, currentWalk.history...)
				nextHistoryPath = append(nextHistoryPath, edge.connectionPath...)
				nextHistoryPath = append(nextHistoryPath, edge.connectionPoint)
                nextHistoryNode := []int{}
				nextHistoryNode = append(nextHistoryNode, currentWalk.historyNodes...)
				nextHistoryNode = append(nextHistoryNode, edge.connectionIndex)


				isVisited := false
				for _, prevNodeIndex := range currentWalk.historyNodes {
                    //fmt.Printf("Dupe? %d %d\n", prevNodeIndex, edge.connectionIndex)
					if prevNodeIndex == edge.connectionIndex {
						isVisited = true
						break
					}
				}
				if isVisited {
					continue
                }

				if nextSteps > nextNode.highestSteps {
					nextNode.highestSteps = nextSteps
                    if nextNode.pos.x == pointEnd.x && nextNode.pos.y == pointEnd.y {
                        // Save nodes for ending point.
                        // No need to save them along the way.
                        nextNode.highestStepsPath = nextHistoryPath
                        fmt.Printf("Ending: %d %v (ln %d)\n", nextSteps, nextHistoryNode, loop_num)
                    }
					graphNodes[edge.connectionIndex] = nextNode
				}

				heap.Push(&pending_positions, PendingWalk{nextNode.pos, nextSteps, nextHistoryPath, nextHistoryNode})

			}
            if loop_num % 1_000_000 == 0 {
                fmt.Printf("Len: %d %d %d %v\n", len(pending_positions), currentWalk.steps, loop_num, currentWalk.historyNodes)
            }
		}
        fmt.Printf("Final loop count: %d\n", loop_num)

		endIndex := mapTiles2d[pointEnd.y][pointEnd.x].nodeIndex
		visualiseSteps = graphNodes[endIndex].highestStepsPath

		total = graphNodes[endIndex].highestSteps

		fmt.Printf("T: %d\n", total)

		fmt.Printf("Path: %s\n", stepsStringify(visualiseSteps))
	}

	if len(*visualise) > 0 {
		visualiseSteps = stepsParse(*visualise)
	}

	printVisitedGraph(mapTiles2d, visualiseSteps)
}
