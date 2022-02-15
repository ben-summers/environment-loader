package test

import (
	"github.com/numbermess/environment-preloader/pkg/preloader"
	"github.com/rs/zerolog/log"
	"testing"
)

func TestPreloadDebug(t *testing.T) {

	preloader.OverrideLogLevel = true
	preloader.OverriddenLogLevel = "debug"

	_, err := preloader.PreloadEnvironment()
	if err != nil {
		t.Fail()
	}

}

func TestPreloadWarning(t *testing.T) {

	preloader.OverrideLogLevel = true
	preloader.OverriddenLogLevel = "warn"

	_, err := preloader.PreloadEnvironment()
	if err != nil {
		t.Fail()
	}

}

func TestPreloadInfo(t *testing.T) {

	preloader.OverrideLogLevel = true
	preloader.OverriddenLogLevel = "info"

	_, err := preloader.PreloadEnvironment()
	if err != nil {
		t.Fail()
	}

}

func TestPreloadBlob(t *testing.T) {
	blob, err := preloader.PreloadEnvironmentBlob()
	if err != nil {
		t.Fail()
	}

	log.Debug().
		Str("blob", blob).
		Msg("exported blob")

}
