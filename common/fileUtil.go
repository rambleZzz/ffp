package common

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

// 读文件
func ReadFile(filename string) (list []string, err error) {
	f, err := os.Open(filename)
	if err != nil {
		log.Printf("打开文件%s失败:%v\n", filename, err)
		os.Exit(0)
	}
	defer f.Close()
	var text_list []string
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text != "" {
			text_list = append(text_list, text)
		}
	}
	return text_list, nil
}

// 写文件
func WriteFile(filename string, s string) (err error) {
	var f *os.File
	defer f.Close()
	if CheckFileIsExist(filename) { //如果文件存在
		f, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0666) //打开文件
		if err != nil {
			log.Printf("%v\n", err)
		}
	} else {
		f, err = os.Create(filename) //创建文件
		if err != nil {
			log.Printf("%v\n", err)
		}
	}
	_, err = io.WriteString(f, s) //写入文件(字符串)
	if err != nil {
		log.Printf("%v\n", err)
	}
	return nil
}

// 检查文件是否存在
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// FolderExists checks if a folder exists
func FolderExists(folderpath string) bool {
	_, err := os.Stat(folderpath)
	return !os.IsNotExist(err)
}

func DirCreate(dirname string) {
	err := os.MkdirAll(dirname, 0755)
	if err != nil {
		log.Fatalln(err)
	}
}
