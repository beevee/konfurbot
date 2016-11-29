package telegram

import (
	"fmt"
	"time"

	"github.com/beevee/konfurbot"
	"github.com/looplab/fsm"
	"github.com/tucnak/telebot"
)

const (
	greetCommand         = "greet"
	returnToStartCommand = "return"
	unknownCommand       = "unknown"

	foodCommand = "🌶 Еда"

	talkCommand       = "🔥 Доклады / МК"
	talkNowCommand    = "Которые идут сейчас"
	talkNextCommand   = "Которые начнутся скоро"
	talkAllCommand    = "Все"
	talkLongCommand   = "С тизерами (простыня!)"
	talkShortCommand  = "Без тизеров (ура! краткость!)"
	talkTalkCommand   = "Доклады"
	talkMasterCommand = "Мастер-классы"

	funCommand = "🍾 Развлечения"

	transferCommand = "🚜 Трансфер"

	welcomeState = "welcome"
	startState   = "start"
	talkState    = "talk"
	talkNowState = "talknow"
	talkAllState = "talkall"
)

var stateMessageOptions = map[string]*telebot.SendOptions{
	startState: &telebot.SendOptions{
		ReplyMarkup: telebot.ReplyMarkup{
			CustomKeyboard: [][]string{
				[]string{foodCommand, talkCommand},
				[]string{funCommand, transferCommand},
			},
			ResizeKeyboard: true,
		},
	},

	talkState: &telebot.SendOptions{
		ReplyMarkup: telebot.ReplyMarkup{
			CustomKeyboard: [][]string{
				[]string{talkNowCommand, talkNextCommand, talkAllCommand},
			},
			ResizeKeyboard: true,
		},
	},

	talkNowState: &telebot.SendOptions{
		ReplyMarkup: telebot.ReplyMarkup{
			CustomKeyboard: [][]string{
				[]string{talkLongCommand, talkShortCommand},
			},
			ResizeKeyboard: true,
		},
	},

	talkAllState: &telebot.SendOptions{
		ReplyMarkup: telebot.ReplyMarkup{
			CustomKeyboard: [][]string{
				[]string{talkTalkCommand, talkMasterCommand},
			},
			ResizeKeyboard: true,
		},
	},
}

func initStateMachine() *fsm.FSM {
	return fsm.NewFSM(
		welcomeState,

		fsm.Events{
			{Name: greetCommand, Src: []string{welcomeState}, Dst: startState},
			{Name: foodCommand, Src: []string{startState}, Dst: startState},
			{Name: talkCommand, Src: []string{startState}, Dst: talkState},
			{Name: talkNowCommand, Src: []string{talkState}, Dst: talkNowState},
			{Name: talkLongCommand, Src: []string{talkNowState}, Dst: startState},
			{Name: talkShortCommand, Src: []string{talkNowState}, Dst: startState},
			{Name: talkAllCommand, Src: []string{talkState}, Dst: talkAllState},
			{Name: talkTalkCommand, Src: []string{talkAllState}, Dst: startState},
			{Name: talkMasterCommand, Src: []string{talkAllState}, Dst: startState},
			{Name: returnToStartCommand, Src: []string{startState}, Dst: startState},
			{Name: unknownCommand, Src: []string{startState}, Dst: startState},
		},

		fsm.Callbacks{
			greetCommand: wrapCallback(func(e *fsm.Event, chat telebot.Chat, bot *Bot) error {
				return bot.telebot.SendMessage(chat, "Добро пожаловать на КонфУР!", stateMessageOptions[e.Dst])
			}),

			foodCommand: wrapCallback(func(e *fsm.Event, chat telebot.Chat, bot *Bot) error {
				events := bot.ScheduleStorage.GetEventsByType("food")
				return bot.telebot.SendMessage(chat, makeResponseFromEvents(events, false), stateMessageOptions[e.Dst])
			}),

			talkCommand: wrapCallback(func(e *fsm.Event, chat telebot.Chat, bot *Bot) error {
				return bot.telebot.SendMessage(chat, "Окей, какие доклады и мастер-классы?", stateMessageOptions[e.Dst])
			}),

			talkNowCommand: wrapCallback(func(e *fsm.Event, chat telebot.Chat, bot *Bot) error {
				return bot.telebot.SendMessage(chat, "Их может оказаться довольно много. Тизеры надо?", stateMessageOptions[e.Dst])
			}),

			talkLongCommand: wrapCallback(func(e *fsm.Event, chat telebot.Chat, bot *Bot) error {
				var events []konfurbot.Event
				switch e.Src {
				case talkNowState:
					events = bot.ScheduleStorage.GetCurrentEventsByType("talk", time.Now().In(bot.Timezone))
				}
				return bot.telebot.SendMessage(chat, makeResponseFromEvents(events, true), stateMessageOptions[e.Dst])
			}),

			talkShortCommand: wrapCallback(func(e *fsm.Event, chat telebot.Chat, bot *Bot) error {
				var events []konfurbot.Event
				switch e.Src {
				case talkNowState:
					events = bot.ScheduleStorage.GetCurrentEventsByType("talk", time.Now().In(bot.Timezone))
				}
				return bot.telebot.SendMessage(chat, makeResponseFromEvents(events, false), stateMessageOptions[e.Dst])
			}),

			talkTalkCommand: wrapCallback(func(e *fsm.Event, chat telebot.Chat, bot *Bot) error {
				events := bot.ScheduleStorage.GetEventsByTypeAndSubtype("talk", "talk")
				return bot.telebot.SendMessage(chat, makeResponseFromEvents(events, false), stateMessageOptions[e.Dst])
			}),

			talkMasterCommand: wrapCallback(func(e *fsm.Event, chat telebot.Chat, bot *Bot) error {
				events := bot.ScheduleStorage.GetEventsByTypeAndSubtype("talk", "master")
				return bot.telebot.SendMessage(chat, makeResponseFromEvents(events, false), stateMessageOptions[e.Dst])
			}),

			talkAllCommand: wrapCallback(func(e *fsm.Event, chat telebot.Chat, bot *Bot) error {
				return bot.telebot.SendMessage(chat, "Полное расписание довольно длинное. Давай посмотрим отдельно, доклады или мастер-классы? С тизерами вообще не буду предлагать :)", stateMessageOptions[e.Dst])
			}),

			unknownCommand: wrapCallback(func(e *fsm.Event, chat telebot.Chat, bot *Bot) error {
				return bot.telebot.SendMessage(chat, "Я не понимаю эту команду. Давай попробуем еще раз с начала.", stateMessageOptions[e.Dst])
			}),
		},
	)
}

func wrapCallback(f func(*fsm.Event, telebot.Chat, *Bot) error) func(*fsm.Event) {
	return func(e *fsm.Event) {
		if len(e.Args) < 2 {
			return
		}
		chat := e.Args[0].(telebot.Chat)
		bot := e.Args[1].(*Bot)
		if err := f(e, chat, bot); err != nil {
			bot.Logger.Log("msg", "error sending message", "chatid", chat.ID, "error", err)
		}
	}
}

func makeResponseFromEvents(events []konfurbot.Event, long bool) string {
	var response string
	for _, event := range events {
		response += fmt.Sprintf("%s — %s: %s\n",
			event.Start.Format("15:04"), event.Finish.Format("15:04"), event.Short)
		if long {
			response += fmt.Sprintf("%s\n\n", event.Long)
		}
	}
	if response == "" {
		response = "Ничего нет :("
	}
	return response
}
