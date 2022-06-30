use std::collections::HashSet;

#[derive(Debug, Clone, PartialEq, Eq, Hash)]
pub struct Ident {
    pub package: String,
    pub name: String,
}

impl Ident {
    pub fn as_str(&self) -> String {
        format!("{}::{}", &self.package, &self.name)
    }
}

pub fn base_ident(name: &str) -> Ident {
    Ident {
        package: "primary".to_string(),
        name: name.to_string(),
    }
}

pub struct Application {
    pub alpha: Option<Box<dyn Expression>>,
    pub app: Ident,
    pub omega: Box<dyn Expression>,
}

impl Expression for Application {
    fn as_str(&self) -> String {
        format!(
            "({}{} {})",
            match &self.alpha {
                Some(a) => format!("{} ", a.as_str()),
                None => "".to_string(),
            },
            self.app.as_str(),
            self.omega.as_str(),
        )
    }

    fn kind(&self) -> Type {
        // TODO - this is wrong, change to app once there is helper
        self.omega.kind()
    }
}

#[derive(Debug, Clone, PartialEq, Hash, Eq)]
pub enum AtomicType {
    Bool,
    Char,
    Int,
    Real,
    Void,
}

impl Expression for Morpheme {
    fn as_str(&self) -> String {
        match self {
            Morpheme::Bool(x) => x.to_string(),
            Morpheme::Char(x) => x.to_string(),
            Morpheme::Int(x) => x.to_string(),
            Morpheme::Real(x) => x.to_string(),
            Morpheme::Void => "Void".to_string(),
        }
    }

    fn kind(&self) -> Type {
        match self {
            Morpheme::Bool(_) => Type::Atomic(AtomicType::Bool),
            Morpheme::Char(_) => Type::Atomic(AtomicType::Char),
            Morpheme::Int(_) => Type::Atomic(AtomicType::Int),
            Morpheme::Real(_) => Type::Atomic(AtomicType::Real),
            Morpheme::Void => Type::Atomic(AtomicType::Void),
        }
    }
}

#[derive(Clone)]
pub enum Morpheme {
    Bool(bool),
    Char(char),
    Int(i64),
    Real(f64),
    Void,
}

pub struct Vector {
    pub body: Vec<Box<dyn Expression>>,
}

impl Expression for Vector {
    fn as_str(&self) -> String {
        format!(
            "{}",
            self.body
                .iter()
                .map(|x| x.as_str())
                .collect::<Vec<_>>()
                .join(" ")
        )
    }

    fn kind(&self) -> Type {
        self.body[0].kind()
    }
}

pub trait Expression {
    fn as_str(&self) -> String;
    fn kind(&self) -> Type;
}

pub struct Function {
    pub package: String,
    pub name: String,
    pub alpha: Option<Type>,
    pub omega: Type,
    pub sigma: Type,
    pub body: std::option::Option<Vec<Box<dyn Expression>>>,
}

impl Expression for Function {
    fn as_str(&self) -> String {
        format!(
            "{}{}::{} {} -> {}\n\t{}",
            match &self.alpha {
                Some(x) => format!("{} ", x.as_str()),
                None => "".to_string(),
            },
            self.package,
            self.name,
            self.omega.as_str(),
            self.sigma.as_str(),
            match &self.body {
                Some(x) => x
                    .iter()
                    .map(|x| x.as_str())
                    .collect::<Vec<_>>()
                    .join("\n\t"),
                None => "".to_string(),
            },
        )
    }

    fn kind(&self) -> Type {
        self.sigma.clone()
    }
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct TypeGroup {
    pub gamma: HashSet<Type>,
}

impl TypeGroup {
    pub fn new() -> TypeGroup {
        TypeGroup {
            gamma: HashSet::new(),
        }
    }

    pub fn universal(&self) -> bool {
        self.gamma.len() == 0
    }
}

#[derive(Debug, Clone, PartialEq, Hash, Eq)]
pub enum Type {
    Atomic(AtomicType),
    Vector(Box<Type>),
}

impl Type {
    pub fn as_str(&self) -> String {
        match self {
            Type::Atomic(x) => x.as_str(),
            Type::Vector(x) => format!("[{}]", x.as_str()),
        }
    }
}

impl AtomicType {
    pub fn as_str(&self) -> String {
        match self {
            AtomicType::Bool => "Bool",
            AtomicType::Char => "Char",
            AtomicType::Int => "Int",
            AtomicType::Real => "Real",
            AtomicType::Void => "Void",
        }
        .to_string()
    }
}

// TODO
//pub type TypeGroup = HashSet<Type>;
