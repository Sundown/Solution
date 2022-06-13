extern crate pest;
#[macro_use]
extern crate pest_derive;

use core::panic;
use pest::error::Error;
use pest::Parser;
use std::ffi::CString;

#[derive(Parser)]
#[grammar = "grammar.pest"]
pub struct Palisade;

#[derive(Debug, Clone, PartialEq)]
pub struct Ident {
    pub package: String,
    pub name: String,
}

#[derive(Debug, Clone, PartialEq)]
pub enum AtomicType {
    Bool,
    Char,
    Int,
    Real,
    Void,
}

#[derive(Debug, Clone, PartialEq)]
pub struct Function {
    pub package: String,
    pub name: String,
    pub alpha: Option<Type>,
    pub omega: Type,
    pub sigma: Type,
    pub body: Box<AstNode>,
}

#[derive(Debug, Clone, PartialEq)]
pub enum Type {
    Atomic(AtomicType),
    Vector(Box<Type>),
}

#[derive(PartialEq, Debug, Clone)]
pub enum AstNode {
    Integer(i64),
    Real(f64),
    Operator {
        verb: Ident,
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
        verb: Ident,
        expr: Box<AstNode>,
    },
    DyadicOp {
        verb: Ident,
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

pub fn parse(source: &str) -> Result<Vec<AstNode>, Error<Rule>> {
    let mut ast = vec![];

    let pairs = Palisade::parse(Rule::program, source)?;
    for pair in pairs {
        match pair.as_rule() {
            Rule::function => ast.push(build_ast_from_expr(pair)),
            _ => {
                println!("FA");
            }
        }
    }

    Ok(ast)
}

fn build_ast_from_expr(pair: pest::iterators::Pair<Rule>) -> AstNode {
    match pair.as_rule() {
        Rule::expr => build_ast_from_expr(pair.into_inner().next().unwrap()),
        Rule::monadicExpr => {
            let mut pair = pair.into_inner();
            let verb = pair.next().unwrap();
            let expr = pair.next().unwrap();
            let expr = build_ast_from_expr(expr);
            parse_monadic_verb(verb, expr)
        }
        Rule::function => parse_function(pair.into_inner()),
        Rule::dyadicExpr => {
            let mut pair = pair.into_inner();
            let lhspair = pair.next().unwrap();
            let lhs = build_ast_from_expr(lhspair);
            let verb = pair.next().unwrap();
            let rhspair = pair.next().unwrap();
            let rhs = build_ast_from_expr(rhspair);
            parse_dyadic_verb(verb, lhs, rhs)
        }
        Rule::morphemes => {
            let terms: Vec<AstNode> = pair.into_inner().map(build_ast_from_term).collect();
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
        verb: Ident {
            package: "".to_string(),
            name: pair.to_string(),
        },
    }
}

fn parse_type(t: pest::iterators::Pair<Rule>) -> Type {
    match t.as_rule() {
        Rule::typeActual => parse_type(t.into_inner().next().unwrap()),
        Rule::atomicType => Type::Atomic(match t.as_str() {
            "Bool" => AtomicType::Bool,
            "Char" => AtomicType::Char,
            "Int" => AtomicType::Int,
            "Real" => AtomicType::Real,
            "Void" => AtomicType::Void,
            _ => panic!("Unexpected type: {:?}", t.as_str()),
        }),
        Rule::vectorType => Type::Vector(Box::new(parse_type(t.into_inner().next().unwrap()))),
        _ => panic!("Unexpected type: {:?} as {:?}", t.as_str(), t.as_rule()),
    }
}

fn parse_function(mut pair: pest::iterators::Pairs<Rule>) -> Function {
    let head = pair.next().unwrap();

    let (alpha_t, ident_s, omega_t) = match head.as_rule() {
        Rule::typedFunctionHead => {
            let front = head.into_inner().next().unwrap();
            match front.as_rule() {
                Rule::typedDyadic => {
                    let mut front = front.into_inner();
                    (
                        Some(parse_type(front.next().unwrap())),
                        front.next().unwrap().as_str(),
                        parse_type(front.next().unwrap()),
                    )
                }
                Rule::typedMonadic => {
                    let mut front = front.into_inner();
                    (
                        None,
                        front.next().unwrap().as_str(),
                        parse_type(front.next().unwrap()),
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

    Function {
        alpha: alpha_t,
        package: "".to_string(),
        name: ident_s.to_string(),
        omega: omega_t,
        body: Box::new(build_ast_from_expr(body)),
    }
}

fn parse_monadic_verb(pair: pest::iterators::Pair<Rule>, expr: AstNode) -> AstNode {
    AstNode::MonadicOp {
        verb: Ident {
            package: "".to_string(),
            name: pair.to_string(),
        },
        expr: Box::new(expr),
    }
}

fn build_ast_from_term(pair: pest::iterators::Pair<Rule>) -> AstNode {
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

        Rule::expr => build_ast_from_expr(pair),
        Rule::ident => AstNode::Ident(String::from(pair.as_str())),
        unknown_term => panic!("Unexpected term: {:?}", unknown_term),
    }
}

fn main() {
    let unparsed_file = std::fs::read_to_string("code.sol").expect("Cannot read file");
    let astnode = parse(&unparsed_file).expect("Unsuccessful parse");
    println!("{:?}", &astnode);
}
