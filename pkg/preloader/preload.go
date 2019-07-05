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
	EnvironmentVariableName = "ENVIRONMENT_PRELOADER_FILE"
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
				for scanner.Scan() {
					line := scanner.Text()
					chunks := strings.SplitN(line, "=", 2)
					if len(chunks) == 2 {
						if setenvErr := os.Setenv(chunks[0], chunks[1]); setenvErr != nil {
							return errors.New(fmt.Sprintf("error setting environment variable: %s => %v", chunks[0], setenvErr))
						}
					}
				}
				if err := file.Close(); err != nil {
					return err
				}
			}
		}
	} else {
		log.Printf("INFO: Environment variable '%s' not specified.", EnvironmentVariableName)
	}
	return nil
}
