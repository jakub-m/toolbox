use regex::Regex;
use std::io::{self, BufRead};

fn main() {
    let re = Regex::new(r"datetime\.datetime\((\d+),\s*(\d+),\s*(\d+),\s*(\d+),\s*(\d+)(?:,\s*(\d+))?(?:,\s*(\d+))?(?:,\s*tzinfo=datetime\.timezone\.utc)?\)").unwrap();
    
    let stdin = io::stdin();
    for line in stdin.lock().lines() {
        let line = line.unwrap();
        let result = re.replace_all(&line, |caps: &regex::Captures| {
            let year: u32 = caps[1].parse().unwrap();
            let month: u32 = caps[2].parse().unwrap();
            let day: u32 = caps[3].parse().unwrap();
            let hour: u32 = caps[4].parse().unwrap();
            let minute: u32 = caps[5].parse().unwrap();
            let second: u32 = caps.get(6).map_or(0, |m| m.as_str().parse().unwrap());
            let microsecond: u32 = caps.get(7).map_or(0, |m| m.as_str().parse().unwrap());
            
            if microsecond > 0 {
                let millisecond = microsecond / 1000;
                format!("{:04}-{:02}-{:02}T{:02}:{:02}:{:02}.{:03}Z", 
                        year, month, day, hour, minute, second, millisecond)
            } else if second > 0 {
                format!("{:04}-{:02}-{:02}T{:02}:{:02}:{:02}Z", 
                        year, month, day, hour, minute, second)
            } else {
                format!("{:04}-{:02}-{:02}T{:02}:{:02}Z", 
                        year, month, day, hour, minute)
            }
        });
        println!("{}", result);
    }
}
