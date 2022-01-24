package test

import (
	"preloader/pkg/preloader"
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
