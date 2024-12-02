use std::io::{self, BufRead};

#[derive(Debug, Clone, Copy)]
struct Report {
    last_value: i32,
    safe: bool,
    increasing: Option<bool>,
}

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

    let passed_reports = input_list.iter().map(|r| {
        let empty_report: Option<Report> = None;
        let current_report: Option<Report> = r.iter().fold(empty_report, |candidate_report, el| {
            if let Some(current_report) = candidate_report {
                if current_report.safe == false {
                    // Already unsafe.
                    return candidate_report;
                } else if *el == current_report.last_value
                    || (*el <= current_report.last_value
                        && current_report.increasing.is_some_and(|v| v == true))
                    || (*el >= current_report.last_value
                        && current_report.increasing.is_some_and(|v| v == false))
                {
                    // Contains invalid value
                    return Some(Report {
                        last_value: *el,
                        safe: false,
                        increasing: current_report.increasing,
                    });
                } else {
                    let current_increasing = current_report.last_value < *el;
                    let difference = i32::abs(current_report.last_value - *el);
                    let current_safe = difference <= 3;
                    return Some(Report {
                        last_value: *el,
                        safe: current_safe,
                        increasing: Some(current_increasing),
                    });
                }
            } else {
                // Initial report value
                Some(Report {
                    last_value: *el,
                    safe: true,
                    increasing: None,
                })
            }
        });

        if current_report.is_some_and(|r| r.safe) {
            return 1;
        } else {
            return 0;
        }
    });

    println!(
        "Passed reports (Part 1): {:?}",
        passed_reports.clone().sum::<i32>()
    );
}
