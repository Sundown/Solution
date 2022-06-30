#[macro_use]
extern crate pest;
#[macro_use]
extern crate pest_derive;

pub mod prism;
pub mod subtle;

fn main() {
    let env = prism::new_environment();

    let _ = env
        .parse_unit(&std::fs::read_to_string("code.sol").expect("Cannot read file"))
        .expect("Unsuccessful parse");

    println!("Compiled {}", env.emit_name);
}
