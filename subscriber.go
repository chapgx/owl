package owl

// Subscriber is how you read data from the watcher
type Subscriber struct {
	channel chan any
	flag    int
}

// Listen returns the [Subscriber] channel
func (s *Subscriber) Listen() chan any {
	return s.channel
}
