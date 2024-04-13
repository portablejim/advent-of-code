package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type Coord struct {
    y int
    x int
}

type Star struct {
    position Coord
    distance_list []int
}

type SkyNode struct {
    symbol string
    weight int
}

func doCoordsMatch(a Coord, b Coord) bool {
    return a.x == b.x && a.y == b.y
}

func main() {
    var filename = flag.String("f", "../inputs/d11.sample.txt", "file to use")
    var empty_weight = flag.Int("w", 2, "weight of empty rows")
    flag.Parse()
    dat, err := os.ReadFile(*filename)
    if err != nil {
        log.Fatalf("unable to read file: %f", err)
    }

    p_star := regexp.MustCompile("#")

    //sky_image_expanded := []string{}
    sky_image_weights := [][]SkyNode{}

    // Parse the data.
    for i_y, f_line := range strings.Split(string(dat), "\n") {
        if len(f_line) == 0 {
            continue
        }
        inital_weight := 1
        sky_image_line := []SkyNode{}
        if !p_star.MatchString(f_line) {
            fmt.Printf("Adding weight at row %d\n", i_y)
            inital_weight = *empty_weight
        }
        for _,f_char := range f_line {
            sky_image_line = append(sky_image_line, SkyNode{ string(f_char), inital_weight })
        }
        sky_image_weights = append(sky_image_weights, sky_image_line)
    }

    // Find empty vertical lines and duplicate.
    for i_x := range len(sky_image_weights[0]) {
        is_empty := true
        for i_y := range len(sky_image_weights) {
            if sky_image_weights[i_y][i_x].symbol == "#" {
                is_empty = false
            }
        }
        if is_empty {
            fmt.Printf("Adding weight at col %d\n", i_x)
            for i_y := range len(sky_image_weights) {
                sky_image_weights[i_y][i_x].weight *= *empty_weight
            }
        }
    }
    // Print image.
    fmt.Printf("Image:\n")
    for _,i_line := range sky_image_weights {
        for _,i_char := range i_line {
            if i_char.weight > 1 {
                if i_char.symbol == "#" {
                    // This should not happen.
                    fmt.Printf("!")
                } else {
                    fmt.Printf("•")
                }
            } else {
                fmt.Printf("%s", i_char.symbol)
            }
        }
        fmt.Printf("\n")
    }

    total := 0
    star_list := []Star{}

    for i_y,img_line := range sky_image_weights {
        for i_x,img_weight := range img_line {
            if img_weight.symbol == "#" {
                current_star := Star{ Coord{i_y, i_x}, []int{} }
                star_list = append(star_list, current_star)
            }
        }
    }

    for current_star_i := range len(star_list) {
        current_star := &star_list[current_star_i]
        for candidate_star_i := range len(star_list){
            if candidate_star_i < current_star_i {
                current_star.distance_list = append(current_star.distance_list, -1)
            } else if candidate_star_i == current_star_i {
                current_star.distance_list = append(current_star.distance_list, 0)
            } else {
                candidate_star := star_list[candidate_star_i]

                distance_x := 0
                x_i_start := current_star.position.x
                x_i_end := candidate_star.position.x
                if x_i_start > x_i_end {
                    x_i_start, x_i_end = x_i_end, x_i_start
                }
                x_i := x_i_start
                for x_i < x_i_end {
                    distance_x += sky_image_weights[candidate_star.position.y][x_i].weight
                    x_i += 1
                }

                distance_y := 0
                y_i_start := current_star.position.y
                y_i_end := candidate_star.position.y
                if y_i_start > y_i_end {
                    y_i_start, y_i_end = y_i_end, y_i_start
                }
                y_i := y_i_start
                for y_i < y_i_end {
                    distance_y += sky_image_weights[y_i][candidate_star.position.x].weight
                    y_i += 1
                }

                total_distance := distance_x + distance_y
                current_star.distance_list = append(current_star.distance_list, total_distance)
                total += total_distance
            }

        }
    }

    /*
    for i,c_star := range star_list {
        fmt.Printf("S %d: %v\n", i, c_star)
    }
    */


    //fmt.Printf("Start: %v\n", start_pos)
    //fmt.Printf("numbers list 1: %v\n", numbers_list)
    //fmt.Printf("numbers map: %v\n", numbers_map)
    fmt.Printf("T: %d\n", total)
}

