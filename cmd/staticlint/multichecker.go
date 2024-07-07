package main

import (
	"github.com/gordonklaus/ineffassign/pkg/ineffassign"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"

	"github.com/Stern-Ritter/metrics-and-alerting-service/cmd/staticlint/osexitanalyzer"
)

func main() {
	var mychecks []*analysis.Analyzer

	// add standard analyzers
	mychecks = append(mychecks,
		assign.Analyzer,       // check for useless assignments
		atomic.Analyzer,       // check for common mistakes using the sync/atomic package
		bools.Analyzer,        // check for common mistakes involving boolean operators
		buildtag.Analyzer,     // check //go:build and // +build directives
		composite.Analyzer,    // check for unkeyed composite literals
		copylock.Analyzer,     // check for locks erroneously passed by value
		errorsas.Analyzer,     // report passing non-pointer or non-error values to errors.As
		httpresponse.Analyzer, // check for mistakes using HTTP responses
		loopclosure.Analyzer,  // check references to loop variables from within nested functions
		nilfunc.Analyzer,      // check for useless comparisons between functions and nil
		nilness.Analyzer,      // check for redundant or impossible nil comparisons
		printf.Analyzer,       // check consistency of Printf format strings and arguments
		sortslice.Analyzer,    // check the argument type of sort.Slice
		stdmethods.Analyzer,   // check signature of methods of well-known interfaces
		structtag.Analyzer,    // check that struct field tags conform to reflect.StructTag.Get
		tests.Analyzer,        // check for common mistaken usages of tests and examples
		unmarshal.Analyzer,    // report passing non-pointer or non-interface values to unmarshal
		unreachable.Analyzer,  // check for unreachable code
	)

	// excluded staticcheck.io analyzers
	excludeChecks := map[string]bool{
		"ST1000": true,
	}

	// add analyzers that find bugs and performance issues (SA*)
	for _, a := range staticcheck.Analyzers {
		if _, in := excludeChecks[a.Analyzer.Name]; !in {
			mychecks = append(mychecks, a.Analyzer)
		}
	}

	// add analyzers that simplify code (S*)
	for _, a := range simple.Analyzers {
		if _, in := excludeChecks[a.Analyzer.Name]; !in {
			mychecks = append(mychecks, a.Analyzer)
		}
	}

	// add style check analyzers (ST*)
	for _, a := range stylecheck.Analyzers {
		if _, in := excludeChecks[a.Analyzer.Name]; !in {
			mychecks = append(mychecks, a.Analyzer)
		}
	}

	// add additional analyzers
	mychecks = append(mychecks,
		ineffassign.Analyzer,               // check for ineffectual assignments
		osexitanalyzer.NewOsExitAnalyzer(), // check for os.Exit calls in main function
	)

	multichecker.Main(
		mychecks...,
	)
}
