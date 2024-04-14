package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type RoadLine struct {
    rocks string
    row_num int
    rock_count int
    weight int
}


func transposeLines(input_lines []string) []string {
    output := []string{}

    if len(input_lines) > 0 {
        for c_num := range len(input_lines[0]) {
            output_line := ""
            for _,current_line := range input_lines {
                output_line += string(current_line[c_num])
            }
            output = append(output, output_line)
        }
    }

    return output
}


func reverseArray(line_split []string) []string {
    output := []string{}
    for i := len(line_split)-1; i >= 0; i -= 1 {
        output = append(output, line_split[i])
    }
    /*
    for i, j := 0, len(line_split)-1; i < j; i, j = i+1, j-1 {
        line_split[i], line_split[j] = line_split[j], line_split[i]
    }
    */

    return output
}

func turnPlatform(input_lines []string) []string {
    output := []string{}

    if len(input_lines) > 0 {
        for c_num_forward := range len(input_lines[0]) {
            c_num := c_num_forward
            output_line := ""
            for l_num_forward := range input_lines {
                current_line := input_lines[len(input_lines) - 1 - l_num_forward]
                output_line += string(current_line[c_num])
            }
            output = append(output, output_line)
        }
    }

    return output
}

func tiltPlatform(line_list []string) []string {
    output := []string{}

    for _,line := range line_list {
        line_arr := strings.Split(line, "")
        line_i := len(line) - 1
        //has_moved_rock := false
        has_solid_base := true
        for line_i >= 0 {
            if line_arr[line_i] == "." {
                look_i := line_i - 1

                // Needs to place rock on solid item.
                if has_solid_base {
                    for look_i >= 0 {
                        if line_arr[look_i] == "O" {
                            // Found a rock, move it.
                            line_arr[look_i], line_arr[line_i] = line_arr[line_i], line_arr[look_i]
                            has_solid_base = true
                            break
                        } else if line_arr[look_i] == "#" {
                            // Found a wall before a rock, jump to 
                            line_i = look_i
                            has_solid_base = true
                            break
                        }
                        has_solid_base = false
                        look_i -= 1
                    }
                }
            } else {
                has_solid_base = true
            }
            line_i -= 1
        }
        output = append(output, strings.Join(line_arr, ""))
    }

    return output
}

func countWeights(input_lines []string) []RoadLine {
    output := []RoadLine{}

    p_rock := regexp.MustCompile("O")
    for l_i, l := range input_lines {
        rock_num := len(p_rock.FindAllString(l, -1))
        line_num := len(input_lines) - l_i
        weight := rock_num * line_num
        output = append(output, RoadLine{l, line_num, rock_num, weight})
    }

    return output
}

func main() {
    var filename = flag.String("f", "../inputs/d14.sample.txt", "file to use")
    var part2 = flag.Bool("part2", false, "do part 2")
    flag.Parse()
    dat, err := os.ReadFile(*filename)
    if err != nil {
        log.Fatalf("unable to read file: %f", err)
    }

    data_lines := []string{}


    // Parse the data.
    for _, f_line := range strings.Split(string(dat), "\n") {
        if len(f_line) == 0 {
            continue
        }
        data_lines = append(data_lines, f_line)
    }

    total := 0

    dirs := []string{ "N", "W", "S", "E" }

    data_lines_turned := turnPlatform(data_lines)
    if *part2 {
        for i := 0; i < 4; i += 1 {
            data_lines_turned = turnPlatform(tiltPlatform(data_lines_turned))
            fmt.Printf("%s: %v\n", dirs[i], data_lines_turned)
        }
    } else {
        data_lines_turned = tiltPlatform(data_lines_turned)
    }
    data_lines_turned = turnPlatform(data_lines_turned)
    data_lines_turned = turnPlatform(data_lines_turned)
    data_lines_turned = turnPlatform(data_lines_turned)
    data_lines_counted := countWeights(data_lines_turned)
    for _,l := range data_lines_counted {
        fmt.Printf("R: %v\n", l)
        total += l.weight
    }


    //fmt.Printf("Start: %v\n", start_pos)
    //fmt.Printf("numbers list 1: %v\n", numbers_list)
    //fmt.Printf("numbers map: %v\n", numbers_map)
    fmt.Printf("T: %d (part2: %t)\n", total, *part2)
}

