package tests_test

import (
	"os"
	"spotify-tray/spotifydata"
	"strings"
	"testing"
)

func TestSpotifyIntegration(t *testing.T) {
	if os.Getenv("SPOTIFY_TRAY_INTEGRATION_TEST") != "1" {
		t.Skip("set SPOTIFY_TRAY_INTEGRATION_TEST=1 to read the running Spotify app")
	}
	for i := 0; i < 120; i++ {
		progress, album := i%2 == 0, i%4 < 2
		data, err := spotifydata.GetData(progress, album)
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
