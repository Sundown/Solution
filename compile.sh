#!/bin/bash

mkdir build
rm -rf build/*
go run solution.go solutions/algebraic.sol
mv dev build
./build/dev
