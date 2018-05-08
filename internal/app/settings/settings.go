package settings

import (
	"github.com/imdario/mergo"
	"os"
	"sync"
)

type Settings struct {
	RootDir string
}

var instance = &Settings{
	RootDir: os.Getenv("TENT_ROOT_DIR"),
}
var setMutex = sync.Mutex{}

func Set(s Settings) {
	setMutex.Lock()
	defer setMutex.Unlock()
	mergo.Merge(instance, s, mergo.WithOverride)
}

func Get() *Settings {
	return instance
}
