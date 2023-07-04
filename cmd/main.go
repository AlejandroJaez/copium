package main

import (
	"errors"
	cp "github.com/otiai10/copy"
	"github.com/radovskyb/watcher"
	"log"
	"os"
	"os/user"
	"regexp"
	"time"
)

func main() {

	if _, err := os.Stat("./test_folder"); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir("./test_folder", os.ModePerm)
		if err != nil {
			log.Println(err)
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

	// Only notify rename and move events.
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
					log.Println(event) // Print the event's info.

					if err := cp.Copy(desktop, "./test_folder/desktop"); err != nil {
						log.Println(err)
					}
				}

			case err := <-w.Error:
				log.Println(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch this folder for changes.

	// Watch test_folder recursively for changes.
	if err := w.AddRecursive("./test_folder"); err != nil {
		log.Fatalln(err)
	}

	// Watch test_folder recursively for changes.
	if err := w.AddRecursive(desktop); err != nil {
		log.Fatalln(err)
	}

	// Print a list of all of the files and folders currently
	// being watched and their paths.
	//for path, f := range w.WatchedFiles() {
	//	fmt.Printf("%s: %s\n", path, f.Name())
	//}

	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}
