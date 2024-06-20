package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/prometheus1400/kel/src/llvm"
	"github.com/prometheus1400/kel/src/parser"
	"github.com/prometheus1400/kel/src/scanner"
)

func runRepl() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, _ := reader.ReadBytes('\n')
		if len(line) == 0 {
			break
		}
		run(line)
	}
}

func runFile(filePath string) {
	src, _ := os.ReadFile(filePath)
	run(src)
}

func run(source []byte) {
	scanner := scanner.NewScanner()
	scanner.Scan(source)

	if scanner.HadError {
		fmt.Println("scanner had errors:")
		scanner.ReportErrors()
		os.Exit(1)
	}
	// scanner.PrintTokens()

	parser := parser.NewParser()
	stmts := parser.Parse(scanner.Tokens)
	// parser.Parse(scanner.Tokens)
	if parser.HadError {
		fmt.Println("parser had errors:")
		parser.ReportErrors()
		os.Exit(1)
	}

	gen := llvm.NewIRGenerator()
	gen.GenerateIR(stmts, "example")

	// interpreter := interpreter.NewTreeWalkInterpreter()
	// interpreter.Interpret(&stmts)
}

func main() {
	args := os.Args

	switch len(args) {
	case 1:
		runRepl()
	case 2:
		runFile(args[1])
	default:
		fmt.Fprintf(os.Stderr, "error")
	}
}
