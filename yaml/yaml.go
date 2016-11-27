package yaml

import (
	"time"

	"gopkg.in/yaml.v2"

	"github.com/beevee/konfurbot"
)

// Event is a YAML-parseable adapter for konfurbot.Event
type Event struct {
	Type   string `yaml:"type"`
	Short  string `yaml:"short"`
	Long   string `yaml:"long"`
	Start  string `yaml:"start"`
	Finish string `yaml:"finish"`
}

// FillScheduleStorage fills schedule storage with events parsed from YAML file
func FillScheduleStorage(storage konfurbot.ScheduleStorage, file []byte) error {
	var parsedEvents []Event
	if err := yaml.Unmarshal(file, &parsedEvents); err != nil {
		return err
	}

	for _, parsedEvent := range parsedEvents {
		start, err := time.Parse("15:04", parsedEvent.Start)
		if err != nil {
			return err
		}

		finish, err := time.Parse("15:04", parsedEvent.Finish)
		if err != nil {
			return err
		}

		event := konfurbot.Event{
			Type:   parsedEvent.Type,
			Short:  parsedEvent.Short,
			Long:   parsedEvent.Long,
			Start:  start,
			Finish: finish,
		}
		storage.AddEvent(event)
	}

	return nil
}
