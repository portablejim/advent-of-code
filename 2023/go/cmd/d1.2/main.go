package main

import "flag"
import "fmt"
import "os"
import "log"
import "strings"

func main() {
    var filename = flag.String("f", "../inputs/d1.2.sample.txt", "file to use")
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
        char_cache := ""
        f_line := f_lines[l_num]
        for c_num := range len(f_line) {
            cur_char := f_line[c_num]
            is_num := int(cur_char) >= int('0') && int(cur_char) <= int('9')
            char_num := -1
            if is_num {
                char_num = int(cur_char) - int('0')
            } else {
                char_cache = char_cache + string(cur_char)

                char_cache_r3 := ""
                if len(char_cache) >= 3 {
                    char_cache_r3 = char_cache[len(char_cache)-3:]
                }
                char_cache_r4 := ""
                if len(char_cache) >= 4 {
                    char_cache_r4 = char_cache[len(char_cache)-4:]
                }
                char_cache_r5 := ""
                if len(char_cache) >= 5 {
                    char_cache_r5 = char_cache[len(char_cache)-5:]
                }
                if char_cache_r3 == "one" {
                    char_num = 1
                }
                if char_cache_r3 == "two" {
                    char_num = 2
                }
                if char_cache_r5 == "three" {
                    char_num = 3
                }
                if char_cache_r4 == "four" {
                    char_num = 4
                }
                if char_cache_r4 == "five" {
                    char_num = 5
                }
                if char_cache_r3 == "six" {
                    char_num = 6
                }
                if char_cache_r5 == "seven" {
                    char_num = 7
                }
                if char_cache_r5 == "eight" {
                    char_num = 8
                }
                if char_cache_r4 == "nine" {
                    char_num = 9
                }
            }

            if char_num >= 0 {
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

