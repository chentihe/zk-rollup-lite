package dbcache

type Subscriber interface {
	Receive()
	Publish(msg interface{})
}
