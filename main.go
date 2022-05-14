package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
)

var scriptsPath = "./scripts/"

// var scriptsPath = "../Resources/"

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

func (s *SpotifyStatus) Format() string {
	if len(s.track) == 0 {
		return fmt.Sprintf("%s  Spotify is not running!", s.status)
	}

	if len(s.artist) == 0 {
		return fmt.Sprintf("%s  %s  |  %d%%", s.status, trimString(s.track, 64), s.progress)
	}

	return fmt.Sprintf("%s  %s - %s  |  %d%%", s.status, trimString(s.artist, 64), trimString(s.track, 64), s.progress)
}

func trimString(s string, maxLength int) string {
	if len(s) > maxLength {
		trimmed := s[:maxLength] + "..."
		return trimmed
	}
	return s
}

func fetchSpotifyStatus() SpotifyStatus {
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
	systray.SetTitle(" Loading...")
	mLyrics := systray.AddMenuItem("Lyrics", "Search lyrics")
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")

	currentSpotifyStatus := fetchSpotifyStatus()

	systray.AddSeparator()
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
			currentSpotifyStatus = fetchSpotifyStatus()
			message := currentSpotifyStatus.Format()
			systray.SetTitle(message)
			time.Sleep(time.Millisecond * 300)
		}
	}()
}
