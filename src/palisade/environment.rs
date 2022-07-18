pub use crate::palisade::*;
use std::collections::HashMap;

pub struct Environment {
    pub types: HashMap<Ident, Type>,
    pub functions: HashMap<Ident, Function>,
}

pub fn new_environment() -> Environment {
    Environment {
        types: {
            let mut types = HashMap::new();
            types.insert(base_ident("Int"), Type::Atomic(AtomicType::Int));
            types.insert(base_ident("Char"), Type::Atomic(AtomicType::Char));
            types.insert(base_ident("Real"), Type::Atomic(AtomicType::Real));
            types.insert(base_ident("Bool"), Type::Atomic(AtomicType::Bool));
            types.insert(base_ident("Void"), Type::Atomic(AtomicType::Void));

            types
        },

        functions: HashMap::new(),
    }
}

impl Environment {
    pub fn get_function(&self, id: &Ident) -> Option<&Function> {
        self.functions.get(&id)
    }
}
