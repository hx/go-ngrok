package ngrok

// Subscription is the handle to a Dispatcher subscription that allows it to be unsubscribed.
type Subscription struct {
	subscriber Subscriber
	onRelease  func()
}

// Unsubscribe terminates the subscription. If a Subscriber is running, it will be allowed to finish. A Subscription may
// call Unsubscribe on itself.
func (s *Subscription) Unsubscribe() {
	s.onRelease()
}
