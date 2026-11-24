package spotifydata

import (
	"encoding/json"
	"fmt"
	"math"
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

func GetData(getProgress bool, getAlbum bool) (*Data, error) {
	snapshot, err := readSpotify(getProgress, getAlbum)
	if err != nil {
		return nil, fmt.Errorf("read Spotify: %w", err)
	}
	return decodeSnapshot(snapshot)
}

func decodeSnapshot(snapshot []byte) (*Data, error) {
	var data Data
	if err := json.Unmarshal(snapshot, &data); err != nil {
		return nil, fmt.Errorf("decode Spotify data: %w", err)
	}

	switch data.PlayerState {
	case "playing":
		data.Status = "▶︎"
	case "paused":
		data.Status = "❚❚"
	case "stopped":
		data.Status = "■"
	}

	if data.Duration > 0 {
		data.Progress = int(math.Max(0, math.Min(100, data.Position/data.Duration*100)))
	}
	return &data, nil
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

func TogglePlayback() error {
	if err := toggleSpotify(); err != nil {
		return fmt.Errorf("toggle Spotify playback: %w", err)
	}
	return nil
}
