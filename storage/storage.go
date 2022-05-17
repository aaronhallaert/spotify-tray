package storage

import (
	"fmt"
	"os"
	"sync"
)

type preferences struct {
	hasProgress bool
	artistFirst bool
	moreSpace   bool
}

var preferencesInstance *preferences

var path = ""
var lock = &sync.Mutex{}

func GetPreferencesInstance() *preferences {
	if preferencesInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if preferencesInstance == nil {
            // create first preferences instance
            preferencesInstance = &preferences{}

            // read preferences from file or write defaults
			dirname, err := os.UserConfigDir()
			if err != nil {
                preferencesInstance.hasProgress = true
                preferencesInstance.artistFirst = true
                preferencesInstance.moreSpace = true
                return preferencesInstance
			}

			path = dirname + "/spotify-tray/preferences.json"
            fmt.Println(path)
			readOrWriteFileIfExist(path)
		}
	}

	return preferencesInstance
}

// ---- GETTERS / SETTERS ---- //

func (prefs *preferences) GetHasProgress() bool {
	return prefs.hasProgress
}
func (prefs *preferences) SetHasProgress(value bool) {
	prefs.hasProgress = value
	writeFile(prefs)
}

func (prefs *preferences) GetArtistFirst() bool {
	return prefs.artistFirst
}
func (prefs *preferences) SetArtistFirst(value bool) {
	prefs.artistFirst = value
	writeFile(prefs)
}

func (prefs *preferences) GetMoreSpace() bool {
	return prefs.moreSpace
}
func (prefs *preferences) SetMoreSpace(value bool) {
	prefs.moreSpace = value
	writeFile(prefs)
}
