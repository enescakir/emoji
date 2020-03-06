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
	"strings"
	"text/template"
	"time"
)

const (
	emojiListUrl = "https://unicode.org/Public/emoji/13.0/emoji-test.txt"
)

func main() {
	emojis, err := fetch()
	if err != nil {
		panic(err)
	}

	constants := generate(emojis)

	if err = save(constants); err != nil {
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

func generate(emojis *groups) string {
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
			return fmt.Sprintf("%s EmojiWithTone = newEmojiWithTone(%+q).withDefaultTone(%+q) // %s\n", basic.Constant, oneTonedCode, defaultTone, basic.Name)
		}

		return fmt.Sprintf("%s EmojiWithTone = newEmojiWithTone(%+q) // %s\n", basic.Constant, oneTonedCode, basic.Name)
	case 26:
		oneTonedCode := replaceTones(emojis[1].Code)
		twoTonedCode := replaceTones(emojis[2].Code)

		return fmt.Sprintf("%s EmojiWithTone = newEmojiWithTone(%+q, %+q) // %s\n", basic.Constant, oneTonedCode, twoTonedCode, basic.Name)
	default:
		panic(fmt.Errorf("not expected emoji count for a constant: %v", len(emojis)))
	}
}

func save(constants string) error {
	tmpl, err := template.ParseFiles("internal/generator/constants.go.tmpl")
	if err != nil {
		return err
	}

	data := struct {
		Link      string
		Date      string
		Constants string
	}{
		Link:      emojiListUrl,
		Date:      time.Now().Format(time.RFC3339),
		Constants: constants,
	}
	var w bytes.Buffer
	if err = tmpl.Execute(&w, data); err != nil {
		return err
	}

	content, err := format.Source(w.Bytes())

	file, err := os.Create("constants.go")
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
