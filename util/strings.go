package util

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func IsNumeric(number string) bool {
	_, err := strconv.Atoi(number)
	return err == nil
}

func ConvertToInt(number string) int {
	value, err := strconv.Atoi(number)
	if err != nil {
		return 0
	}
	return value
}

func ConvertToInt32(number string) int32 {
	return int32(ConvertToInt(number))
}

func ConvertToUint(number string) uint {
	return uint(ConvertToInt(number))
}

func ConvertToInt64(number string) int64 {
	value, err := strconv.ParseInt(number, 10, 0)
	if err != nil {
		return 0
	}
	return value
}

func ConvertToFloat32(number string) float32 {
	value, err := strconv.ParseFloat(number, 32)
	if err != nil {
		return 0
	}
	return float32(value)
}

func ConvertToFloat64(number string) float64 {
	value, err := strconv.ParseFloat(number, 64)
	if err != nil {
		return 0
	}
	return value
}

func GetCallerFrame(skip int) (frame runtime.Frame, ok bool) {
	const skipOffset = 2 // skip GetCallerFrame and Callers

	pc := make([]uintptr, 1)
	numFrames := runtime.Callers(skip+skipOffset, pc)
	if numFrames < 1 {
		return
	}

	frame, _ = runtime.CallersFrames(pc).Next()
	return frame, frame.PC != 0
}

func SFormatMessage(template string, fmtArgs []interface{}) string {
	if len(fmtArgs) == 0 {
		return template
	}

	if template != "" {
		return fmt.Sprintf(template, fmtArgs...)
	}

	if len(fmtArgs) == 1 {
		if str, ok := fmtArgs[0].(string); ok {
			return str
		}
	}
	return fmt.Sprint(fmtArgs...)
}

func FormatMessage(template string, args ...interface{}) string {
	return SFormatMessage(template, args)
}

func SplitString(s string, sep string) []string {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return []string{}
	}
	s = strings.Trim(s, sep)
	return strings.Split(s, sep)
}

func SplitMultiLineString(is string, removeEmpty bool) []string {
	var r []string
	s := strings.Split(is, "\r")
	for _, n := range s {
		nr := strings.Split(n, "\n")
		if len(nr) > 0 {
			r = append(r, nr...)
		}
	}
	if removeEmpty {
		var re []string
		for _, n := range r {
			nr := strings.TrimSpace(n)
			if len(nr) > 0 {
				re = append(re, nr)
			}
		}
		return re
	}
	return r
}

func ParseTranscoderProgress(s string, prefix string) (int32, float32, bool) {
	progress := int32(0)
	speed := float32(0.0)
	found := false
	i := strings.Index(s, prefix)
	if i > 0 {
		progressStr := s[0:i]
		k := strings.LastIndex(progressStr, "[")
		if k > 0 {
			progressStr = progressStr[k+1:]
			progressStr = strings.TrimSpace(progressStr)
			progressStr = strings.ReplaceAll(progressStr, "%", "")
			progress = ConvertToInt32(progressStr)
			speedStr := s[i+len(prefix):]
			k = strings.Index(speedStr, "]")
			if k > 0 {
				speedStr = speedStr[:k]
				speedStr = strings.TrimSpace(speedStr)
				speed = ConvertToFloat32(speedStr)
				found = true
			}
		}
	}
	return progress, speed, found
}

func GetUrlImgBase64(path string) (baseImg string, err error) {
	//获取本地文件
	file, err := os.Open(path)
	if err != nil {
		err = errors.New("获取本地图片失败")
		return
	}
	defer file.Close()
	imgByte, _ := io.ReadAll(file)
	baseImg = base64.StdEncoding.EncodeToString(imgByte)
	return
}
