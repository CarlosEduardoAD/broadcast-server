package pub

import (
	"slices"

	"github.com/CarlosEduardoAD/broadcast-server/internal/sub"
)

type Pub struct {
	Subscribers []sub.Subscriber
}

func NewPublisher(subs []sub.Subscriber) *Pub {
	return &Pub{
		Subscribers: subs,
	}
}

func (p *Pub) Publish(message []byte, mt int) {
	for _, sub := range p.Subscribers {
		sub.Call(mt, message, sub.Conn)
	}
}

func (p *Pub) Subscribe(subscriber_to_add sub.Subscriber) []sub.Subscriber {
	p.Subscribers = append(p.Subscribers, subscriber_to_add)

	return p.Subscribers
}

func (p *Pub) find_subscriber(sub_to_find sub.Subscriber) int {
	for k := range p.Subscribers {
		if p.Subscribers[k] == sub_to_find {
			return int(slices.Index(p.Subscribers, p.Subscribers[k]))
		}
	}

	return -1
}

func (p *Pub) Remove(subscriber_to_remove sub.Subscriber) {
	sub_index := p.find_subscriber(subscriber_to_remove)

	if sub_index < 0 {
		panic("sub not found")
	}

	p.Subscribers = slices.Delete(p.Subscribers, sub_index, sub_index+1)
}
