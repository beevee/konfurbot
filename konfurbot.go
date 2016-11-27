package konfurbot

import "time"

// ScheduleStorage provides searching and filtering capabilities over schedule
type ScheduleStorage interface {
	AddEvent(Event)
	GetEventsByType(string) []Event
}

// Event is a single event in conference
type Event struct {
	Type   string
	Short  string
	Long   string
	Start  time.Time
	Finish time.Time
}

// Schedule is an implementation of ScheduleStorage
type Schedule struct {
	events map[string][]Event
}

// AddEvent adds event to storage, preserving order of events
func (s *Schedule) AddEvent(event Event) {
	if s.events == nil {
		s.events = make(map[string][]Event)
	}
	s.events[event.Type] = append(s.events[event.Type], event)
}

// GetEventsByType returns unfiltered list of events by type
func (s *Schedule) GetEventsByType(kind string) []Event {
	return s.events[kind]
}
