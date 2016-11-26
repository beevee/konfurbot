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

// ParseSchedule returns schedule struct parsed from YAML file
func ParseSchedule(file []byte) (*konfurbot.Schedule, error) {
	var parsedEvents []Event
	if err := yaml.Unmarshal(file, &parsedEvents); err != nil {
		return nil, err
	}

	schedule := &konfurbot.Schedule{
		Events: make(map[string][]konfurbot.Event),
	}
	for _, parsedEvent := range parsedEvents {
		start, err := time.Parse("15:04", parsedEvent.Start)
		if err != nil {
			return nil, err
		}

		finish, err := time.Parse("15:04", parsedEvent.Finish)
		if err != nil {
			return nil, err
		}

		event := konfurbot.Event{
			Type:   parsedEvent.Type,
			Short:  parsedEvent.Short,
			Long:   parsedEvent.Long,
			Start:  start,
			Finish: finish,
		}
		schedule.Events[parsedEvent.Type] = append(schedule.Events[parsedEvent.Type], event)
	}

	return schedule, nil
}
