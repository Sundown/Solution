#[macro_use]
extern crate pest;
#[macro_use]
extern crate pest_derive;

pub mod subtle;

fn main() {
    let mut env = subtle::new_environment();

    let _ = env
        .parse_unit(&std::fs::read_to_string("code.sol").expect("Cannot read file"))
        .expect("Unsuccessful parse");

    for (_, fn_) in env.dya_fns.iter() {
        println!("{}", fn_.as_str());
    }
}
