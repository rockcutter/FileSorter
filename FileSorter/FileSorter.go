package FileSorter

import (
	"errors"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

const (
	FILE_SORTER_OPTION_COPY = iota
	FILE_SORTER_OPTION_MOVE
	FILE_SORTER_OPTION_COUNT
)

type FileSorter struct {
	directoryNameFormat string
	dateFileMap         map[string][]string
	targetDirectory     string
}

func MakeSorter() *FileSorter {
	var sorterPtr = new(FileSorter)
	sorterPtr.dateFileMap = make(map[string][]string)
	sorterPtr.directoryNameFormat = "2006_01_02"
	sorterPtr.targetDirectory = "."
	return sorterPtr
}

//setter
func (this *FileSorter) SetDirectoryNameFormat(format string) {
	this.directoryNameFormat = format
}

//setter
func (this *FileSorter) SetTargetDirectory(directoryPath string) {
	this.targetDirectory = directoryPath
}

func (this *FileSorter) AppendFile(file fs.DirEntry) {
	var fileName = file.Name()
	var fileInfo, err = file.Info()
	var fileModifiedDateString string

	if err != nil {
		fileModifiedDateString = "ファイル情報を取得できません"
	} else {
		fileModifiedDateString = fileInfo.ModTime().Format(this.directoryNameFormat)
	}

	log.Printf("Sorter.AppendFile: ファイル(%s 更新日時: %s)を追加", fileName, fileModifiedDateString)
	this.dateFileMap[fileModifiedDateString] =
		append(this.dateFileMap[fileModifiedDateString], file.Name())
}

func (this *FileSorter) GetDateFileMap() map[string][]string {
	return this.dateFileMap
}

func (this *FileSorter) IsFileSorterOptionValid(fileSorterOption int) bool {
	if FILE_SORTER_OPTION_COUNT > fileSorterOption {
		return true
	}
	return false
}

func (this *FileSorter) CopyFile(src string, dst string) error {
	srcfp, err := os.Open(src)

	if err != nil {
		return err
	}
	defer srcfp.Close()

	dstfp, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstfp.Close()

	_, err = io.Copy(dstfp, srcfp)

	if err != nil {
		return err
	}
	return nil
}

func (this *FileSorter) MoveFile(src string, dst string) error {
	return os.Rename(src, dst)
}

func (this *FileSorter) ExecuteOperation(fileSorterOption int, src string, dst string) error {
	if !this.IsFileSorterOptionValid(fileSorterOption) {
		return errors.New("Sorter.AppendFile: invalid FileSorter option")
	}

	var err error
	switch fileSorterOption {
	case FILE_SORTER_OPTION_COPY:
		log.Println("Sorter.AppendFile: operation: copy")
		err = this.CopyFile(src, dst)
	case FILE_SORTER_OPTION_MOVE:
		log.Println("Sorter.AppendFile: operation: move")
		err = this.MoveFile(src, dst)
	}

	if err != nil {
		return err
	}

	return nil
}

func (this *FileSorter) Commit(fileSorterOption int) {

	for key, value := range this.dateFileMap {
		log.Printf("Sorter.AppendFile: mkdir %s", key)
		if os.MkdirAll(key, 0775) != nil {
			log.Println("Sorter.AppendFile: mkdir failed")
			continue
		}
		for _, file := range value {
			err := this.ExecuteOperation(fileSorterOption, file, filepath.Join(key, file))
			if err != nil {
				log.Println("Sorter.AppendFile: operation failed: ", err)
				continue
			}
			log.Println("Sorter.AppendFile: operation successed")
		}
	}
}
