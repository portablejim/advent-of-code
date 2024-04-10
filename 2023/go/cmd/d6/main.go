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

type RaceInstance struct {
    time int
    distance int
}

// tr ( rl - tr) = distance
// (tr * rl) - (tr * tr) = distance
// (tr * rl) = (tr * tr) + distance
// 0 = (1)(tr*tr) - (rl)*tr + distance
func getSolutions(ca int, cb int, cc int) (int, int, bool) {
    discrim := math.Sqrt(float64((cb * cb) - (4 * ca * cc)))

    if discrim == math.NaN() {
        return 0, 0, false
    }

    init_num := float64(-cb)/float64(2*ca)

    sol_a_f := init_num - (discrim / float64(2*ca))
    sol_b_f := init_num + (discrim / float64(2*ca))

    sol_a := int(math.Ceil(sol_a_f))
    if float64(sol_a) == sol_a_f {
        sol_a += 1
    }
    sol_b := int(math.Floor(sol_b_f))
    if float64(sol_b) == sol_b_f {
        sol_b -= 1
    }

    //fmt.Printf("i: %f | d: %f | a: %d | b: %d | af: %f | bf: %f\n", init_num, discrim, sol_a, sol_b, sol_a_f, sol_b_f)

    return sol_a, sol_b, true
}

func main() {
    var filename = flag.String("f", "../inputs/d6.sample.txt", "file to use")
    var part2 = flag.Bool("part2", false, "do part 2")
    flag.Parse()
    dat, err := os.ReadFile(*filename)
    if err != nil {
        log.Fatalf("unable to read file: %f\n", err)
    }


    total := 0

    str_time_full, str_distance_full, splits_found := strings.Cut(string(dat), "\n")
    if !splits_found {
        fmt.Printf("Error splitting file\n")
    }

    if *part2 {
        str_time_full = strings.ReplaceAll(str_time_full, " ", "")
        str_distance_full = strings.ReplaceAll(str_distance_full, " ", "")
    }

    r, _ := regexp.Compile("[0-9]+")
    time_str_split := r.FindAllString(str_time_full, -1)
    distance_str_split := r.FindAllString(str_distance_full, -1)

    fmt.Printf("Times: %d %v\n", len(time_str_split), time_str_split)
    fmt.Printf("Distances: %d %v\n", len(distance_str_split), distance_str_split)

    race_list := []RaceInstance{}

    for i,t := range time_str_split {
        t_parsed, t_parsed_err := strconv.ParseInt(t, 10, 64)
        if t_parsed_err != nil {
            continue
        }

        dist := -1
        if i < len(distance_str_split) {
            dist_temp, d_parsed_err := strconv.ParseInt(distance_str_split[i], 10, 64)
            if d_parsed_err == nil {
                dist = int(dist_temp)
            }
        }

        current_race := RaceInstance{ int(t_parsed), dist }
        race_list = append(race_list, current_race)

    }

    for i,r := range race_list {
        if total == 0 {
            total = 1
        }
        a,b,s := getSolutions(1, -r.time, r.distance)
        num_solutions := 0
        if s && a <= b {
            num_solutions = b - a + 1
        }
        total = total * num_solutions
        fmt.Printf("Race %d [%d,%d]: %d\n", i, r.time, r.distance, num_solutions)
    }

    fmt.Printf("T: %d\n", total)
}

