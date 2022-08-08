pub use crate::pest::Parser;
use core::panic;
pub use pest::error::Error;

pub use crate::palisade::*;
pub use crate::prism;

#[derive(Parser)]
#[grammar = "grammar.pest"]
pub struct Palisade;

impl prism::Environment {
    pub fn parse_unit(&mut self, source: &str) -> Result<&prism::Environment, Error<String>> {
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

    fn parse_expression(&mut self, pair: pest::iterators::Pair<Rule>) -> Box<Expression> {
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

    fn parse_dyadic_app(&mut self, mut pair: pest::iterators::Pairs<Rule>) -> Box<Expression> {
        let lhspair = pair.next().unwrap();
        let verb = pair.next().unwrap();
        let rhspair = pair.next().unwrap();

        Box::new(
            Application {
                alpha: Some(self.parse_expression(lhspair)),
                app: self.parse_ident(verb.into_inner().next().unwrap()),
                omega: self.parse_expression(rhspair),
            }
            .expr(),
        )
    }

    fn parse_monadic_app(&mut self, mut pair: pest::iterators::Pairs<Rule>) -> Box<Expression> {
        let verb = pair.next().unwrap();
        let rhspair = pair.next().unwrap();

        Box::new(
            Application {
                alpha: None,
                app: self.parse_ident(verb),
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

    pub fn parse_function_head(&mut self, mut pair: pest::iterators::Pairs<Rule>) {
        let head = pair.next().unwrap();

        let (alpha_t, ident_s, omega_t, sigma_t) = match head.as_rule() {
            Rule::typedFunctionHead => {
                let mut head = head.into_inner();
                let front = head.next().unwrap();
                let rule = front.as_rule();
                let mut front = front.into_inner();
                (
                    match rule {
                        Rule::typedDyadic => {
                            Some(prism::Type::of(parse_type(front.next().unwrap())))
                        }
                        _ => None, // monadic, no alpha type to parse
                    },
                    front.next().unwrap(),
                    prism::Type::of(parse_type(front.next().unwrap())),
                    prism::Type::of(parse_type(head.next().unwrap())),
                )
            }

            Rule::ambiguousFunctionHead => {
                let mut head = head.into_inner();
                (
                    Some(prism::Type::new_any()),
                    head.next().unwrap(),
                    prism::Type::new_any(),
                    match head.next() {
                        Some(t) => prism::Type::of(parse_type(t)),
                        _ => prism::Type::new_any(),
                    },
                )
            }
            _ => panic!("Unexpected function head: {:?}", head.as_rule()),
        };

        let id = self.parse_ident(ident_s);
        let f = Function {
            alpha: alpha_t,
            ident: id.clone(),
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
        };

        self.pre_functions.insert(id, f);
    }

    fn parse_ident(&mut self, pair: pest::iterators::Pair<Rule>) -> prism::Ident {
        let p = pair.clone().into_inner().next().unwrap();
        let mut q = p.clone().into_inner();

        // Normalise identifiers from user
        // "x" -> ( , x) -> ::x
        // "y::x" -> (y, x) -> y::x
        let (package, name) = match &q.clone().count() {
            0 => ("", p.as_str()),
            _ => (q.next().unwrap().as_str(), q.next().unwrap().as_str()),
        };

        prism::Ident::new(package, name)
    }
}
