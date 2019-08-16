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
		if _, stat := os.Stat(envFile); !os.IsNotExist(stat) {
			if file, err := os.Open(envFile); err != nil {
				return errors.New(fmt.Sprintf("could not open environment loader file: %s, %v", envFile, err))
			} else {
				scanner := bufio.NewScanner(file)
				firstScan := true
				keepScanning := false
				var environmentPreloaderPrefix string
				for scanner.Scan() {
					if firstScan {
						line := scanner.Text()
						chunks := strings.SplitN(line, "=", 2)
						if len(chunks) == 2 {
							if EnvironmentPreloaderPrefix == chunks[0] {
								environmentPreloaderPrefix = chunks[1]
							}
							if environmentPreloaderPrefix != "" {
								log.Printf("Found environment preloader prefix: %s", environmentPreloaderPrefix)
								keepScanning = true
							}
						}
						firstScan = false
						if keepScanning {
							log.Printf("I will continue to load environment variables from this file.")
						}
					} else {
						line := scanner.Text()
						chunks := strings.SplitN(line, "=", 2)
						if len(chunks) == 2 {
							prefixedName := fmt.Sprintf("%s_%s", environmentPreloaderPrefix, chunks[0])

							log.Printf("Setting %s...", prefixedName)

							if setenvErr := os.Setenv(prefixedName, chunks[1]); setenvErr != nil {
								return errors.New(fmt.Sprintf("error setting environment variable: %s => %v", chunks[0], setenvErr))
							}
						}
					}

				}
				if err := file.Close(); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
