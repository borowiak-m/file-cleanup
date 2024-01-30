package main

import (
	"bufio"
	"fmt"
	"os"
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
