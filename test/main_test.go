package test

import (
	"preloader/pkg/preloader"
	"testing"
)

func TestPreloadSilent(t *testing.T) {

	preloader.OverrideVerbosity = true
	preloader.OverriddenVerbosity = false

	_, err := preloader.PreloadEnvironment()
	if err != nil {
		t.Fail()
	}

}

func TestPreloadVerbose(t *testing.T) {

	preloader.OverrideVerbosity = true
	preloader.OverriddenVerbosity = true

	_, err := preloader.PreloadEnvironment()
	if err != nil {
		t.Fail()
	}

}
