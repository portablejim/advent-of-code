use std::{
    cmp::Ordering,
    io::{self},
};

fn main() {
    let stdin = io::stdin();
    let input_str = io::read_to_string(stdin)
        .unwrap_or("".to_owned())
        .replace("\r\n", "\n");
    let split_index = (&input_str)
        .find("\n\n")
        .expect("There should be a a split in the input");
    let (pg_rules_str, pg_updates_str) = input_str.split_at(split_index);
    //let pg_rules = input_str.spl[:4];
    //let pg_updates = input_str[split_index+2:];

    let pg_rules: Vec<Vec<u16>> = pg_rules_str
        .to_owned()
        .split("\n")
        .map(|l| {
            l.to_owned()
                .replace("\n", "")
                .split("|")
                .map(|r| r.parse::<u16>().expect("Bad number in input"))
                .collect()
        })
        .filter(|r: &Vec<u16>| r.len() == 2)
        .collect();
    let pg_updates: Vec<Vec<u16>> = pg_updates_str
        .to_owned()
        .split("\n")
        .filter(|l| l.len() >= 3)
        .map(|l| {
            l.to_owned()
                .replace("\n", "")
                .split(",")
                .map(|r| r.parse::<u16>().expect("Bad number in input"))
                .collect()
        })
        .collect();

    let valid_updates: Vec<Vec<u16>> = pg_updates
        .iter()
        .filter(|pg_list| {
            let failed_rule_count = pg_rules
                .iter()
                .filter(|rule| {
                    let index_a = pg_list.iter().position(|it| *it == rule[0]);
                    let index_b = pg_list.iter().position(|it| *it == rule[1]);

                    if let (Some(idx_a), Some(idx_b)) = (index_a, index_b) {
                        // If pages in wrong order
                        return idx_b <= idx_a;
                    }

                    return false;
                })
                .count();

            // Is valid if no falied rules.
            return failed_rule_count == 0;
        })
        .cloned()
        .collect();

    let bad_updates: Vec<Vec<u16>> = valid_updates
        .iter()
        .filter(|u| u.len() % 2 == 0)
        .cloned()
        .collect();
    println!("bad updates: {:?}", bad_updates);

    let total_1: i32 = valid_updates
        .iter()
        .filter(|u| u.len() % 2 == 1)
        .map(|u| u.get((u.len() + 1) / 2 - 1).unwrap_or(&0).clone() as i32)
        .sum();

    let invalid_updates: Vec<Vec<u16>> = pg_updates
        .iter()
        .filter(|pg_list| {
            let failed_rule_count = pg_rules
                .iter()
                .filter(|rule| {
                    let index_a = pg_list.iter().position(|it| *it == rule[0]);
                    let index_b = pg_list.iter().position(|it| *it == rule[1]);

                    if let (Some(idx_a), Some(idx_b)) = (index_a, index_b) {
                        // If pages in wrong order
                        return idx_b <= idx_a;
                    }

                    return false;
                })
                .count();

            // Is valid if no falied rules.
            return failed_rule_count > 0;
        })
        .cloned()
        .collect();
    let invalid_updates_sorted: Vec<Vec<u16>> = invalid_updates
        .iter()
        .map(|pg_list| {
            let mut sorted_list = pg_list.clone();
            sorted_list.sort_by(|a, b| {
                let matching_rule_q = pg_rules
                    .iter()
                    .find(|r| (r[0] == *a && r[1] == *b) || (r[0] == *b && r[1] == *a));
                if let Some(matching_rule) = matching_rule_q {
                    if matching_rule[0] == *a {
                        return Ordering::Less;
                    } else {
                        return Ordering::Greater;
                    }
                } else {
                    return Ordering::Equal;
                }
            });
            return sorted_list;
        })
        .collect();

    let total_2: i32 = invalid_updates_sorted
        .iter()
        .filter(|u| u.len() % 2 == 1)
        .map(|u| u.get((u.len() + 1) / 2 - 1).unwrap_or(&0).clone() as i32)
        .sum();

    println!("total 1: {:?}", total_1);
    println!("total 2: {:?}", total_2);
}
