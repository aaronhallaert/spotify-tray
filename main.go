package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
)

func main() {
	onExit := func() {
	}

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
	mLyrics := systray.AddMenuItem("Lyrics", "Search lyrics")
	systray.AddSeparator()
	mProgress := systray.AddMenuItemCheckbox("Show progress?", "Show Progress", true)
	mArtistFirst := systray.AddMenuItemCheckbox("Show artist first?", "Show artist first", true)
	mMoreSpace := systray.AddMenuItemCheckbox("Use more space?", "Use more space", true)
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
		open.Run("https://www.google.be/search?q=" + currentSpotifyStatus.track + " lyrics")
	}()

	go func() {
		for {
			select {
			case <-mProgress.ClickedCh:
				if mProgress.Checked() {
					mProgress.Uncheck()
				} else {
					mProgress.Check()
				}

			case <-mArtistFirst.ClickedCh:
				if mArtistFirst.Checked() {
					mArtistFirst.Uncheck()
				} else {
					mArtistFirst.Check()
				}

			case <-mMoreSpace.ClickedCh:
				if mMoreSpace.Checked() {
					mMoreSpace.Uncheck()
				} else {
					mMoreSpace.Check()
				}
			}
		}
	}()

	go func() {
		for {
			currentSpotifyStatus = fetchSpotifyStatus()
			message := currentSpotifyStatus.Format(mProgress.Checked(), mArtistFirst.Checked(), mMoreSpace.Checked())
			systray.SetTitle(message)
			time.Sleep(time.Millisecond * 300)
		}
	}()
}
