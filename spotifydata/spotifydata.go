package spotifydata

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
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

func Init() *Data {
	return getData()
}

func (d *Data) Update() {
	newStatus := getData()

	d.Track = newStatus.Track
	d.Artist = newStatus.Artist
	d.Album = newStatus.Album
	d.Status = newStatus.Status
	d.Duration = newStatus.Duration
	d.Position = newStatus.Position
	d.Progress = newStatus.Progress
}

func getData() *Data {
	track := getValueFromScript("name of current track")
	artist := getValueFromScript("artist of current track")
	status := getValueFromScript("player state")
	album := getValueFromScript("album of current track")
	duration := getValueFromScript("duration of current track")
	position := strings.ReplaceAll(getValueFromScript("player position"), ",", ".")

	durationFloat, _ := strconv.ParseFloat(duration, 64)
	durationFloat = durationFloat / 1000
	positionFloat, _ := strconv.ParseFloat(position, 64)
	progress := int((positionFloat / durationFloat) * 100)

	return &Data{
		Track:    track,
		Artist:   artist,
		Album:    album,
		Status:   status,
		Duration: durationFloat,
		Position: positionFloat,
		Progress: progress,
	}
}

func (d *Data) Format(showProgress bool, isArtistFirst bool, isMoreSpace bool) string {
	if len(d.Track) == 0 {
		return fmt.Sprintf("%s  Spotify is not playing!", d.Status)
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
		return fmt.Sprintf("%s  %s%s", d.Status, trimString(d.Track, formatStrLength), formatProgres)
	}

	artistAndTrack := [2]string{trimString(d.Artist, formatStrLength), trimString(d.Track, formatStrLength)}
	if !isArtistFirst {
		artistAndTrack = [2]string{trimString(d.Track, formatStrLength), trimString(d.Artist, formatStrLength)}
	}

	return fmt.Sprintf("%s  %s - %s%s", d.Status, artistAndTrack[0], artistAndTrack[1], formatProgres)
}

func trimString(s string, maxLength int) string {
	if len(s) > maxLength {
		trimmed := s[:maxLength] + "..."
		return trimmed
	}
	return s
}

func getValueFromScript(prop string) string {
	nValue, err := exec.Command("osascript", "-e", "if application \"Spotify\" is running then\n tell application \"Spotify\"\n return "+prop+" as string \nend tell \nend if").Output()

	if err != nil {
		fmt.Printf("error %s", err)
	}

	return strings.TrimSuffix(string(nValue), "\n")
}
