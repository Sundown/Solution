#[derive(Debug, Clone, PartialEq, Eq, Hash)]
pub struct Ident {
    pub package: String,
    pub name: String,
}

pub fn base_ident(name: &str) -> Ident {
    Ident {
        package: "primary".to_string(),
        name: name.to_string(),
    }
}

pub struct Application {
    pub alpha: Option<Box<dyn Expression>>,
    pub app: Box<dyn Expression>,
    pub omega: Box<dyn Expression>,
}

impl Expression for Application {
    fn as_str(&self) -> &str {
        if self.alpha.is_some() {
            format!(
                "({} {} {})",
                self.alpha.unwrap().as_str(),
                self.app.as_str(),
                self.omega.as_str(),
            )
            .as_str()
        } else {
            format!("({} {})", self.app.as_str(), self.omega.as_str(),).as_str()
        }
    }

    fn kind(&self) -> Type {
        self.app.kind()
    }
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

    fn kind(&self) -> Type {
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

pub struct Vector {
    pub body: Vec<Box<dyn Expression>>,
}

impl Expression for Vector {
    fn as_str(&self) -> &'static str {
        format!("[{}]", self.body.into_iter().map(|x| x.as_str())).as_str()
    }

    fn kind(&self) -> Type {
        self.body[0].kind()
    }
}

pub trait Expression {
    fn as_str(&self) -> &str;
    fn kind(&self) -> Type;
}

pub struct Function {
    pub package: String,
    pub name: String,
    pub alpha: Option<Type>,
    pub omega: Type,
    pub sigma: Type,
    pub body: std::option::Option<Box<dyn Expression>>,
}

impl Expression for Function {
    fn as_str(&self) -> &'static str {
        if self.alpha.is_some() {
            format!(
                "{} {}::{} {} -> {}",
                self.alpha.unwrap().as_str(),
                self.package,
                self.name,
                self.omega.as_str(),
                self.sigma.as_str()
            )
            .as_str()
        } else {
            format!(
                "{}::{} {} -> {}",
                self.package,
                self.name,
                self.omega.as_str(),
                self.sigma.as_str()
            )
            .as_str()
        }
    }

    fn kind(&self) -> Type {
        self.sigma.clone()
    }
}

#[derive(Debug, Clone, PartialEq)]
pub enum Type {
    Atomic(AtomicType),
    Vector(Box<Type>),
}

impl Type {
    pub fn as_str(&self) -> &'static str {
        match self {
            Type::Atomic(x) => x.as_str(),
            Type::Vector(x) => format!("[{}]", x.as_str()).as_str(),
        }
    }
}

impl AtomicType {
    pub fn as_str(&self) -> &'static str {
        match self {
            AtomicType::Bool => "Bool",
            AtomicType::Char => "Char",
            AtomicType::Int => "Int",
            AtomicType::Real => "Real",
            AtomicType::Void => "Void",
        }
    }
}

// TODO
//pub type TypeGroup = HashSet<Type>;
