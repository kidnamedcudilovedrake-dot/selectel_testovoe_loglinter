package analyzer_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
	"github.com/selectel-tasks/loglint/pkg/analyzer"
)

func TestAll(t *testing.T) {
	// Restore default config in case other tests ran first
	analyzer.GlobalConfig = analyzer.DefaultConfig()
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer.Analyzer, "a")
}

func TestCustomPatterns(t *testing.T) {
	// Setup custom sensitive regex patterns
	analyzer.GlobalConfig = analyzer.DefaultConfig()
	analyzer.GlobalConfig.SensitivePatterns = []string{`^secret_.*$`, `^key_[0-9]+$`}

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer.Analyzer, "b")
}
