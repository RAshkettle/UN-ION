package main

// EventType represents different types of game events
type EventType int

const (
	EventBlocksRemoved EventType = iota
	EventPiecePlaced
	EventHardDrop
	EventLineCleared
	EventGameOver
)

// GameEvent represents a game event with associated data
type GameEvent struct {
	Type EventType
	Data interface{}
}

// BlocksRemovedData contains data for blocks removed event
type BlocksRemovedData struct {
	Count     int
	Positions []Position
}

// PiecePlacedData contains data for piece placed event
type PiecePlacedData struct {
	Position Position
}

// HardDropData contains data for hard drop event
type HardDropData struct {
	DropHeight int
}

// Position represents a world position
type Position struct {
	X, Y float64
}

// EventSystem manages game events
type EventSystem struct {
	listeners map[EventType][]func(GameEvent)
}

// NewEventSystem creates a new event system
func NewEventSystem() *EventSystem {
	return &EventSystem{
		listeners: make(map[EventType][]func(GameEvent)),
	}
}

// Subscribe adds a listener for a specific event type
func (es *EventSystem) Subscribe(eventType EventType, callback func(GameEvent)) {
	es.listeners[eventType] = append(es.listeners[eventType], callback)
}

// Emit sends an event to all listeners
func (es *EventSystem) Emit(event GameEvent) {
	if callbacks, exists := es.listeners[event.Type]; exists {
		for _, callback := range callbacks {
			callback(event)
		}
	}
}
