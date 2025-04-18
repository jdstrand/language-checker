package printer

import (
	"go/token"

	"github.com/jdstrand/language-checker/pkg/result"
	"github.com/jdstrand/language-checker/pkg/rule"
)

func generateFileResult() *result.FileResults {
	r := result.FileResults{Filename: "foo.txt"}
	r.Results = generateResults(r.Filename)
	return &r
}

func generateResults(filename string) []result.Result {
	return []result.Result{
		result.LineResult{
			Rule:    &rule.TestRule,
			Finding: "whitelist",                  // langcheckignore:rule=whitelist
			Line:    "this whitelist must change", // langcheckignore:rule=whitelist
			StartPosition: &token.Position{
				Filename: filename,
				Offset:   0,
				Line:     1,
				Column:   6,
			},
			EndPosition: &token.Position{
				Filename: filename,
				Offset:   0,
				Line:     1,
				Column:   15,
			},
		},
	}
}

func generateSecondFileResult() *result.FileResults {
	r := result.FileResults{Filename: "bar.txt"}
	r.Results = generateSecondResults(r.Filename)
	return &r
}

func generateSecondResults(filename string) []result.Result {
	return []result.Result{
		result.LineResult{
			Rule:    &rule.TestErrorRule,
			Finding: "slave",                       // langcheckignore:rule=slave
			Line:    "this slave term must change", // langcheckignore:rule=slave
			StartPosition: &token.Position{
				Filename: filename,
				Offset:   0,
				Line:     1,
				Column:   6,
			},
			EndPosition: &token.Position{
				Filename: filename,
				Offset:   0,
				Line:     1,
				Column:   15,
			},
		},
	}
}

func generateThirdFileResult() *result.FileResults {
	r := result.FileResults{Filename: "barfoo.txt"}
	r.Results = generateThirdResults(r.Filename)
	return &r
}

func generateThirdResults(filename string) []result.Result {
	return []result.Result{
		result.LineResult{
			Rule:    &rule.TestInfoRule,
			Finding: "test",
			Line:    "this test must change",
			StartPosition: &token.Position{
				Filename: filename,
				Offset:   0,
				Line:     1,
				Column:   6,
			},
			EndPosition: &token.Position{
				Filename: filename,
				Offset:   0,
				Line:     1,
				Column:   15,
			},
		},
	}
}

func generateFilePathResult() *result.FileResults {
	r := result.FileResults{Filename: "whitelist.txt"}
	r.Results = generatePathResults(r.Filename)
	return &r
}

func generatePathResults(filename string) []result.Result {
	return []result.Result{
		result.LineResult{
			Rule:    &rule.TestRule,
			Finding: "whitelist",
			Line:    "this whitelist must change",
			StartPosition: &token.Position{
				Filename: filename,
				Offset:   0,
				Line:     1,
				Column:   1,
			},
			EndPosition: &token.Position{
				Filename: filename,
				Offset:   0,
				Line:     1,
				Column:   1,
			},
		},
	}
}

func newPosition(f string, l, c int) *token.Position {
	return &token.Position{
		Filename: f,
		Offset:   0,
		Line:     l,
		Column:   c,
	}
}
