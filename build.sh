#!/bin/sh

rm ./build/main
go build --tags=llvm17 ./src/main.go #-o ./build/main
mv main ./build/
