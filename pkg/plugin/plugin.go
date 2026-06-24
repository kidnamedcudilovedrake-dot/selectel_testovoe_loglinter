package plugin

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
	"github.com/selectel-tasks/loglint/pkg/analyzer"
)

func init() {
	register.Plugin("loglint", newAnalyzer)
}

type loglintPlugin struct {
	cfg analyzer.Config
}

func (p *loglintPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	analyzer.GlobalConfig = p.cfg
	return []*analysis.Analyzer{analyzer.Analyzer}, nil
}

func (p *loglintPlugin) GetLoadMode() string {
	return "types"
}

func newAnalyzer(conf any) (register.LinterPlugin, error) {
	cfg, err := analyzer.ParseConfig(conf)
	if err != nil {
		return nil, err
	}
	return &loglintPlugin{cfg: cfg}, nil
}
