package ext

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

type Updater struct {
	Bot        gotgbot.Bot
	Dispatcher *Dispatcher
	UpdateChan chan json.RawMessage
}

// NewUpdater Creates a new Updater, as well as the necessary structures for
func NewUpdater(bot *gotgbot.Bot) Updater {
	updateChan := make(chan json.RawMessage)
	return Updater{
		Bot: gotgbot.Bot{
			Token:  bot.Token,
			APIURL: bot.APIURL,
			Client: http.Client{},
		}, // create new bot client to allow for independent timeout changes
		Dispatcher: NewDispatcher(updateChan, DefaultMaxRoutines),
		UpdateChan: updateChan,
	}
}

// StartPolling Starts the polling logic
func (u Updater) StartPolling(b *gotgbot.Bot) {
	go u.Dispatcher.Start(b)
	go u.startPolling(false)
}

// StartCleanPolling Starts clean polling (ignoring stale updates)
func (u Updater) StartCleanPolling(b *gotgbot.Bot) {
	go u.Dispatcher.Start(b)
	go u.startPolling(true)
}

func (u Updater) startPolling(clean bool) {
	v := url.Values{}
	v.Add("offset", strconv.Itoa(0))
	v.Add("timeout", strconv.Itoa(0))
	var offset int64

	for {
		// note: this bot has a custom client with longer timeouts
		r, err := u.Bot.Get("getUpdates", v)
		if err != nil {
			time.Sleep(time.Second)
			continue

		} else if r != nil {
			var rawUpdates []json.RawMessage
			if err := json.Unmarshal(r, &rawUpdates); err != nil {
				os.Stderr.WriteString("failed to unmarshal updates: " + err.Error())
				continue
			}
			if len(rawUpdates) > 0 {

				var lastUpdate gotgbot.Update
				if err := json.Unmarshal(rawUpdates[len(rawUpdates)-1], &lastUpdate); err != nil {
					os.Stderr.WriteString("failed to unmarshal last update: " + err.Error())
					continue
				}

				offset = lastUpdate.UpdateId + 1
				v.Set("offset", strconv.FormatInt(offset, 10))
				if clean {
					continue
				}
			} else if len(rawUpdates) == 0 { // TODO: check this is fine on high loads
				clean = false
			}

			for _, updData := range rawUpdates {
				temp := updData // use new mem address to avoid loop conflicts
				u.UpdateChan <- temp
			}
		}
	}
}

// Idle sets the main thread to idle, allowing the background processes to run as expected (dispatcher and update handlers)
func (u Updater) Idle() {
	for {
		time.Sleep(1 * time.Second)
	}
}

//type WebhookOpts struct {
//}
//
//// StartWebhook Starts the webhook server
//func (u Updater) StartWebhook(b *gotgbot.Bot, opts WebhookOpts) {
//	go u.Dispatcher.Start(b)
//
//	mux := http.NewServeMux()
//	mux.HandleFunc("/"+webhook.ServePath, func(w http.ResponseWriter, r *http.Request) {
//		bytes, _ := ioutil.ReadAll(r.Body)
//		u.UpdateChan <- bytes
//	})
//
//	server := http.Server{
//		Addr:    opts.Addr,
//		Handler: mux,
//	}
//
//	go func() {
//		// todo: TLS when using certs
//		err := server.ListenAndServe()
//		if err != nil {
//			os.Stderr.WriteString("http server failed: " + err.Error())
//		}
//	}()
//}
