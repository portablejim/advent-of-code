use std::io::{self, BufRead};
use std::collections::HashSet;

fn main() {
    let stdin = io::stdin();
    let input_list: Vec<i32> = stdin.lock().lines()
        .map(|s| s.unwrap_or(String::from("0")).parse::<i32>().unwrap())
        .collect();

    println!("Final value: {}", step1(input_list.clone()));

    println!("First repeated value: {}", step2(input_list.clone()));
}

fn step1(input_list: Vec<i32>) -> i32
{
    input_list.iter().fold(0i32, |sum, val| sum + val)
}

fn step2(input_list: Vec<i32>) -> i32
{
    let mut current_val = 0i32;
    let mut previous_vals: HashSet<i32> = HashSet::new();
    loop {
        for item in &input_list {
            if previous_vals.contains(&current_val)  {
                return current_val;
            }
            else {
                previous_vals.insert(current_val);
                current_val += item;
            }
        }
    }
}

