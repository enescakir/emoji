# emoji :wolf: :evergreen_tree: :school_satchel:
`emoji` is a minimalistic emoji library for Go. It lets you use emoji characters in strings.

Inspired by [spatie/emoji](https://github.com/spatie/emoji)

## Install
``` bash
go get github.com/enescakir/emoji
```

## Usage
```go
    package main
    
    import (
        "fmt"
    
        "github.com/enescakir/emoji"
    )
    
    func main() {
        fmt.Printf("hello %v from %v\n", emoji.WavingHand, emoji.FlagsForFlagTurkey)
        fmt.Printf("different skin tones. default: %v light: %v dark: %v\n", 
            emoji.ThumbsUp,
            emoji.OkHand.Tone(emoji.Light),
            emoji.CallMeHand.Tone(emoji.Dark),
        )
        fmt.Printf("emoji with multiple skins: %v\n", emoji.PeopleHoldingHands.Tone(emoji.Light, emoji.Dark))
    }

    /* OUTPUT

        hello 👋 from 🇹🇷
        different skin tones. default: 👍 light: 👌🏻 dark: 🤙🏿
        emoji with multiple skins: 🧑🏻‍🤝‍🧑🏿

    */
    
```

This package contains Full Emoji List v12.0 based on [https://unicode.org/Public/emoji/12.0/emoji-test.txt](https://unicode.org/Public/emoji/12.0/emoji-test.txt).

## Testing
``` bash
go test
```

## Todo
* Add `godoc`
* Add country code to flag emoji converter
* Add badges to README
* Add tests
* Add emoji constant generator

## Contributing
I am accepting PRs that add characters to the class.
Please use [this list](http://unicode.org/emoji/charts/full-emoji-list.html) to look up the unicode value and the name of the character.

## Credits
- [Enes Çakır](https://github.com/enescakir)

## License
The MIT License (MIT). Please see [License File](LICENSE.md) for more information.