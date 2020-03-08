package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"text/template"
	"time"
)

const (
	constantsFile = "constants.go"
	aliasesFile   = "map.go"
)

// customEmojis is the list of emojis which unicode and gemoji databases don't have.
var customEmojis = map[string]string{
	":robot_face:": "\U0001f916", // slack
}

func main() {
	emojis, err := fetchEmojis()
	if err != nil {
		panic(err)
	}

	gemojis, err := fetchGemojis()
	if err != nil {
		panic(err)
	}

	constants := generateConstants(emojis)
	aliases := generateAliases(emojis, gemojis)

	if err = save(constantsFile, emojiListURL, constants); err != nil {
		panic(err)
	}

	if err = save(aliasesFile, gemojiURL, aliases); err != nil {
		panic(err)
	}
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

func generateAliases(emojis *groups, gemojis map[string]string) string {
	var aliases []string
	var emojiMap = make(map[string]string)

	for _, grp := range emojis.Groups {
		for _, subgrp := range grp.Subgroups {
			for _, c := range subgrp.Constants {
				emoji := subgrp.Emojis[c][0]
				alias := makeAlias(snakeCase(emoji.Constant))
				aliases = append(aliases, alias)
				emojiMap[alias] = emoji.Code
			}
		}
	}

	// add gemoji aliases
	{
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

	var r string
	sort.Strings(aliases)
	for _, alias := range aliases {
		r += fmt.Sprintf("%q: %+q,\n", alias, emojiMap[alias])
	}

	return r
}
func save(filename, url, data string) error {
	tmpl, err := template.ParseFiles(fmt.Sprintf("internal/generator/%v.tmpl", filename))
	if err != nil {
		return err
	}

	d := struct {
		Link string
		Date string
		Data string
	}{
		Link: url,
		Date: time.Now().Format(time.RFC3339),
		Data: data,
	}

	var w bytes.Buffer
	if err = tmpl.Execute(&w, d); err != nil {
		return err
	}

	content, err := format.Source(w.Bytes())
	if err != nil {
		return fmt.Errorf("could not format file: %v", err)
	}

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
