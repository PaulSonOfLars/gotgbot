# Go Telegram Bot

**This library is WIP; it does not currently support all of the telegram api methods.**

This library attempts to create a user-friendly wrapper around the telegram bot api.

Heavily inspired by the [python-telegram-bot library](github.com/python-telegram-bot/python-telegram-bot),
this aims to create a simple way to manage a concurrent and scalable bot.

## Getting started
Install it as you would install your usual go library: `go get github.com/PaulSonOfLars/gotgbot`

A sample bot can be found in `sampleBot/`. This bot covers the basics of adding a command, a filter, and a regex handler.


An interesting feature to take note of is that due to go's
handling of exceptions, if you choose not to handle an exception, your bot
will simply keep on going happily and ignore any issues.

All handlers are async; they're all executed in their own go routine,
so can communicate accross channels if needed.
The reason for the `error` return for the methods is to allow for passing `gotgbot.ContinueGroups{}`
or `gotgbot.EndGroups{}`; which will determine whether or not to keep handling methods in that handler group,
or stop handling further groups entirely.

## Message sending

As seen in the example, message sending can be done in two ways; via each received message's
ReplyText() function, or by building your own; and calling msg.Send(). This allows for
ease of use by having the most commonly used shortcuts readily available, while
retaining the flexibility of building each message yourself, which wouldnt be
available otherwise.

