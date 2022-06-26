pub use crate::pest::Parser;

pub use crate::prism;
use core::panic;
pub use pest::error::Error;
use std::ffi::CString;

pub use crate::subtle::*;

#[derive(Parser)]
#[grammar = "grammar.pest"]
pub struct Palisade;

impl prism::Environment {
    pub fn parse_unit(&self, source: &str) -> Option<Error<String>> {
        let pairs = Palisade::parse(Rule::program, source).unwrap().into_iter();

        for pair in pairs {
            match pair.as_rule() {
                Rule::function => {
                    self.parse_expression(pair);
                }
                _ => {
                    panic!("Not implemeneted");
                }
            }
        }

        None
    }

    fn parse_expression(&self, pair: pest::iterators::Pair<Rule>) -> Box<dyn prism::Expression> {
        match pair.as_rule() {
            Rule::expr => self.parse_expression(pair.into_inner().next().unwrap()),
            Rule::monadicExpr => self.parse_monadic_verb(pair.into_inner()),
            Rule::function => self.parse_function(pair.into_inner()),
            Rule::dyadicExpr => {
                let mut pair = pair.into_inner();
                self.parse_dyadic_verb(pair)
            }
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

    fn parse_dyadic_verb(&self, pair: pest::iterators::Pairs<Rule>) -> Box<dyn prism::Expression> {
        let lhspair = pair.next().unwrap();
        let verb = pair.next().unwrap();
        let rhspair = pair.next().unwrap();

        Box::new(prism::Application {
            alpha: Some(self.parse_expression(lhspair)),
            app: Box::new(self.get_function("".to_string(), verb.as_str()).unwrap()),
            omega: self.parse_expression(rhspair),
        })
    }

    fn parse_monadic_verb(&self, pair: pest::iterators::Pairs<Rule>) -> Box<dyn prism::Expression> {
        let verb = pair.next().unwrap();
        let rhspair = pair.next().unwrap();

        Box::new(prism::Application {
            alpha: None,
            app: prism::Ident {
                package: "".to_string(),
                name: pair.to_string(),
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
                match front.as_rule() {
                    Rule::typedDyadic => {
                        let mut front = front.into_inner();
                        (
                            Some(parse_type(front.next().unwrap())),
                            front.next().unwrap().as_str(),
                            parse_type(front.next().unwrap()),
                            parse_type(head.next().unwrap()),
                        )
                    }
                    Rule::typedMonadic => {
                        let mut front = front.into_inner();
                        (
                            None,
                            front.next().unwrap().as_str(),
                            parse_type(front.next().unwrap()),
                            parse_type(head.next().unwrap()),
                        )
                    }
                    _ => {
                        panic!("Unexpected typed function head: {:?}", front.as_rule())
                    }
                }
            }
            Rule::ambiguousFunctionHead => {
                panic!("TODO")
            }
            _ => panic!("Unexpected function head: {:?}", head.as_rule()),
        };

        //let body = pair.next().unwrap().into_inner();

        Box::new(prism::Function {
            alpha: alpha_t,
            package: "".to_string(),
            name: ident_s.to_string(),
            omega: omega_t,
            sigma: sigma_t,
            body: None,
        })
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
