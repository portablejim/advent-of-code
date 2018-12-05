extern crate regex;

use regex::Regex;
use std::io::{self, BufRead};
use std::collections::HashMap;
use std::clone::Clone;

// Hint part 1: Answer is 2851 * 44 = 125444
// 733 * 25

#[derive(Clone)]
struct Guard {
    num: u32,
    minutes_asleep: Vec<u32>
}

impl Guard {
    fn calc_sleep_time(self: &Self) -> u32 {
        self.minutes_asleep.iter().sum()
    }
}

fn main() {
    let re1 = Regex::new(r"\[(\d\d\d\d)-(\d\d)-(\d\d) (\d\d):(\d\d)\] Guard #(\d+) begins shift").unwrap();
    let re2 = Regex::new(r"\[(\d\d\d\d)-(\d\d)-(\d\d) (\d\d):(\d\d)\] falls asleep").unwrap();
    let re3 = Regex::new(r"\[(\d\d\d\d)-(\d\d)-(\d\d) (\d\d):(\d\d)\] wakes up").unwrap();

    let mut saved_guards: HashMap<u32, Guard> = HashMap::new();

    let stdin = io::stdin();
    let input_list: Vec<String> = stdin.lock().lines().filter_map(|l| l.ok()).collect();

    /*let input_list = vec!["[1518-11-01 00:00] Guard #10 begins shift",
"[1518-11-01 00:05] falls asleep",
"[1518-11-01 00:25] wakes up",
"[1518-11-01 00:30] falls asleep",
"[1518-11-01 00:55] wakes up",
"[1518-11-01 23:58] Guard #11 begins shift",
"[1518-11-01 23:58] Guard #99 begins shift",
"[1518-11-02 00:40] falls asleep",
"[1518-11-02 00:50] wakes up",
"[1518-11-03 00:05] Guard #10 begins shift",
"[1518-11-03 00:24] falls asleep",
"[1518-11-03 00:29] wakes up",
"[1518-11-04 00:02] Guard #99 begins shift",
"[1518-11-04 00:36] falls asleep",
"[1518-11-04 00:46] wakes up",
"[1518-11-05 00:03] Guard #99 begins shift",
"[1518-11-05 00:45] falls asleep",
"[1518-11-05 00:55] wakes up" ];*/

    let mut guard_number = 0;
    let mut guard_awake = false;
    let mut last_minute = 0;
    let mut minutes_asleep = vec![0; 60];

    let mut first = true;

    for input_line in input_list {
        if let Some(caps) = re1.captures(&input_line) {
            if first {
                first = false;
            }
            else {
                saved_guards.insert(guard_number, Guard{ num: guard_number, minutes_asleep: minutes_asleep });

            };
            guard_number = caps.get(6).expect("Unwrap guard number").as_str().parse::<u32>().expect("Parse guard number string");

            let current_guard = saved_guards.get(&guard_number).cloned().unwrap_or(Guard{ num: guard_number, minutes_asleep: vec![0; 60] });
            minutes_asleep = current_guard.minutes_asleep.clone();
            guard_awake = true;
        }
        else if let Some(caps) = re2.captures(&input_line) {
            if guard_awake
            {
                guard_awake = false;
                last_minute = caps.get(5).expect("Unwrap minute number").as_str().parse::<u32>().expect("Parse minute number string");
            }
        }
        else if let Some(caps) = re3.captures(&input_line) {
            let new_time = caps.get(5).expect("Unwrap minute number").as_str().parse::<u32>().expect("Parse minute number string");

            if !guard_awake
            {
                while last_minute != new_time {
                    minutes_asleep[last_minute as usize] += 1;
                    last_minute = (last_minute + 1) % 60;
                }
            }
            guard_awake = true;
        }
    }

    let mut guard_list: Vec<&Guard> = saved_guards.values().collect();
    guard_list.sort_unstable_by(|a,b| a.calc_sleep_time().cmp(&b.calc_sleep_time()).clone());

    let last_guard = guard_list[guard_list.len() - 1];
    let mapped_mins: Vec<(&u32, u32)> = last_guard.minutes_asleep.iter().zip(0..60u32).collect();
    let (_max_val, max_minute) = mapped_mins.iter().max_by_key(|(a,_b)| a).expect("Max val");
    println!("Guard: {} {:?}, {}", last_guard.num, max_minute, last_guard.num * max_minute);

    let max_guard = guard_list.iter().max_by_key(|ga| ga.minutes_asleep.iter().max()).expect("max_guard");
    let mapped_max_mins: Vec<(&u32, u32)> = max_guard.minutes_asleep.iter().zip(0..60u32).collect();
    let (_max_max_val, max_max_minute) = mapped_max_mins.iter().max_by_key(|(a,_b)| a).expect("Max val");
    println!("Guard: {} {:?}, {}", max_guard.num, max_max_minute, max_guard.num * max_max_minute);
}
