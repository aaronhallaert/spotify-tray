package main

import (
	"fmt"
	"os/exec"
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
    track string
    artist string
    album string
    status string
}

func (s *SpotifyStatus) Format() string {
    // return fmt.Sprintf("%s - %s", s.status, s.track)
    return fmt.Sprintf("%s  %s - %s", s.status, s.track, s.artist)
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

    return SpotifyStatus {
        track: track,
        artist: artist,
        album: album,
        status: status,
    }
}

func onReady() {
	systray.SetTitle("Loading...")
    mUrl := systray.AddMenuItem("Spotify", "my home")
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")

    systray.AddSeparator()

	go func() {
        select {
        case <-mQuitOrig.ClickedCh:
            fmt.Println("Requesting quit")
            systray.Quit()
            fmt.Println("Finished quitting")
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


	// We can manipulate the systray in other goroutines
	go func() {
		systray.SetTooltip("")
	}()
}
