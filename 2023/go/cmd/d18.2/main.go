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
    y int64
    x int64
}

type HolePoint struct {
    pos Coord
    length int
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

    return output;
}

func doCoordsMatch(a Coord, b Coord) bool {
    return a.x == b.x && a.y == b.y
}

func decodeHex(input_hex string) (string, int64, bool) {
    parsed_hex, err := strconv.ParseInt(input_hex, 16, 64)
    if err != nil {
        fmt.Printf("Error parsing hex %d: %v\n", parsed_hex, err)
        return "", 0, false
    }

    output_direction_num := parsed_hex & 0xf
    output_count := (parsed_hex >> 4)

    var direction_lookup = []string{"R", "D", "L", "U"}
    output_direction := direction_lookup[output_direction_num]

    return output_direction, output_count, true
}

func main() {
    var filename = flag.String("f", "../inputs/d10.sample1.txt", "file to use")
    flag.Parse()
    dat, err := os.ReadFile(*filename)
    if err != nil {
        log.Fatalf("unable to read file: %f", err)
    }

    start_pos := Coord{ 0, 0 }

    min_y := 0
    min_x := 0
    max_y := 0
    max_x := 0
    dug_blocks_list := []HolePoint{}
    dug_blocks_list = append(dug_blocks_list, HolePoint{start_pos, 0})
    polygon_points := []Coord{}
    polygon_points = append(polygon_points, Coord{0, 0})

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

        color_hex := line_parts[2][2:len(line_parts[2])-1]
        line_direction := line_parts[0]

        part2 := true
        if part2 {
            var decode_worked bool
            line_direction, block_count, decode_worked = decodeHex(color_hex)
            if !decode_worked {
                continue
            }
            fmt.Printf("I: %s = %s %d\n", color_hex, line_direction, block_count)
        }

        current_pos = move(line_direction, current_pos, block_count)

        polygon_points = append(polygon_points, current_pos)
        dug_blocks_list = append(dug_blocks_list, HolePoint{current_pos, int(block_count)})
    }
    polygon_points = append(polygon_points, Coord{0, 0})
    // One will be 0
    temp_length := dug_blocks_list[len(dug_blocks_list)-1].pos.x
    temp_length += dug_blocks_list[len(dug_blocks_list)-1].pos.y
    if temp_length < 0 {
        temp_length *= -1
    }
    dug_blocks_list = append(dug_blocks_list, HolePoint{Coord{0, 0}, int(temp_length)})

    total := int64(0)

    temp_total := int64(0)
    for pt_i := 0; pt_i < (len(dug_blocks_list)-1); pt_i += 1 {
        fmt.Printf("%d, %d\n", dug_blocks_list[pt_i].pos.x, dug_blocks_list[pt_i].pos.y)
    }
    for pt_i := 0; pt_i < (len(dug_blocks_list)-1); pt_i += 1 {
        temp_total += (dug_blocks_list[pt_i].pos.x * dug_blocks_list[pt_i+1].pos.y) - (dug_blocks_list[pt_i].pos.y * dug_blocks_list[pt_i+1].pos.x)
        temp_total += int64(dug_blocks_list[pt_i].length)
        fmt.Printf("(%d * %d) - (%d * %d) +", dug_blocks_list[pt_i].pos.x, dug_blocks_list[pt_i+1].pos.y, dug_blocks_list[pt_i].pos.y, dug_blocks_list[pt_i+1].pos.x)
    }
    fmt.Printf("\n")
    total = temp_total / 2 + 1


    fmt.Printf("Range: %d - %d | %d - %d\n", min_y, max_y, min_x, max_x)

    // Process the data.


    //fmt.Printf("numbers list 1: %v\n", numbers_list)
    //fmt.Printf("numbers map: %v\n", numbers_map)
    fmt.Printf("T: %d\n", total)
}

