package events

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var instance *EventManager = nil

type EventManager struct {
	queue     chan IEvent
	mu        sync.RWMutex
	listeners map[EventType][]IEventHandler
}

// the callback is run inside it's own goroutine
func (e *EventManager) Listen(to EventType, callback IEventHandler) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.listeners[to] = append(e.listeners[to], callback)
}

func (e *EventManager) Emit(event IEvent) {
	e.queue <- event
}

func (e *EventManager) listenersFor(eventType EventType) []IEventHandler {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return append([]IEventHandler(nil), e.listeners[eventType]...)
}

func GetEventManager() *EventManager {
	return instance
}

func init() {
	instance = new(EventManager)
	instance.queue = make(chan IEvent, 10)
	instance.listeners = make(map[EventType][]IEventHandler)
	go func() {
		e := instance
		// SIGTERM (and Ctrl+C) are catchable; SIGKILL is not, so os.Kill would
		// make ctx.Done() dead code. Use SIGTERM for graceful shutdown.
		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
		defer cancel()
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-e.queue:
				for _, cb := range e.listenersFor(event.Type()) {
					// NOTE:  this should be a bounded fan-out -- this is temp
					go cb(event)
				}
			}
		}
	}()
}
