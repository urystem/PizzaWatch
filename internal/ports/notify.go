package ports

type NotifyRabbit interface {
	GiveChannel() <-chan []byte
	CloseRabbit() error
}

type NotifyService interface {
	StartNotify()
}
