package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type CamelHand struct {
    hand_str string
    bid_amount int64
    hand_type_strength int
    hand_type_name string
    sort_key string
}

type BySortKey []CamelHand

func (a BySortKey) Len() int { return len(a) }
func (a BySortKey) Less(i, j int) bool { return a[i].sort_key < a[j].sort_key }
func (a BySortKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }


func splitSpaceSepNums(num_list_str string) []int64 {
    num_list_trimmed := strings.Trim(num_list_str, " ")
    num_str_list := strings.Split(num_list_trimmed, " ")
    output := []int64{}

    for _,num_str := range num_str_list {
        num_int, err := strconv.ParseInt(strings.Trim(num_str, " "), 10, 64)
        if err == nil {
            output = append(output, num_int)
        } else {
            fmt.Fprintf(os.Stderr, "Error parsing string as number: '%s' (%s)\n", num_str, num_list_str)
        }
    }

    return output
}

func toSortableHand(input_hand string) string {
    output_str := ""
    for _,c := range input_hand {
        if c == 'A' {
            output_str += "D"
        } else if c == 'K' {
            output_str += "C"
        } else if c == 'Q' {
            output_str += "B"
        } else if c == 'T' {
            output_str += "A"
        } else if c == '9' {
            output_str += "9"
        } else if c == '8' {
            output_str += "8"
        } else if c == '7' {
            output_str += "7"
        } else if c == '6' {
            output_str += "6"
        } else if c == '5' {
            output_str += "5"
        } else if c == '4' {
            output_str += "4"
        } else if c == '3' {
            output_str += "3"
        } else if c == '2' {
            output_str += "2"
        } else if c == 'J' {
            output_str += "1"
        } else {
            output_str += "0"
        }
    }

    return output_str
}

func SortString(w string) string {
    s := strings.Split(w, "")
    sort.Strings(s)
    return strings.Join(s, "")
}

func determineHandType(hand_sortable string) (int, string) {
    type_count_str := "0000000000000000"

    p_J, _ := regexp.Compile("1")
    count_Js := len(p_J.FindAllString(hand_sortable, -1))

    for _,card := range strings.Split(hand_sortable, "") {
        card_index, parse_err := strconv.ParseInt(card, 16, 8)
        if parse_err != nil {
            fmt.Fprintf(os.Stderr, "Error parsing: %v | %v\n", card, parse_err)
            continue
        }
        if card_index < 0 || card_index > 15 {
            fmt.Fprintf(os.Stderr, "Index out of range: %d\n", card_index)
            continue
        }

        target_char := type_count_str[card_index]
        target_char += 1 // Ascii magic, there will be at max 5 (Hopefully)
        type_count_str = type_count_str[:card_index] + string(target_char) + type_count_str[card_index+1:]
    }

    p_1, _ := regexp.Compile("1")
    p_2, _ := regexp.Compile("2")
    p_3, _ := regexp.Compile("3")
    p_4, _ := regexp.Compile("4")
    p_5, _ := regexp.Compile("5")
    card_counts_1s := len(p_1.FindAllString(type_count_str, -1))
    card_counts_2s := len(p_2.FindAllString(type_count_str, -1))
    card_counts_3s := len(p_3.FindAllString(type_count_str, -1))
    card_counts_4s := len(p_4.FindAllString(type_count_str, -1))
    card_counts_5s := len(p_5.FindAllString(type_count_str, -1))

    // Five of a kind.
    is_five_kind := card_counts_5s == 1
    is_five_kind = is_five_kind || (card_counts_4s == 1 && count_Js == 1)
    is_five_kind = is_five_kind || (card_counts_3s == 1 && count_Js == 2)
    is_five_kind = is_five_kind || (card_counts_2s == 1 && count_Js == 3)
    is_five_kind = is_five_kind || (card_counts_1s == 1 && count_Js == 4)
    if is_five_kind {
        return 7, "Five of a kind"
    }

    // Four of a kind
    is_four_kind := card_counts_4s == 1
    is_four_kind = is_four_kind || (card_counts_3s == 1 && count_Js == 1)
    is_four_kind = is_four_kind || (card_counts_2s == 2 && count_Js == 2)
    is_four_kind = is_four_kind || (count_Js == 3)
    if is_four_kind {
        return 6, "Four of a kind"
    }

    // Full house
    is_full_house := (card_counts_3s == 1 && card_counts_2s == 1)
    is_full_house = is_full_house || (card_counts_2s == 2 && count_Js == 1)
    is_full_house = is_full_house || (card_counts_3s == 1 && card_counts_1s == 1 && count_Js == 1)
    is_full_house = is_full_house || (card_counts_2s == 1 && card_counts_1s == 1 && count_Js == 2)
    if is_full_house {
        return 5, "Full house"
    }

    // Three of a kind.
    is_three_kind := card_counts_3s == 1 && card_counts_1s == 2
    is_three_kind = is_three_kind || (card_counts_2s == 1 && count_Js == 1)
    is_three_kind = is_three_kind ||  count_Js == 2
    if is_three_kind {
        return 4, "Three of a kind"
    }

    // Two pair
    is_two_pair := card_counts_2s == 2
    is_two_pair = is_two_pair || (card_counts_2s == 1 && card_counts_1s == 2 && count_Js == 1)
    if is_two_pair {
        return 3, "Two Pair"
    }

    // One pair
    is_one_pair := card_counts_2s == 1 && card_counts_1s == 3
    is_one_pair = is_one_pair || (card_counts_1s == 5 && count_Js == 1)
    if is_one_pair {
        return 2, "One Pair"
    }
    fmt.Printf("TC: %s %s %d %d %d %d %d %d\n", hand_sortable, type_count_str, card_counts_1s, card_counts_2s, card_counts_3s, card_counts_4s, card_counts_5s, count_Js)

    // High card.
    if card_counts_1s == 5 {
        return 1, "High Card"
    }

    return 0, "None"
}

func main() {
    var filename = flag.String("f", "../inputs/d7.sample.txt", "file to use")
    flag.Parse()
    dat, err := os.ReadFile(*filename)
    if err != nil {
        log.Fatalf("unable to read file: %f", err)
    }

    hand_list := []CamelHand{}

    // First, parse the file into data structure.
    f_line_list := strings.Split(string(dat), "\n")
    total := 0


    for _,f_line := range f_line_list {
        hand_str, bet_str, line_is_proper := strings.Cut(f_line, " ")
        if !line_is_proper {
            continue
        }

        bet_num, bet_parse_err := strconv.ParseInt(bet_str, 10, 64)
        if bet_parse_err != nil {
            fmt.Fprintf(os.Stderr, "Error parsing bet num: %v", bet_parse_err)
        }

        hand_str_sortable := toSortableHand(hand_str)
        hand_type, hand_type_name := determineHandType(hand_str_sortable)

        sort_key := fmt.Sprintf("%d%s", hand_type, hand_str_sortable)

        current_hand := CamelHand{ hand_str, bet_num, hand_type, hand_type_name, sort_key }

        hand_list = append(hand_list, current_hand)
    }

    sort.Sort(BySortKey(hand_list))

    for i,hand := range hand_list {
        ith := i + 1
        total_winnings := (ith * int(hand.bid_amount))
        total += total_winnings
        fmt.Printf("Hand: %s %s:%d (%s) = %d * %d = %d | %d\n", hand.hand_str, hand.sort_key, hand.hand_type_strength, hand.hand_type_name, hand.bid_amount, ith, total_winnings, total)
    }


    //fmt.Printf("hand_list: %v\n", hand_list)
    //fmt.Printf("numbers list 1: %v\n", numbers_list)
    //fmt.Printf("numbers map: %v\n", numbers_map)
    fmt.Printf("T: %d\n", total)
}

