package emoji

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unsafe"
)

var (
	flagRegex = regexp.MustCompile(`^:flag-([a-zA-Z]{2}):$`)
)

type Parser struct {
	matched bytes.Buffer
}

func NewParser() *Parser {
	return &Parser{matched: bytes.Buffer{}}
}

// Parse replaces emoji aliases (:pizza:) with unicode representation.
func (p *Parser) Parse(input string) string {
	p.matched.Reset()
	return parseInternal(input, &p.matched)
}

// Parse replaces emoji aliases (:pizza:) with unicode representation.
func Parse(input string) string {
	return parseInternal(input, &bytes.Buffer{})
}

// Parse replaces emoji aliases (:pizza:) with unicode representation.
func parseInternal(input string, matched *bytes.Buffer) string {
	var output strings.Builder
	output.Grow(len(input))

	for _, r := range input {
		// when it's not `:`, it might be inner or outer of the emoji alias
		if r != ':' {
			// if matched is empty, it's the outer of the emoji alias
			if matched.Len() == 0 {
				output.WriteRune(r)
				continue
			}

			matched.WriteRune(r)

			// if it's space, the alias's not valid.
			// reset matched for breaking the emoji alias
			if unicode.IsSpace(r) {
				output.WriteString(unsafeString(matched))
				matched.Reset()
			}
			continue
		}

		// r is `:` now
		// if matched is empty, it's the beginning of the emoji alias
		if matched.Len() == 0 {
			matched.WriteByte(':')
			continue
		}

		// it's the end of the emoji alias
		match := unsafeString(matched)
		matched.WriteByte(':')
		alias := unsafeString(matched)

		// check for emoji alias
		if code, ok := Find(alias); ok {
			output.WriteString(code)
			matched.Reset()
			continue
		}

		// not found any emoji
		output.WriteString(match)
		// it might be the beginning of the another emoji alias
		matched.Reset()
		matched.WriteByte(':')
	}

	// if matched not empty, add it to output
	if matched.Len() != 0 {
		output.WriteString(unsafeString(matched))
		matched.Reset()
	}

	return output.String()
}

// Map returns the emojis map.
// Key is the alias of the emoji.
// Value is the code of the emoji.
func Map() map[string]string {
	return emojiMap
}

// AppendAlias adds new emoji pair to the emojis map.
func AppendAlias(alias, code string) error {
	if c, ok := emojiMap[alias]; ok {
		return fmt.Errorf("emoji already exist: %q => %+q", alias, c)
	}

	for _, r := range alias {
		if unicode.IsSpace(r) {
			return fmt.Errorf("emoji alias is not valid: %q", alias)
		}
	}

	emojiMap[alias] = code

	return nil
}

// Exist checks existence of the emoji by alias.
func Exist(alias string) bool {
	_, ok := Find(alias)

	return ok
}

// Find returns the emoji code by alias.
func Find(alias string) (string, bool) {
	if code, ok := emojiMap[alias]; ok {
		return code, true
	}

	if flag := checkFlag(alias); len(flag) > 0 {
		return flag, true
	}

	return "", false
}

// checkFlag finds flag emoji for `flag-[CODE]` pattern
func checkFlag(alias string) string {
	if matches := flagRegex.FindStringSubmatch(alias); len(matches) == 2 {
		flag, _ := CountryFlag(matches[1])

		return flag.String()
	}

	return ""
}

func unsafeString(matched *bytes.Buffer) string {
	buf := matched.Bytes()
	return *(*string)(unsafe.Pointer(&buf))
}
