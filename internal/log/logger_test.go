package log

import "testing"

func TestNew(t *testing.T) {
	t.Parallel()

	logger, err := New()
	if err != nil {
		t.Fatal(err)
	}

	if logger == nil {
		t.Fatal("expected logger to be not nil")
	}
}
