pub use crate::pest::Parser;
pub use crate::prism;
use crate::prism::Expression;
use core::panic;
pub use pest::error::Error;

use std::ffi::CString;

pub use crate::subtle::*;

#[derive(Parser)]
#[grammar = "grammar.pest"]
pub struct Palisade;

impl prism::Environment {
    pub fn parse_unit(&self, source: &str) -> Result<&prism::Environment, Error<String>> {
        let pairs = Palisade::parse(Rule::program, source).unwrap().into_iter();

        for pair in pairs {
            match pair.as_rule() {
                Rule::function => {
                    self.parse_expression(pair);
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

        panic!("Not implemented");
    }

    fn parse_expression(&self, pair: pest::iterators::Pair<Rule>) -> Box<dyn prism::Expression> {
        match pair.as_rule() {
            Rule::expr => self.parse_expression(pair.into_inner().next().unwrap()),
            Rule::function => self.parse_function(pair.into_inner()),
            Rule::monadicExpr => self.parse_monadic_app(pair.into_inner()),
            Rule::dyadicExpr => self.parse_dyadic_app(pair.into_inner()),
            Rule::morphemes => {
                let terms: Vec<Box<dyn prism::Expression>> =
                    pair.into_inner().map(|x| self.parse_morpheme(x)).collect();

                match terms.len() {
                    //1 => Box::new(terms.get(0).unwrap()),
                    _ => Box::new(prism::Vector { body: terms }),
                }
            }
            unknown_expr => panic!("Unexpected expression: {:?}", unknown_expr),
        }
    }

    fn parse_dyadic_app(
        &self,
        mut pair: pest::iterators::Pairs<Rule>,
    ) -> Box<dyn prism::Expression> {
        let lhspair = pair.next().unwrap();
        let verb = pair.next().unwrap();
        let rhspair = pair.next().unwrap();

        Box::new(prism::Application {
            alpha: Some(self.parse_expression(lhspair)),
            app: prism::Ident {
                // TODO check this is valid type
                package: "".to_string(),
                name: verb.into_inner().as_str().to_string(),
            },
            omega: self.parse_expression(rhspair),
        })
    }

    fn parse_monadic_app(
        &self,
        mut pair: pest::iterators::Pairs<Rule>,
    ) -> Box<dyn prism::Expression> {
        let verb = pair.next().unwrap();
        let rhspair = pair.next().unwrap();

        Box::new(prism::Application {
            alpha: None,
            app: prism::Ident {
                // TODO check this is valid type
                package: "".to_string(),
                name: verb.into_inner().as_str().to_string(),
            },
            omega: self.parse_expression(rhspair),
        })
    }

    fn parse_morpheme(&self, pair: pest::iterators::Pair<Rule>) -> Box<dyn prism::Expression> {
        match pair.as_rule() {
            Rule::integer => {
                let istr = pair.as_str();
                let (sign, istr) = match &istr[..1] {
                    "_" => (-1, &istr[1..]),
                    _ => (1, &istr[..]),
                };
                let integer: i64 = istr.parse().unwrap();
                Box::new(prism::Morpheme::Int(sign * integer))
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

                Box::new(prism::Morpheme::Real(flt))
            }
            unknown_term => panic!("Unexpected term: {:?}", unknown_term),
        }
    }

    pub fn parse_function(
        &self,
        mut pair: pest::iterators::Pairs<Rule>,
    ) -> Box<dyn prism::Expression> {
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
                        _ => None,
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

        let f = Box::new(prism::Function {
            alpha: alpha_t,
            package: "".to_string(),
            name: ident_s.to_string(),
            omega: omega_t,
            sigma: sigma_t,
            body: Some(
                pair.next()
                    .unwrap()
                    .into_inner()
                    .into_iter()
                    .map(|e| self.parse_expression(e))
                    .collect::<Vec<_>>(),
            ),
        });
        println!("{}", f.as_str());
        f
    }
}

#[derive(PartialEq, Debug, Clone)]
pub enum AstNode {
    Integer(i64),
    Real(f64),
    Operator {
        verb: prism::Ident,
        op: Box<AstNode>,
    },
    Block {
        body: Vec<AstNode>,
    },
    Applicable {
        expr: Box<AstNode>,
        ident: String,
        operator: Box<AstNode>,
    },
    MonadicOp {
        verb: prism::Ident,
        expr: Box<AstNode>,
    },
    DyadicOp {
        verb: prism::Ident,
        lhs: Box<AstNode>,
        rhs: Box<AstNode>,
    },
    Morphemes(Vec<AstNode>),
    IsGlobal {
        ident: String,
        expr: Box<AstNode>,
    },
    Ident(String),
    Str(CString),
}
