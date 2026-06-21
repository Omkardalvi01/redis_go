package engine

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExecuteLineBlankInput(t *testing.T) {
	eng, err := New(Config{EnableAOF: false})
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}

	result := eng.ExecuteLine("   ", false)
	if result.Err != ErrEmptyInput {
		t.Fatalf("expected ErrEmptyInput, got %v", result.Err)
	}
}

func TestSetGetDelRenameAndKeys(t *testing.T) {
	eng, err := New(Config{EnableAOF: false})
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}

	checkOK(t, eng.ExecuteLine("set alpha one", false))
	checkOK(t, eng.ExecuteLine("set beta two", false))

	if got := eng.ExecuteLine("get alpha", false).Response; got != "one\n" {
		t.Fatalf("get alpha = %q", got)
	}

	if got := eng.ExecuteLine("keys alpha", false).Response; got != "0) alpha\n" {
		t.Fatalf("exact keys should return the key, got %q", got)
	}

	if got := eng.ExecuteLine("keys b*", false).Response; got != "0) beta\n" {
		t.Fatalf("glob keys should return beta, got %q", got)
	}

	checkOK(t, eng.ExecuteLine("rename alpha gamma", false))
	if got := eng.ExecuteLine("get gamma", false).Response; got != "one\n" {
		t.Fatalf("rename did not move value, got %q", got)
	}

	if got := eng.ExecuteLine("del beta gamma", false).Response; got != "(integer) 2\n" {
		t.Fatalf("del count mismatch: %q", got)
	}
}

func TestReplayMissingLogIsNotAnError(t *testing.T) {
	dir := t.TempDir()
	eng, err := New(Config{LogPath: filepath.Join(dir, "log.txt"), EnableAOF: true})
	if err != nil {
		t.Fatalf("new engine with missing log: %v", err)
	}

	checkOK(t, eng.ExecuteLine("set a 1", true))

	data, err := os.ReadFile(filepath.Join(dir, "log.txt"))
	if err != nil {
		t.Fatalf("read log: %v", err)
	}
	if !strings.Contains(string(data), "set a 1") {
		t.Fatalf("log did not record command: %q", string(data))
	}
}

func checkOK(t *testing.T, result Result) {
	t.Helper()
	if result.Err != nil {
		t.Fatalf("unexpected error: %v", result.Err)
	}
}
