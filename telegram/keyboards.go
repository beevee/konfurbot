package telegram

import "github.com/tucnak/telebot"

var stateKeyboards = map[string][][]string{
	startState: [][]string{
		[]string{talkCommand, masterCommand},
		[]string{funCommand, foodCommand},
		[]string{transferCommand},
	},

	talkState: [][]string{
		[]string{talkNowCommand, talkNextCommand, talkAllCommand},
	},

	talkNowState: [][]string{
		[]string{talkLongCommand, talkShortCommand},
	},

	talkNextState: [][]string{
		[]string{talkLongCommand, talkShortCommand},
	},

	masterState: [][]string{
		[]string{masterNowCommand, masterNextCommand, masterAllCommand},
	},

	masterNowState: [][]string{
		[]string{masterLongCommand, masterShortCommand},
	},

	masterNextState: [][]string{
		[]string{masterLongCommand, masterShortCommand},
	},

	transferState: [][]string{
		[]string{transferMainCommand, transferColorCommand},
	},

	transferMainState: [][]string{
		[]string{transferNextCommand, transferAllCommand},
	},

	transferColorState: [][]string{
		[]string{transferNextCommand, transferAllCommand},
	},

	funState: [][]string{
		[]string{funDayCommand, funNightCommand},
	},
}

func makeMessageOptionsForState(state string) *telebot.SendOptions {
	return &telebot.SendOptions{
		ReplyMarkup: telebot.ReplyMarkup{
			CustomKeyboard: stateKeyboards[state],
			ResizeKeyboard: true,
		},
		ParseMode: telebot.ModeMarkdown,
	}
}
