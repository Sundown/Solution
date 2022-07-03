pub use crate::subtle::*;

pub fn parse_type(t: pest::iterators::Pair<Rule>) -> Type {
    match t.as_rule() {
        Rule::typeActual => parse_type(t.into_inner().next().unwrap()),
        Rule::atomicType => Type::Atomic(match t.as_str() {
            "Bool" => AtomicType::Bool,
            "Char" => AtomicType::Char,
            "Int" => AtomicType::Int,
            "Real" => AtomicType::Real,
            "Void" => AtomicType::Void,
            _ => {
                return Type::Unknown(base_ident(t.as_str()));
            }
        }),
        Rule::vectorType => Type::Vector(Box::new(parse_type(t.into_inner().next().unwrap()))),
        _ => panic!("Unexpected type: {:?} as {:?}", t.as_str(), t.as_rule()),
    }
}
