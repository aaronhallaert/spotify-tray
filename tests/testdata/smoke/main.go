package main

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework AppKit
int smoke_check(int checked);
int smoke_hidden(void);
void smoke_click(void);
*/
import "C"

import (
	"fmt"
	"os"
	"spotify-tray/internal/tray"
	"time"
)

func main() {
	result := make(chan error, 1)
	tray.Run(func() {
		tray.SetTitle("♫ demo")
		setting := tray.AddMenuItemCheckbox("Setting", "hint", true)
		tray.AddSeparator()
		playback := tray.AddMenuItem("Play", "playback")
		playback.Hide()
		if C.smoke_hidden() == 0 {
			result <- fmt.Errorf("playback item was not hidden")
			tray.Quit()
			return
		}
		playback.SetTitle("Pause")
		playback.Show()
		if code := C.smoke_check(1); code != 0 {
			result <- fmt.Errorf("initial AppKit state: %d", code)
			tray.Quit()
			return
		}

		clicked := make(chan struct{})
		go func() {
			<-setting.ClickedCh
			setting.Uncheck()
			close(clicked)
		}()
		deadline := time.After(3 * time.Second)
		ticker := time.NewTicker(20 * time.Millisecond)
		defer ticker.Stop()
		for {
			C.smoke_click()
			select {
			case <-clicked:
				if code := C.smoke_check(0); code != 0 {
					result <- fmt.Errorf("updated AppKit state: %d", code)
				} else {
					result <- nil
				}
				tray.Quit()
				return
			case <-deadline:
				result <- fmt.Errorf("menu click was not delivered")
				tray.Quit()
				return
			case <-ticker.C:
			}
		}
	}, nil)
	if err := <-result; err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("native tray smoke test passed")
}
