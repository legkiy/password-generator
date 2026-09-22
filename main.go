package main

import (
	"fmt"
	"os"

	"github.com/legkiy/password-generator-cli/internal/cli"
	"github.com/legkiy/password-generator-cli/internal/settings"
)

func main() {
	configPath, err := settings.UserPath()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Не удалось определить путь к настройкам:", err)
		os.Exit(1)
	}

	app := cli.New(settings.New(configPath))

	if err := app.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка работы приложения:", err)
		os.Exit(1)
	}
}
