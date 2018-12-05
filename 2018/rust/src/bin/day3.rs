extern crate regex;

use regex::Regex;
use std::io::{self, BufRead};
use std::collections::HashMap;
use std::clone::Clone;

#[derive(PartialEq, Eq, Hash, Clone)]
struct Claim {
    number: i32,
    x: i32,
    y: i32,
    width: i32,
    height: i32
}

#[derive(PartialEq, Eq, Hash, Clone)]
struct Square {
    x: i32,
    y: i32
}

#[derive(PartialEq, Eq, Hash, Clone)]
struct SquareValue {
    count: i32,
    claims: Vec<i32>
}

fn main() {
    let re = Regex::new(r"#(\d+) @ (\d+),(\d+): (\d+)x(\d+)").unwrap();

    let stdin = io::stdin();
    let input_list: Vec<Claim> = stdin.lock().lines()
        .filter_map(|s| s.ok())
        .map(|s| {
            let caps = re.captures(&s).expect("Caps");
            Claim {
                number: caps.get(1).expect("Extract 1").as_str().parse::<i32>().expect("Parse 1"),
                x: caps.get(2).expect("Extract 2").as_str().parse::<i32>().expect("Parse 2"),
                y: caps.get(3).expect("Extract 3").as_str().parse::<i32>().expect("Parse 3"),
                width: caps.get(4).expect("Extract 4").as_str().parse::<i32>().expect("Parse 4"),
                height: caps.get(5).expect("Extract 5").as_str().parse::<i32>().expect("Parse 5")
            }
        })
        .collect();

    let mut fabric_squares: HashMap<Square, SquareValue>= HashMap::new();
    let mut not_overlaps: Vec<i32> = Vec::new();

    for input_item in &input_list {
        let mut overlaps = false;
        for cy in input_item.y..(input_item.y + input_item.height) {
            for cx in input_item.x..(input_item.x + input_item.width) {
                let target_square = Square{ x: cx, y: cy };
                let mut old_value = if let Some(val) = fabric_squares.get_mut(&target_square).cloned() {
                    overlaps = true;
                    val
                }
                else {
                    SquareValue{ count: 0, claims: Vec::new() }
                };

                old_value.count += 1;
                old_value.claims.push(input_item.number);
                
                fabric_squares.insert(target_square, old_value);
            }
        }
    }

    for input_item in &input_list {
        let num_overlapping = fabric_squares.values().filter(|sv| sv.count > 1 && sv.claims.iter().filter(|c| **c == input_item.number).count() > 0).count();
        if num_overlapping == 0
        {
            not_overlaps.push(input_item.number);
        }
    }
    let num_above_2 = fabric_squares.values().map(|sv| sv.count).filter(|v| *v >= 2).count();
    //let num_one_claim: Vec<Vec<i32>> = fabric_squares.values().filter(|sv| sv.count == 1).map(|sv| sv.claims.clone()).collect();
    
    println!("Number of squares with 2 or more claims: {}", num_above_2);
    println!("Squares that don't overlap: {:?}", not_overlaps);

}