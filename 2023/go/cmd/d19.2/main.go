package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type Range struct {
    min int64
    max int64
}

type Rule struct {
    cond_x Range
    cond_m Range
    cond_a Range
    cond_s Range
    action_type string
    action_data string
    ignore_mask int8
}

func (r Rule) partMatches(p Part) bool {
    output := true
    output = output && p.x > r.cond_x.min && p.x < r.cond_x.max
    output = output && p.m > r.cond_m.min && p.m < r.cond_m.max
    output = output && p.a > r.cond_a.min && p.a < r.cond_a.max
    output = output && p.s > r.cond_s.min && p.s < r.cond_s.max
    return output
}

type Workflow struct {
    name string
    rules []Rule
}

type Part struct {
    x int64
    m int64
    a int64
    s int64
}

func createRange(comp_str string, limit_str string) (Range, int64) {
    limit, err := strconv.ParseInt(limit_str, 10, 64)
    if err == nil {
        if comp_str == ">" {
            return Range{limit, math.MaxInt64}, limit + 1
        } else if comp_str == "<" {
            return Range{math.MinInt64, limit}, limit
        }
    }

    return Range{math.MinInt64, math.MaxInt64}, 1
}

func appendArrayUnique(tgt_array []int64, tgt_val int64) []int64 {
    for i := 0; i < len(tgt_array); i += 1 {
        if tgt_array[i] == tgt_val {
            return tgt_array
        }
    }

    tgt_array = append(tgt_array, tgt_val)

    return tgt_array
}

func evaluateWorkflow(workflow_map map[string]Workflow, current_part *Part, workflow_name string, ignore_mask int8) string {
    current_workflow := workflow_map[workflow_name]
    for _,cur_rule := range current_workflow.rules {
        if ignore_mask & cur_rule.ignore_mask > 0 {
            continue
        }
        if cur_rule.partMatches(*current_part) {
            //fmt.Printf("Part matches: %v %v\n", *current_part, cur_rule)
            if cur_rule.action_type == "accept" || cur_rule.action_type == "reject" {
                return cur_rule.action_type
            } else if cur_rule.action_type == "chain" {
                return evaluateWorkflow(workflow_map, current_part, cur_rule.action_data, ignore_mask)
            } else {
                fmt.Printf("Unknown action type %s\n", cur_rule.action_type)
                return "unknown"
            }
        }
    }

    fmt.Printf("No matching action in %s for %v\n", workflow_name, current_part)
    return "nomatch"
}

func main() {
	var filename = flag.String("f", "../inputs/d19.sample1.txt", "file to use")
	//var is_verbose = flag.Bool("v", false, "verbose")
	flag.Parse()
	dat, err := os.ReadFile(*filename)
	if err != nil {
		log.Fatalf("unable to read file: %f", err)
	}

    workflow_strs, part_strs, found := strings.Cut(strings.ReplaceAll(strings.Trim(string(dat), " \r\n"), "\r\n", "\n"), "\n\n")
    if !found {
		log.Fatalf("unable to read file parts: %f", err)
    }

    p_range_rule := regexp.MustCompile("(x|m|a|s)(<|>)([0-9]+):([A-Za-z]+)")

    workflow_map := map[string]Workflow{}

    ranges_x := []int64{1, 4001}
    ranges_m := []int64{1, 4001}
    ranges_a := []int64{1, 4001}
    ranges_s := []int64{1, 4001}

    default_range := Range{math.MinInt64, math.MaxInt64}
    for _,wf_lne := range strings.Split(workflow_strs, "\n") {
        wf_lb, wf_data_raw, workflow_found := strings.Cut(wf_lne, "{")
        if !workflow_found {
            fmt.Printf("No WF: %v %v\n", wf_lb, wf_data_raw)
            continue
        }
        target_workflow := Workflow{wf_lb, []Rule{}}
        for _,wf_rule_raw := range strings.Split(strings.Trim(wf_data_raw[:len(wf_data_raw)-1], " "), ",") {
            target_rule := Rule{default_range, default_range, default_range, default_range, "", "", 0}
            if p_range_rule.MatchString(wf_rule_raw) {
                wf_parts := p_range_rule.FindStringSubmatch(wf_rule_raw)

                var cutoff_val int64
                if wf_parts[1] == "x" {
                    target_rule.ignore_mask = 8
                    target_rule.cond_x, cutoff_val = createRange(wf_parts[2], wf_parts[3])
                    //ranges_x = appendArrayUnique(ranges_x, cutoff_val-1)
                    ranges_x = appendArrayUnique(ranges_x, cutoff_val)
                    //ranges_x = appendArrayUnique(ranges_x, cutoff_val+1)
                }
                if wf_parts[1] == "m" {
                    target_rule.ignore_mask = 4
                    target_rule.cond_m, cutoff_val = createRange(wf_parts[2], wf_parts[3])
                    //ranges_m = appendArrayUnique(ranges_m, cutoff_val-1)
                    ranges_m = appendArrayUnique(ranges_m, cutoff_val)
                    //ranges_m = appendArrayUnique(ranges_m, cutoff_val+1)
                }
                if wf_parts[1] == "a" {
                    target_rule.ignore_mask = 2
                    target_rule.cond_a, cutoff_val = createRange(wf_parts[2], wf_parts[3])
                    //ranges_a = appendArrayUnique(ranges_a, cutoff_val-1)
                    ranges_a = appendArrayUnique(ranges_a, cutoff_val)
                    //ranges_a = appendArrayUnique(ranges_a, cutoff_val+1)
                }
                if wf_parts[1] == "s" {
                    target_rule.ignore_mask = 1
                    target_rule.cond_s, cutoff_val = createRange(wf_parts[2], wf_parts[3])
                    //ranges_s = appendArrayUnique(ranges_s, cutoff_val -1)
                    ranges_s = appendArrayUnique(ranges_s, cutoff_val)
                    //ranges_s = appendArrayUnique(ranges_s, cutoff_val + 1)
                }

                if wf_parts[4] == "A" {
                    target_rule.action_type = "accept"
                } else if wf_parts[4] == "R" {
                    target_rule.action_type = "reject"
                } else {
                    target_rule.action_type = "chain"
                    target_rule.action_data = wf_parts[4]
                }
            } else {
                if wf_rule_raw == "A" {

                }
                if wf_rule_raw == "A" {
                    target_rule.action_type = "accept"
                } else if wf_rule_raw == "R" {
                    target_rule.action_type = "reject"
                } else {
                    target_rule.action_type = "chain"
                    target_rule.action_data = wf_rule_raw
                }
            }
            target_workflow.rules = append(target_workflow.rules, target_rule)
        }
        workflow_map[wf_lb] = target_workflow
    }

    p_part_str := regexp.MustCompile("{x=([0-9]+),m=([0-9]+),a=([0-9]+),s=([0-9]+)}")

    part_list := []Part{}

    for _,part_arr := range p_part_str.FindAllStringSubmatch(part_strs, -1) {
        val_x, val_x_err := strconv.ParseInt(part_arr[1], 10, 64)
        val_m, val_m_err := strconv.ParseInt(part_arr[2], 10, 64)
        val_a, val_a_err := strconv.ParseInt(part_arr[3], 10, 64)
        val_s, val_s_err := strconv.ParseInt(part_arr[4], 10, 64)
        if val_x_err != nil || val_m_err != nil || val_a_err != nil || val_s_err != nil {
            continue
        }
        part_list = append(part_list, Part{val_x, val_m, val_a, val_s})
    }

	total := 0

    for cp_i,current_part := range part_list {
        eval_result := evaluateWorkflow(workflow_map, &current_part, "in", 0)
        part_total := 0
        if eval_result == "accept" {
            part_total = int(current_part.x) + int(current_part.m) + int(current_part.a) + int(current_part.s)
        }
        total += part_total
        fmt.Printf("Current part %d: %s | %d %d\n", cp_i, eval_result, part_total, total)
    }

    slices.Sort(ranges_x)
    slices.Sort(ranges_m)
    slices.Sort(ranges_a)
    slices.Sort(ranges_s)

    fmt.Printf("Ranges x: %v\n", ranges_x)
    fmt.Printf("Ranges m: %v\n", ranges_m)
    fmt.Printf("Ranges a: %v\n", ranges_a)
    fmt.Printf("Ranges s: %v\n", ranges_s)

    range_other := int64(0)

    for i_x := 0; i_x < len(ranges_x)-1; i_x += 1 {
        r_start_x := ranges_x[i_x]
        r_end_x := ranges_x[i_x+1]
        r_range_x := r_end_x - r_start_x
        fmt.Printf("X %d / %d\n", i_x, len(ranges_x)-1)

        for i_m := 0; i_m < len(ranges_m)-1; i_m += 1 {
            r_start_m := ranges_m[i_m]
            r_end_m := ranges_m[i_m+1]
            r_range_m := r_end_m - r_start_m

            for i_a := 0; i_a < len(ranges_a)-1; i_a += 1 {
                r_start_a := ranges_a[i_a]
                r_end_a := ranges_a[i_a+1]
                r_range_a := r_end_a - r_start_a

                for i_s := 0; i_s < len(ranges_s)-1; i_s += 1 {
                    r_start_s := ranges_s[i_s]
                    r_end_s := ranges_s[i_s+1]
                    r_range_s := r_end_s - r_start_s

                    eval_result := evaluateWorkflow(workflow_map, &Part{r_start_x, r_start_m, r_start_a, r_start_s}, "in", 0)
                    if eval_result == "accept" {
                        range_other += r_range_x * r_range_m * r_range_a * r_range_s
                        //fmt.Printf("Rng O: %d * %d (%d - %d) * %d * %d\n", r_range_x, r_range_m, r_start_m, r_end_m, r_range_a, r_range_m)
                    }
                }
            }
        }
    }

	//fmt.Printf("numbers list 1: %v\n", numbers_list)
	//fmt.Printf("numbers map: %v\n", numbers_map)
	fmt.Printf("T: pt1 %d pt2 %d \n", total, range_other)
}
