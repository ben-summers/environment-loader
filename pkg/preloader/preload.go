package preloader

import (
	"bufio"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

const (
	EnvironmentVariableName = "ENVIRONMENT_PRELOADER_FILE"
	Prefix                  = "ENVIRONMENT_PRELOADER_PREFIX"
	LogLevel                = "ENVIRONMENT_PRELOADER_LOG_LEVEL"
)

const (
	Error   = "error"
	Warning = "warning"
	Warn    = "warn"
	Info    = "info"
	Debug   = "debug"
)

var (
	OverriddenLogLevel string
	OverrideLogLevel   bool
)

// PreloadEnvironment looks for a file specified by an environment variable named ENVIRONMENT_PRELOADER
func PreloadEnvironment() (string, error) {
	var environmentPreloaderPrefix string

	envFile := os.Getenv(EnvironmentVariableName)
	if envFile != "" {
		log.Debug().
			Str("filename", envFile).
			Msg("opening environment file")

		if _, stat := os.Stat(envFile); !os.IsNotExist(stat) {
			if file, err := os.Open(envFile); err != nil {
				log.Err(err).Msg("could not open file")
				return "", err
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

				if prefix, ok := env[Prefix]; !ok {
					log.Warn().
						Msg("no environment preloader prefix specified; environment variables named in this file will be exported as-is")
				} else {
					environmentPreloaderPrefix = prefix
				}

				logLevel, ok := env[LogLevel]
				if !ok {
					logLevel = Error
					zerolog.SetGlobalLevel(zerolog.ErrorLevel)
				}

				delete(env, Prefix)
				delete(env, LogLevel)

				if OverrideLogLevel {
					logLevel = OverriddenLogLevel
				}

				switch logLevel {
				case Error:
					zerolog.SetGlobalLevel(zerolog.ErrorLevel)
				case Warning:
					fallthrough
				case Warn:
					zerolog.SetGlobalLevel(zerolog.WarnLevel)
				case Info:
					zerolog.SetGlobalLevel(zerolog.InfoLevel)
				case Debug:
					zerolog.SetGlobalLevel(zerolog.DebugLevel)
				}

				for key, value := range env {
					prefixedName := fmt.Sprintf("%s_%s", environmentPreloaderPrefix, key)

					log.Info().
						Str("key", prefixedName).
						Msg("setting variable")
					log.Debug().
						Str("key", prefixedName).
						Str("value", value).
						Msg("setting variable")

					if setenvErr := os.Setenv(prefixedName, value); setenvErr != nil {
						log.Error().
							Err(setenvErr)
						return "", setenvErr
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
