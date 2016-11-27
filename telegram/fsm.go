package telegram

import (
	"fmt"

	"github.com/looplab/fsm"
	"github.com/tucnak/telebot"
)

const (
	greetCommand         = "greet"
	returnToStartCommand = "return"
	unknownCommand       = "unknown"

	foodCommand     = "🌶 Еда"
	talkCommand     = "🔥 Доклады / МК"
	funCommand      = "🍾 Развлечения"
	transferCommand = "🚜 Трансфер"

	welcomeState = "welcome"
	startState   = "start"
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
}

func initStateMachine() *fsm.FSM {
	return fsm.NewFSM(
		welcomeState,

		fsm.Events{
			{Name: greetCommand, Src: []string{welcomeState}, Dst: startState},
			{Name: returnToStartCommand, Src: []string{startState}, Dst: startState},
			{Name: unknownCommand, Src: []string{startState}, Dst: startState},
			{Name: foodCommand, Src: []string{startState}, Dst: startState},
		},

		fsm.Callbacks{
			greetCommand: extractCallbackParams(func(e *fsm.Event, chat telebot.Chat, bot *Bot) {
				bot.telebot.SendMessage(chat, "Добро пожаловать на КонфУР!", stateMessageOptions[e.Dst])
			}),

			foodCommand: extractCallbackParams(func(e *fsm.Event, chat telebot.Chat, bot *Bot) {
				var response string
				for _, event := range bot.ScheduleStorage.GetEventsByType("food") {
					response += fmt.Sprintf("%s — %s: %s\n",
						event.Start.Format("15:04"), event.Finish.Format("15:04"), event.Short)
				}
				bot.telebot.SendMessage(chat, response, stateMessageOptions[e.Dst])
			}),

			unknownCommand: extractCallbackParams(func(e *fsm.Event, chat telebot.Chat, bot *Bot) {
				bot.telebot.SendMessage(chat, "Я не понимаю эту команду. Давай попробуем еще раз с начала.", stateMessageOptions[e.Dst])
			}),
		},
	)
}

func extractCallbackParams(f func(*fsm.Event, telebot.Chat, *Bot)) func(*fsm.Event) {
	return func(e *fsm.Event) {
		if len(e.Args) < 2 {
			return
		}
		f(e, e.Args[0].(telebot.Chat), e.Args[0].(*Bot))
	}
}
