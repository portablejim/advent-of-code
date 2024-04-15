package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type FloorTile struct {
    char rune
    visited_count int
}

type LaserPosition struct {
    direction rune
    y int
    x int
}

func moveDirection(direction rune, pos LaserPosition) LaserPosition {
    if direction == rune('N') {
        return LaserPosition{'N', pos.y - 1, pos.x}
    } else if direction == rune('E') {
        return LaserPosition{'E', pos.y, pos.x + 1}
    } else if direction == rune('S') {
        return LaserPosition{'S', pos.y + 1, pos.x}
    } else if direction == rune('W') {
        return LaserPosition{'W', pos.y, pos.x - 1}
    } else {
        return LaserPosition{'0', pos.y, pos.x}
    } 
}


func main() {
    var filename = flag.String("f", "../inputs/d16.sample.txt", "file to use")
    flag.Parse()
    dat, err := os.ReadFile(*filename)
    if err != nil {
        log.Fatalf("unable to read file: %f", err)
    }

    floor_tiles := [][]FloorTile{}

    // Parse the data.
    for _, f_line := range strings.Split(strings.Replace(string(dat), "\r\n", "\n", -1), "\n") {
        if len(f_line) == 0 {
            continue
        }
        f_line = strings.Trim(f_line, " \n")

        floor_tile_row := []FloorTile{}
        for c_num := 0; c_num < len(f_line); c_num += 1 {
            floor_tile_row = append(floor_tile_row, FloorTile{ rune(f_line[c_num]), 0 })
        }
        floor_tiles = append(floor_tiles, floor_tile_row)
    }

    if len(floor_tiles) == 0 || len(floor_tiles[0]) == 0 {
        fmt.Fprintf(os.Stderr, "Error when loading floor tiles\n")
        return
    }
    //fmt.Printf("floor_tile_rows: %v\n", floor_tiles)

    starting_laser_ends := []LaserPosition{}
    floor_tile_width := len(floor_tiles[0])
    floor_tile_height := len(floor_tiles[0])
    for y_i := 0; y_i < floor_tile_height; y_i += 1 {
        starting_laser_ends = append(starting_laser_ends, LaserPosition{'E', y_i, -1})
        starting_laser_ends = append(starting_laser_ends, LaserPosition{'W', y_i, floor_tile_width})
    }
    for x_i := 0; x_i < floor_tile_width; x_i += 1 {
        starting_laser_ends = append(starting_laser_ends, LaserPosition{'S', -1, x_i})
        starting_laser_ends = append(starting_laser_ends, LaserPosition{'N', floor_tile_height, x_i})
    }

    total := 0
    first_total := -1

    for s_i,starting_laser_end := range starting_laser_ends {
        visited_locations := map[string]bool{}
        local_total := 0

        laser_ends := []LaserPosition{}
        laser_ends = append(laser_ends, starting_laser_end)

        // Handle the positions
        for len(laser_ends) > 0 {
            prev_end := laser_ends[0]
            laser_ends = laser_ends[1:]

            current_end := moveDirection(prev_end.direction, prev_end)
            if current_end.y >= len(floor_tiles) || current_end.x >= len(floor_tiles[0]) || current_end.y < 0 || current_end.x < 0 || current_end.direction == '0' {
                // Out of range, it expires.
                //fmt.Printf("expired r: %v\n", current_end)
                continue
            }
            location_key := fmt.Sprintf("%v", current_end)
            _,is_visited := visited_locations[location_key]
            if is_visited {
                // Direction & position match, expire.
                //fmt.Printf("expired v: %v\n", current_end)
                continue
            }
            visited_locations[location_key] = true
            current_tile := &floor_tiles[current_end.y][current_end.x]
            //fmt.Printf("current: %v %v %c\n", prev_end, current_end, current_tile.char)
            floor_tiles[current_end.y][current_end.x].visited_count += 1
            next_end := LaserPosition{current_end.direction, current_end.y, current_end.x}
            if current_tile.char == '.' || (current_tile.char == '-' && (current_end.direction == 'E' || current_end.direction == 'W')) || (current_tile.char == '|' && (current_end.direction == 'N' || current_end.direction == 'S')) {
                // Empty tile (or non-action tile)
                laser_ends = append(laser_ends, next_end)
            } else if current_tile.char == '|' && (current_end.direction == 'E' || current_end.direction == 'W') {
                laser_ends = append(laser_ends, LaserPosition{'N', current_end.y, current_end.x})
                laser_ends = append(laser_ends, LaserPosition{'S', current_end.y, current_end.x})
            } else if current_tile.char == '-' && (current_end.direction == 'N' || current_end.direction == 'S') {
                laser_ends = append(laser_ends, LaserPosition{'E', current_end.y, current_end.x})
                laser_ends = append(laser_ends, LaserPosition{'W', current_end.y, current_end.x})
            } else if current_tile.char == '/' && current_end.direction == 'N' {
                laser_ends = append(laser_ends, LaserPosition{'E', current_end.y, current_end.x})
            } else if current_tile.char == '/' && current_end.direction == 'E' {
                laser_ends = append(laser_ends, LaserPosition{'N', current_end.y, current_end.x})
            } else if current_tile.char == '/' && current_end.direction == 'W' {
                laser_ends = append(laser_ends, LaserPosition{'S', current_end.y, current_end.x})
            } else if current_tile.char == '/' && current_end.direction == 'S' {
                laser_ends = append(laser_ends, LaserPosition{'W', current_end.y, current_end.x})
            } else if current_tile.char == '\\' && current_end.direction == 'N' {
                laser_ends = append(laser_ends, LaserPosition{'W', current_end.y, current_end.x})
            } else if current_tile.char == '\\' && current_end.direction == 'W' {
                laser_ends = append(laser_ends, LaserPosition{'N', current_end.y, current_end.x})
            } else if current_tile.char == '\\' && current_end.direction == 'S' {
                laser_ends = append(laser_ends, LaserPosition{'E', current_end.y, current_end.x})
            } else if current_tile.char == '\\' && current_end.direction == 'E' {
                laser_ends = append(laser_ends, LaserPosition{'S', current_end.y, current_end.x})
            }
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
        for ft_i,floor_tile_line := range floor_tiles {
            for ft_j,floor_tile := range floor_tile_line {
                if floor_tile.visited_count > 0 {
                    local_total += 1
                    floor_tiles[ft_i][ft_j].visited_count = 0
                }
            }
        }
        if local_total >= total {
            fmt.Printf("Max local_total: [%d] %d %d %v\n", s_i, total, local_total, starting_laser_end)
            total = local_total
        }
        if first_total < 0 {
            first_total = local_total
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

