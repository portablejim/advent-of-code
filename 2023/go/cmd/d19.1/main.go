package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
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

func createRange(comp_str string, limit_str string) Range {
    limit, err := strconv.ParseInt(limit_str, 10, 64)
    if err == nil {
        if comp_str == ">" {
            return Range{limit, math.MaxInt64}
        } else if comp_str == "<" {
            return Range{math.MinInt64, limit}
        }
    }

    return Range{math.MinInt64, math.MaxInt64}
}

func evaluateWorkflow(workflow_map map[string]Workflow, current_part *Part, workflow_name string) string {
    fmt.Printf("Workflow %s | part %v\n", workflow_name, *current_part)
    current_workflow := workflow_map[workflow_name]
    for _,cur_rule := range current_workflow.rules {
        if cur_rule.partMatches(*current_part) {
            if cur_rule.action_type == "accept" || cur_rule.action_type == "reject" {
                return cur_rule.action_type
            } else if cur_rule.action_type == "chain" {
                return evaluateWorkflow(workflow_map, current_part, cur_rule.action_data)
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

    default_range := Range{math.MinInt64, math.MaxInt64}
    for _,wf_lne := range strings.Split(workflow_strs, "\n") {
        wf_lb, wf_data_raw, workflow_found := strings.Cut(wf_lne, "{")
        if !workflow_found {
            fmt.Printf("No WF: %v %v\n", wf_lb, wf_data_raw)
            continue
        }
        target_workflow := Workflow{wf_lb, []Rule{}}
        for _,wf_rule_raw := range strings.Split(strings.Trim(wf_data_raw[:len(wf_data_raw)-1], " "), ",") {
            target_rule := Rule{default_range, default_range, default_range, default_range, "", ""}
            if p_range_rule.MatchString(wf_rule_raw) {
                wf_parts := p_range_rule.FindStringSubmatch(wf_rule_raw)

                if wf_parts[1] == "x" {
                    target_rule.cond_x = createRange(wf_parts[2], wf_parts[3])
                }
                if wf_parts[1] == "m" {
                    target_rule.cond_m = createRange(wf_parts[2], wf_parts[3])
                }
                if wf_parts[1] == "a" {
                    target_rule.cond_a = createRange(wf_parts[2], wf_parts[3])
                }
                if wf_parts[1] == "s" {
                    target_rule.cond_s = createRange(wf_parts[2], wf_parts[3])
                }

                if wf_parts[4] == "A" {
                    target_rule.action_type = "accept"
                } else if wf_parts[4] == "R" {
                    target_rule.action_type = "reject"
                } else {
                    target_rule.action_type = "chain"
                    target_rule.action_data = wf_parts[4]
                }
                fmt.Printf("WFR: %v %v | %v\n", wf_lb, wf_parts, target_rule)
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
                fmt.Printf("WFO: %v %v\n", wf_lb, target_rule)
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
        fmt.Printf("Part: %v\n", part_arr)
        part_list = append(part_list, Part{val_x, val_m, val_a, val_s})
    }

	total := 0

    for cp_i,current_part := range part_list {
        eval_result := evaluateWorkflow(workflow_map, &current_part, "in")
        part_total := 0
        if eval_result == "accept" {
            part_total = int(current_part.x) + int(current_part.m) + int(current_part.a) + int(current_part.s)
        }
        total += part_total
        fmt.Printf("Current part %d: %s | %d %d\n", cp_i, eval_result, part_total, total)
    }


	//fmt.Printf("numbers list 1: %v\n", numbers_list)
	//fmt.Printf("numbers map: %v\n", numbers_map)
	fmt.Printf("T: %d\n", total)
}
