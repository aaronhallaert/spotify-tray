package spotifydata

import (
	"fmt"
	"os/exec"
	"spotify-tray/icons"
	"strconv"
	"strings"
	"sync"
)

type Data struct {
	Track    string
	Artist   string
	Album    string
	Status   string
	Duration float64
	Position float64
	Progress int
}

func GetData() *Data {
	data := &Data{}
	var wg sync.WaitGroup
	wg.Add(6)

	go func() {
		defer wg.Done()
		data.Track = getValueFromScript("name of current track")
	}()
	go func() {
		defer wg.Done()
		data.Artist = getValueFromScript("artist of current track")
	}()
	go func() {
		defer wg.Done()
		data.Album = getValueFromScript("album of current track")
	}()
	go func() {
		defer wg.Done()
		data.Status = getValueFromScript("player state")
	}()
	go func() {
		defer wg.Done()
		duration := getValueFromScript("duration of current track")
		durationFloat, _ := strconv.ParseFloat(duration, 64)
		data.Duration = durationFloat / 1000
	}()
	go func() {
		defer wg.Done()
		position := strings.ReplaceAll(getValueFromScript("player position"), ",", ".")
		positionFloat, _ := strconv.ParseFloat(position, 64)
		data.Position = positionFloat
	}()

	wg.Wait()
	data.Progress = int((data.Position / data.Duration) * 100)

	return data
}

func getValueFromScript(prop string) string {
	nValue, err := exec.Command("osascript", "-e", "if application \"Spotify\" is running then\n tell application \"Spotify\"\n return "+prop+" as string \nend tell \nend if").Output()

	if err != nil {
		fmt.Printf("error %s", err)
	}

	return strings.TrimSuffix(string(nValue), "\n")
}

func (d *Data) GetIcon() []byte {
	statusIcon := icons.StopIcon
	if d.Status == "playing" {
		statusIcon = icons.PlayIcon
	} else if d.Status == "paused" {
		statusIcon = icons.PauseIcon
	}

	return statusIcon
}

func (d *Data) Format(showProgress bool, isArtistFirst bool, isMoreSpace bool) string {
	if len(d.Track) == 0 {
		return " Spotify is not playing!"
	}

	formatProgres := fmt.Sprintf("  |  %d%%", d.Progress)
	if !showProgress {
		formatProgres = ""
	}

	formatStrLength := 64
	if !isMoreSpace {
		formatStrLength = 20
	}

	if len(d.Artist) == 0 {
		return fmt.Sprintf(" %s%s", trimString(d.Track, formatStrLength), formatProgres)
	}

	artistAndTrack := [2]string{trimString(d.Artist, formatStrLength), trimString(d.Track, formatStrLength)}
	if !isArtistFirst {
		artistAndTrack = [2]string{trimString(d.Track, formatStrLength), trimString(d.Artist, formatStrLength)}
	}

	return fmt.Sprintf(" %s - %s%s", artistAndTrack[0], artistAndTrack[1], formatProgres)
}

func trimString(s string, maxLength int) string {
	if len(s) > maxLength {
		trimmed := s[:maxLength] + "..."
		return trimmed
	}
	return s
}

func IsSpotifyRunning() bool {
	nValue, _ := exec.Command("osascript", "-e", "if application \"Spotify\" is running then\n return true as string \nelse\n return false as string\nend if").Output()
	isSpotifyRunning, err := strconv.ParseBool(strings.TrimSuffix(string(nValue), "\n"))

	if err != nil {
		return false
	}

	return isSpotifyRunning
}
