//http://openmymind.net/Golang-Hot-Configuration-Reload/
//altered to use fsnotify for cross platform compat

package conf

import (
	"encoding/json"
	"github.com/howeyc/fsnotify"
	"io/ioutil"
	"log"
	"os"
	//"os/signal"
	"sync"
	//"syscall"
)

var (
	config       *Config
	configLock   = new(sync.RWMutex)
	watcher, err = fsnotify.NewWatcher()
)

func loadConfig(fail bool) {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Println("open config: ", err)
		if fail {
			os.Exit(1)
		}
	}

	temp := new(Config)
	if err = json.Unmarshal(file, temp); err != nil {
		log.Println("parse config: ", err)
		if fail {
			os.Exit(1)
		}
	}
	configLock.Lock()
	config = temp
	configLock.Unlock()
}

func GetConfig() *Config {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

// go calls init on start
func init() {
	loadConfig(true)
	//s := make(chan os.Signal, 1)
	//signal.Notify(s, syscall.SIGUSR2)

	go func() {
		for {
			//<-s
			select {
			case ev := <-watcher.Event:
				log.Println("ev: ", ev)
				loadConfig(false)
				log.Println("Reloaded")
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Watch("config.json")

	if err != nil {
		log.Fatal(err)
	}
}
