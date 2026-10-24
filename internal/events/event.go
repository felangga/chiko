package events

type (
	IEventHandler func(IEvent)
	EventType     string
	IEvent        interface {
		Type() EventType
	}
)
