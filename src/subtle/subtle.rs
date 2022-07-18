pub use crate::prism;
pub use crate::subtle::*;

pub fn regulate(env: &prism::Environment) -> &prism::Environment {
    for f in &env.monadic_functions {
        f;
    }

    env
}

impl prism::Environment {
	fn regulate_monadic_function(&self, )
}
