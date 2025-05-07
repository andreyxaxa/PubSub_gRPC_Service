package subpub

import "context"

// Callback function that processes messages delivered to subscribers.
type MessageHandler func(msg interface{})

type Subscription interface {
	// Will remove interest in the current subject subscription is for.
	Unsubscribe()
}

type SubPub interface {
	// Creates an async queue subscriber on the given subject.
	Subscribe(subject string, cb MessageHandler) (Subscription, error)

	// Publishies the msg argument to the given subject.
	Publish(subject string, msg interface{}) error

	// Shutdown sub-pub system.
	// May be blocked by data delivery until the context is canceled.
	Close(ctx context.Context) error
}

func NewSubPub() SubPub {
	return &subPubImpl{
		subscribers: make(map[string]map[*subscriber]struct{}),
		closeCh:     make(chan struct{}),
	}
}
