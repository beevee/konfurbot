package konfurbot

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSchedule(t *testing.T) {
	Convey("We fill the schedule with events", t, func() {
		var start, finish time.Time
		schedule := Schedule{}

		start, _ = time.Parse("15:04", "14:30")
		finish, _ = time.Parse("15:04", "15:00")
		food0 := Event{
			Type:   "food",
			Short:  "Кофе-брейк",
			Long:   "",
			Start:  start,
			Finish: finish,
		}
		schedule.AddEvent(food0)

		start, _ = time.Parse("15:04", "19:30")
		finish, _ = time.Parse("15:04", "20:00")
		food1 := Event{
			Type:   "food",
			Short:  "Ужин",
			Long:   "",
			Start:  start,
			Finish: finish,
		}
		schedule.AddEvent(food1)

		start, _ = time.Parse("15:04", "12:00")
		finish, _ = time.Parse("15:04", "13:00")
		schedule.AddEvent(Event{
			Type:   "talk",
			Short:  "WAT",
			Long:   "В докладе пойдет речь о том,\nчему равна сумма объекта и пустой строки.\n",
			Start:  start,
			Finish: finish,
		})

		Convey("so we can retrieve all events of a single type in the same order", func() {
			foodEvents := schedule.GetEventsByType("food")

			So(foodEvents, ShouldHaveLength, 2)
			So(foodEvents[0], ShouldResemble, food0)
			So(foodEvents[1], ShouldResemble, food1)
		})
	})
}
