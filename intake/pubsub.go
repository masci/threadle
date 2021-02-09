package intake

import "sync"

// PubSub is a ridicoulously simple message broker
type PubSub struct {
	sync.RWMutex

	subscribers map[string][]chan []byte
	closed      bool
}

// NewPubsub creates an instance of the broker
func NewPubsub() *PubSub {
	ps := &PubSub{}
	ps.subscribers = make(map[string][]chan []byte)
	return ps
}

// Subscribe lets clients subscribe to a topic
func (ps *PubSub) Subscribe(topic string) <-chan []byte {
	ps.Lock()
	defer ps.Unlock()

	ch := make(chan []byte, 1)
	ps.subscribers[topic] = append(ps.subscribers[topic], ch)
	return ch
}

// Publish sends a message to a topic
func (ps *PubSub) Publish(topic string, msg []byte) {
	ps.RLock()
	defer ps.RUnlock()

	if ps.closed {
		return
	}

	for _, ch := range ps.subscribers[topic] {
		ch <- msg
	}
}

// Close shuts down the broker and cleans up the channels
func (ps *PubSub) Close() {
	ps.Lock()
	defer ps.Unlock()

	if ps.closed {
		return
	}

	ps.closed = true
	for _, subs := range ps.subscribers {
		for _, ch := range subs {
			close(ch)
		}
	}
}
