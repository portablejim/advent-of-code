use std::{collections::BinaryHeap, fmt::{write, Display}, io::{self, BufRead}, usize};


#[derive(Debug,Clone, Copy,PartialEq, Eq, PartialOrd, Ord)]
enum Direction {
    NORTH(i8),
    EAST(i8),
    SOUTH(i8),
    WEST(i8),
}

impl Display for Direction {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            Direction::NORTH(n) => write(f, format_args!("N({})", n)),
            Direction::EAST(n) => write(f, format_args!("E({})", n)),
            Direction::SOUTH(n) => write(f, format_args!("S({})", n)),
            Direction::WEST(n) => write(f, format_args!("W({})", n)),
        }
    }
}

#[derive(Debug,Clone, Copy,PartialEq, Eq, PartialOrd, Ord)]
struct GraphNodeLowest {
    cost: i16,
    prev_dir: Option<Direction>
}

#[derive(Debug,Clone, Copy,PartialEq, Eq)]
struct GraphState {
    y: usize,
    x: usize,
    is_horizontal: bool,
    lowest_cost: i16,
}

impl Ord for GraphState {
    fn cmp(&self, other: &Self) -> std::cmp::Ordering {
        other.lowest_cost.cmp(&self.lowest_cost)
            .then_with(|| self.is_horizontal.cmp(&other.is_horizontal))
            .then_with(|| self.x.cmp(&other.x))
            .then_with(|| self.y.cmp(&other.y))
    }
    
}

impl GraphState {
    fn get_graph_item<'a>(self: &'a Self, graph_data: &'a Vec<Vec<GraphWeight>>) -> Option<&GraphWeight> {
        graph_data.get(self.y.clone()).as_ref().and_then(|dr| dr.get(self.x.clone()))
    }
}

impl PartialOrd for GraphState {
    fn partial_cmp(&self, other: &Self) -> Option<std::cmp::Ordering> {
        Some(self.cmp(other))
    }
}

#[derive(Debug,Clone, Copy,PartialEq, Eq, PartialOrd, Ord)]
struct GraphWeight {
    y: u16,
    x: u16,
    cost: i16,
    lowest_horizontal: GraphNodeLowest,
    lowest_vertical: GraphNodeLowest,
    final_path: bool
}

fn parse_answer(mut passed_weights: Vec<Vec<GraphWeight>>, finish_y: usize, finish_x: usize) -> (String, Vec<Vec<GraphWeight>>) {
    let mut current_y = finish_y;
    let mut current_x = finish_x;
    let mut path = String::new();
    let mut is_horizontal = if let finish_weight = passed_weights.get(current_y).unwrap().get(current_x).unwrap() {
        finish_weight.lowest_horizontal.cost < finish_weight.lowest_vertical.cost
    }
    else {
        false
    };

    passed_weights.get_mut(current_y).unwrap().get_mut(current_x).unwrap().final_path = true;

    while current_y != 0 || current_x != 0 {
        let current_weight_row = passed_weights.get(current_y).unwrap();
        let current_weight = current_weight_row.get(current_x).unwrap();
        let current_lowest = if is_horizontal {
            current_weight.lowest_horizontal
        } else {
            current_weight.lowest_vertical
        };
        is_horizontal = !is_horizontal;

        if let Some(current_lowest_prevdir) = current_lowest.prev_dir {
            let (dx, dy, d_char, d_count) = match current_lowest_prevdir {
                Direction::NORTH(c) => (0i16, 1i16, "N", c),
                Direction::EAST(c) => (-1i16, 0i16, "E", c),
                Direction::SOUTH(c) => (0i16, -1i16, "S", c),
                Direction::WEST(c) => (1i16, 0i16, "W", c),
            };

            for i in 0..d_count {
                current_x = usize::try_from(i16::try_from(current_x).unwrap() + dx).unwrap();
                current_y = usize::try_from(i16::try_from(current_y).unwrap() + dy).unwrap();
                passed_weights.get_mut(current_y).unwrap().get_mut(current_x).unwrap().final_path = true;
                path = d_char.to_string() + path.as_str();
            }
        }
        else {
            break
        }
    }

    return (path, passed_weights)
}

fn cheapest_path(passed_weights: Vec<Vec<GraphWeight>>, min_before_turn: i32, max_before_turn: i32) -> Option<(i16, String, Vec<Vec<GraphWeight>>)> {
    let mut weights = passed_weights.clone();

    let data_height = weights.len();
    let data_width = weights.get(0).expect("Error: No data").len();
    let finish_y: usize = data_height - 1;
    let finish_x: usize = data_width - 1;

    let mut state_heap: BinaryHeap<GraphState> = BinaryHeap::new();
    state_heap.push(GraphState{ y: 0, x: 0, is_horizontal: true, lowest_cost: 0 });
    state_heap.push(GraphState{ y: 0, x: 0, is_horizontal: false, lowest_cost: 0 });

    while let Some(current_state) = state_heap.pop() {
        if current_state.x == finish_x && current_state.y == finish_y {
            let (path, final_weights) = parse_answer(weights, finish_y, finish_x);
            return Some((current_state.lowest_cost, path, final_weights));
        }

        let current_weight_item: GraphWeight = if let Some(weight_item) = current_state.get_graph_item(weights.as_ref()) {
            weight_item.clone()
        } else {
            continue
        };

        let current_weight_item_lowest = if current_state.is_horizontal {
            current_weight_item.lowest_horizontal
        }
        else {
            current_weight_item.lowest_vertical
        };

        if current_weight_item_lowest.cost < current_state.lowest_cost {
            continue;
        }

        let valid_direction_list = if current_state.is_horizontal {
            vec![Direction::NORTH(0), Direction::SOUTH(0)]
        }
        else {
            vec![Direction::EAST(0), Direction::WEST(0)]
        };

        for valid_direction in valid_direction_list {
            let mut additional_cost = 0;
            for num_steps in 0i32..max_before_turn {
                let num_steps_i = usize::try_from(num_steps + 1).unwrap();

                let next_state_opt = match valid_direction {
                    Direction::NORTH(_) => if current_state.y >= num_steps_i { Some(GraphState{y: current_state.y-num_steps_i, x: current_state.x+0, is_horizontal: false, lowest_cost: current_state.lowest_cost + current_weight_item.cost}) } else { None },
                    Direction::EAST(_) => Some(GraphState{y: current_state.y+0, x: current_state.x+num_steps_i, is_horizontal: true, lowest_cost: current_state.lowest_cost + current_weight_item.cost}),
                    Direction::SOUTH(_) => Some(GraphState{y: current_state.y+num_steps_i, x: current_state.x+0, is_horizontal: false, lowest_cost: current_state.lowest_cost + current_weight_item.cost}),
                    Direction::WEST(_) => if current_state.x >= num_steps_i { Some(GraphState{y: current_state.y+0, x: current_state.x-num_steps_i, is_horizontal: true, lowest_cost: current_state.lowest_cost + current_weight_item.cost}) } else { None },
                };

                let mut next_state = if let Some(nxt_state) = next_state_opt { nxt_state } else { continue };

                let next_weight = if let Some(x) = next_state.clone().get_graph_item(&weights) {
                    x.to_owned()
                } else {
                    continue
                };

                additional_cost += next_weight.cost;

                next_state.lowest_cost = current_state.lowest_cost + additional_cost;

                let next_direction = match valid_direction {
                    Direction::NORTH(_) => Direction::NORTH(i8::try_from(num_steps_i).unwrap()),
                    Direction::EAST(_) => Direction::EAST(i8::try_from(num_steps_i).unwrap()),
                    Direction::SOUTH(_) => Direction::SOUTH(i8::try_from(num_steps_i).unwrap()),
                    Direction::WEST(_) => Direction::WEST(i8::try_from(num_steps_i).unwrap()),
                };

                let next_weight_optional = next_state.get_graph_item(&weights);
                if next_weight_optional.is_none() {
                    continue;
                }

                if (num_steps+1) >= min_before_turn {
                    let next_weight = next_weight_optional.unwrap();
                    let next_weight_lowest = if next_state.is_horizontal {
                        next_weight.lowest_horizontal.clone()
                    } else {
                        next_weight.lowest_vertical.clone()
                    };

                    if next_state.lowest_cost < next_weight_lowest.cost {
                        state_heap.push(next_state);

                        let new_lowest = GraphNodeLowest{cost: next_state.lowest_cost, prev_dir: Some(next_direction) };

                        if let Some(weight_item) = weights.get_mut(next_state.y).and_then(|weight_row| weight_row.get_mut(next_state.x)) {
                            if next_state.is_horizontal {
                                weight_item.lowest_horizontal = new_lowest;
                            }
                            else {
                                weight_item.lowest_vertical = new_lowest;
                            }
                        }
                    }
                }

            }
        }
    }

    println!("End without solution:");
    for wl in weights {
        for wc in wl {
            println!("{:?}", wc);
        }
    }

    None
}

fn main() {
    let stdin = io::stdin();
    let input_list: Vec<String> = stdin.lock().lines()
        .filter_map(|s| s.ok())
        .collect();

    let num_offset: u8 = u8::try_from('0').unwrap();
    let initial_lowest = GraphNodeLowest{ cost: i16::MAX, prev_dir: None };
    let initial_initial = GraphNodeLowest{ cost: 0, prev_dir: None };

    let part2: bool = false;

    let (min_before_turn,max_before_turn) = if part2 {
        (4, 10)
    } else {
        (1, 3)
    };

    let weights: Vec<Vec<GraphWeight>> = input_list.iter().enumerate().map(|(line_num, line_str)| {
        line_str.bytes().enumerate().map(|(char_num, char_byte)| {
            GraphWeight{ 
                y: u16::try_from(line_num).expect("Too many rows in input."),
                x: u16::try_from(char_num).expect("Too many columns in file"),
                cost: i16::from(char_byte - num_offset),
                lowest_horizontal: if line_num == 0 && char_num == 0 { initial_initial } else { initial_lowest },
                lowest_vertical: if line_num == 0 && char_num == 0 { initial_initial } else { initial_lowest },
                final_path: false
            }
        })
        .collect()
    })
    .collect();

    let final_cost = cheapest_path(weights, min_before_turn, max_before_turn);

    if let Some((final_cost, path, final_weights)) = final_cost {
        let total = final_cost;
        println!("Final score: {}, {}", total, path);
        for wl in final_weights {
            for wc in wl {
                if wc.final_path {
                    print!("{}", wc.cost)
                }
                else {
                    print!("\u{001b}[31m{}\u{001b}[0m", wc.cost)
                }
            }
            print!("\n")
        }

    } else {
        println!("No result")
    }


}
