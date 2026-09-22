package jmap

import (
	"fmt"
	"regexp"
	"strings"
)

var level1Var = regexp.MustCompile(`^\{[A-Za-z0-9_.%-]+\}$`)

// ValidateURITemplateLevel1 rejects any expression other than a single
// RFC 6570 Level 1 {varname} (no operators, modifiers, or lists).
func ValidateURITemplateLevel1(tmpl string) error {
	for i := 0; i < len(tmpl); {
		open := strings.IndexByte(tmpl[i:], '{')
		if open < 0 {
			return nil
		}
		open += i
		close := strings.IndexByte(tmpl[open:], '}')
		if close < 0 {
			return fmt.Errorf("jmap: unterminated expression in %q", tmpl)
		}
		close += open
		if !level1Var.MatchString(tmpl[open : close+1]) {
			return fmt.Errorf("jmap: %q is not an RFC 6570 Level 1 expression", tmpl[open:close+1])
		}
		i = close + 1
	}
	return nil
}

// ExpandURITemplateLevel1 substitutes {name} variables using RFC 6570 §3.2.2
// simple string expansion: every octet outside the unreserved set
// (ALPHA / DIGIT / "-" / "." / "_" / "~") is percent-encoded.
func ExpandURITemplateLevel1(tmpl string, vars map[string]string) string {
	var b strings.Builder
	b.Grow(len(tmpl))
	for i := 0; i < len(tmpl); {
		if tmpl[i] == '{' {
			end := strings.IndexByte(tmpl[i:], '}')
			if end < 0 {
				b.WriteString(tmpl[i:])
				break
			}
			name := tmpl[i+1 : i+end]
			b.WriteString(encodeURITemplateValue(vars[name]))
			i += end + 1
			continue
		}
		b.WriteByte(tmpl[i])
		i++
	}
	return b.String()
}

func encodeURITemplateValue(s string) string {
	const hex = "0123456789ABCDEF"
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if isUnreservedByte(c) {
			b.WriteByte(c)
		} else {
			b.WriteByte('%')
			b.WriteByte(hex[c>>4])
			b.WriteByte(hex[c&0x0F])
		}
	}
	return b.String()
}

func isUnreservedByte(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') ||
		c == '-' || c == '.' || c == '_' || c == '~'
}
