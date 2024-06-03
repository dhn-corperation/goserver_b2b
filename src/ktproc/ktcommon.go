package ktproc

import (
	"crypto/hmac"
	"crypto/sha256"
    "encoding/hex"
    "strings"
    "fmt"
    "encoding/json"
    "time"

    config "mycs/src/kaoconfig"
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
}

func (m *Message) setHeader(param SendReqTable, isMulti bool) []string {
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

	return headers
}

func (m *Message) getHash(data, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return hex.EncodeToString(h.sum(nil))
}

func (m *Message) getBoundary() string {
	hash := md5.Sum([]byte(m.getHashKey()))
	return "--------------------------" + hex.EncodeToString(hash[:])[:25]
}

func (m *Message) execSMS(apiUrl string, param SendReqTable) (*http.Response, error) {
	headers := m.setHeader(param, false)
	body, err := json.Marchal(param)
	if err != nil {
		return nil, err
	}
	return m.requestAPI(apiUrl, headers, body)
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
			req.Header.set(parts[0], parts[1])
		}
	}

	return client.Do(req)
}

































