package sync_notification

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestSyncNotification(t *testing.T) {
	sn := CreateSyncNotification()

	ch := make(chan string)
	defer close(ch)
	wg := sync.WaitGroup{}

	f := func(i uint32, sn ISyncNotification) {
		waiting, err := sn.GetWaiting()
		if err != nil {
			ch <- "err"
			return
		}
		defer waiting.Done()
		wg.Done()

		<-waiting.GetSignalChan()
		ch <- fmt.Sprint(i)

	}

	t.Log("sn.Signal() check ---------------------------")
	wg.Add(2)
	go f(1, sn)
	go f(2, sn)
	wg.Wait()

	sn.Signal(1)

	select {
	case <-time.After(3 * time.Second):
		t.Fatal("a value was expected in the channel, but it was not received")
	case str := <-ch:
		if str == "err" {
			t.Fatal("Received error from SyncNotification")
		}
		t.Log("output:", str)
	}

	select {
	case <-time.After(3 * time.Second):
	case str := <-ch:
		t.Fatal("an unexpected message was received in the channel", str)
	}

	t.Log("sn.Broadcast() check ---------------------------")
	wg.Add(3)
	go f(3, sn)
	go f(4, sn)
	go f(5, sn)
	wg.Wait()

	sn.Broadcast()

	for range 4 {
		select {
		case <-time.After(3 * time.Second):
			t.Fatal("a value was expected in the channel, but it was not received")
		case str := <-ch:
			if str == "err" {
				t.Fatal("received error from SyncNotification")
			}
			t.Log("output:", str)
		}
	}

	t.Log("sn.Done(id) check ---------------------------")
	waiting, err := sn.GetWaiting()
	if err != nil {
		t.Fatal("received error from SyncNotification")
	}
	id := waiting.GetId()
	sn.Done(id)
	_, ok := <-waiting.GetSignalChan()
	if ok {
		t.Fatal("done function does not work")
	}
}
