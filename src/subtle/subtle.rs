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
            // TODO make this check function isn't being redefined differently
            // if self.monadic_functions.contains_key(i) || self.monadic_functions.contains_key(i) {
            //     continue;
            // }

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
            palisade::Expression::Application(_) => panic!("TODO app"),
            palisade::Expression::Vector(_) => panic!("TODO vec"),
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

    fn regulate_monadic_application(
        &mut self,
        a: &palisade::Application,
    ) -> prism::MonadicApplication {
        // Check if function is in level 2 functions
        // If yes, use type information to continue

        let f = match self.monadic_functions.get(&a.app) {
            Some(f) => f,
            None => {
                let n = self.pre_functions.get(&a.app).unwrap();
                let n = self.regulate_monadic_function(n);

                self.monadic_functions.insert(n.ident.clone(), n);
                self.monadic_functions.get(&a.app).unwrap() // &n
            }
        };

        panic!();
    }
}
