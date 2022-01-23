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
	EnvironmentVariableName     = "ENVIRONMENT_PRELOADER_FILE"
	EnvironmentPreloaderPrefix  = "ENVIRONMENT_PRELOADER_PREFIX"
	EnvironmentPreloaderVerbose = "ENVIRONMENT_PRELOADER_VERBOSE"
)

// PreloadEnvironment looks for a file specified by an environment variable named ENVIRONMENT_PRELOADER
func PreloadEnvironment() (string, error) {
	var environmentPreloaderPrefix string
	verbose := false

	envFile := os.Getenv(EnvironmentVariableName)
	if envFile != "" {
		log.Printf("loading environment file: %s", envFile)
		if _, stat := os.Stat(envFile); !os.IsNotExist(stat) {
			if file, err := os.Open(envFile); err != nil {
				return "", errors.New(fmt.Sprintf("Could not open environment loader file: %s, %v", envFile, err))
			} else {
				scanner := bufio.NewScanner(file)

				// Try to find lines relevant to the preloader first.
				for scanner.Scan() {
					line := scanner.Text()
					chunks := strings.SplitN(line, "=", 2)
					if len(chunks) == 2 {
						if EnvironmentPreloaderPrefix == chunks[0] {
							environmentPreloaderPrefix = chunks[1]
						} else if EnvironmentPreloaderVerbose == chunks[0] {
							value := strings.ToLower(chunks[1])
							if "yes" == value || "true" == value {
								verbose = true
							}
						}
					}
				}

				for scanner.Scan() {
					line := scanner.Text()
					chunks := strings.SplitN(line, "=", 2)
					if len(chunks) == 2 {

						chunks[0] = strings.TrimSpace(chunks[0])
						chunks[1] = strings.TrimSpace(chunks[1])

						if chunks[0] == EnvironmentPreloaderVerbose || chunks[0] == EnvironmentPreloaderPrefix {
							continue
						}

						prefixedName := fmt.Sprintf("%s_%s", environmentPreloaderPrefix, chunks[0])

						if verbose {
							log.Printf("Setting %s", prefixedName)
						}

						if setenvErr := os.Setenv(prefixedName, chunks[1]); setenvErr != nil {
							return "", errors.New(fmt.Sprintf("error setting environment variable: %s => %v", chunks[0], setenvErr))
						}
					}
				}
				if err := file.Close(); err != nil {
					return "", err
				}
				if verbose {
					log.Println("All done setting prefixed environment variables from your file. Have a nice day!")
				}
			}
		}
	}
	return environmentPreloaderPrefix, nil
}
