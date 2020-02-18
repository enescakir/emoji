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

// String returns string representation of the emoji.
func (e Emoji) String() string {
	return string(e)
}

// EmojiWithTone defines an emoji object that has skin tone options.
type EmojiWithTone Emoji

// String returns string representation of the emoji with default skin tone.
func (e EmojiWithTone) String() string {
	return strings.ReplaceAll(string(e), "@", Default.String())
}

// Tone returns string representation of the emoji with given skin tones.
func (e EmojiWithTone) Tone(tones ...Tone) string {
	str := string(e)

	// if no given tones, return with default skin tone
	if len(tones) == 0 {
		return e.String()
	}

	// replace tone one by one
	for _, t := range tones {
		str = strings.Replace(str, "@", t.String(), 1)
	}

	// if skin tone count is not enough, fill with last tone.
	if strings.Count(str, "@") > 0 {
		last := tones[len(tones)-1]
		str = strings.ReplaceAll(str, "@", last.String())
	}

	return str
}

// Tone defines skin tone options for emojis.
type Tone string

// String returns string representation of the skin tone.
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
	shifted := unicodeFlagBaseIndex + int(l)

	return html.UnescapeString(fmt.Sprintf("&#%v;", shifted))
}
