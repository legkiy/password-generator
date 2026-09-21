package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
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
		fmt.Printf(`1: Скопировать в буфер обмена
2: Сохранить в файл
3: Вернуться в главное меню

0: Закрыть приложение
`)

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return false, err
		}
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			fmt.Println("Скоро добавим")
		case "2":
			fmt.Println("Скоро добавим")
		case "3":
				fmt.Println("Скоро добавим")
		case "0":
			return true, nil
		default:
			fmt.Println("Неизвестный пункт меню")
		}
		}
	}
}

func main() {
	config := Config{
		Length:  16,
		Exclude: "",
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf(`1: Создать пароль
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
			fmt.Println(password)
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
