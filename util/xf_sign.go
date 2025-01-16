package util

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	// Algorithm 支持的算法
	Algorithm = "hmac-sha256"
	// HttpProto 版本协议
	HttpProto = "HTTP/1.1"
)

func BuildXfSignHeader(req *http.Request, apiKey, secret, host, uri string, body []byte) *http.Request {
	host = strings.ReplaceAll(host, "http://", "")
	host = strings.ReplaceAll(host, "https://", "")
	req.Header.Set("Content-Type", "application/json")
	// 设置请求头 其中Host Date 必须有
	req.Header.Set("Host", host)
	// date必须是utc时区，且不能和服务器时间相差300s
	currentTime := time.Now().UTC().Format(time.RFC1123)
	req.Header.Set("Date", currentTime)
	// 对body进行sha256签名,生成digest头部，POST请求必须对body验证
	digest := "SHA-256=" + signBody(body)
	req.Header.Set("Digest", digest)
	// 根据请求头部内容，生成签名
	sign := generateSignature(host, currentTime, "POST", uri, HttpProto, digest, secret)
	// 组装Authorization头部
	authHeader := fmt.Sprintf(`api_key="%s", algorithm="%s", headers="host date request-line digest", signature="%s"`, apiKey, Algorithm, sign)
	req.Header.Set("Authorization", authHeader)
	return req
}
func generateSignature(host, date, httpMethod, requestUri, httpProto, digest string, secret string) string {
	// 不是request-line的话，则以header名称,后跟ASCII冒号:和ASCII空格，再附加header值
	var signatureStr string
	if len(host) != 0 {
		signatureStr = "host: " + host + "\n"
	}
	signatureStr += "date: " + date + "\n"
	// 如果是request-line的话，则以 http_method request_uri http_proto
	signatureStr += httpMethod + " " + requestUri + " " + httpProto + "\n"
	signatureStr += "digest: " + digest
	return hmacsign(signatureStr, secret)
}
func hmacsign(data, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(data))
	encodeData := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(encodeData)
}
func signBody(data []byte) string {
	// 进行sha256签名
	sha := sha256.New()
	sha.Write(data)
	encodeData := sha.Sum(nil)
	// 经过base64转换
	return base64.StdEncoding.EncodeToString(encodeData)
}

func Md5Hash(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
func HmacSHA1(key, data string) string {
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(data))
	return string(h.Sum(nil))
}
