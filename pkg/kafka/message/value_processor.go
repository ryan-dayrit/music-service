package message

type MessageValueProcessor interface {
	Process(msg []byte)
}
