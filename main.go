package main

import (
	"fmt"
	"pass-generator/internal/cli"
	"pass-generator/internal/settings"
)

func main() {

	configPath, err := settings.UserPath()

	if err != nil {
		fmt.Println("Не удалось определить путь к настройкам: ", err)
		return
	}

	settingsInst := settings.New(configPath)

	app := cli.New(settingsInst)

	if err := app.Run(); err != nil {
		fmt.Println("Ошибка работы приложения:", err)
	}

}
