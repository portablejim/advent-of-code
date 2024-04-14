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
	count     int
}

type SpringRow struct {
	condition_input    string
	condition_groups   []SpringGroup
    broken_groups []int
	valid_combinations int
}


func countCombinations(input_str string, broken_group_list []int, cache map[string]int) int {
    // Check cache.
    cache_key := fmt.Sprintf("%v|%v", input_str, broken_group_list)
    cache_entry, has_cache_entry := cache[cache_key]
    if has_cache_entry {
        return cache_entry
    }

    // If input starts with unbroken, it can be ignored.
    for len(input_str) > 0 && input_str[0] == '.' {
        input_str = input_str[1:]
    }

    // Base case: empty string. - no string means no ways.
    if len(input_str) == 0 {
        if len(broken_group_list) == 0 {
            return 1
        }
        return 0
    }

    // Base case: no groups to process (the rest are unbroken)
    if len(broken_group_list) == 0 {
        if strings.Contains(input_str, "#") {
            // If input contains a broken value, the input is invalid.
            //fmt.Printf("countCombinations: %s = %d\n", cache_key, 0)
            return 0
        } else {
            // No broken values, thus a single possible solution.
            //fmt.Printf("countCombinations: %s = %d\n", cache_key, 1)
            return 1
        }
    }

    // Count number of 'broken' values left to find
    broken_size := 0
    for bg_i,broken_group := range broken_group_list {
        if bg_i > 0 {
            // If more than 1 group of broken values, there needs to be an unbroken value.
            broken_size += 1
        }
        broken_size += broken_group
    }

    // Base case: more required values than present values - 0 ways that is valid
    if len(input_str) < broken_size {
        return 0
    }

    // If the input string is the size of the minimal representation of the broken groups
    // then, the two options are that it matches the minimal representation (1 match)
    // or, it doesn't match the minimal represenatiton (0 matches)
    if len(input_str) == broken_size {
        test_str := ""
        for t_i,broken_group := range broken_group_list {
            if t_i > 0 {
                test_str += "."
            }
            test_str += strings.Repeat("#", broken_group)
        }
        if len(input_str) != len(test_str) {
            return 0
        }

        for c_num := range len(input_str) {
            if (input_str[c_num] == '#' && test_str[c_num] == '.') || (input_str[c_num] == '.' && test_str[c_num] == '#') {
                //fmt.Printf("countCombinations[eq]: %s, %s | %s = 0\n", cache_key, input_str, test_str)
                return 0
            }
        }
        //fmt.Printf("countCombinations[eq]: %s, %s | %s = 1\n", cache_key, input_str, test_str)
        return 1
    }

    // Now, the two options
    // 1) Consume the input (if the first char is broken)
    // 2) Don't consumer the input (the first char is unbroken)

    // If input starts with broken, it uses the group
    able_to_not_consume := true
    if input_str[0] == '#' {
        able_to_not_consume = false
    }

    target_count := 0
    count_not_consume := 0
    count_consume := 0
    if able_to_not_consume {
        // Count ways when not consuming the first broken group
        count_not_consume = countCombinations(input_str[1:], broken_group_list, cache)
        target_count += count_not_consume
    }

    // Count ways when consuming group.
    next_broken_group := broken_group_list[0]
    input_consumed := input_str[:next_broken_group]
    is_after_consume_bad := false
    if len(input_str) > next_broken_group {
        is_after_consume_bad = input_str[next_broken_group] == '#'
        //fmt.Printf("countConsumeA: %s %c %v, %v, %d\n", input_str, input_str[next_broken_group], input_consumed, is_after_consume_bad, target_count)
    }
    //fmt.Printf("countConsumeB: %s %v, %v, %d, %d\n", input_str, input_consumed, is_after_consume_bad, target_count, next_broken_group)
    if !strings.Contains(input_consumed, ".") && !is_after_consume_bad {
        // If using the broken group, the next N inputs must not contain an unbroken.
        // The input N+1 must not be a '#' (otherwise the group would be larger)
        count_consume = countCombinations(input_str[next_broken_group+1:], broken_group_list[1:], cache)
        target_count += count_consume
    }

    //fmt.Printf("fn: [%s|%s|%t] %s = %d (%d, %d)\n", input_consumed, input_str[broken_group_list[0]:], able_to_not_consume, cache_key, target_count, count_not_consume, count_consume)

    cache[cache_key] = target_count

    return target_count
}

func splitCommaSepNums(num_list_str string) []int {
	num_list_trimmed := strings.Trim(num_list_str, " ")
	num_str_list := strings.Split(num_list_trimmed, ",")
	output := []int{}

	for _, num_str := range num_str_list {
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
			condition_str = strings.Repeat(condition_str+"?", *num_copies)
			condition_str = condition_str[:len(condition_str)-1]

			groupings_str = strings.Repeat(groupings_str+",", *num_copies)
			groupings_str = groupings_str[:len(groupings_str)-1]
		}

		grouping_ints := splitCommaSepNums(groupings_str)
		grouping_objs := []SpringGroup{}
        broken_group_list := []int{}
		for _, g_i := range grouping_ints {
			//grouping_objs = append(grouping_objs, SpringGroup{ "O", -1 })
			grouping_objs = append(grouping_objs, SpringGroup{"D", g_i})
            broken_group_list = append(broken_group_list, g_i)
		}
		//grouping_objs = append(grouping_objs, SpringGroup{ "O", -1 })

		data_rows = append(data_rows, SpringRow{condition_str, grouping_objs, broken_group_list, 0})
	}

	count_cache := map[string]int{}
    for d_i, d_row := range data_rows {
		//fmt.Printf("R %d: %v\n", d_i, d_row)
        data_rows[d_i].valid_combinations = countCombinations(d_row.condition_input, d_row.broken_groups, count_cache)
	}

	total := 0

	for _, p_row := range data_rows {
		total += p_row.valid_combinations
		//fmt.Printf("R: %v\n", p_row)
	}

	//fmt.Printf("Start: %v\n", start_pos)
	//fmt.Printf("numbers list 1: %v\n", numbers_list)
	//fmt.Printf("numbers map: %v\n", numbers_map)
	fmt.Printf("T: %d\n", total)
}
