# Spotify Tray

_MacOS only_

![preview](./preview.png)

## Description

This application fetches spotify status with osascript from the Spotify Desktop application and displays the result in the system tray.

## How to use

- Install go on your machine (`brew install go`)
- Build with `go build`
- Execute with `./spotify-tray`

## Status

- [ ] Package in an MacOS application

  - workaround possible with `Automator`:

    `cd ${project_folder} && ${project_folder}/spotify-tray`
