package spotifydata

import (
	"encoding/json"
	"fmt"
	"math"
	"sync"
	"time"
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

const metadataRefreshInterval = 30 * time.Second

type spotifySnapshot struct {
	Data
	TrackID           string
	ProcessID         int
	MetadataRefreshed bool
}

type spotifyReader struct {
	mu          sync.Mutex
	read        func(bool, bool, string, int) ([]byte, error)
	cached      spotifySnapshot
	cachedAt    time.Time
	getProgress bool
	getAlbum    bool
}

var defaultReader = spotifyReader{read: readSpotify}

func GetData(getProgress bool, getAlbum bool) (*Data, error) {
	return defaultReader.getData(getProgress, getAlbum, time.Now())
}

func (r *spotifyReader) getData(getProgress, getAlbum bool, now time.Time) (*Data, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	cachedID, cachedPID := r.cached.TrackID, r.cached.ProcessID
	if now.Sub(r.cachedAt) >= metadataRefreshInterval ||
		getProgress != r.getProgress || getAlbum != r.getAlbum {
		cachedID, cachedPID = "", 0
	}
	snapshot, err := r.read(getProgress, getAlbum, cachedID, cachedPID)
	if err != nil {
		r.cached = spotifySnapshot{}
		return nil, fmt.Errorf("read Spotify: %w", err)
	}

	data, err := decodeSnapshot(snapshot)
	if err != nil {
		r.cached = spotifySnapshot{}
		return nil, err
	}
	if data.PlayerState == "playing" || data.PlayerState == "paused" {
		if !data.MetadataRefreshed {
			if cachedID == "" || data.TrackID != cachedID || data.ProcessID != cachedPID {
				r.cached = spotifySnapshot{}
				return nil, fmt.Errorf("Spotify returned incomplete track metadata")
			}
			data.Track, data.Artist = r.cached.Track, r.cached.Artist
			data.Album, data.Duration = r.cached.Album, r.cached.Duration
		} else {
			r.cachedAt = now
		}
		r.cached = *data
		r.getProgress, r.getAlbum = getProgress, getAlbum
	} else {
		r.cached = spotifySnapshot{}
	}
	return finishData(&data.Data), nil
}

func decodeSnapshot(snapshot []byte) (*spotifySnapshot, error) {
	var data spotifySnapshot
	if err := json.Unmarshal(snapshot, &data); err != nil {
		return nil, fmt.Errorf("decode Spotify data: %w", err)
	}
	return &data, nil
}

func finishData(data *Data) *Data {
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
	return data
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
