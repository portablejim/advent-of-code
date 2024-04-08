package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
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

func main() {
    var filename = flag.String("f", "../inputs/d2.1.sample.txt", "file to use")
    flag.Parse()
    dat, err := os.ReadFile(*filename)
    if err != nil {
        log.Fatalf("unable to read file: %f", err)
    }

    game_list := []GameInstance{}

    f_lines := strings.Split(string(dat), "\n")
    total := 0
    println("Parsing games")
    for l_num := range len(f_lines) {
        // Parse line.
        f_line := f_lines[l_num]
        game_num_str, game_def_str, found := strings.Cut(f_line, ":")
        if !found {
            continue
        }

        current_game := GameInstance{ 0, []Handful{}, Handful{} }

        game_num, err := strconv.ParseInt(game_num_str[len("Game "):], 10, 64)
        if err != nil {
            continue;
        }

        current_game.game_num = game_num

        for _, handful_str := range strings.Split(game_def_str, ";") {
            handful_ob := Handful{ 0, 0, 0 }
            for _, handful_parts := range strings.Split(handful_str, ",") {
                num_cubes_str, cube_color, parts_split := strings.Cut(strings.Trim(handful_parts, " "), " ")
                num_cubes, int_parse_err := strconv.ParseInt(strings.Trim(num_cubes_str, " "), 10, 32)

                if !parts_split || int_parse_err != nil {
                    continue
                }
                if strings.ToLower(cube_color) == "red" {
                    handful_ob.num_red = handful_ob.num_red + int32(num_cubes)
                }
                if strings.ToLower(cube_color) == "green" {
                    handful_ob.num_green = handful_ob.num_green + int32(num_cubes)
                }
                if strings.ToLower(cube_color) == "blue" {
                    handful_ob.num_blue = handful_ob.num_blue + int32(num_cubes)
                }
            }

            if handful_ob.num_red > current_game.max_nums.num_red {
                current_game.max_nums.num_red = handful_ob.num_red
            }
            if handful_ob.num_green > current_game.max_nums.num_green {
                current_game.max_nums.num_green = handful_ob.num_green
            }
            if handful_ob.num_blue > current_game.max_nums.num_blue {
                current_game.max_nums.num_blue = handful_ob.num_blue
            }
        }

        game_list = append(game_list, current_game)
    }

    for _, game_record := range game_list {
        valid_red := game_record.max_nums.num_red <= 12
        valid_green := game_record.max_nums.num_green <= 13
        valid_blue := game_record.max_nums.num_blue <= 14

        if valid_red && valid_green && valid_blue {
            fmt.Printf("Game %d: valid (%d, %d, %d)\n", game_record.game_num, game_record.max_nums.num_red, game_record.max_nums.num_green, game_record.max_nums.num_blue)
            total = total + int(game_record.game_num)
        } else {
            fmt.Printf("Game %d: invalid (%d, %d, %d)\n", game_record.game_num, game_record.max_nums.num_red, game_record.max_nums.num_green, game_record.max_nums.num_blue)
        }
    }
    fmt.Printf("T: %d\n", total)
}

