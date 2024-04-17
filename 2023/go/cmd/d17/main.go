package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"slices"
	"strings"
)

type GraphNode struct {
	y                  int
	x                  int
	cost               int
	visited            bool
	isPending          bool
	lowest_cost        int
	lowest_cost_path   string
	lowest_cost_h      int
	lowest_cost_path_h string
	lowest_cost_v      int
	lowest_cost_path_v string
}

type CruciblePosition struct {
	direction   string
	history     string
	cost        int
	y           int
	x           int
	is_starting bool
}

func generatePositionKey(current_position CruciblePosition) string {
    var current_position_key string
    num_in_direction := 0
    if len(current_position.history) > 0 {
        for i := len(current_position.history)-1; i >= 0; i -= 1 {
            if current_position.direction[0] == current_position.history[i] {
                num_in_direction += 1
            }
        }
    }
    current_position_key = fmt.Sprintf("%s,%d,%d,%d", current_position.direction, num_in_direction, current_position.y, current_position.x)
    /*
    if current_position.direction == "N" || current_position.direction == "S" {
        current_position_key = fmt.Sprintf("V,%d,%d", current_position.y, current_position.x)
    } else {
        current_position_key = fmt.Sprintf("H,%d,%d", current_position.y, current_position.x)
    }
    */
    return current_position_key
}

func moveDirection(direction string, pos CruciblePosition) CruciblePosition {
	next_history := pos.history + string(direction)
	if direction == "N" {
		return CruciblePosition{"N", next_history, pos.cost, pos.y - 1, pos.x, false}
	} else if direction == "E" {
		return CruciblePosition{"E", next_history, pos.cost, pos.y, pos.x + 1, false}
	} else if direction == "S" {
		return CruciblePosition{"S", next_history, pos.cost, pos.y + 1, pos.x, false}
	} else if direction == "W" {
		return CruciblePosition{"W", next_history, pos.cost, pos.y, pos.x - 1, false}
	} else {
		return CruciblePosition{"0", next_history, pos.cost, pos.y, pos.x, false}
	}
}

func main() {
	var filename = flag.String("f", "../inputs/d17.sample.txt", "file to use")
	var part2 = flag.Bool("part2", false, "do part 2")
	var visualise = flag.String("visualise", "", "check a string")
	flag.Parse()
	dat, err := os.ReadFile(*filename)
	if err != nil {
		log.Fatalf("unable to read file: %f", err)
	}

	graph_nodes := [][]GraphNode{}

	// Parse the data.
	for l_num, f_line := range strings.Split(strings.Replace(string(dat), "\r\n", "\n", -1), "\n") {
		if len(f_line) == 0 {
			continue
		}
		f_line = strings.Trim(f_line, " \n")

		graph_node_row := []GraphNode{}
		for c_num := 0; c_num < len(f_line); c_num += 1 {
			node_weight := f_line[c_num] - '0'
			graph_node_row = append(graph_node_row, GraphNode{l_num, c_num, int(node_weight), false, false, math.MaxInt64, "", math.MaxInt64, "", math.MaxInt64, ""})
		}
		graph_nodes = append(graph_nodes, graph_node_row)
	}

	if len(graph_nodes) == 0 || len(graph_nodes[0]) == 0 {
		fmt.Fprintf(os.Stderr, "Error when reading file.\n")
		return
	}
	graph_nodes[0][0].lowest_cost_h = 0
	graph_nodes[0][0].lowest_cost_v = 0

	max_y := len(graph_nodes) - 1
	max_x := len(graph_nodes[0]) - 1

	starting_node := CruciblePosition{"0", "", 0, 0, 0, true}
	ending_node := CruciblePosition{"0", "", 0, max_y, max_x, false}

	total := math.MaxInt64

	min_before_turn := 1
	max_before_turn := 3
	if *part2 {
		min_before_turn = 4
		max_before_turn = 10
	}

	pending_positions := []CruciblePosition{}
	pending_positions = append(pending_positions, starting_node)

	pending_pos_map := map[string]bool{}

	valid_dirs := map[string][]string{"N": {"E", "W"}, "S": {"E", "W"}, "E": {"N", "S"}, "W": {"N", "S"}}

	loop_num := 0
	if len(*visualise) == 0 {
		for len(pending_positions) > 0 {
			loop_num += 1

			current_position := pending_positions[0]
			pending_positions = pending_positions[1:]

			//current_position_key := fmt.Sprintf("%v", current_position)
			current_position_key := generatePositionKey(current_position)

			graph_nodes[current_position.y][current_position.x].isPending = false
			pending_pos_map[current_position_key] = false
			current_node := graph_nodes[current_position.y][current_position.x]

			try_dirs, is_valid_dir := valid_dirs[current_position.direction]
			if !is_valid_dir {
				try_dirs = []string{"E", "S", "N", "W"}
			}
			for _, dir := range try_dirs {
				//next_cost := current_node.lowest_cost
				next_cost := current_position.cost

				next_position := current_position
				next_position.direction = dir
				for i := 0; i < max_before_turn; i += 1 {
					next_position = moveDirection(dir, next_position)

					isPosRangeBad := next_position.y >= len(graph_nodes)
					isPosRangeBad = isPosRangeBad || next_position.x >= len(graph_nodes[0])
					isPosRangeBad = isPosRangeBad || next_position.y < 0
					isPosRangeBad = isPosRangeBad || next_position.x < 0
					isPosRangeBad = isPosRangeBad || next_position.direction == "0"
					if isPosRangeBad {
						// Out of range, it expires.
						break
					}

					next_node := &graph_nodes[next_position.y][next_position.x]
					next_cost += next_node.cost
					// If minimum distance before turning has been reached
					if (i + 1) >= min_before_turn {
						next_position.cost = next_cost

						if next_node.visited {
							continue
						}

                        next_position_key := generatePositionKey(next_position)
						next_node_pending, next_node_pending_found := pending_pos_map[next_position_key]
						if !next_node_pending_found {
							next_node_pending = false
						}

						var node_lowest_cost int
						if next_position.direction == "N" || next_position.direction == "S" {
							node_lowest_cost = next_node.lowest_cost_v
						} else {
							node_lowest_cost = next_node.lowest_cost_h
						}
						if next_cost < node_lowest_cost {
							next_node.lowest_cost = next_cost
							if next_position.direction == "N" || next_position.direction == "S" {
								next_node.lowest_cost_v = next_cost
							} else {
								next_node.lowest_cost_h = next_cost
							}

                            if next_node_pending {
                                // If pending, replace the position.
                                for pos_i := 0; pos_i < len(pending_positions); pos_i += 1 {
                                    candidate_pos := pending_positions[pos_i]
                                    candidate_pos_key := generatePositionKey(candidate_pos)
                                    if next_position_key == candidate_pos_key {
                                        pending_positions[pos_i] = next_position
                                    }
                                }
                            }

							if next_node.y == ending_node.y && next_node.x == ending_node.x {
								next_node.lowest_cost_path = next_position.history
								if next_position.direction == "N" || next_position.direction == "S" {
									next_node.lowest_cost_path_v = next_position.history
								} else {
									next_node.lowest_cost_path_h = next_position.history
								}
								fmt.Printf("[%d] Candidate: %d | %s\n", i, next_cost, next_position.history)
							}
						}

						if next_position.direction == "N" || next_position.direction == "S" {
							if next_node.lowest_cost_v < math.MaxInt64 && !next_node_pending {
								pending_positions = append(pending_positions, next_position)
								pending_pos_map[next_position_key] = true
							}
						} else {
							if next_node.lowest_cost_h < math.MaxInt64 && !next_node_pending {
								pending_positions = append(pending_positions, next_position)
								pending_pos_map[next_position_key] = true
							}
						}
					}
				}
			}
			graph_nodes[current_position.y][current_position.x].visited = true

			if current_node.y == ending_node.y && current_node.x == ending_node.x {
                current_node.lowest_cost = current_node.lowest_cost_v
                current_node.lowest_cost_path = current_node.lowest_cost_path_v
                if current_node.lowest_cost > current_node.lowest_cost_h {
                    current_node.lowest_cost = current_node.lowest_cost_h
                    current_node.lowest_cost_path = current_node.lowest_cost_path_h
                }
                graph_nodes[current_position.y][current_position.x] = current_node

				fmt.Printf("Finished %v\n", current_node)
				break
			}

			slices.SortFunc(pending_positions, func(a CruciblePosition, b CruciblePosition) int {
				if a.cost == b.cost {
					return 0
				} else if a.cost < b.cost {
					return -1
				} else {
					return 1
				}
			})
		}

		for i_y, graph_line := range graph_nodes {
			for i_x, graph_nde := range graph_line {
				if graph_nde.lowest_cost < 10_000 {
					fmt.Printf("%d:%3d|", graph_nde.cost, graph_nde.lowest_cost)
				} else {
					fmt.Printf("%d:   |", graph_nde.cost)
				}
				graph_nodes[i_y][i_x].visited = false
			}
			fmt.Printf("\n")
		}

		total = graph_nodes[ending_node.y][ending_node.x].lowest_cost
		*visualise = graph_nodes[ending_node.y][ending_node.x].lowest_cost_path

		fmt.Printf("T: %d\n", total)
	}

	if len(*visualise) > 0 {
		vis_node := starting_node
		vis_total := 0
		for _, dir := range *visualise {
			vis_node = moveDirection(string(dir), vis_node)
			if vis_node.x >= 0 && vis_node.y >= 0 && vis_node.y < len(graph_nodes) && vis_node.x < len(graph_nodes[0]) {
				graph_nodes[vis_node.y][vis_node.x].visited = true
				add_cost := graph_nodes[vis_node.y][vis_node.x].cost
				fmt.Printf("CT: %4d + %4d = %4d\n", vis_total, add_cost, vis_total+add_cost)
				vis_total += add_cost
			} else {
				fmt.Printf("ERROR: %v\n", vis_node)
			}
		}
	}

	for _, graph_line := range graph_nodes {
		for _, graph_nde := range graph_line {
			if graph_nde.visited {
				fmt.Printf("%d", graph_nde.cost)
			} else {
				//fmt.Printf(" ")
				fmt.Printf("\u001b[31m%d\u001b[0m", graph_nde.cost)
			}
		}
		fmt.Printf("\n")
	}
}
