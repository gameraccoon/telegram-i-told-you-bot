package main

import (
	// "github.com/gameraccoon/telegram-i-told-you-bot/database"
	"github.com/gameraccoon/telegram-i-told-you-bot/processing"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

type ProcessorFunc func(*processing.ProcessData)

type ProcessorFuncMap map[string]ProcessorFunc

type Processors struct {
	Group   ProcessorFuncMap
	Private ProcessorFuncMap
}

func sendResults(staticData *processing.StaticProccessStructs, questionId int64, chatIds []int64) {
	// variants := staticData.Db.GetQuestionVariants(questionId)
	// answers := staticData.Db.GetQuestionAnswers(questionId)
	// answersCount := staticData.Db.GetQuestionAnswersCount(questionId)

	// var buffer bytes.Buffer
	// buffer.WriteString(staticData.Trans("results_header"))
	// buffer.WriteString(fmt.Sprintf("<i>%s</i>", staticData.Db.GetQuestionText(questionId)))

	// for i, variant := range variants {
	// 	buffer.WriteString(fmt.Sprintf("\n%s - %d (%d%%)", variant, answers[i], int64(100.0*float32(answers[i])/float32(answersCount))))
	// }
	// resultText := buffer.String()

	// for _, chatId := range chatIds {
	// 	staticData.Chat.SendMessage(chatId, resultText)
	// }
}

func completeBet(staticData *processing.StaticProccessStructs, questionId int64) {
	// sendResults(staticData, questionId, chatIds)
}

func isBetExpired(staticData *processing.StaticProccessStructs, questionId int64) bool {
	if _, ok := staticData.Timers[questionId]; !ok {
		return true
	}

	return false
}

func processCompleteness(staticData *processing.StaticProccessStructs, questionId int64) {
	if isBetExpired(staticData, questionId) {
		completeBet(staticData, questionId)
	}
}

func betCommand(data *processing.ProcessData) {
	timeSeparator := strings.Index(data.Message, " ")

	var timeStr string
	var message string

	if timeSeparator != -1 {
		timeStr = data.Message[:timeSeparator]
		message = data.Message[timeSeparator+1:]
	} else {
		timeStr = data.Message
	}

	duration, isSuccessful, errorMessage := processing.ParseBetTime(timeStr)
	if isSuccessful {
		data.Static.Chat.SendMessage(data.ChatId, processing.GetBetDurationText(duration, data.Static.Trans)+" "+message)
	} else {
		data.Static.Chat.SendMessage(data.ChatId, errorMessage)
	}
}
func betsCommand(data *processing.ProcessData) {
	data.Static.Chat.SendMessage(data.ChatId, data.Static.Trans("test2"))
}
func mybetsCommand(data *processing.ProcessData) {
	data.Static.Chat.SendMessage(data.ChatId, data.Static.Trans("test3"))
}
func startCommand(data *processing.ProcessData) {
	data.Static.Chat.SendMessage(data.ChatId, data.Static.Trans("hello_message"))
}
func commandsListCommand(data *processing.ProcessData) {
	data.Static.Chat.SendMessage(data.ChatId, data.Static.Trans("commands_list"))
}

func makeGroupCommandProcessors() ProcessorFuncMap {
	return map[string]ProcessorFunc{
		"bet":    betCommand,
		"bets":   betsCommand,
		"mybets": mybetsCommand,
	}
}

func makePrivateCommandProcessors() ProcessorFuncMap {
	return map[string]ProcessorFunc{
		"start":    startCommand,
		"commands": commandsListCommand,
	}
}

func processCommand(data *processing.ProcessData, processors ProcessorFuncMap) bool {
	processor, ok := processors[data.Command]
	if ok {
		processor(data)
		return true
	}
	return false
}

func processUpdate(update *tgbotapi.Update, staticData *processing.StaticProccessStructs, processors *Processors) {
	data := processing.ProcessData{
		Static: staticData,
		ChatId: update.Message.Chat.ID,
		UserId: int64(update.Message.From.ID),
	}

	message := update.Message.Text

	var prefix string
	var selectedProcessors ProcessorFuncMap

	isPrivate := update.Message.Chat.IsPrivate()

	if isPrivate {
		prefix = "/"
		selectedProcessors = processors.Private
	} else {
		prefix = "#"
		selectedProcessors = processors.Group
	}

	if strings.HasPrefix(message, prefix) {
		commandLen := strings.Index(message, " ")
		if commandLen != -1 {
			data.Command = message[1:commandLen]
			data.Message = message[commandLen+1:]
		} else {
			data.Command = message[1:]
		}

		isProcessed := processCommand(&data, selectedProcessors)

		if isPrivate && !isProcessed {
			data.Static.Chat.SendMessage(data.ChatId, data.Static.Trans("warn_unknown_command"))
		}
	}
}

func processTimer(staticData *processing.StaticProccessStructs, questionId int64) {
	processCompleteness(staticData, questionId)
}
