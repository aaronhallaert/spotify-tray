# Spotify Tray

_MacOS only_

![preview](./preview.png)

## Description

This application fetches spotify status with osascript from the Spotify Desktop application and displays the result in the system tray.

## How to use

- Install go on your machine (`brew install go`)

### CLI

- Build with `make build`
- Execute with `./lib/spotify-tray`

### App

- Build with `make bundle`
- Open the `Spotify Tray.dmg` that was created in the `dist` folder

### Generate icon

You can generate a custom tray icon with `make generate_icon 'path-to-file'` or you can generate the default icon with `make generate_icon_default`.
