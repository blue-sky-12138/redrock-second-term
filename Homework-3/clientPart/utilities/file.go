package utilities

import (
	"encoding/base64"
	"os"
	"strconv"
	"strings"
	"time"
)

func Base64Encoding(file *os.File) (string, error) {
	var (
		source = make([]byte, 10000000)
	)

	suffix := GetSuffix(file.Name())

	n, err := file.Read(source)
	if err != nil {
		return "", err
	}

	res := base64.StdEncoding.EncodeToString(source[:n])
	return res + suffix, nil
}

func Base64DeEncoding(s string) error {
	suffix := GetSuffix(s)
	s = strings.Replace(s, suffix, "", -1)

	bytes, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return err
	}

	f, _ := os.Create(strconv.FormatInt(time.Now().Unix(), 10) + suffix)
	_, err = f.Write(bytes)
	defer f.Close()

	if err != nil {
		return err
	}

	return nil
}

func GetSuffix(fileName string) string {
	//获取文件后缀名
	var (
		fileNameLen = len(fileName) //文件名总长
		suffix      string
	)
	for i := fileNameLen; ; i-- {
		if fileName[i-1:i] == "." || fileNameLen-i > 1000 {
			suffix = fileName[i-1 : fileNameLen]
			break
		}
	}
	return suffix
}
