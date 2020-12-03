package bus

import (
	"testing"
)

var eb = &EventBus{
	subscribers: map[string]DataChannelSlice{},
}

func TestPubSub(t *testing.T) {
	ch := make(chan DataEvent)
	eb.Subscribe("topic", ch)
	eb.Publish("topic", "this is a test")
	d := <-ch
	if d.Data != "this is a test" {
		t.Errorf("Unexpected data %s for topic %s", d.Data, d.Topic)
	}

}

func TestMultiPubSub(t *testing.T) {
	ch := make(chan DataEvent)
	ch2 := make(chan DataEvent)
	ch3 := make(chan DataEvent)

	eb.Subscribe("topic1", ch)
	eb.Subscribe("topic2", ch2)
	eb.Subscribe("topic1", ch3)

	eb.Publish("topic1", "data1")
	eb.Publish("topic2", "data2")
	d := <-ch
	d2 := <-ch2
	d3 := <-ch3
	if d.Data != "data1" {
		t.Errorf("Unexpected data from d1 for topic %s", d.Topic)
	}
	if d2.Data != "data2" {
		t.Errorf("Unexpected data from d2 for topic %s", d2.Topic)
	}
	if d3.Data != d.Data {
		t.Errorf("Unexpected data from d3 for topic %s", d3.Topic)
	}

}
