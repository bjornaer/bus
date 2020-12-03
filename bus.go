package bus

import "sync"

// DataEvent represents an event of varied Data
// and to which Topic it is addressed
type DataEvent struct {
	Data  interface{}
	Topic string
}

// DataChannel is a channel which can accept a DataEvent
type DataChannel chan DataEvent

// DataChannelSlice is a set of DataChannels
type DataChannelSlice []DataChannel

// EventBus stores the information about subscribers interested for
// a particular topic
type EventBus struct {
	subscribers map[string]DataChannelSlice
	rm          sync.RWMutex
}

// Subscribe to a Topic
func (eb *EventBus) Subscribe(topic string, ch DataChannel) {
	eb.rm.Lock()
	if prev, found := eb.subscribers[topic]; found {
		eb.subscribers[topic] = append(prev, ch)
	} else {
		eb.subscribers[topic] = append([]DataChannel{}, ch)
	}
	eb.rm.Unlock()
}

// Publish data to a topic
func (eb *EventBus) Publish(topic string, data interface{}) {
	eb.rm.RLock()
	if chans, found := eb.subscribers[topic]; found {
		// this is done because the slices refer to same array even though they are passed by value
		// thus we are creating a new slice with our elements thus preserve locking correctly.
		channels := append(DataChannelSlice{}, chans...)
		go func(data DataEvent, dataChannelSlices DataChannelSlice) {
			for _, ch := range dataChannelSlices {
				ch <- data
			}
		}(DataEvent{Data: data, Topic: topic}, channels)
	}
	eb.rm.RUnlock()
}
