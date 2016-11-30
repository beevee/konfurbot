package telegram

import (
	"fmt"
	"time"

	"github.com/beevee/konfurbot"
	"github.com/looplab/fsm"
	"github.com/tucnak/telebot"
)

const (
	welcomeCommand       = "greet"
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

	funCommand      = "🍾 Развлечения"
	funDayCommand   = "🍼 Утром"
	funNightCommand = "🍸 Вечером"

	transferCommand = "🚜 Трансфер"

	welcomeState  = "welcome"
	startState    = "start"
	talkState     = "talk"
	talkNowState  = "talknow"
	talkNextState = "talknext"
	talkAllState  = "talkall"
	funState      = "fun"
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
				[]string{talkNowCommand},
				[]string{talkNextCommand},
				[]string{talkAllCommand},
			},
			ResizeKeyboard: true,
		},
	},

	talkNowState: &telebot.SendOptions{
		ReplyMarkup: telebot.ReplyMarkup{
			CustomKeyboard: [][]string{
				[]string{talkLongCommand},
				[]string{talkShortCommand},
			},
			ResizeKeyboard: true,
		},
	},

	talkNextState: &telebot.SendOptions{
		ReplyMarkup: telebot.ReplyMarkup{
			CustomKeyboard: [][]string{
				[]string{talkLongCommand},
				[]string{talkShortCommand},
			},
			ResizeKeyboard: true,
		},
	},

	talkAllState: &telebot.SendOptions{
		ReplyMarkup: telebot.ReplyMarkup{
			CustomKeyboard: [][]string{
				[]string{talkTalkCommand},
				[]string{talkMasterCommand},
			},
			ResizeKeyboard: true,
		},
	},

	funState: &telebot.SendOptions{
		ReplyMarkup: telebot.ReplyMarkup{
			CustomKeyboard: [][]string{
				[]string{funDayCommand, funNightCommand},
			},
			ResizeKeyboard: true,
		},
	},
}

func initStateMachine() *fsm.FSM {
	return fsm.NewFSM(
		welcomeState,

		fsm.Events{
			{Name: welcomeCommand, Src: []string{welcomeState}, Dst: startState},
			{Name: foodCommand, Src: []string{startState}, Dst: startState},
			{Name: talkCommand, Src: []string{startState}, Dst: talkState},
			{Name: talkNowCommand, Src: []string{talkState}, Dst: talkNowState},
			{Name: talkNextCommand, Src: []string{talkState}, Dst: talkNextState},
			{Name: talkLongCommand, Src: []string{talkNowState, talkNextState}, Dst: startState},
			{Name: talkShortCommand, Src: []string{talkNowState, talkNextState}, Dst: startState},
			{Name: talkAllCommand, Src: []string{talkState}, Dst: talkAllState},
			{Name: talkTalkCommand, Src: []string{talkAllState}, Dst: startState},
			{Name: talkMasterCommand, Src: []string{talkAllState}, Dst: startState},
			{Name: funCommand, Src: []string{startState}, Dst: funState},
			{Name: funDayCommand, Src: []string{funState}, Dst: startState},
			{Name: funNightCommand, Src: []string{funState}, Dst: startState},
			{Name: returnToStartCommand, Src: []string{startState}, Dst: startState},
			{
				Name: unknownCommand,
				Src:  []string{welcomeState, startState, talkState, talkNowState, talkNextState, talkAllState, funState},
				Dst:  startState,
			},
		},

		fsm.Callbacks{
			welcomeCommand: wrapCallback(func(e *fsm.Event, chat telebot.Chat, bot *Bot) error {
				return bot.telebot.SendMessage(chat, "Добро пожаловать на КонфУР!", stateMessageOptions[e.Dst])
			}),

			foodCommand: wrapCallback(func(e *fsm.Event, chat telebot.Chat, bot *Bot) error {
				events := bot.ScheduleStorage.GetEventsByType("food")
				return bot.telebot.SendMessage(chat, makeResponseFromEvents(events, true), stateMessageOptions[e.Dst])
			}),

			talkCommand: wrapCallback(func(e *fsm.Event, chat telebot.Chat, bot *Bot) error {
				return bot.telebot.SendMessage(chat, "Окей, какие доклады и мастер-классы?", stateMessageOptions[e.Dst])
			}),

			talkNowCommand: wrapCallback(func(e *fsm.Event, chat telebot.Chat, bot *Bot) error {
				return bot.telebot.SendMessage(chat, "Их может оказаться довольно много. Тизеры надо?", stateMessageOptions[e.Dst])
			}),

			talkNextCommand: wrapCallback(func(e *fsm.Event, chat telebot.Chat, bot *Bot) error {
				return bot.telebot.SendMessage(chat, "Их может оказаться довольно много. Тизеры надо?", stateMessageOptions[e.Dst])
			}),

			talkLongCommand: wrapCallback(func(e *fsm.Event, chat telebot.Chat, bot *Bot) error {
				var events []konfurbot.Event
				switch e.Src {
				case talkNowState:
					events = bot.ScheduleStorage.GetCurrentEventsByType("talk", time.Now().In(bot.Timezone))
				case talkNextState:
					events = bot.ScheduleStorage.GetNextEventsByType("talk", time.Now().In(bot.Timezone), time.Hour)
				}
				return bot.telebot.SendMessage(chat, makeResponseFromEvents(events, true), stateMessageOptions[e.Dst])
			}),

			talkShortCommand: wrapCallback(func(e *fsm.Event, chat telebot.Chat, bot *Bot) error {
				var events []konfurbot.Event
				switch e.Src {
				case talkNowState:
					events = bot.ScheduleStorage.GetCurrentEventsByType("talk", time.Now().In(bot.Timezone))
				case talkNextState:
					events = bot.ScheduleStorage.GetNextEventsByType("talk", time.Now().In(bot.Timezone), time.Hour)
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

			funCommand: wrapCallback(func(e *fsm.Event, chat telebot.Chat, bot *Bot) error {
				return bot.telebot.SendMessage(chat, "Утром или вечером?", stateMessageOptions[e.Dst])
			}),

			funDayCommand: wrapCallback(func(e *fsm.Event, chat telebot.Chat, bot *Bot) error {
				events := bot.ScheduleStorage.GetDayEventsByType("fun")
				return bot.telebot.SendMessage(chat, makeResponseFromEvents(events, false), stateMessageOptions[e.Dst])
			}),

			funNightCommand: wrapCallback(func(e *fsm.Event, chat telebot.Chat, bot *Bot) error {
				events := bot.ScheduleStorage.GetNightEventsByType("fun")
				return bot.telebot.SendMessage(chat, makeResponseFromEvents(events, false), stateMessageOptions[e.Dst])
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
		var eventStart, eventFinish string
		if event.Start != nil {
			eventStart = event.Start.Format("15:04")
		}
		if event.Finish != nil {
			eventFinish = event.Finish.Format("15:04")
		}
		if eventStart == "" && eventFinish == "" {
			response += "весь день"
		} else {
			response += fmt.Sprintf("%s — %s", eventStart, eventFinish)
		}

		if event.Venue != "" {
			response += fmt.Sprintf(" [%s]", event.Venue)
		}

		response += fmt.Sprintf(": %s", event.Short)

		if event.Speaker != "" {
			response += fmt.Sprintf(" (%s)", event.Speaker)
		}

		response += "\n"

		if long && event.Long != "" {
			response += fmt.Sprintf("%s\n\n", event.Long)
		}
	}

	if response == "" {
		response = "Ничего нет :("
	}

	return response
}
