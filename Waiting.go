package sync_notification

type IWaiting interface {
	Done()
	GetSignalChan() (signal <-chan struct{})
	GetId() uint32
}

func CreateWaiting(id uint32, signal <-chan struct{}, sn ISyncNotification) IWaiting {
	w := Waiting{id: id, signal: signal, sn: sn}
	return w
}

type Waiting struct {
	id     uint32
	signal <-chan struct{}
	sn     ISyncNotification
}

func (w Waiting) Done() {
	w.sn.Done(w.id)
}

func (w Waiting) GetSignalChan() (signal <-chan struct{}) {
	return w.signal
}

func (w Waiting) GetId() uint32 {
	return w.id
}
