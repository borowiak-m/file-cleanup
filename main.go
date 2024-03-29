package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	ConfigFilePath         = "C:/Users/mariusz.borowiak/Documents/Dev/GO/file-cleanup/config/config.txt"
	LogFilePath            = "C:/Users/mariusz.borowiak/Documents/Dev/GO/file-cleanup/log/log.txt"
	ArchiveSubDirName      = "Archive"
	ArchivedFileNamePrefix = "archived_"
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
	archivePath := filepath.Join(folderPath, ArchiveSubDirName)
	if _, err := os.Stat(archivePath); os.IsNotExist(err) {
		if err := os.Mkdir(archivePath, 0755); err != nil {
			return err
		}
	}

	filesInFolder, err := os.ReadDir(folderPath)
	if err != nil {
		return err
	}

	for _, file := range filesInFolder {
		fileInfo, err := file.Info()
		if err != nil {
			return err
		}
		if fileInfo.IsDir() {
			continue
		}
		if time.Since(fileInfo.ModTime()).Hours() > (24 * 30 * 3) {
			oldPath := filepath.Join(folderPath, file.Name())
			newPath := filepath.Join(archivePath, ArchivedFileNamePrefix+file.Name())

			if err := os.Rename(oldPath, newPath); err != nil {
				return err
			}
			currentTime := time.Now().Format("2006-01-02 15:04:05")
			logMessage := fmt.Sprintf("[%s] Moved file: %s to %s", currentTime, oldPath, newPath)
			logChan <- logMessage
		}
	}

	return nil
}

func deleteEmptyFolders(folderPath string, logChan chan string) error {
	dirs, err := os.ReadDir(folderPath)
	if err != nil {
		return err
	}
	logChan <- fmt.Sprintf("*** Folder %s has %d elements", folderPath, len(dirs))

	// if it does contain something, loop over it
	for _, dir := range dirs {
		if !dir.IsDir() {
			logChan <- dir.Name() + " is not a directory"
			// if not a directory, move on
			continue
		}
		logChan <- dir.Name() + " is a directory"
		dirPath := filepath.Join(folderPath, dir.Name())
		// for each directory recurse over
		if err := deleteEmptyFolders(dirPath, logChan); err != nil {
			logChan <- "Error when recursing over dir: " + dirPath + "with error: " + err.Error()
		}
	}
	// if the folder does not contain anything, remove it
	if len(dirs) == 0 {
		logChan <- "Empty folder found:" + folderPath
		if err := os.Remove(folderPath); err != nil {
			logChan <- "Error trying to remove dir: " + folderPath + " with error: " + err.Error()
			return err
		}
		logChan <- "Deleted empty folder: " + folderPath
		return nil
	}
	return nil
}

func logActivity(logChan chan string, doneChan chan bool) {
	logFile, err := os.OpenFile(LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		doneChan <- true
		return
	}
	defer logFile.Close()

	for logMsg := range logChan {
		if _, err := logFile.WriteString(logMsg + "\n"); err != nil {
			fmt.Println("Error writing to log file:", err)
			continue
		}
	}
	doneChan <- true
}

func main() {
	filePath := ConfigFilePath
	fmt.Println("Attempting to read file: ", filePath)
	folders, err := readConfig(filePath)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}
	fmt.Printf("Read %d folders from the config file.\n", len(folders))

	logChan := make(chan string)
	doneChan := make(chan bool)

	go logActivity(logChan, doneChan)
	for _, folder := range folders {
		if err := processFolder(folder, logChan); err != nil {
			fmt.Println("Error processing folder:", err)
		}
		if err := deleteEmptyFolders(folder, logChan); err != nil {
			fmt.Println("Error during deletion checks:", err)
		}
	}

	close(logChan)
	<-doneChan
}
