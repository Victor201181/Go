//go:build gocp
// +build gocp

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

const version = "1.0.0"

// copyFile копирует файл из sourcePath в destPath
func copyFile(sourcePath, destPath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("ошибка открытия исходного файла: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("ошибка создания целевого файла: %w", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("ошибка копирования данных: %w", err)
	}

	err = destFile.Sync()
	if err != nil {
		return fmt.Errorf("ошибка синхронизации данных на диск: %w", err)
	}

	return nil
}

// moveFile перемещает файл из sourcePath в destPath
func moveFile(sourcePath, destPath string) error {
	if err := copyFile(sourcePath, destPath); err != nil {
		return err
	}
	return os.Remove(sourcePath)
}

// promptUser запрашивает у пользователя путь к файлу-источнику и файлу-назначению
func promptUser() (string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Введите путь к файлу-источнику: ")
	sourcePath, _ := reader.ReadString('\n')
	sourcePath = sourcePath[:len(sourcePath)-1] // Удаляем символ новой строки

	fmt.Print("Введите путь к файлу-назначению: ")
	destPath, _ := reader.ReadString('\n')
	destPath = destPath[:len(destPath)-1] // Удаляем символ новой строки

	return sourcePath, destPath
}

func main() {
	helpFlag := flag.Bool("h", false, "Показать справку")
	versionFlag := flag.Bool("v", false, "Показать версию программы")

	// Разбор флагов командной строки
	flag.Parse()

	// Обработка флага -h
	if *helpFlag {
		fmt.Println("Использование: gocp <команда> [опции] [файл-источник] [файл-назначение]")
		fmt.Println("Команды:")
		fmt.Println("  copy    Копирует файл из одного места в другое")
		fmt.Println("  move    Перемещает файл, копируя его и удаляя оригинал")
		fmt.Println("  input   Запрашивает пути к файлам для копирования")
		fmt.Println("Опции:")
		fmt.Println("  -h      Показать справку")
		fmt.Println("  -v      Показать версию программы")
		return
	}

	// Обработка флага -v
	if *versionFlag {
		fmt.Printf("gocp версия %s\n", version)
		return
	}

	// Обработка подкоманды
	if flag.NArg() < 1 {
		fmt.Println("Ошибка: недостаточно аргументов")
		fmt.Println("Использование: gocp <команда> [опции] [файл-источник] [файл-назначение]")
		return
	}

	command := flag.Arg(0)

	switch command {
	case "copy":
		if flag.NArg() != 3 {
			fmt.Println("Ошибка: необходимо указать путь к файлу-источнику и файлу-назначению")
			return
		}
		sourcePath := flag.Arg(1)
		destPath := flag.Arg(2)
		if err := copyFile(sourcePath, destPath); err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка при копировании: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Файл успешно скопирован.")

	case "move":
		if flag.NArg() != 3 {
			fmt.Println("Ошибка: необходимо указать путь к файлу-источнику и файлу-назначению")
			return
		}
		sourcePath := flag.Arg(1)
		destPath := flag.Arg(2)
		if err := moveFile(sourcePath, destPath); err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка при перемещении: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Файл успешно перемещен.")

	case "input":
		sourcePath, destPath := promptUser()
		if err := copyFile(sourcePath, destPath); err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка при копировании: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Файл успешно скопирован.")

	default:
		fmt.Printf("Неизвестная команда: %s\n", command)
		fmt.Println("Доступные команды: copy, move, input")
		os.Exit(1)
	}
}
