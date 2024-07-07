// Package osexitanalyzer provides an analyzer to check for os.Exit calls in the main function.
package osexitanalyzer

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const (
	targetPkgName = "main"
	targetFuncName
	generatedCodeCommentPrefix = "// Code generated"
)

func NewOsExitAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "osexit",
		Doc:  "check for os.Exit calls in main function",
		Run:  run,
	}
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		if isGeneratedFile(pass.Fset, f) {
			continue
		}

		if pass.Pkg.Name() != targetPkgName {
			continue
		}

		ast.Inspect(f, func(n ast.Node) bool {
			if fn, ok := isTargetFunc(n, targetFuncName); ok {
				checkOsExitCall(pass, pass.Fset, f, fn)
			}
			return true
		})
	}
	return nil, nil
}

func isGeneratedFile(fset *token.FileSet, file *ast.File) bool {
	position := fset.Position(file.Pos())
	filePath := position.Filename

	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return false
	}

	for _, gr := range f.Comments {
		for _, c := range gr.List {
			if strings.HasPrefix(c.Text, generatedCodeCommentPrefix) {
				return true
			}
		}
	}
	return false
}

func isTargetFunc(node ast.Node, funcName string) (*ast.FuncDecl, bool) {
	fn, ok := node.(*ast.FuncDecl)
	return fn, ok && fn.Name.Name == funcName
}

func checkOsExitCall(pass *analysis.Pass, fset *token.FileSet, file *ast.File, fn *ast.FuncDecl) {
	position := fset.Position(file.Pos())
	filePath := position.Filename
	log.Printf("Inspecting file: %s function: %s.%s\n", filePath, pass.Pkg.Name(), fn.Name.Name)
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if isOsExitCall(callExpr) {
			msg := formatErrorReportMsg(pass, fn.Name.Name, callExpr)
			log.Printf("Found os.Exit call: %s\n", msg)
			pass.Reportf(callExpr.Pos(), msg)
			return false
		}
		return true
	})
}

func isOsExitCall(callExpr *ast.CallExpr) bool {
	selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	pkgIdent, ok := selExpr.X.(*ast.Ident)
	if !ok {
		return false
	}

	pkgName := pkgIdent.Name
	funcName := selExpr.Sel.Name
	isOsExit := pkgName == "os" && funcName == "Exit"
	log.Printf("Checking call: %s.%s - is os.Exit: %t\n", pkgName, funcName, isOsExit)
	return isOsExit
}

func formatErrorReportMsg(pass *analysis.Pass, funcName string, callExpr *ast.CallExpr) string {
	pos := pass.Fset.Position(callExpr.Pos())
	return fmt.Sprintf("using os.Exit call in %s func at line %d, column %d, full call expr: %s",
		funcName, pos.Line, pos.Column, exprToString(callExpr))
}

func exprToString(expr ast.Expr) string {
	var buf bytes.Buffer
	err := format.Node(&buf, token.NewFileSet(), expr)
	if err != nil {
		return "<error>"
	}
	return buf.String()
}
