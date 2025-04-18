package parser

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jdstrand/language-checker/pkg/result"
	"github.com/jdstrand/language-checker/pkg/rule"
	"github.com/jdstrand/language-checker/pkg/util"

	"github.com/rs/zerolog/log"
)

func (p *Parser) generateFileFindingsFromFilename(filename string) (*result.FileResults, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return p.generateFileFindings(file)
}

// generateFileFindings reads the file and returns results of places where rules are broken
// this function will not close the file, that should be handled by the caller
func (p *Parser) generateFileFindings(file *os.File) (*result.FileResults, error) {
	filename := filepath.ToSlash(file.Name())
	start := time.Now()
	defer func() {
		log.Debug().
			TimeDiff("durationMS", time.Now(), start).
			Str("file", filename).
			Msg("finished processing findings")
	}()

	results := &result.FileResults{
		Filename: filename,
	}

	// Check for findings in the filename itself
	for _, pathResult := range result.MatchPathRules(p.Rules, file.Name()) {
		results.Results = append(results.Results, pathResult)
	}

	// Don't check file content if it's not a text file or file is empty
	if err := util.IsTextFileFromFilename(filename); err != nil {
		log.Debug().Str("file", filename).Str("reason", err.Error()).Msg("skipping content")
		return results, nil
	}

	reader := bufio.NewReader(file)

	var ignoreNextLineText string
	line := 1

Loop:
	for {
		switch text, err := reader.ReadString('\n'); {
		case err == nil || (err == io.EOF && text != ""):
			text = strings.TrimSuffix(text, "\n")

			// Store current line's langcheckignore text if ignoring next line
			if rule.IsDirectiveOnlyLine(text) {
				ignoreNextLineText = text
				line++
				continue
			}

			for _, r := range p.Rules {
				if p.Ignorer != nil {
					if ignoreNextLineText == "" && r.CanIgnoreLine(text) {
						log.Debug().
							Str("rule", r.Name).
							Str("file", filename).
							Int("line", line).
							Msg("ignoring via in-line")
						continue
					} else if r.CanIgnoreLine(ignoreNextLineText) {
						// Check current rule against prev line's next-line langcheckignore text (if applicable)
						log.Debug().
							Str("rule", r.Name).
							Str("file", filename).
							Int("line", line).
							Msg("ignoring via next-line")
						continue
					}
				}

				lineResults := result.FindResults(r, results.Filename, text, line)
				results.Results = append(results.Results, lineResults...)
			}

			ignoreNextLineText = ""
			line++
		case err == io.EOF:
			break Loop
		case err != nil:
			return nil, err
		}
	}

	return results, nil
}
