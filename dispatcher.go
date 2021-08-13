package ngrok

import (
	"context"
	"sync"
)

// Dispatcher is used to manage subscriptions to a process's LogMessage stream.
type Dispatcher struct {
	subscriptions []*Subscription
	mutex         sync.RWMutex
	finished      bool
}

// Subscribe adds a Subscriber, and returns a Subscription, which can be used to unsubscribe.
func (d *Dispatcher) Subscribe(subscriber Subscriber) (subscription *Subscription) {
	return d.subscribe(subscriber, nil)
}

// WaitOne blocks until a LogMessage is received for which selector returns true, or the underlying stream is closed.
func (d *Dispatcher) WaitOne(selector func(message *LogMessage) bool) (err error) {
	return d.WaitOneContext(context.Background(), selector)
}

// WaitOneContext blocks until a LogMessage is received for which selector returns true, or the underlying stream is
// closed, or ctx expires.
func (d *Dispatcher) WaitOneContext(ctx context.Context, selector func(message *LogMessage) bool) (err error) {
	var (
		localCtx, cancel = context.WithCancel(ctx)
		interrupt        = make(chan struct{})
		result           = make(chan bool)
	)

	subscription := d.subscribe(func(message *LogMessage) {
		if _, ok := <-interrupt; ok {
			result <- selector(message)
		}
	}, cancel)

	for done := false; !done; {
		select {
		case <-localCtx.Done():
			done = true
			err = localCtx.Err()
		case interrupt <- struct{}{}:
			done = <-result
		}
	}

	subscription.Unsubscribe()
	close(interrupt)
	close(result)

	return
}

func (d *Dispatcher) subscribe(subscriber Subscriber, onRelease func()) (subscription *Subscription) {
	subscription = &Subscription{subscriber: subscriber}
	subscription.onRelease = func() {
		d.unsubscribe(subscription)
		if onRelease != nil {
			onRelease()
		}
	}
	d.mutex.Lock()
	if d.finished {
		panic("cannot subscribe to released dispatcher")
	}
	d.subscriptions = append(d.subscriptions, subscription)
	d.mutex.Unlock()
	return
}

func (d *Dispatcher) dispatch(message *LogMessage) {
	d.mutex.RLock()
	subscriptions := d.subscriptions
	d.mutex.RUnlock()
	for _, subscription := range subscriptions {
		subscription.subscriber(message)
	}
}

func (d *Dispatcher) unsubscribe(subscription *Subscription) {
	d.mutex.Lock()
	for i, s := range d.subscriptions {
		if s == subscription {
			d.subscriptions = append(d.subscriptions[:i], d.subscriptions[i+1:]...)
		}
	}
	d.mutex.Unlock()
}

func (d *Dispatcher) release() {
	d.mutex.Lock()
	if d.finished {
		panic("dispatcher already released")
	}
	d.finished = true
	subscriptions := d.subscriptions
	d.mutex.Unlock()
	for _, subscription := range subscriptions {
		subscription.Unsubscribe()
	}
}
