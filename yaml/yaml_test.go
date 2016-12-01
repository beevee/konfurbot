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

	tz, _ := time.LoadLocation("Asia/Yekaterinburg")

	Convey("Нам передали валидный YAML с расписанием", t, func() {
		sYaml := []byte(`date: 29.11.2016
night: 20:00
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
    speaker: Александр Казаков
    venue: Кафе
    short: Что делает умный инженер, когда у него падает сервис
    long: |
      "Работает же, чего ты начал" — такой фразой обычно встречают инициативу зафиксировать хронологию и причины факапа.
      Действительно, когда пожар потушен, у нас резко падает мотивация проследить историю развития факапа и покопаться в деталях. А между тем, именно вдумчивый анализ хронологии инцидента позволяет найти глубинные и самые важные проблемы в разработке и эксплуатации, сделать неожиданные выводы и сформулировать задачи, которые качественно повысят стабильность сервиса.
    start: 19:30
    finish: 20:00

  - type: master
    speaker: Василий Петров
    venue: Переговорка 711
    short: Мастер-класс по приготовлению пиццы.
    long: |
      Пицца пицца пицца
      Пицца пицца
    start: 11:30
    finish: 14:25

  - type: fun
    short: Боулинг
    venue: Возле лифтов

  - type: fun
    short: Ночная забава
    start: 01:00+
    finish: 02:00+`)

		Convey("и мы наполнили этим расписанием сторадж", func() {
			mockStorage := mock.NewMockScheduleStorage(ctrl)

			nightCutoff, _ := time.ParseInLocation("15:04 02.01.2006", "20:00 29.11.2016", tz)
			mockStorage.EXPECT().SetNightCutoff(nightCutoff)

			start, _ := time.ParseInLocation("15:04 02.01.2006", "09:00 29.11.2016", tz)
			finish, _ := time.ParseInLocation("15:04 02.01.2006", "10:00 29.11.2016", tz)
			mockStorage.EXPECT().AddEvent(konfurbot.Event{
				Type:   "food",
				Short:  "Утренний кофе",
				Long:   "",
				Start:  &start,
				Finish: &finish,
			})

			start2, _ := time.ParseInLocation("15:04 02.01.2006", "14:30 29.11.2016", tz)
			finish2, _ := time.ParseInLocation("15:04 02.01.2006", "15:00 29.11.2016", tz)
			mockStorage.EXPECT().AddEvent(konfurbot.Event{
				Type:   "food",
				Short:  "Кофе-брейк",
				Long:   "",
				Start:  &start2,
				Finish: &finish2,
			})

			start3, _ := time.ParseInLocation("15:04 02.01.2006", "19:30 29.11.2016", tz)
			finish3, _ := time.ParseInLocation("15:04 02.01.2006", "20:00 29.11.2016", tz)
			mockStorage.EXPECT().AddEvent(konfurbot.Event{
				Type:   "food",
				Short:  "Ужин",
				Long:   "",
				Start:  &start3,
				Finish: &finish3,
			})

			start4, _ := time.ParseInLocation("15:04 02.01.2006", "19:30 29.11.2016", tz)
			finish4, _ := time.ParseInLocation("15:04 02.01.2006", "20:00 29.11.2016", tz)
			mockStorage.EXPECT().AddEvent(konfurbot.Event{
				Type:    "talk",
				Speaker: "Александр Казаков",
				Venue:   "Кафе",
				Short:   "Что делает умный инженер, когда у него падает сервис",
				Long:    "\"Работает же, чего ты начал\" — такой фразой обычно встречают инициативу зафиксировать хронологию и причины факапа.\nДействительно, когда пожар потушен, у нас резко падает мотивация проследить историю развития факапа и покопаться в деталях. А между тем, именно вдумчивый анализ хронологии инцидента позволяет найти глубинные и самые важные проблемы в разработке и эксплуатации, сделать неожиданные выводы и сформулировать задачи, которые качественно повысят стабильность сервиса.\n",
				Start:   &start4,
				Finish:  &finish4,
			})

			start5, _ := time.ParseInLocation("15:04 02.01.2006", "11:30 29.11.2016", tz)
			finish5, _ := time.ParseInLocation("15:04 02.01.2006", "14:25 29.11.2016", tz)
			mockStorage.EXPECT().AddEvent(konfurbot.Event{
				Type:    "master",
				Speaker: "Василий Петров",
				Venue:   "Переговорка 711",
				Short:   "Мастер-класс по приготовлению пиццы.",
				Long:    "Пицца пицца пицца\nПицца пицца\n",
				Start:   &start5,
				Finish:  &finish5,
			})

			mockStorage.EXPECT().AddEvent(konfurbot.Event{
				Type:  "fun",
				Venue: "Возле лифтов",
				Short: "Боулинг",
			})

			start6, _ := time.ParseInLocation("15:04 02.01.2006", "01:00 30.11.2016", tz)
			finish6, _ := time.ParseInLocation("15:04 02.01.2006", "02:00 30.11.2016", tz)
			mockStorage.EXPECT().AddEvent(konfurbot.Event{
				Type:   "fun",
				Short:  "Ночная забава",
				Start:  &start6,
				Finish: &finish6,
			})

			err := FillScheduleStorage(mockStorage, sYaml, tz)
			So(err, ShouldBeNil)
		})
	})

	Convey("Нам передали валидный YAML, но в нем некорректное время начала ночи", t, func() {
		sYaml := []byte(`date: 29.11.2016
night: AAAAAAA
events:
  - type: food
    short: Утренний кофе
    start: 09:00
    finish: 10:00`)

		Convey("поэтому мы не смогли наполнить сторадж", func() {
			mockStorage := mock.NewMockScheduleStorage(ctrl)
			err := FillScheduleStorage(mockStorage, sYaml, tz)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Нам передали валидный YAML, но в нем некорректное время завершения события", t, func() {
		sYaml := []byte(`date: 29.11.2016
night: 20:00
events:
  - type: food
    short: Утренний кофе
    start: 09:00
    finish: 嘘`)

		Convey("поэтому мы не смогли наполнить сторадж", func() {
			mockStorage := mock.NewMockScheduleStorage(ctrl)

			nightCutoff, _ := time.ParseInLocation("15:04 02.01.2006", "20:00 29.11.2016", tz)
			mockStorage.EXPECT().SetNightCutoff(nightCutoff)

			err := FillScheduleStorage(mockStorage, sYaml, tz)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Нам передали валидный YAML, но в нем некорректное время начала события", t, func() {
		sYaml := []byte(`date: 29.11.2016
night: 20:00
events:
  - type: food
    short: Утренний кофе
    start: 不正解
    finish: 10:00`)

		Convey("поэтому мы не смогли наполнить сторадж", func() {
			mockStorage := mock.NewMockScheduleStorage(ctrl)

			nightCutoff, _ := time.ParseInLocation("15:04 02.01.2006", "20:00 29.11.2016", tz)
			mockStorage.EXPECT().SetNightCutoff(nightCutoff)

			err := FillScheduleStorage(mockStorage, sYaml, tz)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Нам передали валидный YAML, но в нем некорректная дата конференции", t, func() {
		sYaml := []byte(`date: かわいい！＾＝＾
night: 20:00
events:
  - type: food
    short: Утренний кофе
    start: 09:00
    finish: 10:00`)

		Convey("поэтому мы не смогли наполнить сторадж", func() {
			mockStorage := mock.NewMockScheduleStorage(ctrl)
			err := FillScheduleStorage(mockStorage, sYaml, tz)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Нам передали полностью невалидный YAML", t, func() {
		sYaml := []byte(`	1123123`)

		Convey("поэтому мы не смогли наполнить сторадж", func() {
			mockStorage := mock.NewMockScheduleStorage(ctrl)
			err := FillScheduleStorage(mockStorage, sYaml, tz)
			So(err, ShouldNotBeNil)
		})
	})
}
