## Compiler :dove:</h3>

A minimalist, compiled, functional language.

---
#### Sample

```java
@Package Example
@Entry Start

Start : Void -> Void = Print ["1/7 = ", String Inverse 7];

Inverse : Real -> Real = Return Quotient (1, @); // @ refers to param
```

```sh
$ compiler run Example.xx
1/7 = 0.143
```
---
#### Running
`go build main.go` provided dependancies are installed.

---

#### TODO
- [ ] Named tuple fields
- [ ] Algebraic types
- [ ] Parametric polymorphism
- [ ] Boxed data/explicit heap allocation
- [ ] CFFI
- [ ] Multiple files with separate namespaces

far future

- [ ] Closures, exceptions, coroutines
