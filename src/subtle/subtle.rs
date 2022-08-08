pub use crate::palisade;
pub use crate::prism;
pub use crate::subtle::*;

impl prism::Environment {
    pub fn regulate(&mut self) -> &prism::Environment {
        for (i, f) in &self.pre_functions {
            match f.alpha.clone() {
                Some(_) => {
                    self.dyadic_functions
                        .insert(i.clone(), self.induct_dyadic_function(f).unwrap());
                }
                None => {
                    self.monadic_functions
                        .insert(i.clone(), self.induct_monadic_function(f).unwrap());
                }
            };
        }

        for (i, f) in &self.pre_functions {
            if self.monadic_functions.contains_key(i) || self.monadic_functions.contains_key(i) {
                panic!();
            }

            match f.alpha.clone() {
                Some(_) => {
                    self.dyadic_functions
                        .insert(i.clone(), self.regulate_dyadic_function(f));
                }
                None => {
                    self.monadic_functions
                        .insert(i.clone(), self.regulate_monadic_function(f));
                }
            };
        }

        self
    }

    fn regulate_morpheme(&self, m: &palisade::Morpheme) -> prism::Morpheme {
        match m {
            palisade::Morpheme::Bool(b) => prism::Morpheme::Bool(*b),
            palisade::Morpheme::Char(c) => prism::Morpheme::Char(*c),
            palisade::Morpheme::Int(i) => prism::Morpheme::Int(*i),
            palisade::Morpheme::Real(r) => prism::Morpheme::Real(*r),
            palisade::Morpheme::Void => prism::Morpheme::Void,
        }
    }

    fn regulate_expression(&self, e: &palisade::Expression) -> prism::Expression {
        match e {
            palisade::Expression::Morpheme(m) => {
                prism::Expression::Morpheme(self.regulate_morpheme(m))
            }
            palisade::Expression::Application(a) => self.regulate_application(a),
            palisade::Expression::Vector(v) => prism::Expression::Vector(self.regulate_vector(v)),
        }
    }

    fn regulate_vector(&self, v: &palisade::Vector) -> prism::Vector {
        prism::Vector {
            element_type: v.kind().clone(),
            body: v.body.iter().map(|e| self.regulate_expression(e)).collect(),
        }
    }

    fn induct_monadic_function(&self, f: &palisade::Function) -> Option<prism::MonadicFunction> {
        let id = f.ident.clone();
        if self.monadic_functions.contains_key(&id) {
            return None;
        }

        Some(prism::MonadicFunction {
            ident: id,
            omega: f.omega.clone(),
            sigma: f.sigma.clone(),
            body: None,
            attrs: prism::FuncAttrs::new(),
        })
    }

    fn induct_dyadic_function(&self, f: &palisade::Function) -> Option<prism::DyadicFunction> {
        let id = f.ident.clone();
        if self.dyadic_functions.contains_key(&id) {
            return None;
        }

        Some(prism::DyadicFunction {
            ident: id,
            alpha: f.alpha.clone().unwrap(),
            omega: f.omega.clone(),
            sigma: f.sigma.clone(),
            body: None,
            attrs: prism::FuncAttrs::new(),
        })
    }

    fn regulate_monadic_function(&self, f: &palisade::Function) -> prism::MonadicFunction {
        prism::MonadicFunction {
            ident: f.ident.clone(),
            omega: f.omega.clone(),
            sigma: f.sigma.clone(),
            body: Some(
                f.body
                    .clone()
                    .unwrap()
                    .into_iter()
                    .map(|e| self.regulate_expression(&e))
                    .collect(),
            ),
            attrs: prism::FuncAttrs::new(),
        }
    }

    fn regulate_dyadic_function(&self, f: &palisade::Function) -> prism::DyadicFunction {
        prism::DyadicFunction {
            ident: f.ident.clone(),
            alpha: f.alpha.clone().unwrap(),
            omega: f.omega.clone(),
            sigma: f.sigma.clone(),
            body: Some(
                f.body
                    .clone()
                    .unwrap()
                    .into_iter()
                    .map(|e| self.regulate_expression(&e))
                    .collect(),
            ),
            attrs: prism::FuncAttrs::new(),
        }
    }

    fn regulate_application(&self, a: &palisade::Application) -> prism::Expression {
        match a.alpha {
            Some(_) => prism::Expression::Dyadic(self.regulate_dyadic_application(a)),
            None => prism::Expression::Monadic(self.regulate_monadic_application(a)),
        }
    }

    fn regulate_dyadic_application(&self, a: &palisade::Application) -> prism::DyadicApplication {
        // Check if function is in level 2 functions
        // If yes, use type information to continue

        let f = self.dyadic_functions.get(&a.app).unwrap();
        let lhs = self.regulate_expression(&a.alpha.as_ref().unwrap());
        let rhs = self.regulate_expression(&a.omega);
        let transformation = inspect_dapp_type(&f, &lhs, &rhs);

        let (_, active_t) = transformation;

        prism::DyadicApplication {
            alpha_t: f.alpha.clone(),
            alpha: Box::new(lhs),
            phi: f.ident.clone(),
            omega: Box::new(rhs),
            omega_t: active_t,
            sigma_t: f.sigma.clone(),
        }
    }

    fn regulate_monadic_application(&self, a: &palisade::Application) -> prism::MonadicApplication {
        // Check if function is in level 2 functions
        // If yes, use type information to continue

        let f = self.monadic_functions.get(&a.app).unwrap();
        let expr = self.regulate_expression(&a.omega);
        let transformation = inspect_mapp_type(&f, &expr);

        let (_, active_t) = transformation;

        prism::MonadicApplication {
            phi: f.ident.clone(),
            omega: Box::new(expr),
            omega_t: active_t,
            sigma_t: f.sigma.clone(),
        }
    }
}
