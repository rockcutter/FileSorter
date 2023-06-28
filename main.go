package main

import (
	"FileSorter/FileSorter"
	"fmt"
	"log"
	"os"
)

var OptionArray [FileSorter.FILE_SORTER_OPTION_COUNT]string

func CheckFileSorterOption(option string) (string, error) {
}

func PrintUsage() {
	fmt.Println("Usage: ./FileSorter.exe [copy | move]")
}

func main() {
	if len(os.Args) < 2 {
		PrintUsage()
		return
	}
	var fileSorterOption int = -1
	if os.Args[1] == "copy" {
		fileSorterOption = FileSorter.FILE_SORTER_OPTION_COPY
	} else if os.Args[1] == "move" {
		fileSorterOption = FileSorter.FILE_SORTER_OPTION_MOVE
	}
	if fileSorterOption == -1 {
		PrintUsage()
		return
	}

	const TARGET_DIR = "."

	var fs = FileSorter.MakeSorter()

	fileArray, err := os.ReadDir(TARGET_DIR)

	if err != nil {
		log.Fatal(TARGET_DIR, "内のファイルを取得できません")
	}
	log.Println("ファイル一覧取得完了")

	for _, file := range fileArray {
		if file.IsDir() {
			continue
		}
		fs.AppendFile(file)
	}

	for key, value := range fs.GetDateFileMap() {
		for _, fileName := range value {
			log.Printf("key: %s, fileName: %s", key, fileName)
		}
	}

	fs.Commit(fileSorterOption)
}
