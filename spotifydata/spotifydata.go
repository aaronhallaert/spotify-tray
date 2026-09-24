package spotifydata

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type Data struct {
	Track       string
	Artist      string
	Album       string
	Status      string
	PlayerState string
	Duration    float64
	Position    float64
	Progress    int
}

func GetData(getProgress bool, getAlbum bool) *Data {
	track, artist, album, status, duration, position := "", "", "", "", 1.0, 0.0

	track = getValueFromScript("name of current track")
	artist = getValueFromScript("artist of current track")
	status = getValueFromScript("player state")

	if getAlbum {
		album = getValueFromScript("album of current track")
	}

	if getProgress {
		durationFloat, _ := strconv.ParseFloat(getValueFromScript("duration of current track"), 64)
		duration = durationFloat / 1000
		position, _ = strconv.ParseFloat(strings.ReplaceAll(getValueFromScript("player position"), ",", "."), 64)
	}

	progress := int((position / duration) * 100)
	statusIcon := "■"
	switch status {
	case "playing":
		statusIcon = "▶︎"
	case "paused":
		statusIcon = "❚❚"
	}

	return &Data{
		track,
		artist,
		album,
		statusIcon,
		status,
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

func (d *Data) Format(showProgress bool, showAlbum bool, isArtistFirst bool, isMoreSpace bool, isAlternateSeparator bool) string {
	if len(d.Track) == 0 {
		return fmt.Sprintf("%s Spotify is not playing!", d.Status)
	}

	formatProgres := fmt.Sprintf("  |  %d%%", d.Progress)
	if !showProgress {
		formatProgres = ""
	}

	formatStrLength, formatStrArtist := 64, 64
	if !isMoreSpace {
		formatStrLength, formatStrArtist = 25, 20
	}

	separator := "-"
	if isAlternateSeparator {
		separator = " -§- "
	}

	formatAlbum := fmt.Sprintf(" %s %s", separator, trimString(d.Album, formatStrLength))
	if !showAlbum {
		formatAlbum = ""
	}

	if len(d.Artist) == 0 {
		return fmt.Sprintf("%s  %s%s", d.Status, trimString(d.Track, formatStrLength), formatProgres)
	}

	artistAndTrack := [2]string{trimString(d.Artist, formatStrArtist), trimString(d.Track, formatStrLength)}
	if !isArtistFirst {
		artistAndTrack = [2]string{trimString(d.Track, formatStrLength), trimString(d.Artist, formatStrArtist)}
	}

	return fmt.Sprintf("%s  %s %s %s%s%s", d.Status, artistAndTrack[0], separator, artistAndTrack[1], formatAlbum, formatProgres)
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
	running, err := strconv.ParseBool(strings.TrimSpace(string(nValue)))
	return err == nil && running
}

func TogglePlayback() error {
	output, err := exec.Command("osascript", "-e", "if application \"Spotify\" is running then\n tell application \"Spotify\" to playpause\nelse\n error \"Spotify is not running\"\nend if").CombinedOutput()
	if err != nil {
		return fmt.Errorf("toggle Spotify playback: %w: %s", err, strings.TrimSpace(string(output)))
	}
	return nil
}
