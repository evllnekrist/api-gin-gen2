package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/scanner"
	"go/token"
	"net/http"
	// "log"
	"os"
)

type AsmController struct{}

var sample_src = []byte(`
		package main
		import "fmt"
		func main(){
			fmt.Println("Hello, hooman")
		} `)

func ASM_Open(title string, subtitle string) {

	fmt.Println("\n\n")
	fmt.Println("assy ~ ------------------------ ~ assy ~ ------------------------ ~ assy ")
	fmt.Println("assy ~                            ____                            ~ assy ")
	fmt.Println("assy ~                                                            ~ assy ")
	fmt.Println("assy ~                            OPEN                            ~ assy ")
	fmt.Println("assy ~                            ____                            ~ assy ")
	fmt.Println("assy ~                                                            ~ assy ")
	fmt.Println("assy ~ ------------------------ " + title + " ------------------------ ~ assy ")
	fmt.Println("assy ~ DETAIL ACT :: " + subtitle)
	fmt.Println("\n")
}

func ASM_Close(title string, subtitle string) {

	fmt.Println("\n")
	fmt.Println("assy ~ DETAIL ACT :: " + subtitle)
	fmt.Println("assy ~ ------------------------ " + title + " ------------------------ ~ assy ")
	fmt.Println("assy ~                            ____                            ~ assy ")
	fmt.Println("assy ~                                                            ~ assy ")
	fmt.Println("assy ~                           CLOSED                           ~ assy ")
	fmt.Println("assy ~                            ____                            ~ assy ")
	fmt.Println("assy ~                                                            ~ assy ")
	fmt.Println("assy ~ ------------------------ ~ assy ~ ------------------------ ~ assy ")
	fmt.Println("\n\n")
}

// ----------------------------------------------------------------------------------- ASSEMBLER (MACHINE LANGUAGE)
// CASE [1]:
// The scanner, which converts the source code into a list of tokens, for use by the parser.
// The parser, which converts the tokens into an Abstract Syntax Tree to be used by code generation.
// The code generation, which converts the Abstract Syntax Tree to machine code.
// url :: https://getstream.io/blog/how-a-go-program-compiles-down-to-machine-code/

func (ctrl AsmController) Scanner(c *gin.Context) {

	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(sample_src))
	s.Init(file, sample_src, nil, 0)

	ASM_Open("SCANNER.", "") //just string to mark the print -- start
	for {
		pos, tok, lit := s.Scan()
		fmt.Printf("%-6s%-8s%q\n", fset.Position(pos), tok, lit)

		if tok == token.EOF {
			break
		}
	}
	ASM_Close("SCANNER.", "") //just string to mark the print -- end
	http.Redirect(c.Writer, c.Request, "localhost:8080/", http.StatusSeeOther)
} // break up the raw source code text into tokens. // the result explains why Go does not need semicolons: they are placed intelligently by the scanner

func (ctrl AsmController) Parser(c *gin.Context) {

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", sample_src, 0)
	helper_basic.Panics(err)

	act := c.Param("act")
	ASM_Open(" PARSER ", act)
	if act == "print" {
		ast.Print(fset, file)
	} else if act == "inspect" {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			printer.Fprint(os.Stdout, fset, call.Fun)
			fmt.Println()
			return false
		})
	} else {
		fmt.Println(" INCORRECT PARAMETER GIVEN ")
	}
	ASM_Close(" PARSER ", act)
	http.Redirect(c.Writer, c.Request, "localhost:8080/", http.StatusSeeOther)
} // a phase of the compiler that converts the tokens into an Abstract Syntax Tree (AST)

func (ctrl AsmController) CodeGenerator(c *gin.Context) {
	http.Redirect(c.Writer, c.Request, "https://github.com/golang/go/tree/master/src/cmd/compile/internal/ssa", http.StatusSeeOther)
}

// CASE [1] END
