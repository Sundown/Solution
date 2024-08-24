#!/bin/bash

mkdir build
rm -rf build/*
go run solution.go solutions/algebraic.sol -emit purellvm
mv dev.ll build
clang libsol.c -S -emit-llvm -o build/libsol.ll


clang build/libsol.ll build/dev.ll -o dev
