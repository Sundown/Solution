# The How of Solution

This document will begin to explain the foundations of the Solution language.

Solution is APL, any knowledge from versions of APL such as that maintained by Dyalog will be easily transferred to Solution, only requiring relearning symbols.

## Syntax

Syntax in Solution is simple, there are monadic and dyadic operators which take one or two arguments respectively. Arguments may be single atoms, or vectors of such. Furthermore there are also operators which take as their arguments functions instead of atoms or vectors.

A monadic function will accept a single argument which consists of all that code which is to the left of it.

A dyadic function will accept two arguments, the first of which is the morpheme directly to its left, and the second of which—like the monadic function—is all that code to the right of it.
