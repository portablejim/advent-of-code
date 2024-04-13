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

func getCombinations(budget int, num_vars int) [][]int {

    if num_vars <= 0 {
        return [][]int{ }
    } else if num_vars == 1 {
        output := [][]int{ { budget } }
        return output
    } else {
        output := [][]int{}
        for n := range budget+1 {
            for _,current_comb := range getCombinations(budget-n, num_vars-1) {
                current := append([]int{ n }, current_comb...)
                output = append(output, current)
            }
        }

        return output
    }
}

func formatCombination(working_nums []int, broken_groups []SpringGroup) string {
    output := ""

    //fmt.Printf("fmt combin: %v %v\n", working_nums, broken_groups)

    for g_i,broken_group := range broken_groups {
        output += strings.Repeat(".", working_nums[g_i])
        output += strings.Repeat("#", broken_group.count)
        //fmt.Printf("pbg: %d %d | %s\n", working_nums[g_i], broken_group.count, output)
    }
    output += strings.Repeat(".", working_nums[len(broken_groups)])

    return output
}

func countBrokenString(input_str string) []int {
    output := []int{}
    active := false
    current_count := 0

    //fmt.Printf("count broken str: %s\n", input_str)
    for _,cur_char := range input_str {
        //fmt.Printf("str: %c %t\n", cur_char, cur_char == '#')
        if cur_char == '#' {
            active = true
            current_count += 1
            //fmt.Printf("active: %d", current_count)
        } else {
            if active {
                active = false
                if current_count > 0 {
                    output = append(output, current_count)
                    current_count = 0
                }
                //fmt.Printf("inactive: %v", output)
            }
        }
    }
    if active && current_count > 0 {
        output = append(output, current_count)
        current_count = 0
    }

    return output
}

func doesCombinationWork(given_str string, working_nums []int, broken_groups []SpringGroup) (bool, string, string) {
    candidate_str := formatCombination(working_nums, broken_groups)
    if len(candidate_str) != len(given_str) {
        
        return false, candidate_str, fmt.Sprintf("Combin: Different lengths - %s | %s\n", given_str, candidate_str)
    }
    for char_i, given_char := range given_str {
        candidate_char := candidate_str[char_i]
        if given_char != '?' && given_char != rune(candidate_char) {
            return false, candidate_str, fmt.Sprintf("Combin: Chars at %d differ - %c | %c\n", char_i, given_char, candidate_char)
        }
    }

    candidate_broken_counts := countBrokenString(candidate_str)
    if len(broken_groups) != len(candidate_broken_counts) {
        err := fmt.Sprintf("Combin: Different broken lengths - %d | %d\n", len(broken_groups), len(candidate_broken_counts))
        return false, candidate_str, err
    }

    for g_i,broken_group := range broken_groups {
        if broken_group.count != candidate_broken_counts[g_i] {
            err := fmt.Sprintf("Combin: Diff broken vals @ %d - %d | %d\n", g_i, broken_group.count, candidate_broken_counts[g_i])
            return false, candidate_str, err
        }
    }

    return true, candidate_str, "";
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


func main() {
    var filename = flag.String("f", "../inputs/d10.sample1.txt", "file to use")
    var num_copies = flag.Int("copies", 1, "number of copies")
    flag.Parse()
    dat, err := os.ReadFile(*filename)
    if err != nil {
        log.Fatalf("unable to read file: %f", err)
    }

    data_rows := []SpringRow{}


    // Parse the data.
    for f_y, f_line := range strings.Split(string(dat), "\n") {
        if len(f_line) == 0 {
            continue
        }
        condition_str, groupings_str, line_correct := strings.Cut(f_line, " ")
        if !line_correct {
            fmt.Printf("Error splitting line %d\n", f_y)
        }

        if *num_copies > 1 {
            condition_str = strings.Repeat(condition_str + "?", *num_copies)
            condition_str = condition_str[:len(condition_str)-1]

            groupings_str = strings.Repeat(groupings_str + ",", *num_copies)
            groupings_str = groupings_str[:len(groupings_str)-1]
            fmt.Printf("Test: %s %s\n", condition_str, groupings_str)
        }

        grouping_ints := splitCommaSepNums(groupings_str)
        grouping_objs := []SpringGroup{}
        for _,g_i := range grouping_ints {
            //grouping_objs = append(grouping_objs, SpringGroup{ "O", -1 })
            grouping_objs = append(grouping_objs, SpringGroup{ "D", g_i })
        }
        //grouping_objs = append(grouping_objs, SpringGroup{ "O", -1 })

        data_rows = append(data_rows, SpringRow{condition_str, grouping_objs, 0})
    }

    for d_i,d_row := range data_rows {
        unknown_budget := len(d_row.condition_input)
        for _,c_group := range d_row.condition_groups {
            unknown_budget -= c_group.count
        }

        combs := getCombinations(unknown_budget, len(d_row.condition_groups) + 1)
        //fmt.Printf("R: %v\n", d_row)
        for _,c := range combs {
            comb_works, _, _ := doesCombinationWork(d_row.condition_input, c, d_row.condition_groups)
            if comb_works {
                data_rows[d_i].valid_combinations += 1
                //fmt.Printf("C: %v %t %s\n", c, comb_works, comb_str)
            }
        }
    }

    total := 0

    for _,p_row := range data_rows {
        total += p_row.valid_combinations
        fmt.Printf("R: %v\n", p_row)
    }

    //fmt.Printf("Start: %v\n", start_pos)
    //fmt.Printf("numbers list 1: %v\n", numbers_list)
    //fmt.Printf("numbers map: %v\n", numbers_map)
    fmt.Printf("T: %d\n", total)
}

