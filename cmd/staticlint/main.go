package main

import (
	"github.com/gordonklaus/ineffassign/pkg/ineffassign"
	"github.com/kisielk/errcheck/errcheck"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	// Standard library analyzers - comprehensive set for learning Go best practices
	checks := []*analysis.Analyzer{
		// Assembly and low-level checks
		asmdecl.Analyzer,   // Check assembly declarations
		unsafeptr.Analyzer, // Check unsafe.Pointer conversions

		// Common mistakes and bugs
		assign.Analyzer,          // Check for useless assignments
		atomic.Analyzer,          // Check for common mistakes using sync/atomic
		bools.Analyzer,           // Check for common mistakes with boolean operators
		copylock.Analyzer,        // Check for locks that are erroneously passed by value
		deepequalerrors.Analyzer, // Check for inappropriate use of reflect.DeepEqual with errors
		errorsas.Analyzer,        // Check for passing non-pointer to errors.As
		httpresponse.Analyzer,    // Check for mistakes using HTTP responses
		loopclosure.Analyzer,     // Check for references to loop variables from nested functions
		lostcancel.Analyzer,      // Check for failure to call a context cancellation function
		nilfunc.Analyzer,         // Check for useless comparisons against nil
		nilness.Analyzer,         // Check for redundant or impossible nil comparisons
		stringintconv.Analyzer,   // Check for string(int) conversions
		unmarshal.Analyzer,       // Check for passing non-pointer or non-interface to unmarshal
		unreachable.Analyzer,     // Check for unreachable code
		unusedresult.Analyzer,    // Check for unused results of calls to certain functions

		// Printf-like format string checks
		printf.Analyzer, // Check consistency of Printf format strings and arguments

		// Struct and interface checks
		structtag.Analyzer,  // Check struct field tags are well formed
		composite.Analyzer,  // Check for unkeyed composite literals
		stdmethods.Analyzer, // Check signature of methods of well-known interfaces

		// Testing-related checks
		tests.Analyzer,            // Check for common mistakes in tests
		testinggoroutine.Analyzer, // Check for calls to Fatal from goroutines started by a test

		// Shift operations
		shift.Analyzer, // Check for shifts that exceed the width of an integer

		// Sorting checks
		sortslice.Analyzer, // Check for calls to sort.Slice that do not use a slice type

		// Shadow variables (commonly caught bug for beginners)
		shadow.Analyzer, // Check for shadowed variables

		// Build tags
		buildtag.Analyzer, // Check build tags are well formed
	}

	// Staticcheck analyzers - industry-standard static analysis
	for _, v := range staticcheck.Analyzers {
		// Add all SA checks - these catch serious bugs and issues
		checks = append(checks, v.Analyzer)
	}

	// Simple analyzers - suggest code simplifications
	for _, v := range simple.Analyzers {
		// Add all S checks - these suggest simpler alternatives
		checks = append(checks, v.Analyzer)
	}

	// NOTE: Stylecheck analyzers are disabled because they enforce documentation requirements
	// including package comments, exported function comments, etc.
	// Uncomment below if you want to enforce Go style guide including documentation:
	//
	// for _, v := range stylecheck.Analyzers {
	//     checks = append(checks, v.Analyzer)
	// }

	// Third-party critical analyzers
	checks = append(checks,
		errcheck.Analyzer,    // Check that error return values are used
		ineffassign.Analyzer, // Detect ineffectual assignments
	)

	multichecker.Main(
		checks...,
	)
}
