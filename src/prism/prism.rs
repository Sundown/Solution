pub use crate::palisade;
pub use crate::prism::*;
use std::collections::{HashMap, HashSet};

pub struct Environment {
    pub types: HashMap<Ident, TypeInstance>,
    pub pre_functions: HashMap<Ident, palisade::Function>,
    pub monadic_functions: HashMap<Ident, MonadicFunction>,
    pub dyadic_functions: HashMap<Ident, DyadicFunction>,
}

impl Environment {
    pub fn new() -> Environment {
        Environment {
            types: {
                let mut h = HashMap::new();
                h.insert(Ident::new("", "Bool"), TypeInstance::Bool);
                h.insert(Ident::new("", "Char"), TypeInstance::Char);
                h.insert(Ident::new("", "Word"), TypeInstance::Word);
                h.insert(Ident::new("", "Nat"), TypeInstance::Nat);
                h.insert(Ident::new("", "Int"), TypeInstance::Int);
                h.insert(Ident::new("", "Real"), TypeInstance::Real);
                h
            },
            pre_functions: HashMap::new(),
            monadic_functions: HashMap::new(),
            dyadic_functions: HashMap::new(),
        }
    }
}

pub enum Expression {
    Morpheme(Morpheme),
    Vector(Vector),
    Monadic(MonadicApplication),
    Dyadic(DyadicApplication),
}

pub struct Vector {
    pub element_type: TypeInstance,
    pub body: Vec<Expression>,
}
impl Vector {
    pub fn as_str(&self) -> String {
        let mut s = String::new();
        s.push_str("[");
        for e in &self.body {
            s.push_str(&format!("{}", e.as_str()));
        }
        s.push_str("]");
        s
    }
}
impl Expression {
    pub fn kind(&self) -> Type {
        match &self {
            Expression::Morpheme(m) => m.kind(),
            Expression::Vector(v) => Type::of(v.element_type.clone()),
            Expression::Monadic(m) => m.kind(),
            Expression::Dyadic(d) => d.kind(),
        }
    }

    pub fn as_str(&self) -> String {
        match &self {
            Expression::Morpheme(m) => m.as_str(),
            Expression::Vector(v) => v.as_str(),
            Expression::Monadic(m) => m.as_str(),
            Expression::Dyadic(d) => d.as_str(),
        }
    }
}

#[derive(Hash, PartialEq, Eq, Debug, Clone, PartialOrd, Ord)]
pub struct Ident {
    pub package: String,
    pub name: String,
}

impl Ident {
    pub fn as_str(&self) -> String {
        format!("{}::{}", self.package, self.name)
            .as_str()
            .to_string()
    }

    pub fn new(package: &str, name: &str) -> Self {
        Ident {
            package: package.to_string(),
            name: name.to_string(),
        }
    }
}

pub struct MonadicApplication {
    pub phi: Ident,
    pub sigma_t: Type,
    pub omega_t: Type,
    pub omega: Box<Expression>,
}

impl MonadicApplication {
    pub fn kind(&self) -> Type {
        self.sigma_t.clone()
    }

    pub fn as_str(&self) -> String {
        format!("{} {}", self.phi.as_str(), self.omega.as_str(),)
    }
}

pub struct DyadicApplication {
    pub alpha: Box<Expression>,
    pub alpha_t: Type,
    pub phi: Ident,
    pub sigma_t: Type,
    pub omega: Box<Expression>,
    pub omega_t: Type,
}

impl DyadicApplication {
    pub fn kind(&self) -> Type {
        self.sigma_t.clone()
    }

    pub fn as_str(&self) -> String {
        format!(
            "{} {} {}",
            self.alpha.as_str(),
            self.phi.as_str(),
            self.omega.as_str(),
        )
    }
}

pub struct MonadicFunction {
    pub ident: Ident,
    pub omega: Type,
    pub sigma: Type,
    pub body: Option<Vec<Expression>>,
    pub attrs: FuncAttrs,
}

impl MonadicFunction {
    pub fn kind(&self) -> Type {
        self.sigma.clone()
    }

    pub fn as_str(&self) -> String {
        format!(
            "{} {} -> {}\n\t{}",
            self.ident.as_str(),
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
}

pub struct DyadicFunction {
    pub ident: Ident,
    pub alpha: Type,
    pub omega: Type,
    pub sigma: Type,
    pub body: Option<Vec<Expression>>,
    pub attrs: FuncAttrs,
}

impl DyadicFunction {
    pub fn kind(&self) -> Type {
        self.sigma.clone()
    }

    pub fn as_str(&self) -> String {
        format!(
            "{} {} {} -> {}\n\t{}",
            self.alpha.as_str(),
            self.ident.as_str(),
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
}

pub struct FuncAttrs {
    pub inline: bool, // For LLVM func attr
    pub elide: bool,  // Subtle stage should ignore
}

impl FuncAttrs {
    pub fn new() -> FuncAttrs {
        FuncAttrs {
            inline: false,
            elide: false,
        }
    }
}

#[derive(Eq, PartialEq, Clone, Debug)]
pub struct Type {
    pub gamma: HashSet<TypeInstance>,
}

impl Type {
    pub fn any(&self) -> bool {
        self.gamma.contains(&TypeInstance::Any)
    }

    pub fn allows(&self, t: &TypeInstance) -> bool {
        self.gamma.contains(t)
    }

    pub fn new_any() -> Type {
        Type {
            gamma: HashSet::from([TypeInstance::Any]),
        }
    }

    pub fn none(&self) -> bool {
        self.gamma.len() > 0
    }

    pub fn single(&self) -> Option<TypeInstance> {
        match self.gamma.iter().next() {
            Some(t) => Some(t.clone()),
            None => None,
        }
    }

    pub fn of(t: TypeInstance) -> Type {
        Type {
            gamma: HashSet::from([t]),
        }
    }

    pub fn as_str(&self) -> String {
        format!(
            "{{{}}}",
            self.gamma
                .iter()
                .map(|x| x.as_str())
                .collect::<Vec<_>>()
                .join(" ")
        )
    }
}

#[derive(Hash, Eq, PartialEq, Clone, PartialOrd, Ord, Debug)]
pub enum TypeInstance {
    Any,
    Void,
    Bool,
    Char,
    Word,
    Nat,
    Int,
    Real,
    Vector(Box<TypeInstance>),
}

impl TypeInstance {
    pub fn is_atomic(&self) -> bool {
        match self {
            TypeInstance::Any | TypeInstance::Void | TypeInstance::Vector(_) => false,
            _ => true,
        }
    }

    pub fn as_str(&self) -> String {
        match self {
            TypeInstance::Any => "T",
            TypeInstance::Bool => "Bool",
            TypeInstance::Char => "Char",
            TypeInstance::Word => "Word",
            TypeInstance::Nat => "Nat",
            TypeInstance::Int => "Int",
            TypeInstance::Real => "Real",
            TypeInstance::Void => "Void",
            TypeInstance::Vector(v) => return format!("[{}]", v.as_str()),
        }
        .to_string()
    }
}
