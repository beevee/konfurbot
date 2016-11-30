package yaml

import (
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
		if parsedEvent.Start == "" {
			parsedEvent.Start = "00:00"
		}
		startTime, err := time.Parse("15:04", parsedEvent.Start)
		if err != nil {
			return err
		}
		start := baseDate.Add(time.Duration(startTime.Hour())*time.Hour + time.Duration(startTime.Minute())*time.Minute)

		if parsedEvent.Finish == "" {
			parsedEvent.Finish = "23:59"
		}
		finishTime, err := time.Parse("15:04", parsedEvent.Finish)
		if err != nil {
			return err
		}
		finish := baseDate.Add(time.Duration(finishTime.Hour())*time.Hour + time.Duration(finishTime.Minute())*time.Minute)
		if finish.Before(start) {
			finish = finish.Add(24 * time.Hour)
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
