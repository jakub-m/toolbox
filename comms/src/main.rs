use regex::Regex;
use std::{
    collections::HashSet,
    env,
    io::{self, BufRead},
};

fn main() {
    match mainerr() {
        Ok(_) => return,
        Err(message) => {
            println!("{}", message);
            std::process::exit(1);
        }
    }
}

enum State {
    BeforeFirstSection,
    CollectingFirstSection,
    BeforeSecondSection,
    CollectingSecondSection,
    Finished,
}

fn mainerr() -> Result<(), String> {
    let args = get_arguments()?;
    if args.print_help {
        print_help();
        return Ok(());
    }

    let (first_section, second_section) = collect_sections()?;
    let first_section_set: HashSet<String> = HashSet::from_iter(first_section.iter().cloned());
    let second_section_set: HashSet<String> = HashSet::from_iter(second_section.iter().cloned());

    let first_only: HashSet<String> = if args.ignore_first_section {
        HashSet::new()
    } else {
        &first_section_set - &second_section_set
    };

    let second_only: HashSet<String> = if args.ignore_second_section {
        HashSet::new()
    } else {
        &second_section_set - &first_section_set
    };

    let common: HashSet<String> = if args.ignore_common {
        HashSet::new()
    } else {
        &first_section_set & &second_section_set
    };

    let mut printed_first = HashSet::new();
    let mut printed_second = HashSet::new();
    let mut printed_common = HashSet::new();

    for line in first_section.iter().chain(second_section.iter()) {
        if first_only.contains(line) {
            if !printed_first.contains(line) {
                printed_first.insert(line);
                println!("{}", line)
            }
        }
        if second_only.contains(line) {
            if !printed_second.contains(line) {
                printed_second.insert(line);
                println!("\t{}", line)
            }
        }
        if common.contains(line) {
            if !printed_common.contains(line) {
                printed_common.insert(line);
                println!("\t\t{}", line)
            }
        }
    }

    Ok(())
}

struct Args {
    ignore_first_section: bool,
    ignore_second_section: bool,
    ignore_common: bool,
    ignore_case: bool,
    print_help: bool,
}

fn get_arguments() -> Result<Args, String> {
    let mut args = Args {
        ignore_first_section: false,
        ignore_second_section: false,
        ignore_common: false,
        ignore_case: false,
        print_help: false,
    };

    let mut args_iter = env::args().into_iter();
    args_iter.next();

    let pat = Regex::new(r"^-[123i]+$").unwrap();
    while let Some(arg) = args_iter.next() {
        if arg == "-h" || arg == "--help" {
            args.print_help = true;
            return Ok(args);
        } else if pat.is_match(&arg) {
            if arg.contains("1") {
                args.ignore_first_section = true;
            }
            if arg.contains("2") {
                args.ignore_second_section = true;
            }
            if arg.contains("3") {
                args.ignore_common = true;
            }
            if arg.contains("i") {
                args.ignore_case = true;
            }
        } else {
            return Err(format!("Bad argument: {}", &arg));
        }
    }
    return Ok(args);
}

fn print_help() {
    let message = "
A version of [comm] command that works with a [s]ingle file (therefore \"comms\"). The input is split into two sections, separated by empty or blank lines. inputs The empty lines at the beginning and the end of the file are ignored. The content for parsing is taken from STDIN.

The sections do not need to be sorted, internally comms works on sets. The output is printed in the order of the original sections.

The lines are printed only once, even if the same line appears many times.

Usage: comms -[123]

Options:

    -1      Ignore lines only in the first section
    -2      Ignore lines only in the second section
    -3      Ignore lines common to both sections

";
    let message = message.trim();
    println!("{}", message);
}

fn collect_sections() -> Result<(Vec<String>, Vec<String>), String> {
    let mut state = State::BeforeFirstSection;

    let mut first_section = vec![];
    let mut second_section = vec![];

    for (i_line, line) in io::stdin().lock().lines().enumerate() {
        let line = line.unwrap();
        if line.trim().is_empty() {
            // Empty line.
            match state {
                State::BeforeFirstSection => continue,
                State::CollectingFirstSection => state = State::BeforeSecondSection,
                State::BeforeSecondSection => continue,
                State::CollectingSecondSection => state = State::Finished,
                State::Finished => continue,
            }
        } else {
            // Non-empty line.
            match state {
                State::BeforeFirstSection => {
                    state = State::CollectingFirstSection;
                    first_section.push(line);
                }
                State::CollectingFirstSection => first_section.push(line),
                State::BeforeSecondSection => {
                    state = State::CollectingSecondSection;
                    second_section.push(line);
                }
                State::CollectingSecondSection => second_section.push(line),
                State::Finished => {
                    return Err(format!(
                        "Found non-empty line {} after the second section: {}",
                        i_line + 1,
                        line
                    ))
                }
            }
        }
    }
    Ok((first_section, second_section))
}
