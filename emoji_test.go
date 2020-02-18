package emoji

import (
	"testing"
)

func TestEmoji(t *testing.T) {
	tt := []struct {
		input    Emoji
		expected string
	}{
		{input: GrinningFace, expected: "\U0001F600"},
		{input: EyeInSpeechBubble, expected: "\U0001F441\uFE0F\u200D\U0001F5E8\uFE0F"},
		{input: ManGenie, expected: "\U0001F9DE\u200D\u2642\uFE0F"},
		{input: Badger, expected: "\U0001F9A1"},
		{input: FlagsForFlagTurkey, expected: "\U0001F1F9\U0001F1F7"},
	}

	for i, tc := range tt {
		got := tc.input.String()
		if got != tc.expected {
			t.Fatalf("test case %v fail: got: %v, expected: %v", i+1, got, tc.expected)
		}
	}
}

func TestEmojiWithTone(t *testing.T) {
	tt := []struct {
		input    EmojiWithTone
		tone     Tone
		expected string
	}{
		{input: WavingHand, tone: Tone(""), expected: "\U0001F44B"},
		{input: WavingHand, tone: Default, expected: "\U0001F44B"},
		{input: WavingHand, tone: Light, expected: "\U0001F44B\U0001F3FB"},
		{input: WavingHand, tone: MediumLight, expected: "\U0001F44B\U0001F3FC"},
		{input: WavingHand, tone: Medium, expected: "\U0001F44B\U0001F3FD"},
		{input: WavingHand, tone: MediumDark, expected: "\U0001F44B\U0001F3FE"},
		{input: WavingHand, tone: Dark, expected: "\U0001F44B\U0001F3FF"},
	}

	for i, tc := range tt {
		got := tc.input.Tone(tc.tone)
		if got != tc.expected {
			t.Fatalf("test case %v fail: got: %v, expected: %v", i+1, got, tc.expected)
		}
	}
}

func TestEmojiWithToneTwo(t *testing.T) {
	tt := []struct {
		input    EmojiWithTone
		tones    []Tone
		expected string
	}{
		{input: WomanAndManHoldingHandsWithTwoTone, tones: []Tone{}, expected: "\U0001F469\u200D\U0001F91D\u200D\U0001F468"},
		{input: WomanAndManHoldingHandsWithTwoTone, tones: []Tone{MediumLight}, expected: "\U0001F469\U0001F3FC\u200D\U0001F91D\u200D\U0001F468\U0001F3FC"},
		{input: WomanAndManHoldingHandsWithTwoTone, tones: []Tone{Medium, Dark}, expected: "\U0001F469\U0001F3FD\u200D\U0001F91D\u200D\U0001F468\U0001F3FF"},
	}

	for i, tc := range tt {
		got := tc.input.Tone(tc.tones...)
		if got != tc.expected {
			t.Fatalf("test case %v fail: got: %v, expected: %v", i+1, got, tc.expected)
		}
	}
}

func TestCountryFlag(t *testing.T) {
	tt := []struct {
		input    string
		expected Emoji
	}{
		{input: "tr", expected: FlagsForFlagTurkey},
		{input: "TR", expected: FlagsForFlagTurkey},
		{input: "us", expected: FlagsForFlagUnitedStates},
		{input: "gb", expected: FlagsForFlagUnitedKingdom},
	}

	for i, tc := range tt {
		got, err := CountryFlag(tc.input)
		if err != nil {
			t.Fatalf("test case %v fail: %v", i+1, err)
		}
		if got != tc.expected {
			t.Fatalf("test case %v fail: got: %v, expected: %v", i+1, got, tc.expected)
		}
	}
}

func TestCountryFlagError(t *testing.T) {
	tt := []struct {
		input string
		fail  bool
	}{
		{input: "tr", fail: false},
		{input: "a", fail: true},
		{input: "tur", fail: true},
	}

	for i, tc := range tt {
		_, err := CountryFlag(tc.input)
		if (err != nil) != tc.fail {
			t.Fatalf("test case %v fail: %v", i+1, err)
		}
	}
}

func BenchmarkEmoji(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = WavingHand.String()
	}
}

func BenchmarkEmojiWithTone(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = WavingHand.Tone(Medium)
	}
}

func BenchmarkEmojiWithToneTwo(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = WomanAndManHoldingHandsWithTwoTone.Tone(Medium, Dark)
	}
}

func BenchmarkCountryFlag(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = CountryFlag("tr")
	}
}
