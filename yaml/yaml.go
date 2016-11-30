package yaml

import (
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/beevee/konfurbot"
)

// Schedule is a YAML-parseable adapter for konfurbot.Schedule
type Schedule struct {
	Date   string `yaml:"date"`
	Night  string `yaml:"night"`
	Events []struct {
		Type    string `yaml:"type"`
		Subtype string `yaml:"subtype"`
		Speaker string `yaml:"speaker"`
		Venue   string `yaml:"venue"`
		Short   string `yaml:"short"`
		Long    string `yaml:"long"`
		Start   string `yaml:"start"`
		Finish  string `yaml:"finish"`
	} `yaml:"events"`
}

// FillScheduleStorage fills schedule storage with events parsed from YAML file
func FillScheduleStorage(storage konfurbot.ScheduleStorage, file []byte) error {
	var parsedSchedule Schedule
	if err := yaml.Unmarshal(file, &parsedSchedule); err != nil {
		return err
	}

	baseDate, err := time.Parse("02.01.2006", parsedSchedule.Date)
	if err != nil {
		return err
	}

	nightCutoff, err := time.Parse("15:04", parsedSchedule.Night)
	if err != nil {
		return err
	}
	storage.SetNightCutoff(baseDate.Add(time.Duration(nightCutoff.Hour())*time.Hour + time.Duration(nightCutoff.Minute())*time.Minute))

	for _, parsedEvent := range parsedSchedule.Events {
		var start, finish *time.Time

		if parsedEvent.Start == "" {
			start = nil
		} else {
			startTime, shift, err := parseRichTime(parsedEvent.Start)
			if err != nil {
				return err
			}
			startTime = baseDate.Add(time.Duration(startTime.Hour())*time.Hour + time.Duration(startTime.Minute())*time.Minute)
			if shift {
				startTime = startTime.Add(24 * time.Hour)
			}
			start = &startTime
		}

		if parsedEvent.Finish == "" {
			finish = nil
		} else {
			finishTime, shift, err := parseRichTime(parsedEvent.Finish)
			if err != nil {
				return err
			}
			finishTime = baseDate.Add(time.Duration(finishTime.Hour())*time.Hour + time.Duration(finishTime.Minute())*time.Minute)
			if shift {
				finishTime = finishTime.Add(24 * time.Hour)
			}
			finish = &finishTime
		}

		event := konfurbot.Event{
			Type:    parsedEvent.Type,
			Subtype: parsedEvent.Subtype,
			Speaker: parsedEvent.Speaker,
			Venue:   parsedEvent.Venue,
			Short:   parsedEvent.Short,
			Long:    parsedEvent.Long,
			Start:   start,
			Finish:  finish,
		}
		storage.AddEvent(event)
	}

	return nil
}

func parseRichTime(timeString string) (time.Time, bool, error) {
	var timeShift bool
	if strings.HasSuffix(timeString, "+") {
		timeString = timeString[:len(timeString)-1]
		timeShift = true
	}
	timeParsed, err := time.Parse("15:04", timeString)
	if err != nil {
		return time.Time{}, false, err
	}
	return timeParsed, timeShift, nil
}
