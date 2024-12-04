use std::io::{self, BufRead};

fn get_2d_char(list_2d: &Vec<Vec<char>>, row_num: usize, column_num: usize) -> Option<char> {
    list_2d
        .get(row_num)
        .and_then(|r| r.get(column_num))
        .cloned()
}

fn checked_get_2d_char(
    list_2d: &Vec<Vec<char>>,
    opt_row_num: Option<usize>,
    opt_column_num: Option<usize>,
) -> Option<char> {
    if let Some(row_num) = opt_row_num {
        if let Some(column_num) = opt_column_num {
            return get_2d_char(list_2d, row_num, column_num);
        }
    }

    return None;
}

fn checked_add(a: usize, b: i8) -> Option<usize> {
    usize::try_from(a as i32 + b as i32).ok()
}

fn find_xmas(
    list_2d: &Vec<Vec<char>>,
    row_num: usize,
    column_num: usize,
    row_delta: i8,
    column_delta: i8,
) -> bool {
    if let Some('M') = checked_get_2d_char(
        list_2d,
        checked_add(row_num, 1 * row_delta),
        checked_add(column_num, 1 * column_delta),
    ) {
        if let Some('A') = checked_get_2d_char(
            list_2d,
            checked_add(row_num, 2 * row_delta),
            checked_add(column_num, 2 * column_delta),
        ) {
            if let Some('S') = checked_get_2d_char(
                list_2d,
                checked_add(row_num, 3 * row_delta),
                checked_add(column_num, 3 * column_delta),
            ) {
                return true;
            }
        }
    }

    return false;
}

fn find_masx(
    list_2d: &Vec<Vec<char>>,
    row_num: usize,
    column_num: usize,
    row_delta: i8,
    column_delta: i8,
) -> bool {
    if let Some('M') = checked_get_2d_char(
        list_2d,
        checked_add(row_num, -1 * row_delta),
        checked_add(column_num, -1 * column_delta),
    ) {
        if let Some('A') =
            checked_get_2d_char(list_2d, checked_add(row_num, 0), checked_add(column_num, 0))
        {
            if let Some('S') = checked_get_2d_char(
                list_2d,
                checked_add(row_num, 1 * row_delta),
                checked_add(column_num, 1 * column_delta),
            ) {
                return true;
            }
        }
    }

    return false;
}

fn main() {
    let stdin = io::stdin();
    let input_list: Vec<Vec<char>> = stdin
        .lock()
        .lines()
        .filter_map(|s| s.ok())
        .map(|l| l.as_bytes().iter().map(|b| char::from(b.clone())).collect())
        .collect();

    let mut total_1: i32 = 0;
    let mut total_2: i32 = 0;

    for i in 0..input_list.len() {
        for j in 0..input_list.len() {
            if let Some('X') = get_2d_char(&input_list, i, j) {
                let xmas_current =
                    0 + if find_xmas(&input_list, i, j, 1, 0) {
                        1
                    } else {
                        0
                    } + if find_xmas(&input_list, i, j, -1, 0) {
                        1
                    } else {
                        0
                    } + if find_xmas(&input_list, i, j, 1, 1) {
                        1
                    } else {
                        0
                    } + if find_xmas(&input_list, i, j, -1, 1) {
                        1
                    } else {
                        0
                    } + if find_xmas(&input_list, i, j, 1, -1) {
                        1
                    } else {
                        0
                    } + if find_xmas(&input_list, i, j, -1, -1) {
                        1
                    } else {
                        0
                    } + if find_xmas(&input_list, i, j, 0, 1) {
                        1
                    } else {
                        0
                    } + if find_xmas(&input_list, i, j, 0, -1) {
                        1
                    } else {
                        0
                    } + 0;
                total_1 += xmas_current;
            }
            if let Some('A') = get_2d_char(&input_list, i, j) {
                let masx_current =
                    0 + if find_masx(&input_list, i, j, 1, 1)
                        && (find_masx(&input_list, i, j, 1, -1)
                            || find_masx(&input_list, i, j, -1, 1))
                    {
                        1
                    } else {
                        0
                    } + if find_masx(&input_list, i, j, -1, -1)
                        && (find_masx(&input_list, i, j, 1, -1)
                            || find_masx(&input_list, i, j, -1, 1))
                    {
                        1
                    } else {
                        0
                    } + 0;
                total_2 += masx_current;
            }
        }
    }

    println!("total 1: {:?}", total_1);
    println!("total 2: {:?}", total_2);
}
