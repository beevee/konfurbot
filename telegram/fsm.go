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

	talkCommand      = "🔥 Доклады"
	talkNowCommand   = "🔛 Сейчас"
	talkNextCommand  = "🔜 Скоро"
	talkAllCommand   = "📜 Все"
	talkLongCommand  = "☠ С тизерами"
	talkShortCommand = "🕊 Без тизеров"

	masterCommand = "💥 Мастер-классы"

	funCommand      = "🍾 Развлечения"
	funDayCommand   = "🍼 Утром"
	funNightCommand = "🍸 Вечером"

	transferCommand      = "🚜 Трансфер"
	transferMainCommand  = "🏎 Дежурный"
	transferColorCommand = "🚲 Цветные"
	transferNextCommand  = "🔜 Ближайшие"
	transferAllCommand   = "📜 Все рейсы"

	welcomeState       = "welcome"
	startState         = "start"
	talkState          = "talk"
	talkNowState       = "talknow"
	talkNextState      = "talknext"
	talkAllState       = "talkall"
	transferState      = "transfer"
	transferMainState  = "transfermain"
	transferColorState = "transfercolor"
	funState           = "fun"
)

var stateMessageOptions = map[string]*telebot.SendOptions{
	startState: &telebot.SendOptions{
		ReplyMarkup: telebot.ReplyMarkup{
			CustomKeyboard: [][]string{
				[]string{talkCommand, masterCommand},
				[]string{funCommand, foodCommand},
				[]string{transferCommand},
			},
			ResizeKeyboard: true,
		},
		ParseMode: telebot.ModeMarkdown,
	},

	talkState: &telebot.SendOptions{
		ReplyMarkup: telebot.ReplyMarkup{
			CustomKeyboard: [][]string{
				[]string{talkNowCommand, talkNextCommand, talkAllCommand},
			},
			ResizeKeyboard: true,
		},
		ParseMode: telebot.ModeMarkdown,
	},

	talkNowState: &telebot.SendOptions{
		ReplyMarkup: telebot.ReplyMarkup{
			CustomKeyboard: [][]string{
				[]string{talkLongCommand, talkShortCommand},
			},
			ResizeKeyboard: true,
		},
		ParseMode: telebot.ModeMarkdown,
	},

	talkNextState: &telebot.SendOptions{
		ReplyMarkup: telebot.ReplyMarkup{
			CustomKeyboard: [][]string{
				[]string{talkLongCommand, talkShortCommand},
			},
			ResizeKeyboard: true,
		},
		ParseMode: telebot.ModeMarkdown,
	},

	transferState: &telebot.SendOptions{
		ReplyMarkup: telebot.ReplyMarkup{
			CustomKeyboard: [][]string{
				[]string{transferMainCommand, transferColorCommand},
			},
			ResizeKeyboard: true,
		},
		ParseMode: telebot.ModeMarkdown,
	},

	transferMainState: &telebot.SendOptions{
		ReplyMarkup: telebot.ReplyMarkup{
			CustomKeyboard: [][]string{
				[]string{transferNextCommand, transferAllCommand},
			},
			ResizeKeyboard: true,
		},
		ParseMode: telebot.ModeMarkdown,
	},

	transferColorState: &telebot.SendOptions{
		ReplyMarkup: telebot.ReplyMarkup{
			CustomKeyboard: [][]string{
				[]string{transferNextCommand, transferAllCommand},
			},
			ResizeKeyboard: true,
		},
		ParseMode: telebot.ModeMarkdown,
	},

	funState: &telebot.SendOptions{
		ReplyMarkup: telebot.ReplyMarkup{
			CustomKeyboard: [][]string{
				[]string{funDayCommand, funNightCommand},
			},
			ResizeKeyboard: true,
		},
		ParseMode: telebot.ModeMarkdown,
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
			{Name: talkAllCommand, Src: []string{talkState}, Dst: startState},
			{Name: transferCommand, Src: []string{startState}, Dst: transferState},
			{Name: transferMainCommand, Src: []string{transferState}, Dst: transferMainState},
			{Name: transferColorCommand, Src: []string{transferState}, Dst: transferColorState},
			{Name: transferNextCommand, Src: []string{transferMainState, transferColorState}, Dst: startState},
			{Name: transferAllCommand, Src: []string{transferMainState, transferColorState}, Dst: startState},
			{Name: funCommand, Src: []string{startState}, Dst: funState},
			{Name: funDayCommand, Src: []string{funState}, Dst: startState},
			{Name: funNightCommand, Src: []string{funState}, Dst: startState},
			{Name: returnToStartCommand, Src: []string{startState}, Dst: startState},
			{
				Name: unknownCommand,
				Src:  []string{welcomeState, startState, talkState, talkNowState, talkNextState, transferState, transferMainState, transferColorState, funState},
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
				return bot.telebot.SendMessage(chat, "Окей, какие доклады?", stateMessageOptions[e.Dst])
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

			talkAllCommand: wrapCallback(func(e *fsm.Event, chat telebot.Chat, bot *Bot) error {
				events := bot.ScheduleStorage.GetEventsByType("talk")
				return bot.telebot.SendMessage(chat, makeResponseFromEvents(events, false), stateMessageOptions[e.Dst])
			}),

			transferCommand: wrapCallback(func(e *fsm.Event, chat telebot.Chat, bot *Bot) error {
				return bot.telebot.SendMessage(chat, "Окей, куда поедем?", stateMessageOptions[e.Dst])
			}),

			transferMainCommand: wrapCallback(func(e *fsm.Event, chat telebot.Chat, bot *Bot) error {
				return bot.telebot.SendMessage(chat, "Расписание довольно большое, может только ближайшие рейсы показать?", stateMessageOptions[e.Dst])
			}),

			transferColorCommand: wrapCallback(func(e *fsm.Event, chat telebot.Chat, bot *Bot) error {
				return bot.telebot.SendMessage(chat, "Расписание довольно большое, может только ближайшие рейсы показать?", stateMessageOptions[e.Dst])
			}),

			transferNextCommand: wrapCallback(func(e *fsm.Event, chat telebot.Chat, bot *Bot) error {
				var events []konfurbot.Event
				switch e.Src {
				case transferMainState:
					events = bot.ScheduleStorage.GetNextEventsByType("transfer_main", time.Now().In(bot.Timezone), time.Hour)
				case transferColorState:
					events = bot.ScheduleStorage.GetNextEventsByType("transfer_color", time.Now().In(bot.Timezone), time.Hour)
				}
				return bot.telebot.SendMessage(chat, makeResponseFromEvents(events, false), stateMessageOptions[e.Dst])
			}),

			transferAllCommand: wrapCallback(func(e *fsm.Event, chat telebot.Chat, bot *Bot) error {
				var events []konfurbot.Event
				switch e.Src {
				case transferMainState:
					events = bot.ScheduleStorage.GetEventsByType("transfer_main")
				case transferColorState:
					events = bot.ScheduleStorage.GetEventsByType("transfer_color")
				}
				return bot.telebot.SendMessage(chat, makeResponseFromEvents(events, false), stateMessageOptions[e.Dst])
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
			response += eventStart
			if eventFinish != "" {
				response += " — " + eventFinish
			}
		}

		if event.Venue != "" {
			response += fmt.Sprintf(" \\[%s]", event.Venue)
		}

		response += fmt.Sprintf(": *%s*", event.Short)

		if event.Speaker != "" {
			response += fmt.Sprintf(" (%s)", event.Speaker)
		}

		if long && event.Long != "" {
			response += fmt.Sprintf("\n%s", event.Long)
		}

		response += "\n\n"
	}

	if response == "" {
		response = "Ничего нет :("
	}

	return response
}
