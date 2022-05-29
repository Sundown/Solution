<h2 align="center"> Solution</h2>
<p align="center">
Solution is a compiler for an array-oriented language, providing the cognition of APL in an accessible, compiled, and open-source platform. The Solution Language is inspired by the work of Kenneth Iverson.
</p>

<h3 align="center">Progress is being made on a Rust rewrite, on a separate branch.</h3>

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


---
### 🔍 Why?

- **Array oriented**: what do our computers process the most? A trivial but accurate answer is *data*, but what is data? Data is the plural form of the Latin word *datum*, which means *a piece of information*. This answer is a joke, yet it is useful; lending the insight that most processing is not performed on a single datum, but rather groups of these, again: data. Despite this striking fact, most languages prefer to focus on a single datum at a time, does this not seem a little strange, considering how much Computer Science is focused on graphs, sets, and other types of groups? Array oriented languages look at the bigger picture, preferring to manipulate lists or matrices of data.

- **Functional**: array-oriented programming is very different to the procedural mindset, once one acclimatises to this, the necessity of immutability and function-purity is trivial. Solution does not enforce functional rules, however, the syntax encourages a functional style.

### ⏰ When?

**Features**: in a short time, Solution will operate correctly with all syntactic features: this includes but is not limited to packages, namespaces, compiler directives, and operators (functions that receive functions). This does not mean the language is complete, simply that future features will extend rather than change the language.

**CUDA**: ability to emit LLVM which can be compiled for NVIDIA GPUs will be added once the lanauge is stable.

**Currently**: the compiler isn't doing much, apart from taking up a lot of my spare time...
