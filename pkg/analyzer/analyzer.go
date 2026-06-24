package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var GlobalConfig = DefaultConfig()

var Analyzer = &analysis.Analyzer{
	Name: "loglint",
	Doc:  "checks log messages for style guidelines and sensitive data leaks",
	Run:  run,
}

func init() {
	Analyzer.Flags.BoolVar(&GlobalConfig.Rules.Lowercase, "lowercase", true, "check that log messages start with a lowercase letter")
	Analyzer.Flags.BoolVar(&GlobalConfig.Rules.EnglishOnly, "english-only", true, "check that log messages are in English only")
	Analyzer.Flags.BoolVar(&GlobalConfig.Rules.NoSpecialChars, "no-special-chars", true, "check that log messages don't contain special characters or emojis")
	Analyzer.Flags.BoolVar(&GlobalConfig.Rules.NoSensitive, "no-sensitive", true, "check that log calls don't leak sensitive data")
}

func run(pass *analysis.Pass) (any, error) {
	var rules []*regexp.Regexp
	for _, pat := range GlobalConfig.SensitivePatterns {
		if re, err := regexp.Compile(pat); err == nil {
			rules = append(rules, re)
		}
	}

	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			fn, pkg, ok := isLogger(pass, call)
			if !ok {
				return true
			}

			if GlobalConfig.Rules.NoSensitive {
				checkSec(pass, call, rules)
			}

			arg := logMsg(call, fn, pkg)
			if arg == nil {
				return true
			}

			val, lit, ok := strLit(pass, arg)
			if !ok {
				return true
			}

			if GlobalConfig.Rules.Lowercase {
				checkLower(pass, lit, val)
			}

			if GlobalConfig.Rules.EnglishOnly {
				checkLang(pass, lit, val)
			}

			if GlobalConfig.Rules.NoSpecialChars {
				checkChars(pass, lit, val)
			}

			return true
		})
	}
	return nil, nil
}

func isLogger(pass *analysis.Pass, call *ast.CallExpr) (string, string, bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", "", false
	}
	obj := pass.TypesInfo.ObjectOf(sel.Sel)
	if obj == nil || obj.Pkg() == nil {
		return "", "", false
	}
	pkg := obj.Pkg().Path()
	if pkg == "log/slog" || pkg == "go.uber.org/zap" || pkg == "log" {
		return obj.Name(), pkg, true
	}
	return "", "", false
}

func logMsg(call *ast.CallExpr, fn string, pkg string) ast.Expr {
	if len(call.Args) == 0 {
		return nil
	}

	switch pkg {
	case "log/slog":
		switch fn {
		case "InfoContext", "ErrorContext", "WarnContext", "DebugContext":
			if len(call.Args) >= 2 {
				return call.Args[1]
			}
		case "Log", "LogAttrs":
			if len(call.Args) >= 3 {
				return call.Args[2]
			}
		default:
			return call.Args[0]
		}
	case "go.uber.org/zap", "log":
		return call.Args[0]
	}
	return nil
}

func strLit(pass *analysis.Pass, expr ast.Expr) (string, *ast.BasicLit, bool) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			if val, err := strconv.Unquote(e.Value); err == nil {
				return val, e, true
			}
		}
	case *ast.BinaryExpr:
		if e.Op == token.ADD {
			return strLit(pass, e.X)
		}
	case *ast.CallExpr:
		sel, ok := e.Fun.(*ast.SelectorExpr)
		if !ok {
			return "", nil, false
		}
		obj := pass.TypesInfo.ObjectOf(sel.Sel)
		if obj != nil && obj.Pkg() != nil && obj.Pkg().Path() == "fmt" && obj.Name() == "Sprintf" {
			if len(e.Args) > 0 {
				return strLit(pass, e.Args[0])
			}
		}
	}
	return "", nil, false
}

func checkLower(pass *analysis.Pass, lit *ast.BasicLit, val string) {
	if !isLowerStart(val) {
		pass.Report(analysis.Diagnostic{
			Pos:     lit.Pos(),
			End:     lit.End(),
			Message: "log message should start with a lowercase letter",
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message: "lowercase the first letter",
					TextEdits: []analysis.TextEdit{
						{
							Pos:     lit.Pos(),
							End:     lit.End(),
							NewText: []byte(quoteString(toLowerStart(val), lit.Value)),
						},
					},
				},
			},
		})
	}
}

func checkLang(pass *analysis.Pass, lit *ast.BasicLit, val string) {
	if !isEnglish(val) {
		pass.Report(analysis.Diagnostic{
			Pos:     lit.Pos(),
			End:     lit.End(),
			Message: "log message should be in English only",
		})
	}
}

func checkChars(pass *analysis.Pass, lit *ast.BasicLit, val string) {
	errs := badChars(val, GlobalConfig.ForbiddenChars)
	if len(errs) > 0 {
		pass.Report(analysis.Diagnostic{
			Pos:     lit.Pos(),
			End:     lit.End(),
			Message: "log message " + strings.Join(errs, ", "),
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message: "clean special characters/emojis",
					TextEdits: []analysis.TextEdit{
						{
							Pos:     lit.Pos(),
							End:     lit.End(),
							NewText: []byte(quoteString(cleanChars(val, GlobalConfig.ForbiddenChars), lit.Value)),
						},
					},
				},
			},
		})
	}
}

func checkSec(pass *analysis.Pass, call *ast.CallExpr, rules []*regexp.Regexp) {
	for _, arg := range call.Args {
		ast.Inspect(arg, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.Ident:
				if isSec(node.Name, GlobalConfig.SensitiveKeywords, rules) {
					if obj := pass.TypesInfo.ObjectOf(node); obj != nil {
						if _, isPkg := obj.(*types.PkgName); isPkg {
							return true
						}
					}
					pass.Reportf(node.Pos(), "log call contains potentially sensitive variable %q", node.Name)
				}
			case *ast.IndexExpr:
				if lit, ok := node.Index.(*ast.BasicLit); ok && lit.Kind == token.STRING {
					if val, err := strconv.Unquote(lit.Value); err == nil {
						if isSec(val, GlobalConfig.SensitiveKeywords, rules) {
							pass.Reportf(node.Pos(), "log call contains potentially sensitive map key %q", val)
						}
					}
				}
			}
			return true
		})
	}
}

func isSec(name string, kws []string, rules []*regexp.Regexp) bool {
	lower := strings.ToLower(name)
	for _, kw := range kws {
		if strings.Contains(lower, strings.ToLower(kw)) {
			return true
		}
	}
	for _, re := range rules {
		if re.MatchString(name) {
			return true
		}
	}
	return false
}

func quoteString(val string, raw string) string {
	if strings.HasPrefix(raw, "`") {
		if !strings.Contains(val, "`") {
			return "`" + val + "`"
		}
	}
	return strconv.Quote(val)
}
