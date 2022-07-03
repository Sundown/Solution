pub use crate::subtle::*;
use std::collections::HashMap;

pub struct Environment {
    pub types: HashMap<Ident, Type>,
    pub mon_fns: HashMap<Ident, Function>,
    pub dya_fns: HashMap<Ident, Function>,
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
        mon_fns: HashMap::new(),
        dya_fns: HashMap::new(),
    }
}

impl Environment {
    pub fn get_function(&self, id: &Ident) -> Option<&Function> {
        if let Some(fn_) = self.mon_fns.get(&id) {
            return Some(fn_);
        }

        if let Some(fn_) = self.dya_fns.get(&id) {
            return Some(fn_);
        }

        None
    }
}
