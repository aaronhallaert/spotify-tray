package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Preferences struct {
	HasProgress bool
	ArtistFirst bool
	MoreSpace   bool
}

var path = ""
var preferences = Preferences{
	HasProgress: true,
	ArtistFirst: true,
	MoreSpace:   true,
}

func Init() {
	dirname, err := os.UserConfigDir()
	if err != nil {
		return
	}

	path = dirname + "/spotify-tray/preferences.json"
	fmt.Println(path)
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

func GetHasProgress() bool {
	return preferences.HasProgress
}
func SetHasProgress(value bool) {
	preferences.HasProgress = value
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

func readFile() {
	if len(path) != 0 {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Println(err)
			return
		}

		preferences = Preferences{}
		_ = json.Unmarshal(content, &preferences)
	}
}

func writeFile() {
	if len(path) != 0 {
		content, err := json.Marshal(preferences)
		if err != nil {
			return
		}

		err = ioutil.WriteFile(path, content, 0644)

		if os.IsNotExist(err) {
			dirname, _ := os.UserConfigDir()

			os.Mkdir(filepath.Join(dirname, "/spotify-tray"), 0755)
			ioutil.WriteFile(path, content, 0644)
		}
	}
}
