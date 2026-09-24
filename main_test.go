package main

import (
	"spotify-tray/spotifydata"
	"testing"
)

func TestPlaybackMenuTitle(t *testing.T) {
	tests := []struct {
		name    string
		data    spotifydata.Data
		title   string
		visible bool
	}{
		{"playing", spotifydata.Data{Track: "Song", PlayerState: "playing"}, "Pause", true},
		{"paused", spotifydata.Data{Track: "Song", PlayerState: "paused"}, "Play", true},
		{"stopped", spotifydata.Data{Track: "Song", PlayerState: "stopped"}, "", false},
		{"closed", spotifydata.Data{}, "", false},
		{"no track", spotifydata.Data{PlayerState: "paused"}, "", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			title, visible := playbackMenuTitle(&test.data)
			if title != test.title || visible != test.visible {
				t.Errorf("playbackMenuTitle(%+v) = (%q, %v), want (%q, %v)", test.data, title, visible, test.title, test.visible)
			}
		})
	}
}
