pub use crate::prism::*;
use std::collections::{HashMap, HashSet};

// Needs Environment to be singleton
pub struct Environment {
    pub base_types: HashMap<Ident, TypeInstance>,
    pub monadic_functions: HashMap<Ident, MonadicFunction>,
    pub dyadic_functions: HashMap<Ident, DyadicFunction>,
}

// Needs Environment to be singleton
impl DyadicApplication {
    pub fn kind(&self, env: &Environment) -> Type {
        env.dyadic_functions.get(&self.phi).unwrap().kind()
    }
}
// Needs Environment to be singleton
impl MonadicApplication {
    pub fn kind(&self, env: &Environment) -> Type {
        env.monadic_functions.get(&self.phi).unwrap().kind()
    }
}

pub enum Expression {
    Morpheme(Morpheme),
    Monadic(MonadicApplication),
    Dyadic(DyadicApplication),
}

#[derive(Hash, PartialEq, Eq)]
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

pub struct MonadicApplication {
    pub phi: Ident,
    pub omega: Box<Expression>,
}

pub struct DyadicApplication {
    pub alpha: Box<Expression>,
    pub phi: Ident,
    pub omega: Box<Expression>,
}

pub struct MonadicFunction {
    pub ident: Ident,
    pub omega: Type,
    pub sigma: Type,
    pub body: Vec<Expression>,
    pub attrs: FuncAttrs,
}

impl MonadicFunction {
    pub fn kind(&self) -> Type {
        self.sigma.clone()
    }
}

impl DyadicFunction {
    pub fn kind(&self) -> Type {
        self.sigma.clone()
    }
}

pub struct DyadicFunction {
    pub ident: Ident,
    pub alpha: Type,
    pub omega: Type,
    pub sigma: Type,
    pub body: Vec<Expression>,
    pub attrs: FuncAttrs,
}

pub struct FuncAttrs {
    pub inline: bool, // For LLVM func attr
    pub elide: bool,  // Subtle stage should ignore
}

#[derive(Eq, PartialEq, Clone)]
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

#[derive(Hash, Eq, PartialEq, Clone)]
pub enum TypeInstance {
    Void,
    Bool,
    Char,
    Nat,
    Int,
    Real,
    Vector(Box<TypeInstance>), // TODO Type, once I figure out how to hash a HashSet
}

impl TypeInstance {
    pub fn is_atomic(&self) -> bool {
        !matches!(self, TypeInstance::Vector(_))
    }
}
