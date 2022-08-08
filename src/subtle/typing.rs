pub use crate::palisade;
pub use crate::prism;
pub use crate::subtle::*;

pub enum Relationship {
    Match,
    FunctionFlexible,
    ExpressionFlexible,
}

pub fn integrate_type(
    model: &prism::TypeInstance,
    by: &prism::TypeInstance,
) -> prism::TypeInstance {
    if model.is_atomic() {
        return model.clone();
    };

    if model == by {
        return model.clone();
    };


    if matches!(model, prism::TypeInstance::Any) {
        return by.clone();
    };

    if let prism::TypeInstance::Vector(inner) = model {
        if let prism::TypeInstance::Vector(by_inner) = by {
            return prism::TypeInstance::Vector(Box::new(integrate_type(inner, by_inner)));
        } else {
            return prism::TypeInstance::Vector(Box::new(integrate_type(inner, by)));
        }
    };

    panic!()
}

pub fn inspect_mapp_type(
    f: &prism::MonadicFunction,
    e: &prism::Expression,
) -> (Relationship, prism::Type) {
    if f.omega == e.kind() {
        return (Relationship::Match, e.kind());
    }

    if f.omega.any() || f.omega.allows(&e.kind().single().unwrap()) {
        return (Relationship::FunctionFlexible, e.kind());
    }

    panic!("Casting not implemented yet")
}
