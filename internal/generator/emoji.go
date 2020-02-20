package main

import (
	"fmt"
	"github.com/enescakir/emoji/strutil"
	"regexp"
	"strconv"
	"strings"
)

var (
	emojiRegex = regexp.MustCompile(`^(?m)(?P<code>[A-Z\d ]+[A-Z\d])\s+;\s+(fully-qualified|component)\s+#\s+.+\s+E\d+\.\d+ (?P<name>.+)$`)
	toneRegex  = regexp.MustCompile(`:\s.*tone,?`)
)

type emoji struct {
	name     string
	constant string
	code     string
	tones    []string
}

func (e *emoji) String() string {
	return fmt.Sprintf("name:%v, constant:%v, code:%v, tones: %v\n", e.name, e.constant, e.code, e.tones)
}

func newEmoji(line string) emoji {
	matches := emojiRegex.FindStringSubmatch(line)
	code := matches[1]
	name := matches[3]

	e := emoji{
		name:     name,
		constant: name,
		code:     code,
		tones:    []string{},
	}
	e.extractAttr()
	e.generateConstant()
	e.generateUnicode()

	return e
}

func (e *emoji) extractAttr() {
	parts := strings.Split(e.constant, ":")
	if len(parts) < 2 {
		// no attributes
		return
	}
	c := parts[0]
	attrs := strings.Split(parts[1], ",")
	for _, attr := range attrs {
		switch {
		case strings.Contains(attr, "tone"):
			e.tones = append(e.tones, attr)
		case strings.Contains(attr, "beard"):
			fallthrough
		case strings.Contains(attr, "hair"):
			c += " with " + attr
		case strings.HasPrefix(c, "flag"):
			c += " for " + attr
		default:
			c += " " + attr
		}
	}
	e.constant = c
}

func (e *emoji) generateConstant() {
	c := e.constant
	c = strutil.Clean(c)
	c = strings.Title(strings.ToLower(c))
	c = strutil.RemoveSpaces(c)
	e.constant = c
}

func (e *emoji) generateUnicode() {
	unicodes := []string{}
	for _, v := range strings.Split(e.code, " ") {
		u, err := strconv.ParseInt(v, 16, 32)
		if err != nil {
			panic(fmt.Errorf("unknown unicode: %v", v))
		}
		unicodes = append(unicodes, string(u))
	}
	e.code = strings.Join(unicodes, "")
}
