<!-- @format -->

<h2 align="center"> Solution</h2>
<p align="center">
<i>An array-oriented language with a modular backend</i>
</p>

A compiled array-oriented language running on LLVM with a modular backend to support development for multiple platforms. This compiler supports optional implicit typing so it remains safe while not getting in your way.

This project relies heavily on:

- [llir/llvm](github.com/llir/llvm) which provides LLVM generation for Go
- [participle](github.com/alecthomas/participle) which is used to define the parser
- [clang](https://clang.llvm.org) for the heavy lifting as linker and assembler

---

### Demo

The following demonstrates Solution's implicit typing system, `Demo` is automatically typed as `Int×Int` and takes the minimum of the two arguments to the power of the maximum using a dyadic train.

The second function calculates the average of a numeric vector using a similar 3-arity dyadic train, however the leftmost side of this train is a function-operator combination, `+/` (add-reduce) which is equivalent to a sum.

```swift
@Package dev;

Main Int → Void {
	8 Demo 2;
	dev::Avg 1 3 99;
}

Demo → Void {
	Println α (⌊*⌈) ω;
}

Avg → Void {
	Println (+/÷≢) ω;
}

```

---

### Usage

```sh
go run solution.go test

go build solution.go

./solution input.sol -emit purellvm

clang libsol.c -S -emit-llvm -O0 -o libsol.ll

clang input.ll libsol.ll -Og -o bin

./bin
```

---

<p align="center">GPL2 </p>
