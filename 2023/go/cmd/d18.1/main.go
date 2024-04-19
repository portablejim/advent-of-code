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
    y int32
    x int32
}

type HolePoint struct {
    pos Coord
    dug int
    dig_passes_bottom bool
    colorHex string
}

func move(dir string, from_pos Coord) Coord {
    output := Coord{-1, -1}
    if dir == "U" {
        output = Coord{from_pos.y - 1, from_pos.x}
    } else if dir == "D" {
        output = Coord{from_pos.y + 1, from_pos.x}
    } else if dir == "L" {
        output = Coord{from_pos.y, from_pos.x - 1}
    } else if dir == "R" {
        output = Coord{from_pos.y, from_pos.x + 1}
    }

    return output;
}

func doCoordsMatch(a Coord, b Coord) bool {
    return a.x == b.x && a.y == b.y
}

func hexToRGB(input_hex string, fallback_r int, fallback_g int, fallback_b int) (int, int, int) {
    parsed_hex, err := strconv.ParseInt(input_hex, 16, 64)
    if err != nil {
        fmt.Printf("Error parsing hex %d: %v\n", parsed_hex, err)
        return fallback_r, fallback_g, fallback_b
    }

    output_b := parsed_hex & 0xff
    output_g := (parsed_hex >> 8) & 0xff
    output_r := (parsed_hex >> 16) & 0xff

    return int(output_r), int(output_g), int(output_b)
}


func main() {
    var filename = flag.String("f", "../inputs/d10.sample1.txt", "file to use")
    var is_verbose = flag.Bool("v", false, "verbose")
    flag.Parse()
    dat, err := os.ReadFile(*filename)
    if err != nil {
        log.Fatalf("unable to read file: %f", err)
    }

    start_pos := Coord{ 0, 0 }

    default_color_hex := "FFFFFF"

    min_y := 0
    min_x := 0
    max_y := 0
    max_x := 0
    dug_blocks_list := []HolePoint{}
    dug_blocks_list = append(dug_blocks_list, HolePoint{start_pos, 1, false, default_color_hex})

    //dug_blocks_map := map[int64]HolePoint{}

    current_pos := start_pos

    // Parse the data.
    for _, f_line := range strings.Split(string(dat), "\n") {
        line_parts := strings.Split(f_line, " ")

        if len(line_parts) != 3 || len(line_parts[2]) < 4 {
            continue
        }

        block_count, err := strconv.ParseInt(strings.Trim(line_parts[1], " "), 10, 32)
        if  err != nil {
            continue
        }

        color_hex := line_parts[2][2:len(line_parts[2])-2]

        for i := 0; i < int(block_count); i += 1 {
            is_first_move := i == 0
            is_last_move := i == int(block_count - 1)

            if line_parts[0] == "D" && is_first_move {
                dug_blocks_list[len(dug_blocks_list)-1].dig_passes_bottom = true
            }

            current_pos = move(line_parts[0], current_pos)

            does_pass_bottom := false
            if line_parts[0] == "D" && !is_last_move {
                does_pass_bottom = true
            }
            if line_parts[0] == "U" {
                does_pass_bottom = true
            }


            dug_blocks_list = append(dug_blocks_list, HolePoint{current_pos, 1, does_pass_bottom, color_hex})

            if min_y > int(current_pos.y) {
                min_y = int(current_pos.y)
            }
            if min_x > int(current_pos.x) {
                min_x = int(current_pos.x)
            }
            if max_y < int(current_pos.y) {
                max_y = int(current_pos.y)
            }
            if max_x < int(current_pos.x) {
                max_x = int(current_pos.x)
            }
        }
    }

    fmt.Printf("Range: %d - %d | %d - %d\n", min_y, max_y, min_x, max_x)

    dug_width := (max_x - min_x) + 1
    dug_height := (max_y - min_y) + 1

    // Process the data.
    hole_array := [][]HolePoint{}
    for hole_array_i := range(dug_height) {
        hole_array_line := []HolePoint{}
        for hole_array_j := range(dug_width) {
            hole_array_line = append(hole_array_line, HolePoint{Coord{int32(hole_array_i), int32(hole_array_j)}, 0, false, default_color_hex})
        }
        hole_array = append(hole_array, hole_array_line)
    }

    for _,dug_block := range dug_blocks_list {
        //hole_index_y := dug_block.pos.y-int32(min_y)
        //hole_index_x := dug_block.pos.x-int32(min_x)
        //fmt.Printf("Get: %d => %d | %d => %d | %d %d\n", dug_block.pos.y, min_y, dug_block.pos.x, min_x, hole_index_y, hole_index_x)
        hole_array[dug_block.pos.y-int32(min_y)][dug_block.pos.x-int32(min_x)] = dug_block
    }

    total := 0

    for hole_i,hole_line := range hole_array {
        is_inside := false
        for hole_j := 0; hole_j < len(hole_line); hole_j += 1 {
            hole_char := &hole_array[hole_i][hole_j]
            if hole_char.dug > 0 {
                // Wall
                hole_char.dug += 1
                if hole_char.dig_passes_bottom {
                    is_inside = !is_inside
                }
                total += 1
            } else if is_inside {
                hole_char.dug += 1
                total += 1
            }
        }
    }

    // Print the data.
    if *is_verbose {
        for _,hole_line := range hole_array {
            for _,hole_char := range hole_line {
                or, og, ob := hexToRGB(hole_char.colorHex, 255, 255, 255)
                if hole_char.dug >= 2 {
                    if hole_char.dig_passes_bottom {
                    fmt.Printf("\x1b[38;2;%d;%d;%dmB\x1b[0m", or, og, ob)
                    } else {
                    fmt.Printf("\x1b[38;2;%d;%d;%dmX\x1b[0m", or, og, ob)
                    }
                } else if hole_char.dug >= 1 {
                    fmt.Printf(".")
                } else {
                    fmt.Printf(" ")
                }
            }
            fmt.Printf("|\n")
        }
    }


    //fmt.Printf("numbers list 1: %v\n", numbers_list)
    //fmt.Printf("numbers map: %v\n", numbers_map)
    fmt.Printf("T: %d\n", total)
}

