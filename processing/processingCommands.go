package processing

import (
	"github.com/nicksnyder/go-i18n/i18n"
)

func GetBetDurationText(time int64, answersTag string, trans i18n.TranslateFunc) string {
	rulesData := map[string]interface{}{
		"Time": trans("hours", time),
	}

	return trans("rules_timer", rulesData)
}
