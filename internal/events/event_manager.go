package events

import (
	"context"
	"os"
	"os/signal"
)

var instance *EventManager = nil

type EventManager struct {
	queue     chan IEvent
	listeners map[EventType][]IEventHandler
}

// the callback is run inside it's own goroutine
func (e *EventManager) Listen(to EventType, callback IEventHandler) {
	e.listeners[to] = append(e.listeners[to], callback)
}

func (e *EventManager) Emit(event IEvent) {
	e.queue <- event
}

func GetEventManager() *EventManager {
	return instance
}

func init() {
	instance = new(EventManager)
	instance.queue = make(chan IEvent, 10)
	instance.listeners = make(map[EventType][]IEventHandler)
	go func() {
		ctx, cancel := signal.NotifyContext(context.Background(), os.Kill)
		defer cancel()
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-instance.queue:
				for _, cb := range instance.listeners[event.Type()] {
					// NOTE:  this should be a bounded fan-out -- this is temp
					go cb(event)
				}
			}
		}
	}()
}
