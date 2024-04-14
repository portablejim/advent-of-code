package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type SpringGroup struct {
    type_name string
    count int
}

type SpringRow struct {
    condition_input string
    condition_groups []SpringGroup
    valid_combinations int
}

type MirrorPattern struct {
    ground_pattern []string
    h_mirror int
    v_mirror int
}

func splitCommaSepNums(num_list_str string) []int {
    num_list_trimmed := strings.Trim(num_list_str, " ")
    num_str_list := strings.Split(num_list_trimmed, ",")
    output := []int{}

    for _,num_str := range num_str_list {
        num_int, err := strconv.ParseInt(strings.Trim(num_str, " "), 10, 64)
        if err == nil {
            output = append(output, int(num_int))
        } else {
            fmt.Fprintf(os.Stderr, "Error parsing string as number: '%s' (%s)\n", num_str, num_list_str)
        }
    }

    return output
}

func findMirror(input_rows []string) int {
    if len(input_rows) == 0 {
        return -1
    }

    // Start with all columns possible
    test_cols := []int{}
    for i := range len(input_rows[0]) - 1 {
        test_cols = append(test_cols, i)
    }

    fmt.Printf("S: %v %v\n", input_rows, test_cols)
    for _,r := range input_rows {
        fmt.Printf("%s\n", r)
    }

    for _,current_row := range input_rows {
        revised_cols := []int{}
        for _,test_col := range test_cols {
            // Test for mirror
            //fmt.Printf("Testing mirror on %d\n", test_col+1)
            mirror_valid := true
            for t_left, t_right := test_col, test_col+1; t_left > -1 && t_right < len(current_row); t_left, t_right = t_left-1, t_right+1 {
                //fmt.Printf("Test: [%d] %c [%d] %c\n", t_left+1, current_row[t_left], t_right+1, current_row[t_right])
                if current_row[t_left] != current_row[t_right] {
                    mirror_valid = false
                    break
                }
            }
            if mirror_valid {
                revised_cols = append(revised_cols, test_col)
            }
        }
        if len(revised_cols) == 0 {
            // No columns can be a mirror.
            return -1
        }
        test_cols = revised_cols
    }

    if len(test_cols) > 1 {
        fmt.Printf("Multi-mirror: %v\n", test_cols)
    } else if len(test_cols) == 0 {
        fmt.Printf("No mirror in %v\n", input_rows)
        return -1
    } else {
        fmt.Printf("Mirror at %d\n", test_cols[0])
    }

    return test_cols[0]
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


func main() {
    var filename = flag.String("f", "../inputs/d10.sample1.txt", "file to use")
    //var num_copies = flag.Int("copies", 1, "number of copies")
    flag.Parse()
    dat, err := os.ReadFile(*filename)
    if err != nil {
        log.Fatalf("unable to read file: %f", err)
    }

    data_patterns := []MirrorPattern{}


    // Parse the data.
    for f_n, f_pattern := range strings.Split(string(dat), "\n\n") {
        fmt.Printf("Pattern %d\n", f_n)
        data_lines := []string{}
        for _, f_line := range strings.Split(f_pattern, "\n") {
            if len(f_line) == 0 {
                continue
            }
            data_lines = append(data_lines, f_line)
        }

        data_patterns = append(data_patterns, MirrorPattern{data_lines, -1, -1})
    }

    total := 0

    for d_n := range len(data_patterns) {
        current_pattern := *&data_patterns[d_n]
        current_pattern.h_mirror = findMirror(current_pattern.ground_pattern)
        current_pattern.v_mirror = findMirror(transposeLines(current_pattern.ground_pattern))
        h_mirror_num := current_pattern.h_mirror + 1
        v_mirror_num := current_pattern.v_mirror + 1

        current_sum := h_mirror_num + 100 * v_mirror_num
        total += current_sum

        fmt.Printf("Pattern %d: %d (%d)\n", d_n, current_sum, total)
    }


    //fmt.Printf("Start: %v\n", start_pos)
    //fmt.Printf("numbers list 1: %v\n", numbers_list)
    //fmt.Printf("numbers map: %v\n", numbers_map)
    fmt.Printf("T: %d\n", total)
}

