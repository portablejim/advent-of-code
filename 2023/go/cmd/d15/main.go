package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type Lens struct {
    label string
    focal_length int
}

func hashString(input_string string) int {
        current_value := 0
        for _, cur_char := range strings.Trim(input_string, " \n") {
            current_value += int(cur_char)
            current_value *= 17
            current_value = current_value % 256
            //fmt.Printf("CV: %c %d\n", cur_char, current_value)
        }
        return current_value
}

func main() {
    var filename = flag.String("f", "../inputs/d15.sample.txt", "file to use")
    flag.Parse()
    dat, err := os.ReadFile(*filename)
    if err != nil {
        log.Fatalf("unable to read file: %f", err)
    }

    total_part1 := 0
    total_part2 := 0

    box_list := [][]Lens{}
    for i := 0; i < 256; i += 1 {
        box_list = append(box_list, []Lens{})
    }

    // Parse the data.
    for _, f_pattern := range strings.Split(string(dat), ",") {
        if len(f_pattern) == 0 {
            continue
        }
        f_pattern = strings.Trim(f_pattern, " \n")

        total_part1 += hashString(f_pattern)

        //fmt.Printf("'%s'\n", f_pattern)
        if len(f_pattern) > 1 && f_pattern[len(f_pattern)-1] == '-' {
            label := f_pattern[:(len(f_pattern)-1)]
            hash_num := hashString(label)
            //fmt.Printf("H-: %s %s %d\n", f_pattern, label, hash_num)
            temp_slot_list := []Lens{}
            current_box := box_list[hash_num]
            for _,slot := range current_box {
                if slot.label != label {
                    temp_slot_list = append(temp_slot_list, slot)
                }
            }
            box_list[hash_num] = temp_slot_list
        } else if len(f_pattern) > 2 && f_pattern[len(f_pattern)-2] == '=' {
            label := f_pattern[:len(f_pattern)-2]
            focal_length := int(f_pattern[len(f_pattern)-1] - '0')
            hash_num := hashString(label)
            //fmt.Printf("H=: %s %s %d %d\n", f_pattern, label, hash_num, focal_length)
            is_added := false
            temp_slot_list := []Lens{}
            current_box := box_list[hash_num]
            for _,slot := range current_box {
                if slot.label == label {
                    slot.focal_length = focal_length
                    //(*current_box)[s_i].focal_length = focal_length
                    //box_list[hash_num][s_i] = Lens{slot.label, focal_length}
                    is_added = true
                }
                temp_slot_list = append(temp_slot_list, slot)
            }
            box_list[hash_num] = temp_slot_list
            if !is_added {
                box_list[hash_num] = append(current_box, Lens{label, focal_length})
            }
        }
        for _,cur_box := range box_list {
            if len(cur_box) > 0 {
                //fmt.Printf("Box %d: %v\n", b_i, cur_box)
            }
        }
    }

    for b_i := range len(box_list) {
        for s_i,current_slot := range box_list[b_i] {
            power := (b_i+1) * (s_i+1) * current_slot.focal_length
            fmt.Printf("FP: %v * %d * %d = %d\n", (b_i+1), s_i+1, current_slot.focal_length, power)
            total_part2 += power
        }
    }

    fmt.Printf("T: p1 %d p2 %d\n", total_part1, total_part2)
}

