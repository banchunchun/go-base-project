package im

import (
	"encoding/json"
	"strconv"
	"strings"
)

const (
	ClientMagicHeader    = "MVCH"
	DelimiterMagicHeader = "$$"
)

type ImMessage struct {
	Header  *ImHeader `json:"header"`
	Command string    `json:"command"`
	Content string    `json:"content"`
}

type ImHeader struct {
	Header string `json:"header"`
	Length int64  `json:"length"`
	Crc    string `json:"crc"`
}

type LoginRequest struct {
	UserId   string `json:"userId"`
	NickName string `json:"nickName"`
	Avatar   string `json:"avatar"`
	Imei     string `json:"imei"`
}

func BuildMessage(command string, content interface{}) *ImMessage {
	message := &ImMessage{}
	message.Command = command
	v, _ := json.Marshal(content)
	length := len(command) + len(v)
	crcBytes := append([]byte(command), v...)
	message.Content = string(v)
	message.Header = &ImHeader{
		Header: ClientMagicHeader,
		Length: int64(length),
		Crc:    BuildCrc1021Sign(crcBytes),
	}
	return &ImMessage{}
}

func BuildByteMessage(command string, content interface{}) []byte {
	message := &ImMessage{}
	message.Command = command
	v, _ := json.Marshal(content)
	length := len(command) + len(v)
	crcBytes := append([]byte(command), v...)
	message.Content = string(v)
	message.Header = &ImHeader{
		Header: ClientMagicHeader,
		Length: int64(length),
		Crc:    BuildCrc1021Sign(crcBytes),
	}
	r, _ := json.Marshal(message)
	return r
}

func BuildCrc1021Sign(v []byte) string {
	crc := 0xFFFF // 初始值
	polynomial := 0x1021
	for _, b := range v {
		for i := 0; i < 8; i++ {
			bit := (b >> (7 - i) & 1) == 1
			c15 := (crc >> 15 & 1) == 1
			crc <<= 1
			if c15 != bit {
				crc ^= polynomial
			}
		}
	}
	crc &= 0xffff
	sb := ""
	hexStr := strconv.FormatInt(int64(crc), 16)
	length := len(hexStr)
	// 对 CRC16 校验码进行补位运算
	for i := 0; i < 4-length; i++ {
		sb += "0"
	}
	sb += hexStr
	return strings.ToUpper(sb)
}
