package password

import (
	"errors"
	"strings"
	"testing"
)

func TestGenerateLength(t *testing.T) {
	for _, length := range []int{1, 16, 128} {
		pass, err := Generate(length, "")
		if err != nil {
			t.Fatalf("длина %d: %v", length, err)
		}
		if len(pass) != length {
			t.Errorf("длина: ожидалось %d, получено %d", length, len(pass))
		}
	}
}

func TestGenerateUsesOnlyAlphabet(t *testing.T) {
	pass, err := Generate(512, "")
	if err != nil {
		t.Fatalf("создать пароль: %v", err)
	}
	for _, char := range pass {
		if !strings.ContainsRune(alphabet, char) {
			t.Fatalf("символа %q нет в алфавите", char)
		}
	}
}

func TestGenerateSkipsExcluded(t *testing.T) {
	const exclude = "0O1lI"
	pass, err := Generate(512, exclude)
	if err != nil {
		t.Fatalf("создать пароль: %v", err)
	}
	if idx := strings.IndexAny(pass, exclude); idx >= 0 {
		t.Errorf("исключённый символ %q попал в пароль", pass[idx])
	}
}

func TestGenerateInvalidLength(t *testing.T) {
	for _, length := range []int{0, -5} {
		if _, err := Generate(length, ""); err == nil {
			t.Errorf("длина %d: ожидалась ошибка, получено nil", length)
		}
	}
}

func TestGenerateAllExcluded(t *testing.T) {
	_, err := Generate(16, alphabet)
	if !errors.Is(err, errAllExcluded) {
		t.Errorf("ожидалась ошибка исключения всех символов, получено: %v", err)
	}
}

func TestGenerateDiffersBetweenCalls(t *testing.T) {
	first, err := Generate(32, "")
	if err != nil {
		t.Fatalf("создать пароль: %v", err)
	}
	second, err := Generate(32, "")
	if err != nil {
		t.Fatalf("создать пароль: %v", err)
	}
	if first == second {
		t.Error("два вызова подряд вернули одинаковый пароль")
	}
}

func TestValidateExclude(t *testing.T) {
	if err := ValidateExclude(""); err != nil {
		t.Errorf("пустые исключения: ожидалось nil, получено %v", err)
	}
	if err := ValidateExclude("0O1lI"); err != nil {
		t.Errorf("частичные исключения: ожидалось nil, получено %v", err)
	}
	if err := ValidateExclude(alphabet); !errors.Is(err, errAllExcluded) {
		t.Errorf("весь алфавит: ожидалась ошибка, получено %v", err)
	}
}
