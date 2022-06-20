#[derive(Debug, Clone, PartialEq)]
pub struct Ident {
    pub package: String,
    pub name: String,
}

#[derive(Debug, Clone, PartialEq)]
pub enum AtomicType {
    Bool,
    Char,
    Int,
    Real,
    Void,
}

impl Expression for Morpheme {
    fn as_str(&self) -> &'static str {
        match self {
            Morpheme::Bool(x) => x.to_string().as_str(),
            Morpheme::Char(x) => x.to_string().as_str(),
            Morpheme::Int(x) => x.to_string().as_str(),
            Morpheme::Real(x) => x.to_string().as_str(),
            Morpheme::Void => "Void",
        }
    }

    fn r#type(&self) -> Type {
        match self {
            Morpheme::Bool(x) => Type::Atomic(AtomicType::Bool),
            Morpheme::Char(x) => Type::Atomic(AtomicType::Char),
            Morpheme::Int(x) => Type::Atomic(AtomicType::Int),
            Morpheme::Real(x) => Type::Atomic(AtomicType::Real),
            Morpheme::Void => Type::Atomic(AtomicType::Void),
        }
    }
}

pub enum Morpheme {
    Bool(bool),
    Char(char),
    Int(i64),
    Real(f64),
    Void,
}

pub trait Expression {
    fn as_str(&self) -> &str;
    fn r#type(&self) -> Type;
}

pub struct Function {
    pub package: String,
    pub name: String,
    pub alpha: Option<Type>,
    pub omega: Type,
    pub sigma: Type,
    pub body: std::option::Option<Box<dyn Expression>>,
}

#[derive(Debug, Clone, PartialEq)]
pub enum Type {
    Atomic(AtomicType),
    Vector(Box<Type>),
}

// TODO
//pub type TypeGroup = HashSet<Type>;
