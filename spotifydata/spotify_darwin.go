package spotifydata

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework AppKit
#include <stdlib.h>
#include "spotify_darwin.h"
*/
import "C"

import (
	"errors"
	"sync"
	"unsafe"
)

var spotifyMu sync.Mutex

func readSpotify(progress, album bool, cachedID string, cachedPID int) ([]byte, error) {
	spotifyMu.Lock()
	defer spotifyMu.Unlock()

	id := C.CString(cachedID)
	defer C.free(unsafe.Pointer(id))
	var message *C.char
	result := C.spotify_read(C.bool(progress), C.bool(album), id, C.int(cachedPID), &message)
	defer C.free(unsafe.Pointer(result))
	if message != nil {
		defer C.free(unsafe.Pointer(message))
		return nil, errors.New(C.GoString(message))
	}
	if result == nil {
		return nil, errors.New("native Spotify read returned no data")
	}
	return []byte(C.GoString(result)), nil
}

func toggleSpotify() error {
	spotifyMu.Lock()
	defer spotifyMu.Unlock()

	if message := C.spotify_toggle(); message != nil {
		defer C.free(unsafe.Pointer(message))
		return errors.New(C.GoString(message))
	}
	return nil
}
