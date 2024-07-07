package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/precheckoutquery"
)

// This bot demonstrates how to provide invoices, checkouts, and successful payments through telegram's in-app purchase
// methods.
// Use this if you want an example of how to sell things through telegram. The example targets Telegram Stars, which
// allows bot developers to sell digital products through Telegram.
func main() {
	// Get token from the environment variable
	token := os.Getenv("TOKEN")
	if token == "" {
		panic("TOKEN environment variable is empty")
	}

	// Create bot from environment value.
	b, err := gotgbot.NewBot(token, nil)
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	// Create updater and dispatcher.
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		// If an error is returned by a handler, log it and continue going.
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})
	updater := ext.NewUpdater(dispatcher, nil)

	// /start command to introduce the bot
	dispatcher.AddHandler(handlers.NewCommand("start", start))
	// PreCheckout to handle the step right before payment. Must be handled within 10s, or the checkout will be abandoned by telegram.
	dispatcher.AddHandler(handlers.NewPreCheckoutQuery(precheckoutquery.All, preCheckout))
	// Payment received; send/provide product to customer.
	dispatcher.AddHandler(handlers.NewMessage(message.SuccessfulPayment, paymentComplete))
	// Bots selling on telegram must be able to provide refunds; do so through the paysupport command, as mentioned in
	// the TOS: https://telegram.org/tos/stars#3-1-disputing-purchases
	dispatcher.AddHandler(handlers.NewCommand("paysupport", paySupport))

	// Start receiving updates.
	err = updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}
	log.Printf("%s has been started...\n", b.User.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
}

// start introduces the bot and sends an initial invoice (in Telegram stars; denoted as XTR).
func start(b *gotgbot.Bot, ctx *ext.Context) error {
	if ctx.EffectiveChat.Type != "private" {
		// Only reply in private chats.
		return nil
	}

	// Introduce the bot.
	_, err := ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Hello, I'm @%s. I demonstrate how telegram payments might work.", b.User.Username), &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}

	// Generate a unique payload for this checkout.
	// In a production environment, you could create a database containing any additional information you might need to
	// complete this transaction, and refer to it at the preCheckout and successful payment stages.
	// For example, you may want to store the invoice creator's ID to ensure that only the creator can pay, and any
	// other necessary data you have collected.
	payload := uuid.NewString()

	// Send the invoice. XTR == Telegram stars, for selling digital products.
	_, err = b.SendInvoice(ctx.EffectiveChat.Id, "Product Name", "Some detailed description", payload, "XTR", []gotgbot.LabeledPrice{{
		Label:  "Some product",
		Amount: 100, // 100 stars.
	}}, &gotgbot.SendInvoiceOpts{
		ProtectContent: true, // Stop people from forwarding this invoice to others.
	})
	if err != nil {
		return fmt.Errorf("failed to generate invoice: %w", err)
	}

	return nil
}

func preCheckout(b *gotgbot.Bot, ctx *ext.Context) error {
	// Do any required preCheckout validation here. If anything failed, we should answer the query with "ok: False",
	// and populate the ErrorMessage field in the opts.
	// For example, you may want to ensure that the user who requested the invoice is the same person as the person who
	// is checking out; but this would require storage, so isn't shown here.

	// Answer true once checks have passed.
	_, err := ctx.PreCheckoutQuery.Answer(b, true, nil)
	if err != nil {
		return fmt.Errorf("failed to answer precheckout query: %w", err)
	}
	return nil
}

func paymentComplete(b *gotgbot.Bot, ctx *ext.Context) error {
	// Payment has been received; a real bot would now provide the user with the product.
	_, err := ctx.EffectiveMessage.Reply(b, "Payment complete - in a real bot, this is where you would provision the product that has been paid for.", nil)
	if err != nil {
		return fmt.Errorf("failed to send payment complete message: %w", err)
	}
	return nil
}

func paySupport(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b, "Explain your refund process here.", nil)
	if err != nil {
		return fmt.Errorf("failed to describe refund process: %w", err)
	}
	return nil
}
