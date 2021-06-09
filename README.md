<div align="center">
	<h3> Girl :dove:</h3>


A compiled language which improves on the elegance of LISP.

Girl generates code via LLVM, hopefully at some point this will be abstracted with a homemade backend.
</div>

---
#### Sample

```haskell
inverse = real -> real:
	return / (1, %)
;
```

:flushed: what's happening!?

Well...

Functions take a single argument, no not like Haskell, this is not currying.

If you, hypothetically, wanted to print "hello world" to stdout, you might write
```haskell
print "hello world"
```

But that's just Python 2! You can find out more about Girl's similarities with Python 2 by pressing <kbd>Ctrl</kbd> + <kbd>W</kbd>.

If you, again hypothetically, wanted to print the square root of 9:
```haskell
print sqrt 9
```
This might be written as
```c
print(sqrt(9))
```
in a C-like language.

To pass 'multiple' arguments, use a tuple or a vector!
```haskell
open ("example.txt", "r")

-- or

sum [2, 3, 5, 7, 11]

```

---
#### Future Goals
- [ ] Handwritten parser
- [ ] C interop.
- [ ] A proper standard library
- [ ] Implement a seperate backend with a funny name
