package konfurbot

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSchedule(t *testing.T) {
	Convey("Мы наполняем расписание событиями", t, func() {
		start := time.Now().Add(-30 * time.Minute)
		finish := time.Now().Add(30 * time.Minute)

		schedule := Schedule{}

		food0 := Event{
			Type:   "food",
			Short:  "Кофе-брейк",
			Long:   "",
			Start:  &start,
			Finish: &finish,
		}
		schedule.AddEvent(food0)

		food1 := Event{
			Type:   "food",
			Short:  "Ужин",
			Long:   "",
			Start:  &start,
			Finish: &finish,
		}
		schedule.AddEvent(food1)

		talk0 := Event{
			Type:   "talk",
			Short:  "WAT",
			Long:   "В докладе пойдет речь о том,\nчему равна сумма объекта и пустой строки.\n",
			Start:  &start,
			Finish: &finish,
		}
		schedule.AddEvent(talk0)

		nightCutoff, _ := time.Parse("02.01.2006 15:04", "26.11.2016 20:00")
		schedule.SetNightCutoff(nightCutoff)

		fun0 := Event{
			Type:  "fun",
			Short: "Боулинг на весь день",
		}
		schedule.AddEvent(fun0)

		startFixed, _ := time.Parse("02.01.2006 15:04", "26.11.2016 10:00")
		finishFixed, _ := time.Parse("02.01.2006 15:04", "26.11.2016 17:00")
		fun1 := Event{
			Type:   "fun",
			Short:  "Клавогонки днем",
			Start:  &startFixed,
			Finish: &finishFixed,
		}
		schedule.AddEvent(fun1)

		startFixed2, _ := time.Parse("02.01.2006 15:04", "26.11.2016 20:00")
		finishFixed2, _ := time.Parse("02.01.2006 15:04", "26.11.2016 23:00")
		fun2 := Event{
			Type:   "fun",
			Short:  "ЧГК вечером",
			Start:  &startFixed2,
			Finish: &finishFixed2,
		}
		schedule.AddEvent(fun2)

		Convey("а потом получаем все события определенного типа в том же порядке", func() {
			foodEvents := schedule.GetEventsByType("food")

			So(foodEvents, ShouldHaveLength, 2)
			So(foodEvents[0], ShouldResemble, food0)
			So(foodEvents[1], ShouldResemble, food1)
		})

		Convey("а потом получаем текущие события определенного типа", func() {
			Convey("если они есть", func() {
				talkEvents := schedule.GetCurrentEventsByType("talk", time.Now())

				So(talkEvents, ShouldHaveLength, 1)
				So(talkEvents[0], ShouldResemble, talk0)
			})

			Convey("если их нет", func() {
				talkEvents := schedule.GetCurrentEventsByType("talk", time.Now().Add(time.Hour))

				So(talkEvents, ShouldHaveLength, 0)
			})
		})

		Convey("а потом получаем ближайшие события определенного типа", func() {
			Convey("если они есть", func() {
				talkEvents := schedule.GetNextEventsByType("talk", time.Now().Add(-time.Hour), time.Hour)

				So(talkEvents, ShouldHaveLength, 1)
				So(talkEvents[0], ShouldResemble, talk0)
			})

			Convey("если их нет", func() {
				talkEvents := schedule.GetNextEventsByType("talk", time.Now(), time.Hour)

				So(talkEvents, ShouldHaveLength, 0)
			})
		})

		Convey("а потом получаем события, которые начинаются раньше определенного времени", func() {
			funDayEvents := schedule.GetDayEventsByType("fun")

			So(funDayEvents, ShouldHaveLength, 2)
			So(funDayEvents[0], ShouldResemble, fun0)
			So(funDayEvents[1], ShouldResemble, fun1)
		})

		Convey("а потом получаем события, которые заканчиваются позже определенного времени", func() {
			funNightEvents := schedule.GetNightEventsByType("fun")

			So(funNightEvents, ShouldHaveLength, 2)
			So(funNightEvents[0], ShouldResemble, fun0)
			So(funNightEvents[1], ShouldResemble, fun2)
		})
	})
}
