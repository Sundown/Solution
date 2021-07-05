<div align="center">
	<h3> Sunday :dove:</h3>

A minimalist, compiled, functional language.

Sunday generates code via LLVM.
</div>

---
#### Sample

```dart
@Package "Example";
@Entry Start;

Start : Void -> Void =
	Print Inverse 7;

Inverse : Real -> Real =
	/* Return 1 / %, where % is the input param */
	Return Quotient (1, %);
```

```sh
$ sunday run Example.xx
0.143
```
