package main

import (
	"github.com/numbermess/environment-preloader/pkg/preloader"
	"log"
)

// Main entrypoint. Doesn't really do anything but make IDEA forget the unused method reference.
func main() {
	if err := preloader.PreloadEnvironment(); err != nil {
		log.Printf("%v", err)
	}
}
