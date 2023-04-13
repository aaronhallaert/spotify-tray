package storage

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Preferences struct {
	ShowProgress       bool
	ShowAlbum          bool
	ArtistFirst        bool
	MoreSpace          bool
	AlternateSeperator bool
}

var path = ""
var preferences = Preferences{
	ShowProgress: true,
	ShowAlbum:    false,
	ArtistFirst:  true,
	MoreSpace:    true,
}

func Init() {
	dirname, err := os.UserConfigDir()
	if err != nil {
		return
	}

	path = dirname + "/spotify-tray/preferences.json"
	readOrWriteFileIfExist(path)
}

func readOrWriteFileIfExist(fileName string) {
	_, error := os.Stat(fileName)

	if os.IsNotExist(error) {
		writeFile()
	} else {
		readFile()
	}
}

func GetShowProgress() bool {
	return preferences.ShowProgress
}
func SetShowProgress(value bool) {
	preferences.ShowProgress = value
	writeFile()
}

func GetShowAlbum() bool {
	return preferences.ShowAlbum
}
func SetShowAlbum(value bool) {
	preferences.ShowAlbum = value
	writeFile()
}

func GetArtistFirst() bool {
	return preferences.ArtistFirst
}
func SetArtistFirst(value bool) {
	preferences.ArtistFirst = value
	writeFile()
}

func GetMoreSpace() bool {
	return preferences.MoreSpace
}
func SetMoreSpace(value bool) {
	preferences.MoreSpace = value
	writeFile()
}

func GetAlternateSeperator() bool {
	return preferences.AlternateSeperator
}
func SetAlternateSeperator(value bool) {
	preferences.AlternateSeperator = value
	writeFile()
}

func GetOpenAtLogin() bool {
	entries, _ := exec.Command("osascript", "-e", "tell application \"System Events\" to get the name of every login item").Output()
	return strings.Contains(string(entries), "Spotify Tray")
}
func SetOpenAtLogin(value bool) {
	if value {
		exec.Command("osascript", "-e", "tell application \"System Events\" to make login item at end with properties {name: \"Spotify Tray\",path:\""+getAppPath()+"\", hidden:false}").Run()
	} else {
		exec.Command("osascript", "-e", "tell application \"System Events\" to delete login item \"Spotify Tray\"").Run()
	}
}

func getAppPath() string {
	executable, _ := os.Executable()

	return filepath.Join(filepath.Dir(executable), "../../")
}

func readFile() {
	if len(path) != 0 {
		content, err := os.ReadFile(path)
		if err != nil {
			return
		}

		preferences = Preferences{}
		_ = json.Unmarshal(content, &preferences)
	}
}

func writeFile() {
	if len(path) != 0 {
		content, err := json.Marshal(&preferences)
		if err != nil {
			return
		}

		err = os.WriteFile(path, content, 0644)

		if os.IsNotExist(err) {
			dirname, _ := os.UserConfigDir()

			os.Mkdir(filepath.Join(dirname, "/spotify-tray"), 0755)
			os.WriteFile(path, content, 0644)
		}
	}
}
