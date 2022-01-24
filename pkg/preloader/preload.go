package preloader

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

const (
	EnvironmentVariableName     = "ENVIRONMENT_PRELOADER_FILE"
	EnvironmentPreloaderPrefix  = "ENVIRONMENT_PRELOADER_PREFIX"
	EnvironmentPreloaderVerbose = "ENVIRONMENT_PRELOADER_VERBOSE"
)

var (
	OverriddenVerbosity bool
	OverrideVerbosity   bool
)

// PreloadEnvironment looks for a file specified by an environment variable named ENVIRONMENT_PRELOADER
func PreloadEnvironment() (string, error) {
	var environmentPreloaderPrefix string
	verbose := false

	envFile := os.Getenv(EnvironmentVariableName)
	if envFile != "" {
		log.Debug().
			Str("filename", envFile).
			Msg("opening environment file")

		if _, stat := os.Stat(envFile); !os.IsNotExist(stat) {
			if file, err := os.Open(envFile); err != nil {
				return "", errors.New(fmt.Sprintf("Could not open environment loader file: %s, %v", envFile, err))
			} else {
				scanner := bufio.NewScanner(file)
				env := make(map[string]string)

				// Try to find lines relevant to the preloader first.
				for scanner.Scan() {

					line := scanner.Text()
					chunks := strings.SplitN(line, "=", 2)
					if len(chunks) == 2 {

						chunks[0] = strings.TrimSpace(chunks[0])
						chunks[1] = strings.TrimSpace(chunks[1])

						env[chunks[0]] = chunks[1]

					}
				}

				if prefix, ok := env[EnvironmentPreloaderPrefix]; !ok {
					log.Warn().
						Msg("no environment preloader prefix specified; environment variables named in this file will be exported as-is")
				} else {
					environmentPreloaderPrefix = prefix
				}
				if shouldBeVerbose, ok := env[EnvironmentPreloaderVerbose]; !ok {
					verbose = false
				} else {
					shouldBeVerbose = strings.ToLower(shouldBeVerbose)
					var maybe bool
					maybe = "yes" == shouldBeVerbose
					maybe = maybe || "true" == shouldBeVerbose
					verbose = maybe
				}

				delete(env, EnvironmentPreloaderPrefix)
				delete(env, EnvironmentPreloaderVerbose)

				if OverrideVerbosity {
					verbose = OverriddenVerbosity
				}

				for key, value := range env {
					prefixedName := fmt.Sprintf("%s_%s", environmentPreloaderPrefix, key)

					if verbose {
						log.Info().
							Str("key", prefixedName).
							Str("value", value).
							Msg("setting environment variable")

					}

					if setenvErr := os.Setenv(prefixedName, value); setenvErr != nil {
						return "", errors.New(fmt.Sprintf("error setting environment variable: %s => %v", key, setenvErr))
					}
				}
				if err := file.Close(); err != nil {
					return "", err
				}

				log.Info().
					Msg("all prefixed environment variables have been set and exported.")
			}
		}
	}
	return environmentPreloaderPrefix, nil
}
