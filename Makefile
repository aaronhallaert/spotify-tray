# For a quick tutorial on `make` visit:
# https://gist.github.com/isaacs/62a2d1825d04437c6f08

build:
	go build -o ./resources/spotify-tray main.go

bundle: build
	go run ./config/macapp.go \
		-assets ./resources \
		-bin spotify-tray \
		-icon ./config/appicon1024.png \
		-identifier com.spotify-tray \
		-name "Spotify Tray" \
		-dmg ./config/template.dmg \
		-o ./dist

generate_icon:
	./config/generate_icon.sh $(icon)