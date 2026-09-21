package password

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

func Generate(length int, exclude string) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("длина должна быть больше нуля")

	}
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Проверяем какие символы разрешены
	allowed := make([]byte, 0, len(alphabet))
	for idx, value := range alphabet {
		if !strings.Contains(exclude, string(value)) {
			allowed = append(allowed, alphabet[idx])
		}
	}
	if len(allowed) == 0 {
		return "", fmt.Errorf("все доступные символы исключены")
	}

	// Генерируем пароль
	bytes := make([]byte, length)
	for i := range bytes {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(allowed))))
		if err != nil {
			return "", err
		}
		bytes[i] = allowed[index.Int64()]
	}
	return string(bytes), nil
}
