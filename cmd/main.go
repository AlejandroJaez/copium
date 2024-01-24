package main

import (
	"errors"
	cp "github.com/otiai10/copy"
	"github.com/radovskyb/watcher"
	"os"
	"os/user"
	"regexp"
	"time"
)

func main() {

	if _, err := os.Stat("./test_folder"); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir("./test_folder", os.ModePerm)
		if err != nil {
			println(err)
		}
	}

	myself, err := user.Current()
	if err != nil {
		panic(err)
	}
	homedir := myself.HomeDir
	desktop := homedir + "/Desktop"

	w := watcher.New()

	// SetMaxEvents to 1 to allow at most 1 event's to be received
	// on the Event channel per watching cycle.
	//
	// If SetMaxEvents is not set, the default is to send all events.
	//w.SetMaxEvents(1)

	// Only notify these events.
	w.FilterOps(watcher.Rename, watcher.Move, watcher.Create, watcher.Write)

	// Only files that match the regular expression during file listings
	// will be watched.
	r := regexp.MustCompile(".*")
	w.AddFilterHook(watcher.RegexFilterHook(r, false))

	go func() {
		for {
			select {
			case event := <-w.Event:
				if event.Op == 0 {

					if err := cp.Copy(desktop, "./test_folder/desktop"); err != nil {
						println(err)
					}
				}

			case err := <-w.Error:
				println(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch test_folder recursively for changes.
	if err := w.AddRecursive(desktop); err != nil {
		println(err)
	}

	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Millisecond * 100); err != nil {
		println(err)
	}
}
