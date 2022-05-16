package spotifydata

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

var scriptsPath = getScriptsPath()

func Init() Data {
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

func getData() Data {
	track := getValueFromScript("track.sh")
	artist := getValueFromScript("artist.sh")
	status := getValueFromScript("status.sh")
	album := getValueFromScript("album.sh")
	duration := getValueFromScript("duration.sh")
	position := strings.ReplaceAll(getValueFromScript("position.sh"), ",", ".")

	durationFloat, _ := strconv.ParseFloat(duration, 64)
	durationFloat = durationFloat / 1000
	positionFloat, _ := strconv.ParseFloat(position, 64)
	progress := int((positionFloat / durationFloat) * 100)

	return Data{
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

func getValueFromScript(file string) string {
	nValue, err := exec.Command("/bin/sh", scriptsPath+file).Output()
	if err != nil {
		fmt.Printf("error %s", err)
	}

	return strings.TrimSuffix(string(nValue), "\n")
}

func getScriptsPath() string {
	executable, _ := os.Executable()
	path := filepath.Join(filepath.Dir(executable), "../Resources/") + "/"
	if !strings.Contains(filepath.Dir(executable), "MacOS") {
		path = filepath.Dir(executable) + "/scripts/"
	}

	return path
}
