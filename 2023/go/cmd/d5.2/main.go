package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type SeedLocation struct {
    seed_num int64
    seed_num_end int64
    soil_num int64
    fertilizer_num int64
    water_num int64
    light_num int64
    temp_num int64
    humidity_num int64
    location_num int64
}

type LookupTableRowRaw struct {
    dest_range_start int64
    source_range_start int64
    range_length int64
}

type BySourceRange []LookupTableRowRaw

func (a BySourceRange) Len() int { return len(a) }
func (a BySourceRange) Less(i, j int) bool { return a[i].source_range_start < a[j].source_range_start }
func (a BySourceRange) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type LookupTableRow struct {
    range_start int64
    range_offset int64
    range_length int64
}

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

func parseMapping(input_str string) []LookupTableRow {
    mapping_raw := []LookupTableRowRaw{}

    str_lines := strings.Split(input_str, "\n")
    if len(str_lines) > 1 {
        for _,str_line := range str_lines[1:] {
            if len(str_line) == 0 {
                continue
            }
            line_parts := splitSpaceSepNums(str_line)
            if len(line_parts) == 3 {
                mapping_raw = append(mapping_raw, LookupTableRowRaw{line_parts[0], line_parts[1], line_parts[2]})
            }
        }
    }

    sort.Sort(BySourceRange(mapping_raw))

    lut_list := []LookupTableRow{}

    lut_list = append(lut_list, LookupTableRow{ 0, 0, -1})

    for _,mapping_row := range mapping_raw {
        lut_row_num := 0
        for (lut_row_num + 1) < len(lut_list) && mapping_row.source_range_start >= lut_list[lut_row_num+1].range_start {
            lut_row_num += 1
        }

        existing_lut_row := lut_list[lut_row_num]

        new_lut_rows := []LookupTableRow{}

        new_continues_after := true
        existing_lut_row_new_length := existing_lut_row.range_length
        if existing_lut_row.range_length > -1 {
            existing_ending_at := existing_lut_row.range_start + existing_lut_row.range_length
            new_ending_at := mapping_row.source_range_start + mapping_row.range_length
            if existing_ending_at < new_ending_at {
                existing_lut_row_new_length = new_ending_at - existing_ending_at
            } else {
                new_continues_after = false
            }

        }

        if existing_lut_row.range_start != mapping_row.source_range_start && new_continues_after {
            // Not hitting exact, append row before
            new_lut_rows = append(new_lut_rows, existing_lut_row)
        }
        new_lut_rows = append(new_lut_rows, LookupTableRow{ mapping_row.source_range_start, mapping_row.dest_range_start - mapping_row.source_range_start, mapping_row.range_length })
        new_lut_rows = append(new_lut_rows, LookupTableRow{ mapping_row.source_range_start + mapping_row.range_length, existing_lut_row.range_offset, existing_lut_row_new_length })

        temp_lut_list := []LookupTableRow{}
        if lut_row_num > 0 {
            temp_lut_list = lut_list[:lut_row_num]
        }
        temp_lut_list = append(temp_lut_list, new_lut_rows...)
        lut_list = append(temp_lut_list, lut_list[lut_row_num+1:]...)
    }

    return lut_list
}

func useLut(lut_rows []LookupTableRow, input_num int64) int64 {
    lut_row_num := 0
    for lut_i,lut_row := range lut_rows {
        if lut_row.range_start <= input_num {
            lut_row_num = lut_i
        } else {
            break
        }
    }

    //fmt.Printf("Lookup lut for %d => %d %v\n", input_num, input_num + lut_rows[lut_row_num].range_offset, lut_rows[lut_row_num])

    return input_num + lut_rows[lut_row_num].range_offset
}

func main() {
    var filename = flag.String("f", "../inputs/d5.sample.txt", "file to use")
    flag.Parse()
    dat, err := os.ReadFile(*filename)
    if err != nil {
        log.Fatalf("unable to read file: %f", err)
    }

    seed_list := []SeedLocation{}

    // First, parse the file into data structure.
    f_sections := strings.Split(string(dat), "\n\n")
    total := 0

    if len(f_sections) != 8 {
        fmt.Printf("Wrong number of sections: %d != 8 - %v\n", len(f_sections), f_sections)
        return
    }

    seed_list_str, _ := strings.CutPrefix(f_sections[0], "seeds: ")
    seed_list_list := splitSpaceSepNums(seed_list_str)
    seed_i := 0;
    for seed_i < len(seed_list_list) {
        seed_num_start := seed_list_list[seed_i]
        seed_num_end := seed_list_list[seed_i+1]
        seed_list = append(seed_list, SeedLocation{ seed_num_start, seed_num_end, -1, -1, -1 , -1, -1, -1, -1 })
        seed_i += 2
    }

    mapping_seed_soil := parseMapping(f_sections[1])
    mapping_soil_fert := parseMapping(f_sections[2])
    mapping_fert_water := parseMapping(f_sections[3])
    mapping_water_light := parseMapping(f_sections[4])
    mapping_light_temp := parseMapping(f_sections[5])
    mapping_temp_humid := parseMapping(f_sections[6])
    mapping_humid_location := parseMapping(f_sections[7])

    fmt.Printf("Seed -> Soil: %v\n", mapping_seed_soil)
    fmt.Printf("Soil -> Fert: %v\n", mapping_soil_fert)
    fmt.Printf("Fert -> Water: %v\n", mapping_fert_water)
    fmt.Printf("Water -> Light: %v\n", mapping_water_light)
    fmt.Printf("Light -> Temp: %v\n", mapping_light_temp)
    fmt.Printf("Temp -> Humid: %v\n", mapping_temp_humid)
    fmt.Printf("Humid -> Location: %v\n", mapping_humid_location)

    min_location := int64(math.MaxInt64)
    for _,current_seed := range seed_list {
        fmt.Printf("Seed: %d\n", current_seed.seed_num)
        for i := range current_seed.seed_num_end {
            if i % (current_seed.seed_num_end / 10) == 0 {
                fmt.Printf("i: %d / %d\n", i, current_seed.seed_num)
            }
            current_seed.soil_num = useLut(mapping_seed_soil, current_seed.seed_num + i)
            current_seed.fertilizer_num = useLut(mapping_soil_fert, current_seed.soil_num)
            current_seed.water_num = useLut(mapping_fert_water, current_seed.fertilizer_num)
            current_seed.light_num = useLut(mapping_water_light, current_seed.water_num)
            current_seed.temp_num = useLut(mapping_light_temp, current_seed.light_num)
            current_seed.humidity_num = useLut(mapping_temp_humid, current_seed.temp_num)
            current_seed.location_num = useLut(mapping_humid_location, current_seed.humidity_num)
            //fmt.Printf("Seed: %v\n", current_seed)

            if min_location > current_seed.location_num {
                min_location = current_seed.location_num
            }
        }
    }

    total = int(min_location)


    //fmt.Printf("numbers list 1: %v\n", numbers_list)
    //fmt.Printf("numbers map: %v\n", numbers_map)
    fmt.Printf("T: %d\n", total)
}

