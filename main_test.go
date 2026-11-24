package main

import (
	"errors"
	"reflect"
	"spotify-tray/spotifydata"
	"testing"
	"time"
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

func TestTitleUpdaterSkipsUnchangedTitles(t *testing.T) {
	var updates []string
	updater := titleUpdater{
		current: "Loading...",
		update:  func(title string) { updates = append(updates, title) },
	}

	for _, title := range []string{"Loading...", "Song", "Song", "Unable to read Spotify", "Unable to read Spotify", "Song"} {
		updater.set(title)
	}
	want := []string{"Song", "Unable to read Spotify", "Song"}
	if !reflect.DeepEqual(updates, want) {
		t.Fatalf("title updates = %q, want %q", updates, want)
	}
}

func TestRetryBackoff(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	var retry retryBackoff
	err := errors.New("permission denied")
	for i, delay := range []time.Duration{1, 2, 4, 8, 16, 30, 30} {
		if !retry.ready(now) {
			t.Fatalf("retry %d should be ready", i)
		}
		if logged := retry.fail(err, now); logged != (i == 0) {
			t.Fatalf("retry %d log = %v", i, logged)
		}
		delay *= time.Second
		if retry.ready(now.Add(delay-time.Nanosecond)) || !retry.ready(now.Add(delay)) {
			t.Fatalf("retry %d did not wait %v", i, delay)
		}
		now = now.Add(delay)
	}
	if !retry.fail(errors.New("Spotify timed out"), now) {
		t.Fatal("changed error should be logged")
	}
	retry.reset()
	if !retry.ready(now) || !retry.fail(err, now) || retry.delay != time.Second {
		t.Fatal("successful read should reset retry delay and log deduplication")
	}
}
