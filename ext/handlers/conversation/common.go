package conversation

import (
	"fmt"
	"strconv"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func StateKey(ctx *ext.Context, strategy KeyStrategy) string {
	switch strategy {
	case KeyStrategySender:
		return strconv.FormatInt(ctx.EffectiveSender.Id(), 10)
	case KeyStrategyChat:
		return strconv.FormatInt(ctx.EffectiveChat.Id, 10)
	case KeyStrategySenderAndChat:
		fallthrough
	default:
		// Default to KeyStrategySenderAndChat if unknown strategy
		return fmt.Sprintf("%d/%d", ctx.EffectiveSender.Id(), ctx.EffectiveChat.Id)
	}
}
