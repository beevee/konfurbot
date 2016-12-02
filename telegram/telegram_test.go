package telegram

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/golang/mock/gomock"
	"github.com/looplab/fsm"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tucnak/telebot"

	"github.com/beevee/konfurbot"
	"github.com/beevee/konfurbot/mock"
)

func TestTelegramIntegration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	Convey("Попробуем запустить и остановить бота (этот тест не будет работать без настоящего токена в переменной окружения KONFURBOT_TEST_TOKEN)", t, func() {
		bot := &Bot{
			ScheduleStorage: mock.NewMockScheduleStorage(ctrl),
			TelegramToken:   os.Getenv("KONFURBOT_TEST_TOKEN"),
			Logger:          log.NewNopLogger(),
		}

		err := bot.Start()
		So(err, ShouldBeNil)

		err = bot.Stop()
		So(err, ShouldBeNil)
	})

	Convey("Попробуем запустить бота с неправильным токеном", t, func() {
		bot := &Bot{
			ScheduleStorage: mock.NewMockScheduleStorage(ctrl),
			TelegramToken:   "123456",
			Logger:          log.NewNopLogger(),
		}

		err := bot.Start()
		So(err, ShouldNotBeNil)
	})
}

func TestUserInteraction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	Convey("Пользователь зашел в наш чатик", t, func() {
		mockStorage := mock.NewMockScheduleStorage(ctrl)
		mockTelebot := mock.NewMockTelebotInterface(ctrl)

		tz, _ := time.LoadLocation("Asia/Yekaterinburg")
		chat := telebot.Chat{ID: 1}
		start, _ := time.Parse("15:04", "17:00")
		finish, _ := time.Parse("15:04", "19:00")

		bot := &Bot{
			ScheduleStorage:   mockStorage,
			Timezone:          tz,
			Logger:            log.NewNopLogger(),
			telebot:           mockTelebot,
			chatStateMachines: make(map[int64]*fsm.FSM, 0),
		}

		Convey("сначала показываем приветствие", func() {
			mockTelebot.EXPECT().SendMessage(chat, "Добро пожаловать на КонфУР!",
				hasButtons("🌶 Еда", "🔥 Доклады", "💥 Мастер-классы", "🍾 Развлечения", "🚜 Трансфер"))
			bot.handleMessage(telebot.Message{Chat: chat, Text: "/start"})

			Convey("пользователь спрашивает про еду (у еды есть место проведения, но нет спикера)", func() {
				mockStorage.EXPECT().GetEventsByType("food").Return([]konfurbot.Event{
					konfurbot.Event{Type: "food", Short: "お好み焼き", Venue: "Бар", Start: &start, Finish: &finish},
					konfurbot.Event{Type: "food", Short: "焼き鳥", Long: "Вегетарианцам накроют на крыше паркинга", Venue: "Кафе", Start: &start, Finish: &finish},
				})
				mockTelebot.EXPECT().SendMessage(chat, "17:00 — 19:00 \\[Бар]: *お好み焼き*\n\n17:00 — 19:00 \\[Кафе]: *焼き鳥*\nВегетарианцам накроют на крыше паркинга\n\n",
					hasButtons("🌶 Еда", "🔥 Доклады", "💥 Мастер-классы", "🍾 Развлечения", "🚜 Трансфер"))
				bot.handleMessage(telebot.Message{Chat: chat, Text: "🌶 Еда"})
			})

			Convey("пользователь спрашивает про доклады", func() {
				mockTelebot.EXPECT().SendMessage(chat, "Окей, какие доклады?",
					hasButtons("🔛 Сейчас", "🔜 Скоро", "📜 Все"))
				bot.handleMessage(telebot.Message{Chat: chat, Text: "🔥 Доклады"})

				Convey("которые идут сейчас", func() {
					mockTelebot.EXPECT().SendMessage(chat, "Их может оказаться довольно много. Тизеры надо?",
						hasButtons("☠ С тизерами", "🕊 Без тизеров"))
					bot.handleMessage(telebot.Message{Chat: chat, Text: "🔛 Сейчас"})

					Convey("с тизерами, и что-то сейчас идет", func() {
						mockStorage.EXPECT().GetCurrentEventsByType("talk", gomock.Any()).Return([]konfurbot.Event{
							konfurbot.Event{Type: "talk", Short: "WAT", Long: "WAAAAT", Start: &start, Finish: &finish},
							konfurbot.Event{Type: "talk", Short: "WAT 2", Long: "WAAAAT 22", Start: &start, Finish: &finish},
						})
						mockTelebot.EXPECT().SendMessage(chat, "17:00 — 19:00: *WAT*\nWAAAAT\n\n17:00 — 19:00: *WAT 2*\nWAAAAT 22\n\n",
							hasButtons("🌶 Еда", "🔥 Доклады", "💥 Мастер-классы", "🍾 Развлечения", "🚜 Трансфер"))
						bot.handleMessage(telebot.Message{Chat: chat, Text: "☠ С тизерами"})
					})

					Convey("с тизерами, и сейчас ничего не идет", func() {
						mockStorage.EXPECT().GetCurrentEventsByType("talk", gomock.Any()).Return([]konfurbot.Event{})
						mockTelebot.EXPECT().SendMessage(chat, "Ничего нет :(",
							hasButtons("🌶 Еда", "🔥 Доклады", "💥 Мастер-классы", "🍾 Развлечения", "🚜 Трансфер"))
						bot.handleMessage(telebot.Message{Chat: chat, Text: "🕊 Без тизеров"})
					})

					Convey("без тизеров, и сейчас что-то идет", func() {
						mockStorage.EXPECT().GetCurrentEventsByType("talk", gomock.Any()).Return([]konfurbot.Event{
							konfurbot.Event{Type: "talk", Short: "WAT", Long: "WAAAAT", Start: &start, Finish: &finish},
							konfurbot.Event{Type: "talk", Short: "WAT 2", Long: "WAAAAT 22", Start: &start, Finish: &finish},
						})
						mockTelebot.EXPECT().SendMessage(chat, "17:00 — 19:00: *WAT*\n\n17:00 — 19:00: *WAT 2*\n\n",
							hasButtons("🌶 Еда", "🔥 Доклады", "💥 Мастер-классы", "🍾 Развлечения", "🚜 Трансфер"))
						bot.handleMessage(telebot.Message{Chat: chat, Text: "🕊 Без тизеров"})
					})

					Convey("пользователь пишет нам ерунду", func() {
						mockTelebot.EXPECT().SendMessage(chat, "Я не понимаю эту команду. Давай попробуем еще раз с начала.",
							hasButtons("🌶 Еда", "🔥 Доклады", "💥 Мастер-классы", "🍾 Развлечения", "🚜 Трансфер"))
						bot.handleMessage(telebot.Message{Chat: chat, Text: "gibberish"})
					})
				})

				Convey("которые начнутся в ближайший час", func() {
					mockTelebot.EXPECT().SendMessage(chat, "Их может оказаться довольно много. Тизеры надо?",
						hasButtons("☠ С тизерами", "🕊 Без тизеров"))
					bot.handleMessage(telebot.Message{Chat: chat, Text: "🔜 Скоро"})

					Convey("с тизерами", func() {
						mockStorage.EXPECT().GetNextEventsByType("talk", gomock.Any(), 2*time.Hour).Return([]konfurbot.Event{
							konfurbot.Event{Type: "talk", Short: "WAT", Long: "WAAAAT", Start: &start, Finish: &finish},
							konfurbot.Event{Type: "talk", Short: "WAT 2", Long: "WAAAAT 22", Start: &start, Finish: &finish},
						})
						mockTelebot.EXPECT().SendMessage(chat, "17:00 — 19:00: *WAT*\nWAAAAT\n\n17:00 — 19:00: *WAT 2*\nWAAAAT 22\n\n",
							hasButtons("🌶 Еда", "🔥 Доклады", "💥 Мастер-классы", "🍾 Развлечения", "🚜 Трансфер"))
						bot.handleMessage(telebot.Message{Chat: chat, Text: "☠ С тизерами"})
					})

					Convey("без тизеров", func() {
						mockStorage.EXPECT().GetNextEventsByType("talk", gomock.Any(), 2*time.Hour).Return([]konfurbot.Event{
							konfurbot.Event{Type: "talk", Short: "WAT", Long: "WAAAAT", Start: &start, Finish: &finish},
							konfurbot.Event{Type: "talk", Short: "WAT 2", Long: "WAAAAT 22", Start: &start, Finish: &finish},
						})
						mockTelebot.EXPECT().SendMessage(chat, "17:00 — 19:00: *WAT*\n\n17:00 — 19:00: *WAT 2*\n\n",
							hasButtons("🌶 Еда", "🔥 Доклады", "💥 Мастер-классы", "🍾 Развлечения", "🚜 Трансфер"))
						bot.handleMessage(telebot.Message{Chat: chat, Text: "🕊 Без тизеров"})
					})

					Convey("пользователь пишет нам ерунду", func() {
						mockTelebot.EXPECT().SendMessage(chat, "Я не понимаю эту команду. Давай попробуем еще раз с начала.",
							hasButtons("🌶 Еда", "🔥 Доклады", "💥 Мастер-классы", "🍾 Развлечения", "🚜 Трансфер"))
						bot.handleMessage(telebot.Message{Chat: chat, Text: "gibberish"})
					})
				})

				Convey("все (у них есть и спикер, и место проведения)", func() {
					mockStorage.EXPECT().GetEventsByType("talk").Return([]konfurbot.Event{
						konfurbot.Event{
							Type:    "talk",
							Speaker: "Александр Казаков",
							Venue:   "Учебный класс 1",
							Short:   "WAT",
							Long:    "WAAAAT",
							Start:   &start,
							Finish:  &finish,
						},
						konfurbot.Event{
							Type:    "talk",
							Speaker: "Василий Петров",
							Venue:   "Учебный класс 2",
							Short:   "WAT 2",
							Long:    "WAAAAT 22",
							Start:   &start,
							Finish:  &finish,
						},
					})
					mockTelebot.EXPECT().SendMessage(chat, "17:00 — 19:00 \\[Учебный класс 1]: *WAT* (Александр Казаков)\n\n17:00 — 19:00 \\[Учебный класс 2]: *WAT 2* (Василий Петров)\n\n",
						hasButtons("🌶 Еда", "🔥 Доклады", "💥 Мастер-классы", "🍾 Развлечения", "🚜 Трансфер"))
					bot.handleMessage(telebot.Message{Chat: chat, Text: "📜 Все"})
				})

				Convey("пользователь пишет нам ерунду", func() {
					mockTelebot.EXPECT().SendMessage(chat, "Я не понимаю эту команду. Давай попробуем еще раз с начала.",
						hasButtons("🌶 Еда", "🔥 Доклады", "💥 Мастер-классы", "🍾 Развлечения", "🚜 Трансфер"))
					bot.handleMessage(telebot.Message{Chat: chat, Text: "gibberish"})
				})
			})

			Convey("пользователь спрашивает про мастер-классы", func() {
				mockTelebot.EXPECT().SendMessage(chat, "Окей, какие мастер-классы?",
					hasButtons("▶️ Сейчас", "⏭ Скоро", "🔢 Все"))
				bot.handleMessage(telebot.Message{Chat: chat, Text: "💥 Мастер-классы"})

				Convey("которые идут сейчас", func() {
					mockTelebot.EXPECT().SendMessage(chat, "Их может оказаться довольно много. Тизеры надо?",
						hasButtons("🌪 С тизерами", "🌴 Без тизеров"))
					bot.handleMessage(telebot.Message{Chat: chat, Text: "▶️ Сейчас"})

					Convey("с тизерами, и что-то сейчас идет", func() {
						mockStorage.EXPECT().GetCurrentEventsByType("master", gomock.Any()).Return([]konfurbot.Event{
							konfurbot.Event{Type: "master", Short: "WAT", Long: "WAAAAT", Start: &start, Finish: &finish},
							konfurbot.Event{Type: "master", Short: "WAT 2", Long: "WAAAAT 22", Start: &start, Finish: &finish},
						})
						mockTelebot.EXPECT().SendMessage(chat, "17:00 — 19:00: *WAT*\nWAAAAT\n\n17:00 — 19:00: *WAT 2*\nWAAAAT 22\n\n",
							hasButtons("🌶 Еда", "🔥 Доклады", "💥 Мастер-классы", "🍾 Развлечения", "🚜 Трансфер"))
						bot.handleMessage(telebot.Message{Chat: chat, Text: "🌪 С тизерами"})
					})

					Convey("с тизерами, и сейчас ничего не идет", func() {
						mockStorage.EXPECT().GetCurrentEventsByType("master", gomock.Any()).Return([]konfurbot.Event{})
						mockTelebot.EXPECT().SendMessage(chat, "Ничего нет :(",
							hasButtons("🌶 Еда", "🔥 Доклады", "💥 Мастер-классы", "🍾 Развлечения", "🚜 Трансфер"))
						bot.handleMessage(telebot.Message{Chat: chat, Text: "🌴 Без тизеров"})
					})

					Convey("без тизеров, и сейчас что-то идет", func() {
						mockStorage.EXPECT().GetCurrentEventsByType("master", gomock.Any()).Return([]konfurbot.Event{
							konfurbot.Event{Type: "master", Short: "WAT", Long: "WAAAAT", Start: &start, Finish: &finish},
							konfurbot.Event{Type: "master", Short: "WAT 2", Long: "WAAAAT 22", Start: &start, Finish: &finish},
						})
						mockTelebot.EXPECT().SendMessage(chat, "17:00 — 19:00: *WAT*\n\n17:00 — 19:00: *WAT 2*\n\n",
							hasButtons("🌶 Еда", "🔥 Доклады", "💥 Мастер-классы", "🍾 Развлечения", "🚜 Трансфер"))
						bot.handleMessage(telebot.Message{Chat: chat, Text: "🌴 Без тизеров"})
					})

					Convey("пользователь пишет нам ерунду", func() {
						mockTelebot.EXPECT().SendMessage(chat, "Я не понимаю эту команду. Давай попробуем еще раз с начала.",
							hasButtons("🌶 Еда", "🔥 Доклады", "💥 Мастер-классы", "🍾 Развлечения", "🚜 Трансфер"))
						bot.handleMessage(telebot.Message{Chat: chat, Text: "gibberish"})
					})
				})

				Convey("которые начнутся в ближайший час", func() {
					mockTelebot.EXPECT().SendMessage(chat, "Их может оказаться довольно много. Тизеры надо?",
						hasButtons("🌪 С тизерами", "🌴 Без тизеров"))
					bot.handleMessage(telebot.Message{Chat: chat, Text: "⏭ Скоро"})

					Convey("с тизерами", func() {
						mockStorage.EXPECT().GetNextEventsByType("master", gomock.Any(), 2*time.Hour).Return([]konfurbot.Event{
							konfurbot.Event{Type: "master", Short: "WAT", Long: "WAAAAT", Start: &start, Finish: &finish},
							konfurbot.Event{Type: "master", Short: "WAT 2", Long: "WAAAAT 22", Start: &start, Finish: &finish},
						})
						mockTelebot.EXPECT().SendMessage(chat, "17:00 — 19:00: *WAT*\nWAAAAT\n\n17:00 — 19:00: *WAT 2*\nWAAAAT 22\n\n",
							hasButtons("🌶 Еда", "🔥 Доклады", "💥 Мастер-классы", "🍾 Развлечения", "🚜 Трансфер"))
						bot.handleMessage(telebot.Message{Chat: chat, Text: "🌪 С тизерами"})
					})

					Convey("без тизеров", func() {
						mockStorage.EXPECT().GetNextEventsByType("master", gomock.Any(), 2*time.Hour).Return([]konfurbot.Event{
							konfurbot.Event{Type: "master", Short: "WAT", Long: "WAAAAT", Start: &start, Finish: &finish},
							konfurbot.Event{Type: "master", Short: "WAT 2", Long: "WAAAAT 22", Start: &start, Finish: &finish},
						})
						mockTelebot.EXPECT().SendMessage(chat, "17:00 — 19:00: *WAT*\n\n17:00 — 19:00: *WAT 2*\n\n",
							hasButtons("🌶 Еда", "🔥 Доклады", "💥 Мастер-классы", "🍾 Развлечения", "🚜 Трансфер"))
						bot.handleMessage(telebot.Message{Chat: chat, Text: "🌴 Без тизеров"})
					})

					Convey("пользователь пишет нам ерунду", func() {
						mockTelebot.EXPECT().SendMessage(chat, "Я не понимаю эту команду. Давай попробуем еще раз с начала.",
							hasButtons("🌶 Еда", "🔥 Доклады", "💥 Мастер-классы", "🍾 Развлечения", "🚜 Трансфер"))
						bot.handleMessage(telebot.Message{Chat: chat, Text: "gibberish"})
					})
				})

				Convey("все (у них есть спикер, но нет места проведения)", func() {
					mockStorage.EXPECT().GetEventsByType("master").Return([]konfurbot.Event{
						konfurbot.Event{
							Type:    "talk",
							Speaker: "Александр Казаков",
							Short:   "WAT",
							Long:    "WAAAAT",
							Start:   &start,
							Finish:  &finish,
						},
						konfurbot.Event{
							Type:    "talk",
							Speaker: "Василий Петров",
							Short:   "WAT 2",
							Long:    "WAAAAT 22",
							Start:   &start,
							Finish:  &finish,
						},
					})
					mockTelebot.EXPECT().SendMessage(chat, "17:00 — 19:00: *WAT* (Александр Казаков)\n\n17:00 — 19:00: *WAT 2* (Василий Петров)\n\n",
						hasButtons("🌶 Еда", "🔥 Доклады", "💥 Мастер-классы", "🍾 Развлечения", "🚜 Трансфер"))
					bot.handleMessage(telebot.Message{Chat: chat, Text: "🔢 Все"})
				})

				Convey("пользователь пишет нам ерунду", func() {
					mockTelebot.EXPECT().SendMessage(chat, "Я не понимаю эту команду. Давай попробуем еще раз с начала.",
						hasButtons("🌶 Еда", "🔥 Доклады", "💥 Мастер-классы", "🍾 Развлечения", "🚜 Трансфер"))
					bot.handleMessage(telebot.Message{Chat: chat, Text: "gibberish"})
				})
			})

			Convey("пользователь спрашивает про развлечения", func() {
				mockTelebot.EXPECT().SendMessage(chat, "Утром или вечером?",
					hasButtons("🍼 Утром", "🍸 Вечером"))
				bot.handleMessage(telebot.Message{Chat: chat, Text: "🍾 Развлечения"})

				Convey("утром", func() {
					mockStorage.EXPECT().GetDayEventsByType("fun").Return([]konfurbot.Event{
						konfurbot.Event{Type: "fun", Short: "WAT", Start: &start, Finish: &finish},
						konfurbot.Event{Type: "fun", Short: "WAT 2"},
					})
					mockTelebot.EXPECT().SendMessage(chat, "17:00 — 19:00: *WAT*\n\nвесь день: *WAT 2*\n\n",
						hasButtons("🌶 Еда", "🔥 Доклады", "💥 Мастер-классы", "🍾 Развлечения", "🚜 Трансфер"))
					bot.handleMessage(telebot.Message{Chat: chat, Text: "🍼 Утром"})
				})

				Convey("вечером", func() {
					mockStorage.EXPECT().GetNightEventsByType("fun").Return([]konfurbot.Event{
						konfurbot.Event{Type: "talk", Short: "WAT"},
						konfurbot.Event{Type: "talk", Short: "WAT 2", Start: &start, Finish: &finish},
					})
					mockTelebot.EXPECT().SendMessage(chat, "весь день: *WAT*\n\n17:00 — 19:00: *WAT 2*\n\n",
						hasButtons("🌶 Еда", "🔥 Доклады", "💥 Мастер-классы", "🍾 Развлечения", "🚜 Трансфер"))
					bot.handleMessage(telebot.Message{Chat: chat, Text: "🍸 Вечером"})
				})

				Convey("пользователь пишет нам ерунду", func() {
					mockTelebot.EXPECT().SendMessage(chat, "Я не понимаю эту команду. Давай попробуем еще раз с начала.",
						hasButtons("🌶 Еда", "🔥 Доклады", "💥 Мастер-классы", "🍾 Развлечения", "🚜 Трансфер"))
					bot.handleMessage(telebot.Message{Chat: chat, Text: "gibberish"})
				})
			})

			Convey("пользователь спрашивает про трансфер", func() {
				mockTelebot.EXPECT().SendMessage(chat, "Окей, на каком маршруте поедем?",
					hasButtons("🏎 Дежурный", "🚲 Цветные"))
				bot.handleMessage(telebot.Message{Chat: chat, Text: "🚜 Трансфер"})

				Convey("дежурный", func() {
					mockTelebot.EXPECT().SendMessage(chat, "Расписание довольно большое, может только ближайшие рейсы показать?",
						hasButtons("🔜 Ближайшие", "📜 Все рейсы"))
					bot.handleMessage(telebot.Message{Chat: chat, Text: "🏎 Дежурный"})

					Convey("ближайшие", func() {
						mockStorage.EXPECT().GetNextEventsByType("transfer_main", gomock.Any(), 2*time.Hour).Return([]konfurbot.Event{
							konfurbot.Event{Type: "transfer", Short: "Куда-то вдаль", Start: &start},
							konfurbot.Event{Type: "transfer", Short: "Куда-то вдаль 2", Start: &start},
						})
						mockTelebot.EXPECT().SendMessage(chat, "17:00: *Куда-то вдаль*\n\n17:00: *Куда-то вдаль 2*\n\n",
							hasButtons("🌶 Еда", "🔥 Доклады", "💥 Мастер-классы", "🍾 Развлечения", "🚜 Трансфер"))
						bot.handleMessage(telebot.Message{Chat: chat, Text: "🔜 Ближайшие"})
					})

					Convey("все", func() {
						mockStorage.EXPECT().GetEventsByType("transfer_main").Return([]konfurbot.Event{
							konfurbot.Event{Type: "transfer", Short: "Куда-то вдаль", Start: &start},
							konfurbot.Event{Type: "transfer", Short: "Куда-то вдаль 2", Start: &start},
						})
						mockTelebot.EXPECT().SendMessage(chat, "17:00: *Куда-то вдаль*\n\n17:00: *Куда-то вдаль 2*\n\n",
							hasButtons("🌶 Еда", "🔥 Доклады", "💥 Мастер-классы", "🍾 Развлечения", "🚜 Трансфер"))
						bot.handleMessage(telebot.Message{Chat: chat, Text: "📜 Все рейсы"})
					})

					Convey("пользователь пишет нам ерунду", func() {
						mockTelebot.EXPECT().SendMessage(chat, "Я не понимаю эту команду. Давай попробуем еще раз с начала.",
							hasButtons("🌶 Еда", "🔥 Доклады", "💥 Мастер-классы", "🍾 Развлечения", "🚜 Трансфер"))
						bot.handleMessage(telebot.Message{Chat: chat, Text: "gibberish"})
					})
				})

				Convey("цветные", func() {
					mockTelebot.EXPECT().SendMessage(chat, "Расписание довольно большое, может только ближайшие рейсы показать?",
						hasButtons("🔜 Ближайшие", "📜 Все рейсы"))
					bot.handleMessage(telebot.Message{Chat: chat, Text: "🚲 Цветные"})

					Convey("ближайшие", func() {
						mockStorage.EXPECT().GetNextEventsByType("transfer_color", gomock.Any(), 2*time.Hour).Return([]konfurbot.Event{
							konfurbot.Event{Type: "transfer", Short: "Куда-то вдаль", Start: &start},
							konfurbot.Event{Type: "transfer", Short: "Куда-то вдаль 2", Start: &start},
						})
						mockTelebot.EXPECT().SendMessage(chat, "17:00: *Куда-то вдаль*\n\n17:00: *Куда-то вдаль 2*\n\n",
							hasButtons("🌶 Еда", "🔥 Доклады", "💥 Мастер-классы", "🍾 Развлечения", "🚜 Трансфер"))
						bot.handleMessage(telebot.Message{Chat: chat, Text: "🔜 Ближайшие"})
					})

					Convey("все", func() {
						mockStorage.EXPECT().GetEventsByType("transfer_color").Return([]konfurbot.Event{
							konfurbot.Event{Type: "transfer", Short: "Куда-то вдаль", Start: &start},
							konfurbot.Event{Type: "transfer", Short: "Куда-то вдаль 2", Start: &start},
						})
						mockTelebot.EXPECT().SendMessage(chat, "17:00: *Куда-то вдаль*\n\n17:00: *Куда-то вдаль 2*\n\n",
							hasButtons("🌶 Еда", "🔥 Доклады", "💥 Мастер-классы", "🍾 Развлечения", "🚜 Трансфер"))
						bot.handleMessage(telebot.Message{Chat: chat, Text: "📜 Все рейсы"})
					})

					Convey("пользователь пишет нам ерунду", func() {
						mockTelebot.EXPECT().SendMessage(chat, "Я не понимаю эту команду. Давай попробуем еще раз с начала.",
							hasButtons("🌶 Еда", "🔥 Доклады", "💥 Мастер-классы", "🍾 Развлечения", "🚜 Трансфер"))
						bot.handleMessage(telebot.Message{Chat: chat, Text: "gibberish"})
					})
				})

				Convey("пользователь пишет нам ерунду", func() {
					mockTelebot.EXPECT().SendMessage(chat, "Я не понимаю эту команду. Давай попробуем еще раз с начала.",
						hasButtons("🌶 Еда", "🔥 Доклады", "💥 Мастер-классы", "🍾 Развлечения", "🚜 Трансфер"))
					bot.handleMessage(telebot.Message{Chat: chat, Text: "gibberish"})
				})
			})

			Convey("пользователь пишет нам ерунду", func() {
				mockTelebot.EXPECT().SendMessage(chat, "Я не понимаю эту команду. Давай попробуем еще раз с начала.",
					hasButtons("🌶 Еда", "🔥 Доклады", "💥 Мастер-классы", "🍾 Развлечения", "🚜 Трансфер"))
				bot.handleMessage(telebot.Message{Chat: chat, Text: "gibberish"})
			})
		})
	})
}

func hasButtons(buttons ...string) hasButtonsMatcher {
	sort.Strings(buttons)
	return hasButtonsMatcher{buttons}
}

type hasButtonsMatcher struct {
	buttons []string
}

func (h hasButtonsMatcher) Matches(x interface{}) bool {
	sendOptions := x.(*telebot.SendOptions)

	var flatButtons []string
	for _, buttonRow := range sendOptions.ReplyMarkup.CustomKeyboard {
		for _, button := range buttonRow {
			flatButtons = append(flatButtons, button)
		}
	}
	sort.Strings(flatButtons)
	return reflect.DeepEqual(h.buttons, flatButtons)
}

func (h hasButtonsMatcher) String() string {
	return fmt.Sprintf("has buttons %v", h.buttons)
}
