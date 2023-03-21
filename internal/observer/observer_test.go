package observer

import "testing"

func TestObserver_First(t *testing.T) {
	t.Parallel()

	observer := New()

	if got, want := observer.First(), ""; got != want {
		t.Fatalf("First() = %q, want %q", got, want)
	}

	observer.Info("Hello!")

	if got, want := observer.First(), "Hello!"; got != want {
		t.Fatalf("First() = %q, want %q", got, want)
	}
}

func TestObserver_Len(t *testing.T) {
	t.Parallel()

	observer := New()

	if got, want := observer.Len(), 0; got != want {
		t.Fatalf("Len() = %d, want %d", got, want)
	}

	observer.Info("Info")
	observer.Infof("Infof")
	observer.Error("Error")
	observer.Errorf("Errorf")

	if got, want := observer.Len(), 4; got != want {
		t.Fatalf("Len() = %d, want %d", got, want)
	}
}
