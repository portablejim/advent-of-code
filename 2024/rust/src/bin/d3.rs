use std::io::{self, BufRead};

use regex::Regex;

fn main() {
    let stdin = io::stdin();
    let input_string: String = stdin.lock().lines().filter_map(|s| s.ok()).collect();

    let mul_regex = Regex::new(r"mul\(([0-9]{1,3}),([0-9]{1,3})\)").unwrap();

    let finds: Vec<String> = mul_regex
        .find_iter(&input_string)
        .map(|caps| caps.as_str().to_owned())
        .collect();
    let parsed_finds = finds
        .iter()
        .map(|f| {
            mul_regex
                .captures_iter(f)
                .map(|n| &n[1].parse::<i32>().unwrap() * &n[2].parse::<i32>().unwrap())
                .collect::<Vec<_>>()
        })
        .flatten()
        .collect::<Vec<_>>();

    //println!("input: {:?}", input_string);
    //println!("finds: {:?}", finds);
    //println!("parsed finds: {:?}", parsed_finds);

    let disabled_regex = Regex::new(r"don't\(\).*?(do\(\))").unwrap();
    let enabled_input = disabled_regex.replace_all(&input_string, "");

    let enabled_finds: Vec<String> = mul_regex
        .find_iter(&enabled_input)
        .map(|caps| caps.as_str().to_owned())
        .collect();
    let enabled_parsed_finds = enabled_finds
        .iter()
        .map(|f| {
            mul_regex
                .captures_iter(f)
                .map(|n| &n[1].parse::<i32>().unwrap() * &n[2].parse::<i32>().unwrap())
                .collect::<Vec<_>>()
        })
        .flatten()
        .collect::<Vec<_>>();

    let total_1: i32 = parsed_finds.iter().sum();
    let total_2: i32 = enabled_parsed_finds.iter().sum();
    //println!("disabled finds: {:?}", disabled_finds);
    //println!("enabled input: {:?}", enabled_input);

    println!("total 1: {:?}", total_1);
    println!("total 2: {:?}", total_2);
}
