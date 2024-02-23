package extra

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
	"sync"
)

func CreateXLSXFromSyncMap(KYCLinks *sync.Map, keys []string, filePath string) {
	f := excelize.NewFile()
	index := 1

	for _, key := range keys {
		value, ok := KYCLinks.Load(key)
		if !ok {
			continue
		}
		cellKey := fmt.Sprintf("A%d", index)
		cellValue := fmt.Sprintf("B%d", index)
		if err := f.SetCellValue("Sheet1", cellKey, key); err != nil {
			Logger{}.Error("Failed to save data to XLSX file. Saving all to TXT file.")
			saveToTxtFile(KYCLinks, keys)
			return
		}
		if err := f.SetCellValue("Sheet1", cellValue, value); err != nil {
			return
		}
		index++
	}

	if err := f.SaveAs(filePath); err != nil {
		Logger{}.Error("Failed to save XLSX file.")
		saveToTxtFile(KYCLinks, keys)
	}
}

func saveToTxtFile(KYCLinks *sync.Map, keys []string) {
	file, err := os.Create("data/collected_links.txt")
	if err != nil {
		Logger{}.Error("Failed to create TXT file.")
		return
	}
	defer file.Close()

	for _, key := range keys {
		value, ok := KYCLinks.Load(key)
		if !ok {
			continue
		}
		line := fmt.Sprintf("%v:%v\n", key, value)
		if _, err := file.WriteString(line); err != nil {
			Logger{}.Error("Failed to save links to TXT file. Please copy them manually.")
			printMap(KYCLinks, keys)
			return
		}
	}
}

func printMap(KYCLinks *sync.Map, keys []string) {
	for _, key := range keys {
		value, ok := KYCLinks.Load(key)
		if ok {
			fmt.Printf("%v:%v\n", key, value)
		}
	}
}
