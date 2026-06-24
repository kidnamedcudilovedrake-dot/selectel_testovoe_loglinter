package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"
	"github.com/selectel-tasks/loglint/pkg/analyzer"
)

func main() {
	singlechecker.Main(analyzer.Analyzer)
}
