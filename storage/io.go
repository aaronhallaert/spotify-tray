package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// --- Marshalling and Unmarshalling --- //

type JsonPrefs struct {
	HasProgress bool
	ArtistFirst bool
	MoreSpace   bool
}


func (prefs *preferences) UnmarshalJsonPrefs(b []byte) {
    fmt.Println("Unmarshalling preferences")
    jsonPrefs := JsonPrefs{}

    fmt.Println("JSON to struct")
    _ = json.Unmarshal(b, &jsonPrefs)

    prefs.hasProgress = jsonPrefs.HasProgress
    prefs.artistFirst = jsonPrefs.ArtistFirst
    prefs.moreSpace = jsonPrefs.MoreSpace
}


func (prefs *preferences) MarshalJSON() ([]byte, error) {
    return json.Marshal(JsonPrefs{
        HasProgress: prefs.hasProgress,
        ArtistFirst: prefs.artistFirst,
        MoreSpace:   prefs.moreSpace,
    })
}


// ---- READ / WRITE PREFERENCES FILE ---- //

func readOrWriteFileIfExist(fileName string) {
	_, error := os.Stat(fileName)

	if os.IsNotExist(error) {
		preferencesInstance.hasProgress = true
        preferencesInstance.artistFirst = true
        preferencesInstance.moreSpace = true
		writeFile(preferencesInstance)
	} else {
		readFile()
	}
}

func readFile() {
	if len(path) != 0 {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Println(err)
			return
		}

        preferencesInstance.UnmarshalJsonPrefs(content)
	}
}


func writeFile(prefs *preferences) {
    fmt.Println("Writing preferences to file: ", path)
	if len(path) != 0 {
		content, err := prefs.MarshalJSON()

		if err != nil {
            fmt.Println(err)
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
