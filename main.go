package main

import (
	"os"
	"spotify-tray/spotifydata"
	"spotify-tray/storage"
	"sync"
	"syscall"
	"time"

	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
)

var wg sync.WaitGroup

func main() {
	selfRestartTimer := time.NewTimer(time.Hour * 4)
	wg.Add(1)

	// Restart every 4 hours to avoid memory build up
	go func() {
		defer wg.Done()
		<-selfRestartTimer.C
		RestartSelf()
	}()

	onExit := func() {
	}

	storage.Init()
	systray.Run(onReady, onExit)

	wg.Wait()
}

func onReady() {
	systray.SetTitle("Loading...")
	mLyrics := systray.AddMenuItem("Lyrics", "Search for lyrics online")
	systray.AddSeparator()
	mArtistFirst := systray.AddMenuItemCheckbox("Show artist first?", "Show artist first", storage.GetArtistFirst())
	mShowAlbum := systray.AddMenuItemCheckbox("Show album?", "Show Album", storage.GetShowAlbum())
	mProgress := systray.AddMenuItemCheckbox("Show progress?", "Show Progress", storage.GetShowProgress())
	systray.AddSeparator()
	mMoreSpace := systray.AddMenuItemCheckbox("Use more space?", "Use more space", storage.GetMoreSpace())
	mAlternateSeparator := systray.AddMenuItemCheckbox("Use alternate separator?", "Show Alternate separator", storage.GetAlternateSeparator())
	mOpenAtLogin := systray.AddMenuItemCheckbox("Open at login?", "Open at login", storage.GetOpenAtLogin())
	systray.AddSeparator()
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")

	currentSpotifyData := &spotifydata.Data{}
	if spotifydata.IsSpotifyRunning() {
		currentSpotifyData = spotifydata.GetData(storage.GetShowProgress(), storage.GetShowAlbum())
	}
	updateTray(currentSpotifyData)

	go func() {
		for {
			select {
			case <-mLyrics.ClickedCh:
				open.Run("https://www.google.com/search?q=" + currentSpotifyData.Track + " - " + currentSpotifyData.Artist + " lyrics")
			case <-mArtistFirst.ClickedCh:
				if mArtistFirst.Checked() {
					mArtistFirst.Uncheck()
					storage.SetArtistFirst(false)
				} else {
					mArtistFirst.Check()
					storage.SetArtistFirst(true)
				}
			case <-mShowAlbum.ClickedCh:
				if mShowAlbum.Checked() {
					mShowAlbum.Uncheck()
					storage.SetShowAlbum(false)
				} else {
					mShowAlbum.Check()
					storage.SetShowAlbum(true)
				}
			case <-mProgress.ClickedCh:
				if mProgress.Checked() {
					mProgress.Uncheck()
					storage.SetShowProgress(false)
				} else {
					mProgress.Check()
					storage.SetShowProgress(true)
				}
			case <-mMoreSpace.ClickedCh:
				if mMoreSpace.Checked() {
					mMoreSpace.Uncheck()
					storage.SetMoreSpace(false)
				} else {
					mMoreSpace.Check()
					storage.SetMoreSpace(true)
				}
			case <-mAlternateSeparator.ClickedCh:
				if mAlternateSeparator.Checked() {
					mAlternateSeparator.Uncheck()
					storage.SetAlternateSeparator(false)
				} else {
					mAlternateSeparator.Check()
					storage.SetAlternateSeparator(true)
				}
			case <-mOpenAtLogin.ClickedCh:
				if mOpenAtLogin.Checked() {
					mOpenAtLogin.Uncheck()
					storage.SetOpenAtLogin(false)
				} else {
					mOpenAtLogin.Check()
					storage.SetOpenAtLogin(true)
				}
			}
		}
	}()

	go func() {
		<-mQuitOrig.ClickedCh
		systray.Quit()
	}()

	go func() {
		for {
			if spotifydata.IsSpotifyRunning() {
				currentSpotifyData = spotifydata.GetData(storage.GetShowProgress(), storage.GetShowAlbum())
				updateTray(currentSpotifyData)
			} else {
				currentSpotifyData.Status = ""
			}
			time.Sleep(time.Second)
		}
	}()
}

func updateTray(d *spotifydata.Data) {
	systray.SetTitle(d.Format(
		storage.GetShowProgress(),
		storage.GetShowAlbum(),
		storage.GetArtistFirst(),
		storage.GetMoreSpace(),
		storage.GetAlternateSeparator(),
	))
}

func RestartSelf() error {
	self, err := os.Executable()
	if err != nil {
		return err
	}
	args := os.Args
	env := os.Environ()

	return syscall.Exec(self, args, env)
}
