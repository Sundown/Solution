#[macro_use]
extern crate pest;
#[macro_use]
extern crate pest_derive;

pub mod prism;
pub mod subtle;

fn main() {
    let mut env = prism::new_environment();

    let astnode = env
        .parse_unit(&std::fs::read_to_string("code.sol").expect("Cannot read file"))
        .expect("Unsuccessful parse");

    println!("{:?}", &astnode);
}
