package util

import (
	"bytes"
	"com.banxiaoxiao.server/config"
	"com.banxiaoxiao.server/logger"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

func HttpPost(url string, content string) error {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(content)))
	if err != nil {
		logger.Log().Errorf("%s", err.Error())
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: config.GetHttpTimeout(),
	}
	resp, err := client.Do(req)
	if err != nil {
		logger.Log().Infof("%s", err.Error())
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return errors.New(resp.Status)
}

func HttpRequest(url string, post bool, content []byte, contentType string, timeout time.Duration) ([]byte, error) {
	var req *http.Request
	var err error
	if post {
		if content != nil {
			req, err = http.NewRequest("POST", url, bytes.NewBuffer(content))
		} else {
			req, err = http.NewRequest("POST", url, nil)
		}
	} else {
		req, err = http.NewRequest("GET", url, nil)
	}
	if err != nil {
		logger.Log().Infof("%s", err.Error())
		return nil, err
	}
	if post {
		req.Header.Set("Content-Type", contentType)
	}
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		logger.Log().Infof("%s", err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		return body, err
	}
	return nil, errors.New(resp.Status)
}

func HttpRequestSkipSSL(url string, post bool, content []byte, contentType string, timeout time.Duration) ([]byte, error) {
	var req *http.Request
	var err error
	if post {
		if content != nil {
			req, err = http.NewRequest("POST", url, bytes.NewBuffer(content))
		} else {
			req, err = http.NewRequest("POST", url, nil)
		}
	} else {
		req, err = http.NewRequest("GET", url, nil)
	}
	if err != nil {
		logger.Log().Infof("%s", err.Error())
		return nil, err
	}
	if post {
		req.Header.Set("Content-Type", contentType)
	}
	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			DisableKeepAlives: true,
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		logger.Log().Infof("%s", err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		return body, err
	}
	return nil, errors.New(resp.Status)
}

func HttpRequestSkipSSLHeaderCallBack(url string, post bool, content []byte, contentType string, timeout time.Duration, headerCallback func(header http.Header)) ([]byte, error) {
	var req *http.Request
	var err error
	if post {
		if content != nil {
			req, err = http.NewRequest("POST", url, bytes.NewBuffer(content))
		} else {
			req, err = http.NewRequest("POST", url, nil)
		}
	} else {
		req, err = http.NewRequest("GET", url, nil)
	}
	if err != nil {
		logger.Log().Infof("%s", err.Error())
		return nil, err
	}
	if post {
		req.Header.Set("Content-Type", contentType)
	}
	if headerCallback != nil {
		headerCallback(req.Header)
	}
	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			DisableKeepAlives: true,
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		logger.Log().Infof("%s", err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		return body, err
	}
	return nil, errors.New(resp.Status)
}

func HttpAddParam(address string, k string, v string) (string, error) {
	u, err := url.Parse(address)
	if err != nil {
		return address, err
	}
	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return address, err
	}
	for key, _ := range q {
		if key == k {
			return address, nil
		}
	}

	q.Add(k, v)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func HttpCommand(url string, method string, content []byte, timeout time.Duration, headerCallback func(header http.Header)) ([]byte, error) {
	var req *http.Request
	var err error
	if content != nil {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(content))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return nil, err
	}
	if headerCallback != nil {
		headerCallback(req.Header)
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			DisableKeepAlives: true,
		},
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusOK {
		return body, err
	}
	return nil, errors.New(resp.Status + " " + string(body))
}

func HttpDownload(url string, timeout time.Duration, headerCallback func(header http.Header)) ([]byte, string, error) {
	var req *http.Request
	var err error
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", err
	}
	if headerCallback != nil {
		headerCallback(req.Header)
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			DisableKeepAlives: true,
		},
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusOK {
		return body, resp.Header.Get("Content-Disposition"), err
	}
	return nil, "", errors.New(resp.Status + " " + string(body))
}

func RequestByte(url string, method string, content []byte, timeout time.Duration, headerCallback func(header http.Header), skipSSL bool) ([]byte, error) {
	var req *http.Request
	var err error
	if content != nil {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(content))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return nil, err
	}
	if headerCallback != nil {
		headerCallback(req.Header)
	}
	client := &http.Client{
		Timeout: timeout,
	}
	if skipSSL {
		client.Transport = &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			DisableKeepAlives: true,
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		return body, err
	}
	return nil, errors.New(resp.Status + " " + string(body))
}

func RequestTIR[I any, R any](url string, method string, content *I, timeout time.Duration, headerCallback func(header http.Header), skipSSL bool) ([]byte, *R, error) {
	var err error
	var jsonContent []byte
	if content != nil {
		jsonContent, err = json.Marshal(content)
		if err != nil {
			return nil, nil, err
		}
	}
	body, err := RequestByte(url, method, jsonContent, timeout, headerCallback, skipSSL)
	if err != nil {
		return nil, nil, err
	}
	t := new(R)
	err = json.Unmarshal(body, t)
	return body, t, err
}

func RequestTR[R any](url string, method string, content []byte, timeout time.Duration, headerCallback func(header http.Header), skipSSL bool) ([]byte, *R, error) {
	var err error
	body, err := RequestByte(url, method, content, timeout, headerCallback, skipSSL)
	if err != nil {
		return nil, nil, err
	}
	t := new(R)
	err = json.Unmarshal(body, t)
	return body, t, err
}

func RequestTI[I any](url string, method string, content *I, timeout time.Duration, headerCallback func(header http.Header), skipSSL bool) ([]byte, error) {
	var err error
	var jsonContent []byte
	if content != nil {
		jsonContent, err = json.Marshal(content)
		if err != nil {
			return nil, err
		}
	}
	return RequestByte(url, method, jsonContent, timeout, headerCallback, skipSSL)
}
