package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type Handful struct {
    num_red int32
    num_green int32
    num_blue int32
}

type GameInstance struct {
    game_num int64
    rounds []Handful
    max_nums Handful
}

type MachinePart struct {
    part_char string
    pos_x int32
    pos_y int32
}

func getNumberNum(num_map [][]int64, i_row int, i_col int) int {
    if num_map == nil || i_row >= len(num_map) {
        return -3
    }
    if i_col >= len(num_map[i_row]) {
        return -3
    }

    return int(num_map[i_row][i_col])
}

func main() {
    var filename = flag.String("f", "../inputs/d3.sample.txt", "file to use")
    flag.Parse()
    dat, err := os.ReadFile(*filename)
    if err != nil {
        log.Fatalf("unable to read file: %f", err)
    }

    // A list of numbers found in the file.
    numbers_list := []int64{}

    // A mapping of the character is the file to the number it associated with.
    numbers_map := [][]int64{}

    // Stores the list of the parts the file.
    part_list := []MachinePart{}

    f_lines := strings.Split(string(dat), "\n")
    total := 0
    for l_num := range len(f_lines) {
        // Parse line.
        f_line := f_lines[l_num]

        is_number := false
        current_number := 0

        numbers_map_line := []int64{}

        for c_num := range len(f_line) {
            l_char := f_line[c_num]
            if l_char >= '0' && l_char <= '9' {
                // Is num
                is_number = true

                char_num := l_char - '0'
                current_number = (current_number * 10) + int(char_num)

                // Store the index the current number will be stored in.
                numbers_map_line = append(numbers_map_line, int64(len(numbers_list)))

            } else {
                // When encountering a non-number char, the end of the number is reached.
                // Save the number.
                if is_number {
                    numbers_list = append(numbers_list, int64(current_number))
                    current_number = 0
                    is_number = false
                }

                if l_char != '.' {
                    // Store a negative number to indicate no number.
                    numbers_map_line = append(numbers_map_line, -1)
                    part_list = append(part_list, MachinePart{ string(l_char), int32(l_num), int32(c_num) })
                } else {
                    // Store a negative number to indicate no number.
                    numbers_map_line = append(numbers_map_line, -2)
                }
            }
        }
        // If the lend of the line is reached, the end of the number is reached.
        // Save the number.
        if is_number {
            numbers_list = append(numbers_list, int64(current_number))
            is_number = false
        }

        numbers_map = append(numbers_map, numbers_map_line)
    }

    confirmed_part_numbers := []int{}
    for range len(numbers_list) {
        confirmed_part_numbers = append(confirmed_part_numbers, -1)
    }

    // Confirm any adjacent numbers.
    for _, current_part := range part_list {
        /* Scan each part in positions round it.
        [0][1][2]
        [3][X][4]
        [5][6][7]
        */
        num_index_list := []int {
            getNumberNum(numbers_map, int(current_part.pos_x - 1), int(current_part.pos_y - 1)),
            getNumberNum(numbers_map, int(current_part.pos_x - 1), int(current_part.pos_y)),
            getNumberNum(numbers_map, int(current_part.pos_x - 1), int(current_part.pos_y + 1)),
            getNumberNum(numbers_map, int(current_part.pos_x), int(current_part.pos_y - 1)),
            getNumberNum(numbers_map, int(current_part.pos_x), int(current_part.pos_y + 1)),
            getNumberNum(numbers_map, int(current_part.pos_x + 1), int(current_part.pos_y - 1)),
            getNumberNum(numbers_map, int(current_part.pos_x + 1), int(current_part.pos_y)),
            getNumberNum(numbers_map, int(current_part.pos_x + 1), int(current_part.pos_y + 1)),
        }
        for _, num_index := range num_index_list {
            // If the numbers map contains a non-negative number,
            // it is an index of a valid number.
            if num_index >= 0 {
                // This is a number adjacent to a part/symbol.
                // So this is a confirmed part number.
                confirmed_part_numbers[num_index] = int(numbers_list[num_index])
                //numbers_list[num_index] = -1
                //total += int(numbers_list[num_index])
            }
        }
    }

    for _, candidate_num := range confirmed_part_numbers {
        if candidate_num >= 0 {
            // Add valid numbers.
            fmt.Printf("Part number: %v\n", candidate_num)
            total += int(candidate_num)
        }
    }
    //fmt.Printf("numbers list 1: %v\n", numbers_list)
    //fmt.Printf("numbers map: %v\n", numbers_map)
    fmt.Printf("T: %d\n", total)
}

