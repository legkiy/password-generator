package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSettings(t *testing.T) {
	tempDir := t.TempDir()
	file := filepath.Join(tempDir, "settings.json")
	err := os.WriteFile(file, []byte(`{"length":24,"exclude":"0O1lI"}`), 0600)
	if err != nil {
		t.Fatalf("не удалось подготовить файл: %v", err)
	}
	settings, err := loadSettings(file)
	if err != nil {
		t.Fatalf("не удалось подготовить файл: %v", err)

	}
	if settings.Length != 24 {
		t.Errorf("длина: ожидалось 24, получено %d", settings.Length)
	}
	if settings.Exclude != "0O1lI" {
		t.Errorf(`исключения: ожидалось "0O1lI", получено %s`, settings.Exclude)
	}
}
