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
	mArtistFirst := systray.AddMenuItemCheckbox("Show artist first?", "Show artist first", storage.GetArtistFirst())
	mShowAlbum := systray.AddMenuItemCheckbox("Show album?", "Show Album", storage.GetShowAlbum())
	mProgress := systray.AddMenuItemCheckbox("Show progress?", "Show Progress", storage.GetShowProgress())
	systray.AddSeparator()
	mMoreSpace := systray.AddMenuItemCheckbox("Use more space?", "Use more space", storage.GetMoreSpace())
	mOpenAtLogin := systray.AddMenuItemCheckbox("Open at login?", "Open at login", storage.GetOpenAtLogin())
	systray.AddSeparator()
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")

	currentSpotifyData := &spotifydata.Data{}
	if spotifydata.IsSpotifyRunning() {
		currentSpotifyData = spotifydata.GetData()
	}
	updateTray(currentSpotifyData)

	go func() {
		<-mLyrics.ClickedCh
		open.Run("https://www.google.be/search?q=" + currentSpotifyData.Track + " - " + currentSpotifyData.Artist + " lyrics")
	}()

	go func() {
		for {
			select {
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
				currentSpotifyData = spotifydata.GetData()
				updateTray(currentSpotifyData)
			} else {
				currentSpotifyData.Status = ""
			}
			time.Sleep(time.Millisecond * 500)
		}
	}()
}

func updateTray(d *spotifydata.Data) {
	message := d.Format(storage.GetShowProgress(), storage.GetShowAlbum(), storage.GetArtistFirst(), storage.GetMoreSpace())
	// systray.SetTemplateIcon(d.GetIcon(), d.GetIcon())
	systray.SetTitle(message)
}
