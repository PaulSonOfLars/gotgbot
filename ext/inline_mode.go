package ext

import (
	"github.com/PaulSonOfLars/gotgbot/types"
	"encoding/json"
	"net/url"
	"github.com/pkg/errors"
	"strconv"
)

type sendableAnswerInlineQuery struct {
	bot               Bot
	InlineQueryId     string
	Results           []types.InlineQueryResult
	CacheTime         int
	IsPersonal        bool
	NextOffset        string
	SwitchPmText      string
	SwitchPmParameter string
}

func (b Bot) NewSendableAnswerInlineQuery(inlineQueryId string, results []types.InlineQueryResult) *sendableAnswerInlineQuery {
	return &sendableAnswerInlineQuery{bot: b, InlineQueryId: inlineQueryId, Results: results}
}

func (aiq sendableAnswerInlineQuery) Send() (bool, error) {
	resultsStr, err := json.Marshal(aiq.Results)
	if err != nil {
		return false, errors.Wrapf(err, "unable to unmarshal answerInlineQuery result")
	}
	v := url.Values{}
	v.Add("inline_query_id", aiq.InlineQueryId)
	v.Add("results", string(resultsStr))
	v.Add("cache_time", strconv.Itoa(aiq.CacheTime))
	v.Add("is_personal", strconv.FormatBool(aiq.IsPersonal))
	v.Add("next_offset", aiq.NextOffset)
	v.Add("switch_pm_text", aiq.SwitchPmText)
	v.Add("switch_pm_parameter", aiq.SwitchPmParameter)

	r, err := Get(aiq.bot, "answerInlineQuery", v)
	if err != nil {
		return false, errors.Wrapf(err, "unable to execute answerInlineQuery request")
	}
	if !r.Ok {
		return false, errors.Wrapf(err, "invalid answerInlineQuery request")
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)
	return bb, nil
}

func (b Bot) AnswerInlineQuery(inlineQueryId string, results []types.InlineQueryResult) (bool, error) {
	return b.NewSendableAnswerInlineQuery(inlineQueryId, results).Send()
}
