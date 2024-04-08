package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type ScratchCard struct {
    cardNumber int
    winning_nums []int
    have_nums []int
    matching_nums int
    points int
    num_copies int
}

func main() {
    var filename = flag.String("f", "../inputs/d4.sample.txt", "file to use")
    flag.Parse()
    dat, err := os.ReadFile(*filename)
    if err != nil {
        log.Fatalf("unable to read file: %f", err)
    }

    cards := []ScratchCard{}

    // First, parse the file into data structure.
    f_lines := strings.Split(string(dat), "\n")
    total := 0
    for l_num := range len(f_lines) {
        // Parse line.
        f_line := f_lines[l_num]

        card_num_full, other_numbers, card_found := strings.Cut(f_line, ": ")
        if !card_found {
            continue
        }

        // Read the card number
        card_num_str, card_num_success := strings.CutPrefix(card_num_full, "Card ")
        if !card_num_success {
            fmt.Fprintf(os.Stderr, "%d|Error reading card num from '%s'\n", l_num, card_num_full)
        }
        card_num, card_num_parse_err := strconv.ParseInt(strings.Trim(card_num_str, " "), 10, 64)
        if card_num_parse_err != nil {
            fmt.Fprintf(os.Stderr, "%d|Error parsing card num from '%s' %v\n", l_num, card_num_str, card_num_parse_err)
        }

        current_card := ScratchCard{ int(card_num), []int{}, []int{}, 0, 0 , 1 }

        // Extract the winning numbers
        winning_numbers_full, have_numbers_full, numbers_cut_ok := strings.Cut(other_numbers, " | ")
        if !numbers_cut_ok {
            fmt.Fprintf(os.Stderr, "%d|Error cutting numbers line from '%s'\n", l_num, other_numbers)
        }

        winning_numbers_strlist := strings.Split(winning_numbers_full, " ")
        for _,num := range winning_numbers_strlist {
            num_str := strings.Trim(num, " ")
            num_int, num_parse_err := strconv.ParseInt(num_str, 10, 64)
            if num_parse_err == nil {
                current_card.winning_nums = append(current_card.winning_nums, int(num_int))
            }
        }

        // Extract the numbers we have.
        have_numbers_strlist := strings.Split(have_numbers_full, " ")
        for _,num := range have_numbers_strlist {
            num_str := strings.Trim(num, " ")
            num_int, num_parse_err := strconv.ParseInt(num_str, 10, 64)
            if num_parse_err == nil {
                current_card.have_nums = append(current_card.have_nums, int(num_int))
            }
        }
        
        // Count numbers and points.
        for _,h_num := range current_card.have_nums {
            for _,w_num := range current_card.winning_nums {
                if h_num == w_num {
                    current_card.matching_nums += 1
                    current_card.points = 1 << (current_card.matching_nums - 1)
                    break
                }
            }
        }

        //fmt.Printf("Card %d: %v\n", l_num, current_card)
        cards = append(cards, current_card)
    }

    // Callculate card copies.
    for card_n := range len(cards) {
        current_card := cards[card_n]
        if current_card.matching_nums > 0 {

            // Add to counts of next cards
            // The number of copies of the current card tell us how many new
            // copies of the next cards we get.
            for i := range current_card.matching_nums {
                add_c_index := card_n + i + 1
                if add_c_index < len(cards) {
                    cards[add_c_index].num_copies += current_card.num_copies
                }
            }
        }
    }

    // Now we have everything calculate total.
    for _,current_card := range cards {
        fmt.Printf("Card %d: %d %d\n", current_card.cardNumber, current_card.matching_nums, current_card.num_copies)
        total += current_card.num_copies
    }


    //fmt.Printf("numbers list 1: %v\n", numbers_list)
    //fmt.Printf("numbers map: %v\n", numbers_map)
    fmt.Printf("T: %d\n", total)
}

