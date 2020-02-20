package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	emojiListUrl = "https://unicode.org/Public/emoji/13.0/emoji-test.txt"
)

func main() {
	_, err := getEmojis()
	if err != nil {
		panic(err)
	}
}

func getEmojis() (map[string]map[string]map[string][]emoji, error) {
	emojis := make(map[string]map[string]map[string][]emoji)
	b, err := getFile(emojiListUrl)
	if err != nil {
		return nil, err
	}

	var group string
	var subgroup string

	parseLine := func(line string) {

		switch {
		case strings.HasPrefix(line, "# group:"):
			group = strings.TrimSpace(strings.ReplaceAll(line, "# group:", ""))
			emojis[group] = make(map[string]map[string][]emoji)
			fmt.Printf("group: %v\n", group)
		case strings.HasPrefix(line, "# subgroup:"):
			subgroup = strings.TrimSpace(strings.ReplaceAll(line, "# subgroup:", ""))
			emojis[group][subgroup] = make(map[string][]emoji)
			fmt.Printf("subgroup: %v\n", subgroup)
		case !strings.HasPrefix(line, "#") && strings.Contains(line, "fully-qualified"):
			e := newEmoji(line)
			if _, ok := emojis[group][subgroup][e.constant]; !ok {
				emojis[group][subgroup][e.constant] = []emoji{}
			}
			emojis[group][subgroup][e.constant] = append(emojis[group][subgroup][e.constant], e)
			fmt.Printf("emoji: %v\n", e)
		}
	}

	err = readLines(b, parseLine)
	if err != nil {
		return nil, err
	}

	return emojis, nil

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
