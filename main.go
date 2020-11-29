package main

import (
	"encoding/json"
	"github.com/gameraccoon/telegram-i-told-you-bot/database"
	"github.com/gameraccoon/telegram-i-told-you-bot/processing"
	"github.com/gameraccoon/telegram-i-told-you-bot/telegramChat"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nicksnyder/go-i18n/i18n"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	i18n.MustLoadTranslationFile("./data/strings/en-us.all.json")
}

func getFileStringContent(filePath string) (content string, err error) {
	fileContent, err := ioutil.ReadFile(filePath)
	if err == nil {
		content = strings.TrimSpace(string(fileContent))
	}
	return
}

func getApiToken() (token string, err error) {
	return getFileStringContent("./telegramApiToken.txt")
}

func loadConfig(path string) (config processing.StaticConfiguration, err error) {
	jsonString, err := getFileStringContent(path)
	if err == nil {
		dec := json.NewDecoder(strings.NewReader(jsonString))
		err = dec.Decode(&config)
	}
	return
}

func updateTimers(staticData *processing.StaticProccessStructs, mutex *sync.Mutex) {
	bets := staticData.Db.GetActiveBets()

	mutex.Lock()
	for _, betId := range bets {
		_, staticData.Timers[betId], _ = staticData.Db.GetBetData(betId)
	}
	mutex.Unlock()

	for {
		currentTime := time.Now()
		mutex.Lock()
		for betId, endTime := range staticData.Timers {
			if endTime.Sub(currentTime).Seconds() < 0.0 {
				delete(staticData.Timers, betId)
				processTimer(staticData, betId)
			}
		}
		mutex.Unlock()
		time.Sleep(time.Duration(staticData.Config.SleepSeconds) * time.Second)
	}
}

func updateBot(bot *tgbotapi.BotAPI, staticData *processing.StaticProccessStructs, mutex *sync.Mutex) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Fatal(err.Error())
	}

	processors := Processors{
		Group:   makeGroupCommandProcessors(),
		Private: makePrivateCommandProcessors(),
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}
		mutex.Lock()
		processUpdate(&update, staticData, &processors)
		mutex.Unlock()
	}
}

func main() {
	apiToken, err := getApiToken()
	if err != nil {
		log.Fatal(err.Error())
	}

	config, err := loadConfig("./config.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	if config.SleepSeconds < 1 {
		log.Fatal("sleepSeconds should be set to a positive integer value")
	}

	trans, err := i18n.Tfunc(config.Language)
	if err != nil {
		log.Fatal(err.Error())
	}

	db := &database.Database{}
	err = db.Connect("./bets-data.db")
	defer db.Disconnect()

	if err != nil {
		log.Fatal("Can't connect database")
	}

	database.UpdateVersion(db)

	timers := make(map[int64]time.Time)

	mutex := &sync.Mutex{}

	chat, err := telegramChat.MakeTelegramChat(apiToken)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("Authorized on account %s", chat.GetBotUsername())

	chat.SetDebugModeEnabled(config.ExtendedLog)

	staticData := &processing.StaticProccessStructs{
		Chat:   chat,
		Db:     db,
		Config: &config,
		Timers: timers,
		Trans:  trans,
	}

	go updateTimers(staticData, mutex)
	updateBot(chat.GetBot(), staticData, mutex)
}
