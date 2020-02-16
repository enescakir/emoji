package emoji

import (
	"fmt"
	"html"
	"strings"
)

// Base attributes
const (
	unicodeFlagBaseIndex = 127397
)

// Skin tone colors
const (
	Default     Tone = ""
	Light       Tone = "\U0001F3FB"
	MediumLight Tone = "\U0001F3FC"
	Medium      Tone = "\U0001F3FD"
	MediumDark  Tone = "\U0001F3FE"
	Dark        Tone = "\U0001F3FF"
)

// Emoji defines an emoji object.
type Emoji string

func (e Emoji) String() string {
	return string(e)
}

// EmojiWithTone defines an emoji object that has skin tone options.
type EmojiWithTone Emoji

func (e EmojiWithTone) String() string {
	return strings.ReplaceAll(string(e), "@", string(Default))
}

// Tone returns an emoji object with given skin tone.
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

// Tone defines skin tone options for emojis.
type Tone string

func (t Tone) String() string {
	return string(t)
}

// CountryFlag returns a country flag emoji from given country code.
// Full list of country codes: https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2
func CountryFlag(code string) (Emoji, error) {
	if len(code) != 2 {
		return "", fmt.Errorf("not valid country code: %q", code)
	}

	code = strings.ToUpper(code)
	flag := countryCodeLetter(code[0]) + countryCodeLetter(code[1])

	return Emoji(flag), nil
}

// countryCodeLetter shifts given letter byte as unicodeFlagBaseIndex and changes encoding
func countryCodeLetter(l byte) string {
	return html.UnescapeString(fmt.Sprintf("&#%v;", unicodeFlagBaseIndex+int(l)))
}
