package spotifydata

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestDecodeSnapshot(t *testing.T) {
	tests := []struct {
		name     string
		snapshot string
		state    string
		status   string
		progress int
	}{
		{"playing", `{"PlayerState":"playing","Duration":200,"Position":51.5}`, "playing", "▶︎", 25},
		{"paused", `{"PlayerState":"paused","Duration":200,"Position":100}`, "paused", "❚❚", 50},
		{"stopped", `{"PlayerState":"stopped"}`, "stopped", "■", 0},
		{"closed", `{}`, "", "", 0},
		{"progress disabled", `{"PlayerState":"playing"}`, "playing", "▶︎", 0},
		{"zero duration", `{"Duration":0,"Position":20}`, "", "", 0},
		{"negative duration", `{"Duration":-1,"Position":20}`, "", "", 0},
		{"past end", `{"Duration":100,"Position":101}`, "", "", 100},
		{"negative position", `{"Duration":100,"Position":-1}`, "", "", 0},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, err := decodeSnapshot([]byte(test.snapshot))
			if err != nil {
				t.Fatal(err)
			}
			finishData(&data.Data)
			if data.PlayerState != test.state || data.Status != test.status || data.Progress != test.progress {
				t.Fatalf("decoded %+v, want state %q, status %q, progress %d",
					data, test.state, test.status, test.progress)
			}
		})
	}
}

func TestSpotifyReaderCachesMetadata(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	id, pid, state := "track-1", 42, "playing"
	position := 10
	metadataReads := 0
	reader := spotifyReader{read: func(progress, album bool, cachedID string, cachedPID int) ([]byte, error) {
		if state == "" || state == "stopped" {
			return fmt.Appendf(nil, `{"PlayerState":%q}`, state), nil
		}
		refreshed := id == "" || id != cachedID || pid != cachedPID
		if refreshed {
			metadataReads++
			duration, albumName := 0, ""
			if progress {
				duration = 100
			}
			if album {
				albumName = "Album"
			}
			return fmt.Appendf(nil, `{"PlayerState":%q,"TrackID":%q,"ProcessID":%d,"MetadataRefreshed":true,"Track":"Song","Artist":"Artist","Album":%q,"Duration":%d,"Position":%d}`,
				state, id, pid, albumName, duration, position), nil
		}
		return fmt.Appendf(nil, `{"PlayerState":%q,"TrackID":%q,"ProcessID":%d,"Position":%d}`,
			state, id, pid, position), nil
	}}

	check := func(progress, album bool, expectedReads int) *Data {
		t.Helper()
		data, err := reader.getData(progress, album, now)
		if err != nil {
			t.Fatal(err)
		}
		if metadataReads != expectedReads {
			t.Fatalf("metadata reads = %d, want %d", metadataReads, expectedReads)
		}
		return data
	}

	if data := check(true, false, 1); data.Track != "Song" || data.Progress != 10 {
		t.Fatalf("unexpected first read: %+v", data)
	}
	now = now.Add(time.Second)
	position, state = 30, "paused"
	if data := check(true, false, 1); data.Track != "Song" || data.Progress != 30 || data.PlayerState != "paused" {
		t.Fatalf("cached metadata or paused seek lost: %+v", data)
	}
	id = "track-2"
	check(true, false, 2)
	if data := check(true, true, 3); data.Album != "Album" {
		t.Fatalf("album not refreshed after enabling: %+v", data)
	}
	if data := check(false, true, 4); data.Duration != 0 || data.Progress != 0 {
		t.Fatalf("progress not cleared after disabling: %+v", data)
	}
	pid++
	check(false, true, 5)
	now = now.Add(metadataRefreshInterval)
	check(false, true, 6)
	id = ""
	check(false, true, 7)
	check(false, true, 8)
	state, id = "stopped", "track-2"
	check(false, true, 8)
	state = "playing"
	check(false, true, 9)
	state = ""
	check(false, true, 9)
	state = "playing"
	check(false, true, 10)
}

func TestSpotifyReaderInvalidatesCacheOnErrors(t *testing.T) {
	now := time.Now()
	call := 0
	reader := spotifyReader{read: func(_, _ bool, cachedID string, _ int) ([]byte, error) {
		call++
		switch call {
		case 1, 4:
			if cachedID != "" {
				t.Fatalf("unexpected cached ID on call %d: %q", call, cachedID)
			}
			return []byte(`{"PlayerState":"playing","TrackID":"id","ProcessID":1,"MetadataRefreshed":true,"Track":"Song"}`), nil
		case 2:
			return nil, errors.New("Spotify unavailable")
		case 3:
			return []byte(`{`), nil
		default:
			t.Fatal("unexpected read")
			return nil, nil
		}
	}}
	if _, err := reader.getData(false, false, now); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if _, err := reader.getData(false, false, now); err == nil {
			t.Fatal("expected read error")
		}
	}
	if _, err := reader.getData(false, false, now); err != nil {
		t.Fatal(err)
	}
}

func TestDecodeSnapshotMetadata(t *testing.T) {
	data, err := decodeSnapshot([]byte(`{"Track":"Song\n\"live\", part 1","Artist":"Björk","Album":"A\tB"}`))
	if err != nil {
		t.Fatal(err)
	}
	if data.Track != "Song\n\"live\", part 1" || data.Artist != "Björk" || data.Album != "A\tB" {
		t.Fatalf("metadata changed: %+v", data)
	}
}

func TestDecodeSnapshotError(t *testing.T) {
	for _, snapshot := range []string{`{`, `{"Duration":"invalid"}`} {
		if _, err := decodeSnapshot([]byte(snapshot)); err == nil {
			t.Fatalf("expected error for %q", snapshot)
		}
	}
}
