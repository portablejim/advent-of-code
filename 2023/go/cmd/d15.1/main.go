package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

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
    var part2 = flag.Bool("part2", false, "do part 2")
    flag.Parse()
    dat, err := os.ReadFile(*filename)
    if err != nil {
        log.Fatalf("unable to read file: %f", err)
    }

    total := 0

    // Parse the data.
    for _, f_pattern := range strings.Split(string(dat), ",") {
        if len(f_pattern) == 0 {
            continue
        }
        current_value := hashString(f_pattern)
        total += current_value
        fmt.Printf("V: %d %d\n", current_value, total)
    }


    fmt.Printf("T: %d (part 2: %t)\n", total, *part2)
}

