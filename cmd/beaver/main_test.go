package main

import (
	"testing"

	"go.uber.org/fx"
)

func TestOptions(t *testing.T) {
	if err := fx.ValidateApp(options()); err != nil {
		t.Fatal(err)
	}
}
