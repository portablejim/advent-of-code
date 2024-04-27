package main

import (
	"container/heap"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type Point struct {
	y int16
	x int16
}

func (p Point) move(dir string, num int16) Point {
	return move(dir, p, num)
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

type MapNode struct {
	pos                Point
	isPath             bool
	tileChar           string
	visited            bool
	highest_steps      int
	highest_steps_path []Point
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

func main() {
	var filename = flag.String("f", "../inputs/d23.sample1.txt", "file to use")
	var part2 = flag.Bool("part2", false, "do part 2")
	var visualise = flag.String("visualise", "", "check a string")
	flag.Parse()
	dat, err := os.ReadFile(*filename)
	if err != nil {
		log.Fatalf("unable to read file: %f", err)
	}

	graph_nodes := [][]MapNode{}

	pointStart := Point{-1, -1}
	pointEnd := Point{-1, -1}

	// Parse the data.
	for l_num, f_line := range strings.Split(strings.Replace(string(dat), "\r\n", "\n", -1), "\n") {
		if len(f_line) == 0 {
			continue
		}
		f_line = strings.Trim(f_line, " \n")

		graph_node_row := []MapNode{}
		for c_num := 0; c_num < len(f_line); c_num += 1 {
			node_char := f_line[c_num]
			isPath := node_char != '#'
			if isPath {
				pointEnd = Point{int16(l_num), int16(c_num)}
				if l_num == 0 {
					pointStart = Point{int16(l_num), int16(c_num)}
				}
			}
			graph_node_row = append(graph_node_row, MapNode{Point{int16(l_num), int16(c_num)}, isPath, string(node_char), false, -1, []Point{}})
		}
		graph_nodes = append(graph_nodes, graph_node_row)
	}

	if len(graph_nodes) == 0 || len(graph_nodes[0]) == 0 {
		fmt.Fprintf(os.Stderr, "Error when reading file.\n")
		return
	}
	if pointStart.x < 0 || pointStart.y < 0 {
		fmt.Fprintf(os.Stderr, "Start point not found.\n")
		return
	}
	fmt.Printf("Start point: %v\n", pointStart)
	fmt.Printf("End point: %v\n", pointEnd)
	graph_nodes[pointStart.y][pointStart.x].highest_steps = 0

	max_y := len(graph_nodes) - 1
	max_x := len(graph_nodes[0]) - 1

	starting_node := PendingWalk{pointStart, 0, []Point{pointStart}}

	total := -1

	pending_positions := make(LongestWalk, 0)
	heap.Init(&pending_positions)
	heap.Push(&pending_positions, starting_node)

	valid_dirs := map[string][]string{".": {"U", "R", "D", "L"}, "^": {"U"}, ">": {"R"}, "v": {"D"}, "<": {"L"}}

	loop_num := 0
	if len(*visualise) == 0 {
		for pending_positions.Len() > 0 {
			loop_num += 1

			current_position := heap.Pop(&pending_positions).(PendingWalk)

			current_node := graph_nodes[current_position.pos.y][current_position.pos.x]
			//fmt.Printf("Cur: %v | %v\n", current_position, current_node)

			if !current_node.isPath {
				continue
			}

			if current_node.pos.x == pointEnd.x && current_node.pos.y == pointEnd.y {
				fmt.Printf("Ending: %d\n", current_node.highest_steps)
				continue
			}

			try_dirs, is_valid_dir := valid_dirs[current_node.tileChar]
			if !is_valid_dir {
				try_dirs = []string{"R", "D", "L", "U"}
			}
			for _, dir := range try_dirs {
				nextPosition := moveDirection(dir, current_position)

				if nextPosition.pos.x < 0 || nextPosition.pos.y < 0 || nextPosition.pos.x > int16(max_x) || nextPosition.pos.y > int16(max_y) {
					// Out of bounds
					continue
				}

				nextNode := graph_nodes[nextPosition.pos.y][nextPosition.pos.x]

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
					graph_nodes[nextPosition.pos.y][nextPosition.pos.x] = nextNode
				}

				heap.Push(&pending_positions, nextPosition)
			}
		}

		for _, step := range graph_nodes[pointEnd.y][pointEnd.x].highest_steps_path {
			if graph_nodes[step.y][step.x].visited {
				fmt.Printf("Error: %v\n", step)
			} else {
				graph_nodes[step.y][step.x].visited = true
			}
		}
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

		total = graph_nodes[pointEnd.y][pointEnd.x].highest_steps
		//*visualise = graph_nodes[ending_node.y][ending_node.x].highest_steps_path

		fmt.Printf("T: %d\n", total)
	}

	if len(*visualise) > 0 {
		/*
			vis_node := starting_node
			vis_total := 0
			for _, dir := range *visualise {
				vis_node = moveDirection(string(dir), vis_node)
				if vis_node.x >= 0 && vis_node.y >= 0 && vis_node.y < int16(len(graph_nodes)) && vis_node.x < int16(len(graph_nodes[0])) {
					graph_nodes[vis_node.y][vis_node.x].visited_v = true
					if vis_node.direction == "N" || vis_node.direction == "S" {
						graph_nodes[vis_node.y][vis_node.x].visited_v = true
					} else {
						graph_nodes[vis_node.y][vis_node.x].visited_h = true
					}
					add_cost := graph_nodes[vis_node.y][vis_node.x].cost
					fmt.Printf("CT: %4d + %4d = %4d\n", vis_total, add_cost, vis_total+add_cost)
					vis_total += add_cost
				} else {
					fmt.Printf("ERROR: %v\n", vis_node)
				}
			}
		*/
	}

	if *part2 {
	}

	for _, graph_line := range graph_nodes {
		for _, graph_nde := range graph_line {
			if graph_nde.visited {
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
