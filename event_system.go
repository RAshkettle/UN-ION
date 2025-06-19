package main

type EventType int

const (
	EventBlocksRemoved EventType = iota
	EventPiecePlaced
	EventHardDrop
	EventLineCleared
	EventGameOver
)

type GameEvent struct {
	Type EventType
	Data interface{}
}

type BlocksRemovedData struct {
	Count     int
	Positions []Position
}

type PiecePlacedData struct {
	Position Position
}

type HardDropData struct {
	DropHeight int
}

type Position struct {
	X, Y float64
}

type EventSystem struct {
	listeners map[EventType][]func(GameEvent)
}

func NewEventSystem() *EventSystem {
	return &EventSystem{
		listeners: make(map[EventType][]func(GameEvent)),
	}
}

func (es *EventSystem) Subscribe(eventType EventType, callback func(GameEvent)) {
	es.listeners[eventType] = append(es.listeners[eventType], callback)
}

func (es *EventSystem) Emit(event GameEvent) {
	if callbacks, exists := es.listeners[event.Type]; exists {
		for _, callback := range callbacks {
			callback(event)
		}
	}
}
