package hub

import (
	"github.com/google/uuid"
	"github.com/krls256/dsd2024additional/pkg/entities"
	"github.com/samber/mo"
	"sync"
)

type Subscriber func()

func NewHub() *Hub {
	return &Hub{
		clients:               map[uuid.UUID]entities.AsyncWriter[entities.LeftWithAction]{},
		unregisterSubscribers: map[uuid.UUID][]Subscriber{},
	}
}

type Hub struct {
	clients               map[uuid.UUID]entities.AsyncWriter[entities.LeftWithAction]
	unregisterSubscribers map[uuid.UUID][]Subscriber

	clientsMu               sync.RWMutex
	unregisterSubscribersMu sync.RWMutex
}

func (h *Hub) Register(id uuid.UUID, writer entities.AsyncWriter[entities.LeftWithAction]) {
	h.clientsMu.Lock()
	defer h.clientsMu.Unlock()

	h.clients[id] = writer
}

func (h *Hub) Unregister(id uuid.UUID) {
	h.clientsMu.Lock()
	defer h.clientsMu.Unlock()

	delete(h.clients, id)

	h.unregisterSubscribersMu.RLock()
	subs, ok := h.unregisterSubscribers[id]
	h.unregisterSubscribersMu.RUnlock()

	if ok {
		h.unregisterSubscribersMu.Lock()
		delete(h.unregisterSubscribers, id)
		h.unregisterSubscribersMu.Unlock()

		for _, s := range subs {
			s()
		}
	}
}

func (h *Hub) SubscribeOnUnregister(id uuid.UUID, fn func()) {
	h.clientsMu.RLock()
	_, ok := h.clients[id]
	h.clientsMu.RUnlock()

	if !ok {
		fn()

		return
	}

	h.unregisterSubscribersMu.Lock()
	h.unregisterSubscribers[id] = append(h.unregisterSubscribers[id], fn)
	h.unregisterSubscribersMu.Unlock()
}

func (h *Hub) Broadcast(action string, payload any) {
	h.clientsMu.RLock()
	defer h.clientsMu.RUnlock()

	for _, c := range h.clients {
		c.AsyncWrite(mo.Left[entities.LeftWithAction, error](entities.NewLeftWithAction(action, payload)))
	}
}
