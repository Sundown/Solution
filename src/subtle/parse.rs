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
    pub fn parse_unit(source: pest::iterators::Pairs<Rule>) -> Option<Error<String>> {
        for component in source {
            match component.as_rule() {
                Rule::function => ast.push(parse_expression(pair)),
                _ => {
                    println!("Not implemeneted or incorrectly placed");
                }
            }
        }

        None
    }
}

pub fn parse_unit(source: &str) -> Result<Vec<AstNode>, Error<Rule>> {
    let mut ast = vec![];

    let pairs = Palisade::parse(Rule::program, source)?;
    for pair in pairs {
        match pair.as_rule() {
            Rule::function => ast.push(parse_expression(pair)),
            _ => {
                println!("Not implemeneted");
            }
        }
    }

    Ok(ast)
}

fn parse_expression(pair: pest::iterators::Pair<Rule>) -> Box<dyn prism::Expression> {
    match pair.as_rule() {
        Rule::expr => parse_expression(pair.into_inner().next().unwrap()),
        Rule::monadicExpr => {
            let mut pair = pair.into_inner();
            let verb = pair.next().unwrap();
            let expr = pair.next().unwrap();
            let expr = parse_expression(expr);
            parse_monadic_verb(verb, expr)
        }
        Rule::function => {
            parse_function(pair.into_inner());
            AstNode::Integer(0)
        }
        Rule::dyadicExpr => {
            let mut pair = pair.into_inner();
            let lhspair = pair.next().unwrap();
            let lhs = parse_expression(lhspair);
            let verb = pair.next().unwrap();
            let rhspair = pair.next().unwrap();
            let rhs = parse_expression(rhspair);
            parse_dyadic_verb(verb, lhs, rhs)
        }
        Rule::morphemes => {
            let terms: Vec<AstNode> = pair.into_inner().map(parse_morpheme).collect();
            // If there's just a single term, return it without
            // wrapping it in a Terms node.
            match terms.len() {
                1 => terms.get(0).unwrap().clone(),
                _ => AstNode::Morphemes(terms),
            }
        }
        Rule::string => {
            let str = &pair.as_str();
            // Strip leading and ending quotes.
            let str = &str[1..str.len() - 1];
            // Escaped string quotes become single quotes here.
            let str = str.replace("''", "'");
            AstNode::Str(CString::new(&str[..]).unwrap())
        }
        unknown_expr => panic!("Unexpected expression: {:?}", unknown_expr),
    }
}

fn parse_dyadic_verb(pair: pest::iterators::Pair<Rule>, lhs: AstNode, rhs: AstNode) -> AstNode {
    AstNode::DyadicOp {
        lhs: Box::new(lhs),
        rhs: Box::new(rhs),
        verb: prism::Ident {
            package: "".to_string(),
            name: pair.to_string(),
        },
    }
}

fn parse_monadic_verb(pair: pest::iterators::Pair<Rule>, expr: AstNode) -> AstNode {
    AstNode::MonadicOp {
        verb: prism::Ident {
            package: "".to_string(),
            name: pair.to_string(),
        },
        expr: Box::new(expr),
    }
}

fn parse_morpheme(pair: pest::iterators::Pair<Rule>) -> AstNode {
    match pair.as_rule() {
        Rule::integer => {
            let istr = pair.as_str();
            let (sign, istr) = match &istr[..1] {
                "_" => (-1, &istr[1..]),
                _ => (1, &istr[..]),
            };
            let integer: i64 = istr.parse().unwrap();
            AstNode::Integer(sign * integer)
        }
        Rule::real => {
            let dstr = pair.as_str();
            let (sign, dstr) = match &dstr[..1] {
                "_" => (-1.0, &dstr[1..]),
                _ => (1.0, &dstr[..]),
            };
            let mut flt: f64 = dstr.parse().unwrap();
            if flt != 0.0 {
                // Avoid negative zeroes; only multiply sign by nonzeroes.
                flt *= sign;
            }
            AstNode::Real(flt)
        }

        Rule::expr => parse_expression(pair),
        Rule::ident => AstNode::Ident(String::from(pair.as_str())),
        unknown_term => panic!("Unexpected term: {:?}", unknown_term),
    }
}

pub fn parse_function(mut pair: pest::iterators::Pairs<Rule>) -> prism::Function {
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

    prism::Function {
        alpha: alpha_t,
        package: "".to_string(),
        name: ident_s.to_string(),
        omega: omega_t,
        sigma: sigma_t,
        body: None,
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
