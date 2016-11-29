package konfurbot

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSchedule(t *testing.T) {
	Convey("Мы наполняем расписание событиями", t, func() {
		start := time.Now().Add(-30 * time.Minute)
		finish := start.Add(time.Hour)

		schedule := Schedule{}

		food0 := Event{
			Type:   "food",
			Short:  "Кофе-брейк",
			Long:   "",
			Start:  start,
			Finish: finish,
		}
		schedule.AddEvent(food0)

		food1 := Event{
			Type:   "food",
			Short:  "Ужин",
			Long:   "",
			Start:  start,
			Finish: finish,
		}
		schedule.AddEvent(food1)

		talk0 := Event{
			Type:    "talk",
			Subtype: "talk",
			Short:   "WAT",
			Long:    "В докладе пойдет речь о том,\nчему равна сумма объекта и пустой строки.\n",
			Start:   start,
			Finish:  finish,
		}
		schedule.AddEvent(talk0)

		Convey("а потом получаем все события определенного типа в том же порядке", func() {
			foodEvents := schedule.GetEventsByType("food")

			So(foodEvents, ShouldHaveLength, 2)
			So(foodEvents[0], ShouldResemble, food0)
			So(foodEvents[1], ShouldResemble, food1)
		})

		Convey("а потом получаем текущие события определенного типа и подтипа", func() {
			Convey("если они есть", func() {
				talkEvents := schedule.GetEventsByTypeAndSubtype("talk", "talk")

				So(talkEvents, ShouldHaveLength, 1)
				So(talkEvents[0], ShouldResemble, talk0)
			})

			Convey("если их нет", func() {
				talkEvents := schedule.GetEventsByTypeAndSubtype("talk", "gibberish")

				So(talkEvents, ShouldHaveLength, 0)
			})
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
	})
}
