package analyser

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/LRz00/cli-lint/common"
)

// AnalyzeJavaScript analyzes JavaScript files and returns lint suggestions.
func AnalyzeJavaScript(files []string) []common.LintSuggestion {
	var suggestions []common.LintSuggestion

	var (
		totalLines    int
		totalSemis    int
		singleQuotes  int
		doubleQuotes  int
		consoleFiles  int
		varFiles      int
		weakCompFiles int
		unusedVars    int
	)

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		code := string(content)

		if usesVar(code) {
			varFiles++
		}

		semiStats := analyzeSemi(code)
		totalLines += semiStats.lines
		totalSemis += semiStats.semis

		quoteStats := analyzeQuotes(code)
		singleQuotes += quoteStats.single
		doubleQuotes += quoteStats.double

		if usesWeakComparisons(code) {
			weakCompFiles++
		}

		if usesConsole(code) {
			consoleFiles++
		}

		unusedVars += countUnusedVars(code)
	}

	fileCount := len(files)

	if varFiles > 0 {
		suggestions = append(suggestions, common.LintSuggestion{
			Rule:   "no-var",
			Reason: fmt.Sprintf("Uso de 'var' detectado em %d arquivo(s)", varFiles),
		})
	}

	if totalLines > 0 {
		ratio := float64(totalSemis) / float64(totalLines)
		if ratio < 0.2 {
			suggestions = append(suggestions, common.LintSuggestion{
				Rule:   "semi",
				Reason: "Pouco uso de ponto e vírgula no projeto",
			})
		}
	}

	totalQuotes := singleQuotes + doubleQuotes
	if totalQuotes > 0 {
		ratio := float64(singleQuotes) / float64(totalQuotes)
		if ratio > 0.7 {
			suggestions = append(suggestions, common.LintSuggestion{
				Rule:   "quotes",
				Reason: "Predominância de aspas simples",
			})
		} else if ratio < 0.3 {
			suggestions = append(suggestions, common.LintSuggestion{
				Rule:   "quotes",
				Reason: "Predominância de aspas duplas",
			})
		}
	}

	if weakCompFiles > 0 {
		suggestions = append(suggestions, common.LintSuggestion{
			Rule:   "eqeqeq",
			Reason: fmt.Sprintf("Comparações fracas (== ou !=) em %d arquivo(s)", weakCompFiles),
		})
	}

	if consoleFiles > 0 && fileCount > 0 {
		ratio := float64(consoleFiles) / float64(fileCount)
		if ratio > 0.3 {
			suggestions = append(suggestions, common.LintSuggestion{
				Rule:   "no-console",
				Reason: fmt.Sprintf("Uso frequente de console em %d arquivo(s)", consoleFiles),
			})
		}
	}

	if unusedVars > 0 {
		suggestions = append(suggestions, common.LintSuggestion{
			Rule:   "no-unused-vars",
			Reason: fmt.Sprintf("%d variável(is) declarada(s) e aparentemente não utilizada(s)", unusedVars),
		})
	}

	return suggestions
}

// --- Rule detectors (private) ---

func usesVar(code string) bool {
	return strings.Contains(code, "var ")
}

func usesWeakComparisons(code string) bool {
	return strings.Contains(code, "==") || strings.Contains(code, "!=")
}

type semiStats struct {
	lines int
	semis int
}

func analyzeSemi(code string) semiStats {
	return semiStats{
		lines: len(strings.Split(code, "\n")),
		semis: strings.Count(code, ";"),
	}
}

type quoteStats struct {
	single int
	double int
}

func analyzeQuotes(code string) quoteStats {
	return quoteStats{
		single: strings.Count(code, "'"),
		double: strings.Count(code, "\""),
	}
}

func usesConsole(code string) bool {
	return strings.Contains(code, "console.")
}

var varDeclRegex = regexp.MustCompile(`\b(const|let)\s+([a-zA-Z_][a-zA-Z0-9_]*)`)

func countUnusedVars(code string) int {
	matches := varDeclRegex.FindAllStringSubmatch(code, -1)
	unused := 0

	for _, m := range matches {
		varName := m[2]
		if strings.Count(code, varName) == 1 {
			unused++
		}
	}

	return unused
}
