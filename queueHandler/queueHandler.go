package queueHandler

type QueueHandler interface {
	Publish(ev any)
}
