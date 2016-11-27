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
			greetCommand: wrapCallback(func(e *fsm.Event, chat telebot.Chat, bot *Bot) error {
				return bot.telebot.SendMessage(chat, "Добро пожаловать на КонфУР!", stateMessageOptions[e.Dst])
			}),

			foodCommand: wrapCallback(func(e *fsm.Event, chat telebot.Chat, bot *Bot) error {
				var response string
				for _, event := range bot.ScheduleStorage.GetEventsByType("food") {
					response += fmt.Sprintf("%s — %s: %s\n",
						event.Start.Format("15:04"), event.Finish.Format("15:04"), event.Short)
				}
				return bot.telebot.SendMessage(chat, response, stateMessageOptions[e.Dst])
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
