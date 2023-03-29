package server

import (
	"testing"

	"go.uber.org/fx"
)

func TestModule(t *testing.T) {
	t.Parallel()

	if err := fx.ValidateApp(Module); err != nil {
		t.Fatal(err)
	}
}
