package preloader

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	EnvironmentVariableName    = "ENVIRONMENT_PRELOADER_FILE"
	EnvironmentPreloaderPrefix = "ENVIRONMENT_PRELOADER_PREFIX"
)

// PreloadEnvironment looks for a file specified by an environment variable named ENVIRONMENT_PRELOADER
func PreloadEnvironment() error {
	envFile := os.Getenv(EnvironmentVariableName)
	if envFile != "" {
		log.Printf("Loading environment preload file from: %s", envFile)
		if _, stat := os.Stat(envFile); !os.IsNotExist(stat) {
			if file, err := os.Open(envFile); err != nil {
				return errors.New(fmt.Sprintf("Could not open environment loader file: %s, %v", envFile, err))
			} else {
				scanner := bufio.NewScanner(file)
				firstScan := true
				var environmentPreloaderPrefix string

			scanning:
				for scanner.Scan() {
					if firstScan {
						log.Printf("Scanning first line. It should look something like %s=%s", EnvironmentPreloaderPrefix, "SOME_PREFIX")
						line := scanner.Text()
						chunks := strings.SplitN(line, "=", 2)
						if len(chunks) == 2 {
							if EnvironmentPreloaderPrefix == chunks[0] {
								environmentPreloaderPrefix = chunks[1]
							} else {
								log.Printf("Did not find %s on the first line of %s. That's OK, but you might want to check on this in a local development environment.", EnvironmentPreloaderPrefix, envFile)
								log.Println("Exiting environment preload script.")
								break scanning
							}
							if environmentPreloaderPrefix != "" {
								log.Printf("Found environment preloader prefix: %s", environmentPreloaderPrefix)
							}
						}
						firstScan = false
					} else {
						line := scanner.Text()
						chunks := strings.SplitN(line, "=", 2)
						if len(chunks) == 2 {
							prefixedName := fmt.Sprintf("%s_%s", environmentPreloaderPrefix, chunks[0])

							log.Printf("Setting %s", prefixedName)

							if setenvErr := os.Setenv(prefixedName, chunks[1]); setenvErr != nil {
								return errors.New(fmt.Sprintf("error setting environment variable: %s => %v", chunks[0], setenvErr))
							}
						}
					}

				}
				if err := file.Close(); err != nil {
					return err
				}
				log.Println("All done setting prefixed environment variables from your file. Have a nice day!")
			}
		}
	}
	return nil
}
