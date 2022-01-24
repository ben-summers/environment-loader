package main

import (
	"github.com/rs/zerolog/log"
	"preloader/pkg/preloader"
)

// Main entrypoint. Doesn't really do anything but make IDEA forget the unused method reference.
func main() {
	if _, err := preloader.PreloadEnvironment(); err != nil {
		log.Err(err).
			Msg("cannot preload environment")
	}
}
