package telegram

import (
	"fmt"
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

func TestTelegram(t *testing.T) {
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
				hasButtons("🌶 Еда", "🔥 Доклады / МК", "🍾 Развлечения", "🚜 Трансфер"))
			bot.handleMessage(telebot.Message{Chat: chat, Text: "/start"})

			Convey("пользователь спрашивает про еду (у еды есть место проведения, но нет спикера)", func() {
				mockStorage.EXPECT().GetEventsByType("food").Return([]konfurbot.Event{
					konfurbot.Event{Type: "food", Short: "お好み焼き", Venue: "Бар", Start: &start, Finish: &finish},
					konfurbot.Event{Type: "food", Short: "焼き鳥", Long: "Вегетарианцам накроют на крыше паркинга", Venue: "Кафе", Start: &start, Finish: &finish},
				})
				mockTelebot.EXPECT().SendMessage(chat, "17:00 — 19:00 [Бар]: お好み焼き\n17:00 — 19:00 [Кафе]: 焼き鳥\nВегетарианцам накроют на крыше паркинга\n\n",
					hasButtons("🌶 Еда", "🔥 Доклады / МК", "🍾 Развлечения", "🚜 Трансфер"))
				bot.handleMessage(telebot.Message{Chat: chat, Text: "🌶 Еда"})
			})

			Convey("пользователь спрашивает про доклады", func() {
				mockTelebot.EXPECT().SendMessage(chat, "Окей, какие доклады и мастер-классы?",
					hasButtons("Которые идут сейчас", "Которые начнутся скоро", "Все"))
				bot.handleMessage(telebot.Message{Chat: chat, Text: "🔥 Доклады / МК"})

				Convey("которые идут сейчас", func() {
					mockTelebot.EXPECT().SendMessage(chat, "Их может оказаться довольно много. Тизеры надо?",
						hasButtons("С тизерами (простыня!)", "Без тизеров (ура! краткость!)"))
					bot.handleMessage(telebot.Message{Chat: chat, Text: "Которые идут сейчас"})

					Convey("с тизерами, и что-то сейчас идет", func() {
						mockStorage.EXPECT().GetCurrentEventsByType("talk", gomock.Any()).Return([]konfurbot.Event{
							konfurbot.Event{Type: "talk", Short: "WAT", Long: "WAAAAT", Start: &start, Finish: &finish},
							konfurbot.Event{Type: "talk", Short: "WAT 2", Long: "WAAAAT 22", Start: &start, Finish: &finish},
						})
						mockTelebot.EXPECT().SendMessage(chat, "17:00 — 19:00: WAT\nWAAAAT\n\n17:00 — 19:00: WAT 2\nWAAAAT 22\n\n",
							hasButtons("🌶 Еда", "🔥 Доклады / МК", "🍾 Развлечения", "🚜 Трансфер"))
						bot.handleMessage(telebot.Message{Chat: chat, Text: "С тизерами (простыня!)"})
					})

					Convey("с тизерами, и сейчас ничего не идет", func() {
						mockStorage.EXPECT().GetCurrentEventsByType("talk", gomock.Any()).Return([]konfurbot.Event{})
						mockTelebot.EXPECT().SendMessage(chat, "Ничего нет :(",
							hasButtons("🌶 Еда", "🔥 Доклады / МК", "🍾 Развлечения", "🚜 Трансфер"))
						bot.handleMessage(telebot.Message{Chat: chat, Text: "С тизерами (простыня!)"})
					})

					Convey("без тизеров, и сейчас что-то идет", func() {
						mockStorage.EXPECT().GetCurrentEventsByType("talk", gomock.Any()).Return([]konfurbot.Event{
							konfurbot.Event{Type: "talk", Short: "WAT", Long: "WAAAAT", Start: &start, Finish: &finish},
							konfurbot.Event{Type: "talk", Short: "WAT 2", Long: "WAAAAT 22", Start: &start, Finish: &finish},
						})
						mockTelebot.EXPECT().SendMessage(chat, "17:00 — 19:00: WAT\n17:00 — 19:00: WAT 2\n",
							hasButtons("🌶 Еда", "🔥 Доклады / МК", "🍾 Развлечения", "🚜 Трансфер"))
						bot.handleMessage(telebot.Message{Chat: chat, Text: "Без тизеров (ура! краткость!)"})
					})

					Convey("пользователь пишет нам ерунду", func() {
						mockTelebot.EXPECT().SendMessage(chat, "Я не понимаю эту команду. Давай попробуем еще раз с начала.",
							hasButtons("🌶 Еда", "🔥 Доклады / МК", "🍾 Развлечения", "🚜 Трансфер"))
						bot.handleMessage(telebot.Message{Chat: chat, Text: "gibberish"})
					})
				})

				Convey("которые начнутся в ближайший час", func() {
					mockTelebot.EXPECT().SendMessage(chat, "Их может оказаться довольно много. Тизеры надо?",
						hasButtons("С тизерами (простыня!)", "Без тизеров (ура! краткость!)"))
					bot.handleMessage(telebot.Message{Chat: chat, Text: "Которые начнутся скоро"})

					Convey("с тизерами", func() {
						mockStorage.EXPECT().GetNextEventsByType("talk", gomock.Any(), time.Hour).Return([]konfurbot.Event{
							konfurbot.Event{Type: "talk", Short: "WAT", Long: "WAAAAT", Start: &start, Finish: &finish},
							konfurbot.Event{Type: "talk", Short: "WAT 2", Long: "WAAAAT 22", Start: &start, Finish: &finish},
						})
						mockTelebot.EXPECT().SendMessage(chat, "17:00 — 19:00: WAT\nWAAAAT\n\n17:00 — 19:00: WAT 2\nWAAAAT 22\n\n",
							hasButtons("🌶 Еда", "🔥 Доклады / МК", "🍾 Развлечения", "🚜 Трансфер"))
						bot.handleMessage(telebot.Message{Chat: chat, Text: "С тизерами (простыня!)"})
					})

					Convey("без тизеров", func() {
						mockStorage.EXPECT().GetNextEventsByType("talk", gomock.Any(), time.Hour).Return([]konfurbot.Event{
							konfurbot.Event{Type: "talk", Short: "WAT", Long: "WAAAAT", Start: &start, Finish: &finish},
							konfurbot.Event{Type: "talk", Short: "WAT 2", Long: "WAAAAT 22", Start: &start, Finish: &finish},
						})
						mockTelebot.EXPECT().SendMessage(chat, "17:00 — 19:00: WAT\n17:00 — 19:00: WAT 2\n",
							hasButtons("🌶 Еда", "🔥 Доклады / МК", "🍾 Развлечения", "🚜 Трансфер"))
						bot.handleMessage(telebot.Message{Chat: chat, Text: "Без тизеров (ура! краткость!)"})
					})

					Convey("пользователь пишет нам ерунду", func() {
						mockTelebot.EXPECT().SendMessage(chat, "Я не понимаю эту команду. Давай попробуем еще раз с начала.",
							hasButtons("🌶 Еда", "🔥 Доклады / МК", "🍾 Развлечения", "🚜 Трансфер"))
						bot.handleMessage(telebot.Message{Chat: chat, Text: "gibberish"})
					})
				})

				Convey("все", func() {
					mockTelebot.EXPECT().SendMessage(chat, "Полное расписание довольно длинное. Давай посмотрим отдельно, доклады или мастер-классы? С тизерами вообще не буду предлагать :)",
						hasButtons("Доклады", "Мастер-классы"))
					bot.handleMessage(telebot.Message{Chat: chat, Text: "Все"})

					Convey("доклады (у докладов есть спикер и место проведения)", func() {
						mockStorage.EXPECT().GetEventsByTypeAndSubtype("talk", "talk").Return([]konfurbot.Event{
							konfurbot.Event{
								Type:    "talk",
								Subtype: "talk",
								Speaker: "Александр Казаков",
								Venue:   "Учебный класс 1",
								Short:   "WAT",
								Long:    "WAAAAT",
								Start:   &start,
								Finish:  &finish,
							},
							konfurbot.Event{
								Type:    "talk",
								Subtype: "talk",
								Speaker: "Василий Петров",
								Venue:   "Учебный класс 2",
								Short:   "WAT 2",
								Long:    "WAAAAT 22",
								Start:   &start,
								Finish:  &finish,
							},
						})
						mockTelebot.EXPECT().SendMessage(chat, "17:00 — 19:00 [Учебный класс 1]: WAT (Александр Казаков)\n17:00 — 19:00 [Учебный класс 2]: WAT 2 (Василий Петров)\n",
							hasButtons("🌶 Еда", "🔥 Доклады / МК", "🍾 Развлечения", "🚜 Трансфер"))
						bot.handleMessage(telebot.Message{Chat: chat, Text: "Доклады"})
					})

					Convey("мастер-классы (у мастер-классов есть спикер, но нет места проведения)", func() {
						mockStorage.EXPECT().GetEventsByTypeAndSubtype("talk", "master").Return([]konfurbot.Event{
							konfurbot.Event{
								Type:    "talk",
								Subtype: "master",
								Short:   "WAT",
								Long:    "WAAAAT",
								Speaker: "Александр Казаков",
								Start:   &start,
								Finish:  &finish,
							},
							konfurbot.Event{
								Type:    "talk",
								Subtype: "master",
								Short:   "WAT 2",
								Long:    "WAAAAT 22",
								Speaker: "Василий Петров",
								Start:   &start,
								Finish:  &finish,
							},
						})
						mockTelebot.EXPECT().SendMessage(chat, "17:00 — 19:00: WAT (Александр Казаков)\n17:00 — 19:00: WAT 2 (Василий Петров)\n",
							hasButtons("🌶 Еда", "🔥 Доклады / МК", "🍾 Развлечения", "🚜 Трансфер"))
						bot.handleMessage(telebot.Message{Chat: chat, Text: "Мастер-классы"})
					})

					Convey("пользователь пишет нам ерунду", func() {
						mockTelebot.EXPECT().SendMessage(chat, "Я не понимаю эту команду. Давай попробуем еще раз с начала.",
							hasButtons("🌶 Еда", "🔥 Доклады / МК", "🍾 Развлечения", "🚜 Трансфер"))
						bot.handleMessage(telebot.Message{Chat: chat, Text: "gibberish"})
					})
				})

				Convey("пользователь пишет нам ерунду", func() {
					mockTelebot.EXPECT().SendMessage(chat, "Я не понимаю эту команду. Давай попробуем еще раз с начала.",
						hasButtons("🌶 Еда", "🔥 Доклады / МК", "🍾 Развлечения", "🚜 Трансфер"))
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
					mockTelebot.EXPECT().SendMessage(chat, "17:00 — 19:00: WAT\nвесь день: WAT 2\n",
						hasButtons("🌶 Еда", "🔥 Доклады / МК", "🍾 Развлечения", "🚜 Трансфер"))
					bot.handleMessage(telebot.Message{Chat: chat, Text: "🍼 Утром"})
				})

				Convey("вечером", func() {
					mockStorage.EXPECT().GetNightEventsByType("fun").Return([]konfurbot.Event{
						konfurbot.Event{Type: "talk", Short: "WAT"},
						konfurbot.Event{Type: "talk", Short: "WAT 2", Start: &start, Finish: &finish},
					})
					mockTelebot.EXPECT().SendMessage(chat, "весь день: WAT\n17:00 — 19:00: WAT 2\n",
						hasButtons("🌶 Еда", "🔥 Доклады / МК", "🍾 Развлечения", "🚜 Трансфер"))
					bot.handleMessage(telebot.Message{Chat: chat, Text: "🍸 Вечером"})
				})

				Convey("пользователь пишет нам ерунду", func() {
					mockTelebot.EXPECT().SendMessage(chat, "Я не понимаю эту команду. Давай попробуем еще раз с начала.",
						hasButtons("🌶 Еда", "🔥 Доклады / МК", "🍾 Развлечения", "🚜 Трансфер"))
					bot.handleMessage(telebot.Message{Chat: chat, Text: "gibberish"})
				})
			})

			Convey("пользователь пишет нам ерунду", func() {
				mockTelebot.EXPECT().SendMessage(chat, "Я не понимаю эту команду. Давай попробуем еще раз с начала.",
					hasButtons("🌶 Еда", "🔥 Доклады / МК", "🍾 Развлечения", "🚜 Трансфер"))
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
