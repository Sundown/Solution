## Sunday :dove:</h3>

:construction: WIP

A minimalist, compiled, functional language.

This project is currently a work in progress and is useful for little apart from learning the basics of LLVM as used in a small-scale project.

---
#### Sample

```java
@Package "Example";
@Entry Start;

Start : Void -> Void =
	Print Inverse 7;

Inverse : Real -> Real =
	/* Return 1 / %, where % is the input param */
	Return Quotient (1, %);
```

```sh
$ sndy run Example.xx
0.143
```
---
#### Running
Simply `go build sndy.go` and everything should work out, provided deps are available.
