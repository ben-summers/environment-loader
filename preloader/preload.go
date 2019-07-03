package preloader

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func Preload() error {
	envFile := os.Getenv("ENVIRONMENT_LOADER")
	if envFile != "" {
		if _, stat := os.Stat(envFile); !os.IsNotExist(stat) {
			if file, err := os.Open(envFile); err != nil {
				return errors.New(fmt.Sprintf("Could not open environment loader file: %s, %v", envFile, err))
			} else {
				defer file.Close()
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					line := scanner.Text()
					chunks := strings.SplitN(line, "=", 2)
					if len(chunks) == 2 {
						if setenvErr := os.Setenv(chunks[0], chunks[1]); setenvErr != nil {
							return errors.New(fmt.Sprintf("Error setting environment variable: %s => %v", chunks[0], setenvErr))
						}
					}
				}
			}
		}
	}
	return nil
}
