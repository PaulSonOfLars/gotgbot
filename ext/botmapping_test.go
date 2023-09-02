package ext

import (
	"encoding/json"
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

	updateChan := make(chan json.RawMessage)
	pollChan := make(chan struct{})

	t.Run("addBot", func(t *testing.T) {
		// check that bots can be added fine
		err := bm.addBot(b, updateChan, pollChan, "")
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
		err := bm.addBot(b, updateChan, pollChan, "")
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
		if bdata.polling != pollChan {
			t.Errorf("polling channel was not the same")
			t.FailNow()
		}
		if bdata.updateChan != updateChan {
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
