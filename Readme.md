<h3 align="center"> Solution</h3>

<p align="center">
Solution is a compiler for an array-oriented language, providing the cognition of APL in an accessible, compiled, and open-source platform. The Solution Language is inspired by the work of Kenneth Iverson and Dyalog.
</p>


### 🔍 Why?

- **Array oriented**: what do our computers process the most? A trivial but accurate answer is *data*, but what is data? Data is the plural form of the Latin word *datum*, which means *a piece of information*. This answer is a joke, yet it is useful; lending the insight that most processing is not performed on a single datum, but rather groups of these, again: data. Despite this striking fact, most languages prefer to focus on a single datum at a time, does this not seem a little strange, considering how much Computer Science is focused on graphs, sets, and other types of groups? Array oriented languages look at the bigger picture, preferring to manipulate lists or matrices of data. 

- **Functional**: array-oriented programming is very different to the procedural mindset, once one acclimatises to this, the necessity of immutability and function-purity is trivial. Solution does not enforce functional rules, however, the syntax encourages a functional style. 

### ⏰ When?

Solution has two major stages ahead. 

- **Features**: in a short time, Solution will operate correctly with all syntactic features: this includes but is not limited to packages, namespaces, compiler directives, and operators (functions that receive functions). This does not mean the language is complete, simply that future features will extend rather than change the language. 

- **Re-write**: the C++ re-write is intended to occur after the aforementioned stage. To keep things simple, no changes will be made to the language during this transition. 

- **Currently**: at present, the language is unable to provide a Solution to many problems, apart from taking up a lot of spare time... 

### 🪜 Implementation

- **Compiled**: many languages of this type are interpreted in some way, this makes learning and debugging much easier, but hurts performance. Being compiled opens the possibility to extend Solution into the GPU world using NVCC. The optimisations performed by LLVM are more advanced than anything possible in a project this size.

- **Golang**: the development of Solution occurring in Golang is an unfortunate fluke, mostly due to the language's simplicity. At a certain point in development, the compiler will be re-written in C++, this will be a linear process and won't change any underlying functionality. 

### 📦 Distribution

- **Packaging**: no efforts will be made to package Solution before the C++ re-write. Currently, git clone and go build do the job. 

- **Libraries**: if this project matures to the state of requiring a package manager, git will do the hard work similar to Go. 

### 🏷 Name

Programming languages are sometimes thought of as tools for solving problems, however, it appears more accurate to say programming languages are a way of describing a problem to a compiler, which will produce a tool to solve your problem, or put simply: a **Solution**. This parlance is also used in the Navy when computers calculate the launching/firing vector of a projectile weapon, the result is called a *firing solution*. 

<p align="center"><a href="https://github.com/sundown/solution/blob/main/LICENSE"><img src="https://img.shields.io/static/v1.svg?style=for-the-badge&label=License&message=GPL-2.0&logoColor=1f1f1f&colorA=1f1f1f&colorB=f0f0f0"/></a></p>
