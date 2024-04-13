package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
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

func doCoordsMatch(a Coord, b Coord) bool {
    return a.x == b.x && a.y == b.y
}

func main() {
    var filename = flag.String("f", "../inputs/d11.sample.txt", "file to use")
    flag.Parse()
    dat, err := os.ReadFile(*filename)
    if err != nil {
        log.Fatalf("unable to read file: %f", err)
    }

    p_star := regexp.MustCompile("#")

    sky_image_expanded := []string{}
    temp_expanded := []string{}

    // Parse the data.
    for f_y, f_line := range strings.Split(string(dat), "\n") {
        if len(f_line) == 0 {
            continue
        }
        sky_image_expanded = append(sky_image_expanded, f_line)
        temp_expanded = append(temp_expanded, "")
        // Expand horizontal lines
        if !p_star.MatchString(f_line) {
            fmt.Printf("Duplicating line %d | %s\n", f_y, f_line)
            sky_image_expanded = append(sky_image_expanded, f_line)
            temp_expanded = append(temp_expanded, "")
        }
    }

    // Find empty vertical lines and duplicate.
    for i_x := range len(sky_image_expanded[0]) {
        is_empty := true
        for i_y := range len(sky_image_expanded) {
            temp_expanded[i_y] += string(sky_image_expanded[i_y][i_x])
            if sky_image_expanded[i_y][i_x] == '#' {
                is_empty = false
            }
        }
        if is_empty {
            fmt.Printf("Adding column at %d\n", i_x)
            for i_y := range len(sky_image_expanded) {
                temp_expanded[i_y] += string(sky_image_expanded[i_y][i_x])
            }
        }
    }
    sky_image_expanded = temp_expanded
    fmt.Printf("Expanded image:\n")
    for _,i_line := range sky_image_expanded {
        fmt.Printf("%s\n", i_line)
    }

    total := 0
    star_list := []Star{}

    for i_y,img_line := range sky_image_expanded {
        for i_x,img_char := range img_line {
            if img_char == '#' {
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

                distance_x := current_star.position.x - candidate_star.position.x
                if distance_x < 0 {
                    distance_x *= -1
                }
                distance_y := current_star.position.y - candidate_star.position.y
                if distance_y < 0 {
                    distance_y *= -1
                }
                total_distance := distance_x + distance_y
                //fmt.Printf("Dist %d,%d => %d,%d = %d\n", current_star.position.x, current_star.position.y, candidate_star.position.x, candidate_star.position.y, total_distance)
                current_star.distance_list = append(current_star.distance_list, total_distance)
                total += total_distance
            }

        }
    }


    //fmt.Printf("Start: %v\n", start_pos)
    //fmt.Printf("numbers list 1: %v\n", numbers_list)
    //fmt.Printf("numbers map: %v\n", numbers_map)
    fmt.Printf("T: %d\n", total)
}

