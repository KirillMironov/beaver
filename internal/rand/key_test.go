package rand

import (
	"bytes"
	"testing"
)

func TestKey(t *testing.T) {
	t.Parallel()

	const length = 32

	key, err := Key(length)
	if err != nil {
		t.Fatal(err)
	}

	if len(key) != length {
		t.Fatalf("expected key length to be %d, got %d", length, len(key))
	}

	key2, err := Key(length)
	if err != nil {
		t.Fatal(err)
	}

	if len(key2) != length {
		t.Fatalf("expected key2 length to be %d, got %d", length, len(key2))
	}

	if bytes.Equal(key, key2) {
		t.Fatalf("expected key to be different from key2, got %q, %q", key, key2)
	}
}
