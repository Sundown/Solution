pub use crate::prism::*;

pub enum Morpheme {
    Bool(bool),
    Char(u8),
    Word(isize),
    Nat(u64),
    Int(i64),
    Real(f64),
    Void,
}

impl Morpheme {
    pub fn kind(&self) -> Type {
        return Type::of(match &self {
            Morpheme::Bool(_) => TypeInstance::Bool,
            Morpheme::Char(_) => TypeInstance::Char,
            Morpheme::Word(_) => TypeInstance::Word,
            Morpheme::Nat(_) => TypeInstance::Nat,
            Morpheme::Int(_) => TypeInstance::Int,
            Morpheme::Real(_) => TypeInstance::Real,
            Morpheme::Void => TypeInstance::Void,
        });
    }
}
