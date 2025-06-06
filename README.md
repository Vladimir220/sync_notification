# sync_notification #

### En ###
A synchronization tool that will notify goroutines about an event that has occurred.

Similar to sync.Cond, but there is no mutex handling and no Wait method. Instead of the Wait method, we have a waiting channel of the IWaiting interface.

### Rus ###
Инструмент для синхронизации, который будет уведомлять горутины о произошедшем событии. 

Похоже на sync.Cond, но нет работы с мьютексом и нет метода Wait. Вместо метода Wait у нас канал ожидания интерфейса IWaiting.

### Interfaces ###

```
type ISyncNotification interface {
	GetWaiting() (waiting IWaiting, err error)
	Signal()
	Broadcast()
	Done(id uint32)
}
```
```
type IWaiting interface {
	Done()
	GetSignalChan() (signal <-chan struct{})
	GetId() uint32
}
```

### Constructors ###
```
func CreateSyncNotification() ISyncNotification
```
```
func CreateWaiting(id uint32, signal <-chan struct{}, sn ISyncNotification) IWaiting
```