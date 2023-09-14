<!-- @format -->

<h2 align="center"> Solution</h2>
<p align="center">
Solution is a compiler for an array-oriented language, providing the cognition of APL in an accessible, compiled, and open-source platform. The Solution Language is inspired by the work of Kenneth Iverson.
</p>

<p align="center">
  <a href="https://github.com/Sundown/Solution/blob/master/go.mod">
		<img alt="Go Version" src="https://img.shields.io/github/go-mod/go-version/sundown/solution?style=for-the-badge&logo=go&color=f1f1f1&logoColor=f1f1f1&labelColor=262D3A">
  </a>
  <a href="https://github.com/sundown/solution/blob/main/LICENSE">
    <img src="https://img.shields.io/static/v1.svg?style=for-the-badge&logo=gnu&label=License&message=GPL-2.0&color=f1f1f1&logoColor=f1f1f1&labelColor=262D3A"/>
  </a>
  <a href="https://llvm.org">
    <img src="https://img.shields.io/static/v1.svg?style=for-the-badge&logo=llvm&label=LLVM&message=v13.0&color=f1f1f1&logoColor=f1f1f1&labelColor=262D3A"/>
  </a>

</p>

The following demonstrates Solution's implicit typing system, `Demo` is automatically typed as `Int×Int` and takes the minimum of the two arguments to the power of the maximum using a dyadic train.

```swift
@Package dev;

Main Int → Void {
	8 Demo 2;
	3.5 dev::Demo 2.0;
}

Demo → Void {
	Println α (⌊*⌈) ω;
}
```
