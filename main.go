package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func readConfig(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("Read line: ", line)
		lines = append(lines, line)
	}
	return lines, scanner.Err()
}

func processFolder(folderPath string, logChan chan string) error {
	archivePath := filepath.Join(folderPath, "Archive")
	if _, err : = os.Stat(archivePath); os.IsNotExist(err) {
		if err := os.Mkdir(archivePath, 0755); err !=nil {
			return err
		}
	}

	filesInFolder, err := ioutil.ReadDir(folderPath)
	if err != nil {return err}

	for _, file := range filesInFolder {
		if file.IsDir() {continue}
		if time.Since(file.ModTime()).Hours() > (24*30*3) {
			oldPath := filepath.Join(folderPath, file.Name())
			newPath := filepath.Join(archivePath, "archived_"+file.Name())

			if err := os.Rename(oldPath,newPath); err != nil {return err}

			logChan <- fmt.Sprintf("Moved file: %s to %s", oldPath, newPath)
		}
	}

	return nil
} 

func main() {
	filePath := "C:/Users/mariusz.borowiak/Documents/Dev/GO/file-cleanup/config/config.txt"
	fmt.Println("Attempting to read file: ", filePath)
	folders, err := readConfig(filePath)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}
	fmt.Printf("Read %d folders from the config file.\n", len(folders))

}
