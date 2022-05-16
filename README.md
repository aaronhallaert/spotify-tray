# Spotify Tray

_MacOS only_

![preview](./preview.png)

## Description

This application fetches spotify status with osascript from the Spotify Desktop application and displays the result in the system tray.

## How to create the app

- Install go on your machine (`brew install go`)

### CLI

- Build with `make build`
- Execute with `./resources/spotify-tray`

### App

- Build with `make bundle`
- Open the `Spotify Tray.dmg` that was created in the `dist` folder

### Generate icon

You can generate a custom tray icon with `make generate_icon 'path-to-file'` or you can generate the default icon with `make generate_icon_default`.

## How to use the app

Launch the app or from the terminal as a process. A system tray should appear with information about current track data from the Spotify desktop app. You can format this information with some options by clicking on the system tray:

- `Show progress`: show a percentage with how far the track has progressed (default: `true`)
- `Show artist first`: show the artist first and then the title (default: `true`)
- `Use more space`: use `64` characters for artist and `64` characters for track title, otherwise use `20` characters (default: `true`)
