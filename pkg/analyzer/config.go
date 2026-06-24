package analyzer

type Config struct {
	Rules             RulesConfig `json:"rules" yaml:"rules"`
	SensitiveKeywords []string    `json:"sensitive_keywords" yaml:"sensitive_keywords"`
	SensitivePatterns []string    `json:"sensitive_patterns" yaml:"sensitive_patterns"`
	AllowedChars      string      `json:"allowed_chars" yaml:"allowed_chars"`
	ForbiddenChars    string      `json:"forbidden_chars" yaml:"forbidden_chars"`
}

type RulesConfig struct {
	Lowercase      bool `json:"lowercase" yaml:"lowercase"`
	EnglishOnly    bool `json:"english_only" yaml:"english_only"`
	NoSpecialChars bool `json:"no_special_chars" yaml:"no_special_chars"`
	NoSensitive    bool `json:"no_sensitive" yaml:"no_sensitive"`
}

func DefaultConfig() Config {
	return Config{
		Rules: RulesConfig{
			Lowercase:      true,
			EnglishOnly:    true,
			NoSpecialChars: true,
			NoSensitive:    true,
		},
		SensitiveKeywords: []string{
			"password", "passwd", "token", "secret", "key",
			"credential", "creds", "auth", "session", "sid",
		},
		ForbiddenChars: "!?",
	}
}

func ParseConfig(conf any) (Config, error) {
	cfg := DefaultConfig()
	if conf == nil {
		return cfg, nil
	}

	rawMap, ok := conf.(map[string]any)
	if !ok {
		anyMap, ok := conf.(map[any]any)
		if !ok {
			return cfg, nil
		}
		rawMap = make(map[string]any)
		for k, v := range anyMap {
			if s, ok := k.(string); ok {
				rawMap[s] = v
			}
		}
	}

	if rulesVal, ok := rawMap["rules"]; ok {
		if rulesMap, ok := rulesVal.(map[string]any); ok {
			if v, ok := rulesMap["lowercase"]; ok {
				if b, ok := v.(bool); ok {
					cfg.Rules.Lowercase = b
				}
			}
			if v, ok := rulesMap["english_only"]; ok {
				if b, ok := v.(bool); ok {
					cfg.Rules.EnglishOnly = b
				}
			}
			if v, ok := rulesMap["no_special_chars"]; ok {
				if b, ok := v.(bool); ok {
					cfg.Rules.NoSpecialChars = b
				}
			}
			if v, ok := rulesMap["no_sensitive"]; ok {
				if b, ok := v.(bool); ok {
					cfg.Rules.NoSensitive = b
				}
			}
		} else if rulesMap, ok := rulesVal.(map[any]any); ok {
			if v, ok := rulesMap["lowercase"]; ok {
				if b, ok := v.(bool); ok {
					cfg.Rules.Lowercase = b
				}
			}
			if v, ok := rulesMap["english_only"]; ok {
				if b, ok := v.(bool); ok {
					cfg.Rules.EnglishOnly = b
				}
			}
			if v, ok := rulesMap["no_special_chars"]; ok {
				if b, ok := v.(bool); ok {
					cfg.Rules.NoSpecialChars = b
				}
			}
			if v, ok := rulesMap["no_sensitive"]; ok {
				if b, ok := v.(bool); ok {
					cfg.Rules.NoSensitive = b
				}
			}
		}
	}

	if kwVal, ok := rawMap["sensitive_keywords"]; ok {
		if kwList, ok := kwVal.([]any); ok {
			var kws []string
			for _, item := range kwList {
				if s, ok := item.(string); ok {
					kws = append(kws, s)
				}
			}
			if len(kws) > 0 {
				cfg.SensitiveKeywords = kws
			}
		} else if kwList, ok := kwVal.([]string); ok {
			cfg.SensitiveKeywords = kwList
		}
	}

	if patVal, ok := rawMap["sensitive_patterns"]; ok {
		if patList, ok := patVal.([]any); ok {
			var pats []string
			for _, item := range patList {
				if s, ok := item.(string); ok {
					pats = append(pats, s)
				}
			}
			if len(pats) > 0 {
				cfg.SensitivePatterns = pats
			}
		} else if patList, ok := patVal.([]string); ok {
			cfg.SensitivePatterns = patList
		}
	}

	if val, ok := rawMap["allowed_chars"]; ok {
		if s, ok := val.(string); ok {
			cfg.AllowedChars = s
		}
	}

	if val, ok := rawMap["forbidden_chars"]; ok {
		if s, ok := val.(string); ok {
			cfg.ForbiddenChars = s
		}
	}

	return cfg, nil
}
