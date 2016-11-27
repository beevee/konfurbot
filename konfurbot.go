package konfurbot

import "time"

// Event is a single event in conference
type Event struct {
	Type   string
	Short  string
	Long   string
	Start  time.Time
	Finish time.Time
}

// Schedule is a collection of events in a conference
type Schedule struct {
	Events map[string][]Event
}
