package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
    "strconv"

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
    track string
    artist string
    album string
    status string
    duration float64
    position float64
    progress int
}

func (s *SpotifyStatus) Format() string {
    // return fmt.Sprintf("%s - %s", s.status, s.track)
    return fmt.Sprintf("%d%% %s  %s - %s", s.progress, s.status, s.track, s.artist)
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
    durationFloat, err := strconv.ParseFloat(duration, 64)
    durationFloat = durationFloat / 1000


    nPosition, err := exec.Command("/bin/sh", scriptsPath+"position.sh").Output()
    if err != nil {
        fmt.Printf("error %s", err)
    }
    position := strings.TrimSuffix(string(nPosition), "\n")
    position = strings.ReplaceAll(position, ",", ".")
    positionFloat, err := strconv.ParseFloat(position, 64)
    progress := int((positionFloat / durationFloat) * 100)

    return SpotifyStatus {
        track: track,
        artist: artist,
        album: album,
        status: status,
        duration: durationFloat,
        position: positionFloat,
        progress: progress,
    }
}

func onReady() {
	systray.SetTitle("Loading...")
    mUrl := systray.AddMenuItem("Spotify", "my home")
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")

    systray.AddSeparator()
    go func() {
        <-mQuitOrig.ClickedCh
        fmt.Println("Requesting quit")
        systray.Quit()
        fmt.Println("Finished quitting")
    }()

	go func() {
        select {
        case <-mUrl.ClickedCh:
            open.Run("https://www.spotify.com")
            return

        }
	}()

	go func() {
        for {
            newSpotifyStatus := fetchSpotifyStatus()
            title := newSpotifyStatus.Format()
            systray.SetTitle(title)
            time.Sleep(time.Millisecond * 300)
        }
	}()
}
