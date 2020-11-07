package processing

import (
	"github.com/gameraccoon/telegram-i-told-you-bot/chat"
	"github.com/gameraccoon/telegram-i-told-you-bot/database"
	"github.com/nicksnyder/go-i18n/i18n"
	"time"
)

type StaticConfiguration struct {
	Language     string
	ExtendedLog  bool
	SleepSeconds int
}

type StaticProccessStructs struct {
	Chat       chat.Chat
	Db         *database.Database
	Timers     map[int64]time.Time
	Config     *StaticConfiguration
	Trans      i18n.TranslateFunc
}
