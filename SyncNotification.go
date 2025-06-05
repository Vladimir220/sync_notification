package sync_notification

import (
	"fmt"
	"sync"
)

type ISyncNotification interface {
	GetWaiting() (waiting IWaiting, err error)
	Signal()
	Broadcast()
	Done(id uint32)
}

func CreateSyncNotification() ISyncNotification {
	return &syncNotification{
		ch:      make(map[uint32]chan struct{}),
		freeIds: make([]uint32, 0),
		mu:      &sync.RWMutex{},
	}
}

type syncNotification struct {
	ch      map[uint32]chan struct{}
	nextId  uint32
	freeIds []uint32
	mu      *sync.RWMutex
}

func (c *syncNotification) getNextId() (uint32, error) {
	maxId := ^uint32(0)

	if len(c.freeIds) != 0 {
		buf := c.freeIds[0]
		c.freeIds = c.freeIds[1:]
		return buf, nil
	}

	if c.nextId == maxId {
		return c.nextId, fmt.Errorf("the maximum number of goroutines has already been created: %d", c.nextId)
	}

	buf := c.nextId
	c.nextId++
	return buf, nil
}

func (sn *syncNotification) GetWaiting() (waiting IWaiting, err error) {
	ch := make(chan struct{})

	sn.mu.Lock()
	defer sn.mu.Unlock()

	id, err := sn.getNextId()
	if err != nil {
		return Waiting{}, err
	}

	waiting = CreateWaiting(id, ch, sn)
	sn.ch[id] = ch

	return
}

func (sn *syncNotification) Done(id uint32) {
	sn.mu.Lock()
	defer sn.mu.Unlock()

	close(sn.ch[id])
	delete(sn.ch, id)
	sn.freeIds = append(sn.freeIds, id)
}

func (sn *syncNotification) Signal() {
	sn.mu.RLock()
	defer sn.mu.RUnlock()

	for _, ch := range sn.ch {
		ch <- struct{}{}
		break
	}
}

func (sn *syncNotification) Broadcast() {
	sn.mu.RLock()
	defer sn.mu.RUnlock()

	for _, ch := range sn.ch {
		ch <- struct{}{}
	}
}
