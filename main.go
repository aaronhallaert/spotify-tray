package main

import (
	"spotify-tray/spotifydata"
	"spotify-tray/storage"
	"time"

	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
)

func main() {
	onExit := func() {
	}

	storage.Init()
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("Loading...")
	mLyrics := systray.AddMenuItem("Lyrics", "Search for lyrics online")
	systray.AddSeparator()
	mProgress := systray.AddMenuItemCheckbox("Show progress?", "Show Progress", storage.GetHasProgress())
	mArtistFirst := systray.AddMenuItemCheckbox("Show artist first?", "Show artist first", storage.GetArtistFirst())
	mMoreSpace := systray.AddMenuItemCheckbox("Use more space?", "Use more space", storage.GetMoreSpace())
	systray.AddSeparator()
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")

	currentSpotifyData := spotifydata.Init()

	go func() {
		<-mLyrics.ClickedCh
		open.Run("https://www.google.be/search?q=" + currentSpotifyData.Track + " - " + currentSpotifyData.Artist + " lyrics")
	}()

	go func() {
		for {
			select {
			case <-mProgress.ClickedCh:
				if mProgress.Checked() {
					mProgress.Uncheck()
					storage.SetHasProgress(false)
				} else {
					mProgress.Check()
					storage.SetHasProgress(true)
				}

			case <-mArtistFirst.ClickedCh:
				if mArtistFirst.Checked() {
					mArtistFirst.Uncheck()
					storage.SetArtistFirst(false)
				} else {
					mArtistFirst.Check()
					storage.SetArtistFirst(true)
				}

			case <-mMoreSpace.ClickedCh:
				if mMoreSpace.Checked() {
					mMoreSpace.Uncheck()
					storage.SetMoreSpace(false)
				} else {
					mMoreSpace.Check()
					storage.SetMoreSpace(true)
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
			currentSpotifyData.Update()
			message := currentSpotifyData.Format(storage.GetHasProgress(), storage.GetArtistFirst(), storage.GetMoreSpace())
			systray.SetTitle(message)
			time.Sleep(time.Millisecond * 300)
		}
	}()
}
