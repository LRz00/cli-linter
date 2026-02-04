package analyser

import (
	"fmt"
	"os"
	"strings"

	"github.com/LRz00/cli-lint/common"
)

// WriteMarkdownReport writes lint suggestions to a markdown file.
func WriteMarkdownReport(suggestions []common.LintSuggestion, outputFile string) error {
	var sb strings.Builder

	sb.WriteString("#Relatório de Lint\n\n")

	if len(suggestions) == 0 {
		sb.WriteString("**Nenhuma sugestão de lint encontrada.**\n")
	} else {
		sb.WriteString(fmt.Sprintf("Foram encontradas **%d** sugestões de regras ESLint para o projeto.\n\n", len(suggestions)))
		sb.WriteString("## Sugestões\n\n")
		sb.WriteString("| Regra | Motivo |\n")
		sb.WriteString("|-------|--------|\n")

		for _, s := range suggestions {
			sb.WriteString(fmt.Sprintf("| `%s` | %s |\n", s.Rule, s.Reason))
		}
	}

	return os.WriteFile(outputFile, []byte(sb.String()), 0644)
}
