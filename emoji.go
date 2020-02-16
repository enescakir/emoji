package emoji

import (
	"fmt"
	"html"
	"strings"
)

const (
	Default     Tone = ""
	Light       Tone = "\U0001F3FB"
	MediumLight Tone = "\U0001F3FC"
	Medium      Tone = "\U0001F3FD"
	MediumDark  Tone = "\U0001F3FE"
	Dark        Tone = "\U0001F3FF"

	unicodeFlagBaseIndex = 127397
)

type Emoji string

func (e Emoji) String() string {
	return string(e)
}

type EmojiWithTone Emoji

func (e EmojiWithTone) String() string {
	return strings.ReplaceAll(string(e), "@", string(Default))
}

func (e EmojiWithTone) Tone(tones ...Tone) EmojiWithTone {
	str := string(e)
	for _, tone := range tones {
		str = strings.Replace(str, "@", string(tone), 1)
	}

	if strings.Count(str, "@") > 0 {
		lastTone := tones[len(tones)-1]
		str = strings.ReplaceAll(str, "@", string(lastTone))
	}

	return EmojiWithTone(str)
}

type Tone string

func (t Tone) String() string {
	return string(t)
}

func CountryFlag(code string) (Emoji, error) {
	if len(code) != 2 {
		return "", fmt.Errorf("not valid country code: %q", code)
	}

	code = strings.ToUpper(code)
	flag := countryCodeLetter(code[0]) + countryCodeLetter(code[1])

	return Emoji(flag), nil
}

func countryCodeLetter(l byte) string {
	return html.UnescapeString(fmt.Sprintf("&#%v;", unicodeFlagBaseIndex+int(l)))
}
