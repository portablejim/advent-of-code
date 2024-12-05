package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	var filename = flag.String("f", "../inputs/d3.sample.txt", "file to use")
	var part2 = flag.Bool("part2", false, "do part 2")
	flag.Parse()
	dat, err := os.ReadFile(*filename)
	if err != nil {
		log.Fatalf("unable to read file: %f\n", err)
	}
	fmt.Printf("Part 2: %v\n", *part2)

	norm_data := strings.Replace(string(dat), "\r\n", "\n", -1) + "          "

	total := 0
	enabled := true

	for i := 0; i < len(norm_data)-10; i++ {
		current_substring := norm_data[i:]
		if *part2 && current_substring[:7] == "don't()" {
			fmt.Println("disabling")
			enabled = false
			i += 6
			continue
		}
		if *part2 && current_substring[:4] == "do()" {
			fmt.Println("enabling")
			enabled = true
			i += 3
			continue
		}
		if enabled && current_substring[:4] == "mul(" {
			if unicode.IsDigit(rune(current_substring[4])) {
				num_a_start := 4
				num_a_offset := 0
				for unicode.IsDigit(rune(current_substring[num_a_start+num_a_offset])) && num_a_offset < 3 {
					num_a_offset += 1
				}
				num_a := current_substring[num_a_start : num_a_start+num_a_offset]
				if current_substring[num_a_start+num_a_offset] == ',' {
					num_b_start := num_a_start + num_a_offset + 1
					num_b_offset := 0
					if unicode.IsDigit(rune(current_substring[num_b_start])) {
						for unicode.IsDigit(rune(current_substring[num_b_start+num_b_offset])) && num_b_offset < 3 {
							num_b_offset += 1
						}
						num_b := current_substring[num_b_start : num_b_start+num_b_offset]
						if current_substring[num_b_start+num_b_offset] == ')' {
							fmt.Printf("nums: %v, %v\n", num_a, num_b)
							num_a_num, _ := strconv.ParseInt(num_a, 10, 16)
							num_b_num, _ := strconv.ParseInt(num_b, 10, 16)

							total += int(num_a_num) * int(num_b_num)
							i += num_b_start + num_b_offset
						}
					}

				}
			}
		}
	}
	fmt.Printf("total: %d", total)
}
