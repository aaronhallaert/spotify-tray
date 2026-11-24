package tray

import (
	"testing"
	"time"
)

func TestDeliverClick(t *testing.T) {
	const id = 100
	item := &MenuItem{ClickedCh: make(chan struct{}, 1)}
	itemsMu.Lock()
	items[id] = item
	itemsMu.Unlock()
	t.Cleanup(func() {
		itemsMu.Lock()
		delete(items, id)
		itemsMu.Unlock()
	})

	deliverClick(id)
	select {
	case <-item.ClickedCh:
	default:
		t.Fatal("click was not delivered to its menu item")
	}

	item.ClickedCh <- struct{}{}
	returned := make(chan struct{})
	go func() {
		deliverClick(id)
		close(returned)
	}()
	select {
	case <-returned:
	case <-time.After(time.Second):
		t.Fatal("click delivery blocked the UI callback")
	}
}

func TestReadyDoesNotBlock(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	finished := make(chan struct{})
	onReady = func() {
		close(started)
		<-release
	}
	t.Cleanup(func() { onReady = nil })

	go func() {
		tray_ready()
		close(finished)
	}()
	select {
	case <-started:
	case <-time.After(time.Second):
		close(release)
		t.Fatal("ready callback did not start")
	}
	select {
	case <-finished:
	case <-time.After(time.Second):
		close(release)
		t.Fatal("ready callback blocked the UI")
	}
	close(release)
}
