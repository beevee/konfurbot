package yaml

import (
	"testing"
	"time"

	"github.com/beevee/konfurbot"
	. "github.com/smartystreets/goconvey/convey"
)

func TestYaml(t *testing.T) {
	Convey("Given valid YAML schedule", t, func() {
		sYaml := []byte(`- type: food
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
  short: WAT
  long: |
    В докладе пойдет речь о том,
    чему равна сумма объекта и пустой строки.
  start: 12:00
  finish: 13:00

- type: fun
  short: ЧГК
  start: 16:30
  finish: 17:00`)

		Convey("It should parse into a schedule structure", func() {
			s, err := ParseSchedule(sYaml)
			So(err, ShouldBeNil)

			Convey("Its entries should be correct and in the same order", func() {
				So(len(s.Events["food"]), ShouldEqual, 3)
				So(len(s.Events["fun"]), ShouldEqual, 1)
				So(len(s.Events["talk"]), ShouldEqual, 1)

				var start, finish time.Time

				start, _ = time.Parse("15:04", "09:00")
				finish, _ = time.Parse("15:04", "10:00")
				So(s.Events["food"][0], ShouldResemble, konfurbot.Event{
					Type:   "food",
					Short:  "Утренний кофе",
					Long:   "",
					Start:  start,
					Finish: finish,
				})

				start, _ = time.Parse("15:04", "14:30")
				finish, _ = time.Parse("15:04", "15:00")
				So(s.Events["food"][1], ShouldResemble, konfurbot.Event{
					Type:   "food",
					Short:  "Кофе-брейк",
					Long:   "",
					Start:  start,
					Finish: finish,
				})

				start, _ = time.Parse("15:04", "19:30")
				finish, _ = time.Parse("15:04", "20:00")
				So(s.Events["food"][2], ShouldResemble, konfurbot.Event{
					Type:   "food",
					Short:  "Ужин",
					Long:   "",
					Start:  start,
					Finish: finish,
				})

				start, _ = time.Parse("15:04", "12:00")
				finish, _ = time.Parse("15:04", "13:00")
				So(s.Events["talk"][0], ShouldResemble, konfurbot.Event{
					Type:   "talk",
					Short:  "WAT",
					Long:   "В докладе пойдет речь о том,\nчему равна сумма объекта и пустой строки.\n",
					Start:  start,
					Finish: finish,
				})

				start, _ = time.Parse("15:04", "16:30")
				finish, _ = time.Parse("15:04", "17:00")
				So(s.Events["fun"][0], ShouldResemble, konfurbot.Event{
					Type:   "fun",
					Short:  "ЧГК",
					Long:   "",
					Start:  start,
					Finish: finish,
				})
			})
		})
	})

	Convey("Given generally valid YAML schedule with invalid start time", t, func() {
		sYaml := []byte(`- type: food
  short: Утренний кофе
  start: 09:00
  finish: 嘘`)

		Convey("It should not parse", func() {
			_, err := ParseSchedule(sYaml)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Given generally valid YAML schedule with invalid finish time", t, func() {
		sYaml := []byte(`- type: food
  short: Утренний кофе
  start: 不正解
  finish: 10:00`)

		Convey("It should not parse", func() {
			_, err := ParseSchedule(sYaml)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Given totally invalid YAML schedule", t, func() {
		sYaml := []byte(`	1123123`)

		Convey("It should not parse", func() {
			_, err := ParseSchedule(sYaml)
			So(err, ShouldNotBeNil)
		})
	})
}
