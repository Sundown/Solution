pub use crate::prism;
pub use crate::subtle::*;

pub fn parse_type(t: pest::iterators::Pair<Rule>) -> prism::Type {
    match t.as_rule() {
        Rule::typeActual => parse_type(t.into_inner().next().unwrap()),
        Rule::atomicType => prism::Type::Atomic(match t.as_str() {
            "Bool" => prism::AtomicType::Bool,
            "Char" => prism::AtomicType::Char,
            "Int" => prism::AtomicType::Int,
            "Real" => prism::AtomicType::Real,
            "Void" => prism::AtomicType::Void,
            _ => panic!("Unexpected type: {:?}", t.as_str()),
        }),
        Rule::vectorType => {
            prism::Type::Vector(Box::new(parse_type(t.into_inner().next().unwrap())))
        }
        _ => panic!("Unexpected type: {:?} as {:?}", t.as_str(), t.as_rule()),
    }
}
