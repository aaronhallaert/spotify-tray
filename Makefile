# For a quick tutorial on `make` visit:
# https://gist.github.com/isaacs/62a2d1825d04437c6f08

build: go build -o scripts/spotify-tray main.go

bundle: build
	go run ./config/macapp.go \
		-assets ./scripts/ \
		-bin spotify-tray \
		-icon ./config/appicon1024.png \
		-identifier com.spotify-tray \
		-name "Spotify Tray" \
		-o ./dist
