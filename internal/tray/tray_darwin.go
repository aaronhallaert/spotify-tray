package tray

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework AppKit
#include <stdlib.h>
#include "tray_darwin.h"
*/
import "C"

import (
	"log"
	"runtime"
	"sync"
	"unsafe"
)

var (
	onReady func()
	itemsMu sync.RWMutex
	items   = make(map[uint32]*MenuItem)
	nextID  uint32
)

func init() {
	runtime.LockOSThread()
}

type MenuItem struct {
	ClickedCh chan struct{}
	id        uint32
	checked   bool
}

func Run(ready, exit func()) {
	onReady = ready
	C.tray_run()
	if exit != nil {
		exit()
	}
}

func Quit() {
	C.tray_quit()
}

func SetTitle(title string) {
	value := C.CString(title)
	defer C.free(unsafe.Pointer(value))
	C.tray_set_title(value)
}

func AddMenuItem(title, tooltip string) *MenuItem {
	return addMenuItem(title, tooltip, false)
}

func AddMenuItemCheckbox(title, tooltip string, checked bool) *MenuItem {
	return addMenuItem(title, tooltip, checked)
}

func addMenuItem(title, tooltip string, checked bool) *MenuItem {
	itemsMu.Lock()
	nextID++
	item := &MenuItem{ClickedCh: make(chan struct{}), id: nextID, checked: checked}
	items[item.id] = item
	itemsMu.Unlock()

	cTitle, cTooltip := C.CString(title), C.CString(tooltip)
	defer C.free(unsafe.Pointer(cTitle))
	defer C.free(unsafe.Pointer(cTooltip))
	C.tray_add_item(C.uint32_t(item.id), cTitle, cTooltip, C.bool(checked))
	return item
}

func AddSeparator() {
	C.tray_add_separator()
}

func (item *MenuItem) SetTitle(title string) {
	value := C.CString(title)
	defer C.free(unsafe.Pointer(value))
	C.tray_set_item_title(C.uint32_t(item.id), value)
}

func (item *MenuItem) Hide() {
	C.tray_set_item_hidden(C.uint32_t(item.id), true)
}

func (item *MenuItem) Show() {
	C.tray_set_item_hidden(C.uint32_t(item.id), false)
}

func (item *MenuItem) Checked() bool {
	return item.checked
}

func (item *MenuItem) Check() {
	item.checked = true
	C.tray_set_item_checked(C.uint32_t(item.id), true)
}

func (item *MenuItem) Uncheck() {
	item.checked = false
	C.tray_set_item_checked(C.uint32_t(item.id), false)
}

//export tray_ready
func tray_ready() {
	if onReady != nil {
		go onReady()
	}
}

//export tray_item_clicked
func tray_item_clicked(id C.uint32_t) {
	deliverClick(uint32(id))
}

func deliverClick(id uint32) {
	itemsMu.RLock()
	item := items[id]
	itemsMu.RUnlock()
	if item == nil {
		log.Printf("tray: unknown menu item %d", id)
		return
	}
	select {
	case item.ClickedCh <- struct{}{}:
	default:
	}
}
