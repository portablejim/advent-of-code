package main

import "flag"
import "fmt"
import "os"
import "log"
import "strings"

func main() {
    var filename = flag.String("f", "../inputs/d1.1.sample.txt", "file to use")
    flag.Parse()
    println(filename)
    dat, err := os.ReadFile(*filename)
    if err != nil {
        log.Fatalf("unable to read file: %f", err)
    }

    f_lines := strings.Split(string(dat), "\n")
    total := 0
    for l_num := range len(f_lines) {
        num_first := -1
        num_last := -1
        f_line := f_lines[l_num]
        for c_num := range len(f_line) {
            cur_char := f_line[c_num]
            is_char := int(cur_char) >= int('0') && int(cur_char) <= int('9')
            if is_char {
                char_num := int(cur_char) - int('0')
                if num_first < 0 {
                    num_first = char_num
                }
                num_last = char_num
            }
        }
        if num_first >= 0 && num_last >= 0 {
            total += (num_first * 10) + num_last
            fmt.Printf("A: %d\n", (num_first * 10) + num_last)
        }
    }
    fmt.Printf("T: %d\n", total)
}

