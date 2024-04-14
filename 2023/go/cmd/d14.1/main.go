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

type CycleNode struct {
    hash string
    value []string
    next_hash string
    visited_at int
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

func runWashCycle(input_lines []string, num_cycles int) []string {
    //dirs := []string{ "N", "W", "S", "E" }
    wash_graph := map[string]CycleNode{}
    has_hashes := false
    prev_hash := ""
    cur_hash := ""
    loop_size := -1
    loop_offset := -1

    for c_num := 0; c_num < num_cycles; c_num += 1 {
        if loop_size > -1 && loop_offset > -1 {
            target_offset := (num_cycles - loop_offset) % loop_size
            current_offset := (c_num - loop_offset) % loop_size

            fmt.Printf("Found loop: %d %d %d %d\n", target_offset, current_offset, loop_size, loop_offset)
            if current_offset == target_offset {
                fmt.Printf("G: %v\n", wash_graph[cur_hash])
                input_lines = wash_graph[cur_hash].value
                break
            }

            prev_hash = cur_hash
            cur_hash = wash_graph[prev_hash].next_hash
            continue
        }
        for i := 0; i < 4; i += 1 {
            input_lines = turnPlatform(tiltPlatform(input_lines))
        }

        prev_hash = cur_hash
        cur_hash = hash_data(input_lines)
        has_hashes = prev_hash != ""
        current_node, hash_in_graph := wash_graph[cur_hash]
        if hash_in_graph {
            if current_node.visited_at < 0 {
                fmt.Printf("Error: node should be visited %d %s->%s\n", c_num, prev_hash, cur_hash)
                panic("Error with nodes")
            }
            loop_offset = current_node.visited_at
            loop_size = c_num - current_node.visited_at
        } else {
            temp_lines := []string{}
            for _,l := range input_lines {
                temp_lines = append(temp_lines, l)
            }
            current_node = CycleNode{cur_hash, temp_lines, "", c_num}
        }
        wash_graph[cur_hash] = current_node
        if has_hashes {
            prev_node,has_prev := wash_graph[prev_hash]
            if has_prev {
                prev_node.next_hash = cur_hash
                wash_graph[prev_hash] = prev_node
            }
        }
        fmt.Printf("Wash %d: %s => %s (%d) %d %d\n", c_num, prev_hash, cur_hash, len(wash_graph), loop_size, loop_offset)
    }
    return input_lines
}

func hash_data(input_lines []string) string {
    output := ""
    current_counter := 0
    num_ops := 0
    for _,input_line := range input_lines {
        for _,input_char := range input_line {
            current_counter *= 2
            if input_char == 'O' {
                current_counter += 1
            }
            num_ops += 1

            if num_ops == 6 {
                current_shifted := '0' + current_counter
                output += string(rune(current_shifted))
                current_counter = 0
                num_ops = 0
            }
        }
    }
    if num_ops > 0 {
        for num_ops < 6 {
            current_counter *= 2
            num_ops += 1
        }
        current_shifted := '0' + current_counter
        output += string(rune(current_shifted))
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

    data_lines_turned := turnPlatform(data_lines)
    if *part2 {
        data_lines_turned = runWashCycle(data_lines_turned, 1_000_000_000)
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

    //fmt.Printf("numbers list 1: %v\n", numbers_list)
    fmt.Printf("T: %d (part2: %t)\n", total, *part2)
}

