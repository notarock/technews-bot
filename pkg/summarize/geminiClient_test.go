package summarize

import (
	"os"
	"testing"
)

func TestModelName(t *testing.T) {
	t.Setenv("GEMINI_MODEL", "")
	if got := modelName(); got != "gemini-2.5-flash" {
		t.Fatalf("empty model: got %q", got)
	}

	if err := os.Unsetenv("GEMINI_MODEL"); err != nil {
		t.Fatal(err)
	}
	if got := modelName(); got != "gemini-2.5-flash" {
		t.Fatalf("unset model: got %q", got)
	}

	t.Setenv("GEMINI_MODEL", "gemini-2.5-pro")
	if got := modelName(); got != "gemini-2.5-pro" {
		t.Fatalf("configured model: got %q", got)
	}
}
