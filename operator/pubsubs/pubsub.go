package pubsubs

type Subscriber interface {
	Receive()
	Publish(msg interface{})
}
