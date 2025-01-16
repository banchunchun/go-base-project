package util

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestUuid8(t *testing.T) {
	//publicKey := "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDcGsUIIAINHfRTdMmgGwLrjzfM\nNSrtgIf4EGsNaYwmC1GjF/bMh0Mcm10oLhNrKNYCTTQVGGIxuc5heKd1gOzb7bdT\nnCDPPZ7oV7p1B9Pud+6zPacoqDz2M24vHFWYY2FbIIJh8fHhKcfXNXOLovdVBE7Z\ny682X1+R1lRK8D+vmQIDAQAB\n-----END PUBLIC KEY-----"
	publicKey := "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCBjtOJC/qkgPxgZLGhdxA9bAPqOSRDiDPXoWEsORbWOP7OAydSH66Nht6GYQqaFKWwXhyI6zomR/K+BKbZKpM3Ygeanr65tXWrohxiu65BNIFEp6we0nTt7eMq3h23OKtFvcutkrJzYnpA5MurqcYnASrqy73dNjE2oYAGuk1V5QIDAQAB"
	//token := "EQ01LjnERgW2nGYUNNA4BAqD7Kz1oEft1/rvzZhgMEswaNBY+rJ4qw=="
	channelId := "7a5472eb-20f9-446f-997a-559a72e428b2"
	//rsa := NewRsa(publicKey, "")

	data := "{\"channelId\":\"7a5472eb-20f9-446f-997a-559a72e428b2\",\"text\":\"会议结束那天\",\"timeStamp\":1706074314780,\"token\":\"EQ01LjnERgW2nGYUNNA4BAqD7Kz1oEft1/rvzZhgMEswaNBY+rJ4qw==\"}"

	fmt.Println(data)

	sign, err := RSABase64Encrypt([]byte(data), publicKey)
	if err != nil {
		return
	}

	param := fmt.Sprintf("{\"channelId\":\"%s\",\"sign\":\"%s\"}", channelId, sign)

	url := "http://49.4.80.3:8881/content_securiy_service/louddt/openApi/checkText"
	post, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(param)))
	if err != nil {
		return
	}
	defer post.Body.Close()
	all, err := io.ReadAll(post.Body)
	if err != nil {
		return
	}

	fmt.Println(string(all))

	decrypt, err := RSABase64Decrypt("NddawFUNiOjBvLsNjYC9BYERS1zkLYVlIAYKCzto1/phR0jlIkphw5HAWhHmP7XthwP8CNL1W99C4cnE947yJYtzDXli4VFUUIAD/Y0ibJgLZCNNt7GaCd6wkU+vO24+ILSw/Ychwa4StJ3LXBarRBtBS2r4nM5kRFFgQ6ZQQ+4=", publicKey)
	if err != nil {
		return
	}
	fmt.Println(decrypt)
}

func Test1(t *testing.T) {
	publicKey := "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAj6XfkmAG2bswz+VDq8ZT/hgDsl3pZfuG/eye3Ev9lQYebultI5eg627n2Vd9OGdUPeJcP1PFj3ouE+S4jtOb17zIzli3xKyiFczwIhwbWyY0392ZJnFx9L41kslxpYd4hfjCh/90pbY71ZBa4E8hSIlqYs1rAOJ0ohgWvpRrQTLVYBewhtWbKE+iIgcrjOXzjyCd3PXyhJ5FC2K66P5tvUtczk4vYrqQETor6E6TnU/LSYmtUvo9Nmwb1+GwWhzx3QGl9SQJj7eV7fhcFV5yBaOOccNrdOnkzBJ/s3u1CpQeGlGRc7GTnoDWdfzBKLDPm8q4B6zXbyAbfGW33gw2swIDAQAB"
	decrypt, err := RSABase64RawUrLDecrypt("aBeUO61ytsNoRpwRDL-Yqipv--4YuQZJTVGse-6FxT0oHyEKXc1LXGiqqOcgGC_b2hDDfsVa5Va6mmTBrSLqtw1-VcMReX-W4ln1lNEKRycLExoGUfk79jUMCKwKJjXwzhWdIMJu7uuwIf578iCww9_qV8_SWawE1E-KS7DUaR77yPk701Zg3f1byewfcSa_sF3OqYc1dhOcDJKzijSfBJb1bx4ILmXBP4oK2BLfMwfb9dPPtLKlkS3P5KQYp3mpP-_xGJwdjxnHV3jfT0fMb2whvlrRYBeyhmp0Kfeh1Q8DYl0OBmPnR9fBvROSZjs9SNpVoJVl-LDMFnCTNKOYRw", publicKey)
	if err != nil {
		return
	}
	fmt.Println(decrypt)
}

func Test2(t *testing.T) {
	// 假设这是你的公钥，存储在字符串中
	publicKey := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAj6XfkmAG2bswz+VDq8ZT/hgDsl3pZfuG/eye3Ev9lQYebultI5eg627n2Vd9OGdUPeJcP1PFj3ouE+S4jtOb17zIzli3xKyiFczwIhwbWyY0392ZJnFx9L41kslxpYd4hfjCh/90pbY71ZBa4E8hSIlqYs1rAOJ0ohgWvpRrQTLVYBewhtWbKE+iIgcrjOXzjyCd3PXyhJ5FC2K66P5tvUtczk4vYrqQETor6E6TnU/LSYmtUvo9Nmwb1+GwWhzx3QGl9SQJj7eV7fhcFV5yBaOOccNrdOnkzBJ/s3u1CpQeGlGRc7GTnoDWdfzBKLDPm8q4B6zXbyAbfGW33gw2swIDAQAB
-----END PUBLIC KEY-----`

	// 解码Base64编码的消息
	encMsg := "aBeUO61ytsNoRpwRDL-Yqipv--4YuQZJTVGse-6FxT0oHyEKXc1LXGiqqOcgGC_b2hDDfsVa5Va6mmTBrSLqtw1-VcMReX-W4ln1lNEKRycLExoGUfk79jUMCKwKJjXwzhWdIMJu7uuwIf578iCww9_qV8_SWawE1E-KS7DUaR77yPk701Zg3f1byewfcSa_sF3OqYc1dhOcDJKzijSfBJb1bx4ILmXBP4oK2BLfMwfb9dPPtLKlkS3P5KQYp3mpP-_xGJwdjxnHV3jfT0fMb2whvlrRYBeyhmp0Kfeh1Q8DYl0OBmPnR9fBvROSZjs9SNpVoJVl-LDMFnCTNKOYRw"
	data, err := base64.RawURLEncoding.DecodeString(encMsg)
	if err != nil {
		panic(err)
	}

	// 解析公钥
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		panic("公钥解析失败")
	}
	// Base64解码公钥字符串
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	// 创建RSA公钥对象
	pubInterface, err := x509.ParsePKIXPublicKey(decodedPublicKey)
	if err != nil {
		return
	}
	rsaPublicKey := pubInterface.(*rsa.PublicKey)
	var result string
	for len(data) > 0 {
		decodePart := data[:128]
		plain := RsaPublicDecrypt(rsaPublicKey, decodePart)
		result += string(plain)
		data = data[128:]
	}
	fmt.Printf("解密后的数据: %s\n", result)
}
