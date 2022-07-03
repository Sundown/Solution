pub use crate::pest::Parser;
use core::panic;
pub use pest::error::Error;

use std::ffi::CString;

pub use crate::subtle::*;

#[derive(Parser)]
#[grammar = "grammar.pest"]
pub struct Palisade;

impl Environment {
    pub fn parse_unit(&mut self, source: &str) -> Result<&Environment, Error<String>> {
        let pairs = Palisade::parse(Rule::program, source).unwrap().into_iter();

        for pair in pairs {
            match pair.as_rule() {
                Rule::function => {
                    self.parse_function_head(pair.into_inner());
                }
                Rule::EOI => {
                    return Ok(self);
                }
                _ => {
                    println!("{:?}", pair);
                    panic!("Not implemeneted");
                }
            }
        }

        Ok(self)
    }

    fn parse_expression(&self, pair: pest::iterators::Pair<Rule>) -> Box<Expression> {
        match pair.as_rule() {
            Rule::expr => self.parse_expression(pair.into_inner().next().unwrap()),
            Rule::monadicExpr => self.parse_monadic_app(pair.into_inner()),
            Rule::dyadicExpr => self.parse_dyadic_app(pair.into_inner()),
            Rule::morphemes => {
                let terms: Vec<Expression> = pair
                    .into_inner()
                    .map(|x| self.parse_morpheme(x).expr())
                    .collect();

                match terms.len() {
                    //1 => Box::new(terms.get(0).unwrap()),
                    _ => Box::new(Vector { body: terms }.expr()),
                }
            }
            unknown_expr => panic!("Unexpected expression: {:?}", unknown_expr),
        }
    }

    fn parse_dyadic_app(&self, mut pair: pest::iterators::Pairs<Rule>) -> Box<Expression> {
        let lhspair = pair.next().unwrap();
        let verb = pair.next().unwrap();
        let rhspair = pair.next().unwrap();

        Box::new(
            Application {
                alpha: Some(self.parse_expression(lhspair)),
                app: Ident {
                    package: "".to_string(),
                    name: verb.into_inner().as_str().to_string(),
                },
                omega: self.parse_expression(rhspair),
            }
            .expr(),
        )
    }

    fn parse_monadic_app(&self, mut pair: pest::iterators::Pairs<Rule>) -> Box<Expression> {
        let verb = pair.next().unwrap();
        let rhspair = pair.next().unwrap();

        Box::new(
            Application {
                alpha: None,
                app: Ident {
                    package: "".to_string(),
                    name: verb.into_inner().as_str().to_string(),
                },
                omega: self.parse_expression(rhspair),
            }
            .expr(),
        )
    }

    fn parse_morpheme(&self, pair: pest::iterators::Pair<Rule>) -> Box<Expression> {
        match pair.as_rule() {
            Rule::integer => {
                let istr = pair.as_str();
                let (sign, istr) = match &istr[..1] {
                    "_" => (-1, &istr[1..]),
                    _ => (1, &istr[..]),
                };
                let integer: i64 = istr.parse().unwrap();
                Box::new(Morpheme::Int(sign * integer).expr())
            }
            Rule::real => {
                let dstr = pair.as_str();
                let (sign, dstr) = match &dstr[..1] {
                    "_" => (-1.0, &dstr[1..]),
                    _ => (1.0, &dstr[..]),
                };
                let mut flt: f64 = dstr.parse().unwrap();
                if flt != 0.0 {
                    // Avoid negative zeroes
                    flt *= sign;
                }

                Box::new(Morpheme::Real(flt).expr())
            }
            Rule::string => {
                let sstr = pair.as_str();

                Box::new(
                    Vector {
                        body: (&sstr[1..sstr.len() - 1])
                            .to_string()
                            .into_bytes()
                            .into_iter()
                            .map(|c| Morpheme::Char(c).expr())
                            .collect(),
                    }
                    .expr(),
                )
            }
            Rule::rune => Box::new(Morpheme::Char(pair.as_str().as_bytes()[1]).expr()),
            Rule::boolean => Box::new(Morpheme::Bool(pair.as_str().parse().unwrap()).expr()),
            unknown_term => panic!("Unexpected term: {:?}", unknown_term),
        }
    }

    pub fn parse_function_head(&mut self, mut pair: pest::iterators::Pairs<Rule>) -> Ident {
        let head = pair.next().unwrap();

        let (alpha_t, ident_s, omega_t, sigma_t) = match head.as_rule() {
            Rule::typedFunctionHead => {
                let mut head = head.into_inner();
                let front = head.next().unwrap();
                let rule = front.as_rule();
                let mut front = front.into_inner();
                (
                    match rule {
                        Rule::typedDyadic => Some(parse_type(front.next().unwrap())),
                        _ => None, // monadic, no alpha type to parse
                    },
                    front.next().unwrap().as_str(),
                    parse_type(front.next().unwrap()),
                    parse_type(head.next().unwrap()),
                )
            }

            Rule::ambiguousFunctionHead => {
                panic!("TODO")
            }
            _ => panic!("Unexpected function head: {:?}", head.as_rule()),
        };
        let id = Ident::new("", ident_s);
        let f = Function {
            alpha: Some(TypeGroup::of(&alpha_t.clone().unwrap())),
            ident: id.clone(),
            omega: TypeGroup::of(&omega_t),
            sigma: TypeGroup::of(&sigma_t),
            body: Some(
                pair.next()
                    .unwrap()
                    .into_inner()
                    .into_iter()
                    .map(|e| self.parse_expression(e))
                    .collect::<Vec<_>>(),
            ),
        };

        match alpha_t {
            Some(_) => &self.dya_fns.insert(id.clone(), f),
            None => &self.mon_fns.insert(id.clone(), f),
        };

        id.clone()
    }
}
