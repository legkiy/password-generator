package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"pass-generator/internal/password"
	"pass-generator/internal/settings"
	"runtime"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
)

type App struct {
	reader *bufio.Reader
	store  settings.Store
}

func New(store settings.Store) *App {
	return &App{
		reader: bufio.NewReader(os.Stdin),
		store:  store,
	}
}

func (app *App) Run() error {
	config := password.Config{
		Length:  16,
		Exclude: "",
	}

	loadedConfig, err := app.store.Load()
	message := ""

	if err == nil {
		config = loadedConfig
	} else if !errors.Is(err, os.ErrNotExist) {
		message = fmt.Sprint("Не удалось загрузить настройки, используются значения по умолчанию: ", err)
	}

	for {
		err := clearScreen()
		if err != nil {
			fmt.Println("Не удалось очистить экран:", err)
			return err
		}
		if message != "" {
			fmt.Println(message)
			message = ""
		}
		fmt.Printf(`
1: Создать пароль
2: Задать длину пароля (%d)
3: Символы исключения (%s)

0: Закрыть приложение

Ваш выбор: `, config.Length, config.Exclude)

		input, err := app.reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return err
		}
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			pass, err := password.Generate(config.Length, config.Exclude)
			if err != nil {
				message = fmt.Sprint("Не удалось создать пароль: ", err)
				continue
			}

			shouldExit, err := app.passwordMenu(pass)
			if err != nil {
				fmt.Println(err)
				return err
			}
			if shouldExit {
				return nil
			}
		case "2":
			err := clearScreen()
			if err != nil {
				fmt.Println("Не удалось очистить экран:", err)
				return err
			}
			fmt.Print("Введите длину пароля (0 — закрыть приложение):")

			input, err := app.reader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return err
			}
			input = strings.TrimSpace(input)
			if input == "0" {
				return nil
			}
			newLength, err := strconv.Atoi(input)
			if err != nil || newLength <= 0 {
				message = "Введите целое число больше нуля"
				continue
			}
			config.Length = newLength
			err = app.store.Save(config)
			if err != nil {
				message = fmt.Sprint("Настройки изменены только для текущего запуска: ", err)
			}

		case "3":
			err := clearScreen()
			if err != nil {
				fmt.Println("Не удалось очистить экран:", err)
				return err
			}
			fmt.Println(`Введите исключаемые символы, например 0O1lI.
Enter — очистить исключения.
0 — закрыть приложение.
Для исключения только нуля введите "0" с кавычками.`)

			input, err := app.reader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return err
			}
			input = strings.TrimSpace(input)
			if input == "0" {
				return nil
			}
			if input == `"0"` {
				input = "0"
			}
			config.Exclude = input
			err = app.store.Save(config)
			if err != nil {
				message = fmt.Sprint("Настройки изменены только для текущего запуска: ", err)
			}

		case "0":
			return nil
		default:
			message = "Неизвестный пункт меню"
		}

	}

}

/*
true, nil — пользователь хочет закрыть приложение;
false, nil — вернуться в главное меню;
false, err — произошла ошибка чтения.
*/
func (app *App) passwordMenu(password string) (bool, error) {
	message := ""
	for {
		err := clearScreen()
		if err != nil {
			return false, err
		}
		fmt.Println("Пароль: ", password)
		if message != "" {
			fmt.Println(message)
			message = ""
		}
		fmt.Print(`
1: Скопировать в буфер обмена
2: Сохранить в файл
3: Вернуться в главное меню

0: Закрыть приложение

Ваш выбор: `)

		input, err := app.reader.ReadString('\n')
		if err != nil {
			return false, err
		}
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			err := copyPassword(password)
			if err != nil {
				message = fmt.Sprint("Не удалось скопировать пароль: ", err)
				continue
			}
			message = "Пароль скопирован"

		case "2":
			err := clearScreen()
			if err != nil {
				fmt.Println("Не удалось очистить экран:", err)
				return false, err
			}
			fmt.Println(`Введите путь к новому файлу, например password.txt.
Enter — вернуться в подменю.
0 — закрыть приложение.`)

			input, err := app.reader.ReadString('\n')
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
				message = fmt.Sprint("Не удалось сохранить пароль: ", err)
				continue
			}
			message = "Пароль сохранён в файл: " + input
		case "3":
			return false, nil
		case "0":
			return true, nil
		default:
			message = "Неизвестный пункт меню"
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

func clearScreen() error {
	goos := runtime.GOOS
	if goos == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		return cmd.Run()

	} else {
		_, err := fmt.Print("\x1b[2J\x1b[H")
		return err
	}
}
