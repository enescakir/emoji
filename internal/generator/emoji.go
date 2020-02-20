package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	emojipkg "github.com/enescakir/emoji"
	"github.com/enescakir/emoji/internal/strutil"
)

var (
	emojiRegex = regexp.MustCompile(`^(?m)(?P<code>[A-Z\d ]+[A-Z\d])\s+;\s+(fully-qualified|component)\s+#\s+.+\s+E\d+\.\d+ (?P<name>.+)$`)
	toneRegex  = regexp.MustCompile(`:\s.*tone,?`)
)

type groups struct {
	Groups []*group
}

func (g *groups) Template() string {
	r := ""
	for _, grp := range g.Groups {
		r += grp.Template()
	}

	return r
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

func (g *group) Template() string {
	r := fmt.Sprintf("\n// GROUP: %v\n", g.Name)
	for _, subs := range g.Subgroups {
		r += subs.Template()
	}

	return r
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

func (s *subgroup) Template() string {
	r := fmt.Sprintf("// SUBGROUP: %v\n", s.Name)
	for _, c := range s.Constants {
		if len(s.Emojis[c]) > 1 {
			for _, e := range s.Emojis[c] {
				fmt.Printf("%+q | %v | %v\n", replaceTones(e.Code), e.Name, e.Tones)
			}
			fmt.Println()
		}

		r += s.Emojis[c][0].Template()
	}

	return r
}

func replaceTones(str string) string {
	for _, tone := range []emojipkg.Tone{emojipkg.Light, emojipkg.MediumLight, emojipkg.Medium, emojipkg.MediumDark, emojipkg.Dark} {
		str = strings.ReplaceAll(str, tone.String(), "@")
	}

	return str
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

func (e *emoji) Template() string {
	return fmt.Sprintf("%s Emoji = %+q // %s\n", e.Constant, e.Code, e.Name)
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
	e.generateConstant()
	e.generateUnicode()

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

func (e *emoji) generateConstant() {
	c := e.Constant
	c = strutil.Clean(c)
	c = strings.Title(strings.ToLower(c))
	c = strutil.RemoveSpaces(c)
	e.Constant = c
}

func (e *emoji) generateUnicode() {
	unicodes := []string{}
	for _, v := range strings.Split(e.Code, " ") {
		u, err := strconv.ParseInt(v, 16, 32)
		if err != nil {
			panic(fmt.Errorf("unknown unicode: %v", v))
		}
		unicodes = append(unicodes, string(u))
	}
	e.Code = strings.Join(unicodes, "")
}
