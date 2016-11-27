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

	Convey("We were given a valid YAML schedule", t, func() {
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

		Convey("so we were able to parse it and fill a storage", func() {
			var start, finish time.Time
			s := mock.NewMockScheduleStorage(ctrl)

			start, _ = time.Parse("15:04", "09:00")
			finish, _ = time.Parse("15:04", "10:00")
			s.EXPECT().AddEvent(konfurbot.Event{
				Type:   "food",
				Short:  "Утренний кофе",
				Long:   "",
				Start:  start,
				Finish: finish,
			})

			start, _ = time.Parse("15:04", "14:30")
			finish, _ = time.Parse("15:04", "15:00")
			s.EXPECT().AddEvent(konfurbot.Event{
				Type:   "food",
				Short:  "Кофе-брейк",
				Long:   "",
				Start:  start,
				Finish: finish,
			})

			start, _ = time.Parse("15:04", "19:30")
			finish, _ = time.Parse("15:04", "20:00")
			s.EXPECT().AddEvent(konfurbot.Event{
				Type:   "food",
				Short:  "Ужин",
				Long:   "",
				Start:  start,
				Finish: finish,
			})

			start, _ = time.Parse("15:04", "12:00")
			finish, _ = time.Parse("15:04", "13:00")
			s.EXPECT().AddEvent(konfurbot.Event{
				Type:   "talk",
				Short:  "WAT",
				Long:   "В докладе пойдет речь о том,\nчему равна сумма объекта и пустой строки.\n",
				Start:  start,
				Finish: finish,
			})

			start, _ = time.Parse("15:04", "16:30")
			finish, _ = time.Parse("15:04", "17:00")
			s.EXPECT().AddEvent(konfurbot.Event{
				Type:   "fun",
				Short:  "ЧГК",
				Long:   "",
				Start:  start,
				Finish: finish,
			})

			err := FillScheduleStorage(s, sYaml)
			So(err, ShouldBeNil)
		})
	})

	Convey("We were given generally valid YAML schedule with invalid start time", t, func() {
		sYaml := []byte(`- type: food
  short: Утренний кофе
  start: 09:00
  finish: 嘘`)

		Convey("so we were not able to parse it", func() {
			s := mock.NewMockScheduleStorage(ctrl)
			err := FillScheduleStorage(s, sYaml)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("We were given generally valid YAML schedule with invalid finish time", t, func() {
		sYaml := []byte(`- type: food
  short: Утренний кофе
  start: 不正解
  finish: 10:00`)

		Convey("so we were not able to parse it", func() {
			s := mock.NewMockScheduleStorage(ctrl)
			err := FillScheduleStorage(s, sYaml)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("We were given totally invalid YAML schedule", t, func() {
		sYaml := []byte(`	1123123`)

		Convey("so we were not able to parse it", func() {
			s := mock.NewMockScheduleStorage(ctrl)
			err := FillScheduleStorage(s, sYaml)
			So(err, ShouldNotBeNil)
		})
	})
}
