package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type HistoryRecord struct {
    values []int64
    diff_values [][]int64
    next_value_set bool
    next_value int64
}

func splitSpaceSepNums(num_list_str string) []int64 {
    num_list_trimmed := strings.Trim(num_list_str, " ")
    num_str_list := strings.Split(num_list_trimmed, " ")
    output := []int64{}

    for _,num_str := range num_str_list {
        if len(num_str) == 0 {
            continue
        }
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
    var filename = flag.String("f", "../inputs/d9.sample.txt", "file to use")
    var part2 = flag.Bool("part2", false, "part 2")
    flag.Parse()
    dat, err := os.ReadFile(*filename)
    if err != nil {
        log.Fatalf("unable to read file: %f", err)
    }

    history_list := []HistoryRecord{}

    f_lines := strings.Split(string(dat), "\n")
    total  := int64(0)
    for l_num := range len(f_lines) {
        // Parse line.
        f_line := f_lines[l_num]

        if len(f_line) <= 0 {
            continue
        }

        line_split := splitSpaceSepNums(f_line)
        if *part2 {
            for i, j := 0, len(line_split)-1; i < j; i, j = i+1, j-1 {
                line_split[i], line_split[j] = line_split[j], line_split[i]
            }
        }

        history_list = append(history_list, HistoryRecord{line_split, [][]int64{}, false, 0})
    }

    for hist_rec_num,hist_rec := range history_list {
        current_diffs := hist_rec.values
        // Calculate diffs.
        for true {
            new_diff := []int64{}

            has_non_zero_value := false
            for d_i := range (len(current_diffs) - 1) {
                diff_value := current_diffs[d_i + 1] - current_diffs[d_i]
                new_diff = append(new_diff, diff_value)
                has_non_zero_value = has_non_zero_value || diff_value != 0
            }
            hist_rec.diff_values = append(hist_rec.diff_values, new_diff)

            if has_non_zero_value {
                current_diffs = new_diff
            } else {
                break
            }
        }

        if len(hist_rec.diff_values) < 2 {
            fmt.Printf("# %d Not erough diff values: %v", hist_rec_num, hist_rec.diff_values)
        }

        // Add value.
        add_i := len(hist_rec.diff_values) - 1

        hist_rec.diff_values[add_i] = append(hist_rec.diff_values[add_i], 0)

        add_i -= 1

        for add_i >= 0 {
            next_num := hist_rec.diff_values[add_i][len(hist_rec.diff_values[add_i])-1] + hist_rec.diff_values[add_i+1][len(hist_rec.diff_values[add_i+1])-1]
            hist_rec.diff_values[add_i] = append(hist_rec.diff_values[add_i], next_num)

            add_i -= 1
        }

        hist_rec.next_value = hist_rec.values[len(hist_rec.values)-1] + hist_rec.diff_values[0][len(hist_rec.diff_values[0])-1]

        total += hist_rec.next_value

        fmt.Printf("V: %v\nD: %v\nN: %d\n", hist_rec.values, hist_rec.diff_values, hist_rec.next_value)
        
    }



    //fmt.Printf("numbers list 1: %v\n", numbers_list)
    //fmt.Printf("numbers map: %v\n", numbers_map)
    fmt.Printf("T: %d\n", total)
}

