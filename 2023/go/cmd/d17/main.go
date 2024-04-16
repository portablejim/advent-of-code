package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

type GraphNode struct {
    cost int
    visited bool
    lowest_cost int
    lowest_cost_path string
}

type FloorTile struct {
    char rune
    visited_count int
}

type CruciblePosition struct {
    direction rune
    history string
    y int
    x int
    is_starting bool
}

func moveDirection(direction rune, pos CruciblePosition) CruciblePosition {
    next_history := pos.history + string(direction)
    if direction == rune('N') {
        return CruciblePosition{'N', next_history, pos.y - 1, pos.x, false}
    } else if direction == rune('E') {
        return CruciblePosition{'E', next_history, pos.y, pos.x + 1, false}
    } else if direction == rune('S') {
        return CruciblePosition{'S', next_history, pos.y + 1, pos.x, false}
    } else if direction == rune('W') {
        return CruciblePosition{'W', next_history, pos.y, pos.x - 1, false}
    } else {
        return CruciblePosition{'0', next_history, pos.y, pos.x, false}
    } 
}

func getNextDirections(direction rune, previous_directions string, min_before_turn int, max_before_turn int) []rune {
    output := []rune{}

    // Not enough directions exist yet before turning.
    if len(previous_directions) < min_before_turn {
        fmt.Printf("Forced straight 1: %d < %d\n", len(previous_directions), min_before_turn)
        output = append(output, direction)
        return output
    }

    // Count num taken in each direction.
    count_same_direction := 0
    for i := len(previous_directions)-1; i >= 0 && count_same_direction <= (max_before_turn+1); i -= 1 {
        if previous_directions[i] != byte(direction) {
            break
        }
        count_same_direction += 1
    }

    // Not enough of the same direction to turn.
    if count_same_direction < min_before_turn {
        output = append(output, direction)
        fmt.Printf("Forced straight 2: %d < %d\n", count_same_direction, min_before_turn)
        return output
    }

    can_move_same := true
    if count_same_direction >= max_before_turn {
        fmt.Printf("Forced turn: %d > %d\n", count_same_direction, max_before_turn)
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

    fmt.Printf("Next dirs: %v\n", output)

    return output
}

func getLowestCost(s_i int, starting_nodes []CruciblePosition, ending_nodes_list []CruciblePosition, graph_nodes [][]GraphNode, min_before_turn int, max_before_turn int) (int, string) {
        ending_node := ending_nodes_list[s_i]

        pending_positions := []CruciblePosition{}
        pending_positions = append(pending_positions, starting_nodes...)

        // Handle the positions
        for len(pending_positions) > 0 {
            prev_position := pending_positions[0]
            pending_positions = pending_positions[1:]

            current_pos := moveDirection(prev_position.direction, prev_position)
            //fmt.Printf("Current pos: %v\n", current_pos)
            if current_pos.y >= len(graph_nodes) || current_pos.x >= len(graph_nodes[0]) || current_pos.y < 0 || current_pos.x < 0 || current_pos.direction == '0' {
                // Out of range, it expires.
                //fmt.Printf("expired r: %v\n", current_end)
                //fmt.Printf("Expire (invalid): %v\n", current_pos)
                continue
            }

            prev_node := graph_nodes[prev_position.y][prev_position.x]
            if prev_position.is_starting {
                // Starting positions get a free first node
                //prev_node.lowest_cost = 0
                //prev_node.visited = true
            }
            current_node := &graph_nodes[current_pos.y][current_pos.x]
            new_cost := prev_node.lowest_cost + current_node.cost
            fmt.Printf("N: %v %v %v %d\n", prev_node, current_pos, *current_node, new_cost)
            if new_cost < 0 || new_cost >= current_node.lowest_cost {
                // Higher cost, don't continue.
                //fmt.Printf("Expire (cost): %d > %d | %s\n", new_cost, current_node.lowest_cost, current_pos.history)
                continue
            }

            current_node.visited = true
            current_node.lowest_cost = new_cost
            current_node.lowest_cost_path = current_pos.history

            if current_pos.y == ending_node.y && current_pos.x == ending_node.x {
                // Found the end.
                //return current_node.lowest_cost, current_node.lowest_cost_path
                fmt.Printf("Possible lowest: %v\n", current_pos)
            }

            for _,next_dir := range getNextDirections(current_pos.direction, current_pos.history, min_before_turn, max_before_turn) {
                pending_positions = append(pending_positions, CruciblePosition{next_dir, current_pos.history, current_pos.y, current_pos.x, false})
            }
            //fmt.Printf("Next: %v\n", pending_positions)
            fmt.Printf("Next: %d\n", len(pending_positions))

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
    flag.Parse()
    dat, err := os.ReadFile(*filename)
    if err != nil {
        log.Fatalf("unable to read file: %f", err)
    }

    graph_nodes := [][]GraphNode{}

    // Parse the data.
    for _, f_line := range strings.Split(strings.Replace(string(dat), "\r\n", "\n", -1), "\n") {
        if len(f_line) == 0 {
            continue
        }
        f_line = strings.Trim(f_line, " \n")

        graph_node_row := []GraphNode{}
        for c_num := 0; c_num < len(f_line); c_num += 1 {
            node_weight := f_line[c_num] - '0'
            graph_node_row = append(graph_node_row, GraphNode{int(node_weight), false,  math.MaxInt64, "" })
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

    starting_nodes_list := [][]CruciblePosition{ {{'E', "", 0, 0, true}, {'S', "", 0, 0, true}} }
    ending_nodes_list := []CruciblePosition{ {'0', "", max_y, max_x, false} }

    total := math.MaxInt64
    first_total := math.MaxInt64

    min_before_turn := 0
    max_before_turn := 3
    if *part2 {
        min_before_turn = 4
        max_before_turn = 10
    }

    for s_i,starting_nodes := range starting_nodes_list {
        current_cost, current_cost_path := getLowestCost(s_i, starting_nodes, ending_nodes_list, graph_nodes, min_before_turn, max_before_turn)
        fmt.Printf("Cost: %d, %s\n", current_cost, current_cost_path)
        if current_cost < total {
            total = current_cost
        }
        if first_total == math.MaxInt64 {
            first_total = current_cost
        }
    }

    /*
    fmt.Printf("Tiles:\n")
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


    fmt.Printf("T: p1 %d p2 %d \n", first_total, total)
}

