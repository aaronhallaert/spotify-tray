package spotifydata

import (
	"os"
	"strings"
	"testing"
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
			if data.PlayerState != test.state || data.Status != test.status || data.Progress != test.progress {
				t.Fatalf("decoded %+v, want state %q, status %q, progress %d",
					data, test.state, test.status, test.progress)
			}
		})
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

func TestSpotifyIntegration(t *testing.T) {
	if os.Getenv("SPOTIFY_TRAY_INTEGRATION_TEST") != "1" {
		t.Skip("set SPOTIFY_TRAY_INTEGRATION_TEST=1 to read the running Spotify app")
	}
	for i := 0; i < 120; i++ {
		progress, album := i%2 == 0, i%4 < 2
		data, err := GetData(progress, album)
		if err != nil {
			t.Fatal(err)
		}
		if data.PlayerState == "" {
			t.Fatal("Spotify must be running for this integration test")
		}
		if !progress && (data.Duration != 0 || data.Position != 0 || data.Progress != 0) {
			t.Fatalf("read progress when disabled: %+v", data)
		}
		if !album && data.Album != "" {
			t.Fatalf("read album when disabled: %+v", data)
		}
		if i == 0 {
			t.Logf("state=%s, track present=%t, duration=%g, position=%g, progress=%d",
				data.PlayerState, strings.TrimSpace(data.Track) != "", data.Duration, data.Position, data.Progress)
		}
	}
}
