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
	emojis, err := getEmojis()
	if err != nil {
		panic(err)
	}

	if err = generateFile(emojis); err != nil {
		panic(err)
	}
}

func getEmojis() (*groups, error) {
	var emojis groups
	b, err := getFile(emojiListUrl)
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

func generateFile(emojis *groups) error {
	tmpl, err := template.ParseFiles("internal/generator/constants.go.tmpl")
	if err != nil {
		return err
	}

	data := struct {
		Link   string
		Date   string
		Emojis *groups
	}{
		Link:   emojiListUrl,
		Date:   time.Now().Format(time.RFC3339),
		Emojis: emojis,
	}
	var w bytes.Buffer
	if err = tmpl.Execute(&w, data); err != nil {
		return err
	}

	content, err := format.Source(w.Bytes())

	file, err := os.Create("test.go")
	if err != nil {
		return fmt.Errorf("could not create file: %v", err)
	}
	defer file.Close()

	if _, err := file.Write(content); err != nil {
		return fmt.Errorf("could not write to file: %v", err)
	}
	return nil
}

func getFile(url string) ([]byte, error) {
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
