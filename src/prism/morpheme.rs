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

    pub fn as_str(&self) -> String {
        return match &self {
            Morpheme::Bool(b) => b.to_string(),
            Morpheme::Char(c) => c.to_string(),
            Morpheme::Word(w) => w.to_string(),
            Morpheme::Nat(n) => n.to_string(),
            Morpheme::Int(i) => i.to_string(),
            Morpheme::Real(r) => r.to_string(),
            Morpheme::Void => "Void".to_string(),
        };
    }
}
