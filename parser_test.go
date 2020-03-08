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
			input:    "too many colon::::closed_book:::: too many colon:",
			expected: fmt.Sprintf("too many colon:::%v::: too many colon:", ClosedBook),
		},
		{
			input:    "emoji with space :angry face_with_horns:anger_symbol:",
			expected: fmt.Sprintf("emoji with space :angry face_with_horns%v", AngerSymbol),
		},
		{
			input:    "flag testing :flag-tr: done",
			expected: fmt.Sprintf("flag testing %v done", FlagForTurkey),
		},
		{
			input:    "not valid flags :flag-tra: :flag-t: testing",
			expected: fmt.Sprintf("not valid flags :flag-tra: :flag-t: testing"),
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

func TestMap(t *testing.T) {
	expected := len(emojiMap)
	got := len(Map())

	if got != expected {
		t.Fatalf("test case fail: got: %v, expected: %v", got, expected)
	}
}

func TestAppendAlias(t *testing.T) {
	tt := []struct {
		alias string
		code  string
		err   bool
	}{
		{alias: ":my_car:", code: "\U0001f3ce\ufe0f", err: false},
		{alias: ":berserker:", code: "\U0001f621", err: false},
		{alias: ":potato:", code: "\U0001f423", err: true},
		{alias: ":not_valid alias:", code: "\U0001f423", err: true},
	}

	for i, tc := range tt {
		err := AppendAlias(tc.alias, tc.code)
		if (err != nil) != tc.err {
			t.Fatalf("test case %v fail: got: %v, expected: %v", i+1, err, tc.err)
		}

		if exist := Exist(tc.alias); !exist && !tc.err {
			t.Fatalf("test case %v fail: got: %v, expected: %v", i+1, !exist, exist)
		}
	}
}

func TestExist(t *testing.T) {
	tt := []struct {
		input    string
		expected bool
	}{
		{input: ":man_technologist:", expected: true},
		{input: ":registered:", expected: true},
		{input: ":robot_face:", expected: true},
		{input: ":wave:", expected: true},
		{input: ":sheaf_of_rice:", expected: true},
		{input: ":random_emoji:", expected: false},
		{input: ":test_emoji:", expected: false},
	}

	for i, tc := range tt {
		got := Exist(tc.input)
		if got != tc.expected {
			t.Fatalf("test case %v fail: got: %v, expected: %v", i+1, got, tc.expected)
		}
	}
}

func TestFind(t *testing.T) {
	tt := []struct {
		input    string
		expected string
		exist    bool
	}{
		{input: ":man_technologist:", expected: ManTechnologist.String(), exist: true},
		{input: ":robot_face:", expected: Robot.String(), exist: true},
		{input: ":wave:", expected: WavingHand.String(), exist: true},
		{input: ":sheaf_of_rice:", expected: SheafOfRice.String(), exist: true},
		{input: ":random_emoji:", expected: "", exist: false},
		{input: ":test_emoji:", expected: "", exist: false},
	}

	for i, tc := range tt {
		got, exist := Find(tc.input)
		if got != tc.expected {
			t.Fatalf("test case %v fail: got: %v, expected: %v", i+1, got, tc.expected)
		}

		if exist != tc.exist {
			t.Fatalf("test case %v fail: got: %v, expected: %v", i+1, exist, tc.exist)
		}
	}
}

func BenchmarkParse(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = Parse("I am :man_technologist: from :flag_for_turkey:. Tests are :thumbs_up:")
	}
}
