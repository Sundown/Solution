use std::collections::HashSet;

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
}

pub struct MonadicFunction {
    pub ident: Ident,
    pub omega: Type,
    pub sigma: Type,
    pub body: Vec<Type>, // TODO Expression
}

#[derive(Eq, PartialEq)]
pub struct Type {
    set: HashSet<TypeInstance>,
    any: bool,
}

impl Type {
    pub fn any(&self) -> bool {
        self.any
    }

    pub fn none(&self) -> bool {
        self.set.len() > 0
    }

    pub fn of(t: TypeInstance) -> Type {
        let mut h = HashSet::new();
        h.insert(t);
        Type { set: h, any: false }
    }
}

#[derive(Hash, Eq, PartialEq)]
pub enum TypeInstance {
    Void,
    Bool,
    Char,
    Int,
    Real,
    Vector(Box<TypeInstance>), // TODO Type, once I figure out how to hash a HashSet
}

impl TypeInstance {
    pub fn is_atomic(&self) -> bool {
        // TODO probably a 1 liner for this
        match self {
            TypeInstance::Vector(_) => false,
            _ => true,
        }
    }
}
