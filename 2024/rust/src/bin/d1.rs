use std::{
    io::{self, BufRead},
    iter,
};

fn main() {
    let stdin = io::stdin();
    let input_list: Vec<Vec<i32>> = stdin
        .lock()
        .lines()
        .filter_map(|s| s.ok())
        .map(|l| {
            l.split(" ")
                .filter_map(|s| s.parse::<i32>().ok())
                .collect::<Vec<i32>>()
        })
        .collect();

    let mut nums1: Vec<i32> = input_list
        .iter()
        .filter_map(|l| l.get(0).clone())
        .cloned()
        .collect();
    let mut nums2: Vec<i32> = input_list
        .iter()
        .filter_map(|l| l.get(1).clone())
        .cloned()
        .collect();

    nums1.sort();
    nums2.sort();

    let total_1: i32 = iter::zip(nums1.iter(), nums2.iter())
        .map(|(a, b)| if a < b { b - a } else { a - b })
        .clone()
        .sum();

    let total_2: i32 = nums1
        .iter()
        .map(|n1| n1 * i32::try_from(nums2.iter().filter(|n2| *n2 == n1).count()).unwrap_or(0))
        .clone()
        .sum();

    println!("total 1: {:?}", total_1);
    println!("total 2: {:?}", total_2);
}
