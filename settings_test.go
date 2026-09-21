package main

import (
	"errors"
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
		t.Fatalf("загрузить настройки: %q", err)

	}
	if settings.Length != 24 {
		t.Errorf("длина: ожидалось 24, получено %d", settings.Length)
	}
	if settings.Exclude != "0O1lI" {
		t.Errorf(`исключения: ожидалось "0O1lI", получено %s`, settings.Exclude)
	}
}

func TestLoadSettingsMissingFile(t *testing.T) {
	tempDir := t.TempDir()
	file := filepath.Join(tempDir, "settings.json")
	_, err := loadSettings(file)
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("ожидалась ошибка отсутствия файла, получено: %v", err)
	}
}

func TestLoadSettingsInvalidJSON(t *testing.T) {
	tempDir := t.TempDir()
	file := filepath.Join(tempDir, "settings.json")
	err := os.WriteFile(file, []byte(`{"length":`), 0600)
	if err != nil {
		t.Fatalf("не удалось подготовить файл: %v", err)
	}
	_, err = loadSettings(file)
	if err == nil {
		t.Fatal("ожидалась ошибка разбора JSON, получено nil")
	}
}

func TestLoadSettingsInvalidLength(t *testing.T) {
	tempDir := t.TempDir()
	file := filepath.Join(tempDir, "settings.json")
	err := os.WriteFile(file, []byte(`{"length":0,"exclude":""}`), 0600)
	if err != nil {
		t.Fatalf("не удалось подготовить файл: %v", err)
	}
	_, err = loadSettings(file)
	if err == nil {
		t.Fatal("ожидалась ошибка для нулевой длины, получено nil")
	}
}

func TestLoadSettingsNegativeLength(t *testing.T) {
	tempDir := t.TempDir()
	file := filepath.Join(tempDir, "settings.json")
	err := os.WriteFile(file, []byte(`{"length":-5,"exclude":""}`), 0600)
	if err != nil {
		t.Fatalf("не удалось подготовить файл: %v", err)
	}
	_, err = loadSettings(file)
	if err == nil {
		t.Fatal("ожидалась ошибка для отрецательной длины, получено nil")
	}
}
