# Spotify Tray

_MacOS only_

![preview](./preview.png)

## Description

This application fetches spotify status with osascript from the Spotify Desktop application and displays the result in the system tray.

## How to use the app

Download the latest version of the app dmg from the [releases page](https://github.com/aaronhallaert/spotify-tray/releases)

Launch the app or from the terminal as a process. A system tray should appear with information about current track data from the Spotify desktop app. You can format this information with some options by clicking on the system tray:

- `Show artist first`: show the artist first and then the title (default: `true`)
- `Show album`: show the album as well (default: `false`, lot of space get's taken in if you also enable the album) 
- `Show progress`: show a percentage with how far the track has progressed (default: `true`)
- `Use more space`: use `64` characters for artist and `64` characters for track title, otherwise use `20` characters (default: `true`)
- `Open at login`: open the app when mac starts up or user logs in

## How to create the app

- Install go on your machine (`brew install go`)

### CLI

- Build with `make build`
- Execute with `./resources/spotify-tray`

### App

- Build with `make bundle`
- Open the `Spotify Tray.dmg` that was created in the `dist` folder

### Generate icon

You can generate a custom tray icon with `make generate_icon icon=path-to-file`
