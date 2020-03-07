package emoji

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	tt := []struct {
		input    string
		expected string
	}{
		{
			input:    "I am :man_technologist: from :flag_for_turkey:. Tests are :thumbs_up:",
			expected: fmt.Sprintf("I am %v from %v. Tests are %v", ManTechnologist, FlagForTurkey, ThumbsUp),
		},
		{
			input:    "consecutive emojis :pizza::sushi::sweat:",
			expected: fmt.Sprintf("consecutive emojis %v%v%v", Pizza, Sushi, DowncastFaceWithSweat),
		},
		{
			input:    ":accordion::anguished_face: \n woman :woman_golfing:",
			expected: fmt.Sprintf("%v%v \n woman %v", Accordion, AnguishedFace, WomanGolfing),
		},
		{
			input:    "shared colon :angry_face_with_horns:anger_symbol:",
			expected: fmt.Sprintf("shared colon %vanger_symbol:", AngryFaceWithHorns),
		},
		{
			input:    ":not_exist_emoji: not exist emoji",
			expected: ":not_exist_emoji: not exist emoji",
		},
		{
			input:    ":dragon::",
			expected: fmt.Sprintf("%v:", Dragon),
		},
		{
			input:    "::+1:",
			expected: fmt.Sprintf(":%v", ThumbsUp),
		},
		{
			input:    "::anchor::",
			expected: fmt.Sprintf(":%v:", Anchor),
		},
		{
			input:    ":anguished:::",
			expected: fmt.Sprintf("%v::", AnguishedFace),
		},
		{
			input:    "dummytext",
			expected: "dummytext",
		},
	}

	for i, tc := range tt {
		got := Parse(tc.input)
		if got != tc.expected {
			t.Fatalf("test case %v fail: got: %v, expected: %v", i+1, got, tc.expected)
		}
	}
}

func BenchmarkParse(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = Parse("I am :man_technologist: from :flag_for_turkey:. Tests are :thumbs_up:")
	}
}
