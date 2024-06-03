package ktproc

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"strings"
	"time"
	
	// config "mycs/src/kaoconfig"
)

const (
	URL_DEV = "https://devopenapi.xroshot.com/V1"
	VERSION = "V1"
)

type Message struct {
	apiKey string
	apiPw string
	userKey string
	isDev bool
	apiHost string
}

func NewMessage(apiKey, apiPw, userKey string, isDev bool, apiCenter int) *Message {
	m := &Message{
		apiKey: apiKey,
		apiPw: apiPw,
		userKey: userKey,
		isDev: isDev,
	}
	switch apiCenter {
	case 1:
		m.apiHost = "https://openapi1.xroshot.com/V1"
	case 2:
		m.apiHost = "https://openapi2.xroshot.com/V1"
	case 3:
		m.apiHost = "https://openapis.xroshot.com/V1"
	}

	return m
}

func (m *Message) setHeader(param interface{}, isMulti bool) ([]string, string) {
	datetime := time.Now().Format("20060102150405")
	hashKey := m.apiPw + "_" + datetime
	hashData, _ := json.Marshal(param)
	hash := m.getHash(string(hashData) + m.userKey, hashKey)

	headers := []string{
		"API-KEY: " + m.apiKey,
		"HASH: " + strings.ToUpper(hash),
		"SALT: " + m.userKey,
		"TIMESTAMP: " + datetime,
		"SECRET-KEY: " + m.apiPw,
	}

	if !isMulti {
		headers = append(headers, "Content-Type: application/json; charset=utf-8")
	}

	return headers, datetime
}

func (m *Message) getHash(data, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func (m *Message) getBoundary(datetime string) string {
	hash := md5.Sum([]byte(m.apiPw + "_" + datetime))
	return "--------------------------" + hex.EncodeToString(hash[:])[:25]
}

func (m *Message) ExecSMS(apiUrl string, param SendReqTable) (*http.Response, error) {
	headers, _ := m.setHeader(param, false)
	body, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	return m.requestAPI(apiUrl, headers, body)
}

func (m *Message) SearchResult(apiUrl string, param SearchReqTable) (*http.Response, error) {
	headers, _ := m.setHeader(param, false)
	body, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	return m.requestAPI(apiUrl, headers, body)
}

func (m *Message) ExecMMS(apiUrl string, param SendReqTable, fileParam []string) (*http.Response, error) {
	cr := "\r\n"
	headers, datetime := m.setHeader(param, true)
	boundary := m.getBoundary(datetime)
	headers = append(headers, "Content-Type: multipart/form-data; boundary="+boundary)

	var msgBody bytes.Buffer
	var fileBody bytes.Buffer

	msgBody.WriteString("--" + boundary + cr)
	msgBody.WriteString("Content-Disposition: form-data; name=\"message\"" + cr)
	msgBody.WriteString("Content-Type: application/json; charset=utf-8" + cr + cr)

	msgJson, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	msgBody.WriteString(string(msgJson) + cr)

	for _, val := range fileParam {
		if _, err := os.Stat(val); err == nil {
			file, err := os.Open(val)
			if err != nil {
				return nil, err
			}
			defer file.Close()

			fileMime := mime.TypeByExtension(val)
			fileBody.WriteString("--" + boundary + cr)
			fileBody.WriteString("Content-Disposition: form-data; name=\"file\"; filename=\"" + file.Name() + "\"" + cr)
			fileBody.WriteString("Content-Type: " + fileMime + cr + cr)

			content, err := ioutil.ReadAll(file)
			if err != nil {
				return nil, err
			}
			fileBody.Write(content)
			fileBody.WriteString(cr)
		}
	}

	var body bytes.Buffer
	body.Write(msgBody.Bytes())
	body.Write(fileBody.Bytes())
	body.WriteString("--" + boundary + "--")

	return m.requestAPI(apiUrl, headers, body.Bytes())
}

func (m *Message) requestAPI(apiUrl string, headers []string, body []byte) (*http.Response, error) {
	client := &http.Client{}
	fullUrl := apiUrl
	if !m.isDev{
		fullUrl = m.apiHost + apiUrl
	} else {
		fullUrl = URL_DEV + apiUrl
	}

	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	for _, header := range headers {
		parts := strings.SplitN(header, ": ", 2)
		if len(parts) == 2 {
			req.Header.Set(parts[0], parts[1])
		}
	}

	return client.Do(req)
}
