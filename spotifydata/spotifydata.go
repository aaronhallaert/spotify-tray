package spotifydata

import (
	"fmt"
	"os/exec"
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
	track, artist, album, status, duration, position := "", "", "", "", 1.0, 0.0
	var wg sync.WaitGroup
	wg.Add(6)

	go func() {
		defer wg.Done()
		track = getValueFromScript("name of current track")
	}()
	go func() {
		defer wg.Done()
		artist = getValueFromScript("artist of current track")
	}()
	go func() {
		defer wg.Done()
		album = getValueFromScript("album of current track")
	}()
	go func() {
		defer wg.Done()
		status = getValueFromScript("player state")
	}()
	go func() {
		defer wg.Done()
		durationFloat, _ := strconv.ParseFloat(getValueFromScript("duration of current track"), 64)
		duration = durationFloat / 1000
	}()
	go func() {
		defer wg.Done()
		position, _ = strconv.ParseFloat(strings.ReplaceAll(getValueFromScript("player position"), ",", "."), 64)
	}()

	wg.Wait()

	progress := int((position / duration) * 100)
	statusIcon := "■"
	if status == "playing" {
		statusIcon = "▶︎"
	} else if status == "paused" {
		statusIcon = "❚❚"
	}

	return &Data{
		track,
		artist,
		album,
		statusIcon,
		duration,
		position,
		progress,
	}
}

func getValueFromScript(prop string) string {
	nValue, err := exec.Command("osascript", "-e", "if application \"Spotify\" is running then\n tell application \"Spotify\"\n return "+prop+" as string \nend tell \nend if").Output()

	if err != nil {
		fmt.Printf("error %s", err)
	}

	return strings.TrimSuffix(string(nValue), "\n")
}

func (d *Data) Format(showProgress bool, showAlbum bool, isArtistFirst bool, isMoreSpace bool, isAlternateSeperator bool) string {
	if len(d.Track) == 0 {
		return fmt.Sprintf("%s Spotify is not playing!", d.Status)
	}

	formatProgres := fmt.Sprintf("  |  %d%%", d.Progress)
	if !showProgress {
		formatProgres = ""
	}

	formatStrLength := 64
	if !isMoreSpace {
		formatStrLength = 20
	}

	seperator := "-"
	if isAlternateSeperator {
		seperator = " -§- "
	}

	formatAlbum := fmt.Sprintf(" %s %s", seperator, trimString(d.Album, formatStrLength))
	if !showAlbum {
		formatAlbum = ""
	}

	if len(d.Artist) == 0 {
		return fmt.Sprintf("%s  %s%s", d.Status, trimString(d.Track, formatStrLength), formatProgres)
	}

	artistAndTrack := [2]string{trimString(d.Artist, formatStrLength), trimString(d.Track, formatStrLength)}
	if !isArtistFirst {
		artistAndTrack = [2]string{trimString(d.Track, formatStrLength), trimString(d.Artist, formatStrLength)}
	}

	return fmt.Sprintf("%s  %s %s %s%s%s", d.Status, artistAndTrack[0], seperator, artistAndTrack[1], formatAlbum, formatProgres)
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
	_, err := strconv.ParseBool(strings.TrimSuffix(string(nValue), "\n"))

	if err == nil {
		return true
	} else {
		return false
	}
}
