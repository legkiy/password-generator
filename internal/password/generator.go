package password

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-=_+;:,.?"

var errAllExcluded = errors.New("все доступные символы исключены")

func allowedChars(exclude string) []byte {
	allowed := make([]byte, 0, len(alphabet))
	for i := 0; i < len(alphabet); i++ {
		if !strings.ContainsRune(exclude, rune(alphabet[i])) {
			allowed = append(allowed, alphabet[i])
		}
	}
	return allowed
}

// ValidateExclude позволяет отклонить исключения до того, как они попадут в настройки.
func ValidateExclude(exclude string) error {
	if len(allowedChars(exclude)) == 0 {
		return errAllExcluded
	}
	return nil
}

func Generate(length int, exclude string) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("длина должна быть больше нуля")
	}

	allowed := allowedChars(exclude)
	if len(allowed) == 0 {
		return "", errAllExcluded
	}

	pass := make([]byte, length)
	for i := range pass {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(allowed))))
		if err != nil {
			return "", err
		}
		pass[i] = allowed[index.Int64()]
	}
	return string(pass), nil
}
