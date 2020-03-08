# emoji :rocket: :school_satchel: :tada:
[![Build Status](https://github.com/enescakir/emoji/workflows/build/badge.svg?branch=master)](https://github.com/enescakir/emoji/actions)
[![godoc](https://godoc.org/github.com/enescakir/emoji?status.svg)](https://godoc.org/github.com/enescakir/emoji)
[![Go Report Card](https://goreportcard.com/badge/github.com/enescakir/emoji)](https://goreportcard.com/report/github.com/enescakir/emoji)
[![Codecov](https://img.shields.io/codecov/c/github/enescakir/emoji)](https://codecov.io/gh/enescakir/emoji)
[![MIT License](https://img.shields.io/github/license/enescakir/emoji)](https://github.com/enescakir/emoji/blob/master/LICENSE)

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
	fmt.Printf("Hello %v\n", emoji.WavingHand)
	fmt.Printf("I am %v from %v\n",
		emoji.ManTechnologist,
		emoji.FlagForTurkey,
	)
	fmt.Printf("Different skin tones.\n  default: %v light: %v dark: %v\n",
		emoji.ThumbsUp,
		emoji.OkHand.Tone(emoji.Light),
		emoji.CallMeHand.Tone(emoji.Dark),
	)
	fmt.Printf("Emojis with multiple skin tones.\n  both medium: %v light and dark: %v\n",
		emoji.PeopleHoldingHands.Tone(emoji.Medium),
		emoji.PeopleHoldingHands.Tone(emoji.Light, emoji.Dark),
	)
	fmt.Println(emoji.Parse("Emoji aliases are :sunglasses:"))
	emoji.Println("Use fmt wrappers :+1: with emoji support :tada:")
}

/* OUTPUT

    Hello 👋
    I am 👨‍💻 from 🇹🇷
    Different skin tones.
      default: 👍 light: 👌🏻 dark: 🤙🏿
    Emojis with multiple skin tones.
      both medium: 🧑🏽‍🤝‍🧑🏽 light and dark: 🧑🏻‍🤝‍🧑🏿
    Emoji aliases are 😎
    Use fmt wrappers 👍 with emoji support 🎉
*/
```

This package contains emojis constants based on [Full Emoji List v13.0](https://unicode.org/Public/emoji/13.0/emoji-test.txt).
```go
    emoji.CallMeHand // 🤙
```
Also, it has additional emoji aliases from [github/gemoji](https://github.com/github/gemoji).
```go
    emoji.Parse(":+1:") // 👍
    emoji.Parse(":100:") // 💯
```

Also, you can generate country flag emoji with [ISO 3166 Alpha2](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2) codes:
```go
emoji.CountryFlag("tr") // 🇹🇷
emoji.CountryFlag("US") // 🇺🇸
emoji.CountryFlag("gb") // 🇬🇧
```

All constants are generated by `internal/generator`.

## Testing :hammer:
``` bash
go test
```

## Todo :pushpin:
* Add emoji string parser

## Contributing :man_technologist:
I am accepting PRs that add aliases to the package.
You have to add it to `customEmojis` list at `internal/generator/main`.

If you think an emoji constant is not correct, open an issue.
Please use [this list](http://unicode.org/emoji/charts/full-emoji-list.html)
to look up the correct unicode value and the name of the character.

## Credits :star:
- [Enes Çakır](https://github.com/enescakir)

## License :scroll:
The MIT License (MIT). Please see [License File](LICENSE.md) for more information.
