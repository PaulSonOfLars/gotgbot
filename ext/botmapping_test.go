package ext

import (
	"errors"
	"testing"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

func Test_botMapping(t *testing.T) {
	bm := botMapping{}
	b := &gotgbot.Bot{
		User:      gotgbot.User{},
		Token:     "SOME_TOKEN",
		BotClient: &gotgbot.BaseBotClient{},
	}

	var origBdata *botData
	t.Run("addBot", func(t *testing.T) {
		// check that bots can be added fine
		var err error
		origBdata, err = bm.addPollingBot(b, nil)
		if err != nil {
			t.Errorf("expected to be able to add a new bot fine: %s", err.Error())
			t.FailNow()
		}
		if len(bm.getBots()) != 1 {
			t.Errorf("expected 1 bot, got %d", len(bm.getBots()))
			t.FailNow()
		}
	})

	t.Run("doubleAdd", func(t *testing.T) {
		// Adding the same bot twice should fail
		_, err := bm.addPollingBot(b, nil)
		if err == nil {
			t.Errorf("adding the same bot twice should throw an error")
			t.FailNow()
		}
		if len(bm.getBots()) != 1 {
			t.Errorf("expected only haveing 1 bot when adding a duplicate, but got %d", len(bm.getBots()))
			t.FailNow()
		}
	})

	t.Run("getBot", func(t *testing.T) {
		// check that bot data is correct
		bdata, ok := bm.getBot(b.Token)
		if !ok {
			t.Errorf("failed to get bot with token %s", b.Token)
			t.FailNow()
		}
		if bdata.stopUpdates != origBdata.stopUpdates {
			t.Errorf("stopUpdates channel was not the same")
			t.FailNow()
		}
		if bdata.updateChan != origBdata.updateChan {
			t.Errorf("update channel was not the same")
			t.FailNow()
		}
	})

	t.Run("removeBot", func(t *testing.T) {
		// check that bot cant be removed
		_, ok := bm.removeBot(b.Token)
		if !ok {
			t.Errorf("failed to remove bot with token %s", b.Token)
			t.FailNow()
		}

		_, ok = bm.getBot(b.Token)
		if ok {
			t.Errorf("bot with token %s should be gone", b.Token)
			t.FailNow()
		}
	})
}

func Test_botMapping_rootWebhookPath(t *testing.T) {
	bm := botMapping{}

	b := &gotgbot.Bot{
		User:      gotgbot.User{},
		Token:     "WEBHOOK_TOKEN",
		BotClient: &gotgbot.BaseBotClient{},
	}

	_, err := bm.addWebhookBot(b, "", "")
	if err != nil {
		t.Fatalf("expected root webhook bot to be added, got error: %s", err.Error())
	}

	bdata, ok := bm.getBotFromURL("")
	if !ok {
		t.Fatal("expected root webhook path to resolve to a bot")
	}

	if bdata.bot.Token != b.Token {
		t.Errorf("expected bot token %s, got %s", b.Token, bdata.bot.Token)
	}
}

func Test_botMapping_webhookPath(t *testing.T) {
	bm := botMapping{}

	b := &gotgbot.Bot{
		User:      gotgbot.User{},
		Token:     "WEBHOOK_TOKEN",
		BotClient: &gotgbot.BaseBotClient{},
	}

	_, err := bm.addWebhookBot(b, "/webhook", "")
	if err != nil {
		t.Fatalf("expected webhook bot to be added, got error: %s", err.Error())
	}

	bdata, ok := bm.getBotFromURL("webhook")
	if !ok {
		t.Fatal("expected webhook path to resolve to a bot")
	}

	if bdata.bot.Token != b.Token {
		t.Errorf("expected bot token %s, got %s", b.Token, bdata.bot.Token)
	}
}

func Test_botMapping_rootWebhookPathAlreadyExists(t *testing.T) {
	bm := botMapping{}

	b1 := &gotgbot.Bot{
		User:      gotgbot.User{},
		Token:     "WEBHOOK_TOKEN_1",
		BotClient: &gotgbot.BaseBotClient{},
	}

	b2 := &gotgbot.Bot{
		User:      gotgbot.User{},
		Token:     "WEBHOOK_TOKEN_2",
		BotClient: &gotgbot.BaseBotClient{},
	}

	_, err := bm.addWebhookBot(b1, "", "")
	if err != nil {
		t.Fatalf("expected first root webhook bot to be added, got error: %v", err)
	}

	_, err = bm.addWebhookBot(b2, "", "")
	if !errors.Is(err, ErrBotUrlPathAlreadyExists) {
		t.Fatalf("expected ErrBotUrlPathAlreadyExists, got %v", err)
	}
}

func Test_botData_isUpdateChannelStopped(t *testing.T) {
	bm := botMapping{}
	b := &gotgbot.Bot{
		User:      gotgbot.User{},
		Token:     "SOME_TOKEN",
		BotClient: &gotgbot.BaseBotClient{},
	}

	ctxCancelled := false
	bData, err := bm.addPollingBot(b, func() {
		ctxCancelled = true
	})
	if err != nil {
		t.Errorf("bot with token %s should not have failed to be added", b.Token)
		return
	}
	if bData.shouldStopUpdates() {
		t.Errorf("bot with token %s should not be stopped yet", b.Token)
		return
	}

	bData.stop()
	if !bData.shouldStopUpdates() {
		t.Errorf("bot with token %s should be stopped", b.Token)
		return
	}
	if !ctxCancelled {
		t.Errorf("bot with token %s should have a cancelled context ", b.Token)
	}
}