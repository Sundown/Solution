use std::collections::HashMap;

pub use crate::prism::*;

use super::base_ident;

pub struct Environment {
    // Prism
    pub iter: u32,
    pub verbose: bool,

    // Subtle
    pub types: HashMap<Ident, Type>,
    pub mon_fns: HashMap<Ident, Function>,
    pub dya_fns: HashMap<Ident, Function>,
    pub current_fn: Option<Function>,

    // Apotheosis
    // TOOD

    // Emit
    emit_format: Format,
    emit_name: String,
}

pub fn new_environment() -> Environment {
    Environment {
        iter: 0,
        verbose: false,
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
        current_fn: None,
        emit_name: "main".to_string(),
        emit_format: Format::LLVM(Opt::Fast),
    }
}

enum Opt {
    Pure,
    None,
    Size,
    Fast,
}

enum Format {
    IR,
    LLVM(Opt),
    Assembly,
}
