package kafka

type MessageValueProcessor interface {
	ProcessMessageValue(msg []byte)
}
