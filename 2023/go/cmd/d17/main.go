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
	y                int
	x                int
	cost             int
	visited          bool
	isPending        bool
	lowest_cost      int
	lowest_cost_path string
}

type Coord struct {
	y int
	x int
}

func (c Coord) move(direction string) Coord {
	switch direction {
	case "N":
		c.y -= 1
	case "S":
		c.y += 1
	case "E":
		c.x -= 1
	case "W":
		c.x += 1
	}
	return c
}

type CruciblePosition struct {
	direction   string
	history     string
	cost        int
	y           int
	x           int
	is_starting bool
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

func getNextDirections(direction rune, previous_directions string, min_before_turn int, max_before_turn int) []rune {
	output := []rune{}

	// Not enough directions exist yet before turning.
	if len(previous_directions) < min_before_turn {
		//fmt.Printf("Forced straight 1: %d < %d\n", len(previous_directions), min_before_turn)
		output = append(output, direction)
		return output
	}

	// Count num taken in each direction.
	count_same_direction := 0
	for i := len(previous_directions) - 1; i >= 0 && count_same_direction <= (max_before_turn+1); i -= 1 {
		if previous_directions[i] != byte(direction) {
			break
		}
		count_same_direction += 1
	}

	// Not enough of the same direction to turn.
	if count_same_direction < min_before_turn {
		output = append(output, direction)
		//fmt.Printf("Forced straight 2: %d < %d\n", count_same_direction, min_before_turn)
		return output
	}

	can_move_same := true
	if count_same_direction >= max_before_turn {
		//fmt.Printf("Forced turn: %d > %d\n", count_same_direction, max_before_turn)
		can_move_same = false
	}

	if can_move_same {
		output = append(output, direction)
	}
	if direction == 'N' || direction == 'S' {
		output = append(output, 'E', 'W')
	} else if direction == 'E' || direction == 'W' {
		output = append(output, 'N', 'S')
	}

	//fmt.Printf("Next dirs: %v\n", output)

	return output
}

func getLowestCost(s_i int, starting_nodes []CruciblePosition, ending_nodes_list []CruciblePosition, graph_nodes [][]GraphNode, min_before_turn int, max_before_turn int) (int, string) {
	ending_node := ending_nodes_list[s_i]

	pending_positions := []CruciblePosition{}
	pending_positions = append(pending_positions, starting_nodes...)

	// Handle the positions
	for len(pending_positions) > 0 {
		prev_position := pending_positions[0]
		prev_node := graph_nodes[prev_position.y][prev_position.x]
		pending_positions = pending_positions[1:]

		min_moves := 1
		if len(prev_position.history) == 0 || prev_position.direction != string(prev_position.history[len(prev_position.history)-1]) {
			min_moves = min_before_turn
		}
		out_of_range := false
		current_pos := prev_position
		// Skip over any positions where we can't turn.
		for i := 0; i < min_moves; i += 1 {
			prev_position = current_pos
			current_pos = moveDirection(prev_position.direction, prev_position)
			//fmt.Printf("Current pos: %v\n", current_pos)
			if current_pos.y >= len(graph_nodes) || current_pos.x >= len(graph_nodes[0]) || current_pos.y < 0 || current_pos.x < 0 || current_pos.direction == "0" {
				// Out of range, it expires.
				//fmt.Printf("expired r: %v\n", current_end)
				//fmt.Printf("Expire (invalid): %v\n", current_pos)
				out_of_range = true
				break
			}
			//fmt.Printf("Add cost %d + %d at %d,%d (%s)\n", current_pos.cost, graph_nodes[current_pos.y][current_pos.x].cost, current_pos.y, current_pos.x, current_pos.history)
			current_pos.cost += graph_nodes[current_pos.y][current_pos.x].cost
		}
		if out_of_range {
			continue
		}
		new_cost := current_pos.cost

		current_node := &graph_nodes[current_pos.y][current_pos.x]
		//fmt.Printf("N: %v %v %v %d\n", prev_node, current_pos, *current_node, new_cost)
		if new_cost < 0 || new_cost >= current_node.lowest_cost {
			// Higher cost, don't continue.
			//fmt.Printf("Expire (cost): %d > %d | %s | %s\n", new_cost, current_node.lowest_cost, current_pos.history, current_node.lowest_cost_path)
			continue
		}

		current_node.visited = true
		current_node.lowest_cost = new_cost
		current_node.lowest_cost_path = current_pos.history

		if current_pos.y == ending_node.y && current_pos.x == ending_node.x {
			// Found the end.
			//return current_node.lowest_cost, current_node.lowest_cost_path
			fmt.Printf("Possible lowest: %v %v %v\n", current_pos, *current_node, prev_node)
			continue
		}

		for _, next_dir := range getNextDirections(rune(current_pos.direction[0]), current_pos.history, min_before_turn, max_before_turn) {
			pending_positions = append(pending_positions, CruciblePosition{string(next_dir), current_pos.history, current_pos.cost, current_pos.y, current_pos.x, false})
		}
		//fmt.Printf("Next: %v\n", pending_positions)
		//fmt.Printf("Next: %d\n", len(pending_positions))

		/*
		   fmt.Printf("next: %v\n", laser_ends)
		   for _,floor_tile_line := range floor_tiles {
		       for _,floor_tile := range floor_tile_line {
		           if floor_tile.visited_count > 0 {
		               if floor_tile.visited_count > 1 {
		                   //fmt.Printf("█")
		                   fmt.Printf("#")
		               } else {
		                   fmt.Printf("#")
		               }
		           } else {
		               fmt.Printf(" ")
		           }
		       }
		       fmt.Printf("\n")
		   }
		*/
	}

	return graph_nodes[ending_node.y][ending_node.x].lowest_cost, graph_nodes[ending_node.y][ending_node.x].lowest_cost_path

	//return math.MaxInt64, "ERROR"
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
			graph_node_row = append(graph_node_row, GraphNode{l_num, c_num, int(node_weight), false, false, math.MaxInt64, ""})
		}
		graph_nodes = append(graph_nodes, graph_node_row)
	}

	if len(graph_nodes) == 0 || len(graph_nodes[0]) == 0 {
		fmt.Fprintf(os.Stderr, "Error when reading file.\n")
		return
	}
	graph_nodes[0][0].lowest_cost = 0
	//fmt.Printf("floor_tile_rows: %v\n", floor_tiles)

	max_y := len(graph_nodes) - 1
	max_x := len(graph_nodes[0]) - 1

	//starting_nodes_list := [][]CruciblePosition{{{'E', "", 0, 0, 0, true}, {'S', "", 0, 0, 0, true}}}
	//ending_nodes_list := []CruciblePosition{{'0', "", 0, max_y, max_x, false}}

	//starting_nodes := []CruciblePosition{{"E", "", 0, 0, 0, true}, {"S", "", 0, 0, 0, true}}
	starting_node := CruciblePosition{"0", "", 0, 0, 0, true}
	ending_node := CruciblePosition{"0", "", 0, max_y, max_x, false}

	total := math.MaxInt64
	first_total := math.MaxInt64

	min_before_turn := 1
	max_before_turn := 3
	if *part2 {
		min_before_turn = 4
		max_before_turn = 10
	}

	pending_positions := []CruciblePosition{}
	pending_positions = append(pending_positions, starting_node)

    pending_pos_map := map[string]bool{}

    valid_dirs := map[string][]string{ "N": {"E", "W"}, "S": {"E", "W"}, "E": {"N", "S"}, "W": {"N", "S"} }

    loop_num := 0
    if len(*visualise) == 0 {
	for len(pending_positions) > 0 {
        loop_num += 1

		current_position := pending_positions[0]
		pending_positions = pending_positions[1:]

        //current_position_key := fmt.Sprintf("%v", current_position)
        current_position_key := fmt.Sprintf("%s,%d,%d,%d", current_position.direction, current_position.y, current_position.x, current_position.cost)
        //fmt.Printf("Pos: %s | p %d\n", current_position_key, len(pending_positions))

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
                //fmt.Printf("Next: %v | p %d\n", next_position, len(pending_positions))

				isPosRangeBad := next_position.y >= len(graph_nodes)
				isPosRangeBad = isPosRangeBad || next_position.x >= len(graph_nodes[0])
				isPosRangeBad = isPosRangeBad || next_position.y < 0
				isPosRangeBad = isPosRangeBad || next_position.x < 0
				isPosRangeBad = isPosRangeBad || next_position.direction == "0"
				if isPosRangeBad {
					// Out of range, it expires.
					//fmt.Printf("[%d] Expire (invalid): %v\n", i, next_position)
                    break
				}

                next_node := &graph_nodes[next_position.y][next_position.x]
                //fmt.Printf("[%d] Cost: %d + %d = %d\n", i, next_cost, next_node.cost, next_cost + next_node.cost)
                next_cost += next_node.cost
				// If minimum distance before turning has been reached
				if (i+1) >= min_before_turn {
                    //fmt.Printf("Reached min\n")
					next_position.cost = next_cost

					if next_node.visited {
                        //fmt.Printf("[%d] Expire (visited): %v\n", i, next_position)
						continue
					}

					if next_cost < next_node.lowest_cost {
						next_node.lowest_cost = next_cost
						if next_node.y == ending_node.y && next_node.x == ending_node.x {
							next_node.lowest_cost_path = next_position.history
                            fmt.Printf("[%d] Candidate: %d | %s", i, next_cost, next_position.history)
						}
					} else {
                        //fmt.Printf("[%d] Not cheaper: %d < %d %s %v\n", i, next_cost, next_node.lowest_cost, next_position.history, next_node)
					}

                    //next_node_key := fmt.Sprintf("%v", *&next_position)
                    //next_node_key := fmt.Sprintf("%s,%d,%d,%d", next_position.direction, next_position.y, next_position.x, next_position.cost)
                    next_node_key := fmt.Sprintf("%s,%d,%d", next_position.direction, next_position.y, next_position.x)
                    //next_node_pending, next_node_pending_found := pending_pos_map[next_node_key]
                    next_node_pending, next_node_pending_found := pending_pos_map[next_node_key]
                    if !next_node_pending_found {
                        next_node_pending = false
                    }
					if next_node.lowest_cost < math.MaxInt64 && !next_node_pending {
                        //fmt.Printf("[%d] Appending next: %v\n", i, next_position)
						pending_positions = append(pending_positions, next_position)
                        pending_pos_map[next_node_key] = true
					} else {
                        //fmt.Printf("Not cheaper: %s %t %v\n", next_node_key, next_node_pending, *next_node)
                    }
				}
			}
		}
		graph_nodes[current_position.y][current_position.x].visited = true

        if current_node.y == ending_node.y && current_node.x == ending_node.x {
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
        //if len(pending_positions) < 10 {
        //    fmt.Fprintf(os.Stderr, "Pending: %v\n", pending_positions)
        //} else {
        //    fmt.Fprintf(os.Stderr, "Pending: %d %v\n", len(pending_positions), pending_positions[:2])
        //}
        if loop_num % 10_000 == 0 {
            if len(pending_positions) < 10 {
                fmt.Printf("Pending: %v\n", pending_positions)
            } else {
                fmt.Printf("Pending: %d %v\n", len(pending_positions), pending_positions[:2])
            }
        }
	}

    for i_y,graph_line := range graph_nodes {
        for i_x,graph_nde := range graph_line {
            if graph_nde.lowest_cost < 10_000 {
                fmt.Printf("%d:%4d | ", graph_nde.cost, graph_nde.lowest_cost)
            } else {
                fmt.Printf("%d:     | ", graph_nde.cost)
            }
            graph_nodes[i_y][i_x].visited = false
        }
        fmt.Printf("\n")
    }

    total = graph_nodes[ending_node.y][ending_node.x].lowest_cost
    *visualise = graph_nodes[ending_node.y][ending_node.x].lowest_cost_path
    first_total = total

	fmt.Printf("T: p1 %d p2 %d \n", first_total, total)
	}

    if len(*visualise) > 0 {
        vis_node := starting_node
        vis_total := 0
        for _,dir := range *visualise {
            vis_node = moveDirection(string(dir), vis_node)
            if vis_node.x >= 0 && vis_node.y >= 0 && vis_node.y < len(graph_nodes) && vis_node.x < len(graph_nodes[0]) {
                graph_nodes[vis_node.y][vis_node.x].visited = true
                add_cost := graph_nodes[vis_node.y][vis_node.x].cost
                fmt.Printf("CT: %4d + %4d = %4d\n", vis_total, add_cost,vis_total + add_cost)
                vis_total += add_cost
            } else {
                fmt.Printf("ERROR: %v\n", vis_node)
            }
        }
    }

    for _,graph_line := range graph_nodes {
        for _,graph_nde := range graph_line {
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
