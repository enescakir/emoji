package emoji

import (
	"fmt"
	"testing"
)

func TestColor(t *testing.T) {
	fmt.Printf("hello %v from %v\n", WavingHand, FlagsForFlagTurkey)
	fmt.Printf("different skin tones. default: %v light: %v dark: %v\n", ThumbsUp, OkHand.Tone(Light), CallMeHand.Tone(Dark))
	fmt.Printf("emoji with multiple skins: %v\n", PeopleHoldingHands.Tone(Light, Dark))

}
