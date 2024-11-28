package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	outputFile   = "project_source.txt"
	thisFileName = "collect_source.go"
)

var (
	// Конфигурация расширений файлов для парсинга
	fileExtensions = []string{
		".go",
		".html",
		".sql",
		".svelte",
		".js",
		// Можно добавить другие расширения
	}

	// Директории для исключения
	excludedDirs = []string{
		".git",
		"vendor",
		"node_modules",
		"dist",
		"build",
		"coverage",
		".next",
	}
)

// Helper function to check if directory should be excluded
func shouldExcludeDir(name string) bool {
	for _, dir := range excludedDirs {
		if name == dir {
			return true
		}
	}
	return false
}

func main() {
	f, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer f.Close()

	// Записываем мета-информацию
	fmt.Fprintf(f, "Project Source Code Export\n")
	fmt.Fprintf(f, "Generated: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(f, "Parsing extensions: %v\n", fileExtensions)
	fmt.Fprintf(f, "Excluded directories: %v\n", excludedDirs)
	fmt.Fprintf(f, "Working Directory: %s\n", getCurrentDir())
	fmt.Fprintf(f, "\n"+strings.Repeat("-", 80)+"\n\n")

	// Выводим дерево проекта
	fmt.Fprintf(f, "Project Tree:\n")
	fmt.Fprintf(f, "=============\n\n")
	printProjectTree(f, ".", 0, make(map[string]bool))
	fmt.Fprintf(f, "\n"+strings.Repeat("-", 80)+"\n\n")

	// Выводим содержимое файлов
	fmt.Fprintf(f, "Source Code:\n")
	fmt.Fprintf(f, "============\n\n")

	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check for excluded directories
		for _, dir := range excludedDirs {
			if strings.Contains(path, dir+string(os.PathSeparator)) {
				return filepath.SkipDir
			}
		}

		// Пропускаем директории
		if info.IsDir() {
			return nil
		}

		// Пропускаем файл самого скрипта
		if strings.HasSuffix(path, thisFileName) {
			return nil
		}

		// Проверяем расширение файла
		ext := filepath.Ext(path)
		shouldParse := false
		for _, allowedExt := range fileExtensions {
			if ext == allowedExt {
				shouldParse = true
				break
			}
		}

		if !shouldParse {
			return nil
		}

		// Читаем содержимое файла
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Записываем путь к файлу и его содержимое
		fmt.Fprintf(f, "// File: %s\n", path)
		fmt.Fprintf(f, "// Size: %d bytes\n", info.Size())
		fmt.Fprintf(f, "// Extension: %s\n", ext)
		fmt.Fprintf(f, strings.Repeat("-", 40)+"\n\n")
		fmt.Fprintf(f, "%s\n\n", string(content))
		fmt.Fprintf(f, strings.Repeat("=", 80)+"\n\n")

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		return
	}

	fmt.Printf("Source code has been exported to %s\n", outputFile)
}

// printProjectTree выводит дерево проекта
func printProjectTree(f *os.File, path string, level int, isLast map[string]bool) {
	files, err := os.ReadDir(path)
	if err != nil {
		return
	}

	for i, file := range files {
		// Skip excluded directories
		if shouldExcludeDir(file.Name()) {
			continue
		}

		// Формируем префикс для текущего уровня
		prefix := ""
		for l := 0; l < level; l++ {
			if isLast[fmt.Sprintf("%d", l)] {
				prefix += "    "
			} else {
				prefix += "│   "
			}
		}

		isLastItem := i == len(files)-1
		if isLastItem {
			prefix += "└── "
		} else {
			prefix += "├── "
		}

		// Записываем текущий файл/директорию
		fullPath := filepath.Join(path, file.Name())
		fmt.Fprintf(f, "%s%s", prefix, file.Name())
		if !file.IsDir() {
			if info, err := file.Info(); err == nil {
				fmt.Fprintf(f, " (%d bytes)", info.Size())
			}
		}
		fmt.Fprintf(f, "\n")

		// Рекурсивно обрабатываем поддиректории
		if file.IsDir() {
			newIsLast := make(map[string]bool)
			for k, v := range isLast {
				newIsLast[k] = v
			}
			newIsLast[fmt.Sprintf("%d", level)] = isLastItem
			printProjectTree(f, fullPath, level+1, newIsLast)
		}
	}
}

func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "unknown"
	}
	return dir
}
