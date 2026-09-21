package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
)

type Config struct {
	Length  int    `json:"length"`
	Exclude string `json:"exclude"`
}

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

/*
true, nil — пользователь хочет закрыть приложение;
false, nil — вернуться в главное меню;
false, err — произошла ошибка чтения.
*/
func passwordMenu(reader *bufio.Reader, password string) (bool, error) {
	fmt.Println(password)
	for {
		fmt.Printf(`
1: Скопировать в буфер обмена
2: Сохранить в файл
3: Вернуться в главное меню

0: Закрыть приложение

Ваш выбор: `)

		input, err := reader.ReadString('\n')
		if err != nil {
			return false, err
		}
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			err := copyPassword(password)
			if err != nil {
				fmt.Println("Не удалось скопировать пароль. ", err)
				continue
			}
			fmt.Println("Пароль скопирован")

		case "2":
			fmt.Println(`Введите путь к новому файлу, например password.txt.
Enter — вернуться в подменю.
0 — закрыть приложение.`)

			input, err := reader.ReadString('\n')
			if err != nil {
				return false, err
			}

			input = strings.TrimSpace(input)
			if input == "0" {
				return true, nil
			}
			if len(input) == 0 {
				continue
			}
			err = savePassword(input, password)
			if err != nil {
				fmt.Println("Не удалось сохранить пароль", err)
				continue
			}
			fmt.Println("Пароль сохранён в файл: ", input)
		case "3":
			return false, nil
		case "0":
			return true, nil
		default:
			fmt.Println("Неизвестный пункт меню")
		}
	}
}

func copyPassword(password string) error {
	err := clipboard.WriteAll(password)
	return err
}

func savePassword(path string, password string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err

	}
	defer file.Close()

	_, err = file.WriteString(password)
	if err != nil {
		return err
	}
	return file.Close()
}

func main() {
	config := Config{
		Length:  16,
		Exclude: "",
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf(`
1: Создать пароль
2: Задать длину пароля (%d)
3: Символы исключения (%s)

0: Закрыть приложение

Ваш выбор: `, config.Length, config.Exclude)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			password, err := generatePassword(config.Length, config.Exclude)
			if err != nil {
				fmt.Println(err)
				continue
			}

			shouldExit, err := passwordMenu(reader, password)
			if err != nil {
				fmt.Println(err)
				return
			}
			if shouldExit {
				return
			}
		case "2":
			fmt.Print("Введите длину пароля (0 — закрыть приложение):")

			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return
			}
			input = strings.TrimSpace(input)
			if input == "0" {
				return
			}
			newLength, err := strconv.Atoi(input)
			if err != nil || newLength <= 0 {
				fmt.Println("Введи корректное число > 0")
				continue
			}
			config.Length = newLength
		case "3":
			fmt.Println(`Введите исключаемые символы, например 0O1lI.
Enter — очистить исключения.
0 — закрыть приложение.
Для исключения только нуля введите "0" с кавычками.`)

			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return
			}
			input = strings.TrimSpace(input)
			if input == "0" {
				return
			}
			if input == `"0"` {
				input = "0"
			}
			config.Exclude = input

		case "0":
			return
		default:
			fmt.Println("Неизвестный пункт меню")
		}

	}
}
