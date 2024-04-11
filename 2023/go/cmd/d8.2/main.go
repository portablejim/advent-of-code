package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type GraphNode struct {
    name string
    left_name string
    right_name string
    target_hash int64
    left_index int64
    right_index int64
}

func getInstruction(ins_list string, ins_num int64) byte {
    ins_len := int64(len(ins_list))
    if ins_len == 0 {
        return 0
    }
    loop_num := ins_num / ins_len
    ins_num_safe := ins_num - (ins_len * loop_num)

    return ins_list[ins_num_safe]
}

// Function to compute the Greatest Common Divisor (GCD) using the Euclidean algorithm.
func findGCD(a, b int64) int64 {
    for b != 0 {
        a, b = b, a%b
    }
    return a
}

// Function to compute the LCM.
func findLCM(a, b int64) int64 {
    return (a * b) / findGCD(a, b)
}

func main() {
    var filename = flag.String("f", "../inputs/d8.sample1.txt", "file to use")
    flag.Parse()
    dat, err := os.ReadFile(*filename)
    if err != nil {
        log.Fatalf("unable to read file: %f", err)
    }

    node_list := []GraphNode{}

    instructions_str, graph_str, instructions_found := strings.Cut(string(dat), "\n\n")
    if !instructions_found {
        fmt.Fprintf(os.Stderr, "Error parsing file")
        return
    }

    fmt.Printf("Instructions: %s\n", instructions_str)

    //p_node := regexp.MustCompile("([A-Z]{3}) = \x28([A-Z]{3}), ([A-Z]{3})\x29")
    p_node_name := regexp.MustCompile("[0-9A-Z]{3}")

    // Parse the file
    for _,ins_line_str := range strings.Split(graph_str, "\n") {
        if len(ins_line_str) == 0 {
            continue
        }
        line_parts := p_node_name.FindAllString(ins_line_str, -1)
        if len(line_parts) != 3 {
            fmt.Fprintf(os.Stderr, "Error parsing line: %v\n", ins_line_str)
            continue
        }
        //fmt.Printf("Matches: %d %s '%s'\n", len(line_parts), line_parts, ins_line_str)

        current_node := GraphNode{}
        current_node.name = line_parts[0]
        current_node.left_name = line_parts[1]
        current_node.right_name = line_parts[2]
        current_node.left_index = -1
        current_node.right_index = -1

        node_list = append(node_list, current_node)
    }

    starting_index_list := []int64{}
    for cur_i,cur_node := range node_list {
        if cur_node.left_index < 0 {
            for search_i := range len(node_list) {
                //fmt.Printf("SL: %s | %s\n", cur_node.left_name, node_list[search_i].name)
                if cur_node.left_name == node_list[search_i].name {
                    node_list[cur_i].left_index = int64(search_i)
                    break
                }
            }
        }
        if cur_node.right_index < 0 {
            for search_i := range len(node_list) {
                //fmt.Printf("SR: %s | %s\n", cur_node.right_name, node_list[search_i].name)
                if cur_node.right_name == node_list[search_i].name {
                    node_list[cur_i].right_index = int64(search_i)
                    break
                }
            }
        }
        if cur_node.name[2] == 'A' {
            starting_index_list = append(starting_index_list, int64(cur_i))
        }
    }

    if len(starting_index_list) == 0 {
        fmt.Printf("Error finding starting index")
    }

    total := -1

    step_counts_list := []int64{}
    for _,starting_index := range starting_index_list {
        fmt.Println("Processing next index")
        current_index := starting_index
        var current_ins byte
        var current_node GraphNode
        for step_i := range 100_000 {
            current_node = node_list[current_index]
            if current_node.name[2] == 'Z' {
                step_counts_list = append(step_counts_list, int64(step_i))
                break
            }
            current_ins = getInstruction(instructions_str, int64(step_i))

            if current_ins == 'L' {
                current_index = current_node.left_index
                //fmt.Printf("Node: %v => %c %s \n", current_node, current_ins, current_node.left_name)
            } else {
                current_index = current_node.right_index
                //fmt.Printf("Node: %v => %c %s \n", current_node, current_ins, current_node.right_name)
            }
            if current_index < 0 {
                current_index = starting_index
                continue
            }
        }
    }

    total_lcm := int64(len(instructions_str))
    fmt.Printf("LCM: %d\n", total_lcm)
    for _,step_count := range step_counts_list {
        total_lcm = findLCM(total_lcm, step_count)
        fmt.Printf("LCM: %d\n", total_lcm)
    }
    total = int(total_lcm)


    fmt.Printf("Step count list: %v\n", step_counts_list)
    //fmt.Printf("node_list: %v\n", node_list)
    fmt.Printf("T: %d\n", total)
}

