package config

import (
	"fmt"
	"github.com/jinzhu/configor"
	"os"
	"path/filepath"
	"strings"
)

type Error struct {
	Warning []struct {
		Id          string `yaml:"id"`
		Description string `yaml:"description"`
	} `yaml:"warning"`
}

var ErrorCfg = &Error{}

var errorMap map[string]string

func processMap() {
	for _, v := range ErrorCfg.Warning {
		v.Id = v.Id[3:]
		v.Id = strings.ToUpper(v.Id)
		errorMap[v.Id] = v.Description
	}
}

func loadErrorCfg(folder string) *Error {
	errorMap = make(map[string]string)
	if err := configor.New(&configor.Config{
		AutoReload: true,
		AutoReloadCallback: func(config interface{}) {
			processMap()
		},
	}).Load(ErrorCfg, filepath.Join(folder, "transcoder.error.yml")); err != nil {
		fmt.Printf("Failed to read transcoder.error.yml: %s", err)
		os.Exit(2)
	}
	processMap()
	return ErrorCfg
}

func FindDescriptionById(id string) string {
	if len(id) > 8 {
		id = id[3:]
		id = strings.ToUpper(id)
		return errorMap[id]
	}
	return ""
}
