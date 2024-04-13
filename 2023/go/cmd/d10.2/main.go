package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Coord struct {
    y int
    x int
}

type PipeNode struct {
    symbol string
    pos Coord
    left Coord
    right Coord
    status int
    distance int
}

func printGraph(graph [][]PipeNode, inside_points []Coord) {
    for _,n_y := range graph {
        for _,n_x := range n_y {
            if n_x.distance >= 0 {
                fmt.Printf("%s", n_x.symbol)
            } else {
                is_inside := false
                for _,inp := range inside_points {
                    if doCoordsMatch(n_x.pos, inp) {
                        is_inside = true;
                        break
                    }

                }
                if is_inside {
                    fmt.Printf("I")
                } else {
                    fmt.Printf(" ")
                }
            }
        }
        fmt.Printf("\n")
    }
}

func getOtherValidDir(symbol string, direction string) string {
    if symbol == "|" {
        if direction == "N" {
            return "S"
        } else {
            return "N"
        }
    }
    if symbol == "-" {
        if direction == "E" {
            return "W"
        } else {
            return "E"
        }
    }
    if symbol == "L" {
        if direction == "N" {
            return "E"
        } else {
            return "N"
        }
    }
    if symbol == "J" {
        if direction == "N" {
            return "W"
        } else {
            return "N"
        }
    }
    if symbol == "7" {
        if direction == "S" {
            return "W"
        } else {
            return "S"
        }
    }
    if symbol == "F" {
        if direction == "S" {
            return "E"
        } else {
            return "S"
        }
    }

    return ""
}

func inverseDir(in_pos string) string {
    if in_pos == "N" {
        return "S"
    }
    if in_pos == "S" {
        return "N"
    }
    if in_pos == "W" {
        return "E"
    }
    if in_pos == "E" {
        return "W"
    }

    return ""
}

func move(dir string, from_pos Coord, len_y int, len_x int) Coord {
    output := Coord{-1, -1}
    if dir == "N" {
        output = Coord{from_pos.y - 1, from_pos.x}
    } else if dir == "S" {
        output = Coord{from_pos.y + 1, from_pos.x}
    }
    if dir == "W" {
        output = Coord{from_pos.y, from_pos.x - 1}
    } else if dir == "E" {
        output = Coord{from_pos.y, from_pos.x + 1}
    }

    if output.x < 0 || output.y < 0 || output.y >= len_y || output.x >= len_x {
        return Coord{-1, -1}
    }

    return output;
}

func doCoordsMatch(a Coord, b Coord) bool {
    return a.x == b.x && a.y == b.y
}

func connectsTo(graph [][]PipeNode, from_pos Coord, to_pos Coord) bool {
    if from_pos.y < 0 || from_pos.y >= len(graph) || from_pos.x < 0 || from_pos.x >= len(graph[from_pos.y]) {
        return false
    }
    if to_pos.y < 0 || to_pos.y >= len(graph) || to_pos.x < 0 || to_pos.x >= len(graph[to_pos.y]) {
        return false
    }
    len_y := len(graph)
    len_x := len(graph[0])

    from_node := graph[from_pos.y][from_pos.x]
    //to_node := graph[to_pos.y][to_pos.x]

    if from_node.symbol == "|" && doCoordsMatch(to_pos, move("N", from_pos, len_y, len_x)) {
        return true
    }
    if from_node.symbol == "|" && doCoordsMatch(to_pos, move("S", from_pos, len_y, len_x)) {
        return true
    }
    if from_node.symbol == "-" && doCoordsMatch(to_pos, move("W", from_pos, len_y, len_x)) {
        return true
    }
    if from_node.symbol == "-" && doCoordsMatch(to_pos, move("E", from_pos, len_y, len_x)) {
        return true
    }
    if from_node.symbol == "L" && doCoordsMatch(to_pos, move("N", from_pos, len_y, len_x)) {
        return true
    }
    if from_node.symbol == "L" && doCoordsMatch(to_pos, move("E", from_pos, len_y, len_x)) {
        return true
    }
    if from_node.symbol == "J" && doCoordsMatch(to_pos, move("N", from_pos, len_y, len_x)) {
        return true
    }
    if from_node.symbol == "J" && doCoordsMatch(to_pos, move("W", from_pos, len_y, len_x)) {
        return true
    }
    if from_node.symbol == "7" && doCoordsMatch(to_pos, move("S", from_pos, len_y, len_x)) {
        return true
    }
    if from_node.symbol == "7" && doCoordsMatch(to_pos, move("W", from_pos, len_y, len_x)) {
        return true
    }
    if from_node.symbol == "F" && doCoordsMatch(to_pos, move("S", from_pos, len_y, len_x)) {
        return true
    }
    if from_node.symbol == "F" && doCoordsMatch(to_pos, move("E", from_pos, len_y, len_x)) {
        return true
    }

    return false
}

func splitSpaceSepNums(num_list_str string) []int64 {
    num_list_trimmed := strings.Trim(num_list_str, " ")
    num_str_list := strings.Split(num_list_trimmed, " ")
    output := []int64{}

    for _,num_str := range num_str_list {
        num_int, err := strconv.ParseInt(strings.Trim(num_str, " "), 10, 64)
        if err == nil {
            output = append(output, num_int)
        } else {
            fmt.Fprintf(os.Stderr, "Error parsing string as number: '%s' (%s)\n", num_str, num_list_str)
        }
    }

    return output
}


func main() {
    var filename = flag.String("f", "../inputs/d10.sample1.txt", "file to use")
    flag.Parse()
    dat, err := os.ReadFile(*filename)
    if err != nil {
        log.Fatalf("unable to read file: %f", err)
    }

    pipe_graph := [][]PipeNode{}
    start_pos := Coord{ -1, -1 }


    // Parse the data.
    for f_y, f_line := range strings.Split(string(dat), "\n") {
        g_line := []PipeNode{}
        for f_x, f_char := range strings.Split(f_line, "") {
            g_line = append(g_line, PipeNode{f_char, Coord{f_y, f_x}, Coord{-1, -1}, Coord{-1, -1}, 0, -1})
            if f_char == "S" {
                start_pos = Coord{f_y, f_x}
            }
        }
        if len(g_line) > 0 {
            pipe_graph = append(pipe_graph, g_line)
        }
    }

    if start_pos.x < 0 || start_pos.y < 0 {
        fmt.Printf("Unable to find start position\n")
        return
    }

    pipe_graph[start_pos.y][start_pos.x].distance = 0

    start_pipes := []Coord{}
    start_pipes_from := []string{}
    dist := 0

    total := 0

    dir_list := []string{ "N", "S", "E", "W" }
    for _,dir := range dir_list {
        candidate_pos := move(dir, start_pos, len(pipe_graph), len(pipe_graph[0]))
        if connectsTo(pipe_graph, candidate_pos, start_pos) {
            start_pipes = append(start_pipes, candidate_pos)
            start_pipes_from = append(start_pipes_from, dir)
        }
    }

    if len(start_pipes) != 2 {
        fmt.Printf("Invalid starting condition: %d connections", len(start_pipes))
    }

    if (start_pipes_from[0] == "N" && start_pipes_from[1] == "S") || (start_pipes_from[0] == "S" && start_pipes_from[1] == "N") {
        pipe_graph[start_pos.y][start_pos.x].symbol = "|"
    }
    if (start_pipes_from[0] == "E" && start_pipes_from[1] == "W") || (start_pipes_from[0] == "W" && start_pipes_from[1] == "E") {
        pipe_graph[start_pos.y][start_pos.x].symbol = "-"
    }
    if (start_pipes_from[0] == "N" && start_pipes_from[1] == "E") || (start_pipes_from[0] == "E" && start_pipes_from[1] == "N") {
        pipe_graph[start_pos.y][start_pos.x].symbol = "L"
    }
    if (start_pipes_from[0] == "N" && start_pipes_from[1] == "W") || (start_pipes_from[0] == "W" && start_pipes_from[1] == "N") {
        pipe_graph[start_pos.y][start_pos.x].symbol = "J"
    }
    if (start_pipes_from[0] == "W" && start_pipes_from[1] == "S") || (start_pipes_from[0] == "S" && start_pipes_from[1] == "W") {
        pipe_graph[start_pos.y][start_pos.x].symbol = "7"
    }
    if (start_pipes_from[0] == "E" && start_pipes_from[1] == "S") || (start_pipes_from[0] == "S" && start_pipes_from[1] == "E") {
        pipe_graph[start_pos.y][start_pos.x].symbol = "F"
    }

    prev_pos_left := start_pos
    current_left_pos := start_pipes[0]
    current_left_from := start_pipes_from[0]
    prev_pos_right := start_pos
    current_right_pos := start_pipes[1]
    current_right_from := start_pipes_from[1]

    for range 100_000 {
        dist += 1

        if current_left_pos.y < 0 || current_right_pos.y < 0 || current_left_pos.x < 0 || current_right_pos.x < 0 {
            fmt.Printf("Invalid positions: %v %v | %d\n", current_left_pos, current_right_pos, dist)
            return
        }
        //fmt.Printf("CLP: %v | CRP: %v\n", current_left_pos, current_right_pos)
        //fmt.Printf("CFL: %s \n", current_left_from)
        //fmt.Printf("CFR: %s \n", current_right_from)

        current_left := &pipe_graph[current_left_pos.y][current_left_pos.x]
        current_right := &pipe_graph[current_right_pos.y][current_right_pos.x]
        current_left.status = 1
        current_left.right = prev_pos_left
        if current_left.distance < 0 {
            current_left.distance = dist
        }
        are_either_nodes_filled := (current_left.left.x > -1 && current_left.right.x > -1) || (current_right.left.x > -1 && current_right.right.x > -1)
        if are_either_nodes_filled {
            fmt.Printf("Either filled at %d\n", dist)
        }
        current_right.status = 1
        current_right.left = prev_pos_right
        if current_right.distance < 0 {
            current_right.distance = dist
        }

        //fmt.Printf("CL: %v | CR: %v\n", *current_left, *current_right)

        are_both_nodes_filled := (current_left.left.x > -1 && current_left.right.x > -1) && (current_right.left.x > -1 && current_right.right.x > -1)
        if doCoordsMatch(current_left_pos, current_right_pos) || are_both_nodes_filled {
            total = dist
            break
        }

        next_dir_left := getOtherValidDir(current_left.symbol, inverseDir(current_left_from))
        //fmt.Printf("NDL: %s => %s\n", current_left_from, next_dir_left)
        next_pos_left := move(next_dir_left, current_left_pos, len(pipe_graph), len(pipe_graph[0]))

        prev_pos_left = current_left_pos
        current_left_pos = next_pos_left
        current_left_from = next_dir_left

        next_dir_right := getOtherValidDir(current_right.symbol, inverseDir(current_right_from))
        //fmt.Printf("NDR: %s => %s\n", current_right_from, next_dir_right)
        next_pos_right := move(next_dir_right, current_right_pos, len(pipe_graph), len(pipe_graph[0]))

        prev_pos_right = current_right_pos
        current_right_pos = next_pos_right
        current_right_from = next_dir_right
    }

    // Part 2: Count inside.
    inside_count := 0
    inside_y := []Coord{}
    for _,n_y := range pipe_graph {
        is_inside := false
        wall_start := ""
        is_wall := false
        for _,n_x := range n_y {
            if n_x.distance >= 0 {
                if is_wall {
                    if (wall_start == "L" && n_x.symbol == "J") || (wall_start == "F" && n_x.symbol == "7") {
                        is_wall = false
                        wall_start = ""
                    }
                    if (wall_start == "L" && n_x.symbol == "7") || (wall_start == "F" && n_x.symbol == "J") {
                        is_wall = false
                        wall_start = ""
                        is_inside = !is_inside
                    }
                } else if n_x.symbol == "|" {
                    is_inside = !is_inside
                    is_wall = false
                    wall_start = ""
                } else {
                    is_wall = true
                    wall_start = n_x.symbol
                }
            } else {
                if is_inside {
                    inside_y = append(inside_y, n_x.pos)
                }
            }
        }
    }

    inside_points := inside_y
    inside_count = len(inside_points)
    fmt.Printf("\n")
    printGraph(pipe_graph, inside_y)


    fmt.Printf("Start: %v\n", start_pos)
    fmt.Printf("T: %d\n", total)
    fmt.Printf("Inside: %d\n", inside_count)
}

