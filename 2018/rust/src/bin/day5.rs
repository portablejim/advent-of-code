
fn main() {
    let input_string = include_str!("../../../05/input").to_owned();
    let mut input_string_chars = input_string.bytes().filter(|b| (*b >= 97 && *b <= 122) || (*b >= 65 && *b <= 90) ).collect::<Vec<u8>>();
    let input_string_chars2 = input_string_chars.clone();
    let mut changes = true;

    while changes {
        let (input_string_chars_out, changes_out) = do_pass(input_string_chars.clone());
        input_string_chars = input_string_chars_out;
        changes = changes_out;
    }
    println!(" ");

    let min_len = (97u8..123).map(|lower| step1(input_string_chars2.clone().iter().map(|n| *n).filter(|n| *n != lower && *n != (lower - 32) ).collect()).len()).min();

    println!("Step 1: {}", input_string_chars.len());
    println!("Step 2: {}", min_len.unwrap_or(0));
}

fn step1(mut input_string_chars: Vec<u8>) -> Vec<u8>
{
    let mut changes = true;

    while changes {
        let (input_string_chars_out, changes_out) = do_pass(input_string_chars.clone());
        input_string_chars = input_string_chars_out;
        changes = changes_out;
    }

    input_string_chars
}

fn step2_for_char(mut input_string_chars: Vec<u8>, in_char: u8) -> Vec<u8>
{
    let mut changes = true;

    while changes {
        let (input_string_chars_out, changes_out) = do_pass_for_char(input_string_chars.clone(), in_char);
        input_string_chars = input_string_chars_out;
        changes = changes_out;
    }

    input_string_chars
}

fn do_pass(input_string_chars: Vec<u8>) -> (Vec<u8>, bool)
{
    let mut changed = false;
    let mut new_char_vec: Vec<u8> = Vec::with_capacity(input_string_chars.len());
    let mut i = 0usize;
    while i < input_string_chars.len() {
        let ca = input_string_chars.get(i).expect("Unwrap ca");
        if let Some(cb) = input_string_chars.get(i + 1) {
            match (*ca as i32) - (*cb as i32) {
                // same char different case => Skip.
                32 | -32 => {
                    // Advance extra one to skip whole pair.
                    i += 1;
                    changed = true;
                },
                _ => {
                    new_char_vec.push(*ca);
                }
            }
        }
        else {
            // At end, so just add last char.
            new_char_vec.push(*ca);
        }

        // Normal iteration
        i += 1;
    }

    (new_char_vec, changed)
}

fn do_pass_for_char(input_string_chars: Vec<u8>, char_lower: u8) -> (Vec<u8>, bool)
{
    let char_upper = char_lower - 32;

    let mut changed = false;
    let mut new_char_vec: Vec<u8> = Vec::with_capacity(input_string_chars.len());
    let mut i = 0usize;
    while i < input_string_chars.len() {
        let ca = input_string_chars.get(i).expect("Unwrap ca");
        if let Some(cb) = input_string_chars.get(i + 1) {
            match (ca, cb) {
                // same char different case => Skip.
                (a,b) if ((char_lower == *a &&  char_upper == *b ) || (char_lower == *b &&  char_upper == *a)) => {
                    // Advance extra one to skip whole pair.
                    i += 1;
                    changed = true;
                },
                _ => {
                    new_char_vec.push(*ca);
                }
            }
        }
        else {
            // At end, so just add last char.
            new_char_vec.push(*ca);
        }

        // Normal iteration
        i += 1;
    }

    (new_char_vec, changed)
}
