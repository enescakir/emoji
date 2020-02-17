# emoji :rocket: :school_satchel: :tada:
`emoji` is a minimalistic emoji library for Go. It lets you use emoji characters in strings.

Inspired by [spatie/emoji](https://github.com/spatie/emoji)

## Install :floppy_disk:
``` bash
go get github.com/enescakir/emoji
```

## Usage :surfer:
```go
package main

import (
    "fmt"

    "github.com/enescakir/emoji"
)

func main() {
    fmt.Printf("hello %v from %v\n", 
        emoji.WavingHand, 
        emoji.FlagsForFlagTurkey,
    )
    fmt.Printf("different skin tones. default: %v light: %v dark: %v\n", 
        emoji.ThumbsUp,
        emoji.OkHand.Tone(emoji.Light),
        emoji.CallMeHand.Tone(emoji.Dark),
    )
    fmt.Printf("emoji with multiple skins: %v\n", 
        emoji.PeopleHoldingHands.Tone(emoji.Light, emoji.Dark),
    )
}

/* OUTPUT

    hello 👋 from 🇹🇷
    different skin tones. default: 👍 light: 👌🏻 dark: 🤙🏿
    emoji with multiple skins: 🧑🏻‍🤝‍🧑🏿

*/
```

This package contains Full Emoji List v12.0 based on [https://unicode.org/Public/emoji/12.0/emoji-test.txt](https://unicode.org/Public/emoji/12.0/emoji-test.txt).

Also, you can generate country flag emoji with [ISO 3166 Alpha2](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2) codes:
```go
emoji.CountryFlag("tr") // 🇹🇷
```

## Testing :hammer:
``` bash
go test
```

## Todo :pushpin:
* Add `godoc`
* Add badges to README
* Add benchmarks
* Add emoji constant generator

## Contributing :man_technologist:
I am accepting PRs that add characters to the class.
Please use [this list](http://unicode.org/emoji/charts/full-emoji-list.html) to look up the unicode value and the name of the character.

## Credits :star:
- [Enes Çakır](https://github.com/enescakir)

## License :scroll:
The MIT License (MIT). Please see [License File](LICENSE.md) for more information.
