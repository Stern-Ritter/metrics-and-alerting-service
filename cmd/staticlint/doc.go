// Package provides a static analysis tool that integrates multiple
// analyzers for checking Go code.
//
// Usage:
//	 go build -o multichecker ./multichecker.go
//   ./multichecker [files]
//
// Standard static analyzers:
// - assign.Analyzer:
//     Checks for useless assignments.
//
// - atomic.Analyzer:
//     Checks for common mistakes using the sync/atomic package.
//
// - bools.Analyzer:
//     Checks for common mistakes involving boolean operators.
//
// - buildtag.Analyzer:
//     Checks go:build and build directives.
//
// - composite.Analyzer:
//     Checks for unkeyed composite literals.
//
// - copylock.Analyzer:
//     Checks for locks erroneously passed by value.
//
// - errorsas.Analyzer:
//     Reports passing non-pointer or non-error values to errors.As.
//
// - httpresponse.Analyzer:
//     Checks for mistakes using HTTP responses.
//
// - loopclosure.Analyzer:
//     Checks references to loop variables from within nested functions.
//
// - nilfunc.Analyzer:
//     Checks for useless comparisons between functions and nil.
//
// - nilness.Analyzer:
//     Checks for redundant or impossible nil comparisons.
//
// - printf.Analyzer:
//     Checks consistency of Printf format strings and arguments.
//
// - sortslice.Analyzer:
//     Checks the argument type of sort.Slice.
//
// - stdmethods.Analyzer:
//     Checks the signature of methods of well-known interfaces.
//
// - structtag.Analyzer:
//     Checks that struct field tags conform to reflect.StructTag.Get.
//
// - tests.Analyzer:
//     Checks for common mistaken usages of tests and examples.
//
// - unmarshal.Analyzer:
//     Reports passing non-pointer or non-interface values to unmarshal.
//
// - unreachable.Analyzer:
//     Checks for unreachable code.
//
// This package uses all analyzers from the following packages:
// - "honnef.co/go/tools/staticcheck":
//     Provides analyzers that find bugs and performance issues.
// - "honnef.co/go/tools/simple":
//     Provides analyzers that simplify code.
// - "honnef.co/go/tools/stylecheck"
//     Provides style check analyzers.
//     Exclude 'ST1000' analyzer that checks at least one file in a non-main package should have a package comment.
//
// Additional Analyzers:
// - ineffassign.Analyzer:
//     Checks for ineffectual assignments.
//
// - osexitanalyzer.OsExitAnalyzer:
//     Checks for os.Exit calls in main function.
//

package main
