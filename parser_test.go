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
			expected: "",
		},
		{
			input:    "I love eating :pizza::sushi::sweat:",
			expected: "",
		},
		{
			input:    " :pizza::sushi: \n mann :woman_golfing:",
			expected: "",
		},
		{
			input:    ":pizza:sushi: mann :sweat:",
			expected: "",
		},
		{
			input:    ":piasfaf: not exist  :sweat:",
			expected: "",
		},
		{
			input:    ":pizza::",
			expected: "",
		},
		{
			input:    ":pizza:::",
			expected: "",
		},
		{
			input:    "afsaff",
			expected: "",
		},
	}

	for _, tc := range tt {
		fmt.Println(Parse(tc.input))
		//got := tc.input
		//if got != tc.expected {
		//	t.Fatalf("test case %v fail: got: %v, expected: %v", i+1, got, tc.expected)
		//}
	}
}

func BenchmarkParse(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = Parse("I am :man_technologist: from :flag_for_turkey:. Tests are :thumbs_up:")
	}
}
