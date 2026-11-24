package main

import (
	"fmt"
	"spotify-tray/internal/tray"
	"spotify-tray/spotifydata"
	"spotify-tray/storage"
	"time"
)

func main() {
	onExit := func() {
	}

	storage.Init()
	tray.Run(onReady, onExit)
}

func onReady() {
	tray.SetTitle("Loading...")
	trayTitle := titleUpdater{current: "Loading...", update: tray.SetTitle}
	mPlayback := tray.AddMenuItem("Play", "Toggle Spotify playback")
	mPlayback.Hide()
	tray.AddSeparator()
	mArtistFirst := tray.AddMenuItemCheckbox("Show artist first?", "Show artist first", storage.GetArtistFirst())
	mShowAlbum := tray.AddMenuItemCheckbox("Show album?", "Show Album", storage.GetShowAlbum())
	mProgress := tray.AddMenuItemCheckbox("Show progress?", "Show Progress", storage.GetShowProgress())
	tray.AddSeparator()
	mMoreSpace := tray.AddMenuItemCheckbox("Use more space?", "Use more space", storage.GetMoreSpace())
	mAlternateSeparator := tray.AddMenuItemCheckbox("Use alternate separator?", "Show Alternate separator", storage.GetAlternateSeparator())
	mOpenAtLogin := tray.AddMenuItemCheckbox("Open at login?", "Open at login", storage.GetOpenAtLogin())
	tray.AddSeparator()
	mQuitOrig := tray.AddMenuItem("Quit", "Quit the whole app")

	playbackTitle := ""
	playbackVisible := false
	var retry retryBackoff
	refresh := func() {
		data, err := spotifydata.GetData(storage.GetShowProgress(), storage.GetShowAlbum())
		if err != nil {
			if retry.fail(err, time.Now()) {
				fmt.Printf("error %s\n", err)
			}
			trayTitle.set("Unable to read Spotify")
			if playbackVisible {
				mPlayback.Hide()
			}
			playbackVisible = false
			return
		}
		retry.reset()
		trayTitle.set(data.Format(
			storage.GetShowProgress(),
			storage.GetShowAlbum(),
			storage.GetArtistFirst(),
			storage.GetMoreSpace(),
			storage.GetAlternateSeparator(),
		))
		title, visible := playbackMenuTitle(data)
		if visible {
			if title != playbackTitle {
				mPlayback.SetTitle(title)
				playbackTitle = title
			}
			if !playbackVisible {
				mPlayback.Show()
			}
		} else if playbackVisible {
			mPlayback.Hide()
		}
		playbackVisible = visible
	}
	refresh()

	playbackChanged := make(chan struct{}, 1)
	go func() {
		for range mPlayback.ClickedCh {
			if err := spotifydata.TogglePlayback(); err != nil {
				fmt.Printf("error %s\n", err)
			}
			select {
			case playbackChanged <- struct{}{}:
			default:
			}
		}
	}()

	go func() {
		<-mQuitOrig.ClickedCh
		tray.Quit()
	}()

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if retry.ready(time.Now()) {
					refresh()
				}
			case <-playbackChanged:
				refresh()
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
}

func playbackMenuTitle(data *spotifydata.Data) (string, bool) {
	if data.Track == "" {
		return "", false
	}
	switch data.PlayerState {
	case "playing":
		return "Pause", true
	case "paused":
		return "Play", true
	default:
		return "", false
	}
}

type titleUpdater struct {
	current string
	update  func(string)
}

type retryBackoff struct {
	next      time.Time
	delay     time.Duration
	lastError string
}

func (r *retryBackoff) ready(now time.Time) bool {
	return !now.Before(r.next)
}

func (r *retryBackoff) fail(err error, now time.Time) bool {
	if r.delay == 0 {
		r.delay = time.Second
	} else {
		r.delay = min(r.delay*2, 30*time.Second)
	}
	r.next = now.Add(r.delay)
	message := err.Error()
	changed := message != r.lastError
	r.lastError = message
	return changed
}

func (r *retryBackoff) reset() {
	*r = retryBackoff{}
}

func (t *titleUpdater) set(title string) {
	if title != t.current {
		t.update(title)
		t.current = title
	}
}
