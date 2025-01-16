package util

import (
	"testing"
	"time"
)

func TestShouldCheckIfFileExists(t *testing.T) {
	exists, err := FileExists("./files.go")
	t.Log(exists, err)
}

func TestShouldCheckIfDirectoryExists(t *testing.T) {
	exists, err := DirectoryExists("./")
	t.Log(exists, err)
	exists, err = DirectoryExists("k:\\abc")
	t.Log(exists, err)
}

func TestShouldCheckIfPathExists(t *testing.T) {
	exists, err := PathExists("./")
	t.Log(exists, err)
	exists, err = PathExists("k:\\abc")
	t.Log(exists, err)
}

func TestReadFileToString(t *testing.T) {
	t.Log(ReadFileToString("..\\feature.tpl"))
}

func TestListFiles(t *testing.T) {
	r, _ := ListFiles("c:\\temp", ".mp4")
	t.Log(r)
	r, _ = ListFiles("c:\\temp", ".jpg")
	t.Log(r)
}

func TestFileWithExt(t *testing.T) {
	t.Log(FileWithExt("/mnt/1.jpg", ".JPG"))
	t.Log(FileWithExt("/mnt/1.jpg", ".JPEG"))
}

func TestListImageFiles(t *testing.T) {
	v, _ := ListImageFiles("C:\\Users\\ywu\\Documents\\WeChat Files\\wygjyw\\FileStorage\\Video\\2022-11\\")
	for _, a := range v {
		a := a
		t.Log("def" + a)
		go func() {
			t.Log("abc" + a)
		}()
	}
	time.Sleep(time.Minute)
}
