package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"wordlang/interpreter"
	"wordlang/lexer"
	"wordlang/object"
	"wordlang/parser"
)

// main is the entry point of the WordLang interpreter.
// It reads a file specified as a command-line argument, parses it,
// and evaluates the resulting program.
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: wordlang <filename>")
		return
	}

	filename := os.Args[1]
	runFile(filename)
}

// runFile reads, parses, and evaluates a WordLang program from a file.
func runFile(filename string) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		return
	}

	input := string(content)
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		printParserErrors(p.Errors())
		return
	}

	fmt.Println("\n--- AST ---")
	fmt.Println(program.String())
	fmt.Println("--- End AST ---\n")

	env := interpreter.NewEnvironment()
	result := interpreter.Eval(program, env)

	if result != nil && result.Type() == object.ERROR_OBJ {
		fmt.Println(result.Inspect())
	}
}

// printParserErrors prints parser error messages to the console.
func printParserErrors(errors []string) {
	fmt.Println("Parser errors:")
	for _, msg := range errors {
		fmt.Println("\t" + msg)
	}
}

// repl starts a Read-Eval-Print Loop for interactive WordLang execution.
// This function is not called in main by default and is intended for testing.
func repl() {
	reader := bufio.NewReader(os.Stdin)
	env := interpreter.NewEnvironment()

	for {
		fmt.Print("WordLang > ")
		line, _ := reader.ReadString('\n')
		if line == "exit\n" {
			break
		}
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(p.Errors())
			continue
		}

		result := interpreter.Eval(program, env)
		if result != nil {
			fmt.Println(result.Inspect())
		}
	}
}
