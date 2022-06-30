pub use crate::prism::*;
use std::collections::HashMap;

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
    pub emit_format: Format,
    pub emit_name: String,
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

impl Environment {
    pub fn get_function(&self, package: &str, name: &str) -> Option<&Function> {
        let package = package.to_string();
        let name = name.to_string();

        if let Some(fn_) = self
            .mon_fns
            .get(&base_ident(&format!("{}.{}", package, name)))
        {
            return Some(fn_);
        }

        if let Some(fn_) = self
            .dya_fns
            .get(&base_ident(&format!("{}.{}", package, name)))
        {
            return Some(fn_);
        }

        None
    }
}

pub enum Opt {
    // Pure,
    // None,
    // Size,
    Fast,
}

pub enum Format {
    //IR,
    LLVM(Opt),
    //Assembly,
}
