package util

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func CopyFile(src string, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	target, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer target.Close()
	_, err = io.Copy(target, source)
	return err
}

func Uuid() string {
	uuidWithHyphen := uuid.New()
	return strings.Replace(uuidWithHyphen.String(), "-", "", -1)
}
func Uuid8() string {
	fullUuid := Uuid()
	id := string([]byte(fullUuid)[:8])
	return id
}

func DumpString(fileName string, data string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(data)
	return err
}

func DumpBytes(fileName string, data []byte) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}

func DumpInterface(fileName string, data interface{}) error {
	r, err := json.Marshal(data)
	if err != nil {
		return err
	}
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(r)
	return err
}

func GoId() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.ParseInt(idField, 10, 0)
	if err != nil {
		fmt.Sprintf("cannot get goroutine id: %v", err)
		return -1
	}
	return id
}

func TruncateString(str string, length int, unicodeMode bool) string {
	if length <= 0 {
		return ""
	}

	if !unicodeMode {
		orgLen := len(str)
		if orgLen <= length {
			return str
		}
		return str[:length]
	}

	truncated := ""
	count := 0
	for _, char := range str {
		truncated += string(char)
		count++
		if count >= length {
			break
		}
	}
	return truncated
}

func ToPrintString(v interface{}) string {
	r, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return ""
	}
	return string(r)
}

func ToString(v interface{}) string {
	r, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(r)
}

func Int64Abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

func Int64Min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func Int64Max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func GetExeDir() string {
	exePath := os.Args[0]
	exeDir := filepath.Dir(exePath)
	exeAbsDir, _ := filepath.Abs(exeDir)
	return exeAbsDir
}

func ChangeWorkPath(wp string) {
	if len(wp) > 0 {
		curWorkPath, _ := os.Getwd()
		fmt.Printf("before current work path: %s\n", curWorkPath)
		if wp != curWorkPath {
			os.Chdir(wp)
		}
		curWorkPath, _ = os.Getwd()
		fmt.Printf("after current work path: %s\n", curWorkPath)
	}
}

func NetworkStream(name string) bool {
	return strings.HasPrefix(name, "udp://") ||
		strings.HasPrefix(name, "rtmp://") ||
		strings.HasPrefix(name, "rtsp://") ||
		strings.HasPrefix(name, "srt://")
}

var __idLock = sync.RWMutex{}
var __id = int64(1)

func GetNextId() int64 {
	defer __idLock.Unlock()
	__idLock.Lock()
	__id++
	return __id
}

func Base64ToFile(fileData string) (*os.File, error) {
	var tempFile *os.File
	fileData = strings.ReplaceAll(fileData, "\r", "")
	fileData = strings.ReplaceAll(fileData, "\n", "")
	found := strings.Index(fileData, ",")
	if found != -1 {
		fileData = fileData[found+1:]
	}
	img, err := base64.StdEncoding.DecodeString(fileData)
	if err != nil {
		return nil, err
	}
	tempFile, err = os.CreateTemp("", "img*.jpg")
	if err != nil {
		return nil, err
	}
	defer tempFile.Close()
	tempFile.Write(img)
	return tempFile, nil
}

func UrlToFile(fullURLFile string, fullName string, fileExt string) (*os.File, error) {
	var err error
	var tempFile *os.File
	if len(fullName) > 0 {
		fullName += fileExt
		tempFile, err = os.Create(fullName)
	} else {
		tempFile, err = os.CreateTemp("", "fs*"+fileExt)
	}
	if err != nil {
		return nil, err
	}
	defer tempFile.Close()
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	resp, err := client.Get(fullURLFile)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return nil, err
	}
	return tempFile, nil
}

func UrlToFile2(fullURLFile string, fullName string) (*os.File, error) {
	var err error
	var tempFile *os.File
	tempFile, err = os.Create(fullName)
	if err != nil {
		return nil, err
	}
	defer tempFile.Close()
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	resp, err := client.Get(fullURLFile)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return nil, err
	}
	return tempFile, nil
}

func XmlEscape(s string) string {
	s = strings.Replace(s, "&", "&amp;", -1)
	s = strings.Replace(s, "\"", "&quot;", -1)
	s = strings.Replace(s, "'", "&apos;", -1)
	s = strings.Replace(s, "<", "&lt;", -1)
	s = strings.Replace(s, ">", "&gt;", -1)
	return s
}

func IsLocalHost(addr string) bool {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		host = addr
	}
	if host == "localhost" || host == "127.0.0.1" || host == "::1" {
		return true
	}
	ifaces, err := net.Interfaces()
	if err != nil {
		return false
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err == nil {
			for _, addr := range addrs {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				}
				if ip != nil && ip.String() == host {
					return true
				}
			}
		}
	}
	return false
}

func GenSignature(params url.Values, secretKey string) string {
	var paramStr string
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		paramStr += key + params[key][0]
	}
	paramStr += secretKey
	md5Reader := md5.New()
	md5Reader.Write([]byte(paramStr))
	return hex.EncodeToString(md5Reader.Sum(nil))
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", m.Alloc/1024/1024)
	fmt.Printf("\tTotalAlloc = %v MiB", m.TotalAlloc/1024/1024)
	fmt.Printf("\tSys = %v MiB", m.Sys/1024/1024)
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func ParseJsonMap(m map[string]interface{}, f func(k string, v interface{})) {
	for k, v := range m {
		switch v.(type) {
		case map[string]interface{}:
			ParseJsonMap(v.(map[string]interface{}), f)

		case []interface{}:
			ParseJsonArray(v.([]interface{}), f)

		default:
			f(k, v)
		}
	}
}

func ParseJsonArray(a []interface{}, f func(k string, v interface{})) {
	for _, val := range a {
		switch val.(type) {
		case map[string]interface{}:
			ParseJsonMap(val.(map[string]interface{}), f)

		case []interface{}:
			ParseJsonArray(val.([]interface{}), f)

		default:
		}
	}
}
