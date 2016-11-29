package yaml

import (
	"time"

	"gopkg.in/yaml.v2"

	"github.com/beevee/konfurbot"
)

// Schedule is a YAML-parseable adapter for konfurbot.Schedule
type Schedule struct {
	Date   string `yaml:"date"`
	Events []struct {
		Type    string `yaml:"type"`
		Subtype string `yaml:"subtype"`
		Speaker string `yaml:"speaker"`
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
	for _, parsedEvent := range parsedSchedule.Events {
		startTime, err := time.Parse("15:04", parsedEvent.Start)
		if err != nil {
			return err
		}

		finishTime, err := time.Parse("15:04", parsedEvent.Finish)
		if err != nil {
			return err
		}
		baseDate.Add(time.Duration(finishTime.Hour())*time.Hour + time.Duration(finishTime.Minute())*time.Minute)

		event := konfurbot.Event{
			Type:    parsedEvent.Type,
			Subtype: parsedEvent.Subtype,
			Speaker: parsedEvent.Speaker,
			Short:   parsedEvent.Short,
			Long:    parsedEvent.Long,
			Start:   baseDate.Add(time.Duration(startTime.Hour())*time.Hour + time.Duration(startTime.Minute())*time.Minute),
			Finish:  baseDate.Add(time.Duration(finishTime.Hour())*time.Hour + time.Duration(finishTime.Minute())*time.Minute),
		}
		storage.AddEvent(event)
	}

	return nil
}
