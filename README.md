# Spotify Tray

_MacOS only_

![preview](./preview.png)

## Description

This application reads Spotify Desktop status once per second using in-process macOS Apple Events and displays it in the system tray.
Track metadata is cached by track ID and refreshed at least every 30 seconds; playback state and progress are still checked every second.
The menu bar uses native AppKit APIs through cgo.

On first use, allow Spotify Tray to control Spotify when macOS asks. You can change this permission in System Settings > Privacy & Security > Automation. The open-at-login setting still uses `osascript`, but only at startup and when you change that setting.
If Spotify cannot be read, retries slow down to at most one every 30 seconds. Regular polling resumes when a read succeeds.

## How to use the app

Download the latest version of the app dmg from the [releases page](https://github.com/aaronhallaert/spotify-tray/releases)

Launch the app or from the terminal as a process. A system tray should appear with information about current track data from the Spotify desktop app. You can format this information with some options by clicking on the system tray:

- `Show artist first`: show the artist first and then the title (default: `true`)
- `Show album`: show the album as well (default: `false`, lot of space get's taken in if you also enable the album)
- `Show progress`: show a percentage with how far the track has progressed (default: `true`)
- `Use more space`: use `64` characters for artist and `64` characters for track title, otherwise use `20` characters (default: `true`)
- `Use alternate separator`: use `-§-` as a separator instead of `-`.
- `Open at login`: open the app when mac starts up or user logs in

## How to create the app

- Install go on your machine (`brew install go`)
- Install the Xcode command line tools for the native macOS build (`xcode-select --install`)

### CLI

- Build with `make build`
- Execute with `./dist/spotify-tray`
- Run the native menu bar integration test in a graphical macOS session with `SPOTIFY_TRAY_TEST_APPKIT=1 go test ./internal/tray`

### App

- Build with `make macapp`
- Open the `Spotify Tray.dmg` that was created in the `release` folder

### Generate icon

You can generate a custom tray icon with `make generate_icon icon=path-to-file`
