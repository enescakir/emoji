package emoji

import (
	"fmt"
	"strings"
	"unicode"
)

// Parse replaces emoji aliases (:pizza:) with unicode representation.
func Parse(input string) string {
	var matched strings.Builder
	var output strings.Builder

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
				output.WriteString(matched.String())
				matched.Reset()
			}
			continue
		}

		// r is `:` now
		// if matched is empty, it's the beginning of the emoji alias
		if matched.Len() == 0 {
			matched.WriteRune(r)
			continue
		}

		// it's the end of the emoji alias
		alias := matched.String() + ":"

		code, ok := emojiMap[alias]
		if ok {
			output.WriteString(code)
			matched.Reset()
		} else {
			// TODO: check for country codes: `flag-[a-z]{2}`
			output.WriteString(matched.String())
			// it might be the beginning of the another emoji alias
			matched.Reset()
			matched.WriteRune(r)
		}
	}

	// if matched not empty, add it to output
	if matched.Len() != 0 {
		output.WriteString(matched.String())
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
	_, ok := emojiMap[alias]

	return ok
}

// Find returns the emoji code by alias.
func Find(alias string) (string, bool) {
	code, ok := emojiMap[alias]
	if !ok {
		return "", false
	}

	return code, true
}
