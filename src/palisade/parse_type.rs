pub use crate::palisade::*;

pub fn parse_type(t: pest::iterators::Pair<Rule>) -> prism::TypeInstance {
    match t.as_rule() {
        Rule::typeActual => parse_type(t.into_inner().next().unwrap()),
        Rule::atomicType => match t.as_str() {
            // TODO Change this to parse_ident
            "Bool" => prism::TypeInstance::Bool,
            "Char" => prism::TypeInstance::Char,
            "Int" => prism::TypeInstance::Int,
            "Real" => prism::TypeInstance::Real,
            "Void" => prism::TypeInstance::Void,
            _ => {
                panic!("Unknown type {}", t.as_str());
            }
        },
        Rule::vectorType => {
            prism::TypeInstance::Vector(Box::new(parse_type(t.into_inner().next().unwrap())))
        }
        _ => panic!("Unexpected type: {:?} as {:?}", t.as_str(), t.as_rule()),
    }
}
