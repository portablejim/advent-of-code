use std::{
    fmt::{self, Display},
    io::{self, BufRead},
    mem,
};

#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord)]
enum Direction {
    NORTH(i8),
    EAST(i8),
    SOUTH(i8),
    WEST(i8),
}

impl Display for Direction {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            Direction::NORTH(n) => fmt::write(f, format_args!("N({})", n)),
            Direction::EAST(n) => fmt::write(f, format_args!("E({})", n)),
            Direction::SOUTH(n) => fmt::write(f, format_args!("S({})", n)),
            Direction::WEST(n) => fmt::write(f, format_args!("W({})", n)),
        }
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord)]
struct Pos {
    x: i16,
    y: i16,
}

fn pos_move(current_pos: Pos, dir: Direction) -> Pos {
    match dir {
        Direction::NORTH(n) => Pos {
            x: current_pos.x + 0,
            y: current_pos.y - n as i16,
        },
        Direction::SOUTH(n) => Pos {
            x: current_pos.x + 0,
            y: current_pos.y + n as i16,
        },
        Direction::EAST(n) => Pos {
            x: current_pos.x + n as i16,
            y: current_pos.y + 0,
        },
        Direction::WEST(n) => Pos {
            x: current_pos.x - n as i16,
            y: current_pos.y + 0,
        },
    }
}

fn find_pos_in_2d<T>(vec2d: &Vec<Vec<T>>, pos: Pos) -> Option<&T> {
    if pos.x < 0 || pos.y < 0 {
        return None;
    }

    return vec2d
        .get(pos.y as usize)
        .and_then(|r: &Vec<T>| r.get(pos.x as usize));
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord)]
struct MapSquare {
    pos: Pos,
    is_obstacle: bool,
    last_guard_dir_n: bool,
    last_guard_dir_s: bool,
    last_guard_dir_e: bool,
    last_guard_dir_w: bool,
}

impl MapSquare {
    fn set_guard_dir(&self, last_dir: Option<Direction>) -> MapSquare {
        MapSquare {
            pos: self.pos,
            is_obstacle: self.is_obstacle,
            last_guard_dir_n: self.last_guard_dir_n
                || last_dir.is_some_and(|d| d == Direction::NORTH(1)),
            last_guard_dir_s: self.last_guard_dir_s
                || last_dir.is_some_and(|d| d == Direction::SOUTH(1)),
            last_guard_dir_e: self.last_guard_dir_e
                || last_dir.is_some_and(|d| d == Direction::EAST(1)),
            last_guard_dir_w: self.last_guard_dir_w
                || last_dir.is_some_and(|d| d == Direction::WEST(1)),
        }
    }

    fn has_visited(&self, last_dir: Option<Direction>) -> bool {
        match last_dir {
            Some(Direction::NORTH(_)) => self.last_guard_dir_n,
            Some(Direction::SOUTH(_)) => self.last_guard_dir_s,
            Some(Direction::EAST(_)) => self.last_guard_dir_e,
            Some(Direction::WEST(_)) => self.last_guard_dir_w,
            None => {
                self.last_guard_dir_n
                    || self.last_guard_dir_e
                    || self.last_guard_dir_s
                    || self.last_guard_dir_w
            }
        }
    }
}

impl Display for MapSquare {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        if self.is_obstacle {
            if self.has_visited(None) {
                fmt::write(f, format_args!("!"))
            } else {
                fmt::write(f, format_args!("#"))
            }
        } else {
            if self.has_visited(None) {
                fmt::write(f, format_args!("X"))
            } else {
                fmt::write(f, format_args!(" "))
            }
        }
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord)]
struct DirectionPosition {
    pos: Pos,
    dir: Direction,
}

#[derive(Debug, Clone, PartialEq, Eq, PartialOrd, Ord)]
struct GuardReport {
    covered_tiles: Vec<Pos>,
    loops: bool,
}

fn get_walked_positions(map_tiles: Vec<Vec<MapSquare>>) -> Vec<Pos> {
    map_tiles
        .iter()
        .flat_map(|r| {
            r.iter()
                .filter(|c| c.has_visited(None))
                .map(|c| c.pos)
                .collect::<Vec<Pos>>()
        })
        .collect()
}

fn do_patrol(
    start_guard_position: &DirectionPosition,
    mut map_tiles: Vec<Vec<MapSquare>>,
) -> GuardReport {
    let mut guard_position = start_guard_position.clone();

    loop {
        let current_tile = find_pos_in_2d::<MapSquare>(&map_tiles, guard_position.pos)
            .expect("Invalid tile")
            .clone();

        // Check if looped.
        if current_tile.has_visited(Some(guard_position.dir)) {
            return GuardReport {
                covered_tiles: get_walked_positions(map_tiles),
                loops: true,
            };
        }

        // Get next square
        let next_position = pos_move(guard_position.pos, guard_position.dir);

        // Turn if against obstale.
        let next_tile_option = find_pos_in_2d(&map_tiles, next_position);
        if let Some(next_tile) = next_tile_option {
            if next_tile.is_obstacle {
                let next_guard_dir = match guard_position.dir {
                    Direction::NORTH(n) => Direction::EAST(n),
                    Direction::SOUTH(n) => Direction::WEST(n),
                    Direction::EAST(n) => Direction::SOUTH(n),
                    Direction::WEST(n) => Direction::NORTH(n),
                };
                guard_position = DirectionPosition {
                    pos: guard_position.pos,
                    dir: next_guard_dir,
                };
            } else {
                // Mark square as visited.
                let new_tile = &current_tile.set_guard_dir(Some(guard_position.dir));
                map_tiles.get_mut(guard_position.pos.y as usize).and_then(
                    |r: &mut Vec<MapSquare>| {
                        Some(mem::replace(
                            &mut r[guard_position.pos.x as usize],
                            *new_tile,
                        ))
                    },
                );

                // Move
                guard_position = DirectionPosition {
                    pos: next_position,
                    dir: guard_position.dir,
                };
            }
        } else {
            // Mark final square as visited.
            let new_tile = &current_tile.set_guard_dir(Some(guard_position.dir));
            map_tiles
                .get_mut(guard_position.pos.y as usize)
                .and_then(|r: &mut Vec<MapSquare>| {
                    Some(mem::replace(
                        &mut r[guard_position.pos.x as usize],
                        *new_tile,
                    ))
                });

            return GuardReport {
                covered_tiles: get_walked_positions(map_tiles),
                loops: false,
            };
        }
    }
}

fn main() {
    let stdin = io::stdin();
    let input_list: Vec<Vec<char>> = stdin
        .lock()
        .lines()
        .filter_map(|s| s.ok())
        .map(|l| l.chars().collect())
        .collect();

    let mut start_guard_position: Option<DirectionPosition> = None;

    let mut map_tiles: Vec<Vec<MapSquare>> = vec![];
    for r in 0..input_list.len() {
        let row_str = &input_list[r];

        let mut row_data: Vec<MapSquare> = vec![];
        for c in 0..row_str.len() {
            let is_obstacle = row_str[c] == '#';
            row_data.push(MapSquare {
                pos: Pos {
                    x: c as i16,
                    y: r as i16,
                },
                is_obstacle: is_obstacle,
                last_guard_dir_n: false,
                last_guard_dir_s: false,
                last_guard_dir_e: false,
                last_guard_dir_w: false,
            });

            let guard_direction = match row_str[c] {
                '^' => Some(Direction::NORTH(1)),
                'v' => Some(Direction::SOUTH(1)),
                '<' => Some(Direction::EAST(1)),
                '>' => Some(Direction::WEST(1)),
                _ => None,
            };
            if guard_direction.is_some() {
                start_guard_position = Some(DirectionPosition {
                    pos: Pos {
                        x: c as i16,
                        y: r as i16,
                    },
                    dir: guard_direction.expect("Guard direction is None"),
                })
            }
        }

        map_tiles.push(row_data);
    }

    let inital_guard_position = &start_guard_position.expect("Freedom! (Guard not found)");
    let initial_patrol = do_patrol(inital_guard_position, map_tiles.clone());

    let mut candidate_locations = initial_patrol.covered_tiles.clone();
    candidate_locations.sort();
    let valid_locations: Vec<Pos> = candidate_locations
        .iter()
        .filter(|candate_loc| **candate_loc != inital_guard_position.pos)
        .filter(|candidate_loc| {
            let mut candidate_tiles = map_tiles.clone();
            let mut tile_to_modify = find_pos_in_2d(&candidate_tiles, **candidate_loc)
                .expect("Returned tiles should be valid")
                .clone();
            assert!(tile_to_modify.is_obstacle == false);
            tile_to_modify.is_obstacle = true;
            candidate_tiles
                .get_mut(candidate_loc.y as usize)
                .and_then(|r: &mut Vec<MapSquare>| {
                    Some(mem::replace(
                        &mut r[candidate_loc.x as usize],
                        tile_to_modify.clone(),
                    ))
                });
            let patrol_outcome = do_patrol(inital_guard_position, candidate_tiles);

            return patrol_outcome.loops;
        })
        .cloned()
        .collect();

    let total_1: i32 = initial_patrol.covered_tiles.len() as i32;

    let total_2: i32 = valid_locations.len() as i32;

    println!("total 1: {:?}", total_1);
    println!("total 2: {:?}", total_2);
}
