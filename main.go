package main

import (
	"flag"
	"fmt"

	"github.com/LRz00/cli-lint/analyser"
	"github.com/LRz00/cli-lint/filesystem"
)

func main() {
	path := flag.String("path", ".", "Path to the file or directory to lint")
	lang := flag.String("lang", "", "Programming language of the CLI application (python or javascript)")
	output := flag.String("output", "lint-report.md", "Output file for the lint report (markdown)")

	flag.Parse()

	files, err := filesystem.CollectFiles(*path, *lang)
	if err != nil {
		fmt.Printf("Error collecting files: %v\n", err)
		return
	}

	fmt.Printf("%s files found: %d\n", *lang, len(files))

	if *lang == "javascript" && len(files) > 0 {
		suggestions := analyser.AnalyzeJavaScript(files)
		if err := analyser.WriteMarkdownReport(suggestions, *output); err != nil {
			fmt.Printf("Error writing report: %v\n", err)
			return
		}
		fmt.Printf("Report generated: %s\n", *output)
	}
}
