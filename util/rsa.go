package util

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/subtle"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"math/big"
)

func RSABase64Encrypt(data []byte, publicKey string) (string, error) {
	// Base64解码公钥字符串
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return "", fmt.Errorf("failed to decode public key: %v", err)
	}

	// 创建RSA公钥对象
	pubInterface, err := x509.ParsePKIXPublicKey(decodedPublicKey)
	if err != nil {
		return "", fmt.Errorf("failed to parse public key: %v", err)
	}
	rsaPublicKey := pubInterface.(*rsa.PublicKey)
	blockLength := rsaPublicKey.N.BitLen()/8 - 11
	if len(data) <= blockLength {
		v15, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, data)
		if err != nil {
			return "", err
		}
		return string(v15), nil
	}
	buffer := bytes.NewBufferString("")
	pages := len(data) / blockLength
	for index := 0; index <= pages; index++ {
		start := index * blockLength
		end := (index + 1) * blockLength
		if index == pages {
			if start == len(data) {
				continue
			}
			end = len(data)
		}
		chunk, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, data[start:end])
		if err != nil {
			return "", err
		}
		buffer.Write(chunk)
	}

	return base64.StdEncoding.EncodeToString(buffer.Bytes()), nil
}

func RSABase64Decrypt(encodedData string, publicKey string) (string, error) {
	//privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	// Base64解码加密后的数据
	decodedData, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return "", fmt.Errorf("failed to decode encrypted data: %v", err)
	}
	// Base64解码公钥字符串
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return "", fmt.Errorf("failed to decode public key: %v", err)
	}
	// 创建RSA公钥对象
	pubInterface, err := x509.ParsePKIXPublicKey(decodedPublicKey)
	if err != nil {
		return "", fmt.Errorf("failed to parse public key: %v", err)
	}
	rsaPublicKey := pubInterface.(*rsa.PublicKey)
	var result string
	for len(decodedData) > 0 {
		decodePart := decodedData[:128]
		plain := RsaPublicDecrypt(rsaPublicKey, decodePart)
		result += string(plain)
		decodedData = decodedData[128:]
	}
	return result, nil
}

// RSABase64RawUrLDecrypt ShowBizAI的RSA解密信息
func RSABase64RawUrLDecrypt(encodedData string, publicKey string) (string, error) {
	// Base64解码加密后的数据
	decodedData, err := base64.RawURLEncoding.DecodeString(encodedData)
	if err != nil {
		return "", fmt.Errorf("failed to decode encrypted data: %v", err)
	}
	// Base64解码公钥字符串
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return "", fmt.Errorf("failed to decode public key: %v", err)
	}
	// 创建RSA公钥对象
	pubInterface, err := x509.ParsePKIXPublicKey(decodedPublicKey)
	if err != nil {
		return "", fmt.Errorf("failed to parse public key: %v", err)
	}
	rsaPublicKey := pubInterface.(*rsa.PublicKey)
	data := decrypt(rsaPublicKey, decodedData)
	return string(data), nil
}
func leftPad(input []byte, size int) (out []byte) {
	n := len(input)
	if n > size {
		n = size
	}
	out = make([]byte, size)
	copy(out[len(out)-n:], input)
	return
}

func decrypt(pub *rsa.PublicKey, sig []byte) []byte {
	k := (pub.N.BitLen() + 7) / 8
	m := new(big.Int)
	c := new(big.Int).SetBytes(sig)
	e := big.NewInt(int64(pub.E))
	m.Exp(c, e, pub.N)
	em := leftPad(m.Bytes(), k)
	firstByteIsZero := subtle.ConstantTimeByteEq(em[0], 0)
	secondByteIsTwo := subtle.ConstantTimeByteEq(em[1], 1)
	lookingForIndex := 1
	index := 0
	for i := 2; i < len(em); i++ {
		equals0 := subtle.ConstantTimeByteEq(em[i], 0)
		index = subtle.ConstantTimeSelect(lookingForIndex&equals0, i, index)
		lookingForIndex = subtle.ConstantTimeSelect(equals0, 0, lookingForIndex)
	}
	validPS := subtle.ConstantTimeLessOrEq(2+8, index)
	valid := firstByteIsZero & secondByteIsTwo & (^lookingForIndex & 1) & validPS
	index = subtle.ConstantTimeSelect(valid, index+1, 0)
	return em[index:]
}

func RsaPublicDecrypt(pubKey *rsa.PublicKey, data []byte) []byte {
	c := new(big.Int)
	m := new(big.Int)
	m.SetBytes(data)
	e := big.NewInt(int64(pubKey.E))
	c.Exp(m, e, pubKey.N)
	out := c.Bytes()
	skip := 0
	for i := 2; i < len(out); i++ {
		if i+1 >= len(out) {
			break
		}
		if out[i] == 0xff && out[i+1] == 0 {
			skip = i + 2
			break
		}
	}
	return out[skip:]
}
