#!/bin/sh

# export LLVM_CONFIG="/opt/homebrew/Cellar/llvm@17/bin/llvm-config"
# export PATH="/opt/homebrew/opt/llvm@17/bin:$PATH"
go run -tags=llvm17 ./src/
# echo $LLVM_CONFIG
