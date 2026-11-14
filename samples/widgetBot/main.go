package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/sharedwidgetquery"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/widgetquery"
)

var (
	startTime    = time.Now()
	queryCounter atomic.Int64
	httpClient   = &http.Client{Timeout: 6 * time.Second}
)

func main() {
	token := os.Getenv("TOKEN")
	if token == "" {
		panic("TOKEN environment variable is empty")
	}

	b, err := gotgbot.NewBot(token, &gotgbot.BotOpts{
		BotClient: &gotgbot.BaseBotClient{
			Client: http.Client{},
			DefaultRequestOpts: &gotgbot.RequestOpts{
				Timeout: gotgbot.DefaultTimeout,
				APIURL:  gotgbot.DefaultAPIURL,
			},
		},
	})
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			logger.Error("handler error", "error", err)
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
		Logger:      logger,
	})
	updater := ext.NewUpdater(dispatcher, &ext.UpdaterOpts{Logger: logger})

	// Command & message handlers
	dispatcher.AddHandler(handlers.NewCommand("start", cmdStart))
	dispatcher.AddHandler(handlers.NewCommand("help", cmdStart))
	dispatcher.AddHandler(handlers.NewCommand("weather", cmdWeather))
	dispatcher.AddHandler(handlers.NewCommand("location", cmdLocation))
	dispatcher.AddHandler(handlers.NewCommand("dice", cmdDice))
	dispatcher.AddHandler(handlers.NewCommand("tabs", cmdTabs))
	dispatcher.AddHandler(handlers.NewCommand("time", cmdTime))
	dispatcher.AddHandler(handlers.NewCommand("counter", cmdCounter))
	dispatcher.AddHandler(handlers.NewCommand("tictactoe", cmdTicTacToe))
	dispatcher.AddHandler(handlers.NewCommand("poll", cmdPoll))
	dispatcher.AddHandler(handlers.NewCommand("dashboard", cmdDashboard))

	// Receive GPS location shared by the user for live local weather
	dispatcher.AddHandler(handlers.NewMessage(message.Location, handleLocationMessage))

	// Widget query handlers
	dispatcher.AddHandler(handlers.NewWidgetQuery(widgetquery.Prefix("weather:"), handleWeatherQuery))
	dispatcher.AddHandler(handlers.NewWidgetQuery(widgetquery.Prefix("dice:"), handleDiceQuery))
	dispatcher.AddHandler(handlers.NewWidgetQuery(widgetquery.Prefix("tab:"), handleTabQuery))
	dispatcher.AddHandler(handlers.NewWidgetQuery(widgetquery.Prefix("time:"), handleTimeQuery))
	dispatcher.AddHandler(handlers.NewWidgetQuery(widgetquery.Prefix("counter:"), handleCounterQuery))
	dispatcher.AddHandler(handlers.NewWidgetQuery(widgetquery.Prefix("ttt:"), handleTTTQuery))
	dispatcher.AddHandler(handlers.NewWidgetQuery(widgetquery.Prefix("poll:"), handlePollQuery))
	dispatcher.AddHandler(handlers.NewWidgetQuery(widgetquery.Prefix("dash:"), handleDashboardQuery))

	// Shared widget queries
	dispatcher.AddHandler(handlers.NewSharedWidgetQuery(sharedwidgetquery.Prefix("poll:"), handleSharedPollQuery))

	// Fallback widget query handler
	dispatcher.AddHandler(handlers.NewWidgetQuery(widgetquery.All, func(b *gotgbot.Bot, ctx *ext.Context) error {
		if ctx.WidgetQuery == nil {
			return nil
		}
		_, err := b.AnswerWidgetQuery(ctx.WidgetQuery.Id, &gotgbot.AnswerWidgetQueryOpts{Data: `{"ok":true}`})
		return err
	}))

	_, _ = b.SetMyCommands([]gotgbot.BotCommand{
		{Command: "start", Description: "List all HTML widget demos"},
		{Command: "weather", Description: "Live weather (supports /weather <city>)"},
		{Command: "location", Description: "Share GPS location for live local weather"},
		{Command: "dice", Description: "Interactive dice roller with SVG pips"},
		{Command: "tabs", Description: "Segmented tabs + clipboard & ping"},
		{Command: "time", Description: "World clock & stopwatch"},
		{Command: "counter", Description: "Persistent stepper counter"},
		{Command: "tictactoe", Description: "Tic-Tac-Toe (vs Bot or 2P)"},
		{Command: "poll", Description: "Live shared poll"},
		{Command: "dashboard", Description: "Go runtime telemetry card"},
	}, nil)

	err = updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			AllowedUpdates: []string{
				"message",
				"widget_query",
				"shared_widget_query",
			},
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: 15 * time.Second,
			},
		},
	})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}
	logger.Info("Widget bot started", "username", b.User.Username)
	updater.Idle()
}

func widgetKey(wq *gotgbot.WidgetQuery) string {
	if wq.InlineMessageId != "" {
		return "inline:" + wq.InlineMessageId
	}
	if wq.Chat != nil {
		return fmt.Sprintf("msg:%d:%d", wq.Chat.Id, wq.MessageId)
	}
	return "user:" + strconv.FormatInt(wq.From.Id, 10)
}

func persistChatData(b *gotgbot.Bot, wq *gotgbot.WidgetQuery, data string) {
	if wq == nil {
		return
	}
	opts := &gotgbot.EditMessageWidgetDataOpts{ChatData: data}
	if wq.InlineMessageId != "" {
		opts.InlineMessageId = wq.InlineMessageId
	} else if wq.Chat != nil && wq.MessageId != 0 {
		opts.ChatId = wq.Chat.Id
		opts.MessageId = wq.MessageId
	} else {
		return
	}
	_, _ = b.EditMessageWidgetData(opts)
}

func answerJSON(b *gotgbot.Bot, queryID string, payload any) error {
	queryCounter.Add(1)
	raw, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = b.AnswerWidgetQuery(queryID, &gotgbot.AnswerWidgetQueryOpts{Data: string(raw)})
	return err
}

// Full-bleed base CSS (keeps heavy paint boxes minimal: avoids CSS opacity<1 and heavy shadows).
const fullBleedCSS = `
*{box-sizing:border-box;margin:0;padding:0}
body{font-family:system-ui,sans-serif;margin:0;padding:0;color:var(--tg-theme-text-color,#fff);background:var(--tg-theme-bg-color,#17212b);-webkit-font-smoothing:antialiased}
.wrap{width:100%;padding:14px 14px max(22px,calc(env(safe-area-inset-bottom,0px) + 10px))}
.row{display:flex;align-items:center;justify-content:space-between;gap:8px}
.muted{color:var(--tg-theme-hint-color,rgba(255,255,255,.65));font-size:12px}
.btn{flex:1;padding:9px 10px;border:none;border-radius:10px;background:var(--tg-theme-button-color,#3390ec);color:var(--tg-theme-button-text-color,#fff);font-size:13px;font-weight:600;cursor:pointer;transition:transform .12s}
.btn.sec{background:rgba(255,255,255,.14);color:inherit}
.btn:active{transform:scale(.96)}
.pill{padding:4px 10px;border-radius:999px;font-size:11px;font-weight:600;background:rgba(255,255,255,.16);color:inherit;border:none;cursor:pointer;transition:transform .12s}
.pill:active{transform:scale(.94)}
`

func cmdStart(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b,
		"<b>HTML Widgets Showcase (Bot API 10.4)</b>\n\n"+
			"/weather — Live weather & animated forecast\n"+
			"/weather Tokyo — Live weather for any city\n"+
			"/location — Share GPS for your local weather\n"+
			"/dice — Interactive SVG dice roller\n"+
			"/tabs — Segmented tabs & clipboard\n"+
			"/time — Live world clock & stopwatch\n"+
			"/counter — Persistent counter\n"+
			"/tictactoe — Tic-Tac-Toe (vs Bot or 2P)\n"+
			"/poll — Shared live poll\n"+
			"/dashboard — Live Go runtime metrics",
		&gotgbot.SendMessageOpts{
			ParseMode: "HTML",
			ReplyMarkup: &gotgbot.ReplyKeyboardMarkup{
				Keyboard: [][]gotgbot.KeyboardButton{{{
					Text:            "📍 Share My Location Weather",
					RequestLocation: true,
				}}},
				ResizeKeyboard:  true,
				OneTimeKeyboard: true,
			},
		},
	)
	return err
}

func cmdLocation(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b,
		"Tap the button below (or attach a 📍 Location) to get live local weather:",
		&gotgbot.SendMessageOpts{
			ReplyMarkup: &gotgbot.ReplyKeyboardMarkup{
				Keyboard: [][]gotgbot.KeyboardButton{{{
					Text:            "📍 Get My Current Location Weather",
					RequestLocation: true,
				}}},
				ResizeKeyboard:  true,
				OneTimeKeyboard: true,
			},
		},
	)
	return err
}

type weatherHour struct {
	Label string `json:"l"`
	Icon  string `json:"ic"`
	Temp  int    `json:"t"`
	Cond  string `json:"c"`
	Sky   string `json:"s"` // rain | sun | cloud | night
}

type weatherState struct {
	Idx   int           `json:"i"`
	City  string        `json:"city"`
	Sky   string        `json:"sky"`
	High  int           `json:"hi"`
	Low   int           `json:"lo"`
	IsGPS bool          `json:"gps,omitempty"`
	Hours []weatherHour `json:"hrs"`
}

type savedPlace struct {
	Name  string
	Lat   float64
	Lon   float64
	IsGPS bool
}

var (
	userLocMu    sync.RWMutex
	userLocation = map[int64]savedPlace{}
)

var defaultCities = []savedPlace{
	{Name: "Barcelona", Lat: 41.3888, Lon: 2.1590},
	{Name: "London", Lat: 51.5074, Lon: -0.1278},
	{Name: "Hyderabad", Lat: 17.3850, Lon: 78.4867},
	{Name: "Tokyo", Lat: 35.6762, Lon: 139.6503},
	{Name: "New York", Lat: 40.7128, Lon: -74.0060},
}

func wmoInfo(code, isDay int) (sky, icon, desc string) {
	switch {
	case code == 0:
		if isDay == 0 {
			return "night", "🌙", "Clear night"
		}
		return "sun", "☀️", "Sunny"
	case code == 1:
		if isDay == 0 {
			return "night", "🌙", "Mostly clear"
		}
		return "sun", "🌤️", "Mostly sunny"
	case code == 2:
		if isDay == 0 {
			return "night", "☁️", "Partly cloudy"
		}
		return "cloud", "⛅", "Partly cloudy"
	case code == 3:
		return "cloud", "☁️", "Overcast"
	case code == 45 || code == 48:
		return "cloud", "🌫️", "Foggy"
	case code >= 51 && code <= 57:
		return "rain", "🌦️", "Drizzle"
	case code == 61 || code == 80:
		return "rain", "🌧️", "Light rain"
	case (code >= 63 && code <= 67) || code == 81 || code == 82:
		return "rain", "🌧️", "Rain"
	case (code >= 71 && code <= 77) || code == 85 || code == 86:
		return "cloud", "❄️", "Snow"
	case code >= 95:
		return "rain", "⛈️", "Thunderstorm"
	default:
		return "cloud", "⛅", "Cloudy"
	}
}

func reverseGeocode(lat, lon float64) string {
	u := fmt.Sprintf("https://api.bigdatacloud.net/data/reverse-geocode-client?latitude=%.5f&longitude=%.5f&localityLanguage=en", lat, lon)
	resp, err := httpClient.Get(u)
	if err != nil {
		return fmt.Sprintf("%.2f, %.2f", lat, lon)
	}
	defer resp.Body.Close()
	var out struct {
		City     string `json:"city"`
		Locality string `json:"locality"`
		Region   string `json:"principalSubdivision"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return fmt.Sprintf("%.2f, %.2f", lat, lon)
	}
	if out.City != "" {
		return out.City
	}
	if out.Locality != "" {
		return out.Locality
	}
	if out.Region != "" {
		return out.Region
	}
	return fmt.Sprintf("%.2f, %.2f", lat, lon)
}

func searchCityCoords(name string) (savedPlace, bool) {
	u := "https://geocoding-api.open-meteo.com/v1/search?name=" + url.QueryEscape(name) + "&count=1&language=en"
	resp, err := httpClient.Get(u)
	if err != nil {
		return savedPlace{}, false
	}
	defer resp.Body.Close()
	var out struct {
		Results []struct {
			Name      string  `json:"name"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil || len(out.Results) == 0 {
		return savedPlace{}, false
	}
	r := out.Results[0]
	return savedPlace{Name: r.Name, Lat: r.Latitude, Lon: r.Longitude, IsGPS: true}, true
}

func fetchLiveWeather(place savedPlace, idx int) weatherState {
	u := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%.5f&longitude=%.5f&current=temperature_2m,weather_code,is_day&hourly=temperature_2m,weather_code,is_day&daily=temperature_2m_max,temperature_2m_min&timezone=auto&forecast_days=2",
		place.Lat, place.Lon,
	)
	resp, err := httpClient.Get(u)
	if err == nil {
		defer resp.Body.Close()
		var om struct {
			Current struct {
				Time  string  `json:"time"`
				Temp  float64 `json:"temperature_2m"`
				Code  int     `json:"weather_code"`
				IsDay int     `json:"is_day"`
			} `json:"current"`
			Hourly struct {
				Time  []string  `json:"time"`
				Temp  []float64 `json:"temperature_2m"`
				Code  []int     `json:"weather_code"`
				IsDay []int     `json:"is_day"`
			} `json:"hourly"`
			Daily struct {
				Max []float64 `json:"temperature_2m_max"`
				Min []float64 `json:"temperature_2m_min"`
			} `json:"daily"`
		}
		if json.NewDecoder(resp.Body).Decode(&om) == nil && len(om.Hourly.Time) > 6 {
			curSky, curIcon, curDesc := wmoInfo(om.Current.Code, om.Current.IsDay)
			hrs := []weatherHour{{
				Label: "Now",
				Icon:  curIcon,
				Temp:  int(math.Round(om.Current.Temp)),
				Cond:  curDesc,
				Sky:   curSky,
			}}
			startIdx := 0
			if len(om.Current.Time) >= 13 {
				prefix := om.Current.Time[:13]
				for i, ht := range om.Hourly.Time {
					if strings.HasPrefix(ht, prefix) {
						startIdx = i
						break
					}
				}
			}
			for step := 1; step <= 5 && startIdx+step < len(om.Hourly.Time); step++ {
				j := startIdx + step
				lbl := fmt.Sprintf("+%dh", step)
				if parsed, perr := time.Parse("2006-01-02T15:04", om.Hourly.Time[j]); perr == nil {
					lbl = parsed.Format("3 PM")
				}
				hSky, hIcon, hDesc := wmoInfo(om.Hourly.Code[j], om.Hourly.IsDay[j])
				hrs = append(hrs, weatherHour{
					Label: lbl,
					Icon:  hIcon,
					Temp:  int(math.Round(om.Hourly.Temp[j])),
					Cond:  hDesc,
					Sky:   hSky,
				})
			}
			hi := int(math.Round(om.Current.Temp))
			lo := hi
			if len(om.Daily.Max) > 0 {
				hi = int(math.Round(om.Daily.Max[0]))
			}
			if len(om.Daily.Min) > 0 {
				lo = int(math.Round(om.Daily.Min[0]))
			}
			return weatherState{
				Idx:   idx,
				City:  place.Name,
				Sky:   curSky,
				High:  hi,
				Low:   lo,
				IsGPS: place.IsGPS,
				Hours: hrs,
			}
		}
	}

	// Fallback if Open-Meteo is unreachable
	return weatherState{
		Idx: idx, City: place.Name, Sky: "rain", High: 24, Low: 15, IsGPS: place.IsGPS,
		Hours: []weatherHour{
			{"Now", "🌧️", 18, "Light rain", "rain"},
			{"1 PM", "🌦️", 19, "Showers", "rain"},
			{"2 PM", "☁️", 21, "Cloudy", "cloud"},
			{"3 PM", "⛅", 23, "Partly sunny", "cloud"},
			{"4 PM", "☀️", 24, "Sunny", "sun"},
			{"5 PM", "☀️", 23, "Clear", "sun"},
		},
	}
}

func userPlaces(userID int64) []savedPlace {
	userLocMu.RLock()
	p, ok := userLocation[userID]
	userLocMu.RUnlock()
	if !ok {
		return defaultCities
	}
	list := make([]savedPlace, 0, len(defaultCities)+1)
	list = append(list, p)
	list = append(list, defaultCities...)
	return list
}

func cmdWeather(b *gotgbot.Bot, ctx *ext.Context) error {
	args := strings.TrimSpace(strings.TrimPrefix(ctx.EffectiveMessage.Text, "/weather"))
	if strings.HasPrefix(args, "@") {
		if idx := strings.IndexByte(args, ' '); idx != -1 {
			args = strings.TrimSpace(args[idx+1:])
		} else {
			args = ""
		}
	}
	if args != "" {
		if p, ok := searchCityCoords(args); ok {
			userLocMu.Lock()
			userLocation[ctx.EffectiveUser.Id] = p
			userLocMu.Unlock()
		}
	}

	places := userPlaces(ctx.EffectiveUser.Id)
	st := fetchLiveWeather(places[0], 0)
	raw, _ := json.Marshal(st)

	opts := &gotgbot.SendWidgetOpts{
		AllowQueries: true,
		ChatData:     string(raw),
	}
	userLocMu.RLock()
	_, hasLoc := userLocation[ctx.EffectiveUser.Id]
	userLocMu.RUnlock()
	if !hasLoc {
		opts.ReplyMarkup = &gotgbot.ReplyKeyboardMarkup{
			Keyboard: [][]gotgbot.KeyboardButton{{{
				Text:            "📍 My Location Weather",
				RequestLocation: true,
			}}},
			ResizeKeyboard:  true,
			OneTimeKeyboard: true,
		}
	}

	_, err := b.SendWidget(ctx.EffectiveChat.Id, weatherHTML(), opts)
	return err
}

func handleLocationMessage(b *gotgbot.Bot, ctx *ext.Context) error {
	loc := ctx.EffectiveMessage.Location
	if loc == nil {
		return nil
	}
	city := reverseGeocode(loc.Latitude, loc.Longitude)
	place := savedPlace{Name: city, Lat: loc.Latitude, Lon: loc.Longitude, IsGPS: true}

	userLocMu.Lock()
	userLocation[ctx.EffectiveUser.Id] = place
	userLocMu.Unlock()

	st := fetchLiveWeather(place, 0)
	raw, _ := json.Marshal(st)
	_, err := b.SendWidget(ctx.EffectiveChat.Id, weatherHTML(), &gotgbot.SendWidgetOpts{
		AllowQueries: true,
		ChatData:     string(raw),
		ReplyMarkup:  &gotgbot.ReplyKeyboardRemove{RemoveKeyboard: true},
	})
	return err
}

func handleWeatherQuery(b *gotgbot.Bot, ctx *ext.Context) error {
	wq := ctx.WidgetQuery
	if wq == nil {
		return nil
	}
	parts := strings.Split(wq.Data, ":")
	action, idx := "", 0
	if len(parts) >= 2 {
		action = parts[1]
	}
	if len(parts) >= 3 {
		idx, _ = strconv.Atoi(parts[2])
	}
	places := userPlaces(wq.From.Id)
	if action == "next" {
		idx = (idx + 1) % len(places)
	} else {
		idx = ((idx % len(places)) + len(places)) % len(places)
	}
	st := fetchLiveWeather(places[idx], idx)
	raw, _ := json.Marshal(st)
	persistChatData(b, wq, string(raw))
	return answerJSON(b, wq.Id, st)
}

func weatherHTML() string {
	// Avoid CSS opacity<1 and opacity animations to stay under Telegram's paint-heavy limit.
	return `<!DOCTYPE html><html><head><meta charset="utf-8"><style>
*{box-sizing:border-box;margin:0;padding:0}
body{font-family:system-ui,sans-serif;margin:0;padding:0;color:#fff;-webkit-font-smoothing:antialiased}
.sky{position:relative;overflow:hidden;padding:16px 14px 26px;background:linear-gradient(180deg,#44546b,#69778c);transition:background .6s}
.sky.sun{background:linear-gradient(180deg,#1d6fd6,#5eb2ff)}
.sky.cloud{background:linear-gradient(180deg,#4e7db8,#8eb2dc)}
.sky.night{background:linear-gradient(180deg,#0c1836,#233a6b)}
.sky svg.bg{position:absolute;inset:0;width:100%;height:100%;pointer-events:none}
.rf1{animation:f1 1.05s linear infinite}
.rf2{animation:f2 1.35s linear infinite}
@keyframes f1{from{transform:translate(8px,0)}to{transform:translate(-14px,280px)}}
@keyframes f2{from{transform:translate(4px,-140px)}to{transform:translate(-18px,140px)}}
.cld{animation:drift 14s ease-in-out infinite alternate}
@keyframes drift{from{transform:translateX(-6px)}to{transform:translateX(14px)}}
.sunorb{transform-origin:278px 66px;animation:spin 32s linear infinite}
@keyframes spin{to{transform:rotate(360deg)}}
.top{position:relative;z-index:2;display:flex;justify-content:space-between;align-items:center}
.citybtn{background:none;border:none;color:#fff;font-size:21px;font-weight:500;cursor:pointer;display:inline-flex;align-items:center;gap:6px;padding:0;text-align:left}
.citybtn:active{transform:scale(.97)}
.tag{font-size:11px;color:rgba(255,255,255,.92);background:rgba(255,255,255,.18);padding:3px 8px;border-radius:999px;font-weight:600}
.acts{display:flex;gap:6px;align-items:center}
.pbtn{background:rgba(255,255,255,.18);border:none;color:#fff;font-size:12px;font-weight:600;padding:5px 10px;border-radius:999px;cursor:pointer;transition:transform .12s,background .2s}
.pbtn:active{transform:scale(.92);background:rgba(255,255,255,.32)}
.bigt{position:relative;z-index:2;font-size:66px;font-weight:300;line-height:1;letter-spacing:-2px;margin:6px 0 12px;font-variant-numeric:tabular-nums}
.sub{position:relative;z-index:2;font-size:14px;color:rgba(255,255,255,.92);margin-bottom:2px}
.hl{position:relative;z-index:2;font-size:13px;color:rgba(255,255,255,.8);margin-bottom:18px;font-variant-numeric:tabular-nums}
.dock{position:relative;z-index:2;display:grid;grid-template-columns:repeat(6,1fr);padding:4px;border-radius:20px;background:rgba(255,255,255,.14);border:1px solid rgba(255,255,255,.2)}
.hr{text-align:center;padding:8px 2px;border-radius:16px;cursor:pointer;border:1px solid transparent;background:none;color:#fff;font:inherit;transition:background .2s,border-color .2s,transform .15s}
.hr:active{transform:scale(.94)}
.hr.on{background:rgba(255,255,255,.22);border-color:rgba(255,255,255,.34);transform:translateY(-2px)}
.hr .l{font-size:11px;color:rgba(255,255,255,.88);font-weight:500}
.hr .ic{font-size:18px;margin:7px 0 6px;display:block;line-height:1}
.hr .t{font-size:15px;font-weight:500;font-variant-numeric:tabular-nums}
</style></head><body>
<div id="sky" class="sky rain">
  <svg class="bg" viewBox="0 0 360 280" preserveAspectRatio="xMidYMid slice">
    <defs>
      <radialGradient id="cg" cx="40%" cy="30%" r="70%"><stop offset="0%" stop-color="#fff"/><stop offset="100%" stop-color="#d7e0ec"/></radialGradient>
      <radialGradient id="sg" cx="50%" cy="50%" r="50%"><stop offset="0%" stop-color="#ffe985"/><stop offset="100%" stop-color="#ffb347"/></radialGradient>
    </defs>
    <g id="g_sun" class="sunorb" style="display:none">
      <circle cx="278" cy="66" r="32" fill="url(#sg)"/>
      <circle cx="278" cy="66" r="46" fill="#ffe985" fill-opacity=".2"/>
      <g stroke="#ffe985" stroke-opacity=".55" stroke-width="2.5" stroke-linecap="round">
        <line x1="278" y1="14" x2="278" y2="24"/><line x1="278" y1="108" x2="278" y2="118"/>
        <line x1="226" y1="66" x2="236" y2="66"/><line x1="320" y1="66" x2="330" y2="66"/>
        <line x1="241" y1="29" x2="248" y2="36"/><line x1="308" y1="96" x2="315" y2="103"/>
        <line x1="315" y1="29" x2="308" y2="36"/><line x1="248" y1="96" x2="241" y2="103"/>
      </g>
    </g>
    <g id="g_cld" class="cld">
      <ellipse cx="165" cy="116" rx="34" ry="12" fill="#fff" fill-opacity=".3"/>
      <circle cx="152" cy="106" r="16" fill="url(#cg)" fill-opacity=".92"/>
      <circle cx="176" cy="110" r="12" fill="url(#cg)" fill-opacity=".9"/>
      <rect x="190" y="112" width="110" height="28" rx="14" fill="url(#cg)"/>
      <circle cx="222" cy="96" r="28" fill="url(#cg)"/>
      <circle cx="264" cy="100" r="22" fill="url(#cg)"/>
    </g>
    <g id="g_rain"><g id="rf1" class="rf1" stroke="#fff" stroke-opacity=".34" stroke-width="1.8" stroke-linecap="round"></g><g id="rf2" class="rf2" stroke="#fff" stroke-opacity=".24" stroke-width="1.5" stroke-linecap="round"></g></g>
  </svg>
  <div class="top">
    <button id="city" class="citybtn">Barcelona <span class="tag">next ↻</span></button>
    <div class="acts">
      <button id="ref" class="pbtn">↻</button>
      <button id="unit" class="pbtn">°C</button>
    </div>
  </div>
  <div id="bigt" class="bigt">18°</div>
  <div id="sub" class="sub">Now · Light rain</div>
  <div id="hl" class="hl">H:24° L:15°</div>
  <div id="dock" class="dock"></div>
</div>
<script>
const T=window.TelegramWidget,$=id=>document.getElementById(id);
let S={i:0,city:'Barcelona',sky:'rain',hi:24,lo:15,gps:false,hrs:[{l:'Now',ic:'🌧️',t:18,c:'Light rain',s:'rain'},{l:'1 PM',ic:'🌦️',t:19,c:'Showers',s:'rain'},{l:'2 PM',ic:'☁️',t:21,c:'Cloudy',s:'cloud'},{l:'3 PM',ic:'⛅',t:23,c:'Partly sunny',s:'cloud'},{l:'4 PM',ic:'☀️',t:24,c:'Sunny',s:'sun'},{l:'5 PM',ic:'☀️',t:23,c:'Clear',s:'sun'}]};
let sel=0,far=false,busy=false,curT=18,animId=0;
const cv=c=>far?Math.round(c*9/5+32):c;
function fillRain(el,n,step){
  let h='';
  for(let i=0;i<n;i++){
    const x=16+(i*step)%340,y=(i*47)%280;
    h+='<line x1="'+x+'" y1="'+y+'" x2="'+(x-3)+'" y2="'+(y+14)+'"/>';
    h+='<line x1="'+x+'" y1="'+(y-280)+'" x2="'+(x-3)+'" y2="'+(y-266)+'"/>';
  }
  el.innerHTML=h;
}
fillRain($('rf1'),12,31);
fillRain($('rf2'),10,37);
function animTemp(to){
  const from=curT,t0=Date.now(),my=++animId;
  if(from===to){$('bigt').textContent=to+'°';return}
  (function step(){
    if(my!==animId)return;
    const p=Math.min(1,(Date.now()-t0)/320),e=1-Math.pow(1-p,3);
    curT=Math.round(from+(to-from)*e);
    $('bigt').textContent=curT+'°';
    if(p<1)requestAnimationFrame(step);else curT=to;
  })();
}
function draw(animate){
  const h=S.hrs[sel]||S.hrs[0],sk=(h&&h.s)||S.sky||'rain';
  $('sky').className='sky '+sk;
  $('g_sun').style.display=(sk==='sun'||sk==='cloud')?'block':'none';
  $('g_cld').style.display=(sk==='sun'||sk==='night')?'none':'block';
  $('g_rain').style.display=sk==='rain'?'block':'none';
  $('city').innerHTML=(S.gps?'📍 ':'')+S.city+' <span class="tag">next</span>';
  $('unit').textContent=far?'°F':'°C';
  const target=cv(h.t);
  if(animate)animTemp(target);else{curT=target;$('bigt').textContent=target+'°'}
  $('sub').textContent=h.l+' · '+h.c;
  $('hl').textContent='H:'+cv(S.hi)+'° L:'+cv(S.lo)+'°';
  const dk=$('dock');dk.innerHTML='';
  S.hrs.forEach((x,i)=>{
    const b=document.createElement('button');
    b.className='hr'+(i===sel?' on':'');
    b.innerHTML='<div class="l">'+x.l+'</div><span class="ic">'+x.ic+'</span><div class="t">'+cv(x.t)+'°</div>';
    b.addEventListener('click',()=>{sel=i;draw(true);try{T.HapticFeedback.selectionChanged()}catch(e){}});
    dk.appendChild(b);
  });
}
function load(raw,anim){try{if(raw){const d=JSON.parse(raw);if(d.city&&Array.isArray(d.hrs)){S=d;sel=0}}}catch(e){}draw(!!anim)}
function ask(cmd){
  if(busy||!T||!T.request)return;busy=true;$('ref').textContent='…';
  try{T.HapticFeedback.impactOccurred('light')}catch(e){}
  T.request(cmd).then(r=>{
    if(r&&r.result)load(r.result,true);
    try{T.HapticFeedback.notificationOccurred('success')}catch(e){}
  }).finally(()=>{busy=false;$('ref').textContent='↻'});
}
$('city').addEventListener('click',()=>ask('weather:next:'+S.i));
$('ref').addEventListener('click',()=>ask('weather:ref:'+S.i));
$('unit').addEventListener('click',()=>{far=!far;draw(true);try{T.HapticFeedback.selectionChanged()}catch(e){}});
if(T){load(T.chatData,false);T.onEvent('chatDataChanged',()=>load(T.chatData,true))}else draw(false);
</script></body></html>`
}

type diceState struct {
	D1    int   `json:"d1"`
	D2    int   `json:"d2"`
	Rolls int   `json:"n"`
	Total int   `json:"tot"`
	Hist  []int `json:"h"`
}

var (
	diceMu     sync.Mutex
	diceStates = map[string]diceState{}
)

func cmdDice(b *gotgbot.Bot, ctx *ext.Context) error {
	raw, _ := json.Marshal(diceState{D1: 5, D2: 0, Rolls: 0, Total: 0, Hist: []int{}})
	_, err := b.SendWidget(ctx.EffectiveChat.Id, diceHTML(), &gotgbot.SendWidgetOpts{
		AllowQueries: true,
		ChatData:     string(raw),
	})
	return err
}

func handleDiceQuery(b *gotgbot.Bot, ctx *ext.Context) error {
	wq := ctx.WidgetQuery
	if wq == nil {
		return nil
	}
	key := widgetKey(wq)
	diceMu.Lock()
	st := diceStates[key]
	if st.D1 == 0 {
		st.D1 = 5
	}
	switch wq.Data {
	case "dice:reset":
		st = diceState{D1: 5, D2: 0, Rolls: 0, Total: 0, Hist: []int{}}
	case "dice:roll2":
		st.D1 = rand.Intn(6) + 1
		st.D2 = rand.Intn(6) + 1
		st.Rolls++
		st.Total += st.D1 + st.D2
		st.Hist = append([]int{st.D1 + st.D2}, st.Hist...)
	default:
		st.D1 = rand.Intn(6) + 1
		st.D2 = 0
		st.Rolls++
		st.Total += st.D1
		st.Hist = append([]int{st.D1}, st.Hist...)
	}
	if len(st.Hist) > 6 {
		st.Hist = st.Hist[:6]
	}
	diceStates[key] = st
	diceMu.Unlock()

	raw, _ := json.Marshal(st)
	persistChatData(b, wq, string(raw))
	return answerJSON(b, wq.Id, st)
}

func diceHTML() string {
	return `<!DOCTYPE html><html><head><meta charset="utf-8"><style>` + fullBleedCSS + `
.felt{background:linear-gradient(155deg,#102a2e,#164e45);color:#fff}
.stage{display:flex;justify-content:center;gap:14px;padding:12px 0 8px}
.die{width:76px;height:76px;border-radius:18px;background:linear-gradient(145deg,#fff,#e2e8f0);border:2px solid rgba(0,0,0,.12);cursor:pointer;transition:transform .16s}
.die:active{transform:scale(.92)}
.die.spin{transform:rotate(14deg) scale(1.06)}
.stats{display:grid;grid-template-columns:1fr 1fr;gap:8px;margin:10px 0}
.box{background:rgba(255,255,255,.11);border-radius:10px;padding:7px;text-align:center}
.hist{display:flex;justify-content:center;gap:5px;min-height:22px;margin-bottom:10px;flex-wrap:wrap}
.chip{background:rgba(255,255,255,.18);padding:2px 8px;border-radius:999px;font-size:11px;font-weight:700}
</style></head><body>
<div class="wrap felt">
  <div class="row">
    <span style="font-weight:700;font-size:15px">🎲 Dice Table</span>
    <span id="res" class="pill">Ready</span>
  </div>
  <div id="stage" class="stage"></div>
  <div class="hist" id="hist"></div>
  <div class="stats">
    <div class="box"><div class="muted" style="color:rgba(255,255,255,.72)">Rolls</div><div id="cnt" style="font-size:17px;font-weight:700">0</div></div>
    <div class="box"><div class="muted" style="color:rgba(255,255,255,.72)">Total</div><div id="tot" style="font-size:17px;font-weight:700">0</div></div>
  </div>
  <div class="row" style="padding-right:42px">
    <button id="r1" class="btn" style="background:#10b981">Roll 1</button>
    <button id="r2" class="btn" style="background:#0ea5e9">Roll 2</button>
    <button id="rst" class="btn sec" style="flex:0 0 auto;padding:9px 12px">↺</button>
  </div>
</div>
<script>
const T=window.TelegramWidget,$=id=>document.getElementById(id);
const P={1:[[50,50]],2:[[28,28],[72,72]],3:[[28,28],[50,50],[72,72]],4:[[28,28],[72,28],[28,72],[72,72]],5:[[28,28],[72,28],[50,50],[28,72],[72,72]],6:[[28,24],[72,24],[28,50],[72,50],[28,76],[72,76]]};
let S={d1:5,d2:0,n:0,tot:0,h:[]},busy=false,spinning=false;
const svgDie=v=>'<svg class="die'+(spinning?' spin':'')+'" viewBox="0 0 100 100">'+(P[v]||P[1]).map(([x,y])=>'<circle cx="'+x+'" cy="'+y+'" r="9" fill="#0f172a"/>').join('')+'</svg>';
function draw(){
  $('stage').innerHTML=svgDie(S.d1)+(S.d2>0?svgDie(S.d2):'');
  $('stage').querySelectorAll('.die').forEach(el=>el.addEventListener('click',()=>roll('dice:roll1')));
  $('res').textContent=S.n===0?'Tap die':(S.d2>0?S.d1+' + '+S.d2+' = '+(S.d1+S.d2):'Rolled '+S.d1);
  $('cnt').textContent=S.n;$('tot').textContent=S.tot;
  $('hist').innerHTML=(S.h||[]).map(v=>'<span class="chip">'+v+'</span>').join('');
}
function load(raw){try{if(raw)S=JSON.parse(raw)}catch(e){}draw()}
function roll(cmd){
  if(busy||!T||!T.request)return;busy=true;spinning=true;
  try{T.HapticFeedback.impactOccurred('medium')}catch(e){}
  let k=0,tm=setInterval(()=>{k++;S.d1=(k%6)+1;if(cmd==='dice:roll2')S.d2=((k+2)%6)+1;draw()},65);
  T.request(cmd).then(r=>{
    clearInterval(tm);spinning=false;if(r&&r.result)load(r.result);
    try{T.HapticFeedback.notificationOccurred('success')}catch(e){}
  }).catch(()=>{clearInterval(tm);spinning=false;draw()}).finally(()=>{busy=false});
}
$('r1').addEventListener('click',()=>roll('dice:roll1'));
$('r2').addEventListener('click',()=>roll('dice:roll2'));
$('rst').addEventListener('click',()=>roll('dice:reset'));
if(T){load(T.chatData);T.onEvent('chatDataChanged',()=>load(T.chatData))}else draw();
</script></body></html>`
}

func cmdTabs(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := b.SendWidget(ctx.EffectiveChat.Id, tabsHTML(), &gotgbot.SendWidgetOpts{
		AllowQueries: true,
		ChatData:     `{"tab":0,"msg":"Tap Ping to test roundtrip"}`,
	})
	return err
}

func handleTabQuery(b *gotgbot.Bot, ctx *ext.Context) error {
	wq := ctx.WidgetQuery
	if wq == nil {
		return nil
	}
	parts := strings.Split(wq.Data, ":")
	tab := 0
	if len(parts) >= 3 {
		tab, _ = strconv.Atoi(parts[2])
	}
	res := map[string]any{
		"tab": tab,
		"msg": "Pong from Go · " + time.Now().Format("15:04:05"),
	}
	raw, _ := json.Marshal(res)
	persistChatData(b, wq, string(raw))
	return answerJSON(b, wq.Id, res)
}

func tabsHTML() string {
	return `<!DOCTYPE html><html><head><meta charset="utf-8"><style>` + fullBleedCSS + `
.card{background:var(--tg-theme-bg-color,#17212b);color:var(--tg-theme-text-color,#f5f5f5)}
.seg{display:grid;grid-template-columns:repeat(4,1fr);gap:3px;padding:3px;border-radius:10px;background:var(--tg-theme-secondary-bg-color,#232e3c);margin-bottom:12px}
.tab{padding:6px 2px;border:none;border-radius:8px;background:transparent;color:var(--tg-theme-hint-color,#8b9bb4);font-size:12px;font-weight:700;cursor:pointer;transition:background .2s,color .2s}
.tab.on{background:var(--tg-theme-button-color,#3390ec);color:var(--tg-theme-button-text-color,#fff)}
.pane{display:none;font-size:13px;line-height:1.5;min-height:105px}
.pane.on{display:block}
.kv{display:flex;justify-content:space-between;padding:7px 0;border-bottom:1px solid rgba(128,128,128,.16)}
.kv:last-child{border-bottom:none}
.code{background:var(--tg-theme-secondary-bg-color,#232e3c);padding:8px 10px;border-radius:9px;font-family:monospace;font-size:11px;margin:8px 0;user-select:text;-webkit-user-select:text}
</style></head><body>
<div class="wrap card">
  <div class="seg">
    <button class="tab on" data-i="0">Overview</button>
    <button class="tab" data-i="1">Specs</button>
    <button class="tab" data-i="2">SDK</button>
    <button class="tab" data-i="3">Status</button>
  </div>
  <div class="pane on" id="p0">
    <div style="font-size:16px;font-weight:700;margin-bottom:4px">Telegram Bot API 10.4</div>
    <p class="muted" style="margin-bottom:10px">Self-contained HTML widgets rendered natively inside chat bubbles with QuickJS sandboxing.</p>
    <div style="padding-right:42px"><button id="ping" class="btn" style="width:100%">📡 Ping Go Backend</button></div>
  </div>
  <div class="pane" id="p1">
    <div class="kv"><span class="muted">Max HTML</span><b>32 KB</b></div>
    <div class="kv"><span class="muted">Max zlib-9</span><b>4,096 B</b></div>
    <div class="kv"><span class="muted">Query Payload</span><b>1–64 B</b></div>
    <div class="kv"><span class="muted">Answer Payload</span><b>0–1,024 B</b></div>
  </div>
  <div class="pane" id="p2">
    <div class="muted">Send a query from JS to Go:</div>
    <div class="code" id="snip">TelegramWidget.request('tab:ping:2')</div>
    <div style="padding-right:42px"><button id="cpy" class="btn sec" style="width:100%;background:var(--tg-theme-secondary-bg-color,#232e3c)">📋 Copy Snippet</button></div>
  </div>
  <div class="pane" id="p3">
    <div style="font-weight:700;margin-bottom:4px">Live Connection</div>
    <div id="st" class="muted">Ready</div>
  </div>
</div>
<script>
const T=window.TelegramWidget,$=id=>document.getElementById(id);
const tabs=[...document.querySelectorAll('.tab')],panes=[...document.querySelectorAll('.pane')];
let cur=0;
function show(i){cur=i;tabs.forEach((b,k)=>b.classList.toggle('on',k===i));panes.forEach((p,k)=>p.classList.toggle('on',k===i))}
tabs.forEach(b=>b.addEventListener('click',()=>{show(+b.dataset.i);try{T.HapticFeedback.selectionChanged()}catch(e){}}));
$('ping').addEventListener('click',()=>{
  if(!T||!T.request)return;$('ping').textContent='Pinging...';
  T.request('tab:ping:'+cur).then(r=>{
    const d=JSON.parse(r.result||'{}');
    $('ping').textContent='✓ '+(d.msg||'Pong!');
    $('st').textContent=d.msg||'Pong!';
    try{T.HapticFeedback.notificationOccurred('success')}catch(e){}
  }).catch(()=>{$('ping').textContent='Retry Ping'});
});
$('cpy').addEventListener('click',()=>{
  try{T.copyText($('snip').textContent);$('cpy').textContent='✓ Copied!'}catch(e){}
});
if(T&&T.chatData){try{const d=JSON.parse(T.chatData);if(typeof d.tab==='number')show(d.tab);if(d.msg)$('st').textContent=d.msg}catch(e){}}
</script></body></html>`
}

func cmdTime(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := b.SendWidget(ctx.EffectiveChat.Id, timeHTML(), &gotgbot.SendWidgetOpts{
		AllowQueries: true,
	})
	return err
}

func handleTimeQuery(b *gotgbot.Bot, ctx *ext.Context) error {
	wq := ctx.WidgetQuery
	if wq == nil {
		return nil
	}
	now := time.Now().UTC()
	return answerJSON(b, wq.Id, map[string]any{
		"ms":  now.UnixMilli(),
		"srv": now.Format("15:04:05 UTC"),
	})
}

func timeHTML() string {
	return `<!DOCTYPE html><html><head><meta charset="utf-8"><style>` + fullBleedCSS + `
.night{background:linear-gradient(150deg,#09142c,#17356e 60%,#2c5fc4);color:#fff}
.big{font-size:46px;font-weight:700;letter-spacing:1px;line-height:1.05;font-family:ui-monospace,monospace;font-variant-numeric:tabular-nums;margin-top:4px}
.zones{display:grid;grid-template-columns:repeat(4,1fr);gap:6px;margin:12px 0}
.z{background:rgba(255,255,255,.11);border-radius:10px;padding:7px 4px;text-align:center}
.z span{display:block;font-size:10px;color:rgba(255,255,255,.72)}
.z b{font-size:13px;font-family:ui-monospace,monospace;font-variant-numeric:tabular-nums}
.sw{background:rgba(0,0,0,.22);border-radius:12px;padding:9px 12px;margin-right:42px}
.swv{font-size:20px;font-weight:700;font-family:ui-monospace,monospace;font-variant-numeric:tabular-nums}
</style></head><body>
<div class="wrap night">
  <div class="row">
    <span id="dt" style="font-size:12px;color:rgba(255,255,255,.78)">Loading...</span>
    <div class="row" style="gap:5px">
      <button id="fmt" class="pill">24H</button>
      <button id="sync" class="pill">Sync</button>
    </div>
  </div>
  <div id="clk" class="big">--:--:--</div>
  <div class="zones">
    <div class="z"><span>UTC</span><b id="z0">--:--</b></div>
    <div class="z"><span>India</span><b id="z1">--:--</b></div>
    <div class="z"><span>New York</span><b id="z2">--:--</b></div>
    <div class="z"><span>Tokyo</span><b id="z3">--:--</b></div>
  </div>
  <div class="sw row">
    <div><div style="font-size:10px;color:rgba(255,255,255,.68)">STOPWATCH</div><div id="swv" class="swv">00:00.0</div></div>
    <div class="row" style="gap:6px">
      <button id="go" class="pill" style="background:#38bdf8;color:#082f49;padding:6px 12px">Start</button>
      <button id="rst" class="pill" style="padding:6px 10px">↺</button>
    </div>
  </div>
</div>
<script>
const T=window.TelegramWidget,$=id=>document.getElementById(id);
let h12=false,skew=0,run=false,el=0,t0=0;
const fz=(d,tz)=>new Intl.DateTimeFormat('en-US',{timeZone:tz,hour:'2-digit',minute:'2-digit',hour12:h12}).format(d);
function tick(){
  const d=new Date(Date.now()+skew);
  $('clk').textContent=new Intl.DateTimeFormat(undefined,{hour:'2-digit',minute:'2-digit',second:'2-digit',hour12:h12}).format(d);
  $('dt').textContent=new Intl.DateTimeFormat(undefined,{weekday:'short',month:'short',day:'numeric'}).format(d);
  $('z0').textContent=fz(d,'UTC');$('z1').textContent=fz(d,'Asia/Kolkata');
  $('z2').textContent=fz(d,'America/New_York');$('z3').textContent=fz(d,'Asia/Tokyo');
  let ms=el+(run?Date.now()-t0:0),cs=Math.floor(ms/100)%10,s=Math.floor(ms/1000)%60,m=Math.floor(ms/60000);
  $('swv').textContent=String(m).padStart(2,'0')+':'+String(s).padStart(2,'0')+'.'+cs;
}
setInterval(tick,100);tick();
$('fmt').addEventListener('click',()=>{h12=!h12;$('fmt').textContent=h12?'12H':'24H';tick();try{T.HapticFeedback.selectionChanged()}catch(e){}});
$('go').addEventListener('click',()=>{
  if(run){el+=Date.now()-t0;run=false;$('go').textContent='Start'}
  else{t0=Date.now();run=true;$('go').textContent='Pause'}
  tick();try{T.HapticFeedback.impactOccurred('light')}catch(e){}
});
$('rst').addEventListener('click',()=>{run=false;el=0;$('go').textContent='Start';tick()});
$('sync').addEventListener('click',()=>{
  if(!T||!T.request)return;const s=Date.now();$('sync').textContent='...';
  T.request('time:sync').then(r=>{
    const d=JSON.parse(r.result||'{}'),rtt=Date.now()-s;
    if(d.ms)skew=d.ms+Math.round(rtt/2)-Date.now();
    $('sync').textContent=rtt+'ms';tick();
    try{T.HapticFeedback.notificationOccurred('success')}catch(e){}
  }).catch(()=>{$('sync').textContent='Sync'});
});
</script></body></html>`
}

var (
	counterMu     sync.Mutex
	counterStates = map[string]int{}
)

func cmdCounter(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := b.SendWidget(ctx.EffectiveChat.Id, counterHTML(), &gotgbot.SendWidgetOpts{
		AllowQueries: true,
		ChatData:     `{"v":0}`,
	})
	return err
}

func handleCounterQuery(b *gotgbot.Bot, ctx *ext.Context) error {
	wq := ctx.WidgetQuery
	if wq == nil {
		return nil
	}
	key := widgetKey(wq)
	counterMu.Lock()
	val := counterStates[key]
	if wq.Data == "counter:reset" {
		val = 0
	} else if strings.HasPrefix(wq.Data, "counter:add:") {
		delta, _ := strconv.Atoi(strings.TrimPrefix(wq.Data, "counter:add:"))
		val += delta
	}
	counterStates[key] = val
	counterMu.Unlock()

	res := map[string]int{"v": val}
	raw, _ := json.Marshal(res)
	persistChatData(b, wq, string(raw))
	return answerJSON(b, wq.Id, res)
}

func counterHTML() string {
	return `<!DOCTYPE html><html><head><meta charset="utf-8"><style>` + fullBleedCSS + `
.card{background:var(--tg-theme-bg-color,#17212b);color:var(--tg-theme-text-color,#fff)}
.val{font-size:62px;font-weight:800;text-align:center;margin:10px 0 14px;font-variant-numeric:tabular-nums;transition:transform .12s}
.pos{color:#22c55e}.neg{color:#ef4444}
</style></head><body>
<div class="wrap card">
  <div class="row"><span class="muted">COUNTER</span><button id="rst" class="pill" style="background:var(--tg-theme-secondary-bg-color,#232e3c)">Reset</button></div>
  <div id="val" class="val">0</div>
  <div class="row" style="padding-right:42px">
    <button class="btn sec" style="background:var(--tg-theme-secondary-bg-color,#232e3c)" data-d="-5">−5</button>
    <button class="btn sec" style="background:var(--tg-theme-secondary-bg-color,#232e3c)" data-d="-1">−1</button>
    <button class="btn" data-d="1">+1</button>
    <button class="btn" data-d="5">+5</button>
  </div>
</div>
<script>
const T=window.TelegramWidget,$=id=>document.getElementById(id);
let v=0,busy=false;
function draw(){const el=$('val');el.textContent=v;el.className='val '+(v>0?'pos':v<0?'neg':'')}
function load(raw){try{if(raw){const d=JSON.parse(raw);if(typeof d.v==='number')v=d.v}}catch(e){}draw()}
function send(cmd){
  if(busy||!T||!T.request)return;busy=true;
  try{T.HapticFeedback.impactOccurred('light')}catch(e){}
  T.request(cmd).then(r=>{if(r&&r.result)load(r.result)}).finally(()=>{busy=false});
}
document.querySelectorAll('[data-d]').forEach(b=>b.addEventListener('click',()=>send('counter:add:'+b.dataset.d)));
$('rst').addEventListener('click',()=>send('counter:reset'));
if(T){load(T.chatData);T.onEvent('chatDataChanged',()=>load(T.chatData))}else draw();
</script></body></html>`
}

type tttState struct {
	Board  [9]string `json:"b"`
	Turn   string    `json:"t"`
	Winner string    `json:"w"`
	VsBot  bool      `json:"ai"`
}

var (
	tttMu     sync.Mutex
	tttStates = map[string]tttState{}
)

func newTTT(vsBot bool) tttState { return tttState{Turn: "X", VsBot: vsBot} }

func checkTTTWinner(b [9]string) string {
	lines := [][3]int{{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, {0, 3, 6}, {1, 4, 7}, {2, 5, 8}, {0, 4, 8}, {2, 4, 6}}
	for _, l := range lines {
		if b[l[0]] != "" && b[l[0]] == b[l[1]] && b[l[1]] == b[l[2]] {
			return b[l[0]]
		}
	}
	for _, c := range b {
		if c == "" {
			return ""
		}
	}
	return "draw"
}

func botMoveTTT(s *tttState) {
	if s.Winner != "" {
		return
	}
	for _, mark := range []string{"O", "X"} {
		for i := range 9 {
			if s.Board[i] == "" {
				s.Board[i] = mark
				if checkTTTWinner(s.Board) == mark {
					s.Board[i] = "O"
					s.Winner = checkTTTWinner(s.Board)
					s.Turn = "X"
					return
				}
				s.Board[i] = ""
			}
		}
	}
	if s.Board[4] == "" {
		s.Board[4] = "O"
	} else {
		var free []int
		for i, c := range s.Board {
			if c == "" {
				free = append(free, i)
			}
		}
		if len(free) > 0 {
			s.Board[free[rand.Intn(len(free))]] = "O"
		}
	}
	s.Winner = checkTTTWinner(s.Board)
	s.Turn = "X"
}

func cmdTicTacToe(b *gotgbot.Bot, ctx *ext.Context) error {
	raw, _ := json.Marshal(newTTT(true))
	_, err := b.SendWidget(ctx.EffectiveChat.Id, tttHTML(), &gotgbot.SendWidgetOpts{
		AllowQueries: true,
		ChatData:     string(raw),
	})
	return err
}

func handleTTTQuery(b *gotgbot.Bot, ctx *ext.Context) error {
	wq := ctx.WidgetQuery
	if wq == nil {
		return nil
	}
	key := widgetKey(wq)
	tttMu.Lock()
	st, ok := tttStates[key]
	if !ok {
		st = newTTT(true)
	}
	switch {
	case wq.Data == "ttt:reset":
		st = newTTT(st.VsBot)
	case wq.Data == "ttt:mode":
		st = newTTT(!st.VsBot)
	case strings.HasPrefix(wq.Data, "ttt:m:"):
		cell, err := strconv.Atoi(strings.TrimPrefix(wq.Data, "ttt:m:"))
		if err == nil && cell >= 0 && cell < 9 && st.Winner == "" && st.Board[cell] == "" {
			st.Board[cell] = st.Turn
			st.Winner = checkTTTWinner(st.Board)
			if st.Winner == "" {
				if st.VsBot {
					botMoveTTT(&st)
				} else if st.Turn == "X" {
					st.Turn = "O"
				} else {
					st.Turn = "X"
				}
			}
		}
	}
	tttStates[key] = st
	tttMu.Unlock()

	raw, _ := json.Marshal(st)
	persistChatData(b, wq, string(raw))
	return answerJSON(b, wq.Id, st)
}

func tttHTML() string {
	return `<!DOCTYPE html><html><head><meta charset="utf-8"><style>` + fullBleedCSS + `
.card{background:var(--tg-theme-bg-color,#17212b);color:var(--tg-theme-text-color,#fff)}
.board{display:grid;grid-template-columns:repeat(3,1fr);gap:8px;max-width:228px;margin:12px auto}
.cell{height:66px;border:none;border-radius:12px;background:var(--tg-theme-secondary-bg-color,#232e3c);color:inherit;font-size:28px;font-weight:800;cursor:pointer;transition:transform .12s}
.cell:active{transform:scale(.93)}
.cell.x{color:#f43f5e}.cell.o{color:#14b8a6}
</style></head><body>
<div class="wrap card">
  <div class="row">
    <b id="info" style="font-size:15px">Turn: X</b>
    <button id="mode" class="pill" style="background:var(--tg-theme-secondary-bg-color,#232e3c)">🤖 vs Bot</button>
  </div>
  <div id="board" class="board"></div>
  <div style="padding-right:42px"><button id="rst" class="btn sec" style="width:100%;background:var(--tg-theme-secondary-bg-color,#232e3c)">🔄 New Game</button></div>
</div>
<script>
const T=window.TelegramWidget,$=id=>document.getElementById(id);
let S={b:['','','','','','','','',''],t:'X',w:'',ai:true},busy=false;
const brd=$('board');
for(let i=0;i<9;i++){const c=document.createElement('button');c.className='cell';c.addEventListener('click',()=>send('ttt:m:'+i));brd.appendChild(c)}
function draw(){
  [...brd.children].forEach((c,i)=>{const v=S.b[i]||'';c.textContent=v;c.className='cell'+(v==='X'?' x':v==='O'?' o':'')});
  $('mode').textContent=S.ai?'🤖 vs Bot':'👥 2 Players';
  $('info').textContent=S.w==='draw'?'🤝 Draw!':S.w?('🏆 '+S.w+' wins!'):'Turn: '+S.t;
}
function load(raw){try{if(raw)S=JSON.parse(raw)}catch(e){}draw()}
function send(cmd){
  if(busy||!T||!T.request)return;busy=true;
  try{T.HapticFeedback.impactOccurred('light')}catch(e){}
  T.request(cmd).then(r=>{if(r&&r.result)load(r.result);if(S.w)try{T.HapticFeedback.notificationOccurred('success')}catch(e){}}).finally(()=>{busy=false});
}
$('mode').addEventListener('click',()=>send('ttt:mode'));
$('rst').addEventListener('click',()=>send('ttt:reset'));
if(T){load(T.chatData);T.onEvent('chatDataChanged',()=>load(T.chatData))}else draw();
</script></body></html>`
}

type pollState struct {
	Question string   `json:"q"`
	Options  []string `json:"o"`
	Votes    []int    `json:"v"`
	Rev      int      `json:"r"`
}

var (
	pollMu        sync.Mutex
	pollStateData = pollState{
		Question: "Tabs or spaces?",
		Options:  []string{"Tabs (go fmt!)", "2 Spaces", "4 Spaces"},
		Votes:    []int{18, 7, 11},
		Rev:      1,
	}
	pollUserVotes = map[int64]int{}
)

func cmdPoll(b *gotgbot.Bot, ctx *ext.Context) error {
	pollMu.Lock()
	raw, _ := json.Marshal(pollStateData)
	pollMu.Unlock()
	_, err := b.SendWidget(ctx.EffectiveChat.Id, pollHTML(), &gotgbot.SendWidgetOpts{
		AllowQueries: true,
		IsShared:     true,
		PublicData:   string(raw),
	})
	return err
}

func processPollVote(b *gotgbot.Bot, userID int64, queryID, widgetID, data string) error {
	idx, err := strconv.Atoi(strings.TrimPrefix(data, "poll:vote:"))
	if err != nil {
		return nil
	}
	pollMu.Lock()
	if idx >= 0 && idx < len(pollStateData.Votes) {
		if prev, had := pollUserVotes[userID]; had {
			if prev == idx {
				pollStateData.Votes[idx]--
				delete(pollUserVotes, userID)
				idx = -1
			} else {
				pollStateData.Votes[prev]--
				pollStateData.Votes[idx]++
				pollUserVotes[userID] = idx
			}
		} else {
			pollStateData.Votes[idx]++
			pollUserVotes[userID] = idx
		}
		pollStateData.Rev++
	}
	pubRaw, _ := json.Marshal(pollStateData)
	reply := map[string]any{"my": idx, "v": append([]int(nil), pollStateData.Votes...), "r": pollStateData.Rev}
	pollMu.Unlock()

	if widgetID != "" {
		_, _ = b.EditSharedWidgetData(widgetID, string(pubRaw), nil)
	}
	return answerJSON(b, queryID, reply)
}

func handlePollQuery(b *gotgbot.Bot, ctx *ext.Context) error {
	if ctx.WidgetQuery == nil {
		return nil
	}
	return processPollVote(b, ctx.WidgetQuery.From.Id, ctx.WidgetQuery.Id, ctx.WidgetQuery.WidgetId, ctx.WidgetQuery.Data)
}

func handleSharedPollQuery(b *gotgbot.Bot, ctx *ext.Context) error {
	if ctx.SharedWidgetQuery == nil {
		return nil
	}
	return processPollVote(b, ctx.SharedWidgetQuery.From.Id, ctx.SharedWidgetQuery.Id, ctx.SharedWidgetQuery.WidgetId, ctx.SharedWidgetQuery.Data)
}

func pollHTML() string {
	return `<!DOCTYPE html><html><head><meta charset="utf-8"><style>` + fullBleedCSS + `
.card{background:var(--tg-theme-bg-color,#17212b);color:var(--tg-theme-text-color,#fff)}
.opt{width:100%;text-align:left;background:var(--tg-theme-secondary-bg-color,#232e3c);border:1.5px solid transparent;border-radius:11px;padding:9px 11px;margin-top:7px;color:inherit;cursor:pointer;transition:transform .12s,border-color .2s}
.opt:active{transform:scale(.98)}
.opt.my{border-color:var(--tg-theme-button-color,#3390ec)}
.bar{height:5px;border-radius:3px;background:rgba(128,128,128,.22);margin-top:6px;overflow:hidden}
.fill{height:100%;width:0%;background:var(--tg-theme-button-color,#3390ec);transition:width .35s ease}
</style></head><body>
<div class="wrap card">
  <div style="font-size:16px;font-weight:700" id="q">Poll</div>
  <div class="muted" id="sub" style="margin-top:2px;margin-bottom:4px">Live Poll · 0 votes</div>
  <div id="opts"></div>
</div>
<script>
const T=window.TelegramWidget,$=id=>document.getElementById(id);
let P={q:'Poll',o:[],v:[],r:0},my=-1,busy=false;
try{my=+(localStorage.getItem('my')??-1)}catch(e){}
function load(raw){try{if(raw){const d=JSON.parse(raw);if(Array.isArray(d.o))P=d}}catch(e){}draw()}
function draw(){
  $('q').textContent=P.q||'Poll';
  const sum=(P.v||[]).reduce((a,b)=>a+b,0);
  $('sub').textContent='Live Shared Poll · '+sum+(sum===1?' vote':' votes');
  const box=$('opts');box.innerHTML='';
  (P.o||[]).forEach((name,i)=>{
    const c=P.v[i]||0,pct=sum?Math.round(c*100/sum):0;
    const b=document.createElement('button');
    b.className='opt'+(i===my?' my':'');
    b.innerHTML='<div class="row"><span>'+(i===my?'✓ ':'')+name+'</span><b>'+pct+'% <span class="muted">('+c+')</span></b></div><div class="bar"><div class="fill" style="width:'+pct+'%"></div></div>';
    b.addEventListener('click',()=>vote(i));
    box.appendChild(b);
  });
}
function vote(i){
  if(busy||!T||!T.request)return;busy=true;
  try{T.HapticFeedback.impactOccurred('light')}catch(e){}
  T.request('poll:vote:'+i).then(r=>{
    const d=JSON.parse(r.result||'{}');
    if(Array.isArray(d.v)){P.v=d.v;P.r=d.r}
    if(typeof d.my==='number'){my=d.my;try{localStorage.setItem('my',my)}catch(e){}}
    draw();try{T.HapticFeedback.notificationOccurred('success')}catch(e){}
  }).finally(()=>{busy=false});
}
if(T){load(T.publicData);T.onEvent('publicDataChanged',()=>load(T.publicData));T.watchPublicData()}else draw();
</script></body></html>`
}

type dashState struct {
	Goroutines int    `json:"g"`
	MemMB      string `json:"m"`
	Uptime     string `json:"u"`
	Queries    int64  `json:"q"`
}

func makeDash() dashState {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return dashState{
		Goroutines: runtime.NumGoroutine(),
		MemMB:      fmt.Sprintf("%.1f MB", float64(ms.Alloc)/(1024*1024)),
		Uptime:     time.Since(startTime).Truncate(time.Second).String(),
		Queries:    queryCounter.Load(),
	}
}

func cmdDashboard(b *gotgbot.Bot, ctx *ext.Context) error {
	raw, _ := json.Marshal(makeDash())
	_, err := b.SendWidget(ctx.EffectiveChat.Id, dashboardHTML(), &gotgbot.SendWidgetOpts{
		AllowQueries: true,
		ChatData:     string(raw),
	})
	return err
}

func handleDashboardQuery(b *gotgbot.Bot, ctx *ext.Context) error {
	if ctx.WidgetQuery == nil {
		return nil
	}
	st := makeDash()
	raw, _ := json.Marshal(st)
	persistChatData(b, ctx.WidgetQuery, string(raw))
	return answerJSON(b, ctx.WidgetQuery.Id, st)
}

func dashboardHTML() string {
	return `<!DOCTYPE html><html><head><meta charset="utf-8"><style>` + fullBleedCSS + `
.card{background:linear-gradient(150deg,#111827,#1e293b);color:#fff}
.badge{background:#166534;color:#86efac;font-size:11px;padding:3px 9px;border-radius:999px;font-weight:700}
.kgrid{display:grid;grid-template-columns:1fr 1fr;gap:8px;margin:12px 0}
.kbox{background:rgba(255,255,255,.07);border-radius:10px;padding:9px;text-align:center}
.kval{font-size:17px;font-weight:800;color:#38bdf8;margin-top:2px;font-variant-numeric:tabular-nums}
</style></head><body>
<div class="wrap card">
  <div class="row"><b style="font-size:15px">⚡ Go Bot Runtime</b><span class="badge">● Online</span></div>
  <div class="kgrid">
    <div class="kbox"><div class="muted">Goroutines</div><div id="g" class="kval">-</div></div>
    <div class="kbox"><div class="muted">Heap Alloc</div><div id="m" class="kval">-</div></div>
    <div class="kbox"><div class="muted">Uptime</div><div id="u" class="kval">-</div></div>
    <div class="kbox"><div class="muted">Queries</div><div id="q" class="kval">-</div></div>
  </div>
  <div class="row" style="padding-right:42px">
    <button id="ref" class="btn">🔄 Refresh</button>
    <button id="png" class="btn sec">📶 Ping</button>
  </div>
</div>
<script>
const T=window.TelegramWidget,$=id=>document.getElementById(id);
let busy=false;
function draw(d){if(!d)return;$('g').textContent=d.g??'-';$('m').textContent=d.m??'-';$('u').textContent=d.u??'-';$('q').textContent=d.q??'-'}
function load(raw){try{if(raw)draw(JSON.parse(raw))}catch(e){}}
function act(c){
  if(busy||!T||!T.request)return;busy=true;
  T.request(c).then(r=>{if(r&&r.result)load(r.result);try{T.HapticFeedback.notificationOccurred('success')}catch(e){}}).finally(()=>{busy=false});
}
$('ref').addEventListener('click',()=>act('dash:ref'));
$('png').addEventListener('click',()=>act('dash:ping'));
if(T){load(T.chatData);T.onEvent('chatDataChanged',()=>load(T.chatData))}
</script></body></html>`
}
