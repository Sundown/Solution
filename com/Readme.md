### Source

#### Layout

Solution compilation is broken into multiple, affectionately named stages.

1. **Palisade** is the stage in which lexical analysis is performed.
2. **Prescience** is is where function declaration is conducted, making the following stages easier due to foreknowledge of appropriate contexts for function names.
3. **Weave** transforms an array of tokens in the function bodies into a tree, no contextual analysis is performed.
4. **Subtle** ensures primatives are used correctly, adds detailed type information, transforms syntax tree into a form ready for optimisation and evaluation.
5. **Monia** performs optimisations, notifies next stage of possible shortcuts, and pre-computes constant expressions.
6. **Apotheosis** performs the now-simple task of code generation.

The first and final stages are handled by **Oversight** which is responsible for input handling and binary emission.
