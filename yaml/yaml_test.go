package yaml

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/beevee/konfurbot"
	"github.com/beevee/konfurbot/mock"
)

func TestYaml(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	Convey("Нам передали валидный YAML с расписанием", t, func() {
		sYaml := []byte(`date: 29.11.2016
events:
  - type: food
    short: Утренний кофе
    start: 09:00
    finish: 10:00

  - type: food
    short: Кофе-брейк
    start: 14:30
    finish: 15:00

  - type: food
    short: Ужин
    start: 19:30
    finish: 20:00

  - type: talk
    subtype: talk
    speaker: Александр Казаков
    short: Что делает умный инженер, когда у него падает сервис
    long: |
      "Работает же, чего ты начал" — такой фразой обычно встречают инициативу зафиксировать хронологию и причины факапа.
      Действительно, когда пожар потушен, у нас резко падает мотивация проследить историю развития факапа и покопаться в деталях. А между тем, именно вдумчивый анализ хронологии инцидента позволяет найти глубинные и самые важные проблемы в разработке и эксплуатации, сделать неожиданные выводы и сформулировать задачи, которые качественно повысят стабильность сервиса.
    start: 19:30
    finish: 20:00

  - type: talk
    subtype: master
    speaker: Василий Петров
    short: Мастер-класс по приготовлению пиццы.
    long: |
      Пицца пицца пицца
      Пицца пицца
    start: 11:30
    finish: 14:25`)

		Convey("и мы наполнили этим расписанием сторадж", func() {
			var start, finish time.Time
			mockStorage := mock.NewMockScheduleStorage(ctrl)

			start, _ = time.Parse("15:04 02.01.2006", "09:00 29.11.2016")
			finish, _ = time.Parse("15:04 02.01.2006", "10:00 29.11.2016")
			mockStorage.EXPECT().AddEvent(konfurbot.Event{
				Type:   "food",
				Short:  "Утренний кофе",
				Long:   "",
				Start:  start,
				Finish: finish,
			})

			start, _ = time.Parse("15:04 02.01.2006", "14:30 29.11.2016")
			finish, _ = time.Parse("15:04 02.01.2006", "15:00 29.11.2016")
			mockStorage.EXPECT().AddEvent(konfurbot.Event{
				Type:   "food",
				Short:  "Кофе-брейк",
				Long:   "",
				Start:  start,
				Finish: finish,
			})

			start, _ = time.Parse("15:04 02.01.2006", "19:30 29.11.2016")
			finish, _ = time.Parse("15:04 02.01.2006", "20:00 29.11.2016")
			mockStorage.EXPECT().AddEvent(konfurbot.Event{
				Type:   "food",
				Short:  "Ужин",
				Long:   "",
				Start:  start,
				Finish: finish,
			})

			start, _ = time.Parse("15:04 02.01.2006", "19:30 29.11.2016")
			finish, _ = time.Parse("15:04 02.01.2006", "20:00 29.11.2016")
			mockStorage.EXPECT().AddEvent(konfurbot.Event{
				Type:    "talk",
				Subtype: "talk",
				Speaker: "Александр Казаков",
				Short:   "Что делает умный инженер, когда у него падает сервис",
				Long:    "\"Работает же, чего ты начал\" — такой фразой обычно встречают инициативу зафиксировать хронологию и причины факапа.\nДействительно, когда пожар потушен, у нас резко падает мотивация проследить историю развития факапа и покопаться в деталях. А между тем, именно вдумчивый анализ хронологии инцидента позволяет найти глубинные и самые важные проблемы в разработке и эксплуатации, сделать неожиданные выводы и сформулировать задачи, которые качественно повысят стабильность сервиса.\n",
				Start:   start,
				Finish:  finish,
			})

			start, _ = time.Parse("15:04 02.01.2006", "11:30 29.11.2016")
			finish, _ = time.Parse("15:04 02.01.2006", "14:25 29.11.2016")
			mockStorage.EXPECT().AddEvent(konfurbot.Event{
				Type:    "talk",
				Subtype: "master",
				Speaker: "Василий Петров",
				Short:   "Мастер-класс по приготовлению пиццы.",
				Long:    "Пицца пицца пицца\nПицца пицца\n",
				Start:   start,
				Finish:  finish,
			})

			err := FillScheduleStorage(mockStorage, sYaml)
			So(err, ShouldBeNil)
		})
	})

	Convey("Нам передали валидный YAML, но в нем некорректное время завершения события", t, func() {
		sYaml := []byte(`date: 29.11.2016
events:
  - type: food
    short: Утренний кофе
    start: 09:00
    finish: 嘘`)

		Convey("поэтому мы не смогли наполнить сторадж", func() {
			s := mock.NewMockScheduleStorage(ctrl)
			err := FillScheduleStorage(s, sYaml)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Нам передали валидный YAML, но в нем некорректное время начала события", t, func() {
		sYaml := []byte(`date: 29.11.2016
events:
  - type: food
    short: Утренний кофе
    start: 不正解
    finish: 10:00`)

		Convey("поэтому мы не смогли наполнить сторадж", func() {
			s := mock.NewMockScheduleStorage(ctrl)
			err := FillScheduleStorage(s, sYaml)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Нам передали валидный YAML, но в нем некорректная дата конференции", t, func() {
		sYaml := []byte(`date: かわいい！＾＝＾
events:
  - type: food
    short: Утренний кофе
    start: 09:00
    finish: 10:00`)

		Convey("поэтому мы не смогли наполнить сторадж", func() {
			s := mock.NewMockScheduleStorage(ctrl)
			err := FillScheduleStorage(s, sYaml)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Нам передали полностью невалидный YAML", t, func() {
		sYaml := []byte(`	1123123`)

		Convey("поэтому мы не смогли наполнить сторадж", func() {
			s := mock.NewMockScheduleStorage(ctrl)
			err := FillScheduleStorage(s, sYaml)
			So(err, ShouldNotBeNil)
		})
	})
}
