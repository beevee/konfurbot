package telegram

import (
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/looplab/fsm"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tucnak/telebot"
	"gopkg.in/tomb.v2"

	"github.com/beevee/konfurbot"
	"github.com/beevee/konfurbot/mock"
)

func TestTelegram(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	Convey("A user is travelling through our menus", t, func() {
		mockLogger := mock.NewMockLogger(ctrl)
		mockStorage := mock.NewMockScheduleStorage(ctrl)
		mockTelebot := mock.NewMockTelebotInterface(ctrl)

		tz, _ := time.LoadLocation("Asia/Yekaterinburg")
		chat := telebot.Chat{ID: 1}
		start, _ := time.Parse("15:04", "17:00")
		finish, _ := time.Parse("15:04", "19:00")

		bot := &Bot{
			ScheduleStorage:   mockStorage,
			Timezone:          tz,
			Logger:            mockLogger,
			telebot:           mockTelebot,
			chatStateMachines: make(map[int64]*fsm.FSM, 0),
			chatStateLock:     sync.RWMutex{},
			tomb:              tomb.Tomb{},
		}

		Convey("and at first she sees the welcome message", func() {
			mockLogger.EXPECT().Log("msg", "message received", "firstname", "",
				"lastname", "", "username", "", "chatid", chat.ID, "command", "/start")
			mockLogger.EXPECT().Log("msg", "no state machine exists, creating new", "chatid", chat.ID)

			mockTelebot.EXPECT().SendMessage(chat, "Добро пожаловать на КонфУР!", gomock.Any())

			bot.handleMessage(telebot.Message{
				Chat: chat,
				Text: "/start",
			})

			Convey("then she decides to see food-related events", func() {
				mockLogger.EXPECT().Log("msg", "message received", "firstname", "",
					"lastname", "", "username", "", "chatid", chat.ID, "command", "🌶 Еда")
				mockLogger.EXPECT().Log("msg", "state machine exists, attempting transition",
					"chatid", chat.ID, "currentstate", startState, "command", "🌶 Еда")

				mockStorage.EXPECT().GetEventsByType("food").Return([]konfurbot.Event{
					konfurbot.Event{Type: "food", Short: "お好み焼き", Start: start, Finish: finish},
					konfurbot.Event{Type: "food", Short: "焼き鳥", Start: start, Finish: finish},
				})

				mockTelebot.EXPECT().SendMessage(chat, "17:00 — 19:00: お好み焼き\n17:00 — 19:00: 焼き鳥\n", gomock.Any())

				bot.handleMessage(telebot.Message{
					Chat: chat,
					Text: "🌶 Еда",
				})
			})

			Convey("then she enters some gibberish", func() {
				mockLogger.EXPECT().Log("msg", "message received", "firstname", "",
					"lastname", "", "username", "", "chatid", chat.ID, "command", "gibberish")
				mockLogger.EXPECT().Log("msg", "state machine exists, attempting transition",
					"chatid", chat.ID, "currentstate", startState, "command", "gibberish")
				mockLogger.EXPECT().Log("msg", "something is wrong with the transition, will return to start",
					"currentstate", startState, "chatid", chat.ID, "error", gomock.Any())

				mockTelebot.EXPECT().SendMessage(chat, "Я не понимаю эту команду. Давай попробуем еще раз с начала.", gomock.Any())

				bot.handleMessage(telebot.Message{
					Chat: chat,
					Text: "gibberish",
				})
			})
		})
	})
}
