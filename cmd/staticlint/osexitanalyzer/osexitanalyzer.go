// Package osexitanalyzer provides an analyzer to check for os.Exit calls in the main function.
package osexitanalyzer

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const (
	targetPkgName = "main"
	targetFuncName
	generatedCodeCommentPrefix = "// Code generated"
	errMsg                     = "using os.Exit in main func"
)

var OsExitAnalyzer = &analysis.Analyzer{
	Name: "osexit",
	Doc:  "check for os.Exit calls in main function",
	Run:  run,
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
				checkOsExitCall(pass, fn)
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

func checkOsExitCall(pass *analysis.Pass, fn *ast.FuncDecl) {
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if isOsExitCall(callExpr) {
			pass.Reportf(callExpr.Pos(), errMsg)
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
	return pkgName == "os" && funcName == "Exit"
}
