pub use crate::prism::*;

pub enum Morpheme {
    Bool(bool),
    Char(u8),
    Nat(u64),
    Int(i64),
    Real(f64),
    Void,
}

impl Morpheme {
    // Type as Type group
    pub fn type_g(&self) -> Type {
        return Type::of(match &self {
            Morpheme::Bool(_) => TypeInstance::Bool,
            Morpheme::Char(_) => TypeInstance::Char,
            Morpheme::Nat(_) => TypeInstance::Nat,
            Morpheme::Int(_) => TypeInstance::Int,
            Morpheme::Real(_) => TypeInstance::Real,
            Morpheme::Void => TypeInstance::Void,
        });
    }
}
