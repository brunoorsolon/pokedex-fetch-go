package state

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAtomicWriteCreatesParentDirectoryAndReplacesFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "state.json")

	if err := AtomicWrite(path, []byte(`{"version":1}`)); err != nil {
		t.Fatalf("AtomicWrite() initial error = %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() initial error = %v", err)
	}
	if string(data) != `{"version":1}` {
		t.Fatalf("initial file content = %q", string(data))
	}

	if err := AtomicWrite(path, []byte(`{"version":2}`)); err != nil {
		t.Fatalf("AtomicWrite() replace error = %v", err)
	}
	data, err = os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() replaced error = %v", err)
	}
	if string(data) != `{"version":2}` {
		t.Fatalf("replaced file content = %q", string(data))
	}
}
