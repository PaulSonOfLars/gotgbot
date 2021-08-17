# Golang Telegram Bot library

Heavily inspired by the [python-telegram-bot library](https://github.com/python-telegram-bot/python-telegram-bot),
this package is a code-generated wrapper for the telegram bot api. We also provide an extensions package which 
defines an updater/dispatcher pattern to avoid having to rewrite update processing.

All the telegram types and methods are present in the `gen_*.go` files. These are all generated from 
[a bot api spec](https://github.com/PaulSonOfLars/telegram-bot-api-spec) and can be recreated by running `go generate` 
in the repo root. This makes it extremely easy to update the library; simply download the latest spec, and regenerate.

If you have any questions, come find us in our [telegram support chat](https://t.me/GotgbotChat)!

## Features:
- All telegram API types and methods are generated from the bot api docs, which makes this library:
    - Guaranteed to match the docs
    - Easy to update
    - Self-documenting (Can simply reuse pre-existing telegram docs)
- Type safe; no weird interface{} logic.
- No third party library bloat; only uses standard library.
- Updates are each processed in their own go routine, encouraging concurrent processing, and keeping your bot responsive.
- Code panics automatically recovered from and logged, avoiding unexpected downtime.

## Getting started

Download the library with the standard `go get` command:

```bash
go get github.com/PaulSonOfLars/gotgbot/v2
```

### Example bots

Sample bots can be found in the [samples directory](samples).

## Docs

Docs can be found [here](https://pkg.go.dev/github.com/PaulSonOfLars/gotgbot/v2).

## Contributing

Contributions are welcome! More information on contributing can be found [here](.github/CONTRIBUTING.md).
