package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	emojipkg "github.com/enescakir/emoji"
)

var (
	emojiRegex = regexp.MustCompile(`^(?m)(?P<code>[A-Z\d ]+[A-Z\d])\s+;\s+(fully-qualified|component)\s+#\s+.+\s+E\d+\.\d+ (?P<name>.+)$`)
	toneRegex  = regexp.MustCompile(`:\s.*tone,?`)
)

type groups struct {
	Groups []*group
}

func (g *groups) Append(grpName string) *group {
	// fmt.Printf("group: %v\n", grpName)
	grp := group{Name: grpName}
	g.Groups = append(g.Groups, &grp)

	return &grp
}

type group struct {
	Name      string
	Subgroups []*subgroup
}

func (g *group) Append(subgrpName string) *subgroup {
	// fmt.Printf("subgroup: %v\n", subgrpName)
	subgrp := subgroup{Name: subgrpName, Emojis: make(map[string][]emoji)}
	g.Subgroups = append(g.Subgroups, &subgrp)

	return &subgrp
}

type subgroup struct {
	Name      string
	Emojis    map[string][]emoji
	Constants []string
}

func (s *subgroup) Append(e emoji) {
	// fmt.Printf("emoji: %v\n", e)
	if _, ok := s.Emojis[e.Constant]; ok {
		s.Emojis[e.Constant] = append(s.Emojis[e.Constant], e)
	} else {
		s.Emojis[e.Constant] = []emoji{e}
		s.Constants = append(s.Constants, e.Constant)
	}
}

type emoji struct {
	Name     string
	Constant string
	Code     string
	Tones    []string
}

func (e *emoji) String() string {
	return fmt.Sprintf("name:%v, constant:%v, code:%v, tones: %v\n", e.Name, e.Constant, e.Code, e.Tones)
}

func newEmoji(line string) *emoji {
	matches := emojiRegex.FindStringSubmatch(line)
	if len(matches) < 4 {
		return nil
	}
	code := matches[1]
	name := matches[3]

	e := emoji{
		Name:     name,
		Constant: name,
		Code:     code,
		Tones:    []string{},
	}
	e.extractAttr()
	e.Constant = generateConstant(e.Constant)
	e.Code = generateUnicode(e.Code)

	return &e
}

func (e *emoji) extractAttr() {
	parts := strings.Split(e.Constant, ":")
	if len(parts) < 2 {
		// no attributes
		return
	}
	c := parts[0]
	attrs := strings.Split(parts[1], ",")
	for _, attr := range attrs {
		switch {
		case strings.Contains(attr, "tone"):
			e.Tones = append(e.Tones, attr)
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
	e.Constant = c
}

func generateConstant(c string) string {
	c = clean(c)
	c = strings.Title(strings.ToLower(c))
	c = removeSpaces(c)

	return c
}

func generateUnicode(code string) string {
	unicodes := []string{}
	for _, v := range strings.Split(code, " ") {
		u, err := strconv.ParseInt(v, 16, 32)
		if err != nil {
			panic(fmt.Errorf("unknown unicode: %v", v))
		}
		unicodes = append(unicodes, string(u))
	}

	return strings.Join(unicodes, "")
}

func defaultTone(basic, toned string) string {
	toneInd := strings.IndexRune(toned, []rune(emojipkg.TonePlaceholder)[0])
	for i, ch := range basic {
		if i != toneInd {
			continue
		}
		if ch == '\ufe0f' {
			return "\ufe0f"
		}
		break
	}

	return ""
}

func replaceTones(code string) string {
	for _, tone := range []emojipkg.Tone{
		emojipkg.Light,
		emojipkg.MediumLight,
		emojipkg.Medium,
		emojipkg.MediumDark,
		emojipkg.Dark,
	} {
		code = strings.ReplaceAll(code, tone.String(), emojipkg.TonePlaceholder)
	}

	return code
}
