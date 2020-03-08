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
		{input: FlagForTurkey, expected: "\U0001F1F9\U0001F1F7"},
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

func TestEmojiWithTones(t *testing.T) {
	tt := []struct {
		input    EmojiWithTone
		tones    []Tone
		expected string
	}{
		{input: WomanAndManHoldingHands, tones: []Tone{}, expected: "\U0001f46b"},
		{input: WomanAndManHoldingHands, tones: []Tone{MediumLight}, expected: "\U0001f46b\U0001F3FC"},
		{input: WomanAndManHoldingHands, tones: []Tone{Medium, Dark}, expected: "\U0001f469\U0001F3FD\u200d\U0001f91d\u200d\U0001f468\U0001F3FF"},
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
		{input: "tr", expected: FlagForTurkey},
		{input: "TR", expected: FlagForTurkey},
		{input: "us", expected: FlagForUnitedStates},
		{input: "gb", expected: FlagForUnitedKingdom},
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

func TestNewEmojiTone(t *testing.T) {
	tt := []struct {
		input    []string
		expected EmojiWithTone
	}{
		{input: nil, expected: EmojiWithTone{}},
		{input: []string{}, expected: EmojiWithTone{}},
		{input: []string{"\U0001f64b@"}, expected: PersonRaisingHand},
		{
			input:    []string{"\U0001f46b@", "\U0001f469@\u200d\U0001f91d\u200d\U0001f468@"},
			expected: WomanAndManHoldingHands,
		},
	}

	for i, tc := range tt {
		got := newEmojiWithTone(tc.input...)
		if got != tc.expected {
			t.Fatalf("test case %v fail: got: %v, expected: %v", i+1, got, tc.expected)
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
		_ = WomanAndManHoldingHands.Tone(Medium, Dark)
	}
}

func BenchmarkCountryFlag(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = CountryFlag("tr")
	}
}
