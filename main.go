package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"math/big"
	"strings"
)

func generatePassword(length int, exclude string) (string, error) {
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

func main() {

	passLength := flag.Int("length", 16, "Length of password")
	exclude := flag.String("exclude", "", "Exclude symbols of password")
	flag.Parse()
	if *passLength <= 0 {
		fmt.Println("Длина должна быть больше нуля")
		return
	}
	fmt.Println("Длина пароля: ", *passLength)

	pass, err := generatePassword(*passLength, *exclude)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(pass)
}
