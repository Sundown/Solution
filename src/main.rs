#[macro_use]
extern crate pest;
#[macro_use]
extern crate pest_derive;

pub mod prism;
pub mod subtle;

fn main() {
    let mut env = prism::new_environment();
    let unparsed_file = std::fs::read_to_string("code.sol").expect("Cannot read file");
    let astnode = subtle::parse_unit(&unparsed_file).expect("Unsuccessful parse");
    println!("{:?}", &astnode);
}
