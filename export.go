package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// https://github.com/zeroidentidad/rawIndex_anime_girls_holding_dev_books
	dirPath := "./rawIndex_anime_girls_holding_dev_books"

	// Base URL to replace local root
	baseURL := "https://raw.githubusercontent.com/zeroidentidad/rawIndex_anime_girls_holding_dev_books/refs/heads/master/"

	outputFile := "output.sql"
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Head SQL
	_, err = file.WriteString("-- Inserts waifus\n")
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	// Get absolute path root folder for precise handling
	absDirPath, err := filepath.Abs(dirPath)
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		return
	}

	absDirPath = filepath.Clean(absDirPath)

	// Browse the folder
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only add files (not folders) to the list
		if !info.IsDir() {
			absFilePath, err := filepath.Abs(path)
			if err != nil {
				return err
			}

			absFilePath = filepath.Clean(absFilePath)

			if !strings.HasPrefix(absFilePath, absDirPath) {
				return fmt.Errorf("file %s is not inside root directory %s", absFilePath, absDirPath)
			}

			relativePath := strings.TrimPrefix(absFilePath, absDirPath)
			relativePath = strings.TrimPrefix(relativePath, string(os.PathSeparator))

			// to UNIX format
			relativePath = filepath.ToSlash(relativePath)

			finalURL := baseURL + relativePath

			sql := fmt.Sprintf("INSERT INTO waifus (file_path) VALUES ('%s');\n", finalURL)

			_, err = file.WriteString(sql)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error while browsing folder:", err)
		return
	}

	fmt.Printf("SQL statements stored in '%s'\n", outputFile)
}
