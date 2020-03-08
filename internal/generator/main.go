package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
	"text/template"
	"time"
	"unicode"
)

const (
	emojiListUrl  = "https://unicode.org/Public/emoji/13.0/emoji-test.txt"
	gemojiURL     = "https://raw.githubusercontent.com/github/gemoji/master/db/emoji.json"
	constantsFile = "constants.go"
	aliasesFile   = "map.go"
)

// unicode and gemoji databases don't have alias like that
var customEmojis = map[string]string{
	":robot_face:": "\U0001f916", // slack
}

func main() {
	emojis, err := fetch()
	if err != nil {
		panic(err)
	}

	constants := generateConstants(emojis)
	aliases := generateAliases(emojis)

	if err = save(constantsFile, constants); err != nil {
		panic(err)
	}

	if err = save(aliasesFile, aliases); err != nil {
		panic(err)
	}
}

func fetch() (*groups, error) {
	var emojis groups
	b, err := fetchData(emojiListUrl)
	if err != nil {
		return nil, err
	}

	var grp *group
	var subgrp *subgroup

	parseLine := func(line string) {
		switch {
		case strings.HasPrefix(line, "# group:"):
			name := strings.TrimSpace(strings.ReplaceAll(line, "# group:", ""))
			grp = emojis.Append(name)
		case strings.HasPrefix(line, "# subgroup:"):
			name := strings.TrimSpace(strings.ReplaceAll(line, "# subgroup:", ""))
			subgrp = grp.Append(name)
		case !strings.HasPrefix(line, "#"):
			if e := newEmoji(line); e != nil {
				subgrp.Append(*e)
			}
		}
	}

	if err = readLines(b, parseLine); err != nil {
		return nil, err
	}

	return &emojis, nil
}

func generateAliases(emojis *groups) string {
	var aliases []string
	var emojiMap = make(map[string]string)

	for _, grp := range emojis.Groups {
		for _, subgrp := range grp.Subgroups {
			for _, c := range subgrp.Constants {
				emoji := subgrp.Emojis[c][0]
				alias := ":" + snakeCase(emoji.Constant) + ":"
				aliases = append(aliases, alias)
				emojiMap[alias] = emoji.Code
			}
		}
	}

	// add gemoji aliases
	{
		gemojis, err := fetchGemoji()
		if err != nil {
			panic(err)
		}

		for alias, code := range gemojis {
			_, ok := emojiMap[alias]
			if !ok {
				aliases = append(aliases, alias)
			}
			emojiMap[alias] = code
		}
	}

	// add custom emoji aliases
	{
		for alias, code := range customEmojis {
			_, ok := emojiMap[alias]
			if !ok {
				aliases = append(aliases, alias)
			}
			emojiMap[alias] = code
		}
	}

	var res string
	sort.Strings(aliases)
	for _, alias := range aliases {
		res += fmt.Sprintf("%q: %+q,\n", alias, emojiMap[alias])
	}

	return res
}

func snakeCase(str string) string {
	var output strings.Builder
	for i, r := range str {
		switch {
		case unicode.IsUpper(r):
			if i != 0 {
				output.WriteRune('_')
			}
			output.WriteRune(unicode.ToLower(r))
		case unicode.IsDigit(r):
			if i != 0 && !unicode.IsDigit(rune(str[i-1])) {
				output.WriteRune('_')
			}
			output.WriteRune(r)
		default:
			output.WriteRune(r)
		}
	}

	return output.String()
}

type gemoji struct {
	Emoji   string   `json:"emoji"`
	Aliases []string `json:"aliases"`
}

func fetchGemoji() (map[string]string, error) {
	b, err := fetchData(gemojiURL)
	if err != nil {
		return nil, err
	}

	var gemojis []gemoji
	r := make(map[string]string)

	if err = json.Unmarshal(b, &gemojis); err != nil {
		return nil, err
	}

	for _, gemoji := range gemojis {
		for _, alias := range gemoji.Aliases {
			if len(alias) == 0 || len(gemoji.Emoji) == 0 {
				continue
			}

			r[":"+alias+":"] = gemoji.Emoji
		}
	}

	return r, nil
}

func generateConstants(emojis *groups) string {
	var res string
	for _, grp := range emojis.Groups {
		res += fmt.Sprintf("\n// GROUP: %v\n", grp.Name)
		for _, subgrp := range grp.Subgroups {
			res += fmt.Sprintf("// SUBGROUP: %v\n", subgrp.Name)
			for _, c := range subgrp.Constants {
				res += emojiConstant(subgrp.Emojis[c])
			}
		}
	}

	return res
}

func emojiConstant(emojis []emoji) string {
	basic := emojis[0]
	switch len(emojis) {
	case 1:
		return fmt.Sprintf("%s Emoji = %+q // %s\n", basic.Constant, basic.Code, basic.Name)
	case 6:
		oneTonedCode := replaceTones(emojis[1].Code)
		defaultTone := defaultTone(basic.Code, oneTonedCode)

		if defaultTone != "" {
			return fmt.Sprintf("%s EmojiWithTone = newEmojiWithTone(%+q).withDefaultTone(%+q) // %s\n",
				basic.Constant, oneTonedCode, defaultTone, basic.Name)
		}

		return fmt.Sprintf("%s EmojiWithTone = newEmojiWithTone(%+q) // %s\n",
			basic.Constant, oneTonedCode, basic.Name)
	case 26:
		oneTonedCode := replaceTones(emojis[1].Code)
		twoTonedCode := replaceTones(emojis[2].Code)

		return fmt.Sprintf("%s EmojiWithTone = newEmojiWithTone(%+q, %+q) // %s\n",
			basic.Constant, oneTonedCode, twoTonedCode, basic.Name)
	default:
		panic(fmt.Errorf("not expected emoji count for a constant: %v", len(emojis)))
	}
}

func save(filename, data string) error {
	tmpl, err := template.ParseFiles(fmt.Sprintf("internal/generator/%v.tmpl", filename))
	if err != nil {
		return err
	}

	d := struct {
		Link string
		Date string
		Data string
	}{
		Link: emojiListUrl,
		Date: time.Now().Format(time.RFC3339),
		Data: data,
	}

	var w bytes.Buffer
	if err = tmpl.Execute(&w, d); err != nil {
		return err
	}

	content, err := format.Source(w.Bytes())

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create file: %v", err)
	}
	defer file.Close()

	if _, err := file.Write(content); err != nil {
		return fmt.Errorf("could not write to file: %v", err)
	}
	return nil
}

func fetchData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func readLines(b []byte, fn func(string)) error {
	reader := bufio.NewReader(bytes.NewReader(b))

	var line string
	var err error
	for {
		line, err = reader.ReadString('\n')
		if err != nil {
			break
		}

		fn(line)
	}

	if err != io.EOF {
		return err
	}

	return nil
}
