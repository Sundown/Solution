#[macro_use]
extern crate pest;
#[macro_use]
extern crate pest_derive;

pub mod palisade;
pub mod prism;
pub mod subtle;

fn main() {
    let mut env = prism::Environment::new();

    let _ = env
        .parse_unit(&std::fs::read_to_string("code.sol").expect("Cannot read file"))
        .expect("Unsuccessful parse");

    env.regulate();

    for (_, fn_) in env.pre_functions.iter() {
        println!("{}", fn_.as_str());
    }

    println!("------");

    for (_, fn_) in env.dyadic_functions.iter() {
        println!("{}", fn_.as_str());
    }

    for (_, fn_) in env.monadic_functions.iter() {
        println!("{}", fn_.as_str());
    }
}
