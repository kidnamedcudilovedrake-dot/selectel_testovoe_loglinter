package main

import (
	"github.com/selectel-tasks/loglint/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

func New(conf any) ([]*analysis.Analyzer, error) {
	cfg, err := analyzer.ParseConfig(conf)
	if err != nil {
		return nil, err
	}
	analyzer.GlobalConfig = cfg
	return []*analysis.Analyzer{analyzer.Analyzer}, nil
}

func main() {}
