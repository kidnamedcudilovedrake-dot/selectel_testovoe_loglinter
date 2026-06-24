package analyzer

import (
	"strings"
	"unicode"
)

func isLowerStart(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) {
			return unicode.IsLower(r)
		}
	}
	return true
}

func toLowerStart(s string) string {
	runes := []rune(s)
	for i, r := range runes {
		if unicode.IsLetter(r) {
			if unicode.IsUpper(r) {
				runes[i] = unicode.ToLower(r)
			}
			break
		}
	}
	return string(runes)
}

func isEnglish(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) {
			if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')) {
				return false
			}
		}
	}
	return true
}

func isAllowedLogChar(r rune) bool {
	if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
		return true
	}
	if r >= '0' && r <= '9' {
		return true
	}
	switch r {
	case ' ', '\t', '-', '_', '=', ':', '.', ',', '/', '\'', '"', '(', ')', '[', ']', '{', '}', '*', '&', '%', '#', '@', '+', ';', '<', '>':
		return true
	}
	return false
}

func badChars(s string, forbidden string) []string {
	var errs []string

	hasEmoji := false
	for _, r := range s {
		if r > 0x7f && (unicode.IsSymbol(r) || unicode.IsPunct(r)) {
			hasEmoji = true
			break
		}
	}
	if hasEmoji {
		errs = append(errs, "contains emoji")
	}

	var found []rune
	for _, r := range s {
		if unicode.IsLetter(r) {
			continue
		}

		if strings.ContainsRune(forbidden, r) {
			dup := false
			for _, f := range found {
				if f == r {
					dup = true
					break
				}
			}
			if !dup {
				found = append(found, r)
			}
			continue
		}

		if !isAllowedLogChar(r) {
			dup := false
			for _, f := range found {
				if f == r {
					dup = true
					break
				}
			}
			if !dup {
				found = append(found, r)
			}
		}
	}

	for _, r := range found {
		if r > 0x7f && (unicode.IsSymbol(r) || unicode.IsPunct(r)) {
			continue
		}
		errs = append(errs, "contains forbidden character '"+string(r)+"'")
	}

	if strings.Contains(s, "..") {
		errs = append(errs, "contains ellipsis or consecutive dots")
	} else if strings.HasSuffix(s, ".") {
		errs = append(errs, "contains trailing dot")
	}

	lower := strings.ToLower(s)
	prefixes := []string{"info:", "warn:", "warning:", "error:", "debug:", "fatal:", "panic:"}
	for _, p := range prefixes {
		if strings.HasPrefix(lower, p) {
			errs = append(errs, "contains redundant level prefix '"+p+"'")
			break
		}
	}

	return errs
}

func cleanChars(s string, forbidden string) string {
	lower := strings.ToLower(s)
	prefixes := []string{"info:", "warn:", "warning:", "error:", "debug:", "fatal:", "panic:"}
	for _, p := range prefixes {
		if strings.HasPrefix(lower, p) {
			s = s[len(p):]
			s = strings.TrimSpace(s)
			break
		}
	}

	var b strings.Builder
	for _, r := range s {
		if r > 0x7f && (unicode.IsSymbol(r) || unicode.IsPunct(r)) {
			continue
		}
		if strings.ContainsRune(forbidden, r) || !isAllowedLogChar(r) {
			continue
		}
		b.WriteRune(r)
	}
	res := b.String()

	res = strings.TrimRight(res, ".")
	return strings.TrimSpace(res)
}
