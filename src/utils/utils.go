package utils

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
)

func GetHomeDir() string {
	userHomeDir, _ := os.UserHomeDir()
	return fmt.Sprintf("%s/%s", userHomeDir, ".youtube-channel-dl")
}

func GetLogDir() string {
	homeDir := GetHomeDir()
	return fmt.Sprintf("%s/%s", homeDir, "logs")
}

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func ByteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

func ParseFilename(keyString string) (filename string) {
	ss := strings.Split(keyString, "/")
	s := ss[len(ss)-1]
	return s
}

func Capitalize(str string) string {
	if str == "id" {
		return "ID"
	} else {
		return strings.ToUpper(str[0:1]) + str[1:]
	}
}

func WriteLog(text string) error {
	pathFile := fmt.Sprintf("%s/%s", GetLogDir(), "youtube-channel-dl.log")
	lock := sync.Mutex{}
	lock.Lock()
	defer lock.Unlock()
	file, errFile := os.OpenFile(pathFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if errFile != nil {
		return errFile
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)
	lineToWrite := text
	lineToWrite += "\n"
	_, errWrite := file.WriteString(lineToWrite)
	if errWrite != nil {
		return errWrite
	}
	return nil
}

func InArray(array []string, value string) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}

func GetTime(strTime string) int {
	infoTime := strings.TrimSpace(strTime)
	if strings.Contains(infoTime, "(") {
		infoTime = strings.Split(infoTime, "(")[1]
		infoTime = strings.Split(infoTime, ")")[0]
		infos := strings.Split(infoTime, ",")
		if len(infos) == 2 {
			from, errFrom := strconv.Atoi(infos[0])
			if errFrom != nil {
				from = 0
			}
			to, errTo := strconv.Atoi(infos[1])
			if errTo != nil {
				to = 0
			}
			seconds := RandomInt(from, to)
			return seconds
		} else {
			seconds, err := strconv.Atoi(infos[0])
			if err != nil {
				seconds = 0
			}
			return seconds
		}
	} else {
		seconds, err := strconv.Atoi(infoTime)
		if err != nil {
			seconds = 0
		}
		return seconds
	}
}

func GetValidFileName(fileName string) string {
	fileName = strings.ReplaceAll(fileName, "\\", "")
	fileName = strings.ReplaceAll(fileName, "/", "")
	fileName = strings.ReplaceAll(fileName, ":", "")
	fileName = strings.ReplaceAll(fileName, "*", "")
	fileName = strings.ReplaceAll(fileName, "?", "")
	fileName = strings.ReplaceAll(fileName, "\"", "")
	fileName = strings.ReplaceAll(fileName, "<", "")
	fileName = strings.ReplaceAll(fileName, ">", "")
	fileName = strings.ReplaceAll(fileName, "|", "")
	return fileName
}
