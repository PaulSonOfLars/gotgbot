package ext

import (
	"gotgbot/types"
	"encoding/json"
	"net/url"
	"github.com/pkg/errors"
)

// TODO: add cache_time, is_personal, next_offset, switch_pm_text, switch_pm_parameter arguments
func (b Bot) AnswerInlineQuery(inlineQueryId string, results []types.InlineQueryResult) (bool, error) {
	resultsStr, err := json.Marshal(results)
	if err != nil {
		return false, errors.Wrapf(err, "unable to unmarshal answerInlineQuery result")
	}
	v := url.Values{}
	v.Add("inline_query_id", inlineQueryId)
	v.Add("results", string(resultsStr))

	r, err := Get(b, "answerInlineQuery", v)
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