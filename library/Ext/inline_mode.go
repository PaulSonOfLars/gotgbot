package Ext

import (
	"bot/library/Types"
	"encoding/json"
	"log"
	"net/url"
)

// TODO: add cache_time, is_personal, next_offset, switch_pm_text, switch_pm_parameter arguments
func (b Bot) AnswerInlineQuery(inline_query_id string, results []Types.InlineQueryResult) bool {
	results_str, err := json.Marshal(results)
	if err != nil {
		log.Println("err in inline query answer")
		log.Fatal(err)
	}
	v := url.Values{}
	v.Add("inline_query_id", inline_query_id)
	v.Add("results", string(results_str))

	r := Get(b, "setGameScore", v)
	if !r.Ok {
		log.Fatal("You done goofed, API Res for setGameScore was not OK")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}