package common

const (
	javascript = "javascript"
	python     = "python"
)

var SupportedFileFormats = map[string]string{
	javascript: ".js",
	python:     ".py",
}

type LintSuggestion struct {
	Rule   string
	Reason string
}
