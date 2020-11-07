package chat

import (
)

type Chat interface {
	SendMessage(chatId int64, message string)
}
