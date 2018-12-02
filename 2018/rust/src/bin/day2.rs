use std::io::{self, BufRead};
use std::collections::HashMap;
use std::collections::HashSet;

struct Counts {
    two: i32,
    three: i32
}

fn main() {
    let stdin = io::stdin();
    let input_list: Vec<String> = stdin.lock().lines()
        .filter_map(|s| s.ok())
        .collect();

    println!("Final value: {}", step1(input_list.clone()));

    println!("Different chars: {}", step2(input_list.clone()));
}

fn get_counts(box_id: &str) -> Counts {
    let mut freq: HashMap<char, i32> = HashMap::new();

    for id_char in box_id.chars() {
        let prev: i32 = freq.get(&id_char).unwrap_or(&0).clone();
        freq.insert(id_char.clone(), prev + 1);
    }

    let has_two = freq.values().any(|v| *v == 2);
    let has_three = freq.values().any(|v| *v == 3);

    Counts{ two: if has_two { 1 } else { 0 }, three: if has_three { 1 } else { 0 } }
}

fn step1(input_list: Vec<String>) -> i32 {
    let total_count = input_list.iter()
        .map(|s| get_counts(&s))
        .fold(Counts{two: 0, three: 0}, |sum, val| Counts{ two: sum.two + val.two, three: sum.three + val.three });

    total_count.two * total_count.three
}

fn get_off_by(a: &String, b: &String) -> i32
{
    let mut count: i32 = 0;
    let a_chars: Vec<char> = a.chars().collect();
    let b_chars: Vec<char> = b.chars().collect();
    let len = std::cmp::min(a_chars.len(), b_chars.len());
    for i in 0..len {
        if a_chars[i] != b_chars[i] {
            count += 1;
        }
    }

    count
}

fn get_same_chars(a: &String, b: &String) -> String
{
    let mut common: Vec<char> = Vec::new();
    let a_chars: Vec<char> = a.chars().collect();
    let b_chars: Vec<char> = b.chars().collect();
    let len = std::cmp::min(a_chars.len(), b_chars.len());
    for i in 0..len {
        if a_chars[i] == b_chars[i] {
            common.push(a_chars[i].clone());
        }
    }

    common.into_iter().collect()
}

fn get_off_by_1(input_list: &Vec<String>) -> Option<(usize, usize)> {
    for i in 0..(input_list.len() - 1) {
        for j in (i+1)..(input_list.len()) {
            if get_off_by(&input_list[i], &input_list[j]) == 1 {
                return Some((i, j));
            }
        }
    }

    None
}

fn step2(input_list: Vec<String>) -> String {
    if let Some((i_a, i_b)) = get_off_by_1(&input_list)
    {

        println!("                 {}", &input_list[i_a]);
        println!("                 {}", &input_list[i_b]);
        get_same_chars(&input_list[i_a], &input_list[i_b])
    }
    else
    {
        String::from("")
    }
}
