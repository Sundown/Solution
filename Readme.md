### Solution

[![Go Report Card](https://goreportcard.com/badge/github.com/sundown/solution)](https://goreportcard.com/report/github.com/sundown/solution)

Array language compiler with a nice typesystem.

Solution files are suffixed with `.sol`.

Some notes about ongoing implementation:

- Code generation is done via LLVM, this is done using [llir/llvm](https://github.com/llir/llvm), which will be replaced with LLVM's Go bindings eventually.
- A (non-JIT) interpreter will be added, which may supplement some optimisations.
- Solution will invoke Clang on a temporary LLVMIR file when ready.
- At present, some programs will be faster written in Solution than in C, if Clang is invoked with default options.
- Basic optimisations such as `leal` combinations for some multiplications, as well as magic numbers for division (once I learn how to have them calculated).
- More complicated optimisations such as loop-jamming when the `Map` operator is invoked successively will be implemented once the interpreter exists.
- Type system is made by combination of various types, including but not limited to: `Int`, `Bool`, `Char`... these may be extended to vectors turning them into `String`, much like C.
- Type system includes sum types (for example a sum type of `Numeric` could be `Int | Real`)
- A single generic type T exists (alternate options in future), the type of which is substituted into other types in the function's signaturse, as well as any generic-accepting functions called within the function body.
- Trains (forks and atops) are available in the Dyalog style:

	`a (fgh) b <=> (a f b) g (a h b)`

	`  (fgh) b <=>   (f b) g (h b)`

	`a (gh) b <=> g (a h b)`

	`  (gh) b <=> g (h b)`
