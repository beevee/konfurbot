package konfurbot

import "time"

// ScheduleStorage provides searching and filtering capabilities over schedule
type ScheduleStorage interface {
	AddEvent(Event)
	GetEventsByType(string) []Event
	GetEventsByTypeAndSubtype(string, string) []Event
	GetCurrentEventsByType(string, time.Time) []Event
	GetNextEventsByType(string, time.Time, time.Duration) []Event
}

// Event is a single event in conference
type Event struct {
	Type    string
	Subtype string
	Speaker string
	Venue   string
	Short   string
	Long    string
	Start   time.Time
	Finish  time.Time
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

// GetEventsByTypeAndSubtype returns unfiltered list of events by type
func (s *Schedule) GetEventsByTypeAndSubtype(kind, subkind string) []Event {
	events := make([]Event, 0, len(s.events[kind]))
	for _, event := range s.events[kind] {
		if event.Subtype == subkind {
			events = append(events, event)
		}
	}
	return events
}

// GetCurrentEventsByType returns list of events by type, and only events that have started and not have finished yet
func (s *Schedule) GetCurrentEventsByType(kind string, now time.Time) []Event {
	events := make([]Event, 0, len(s.events[kind]))
	for _, event := range s.events[kind] {
		if event.Start.Before(now) && event.Finish.After(now) {
			events = append(events, event)
		}
	}
	return events
}

// GetNextEventsByType returns list of events by type, and only events that will start in the next interval
func (s *Schedule) GetNextEventsByType(kind string, now time.Time, interval time.Duration) []Event {
	events := make([]Event, 0, len(s.events[kind]))
	for _, event := range s.events[kind] {
		if event.Start.After(now) && event.Start.Before(now.Add(interval)) {
			events = append(events, event)
		}
	}
	return events
}
