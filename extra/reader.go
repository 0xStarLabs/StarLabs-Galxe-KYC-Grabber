package extra

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func ReadTxtFile(fileName string, filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
		return []string{}
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	var items []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		items = append(items, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed to read file: %v", err)
		return []string{}
	}

	Logger{}.Info("Successfully loaded %d %s.", len(items), fileName)
	return items
}

func ReadPrivateKeys() []string {
	file, err := os.Open("data/wallets.txt")
	if err != nil {
		log.Fatalf("Failed to open wallets.txt: %v", err)
		return []string{}
	}
	defer file.Close()

	var items []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Count(line, " ") >= 2 {
			privateKey := ExtractPrivateKeyFromMnemo(line)
			items = append(items, privateKey)
		} else {
			items = append(items, line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed to read file: %v", err)
		return []string{}
	}

	Logger{}.Info("Successfully loaded %d wallets from wallets.txt.", len(items))
	return items
}

func NoProxies() bool {
	var userChoice int
	fmt.Println("No proxies were detected. Do you want to continue without proxies? (1 or 2)")
	fmt.Println("[1] Yes")
	fmt.Println("[2] No")
	fmt.Print(">> ")
	_, err := fmt.Scan(&userChoice)
	if err != nil {
		Logger{}.Error("Wrong input. Enter the number.")
		panic(err)
	}

	return userChoice == 1
}

func UserInputInteger(textToShow string) (int, bool) {
	var input int

	fmt.Print(textToShow + ": ")
	reader := bufio.NewReader(os.Stdin)
	userInput, _ := reader.ReadString('\n')
	userInput = strings.TrimSpace(userInput)

	input, err := strconv.Atoi(userInput)
	if err != nil {
		Logger{}.Error("Wrong input.")
		return input, false
	}
	return input, true
}
