package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jinzhu/configor"
)

type Config struct {
	Web struct {
		Host  string `yaml:"host" default:"0.0.0.0"`
		Port  string `yaml:"port" default:"17500"`
		SPort string `yaml:"sport" default:"17600"`
	} `yaml:"web"`
	Database struct {
		Dialect   string `yaml:"dialect" default:"sqlite3"`
		Host      string `yaml:"host" default:"migu_ai_server.db"`
		Port      string `yaml:"port" default:"3306"`
		Dbname    string `yaml:"dbname" default:"server"`
		Username  string `yaml:"username" default:"root"`
		Password  string `yaml:"password" default:"root"`
		Migration bool   `default:"true"`
	} `yaml:"database"`
	Extension struct {
		CorsEnabled bool `yaml:"cors_enabled" default:"true"`
	} `yaml:"extension"`
	Log struct {
		RequestLogFormat string `yaml:"request_log_format" default:"${remote_ip} ${uri} ${method} ${status} ${content_length} ${response_length}"`
	} `yaml:"log"`
	ConfigPath  string
	HttpTimeout int64 `yaml:"httpTimeout" default:"120"`
	Etcd        struct {
		NewId      int64    `yaml:"newId" default:"2"` // 采用新的雪花算法
		Election   bool     `yaml:"election" default:"true"`
		ServerList []string `yaml:"serverList"`
		Leader     string   `yaml:"leader" default:"/migu/ai/leader"`
		Alarm      string   `yaml:"alarm" default:"/migu/ai/alarm"`
		Storage    string   `yaml:"storage" default:"/migu/ai/storage"`
	} `yaml:"etcd"`
	QuerySecret struct {
		Key   string `yaml:"key" default:"api-security-key"`
		Value string `yaml:"value"`
	} `yaml:"querySecret"`

	AiAddress string `yaml:"aiAddress" default:"http://172.17.81.10:18900"`
}

type AppConfig struct {
	ImageExtList []struct {
		Ext    string `yaml:"ext"`
		Format string `yaml:"format"`
	} `yaml:"imageExtList"`
	AudioExtList []struct {
		Ext    string `yaml:"ext"`
		Format string `yaml:"format"`
	} `yaml:"audioExtList"`
	AppId []struct {
		Id                   int64    `yaml:"id"`
		Name                 string   `yaml:"name"`
		MediaAnalyzer        string   `yaml:"mediaAnalyzer"`
		MediaAnalyzerCommand string   `yaml:"mediaAnalyzerCommand"`
		MediaAnalyzerTimeout int64    `yaml:"mediaAnalyzerTimeout" default:"2"`
		CutTranscoder        string   `yaml:"cutTranscoder"`
		CutTranscoderCommand []string `yaml:"cutTranscoderCommand"`
		Transcoder           string   `yaml:"transcoder"`
		TranscoderCommand    []string `yaml:"transcoderCommand"`
		Editor               string   `yaml:"editor"`
		EditorCommand        string   `yaml:"editorCommand"`
	} `yaml:"appId"`
}

var Cfg = &Config{}
var AppCfg = &AppConfig{}

func processPath() {
}

func Load(folder string) *Config {
	Cfg.ConfigPath = folder
	serverConfigFileName := "server.application.yml"
	if os.Getenv("ENV") == "dev" {
		serverConfigFileName = "server.application.dev.yml"
	}
	if err := configor.New(&configor.Config{
		AutoReload: true,
		AutoReloadCallback: func(config interface{}) {
		},
	}).Load(Cfg, filepath.Join(folder, serverConfigFileName)); err != nil {
		fmt.Printf("Failed to read server.application.yml: %s", err)
		os.Exit(2)
	}
	//if err := configor.New(&configor.Config{
	//	AutoReload: true,
	//}).Load(AppCfg, filepath.Join(folder, "server.config.yml")); err != nil {
	//	fmt.Printf("Failed to read server.config.yml: %s", err)
	//	os.Exit(2)
	//}
	processPath()
	//loadErrorCfg(folder)
	return Cfg
}

func GetHttpTimeout() time.Duration {
	return (time.Duration)(Cfg.HttpTimeout) * time.Second
}

func GetAppNameById(appId int64) string {
	for _, v := range AppCfg.AppId {
		if v.Id == appId {
			return v.Name
		}
	}
	return ""
}
