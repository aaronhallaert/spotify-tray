package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"spotify-tray/storage"
	"strconv"
	"strings"
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

type SpotifyStatus struct {
	track    string
	artist   string
	album    string
	status   string
	duration float64
	position float64
	progress int
}

func (s *SpotifyStatus) Format(showProgress bool, isArtistFirst bool, isMoreSpace bool) string {
	if len(s.track) == 0 {
		return fmt.Sprintf("%s  Spotify is not playing!", s.status)
	}

	formatProgres := fmt.Sprintf("  |  %d%%", s.progress)
	if !showProgress {
		formatProgres = ""
	}

	formatStrLength := 64
	if !isMoreSpace {
		formatStrLength = 20
	}

	if len(s.artist) == 0 {
		return fmt.Sprintf("%s  %s%s", s.status, trimString(s.track, formatStrLength), formatProgres)
	}

	artistAndTrack := [2]string{trimString(s.artist, formatStrLength), trimString(s.track, formatStrLength)}
	if !isArtistFirst {
		artistAndTrack = [2]string{trimString(s.track, formatStrLength), trimString(s.artist, formatStrLength)}
	}

	return fmt.Sprintf("%s  %s - %s%s", s.status, artistAndTrack[0], artistAndTrack[1], formatProgres)
}

func trimString(s string, maxLength int) string {
	if len(s) > maxLength {
		trimmed := s[:maxLength] + "..."
		return trimmed
	}
	return s
}

func fetchSpotifyStatus() SpotifyStatus {
	executable, _ := os.Executable()
	scriptsPath := filepath.Join(filepath.Dir(executable), "../Resources/") + "/"
	if !strings.Contains(filepath.Dir(executable), "MacOS") {
		scriptsPath = filepath.Dir(executable) + "/"
	}

	nTrack, err := exec.Command("/bin/sh", scriptsPath+"track.sh").Output()
	if err != nil {
		fmt.Printf("error %s", err)
	}
	track := strings.TrimSuffix(string(nTrack), "\n")

	nArtist, err := exec.Command("/bin/sh", scriptsPath+"artist.sh").Output()
	if err != nil {
		fmt.Printf("error %s", err)
	}
	artist := strings.TrimSuffix(string(nArtist), "\n")

	nStatus, err := exec.Command("/bin/sh", scriptsPath+"status.sh").Output()
	if err != nil {
		fmt.Printf("error %s", err)
	}
	status := strings.TrimSuffix(string(nStatus), "\n")

	nAlbum, err := exec.Command("/bin/sh", scriptsPath+"album.sh").Output()
	if err != nil {
		fmt.Printf("error %s", err)
	}
	album := strings.TrimSuffix(string(nAlbum), "\n")

	nDuration, err := exec.Command("/bin/sh", scriptsPath+"duration.sh").Output()
	if err != nil {
		fmt.Printf("error %s", err)
	}
	duration := strings.TrimSuffix(string(nDuration), "\n")
	durationFloat, _ := strconv.ParseFloat(duration, 64)
	durationFloat = durationFloat / 1000

	nPosition, err := exec.Command("/bin/sh", scriptsPath+"position.sh").Output()
	if err != nil {
		fmt.Printf("error %s", err)
	}
	position := strings.TrimSuffix(string(nPosition), "\n")
	position = strings.ReplaceAll(position, ",", ".")
	positionFloat, _ := strconv.ParseFloat(position, 64)
	progress := int((positionFloat / durationFloat) * 100)

	return SpotifyStatus{
		track:    track,
		artist:   artist,
		album:    album,
		status:   status,
		duration: durationFloat,
		position: positionFloat,
		progress: progress,
	}
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

	currentSpotifyStatus := fetchSpotifyStatus()

	go func() {
		<-mQuitOrig.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()

	go func() {
		<-mLyrics.ClickedCh
		open.Run("https://www.google.be/search?q=" + currentSpotifyStatus.track + "-" + currentSpotifyStatus.artist + " lyrics")
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
		for {
			currentSpotifyStatus = fetchSpotifyStatus()
			message := currentSpotifyStatus.Format(storage.GetHasProgress(), storage.GetArtistFirst(), storage.GetMoreSpace())
			systray.SetTitle(message)
			time.Sleep(time.Millisecond * 300)
		}
	}()
}
