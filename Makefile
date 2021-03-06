# For a quick tutorial on `make` visit:
# https://gist.github.com/isaacs/62a2d1825d04437c6f08

build:
	go build -o ./dist/spotify-tray main.go

macapp: build
	go run ./config/macapp.go \
		-assets ./dist \
		-bin spotify-tray \
		-icon ./assets/appicon1024.png \
		-identifier com.spotify-tray \
		-name "Spotify Tray" \
		-dmg ./config/template.dmg \
		-o ./release

generate_icon:
	./config/generate_icon.sh $(icon)