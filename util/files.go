package util

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

func FileExists(path string) (exists bool, err error) {
	info, err := os.Stat(path)
	if err == nil {
		if info.IsDir() {
			return false, errors.New("path is a directory")
		}
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func FileSize(path string) (size int64, err error) {
	info, err := os.Stat(path)
	if err == nil {
		if info.IsDir() {
			return 0, errors.New("path is a directory")
		}
		return info.Size(), nil
	}
	if os.IsNotExist(err) {
		return 0, nil
	}
	return 0, err
}

func DirectoryExists(path string) (exists bool, err error) {
	info, err := os.Stat(path)
	if err == nil {
		if info.IsDir() {
			return true, nil
		}
		return false, errors.New("path is a file")
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func PathExists(path string) (exists bool, err error) {
	_, err = os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func ReadFileToString(fileName string) string {
	if f, err := os.Open(fileName); err == nil {
		defer f.Close()
		b, err := io.ReadAll(f)
		if err == nil {
			return string(b)
		}
	}
	return ""
}

func ReadFileToByte(fileName string) []byte {
	if f, err := os.Open(fileName); err == nil {
		defer f.Close()
		b, err := io.ReadAll(f)
		if err == nil {
			return b
		}
	}
	return nil
}

func WriteJsonToFile(fileName string, v interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return DumpBytes(fileName, bytes)
}

func ListFiles(folderName string, end string) ([]string, error) {
	var result []string
	hasExt := len(end) > 0
	err := filepath.Walk(folderName,
		func(fullPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				name := info.Name()
				if !hasExt || strings.LastIndex(name, end) != -1 {
					result = append(result, fullPath)
				}
			}
			return nil
		})
	return result, err
}

func ListImageFiles(folderName string) ([]string, error) {
	var result []string
	info, err := os.Stat(folderName)
	if err == nil {
		if info.IsDir() {
			files, _ := os.ReadDir(folderName)
			for _, file := range files {
				if !file.IsDir() {
					name := file.Name()
					ext := filepath.Ext(name)
					ext = strings.ToUpper(ext)
					if ext == ".PNG" || ext == ".JPG" || ext == ".JPEG" {
						imgFileName := filepath.Join(folderName, name)
						result = append(result, imgFileName)
					}
				}
			}
		} else {
			result = append(result, folderName)
		}
	}
	return result, err
}

func FileWithExt(fileName string, ext string) bool {
	fileExt := filepath.Ext(fileName)
	return strings.ToLower(fileExt) == strings.ToLower(ext)
}

func DumpFormToFile(file *multipart.FileHeader, fullPath string, fileName string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	var tempFile *os.File
	if len(fullPath) > 0 {
		tempFile, err = os.Create(fullPath)
	} else {
		tempFile, err = os.CreateTemp("", fileName)
	}
	if err != nil {
		return "", err
	}
	defer tempFile.Close()
	_, err = io.Copy(tempFile, src)
	src.Close()
	tempFile.Close()
	return tempFile.Name(), nil
}

func GetTempFileName(pattern string) (string, error) {
	tempFile, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", err
	}
	defer tempFile.Close()
	fileName := tempFile.Name()
	return fileName, nil
}

func GetLocalImageBase64(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	data := base64.StdEncoding.EncodeToString(fileBytes)
	return data, nil
}

func FileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	hash := md5.New()
	_, _ = io.Copy(hash, file)
	return hex.EncodeToString(hash.Sum(nil)), nil
}
