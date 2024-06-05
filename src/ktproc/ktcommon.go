package ktproc

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
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
	apiKey  string
	apiPw   string
	userKey string
	isDev   bool
	apiHost string
}

func NewMessage(apiKey, apiPw, userKey string, isDev bool, apiCenter int) *Message {
	m := &Message{
		apiKey:  apiKey,
		apiPw:   apiPw,
		userKey: userKey,
		isDev:   isDev,
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
	hash := m.getHash(string(hashData)+m.userKey, hashKey)

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
	if !m.isDev {
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

func KTCode(code string) string {
	errmsg := map[string]string{
		"10000": "7006",
		"10001": "7099",
		"10002": "7011",
		"10003": "703",
		"10004": "7011",
		"10005": "7011",
		"10006": "7011",
		"10007": "7011",
		"10008": "7011",
		"10009": "7011",
		"10010": "7011",
		"10011": "7011",
		"10012": "7011",
		"10013": "7011",
		"10014": "7011",
		"10015": "7011",
		"10016": "7011",
		"10017": "7011",
		"10018": "7011",
		"10019": "7011",
		"10020": "7011",
		"10021": "7060",
		"10022": "7011",
		"10023": "7011",
		"10024": "7011",
		"10025": "7011",
		"10026": "7011",
		"10027": "7011",
		"10028": "7011",
		"10029": "7011",
		"10030": "7011",
		"10031": "7011",
		"10032": "7011",
		"10033": "7093",
		"10034": "7061",
		"10035": "7036",
		"10036": "7011",
		"10037": "7011",
		"10038": "7011",
		"10039": "7098",
		"10050": "7011",
		"10100": "7011",
		"10101": "7015",
		"10102": "7013",
		"10103": "7014",
		"10104": "7014",
		"10105": "7061",
		"10106": "7061",
		"10107": "7022",
		"10108": "7011",
		"10109": "7016",
		"10110": "7011",
		"10111": "7084",
		"10112": "7011",
		"10113": "7011",
		"10114": "7053",
		"10115": "7052",
		"10116": "7051",
		"10117": "7051",
		"10200": "7023",
		"10201": "7023",
		"10202": "7059",
		"10203": "7021",
		"10212": "7021",
		"10253": "7021",
		"12002": "7059",
		"12003": "7059",
		"12107": "7059",
		"14005": "7059",
		"14007": "7059",
		"14301": "7036",
		"14307": "7036",
		"20000": "7059",
		"40000": "7097",
		"40002": "7099",
		"40003": "7099",
		"41000": "7099",
		"41001": "7099",
		"41002": "7097",
		"41003": "7097",
		"42000": "7099",
		"42001": "7095",
		"50000": "7099",
	}
	val, ex := errmsg[code]
	if !ex {
		val = "7011"
	}
	return val
}

func KTCodeMessage(code string) string {
	errmsg := map[string]string{
		"7000": "초기 입력 상태 (default)",
		"7001": "전송 요청 완료(결과수신대기)",
		"7003": "메시지 형식 오류",
		"7005": "휴대폰번호 가입자 없음(미등록)",
		"7006": "전송 성공",
		"7007": "결번(or 서비스 정지)",
		"7008": "단말기 전원 꺼짐",
		"7009": "단말기 음영지역",
		"7010": "단말기내 수신메시지함 FULL로 전송 실패 (구:단말 Busy, 기타 단말문제)",
		"7011": "기타 전송실패",
		"7013": "스팸차단 발신번호",
		"7014": "스팸차단 수신번호",
		"7015": "스팸차단 메시지내용",
		"7016": "스팸차단 기타",
		"7020": "*단말기 서비스 불가",
		"7021": "단말기 서비스 일시정지",
		"7022": "단말기 착신 거절",
		"7023": "단말기 무응답 및 통화중 (busy)",
		"7028": "단말기 MMS 미지원",
		"7029": "기타 단말기 문제",
		"7036": "유효하지 않은 수신번호(망)",
		"7037": "유효하지 않은 발신번호(망)",
		"7050": "이통사 컨텐츠 에러",
		"7051": "이통사 전화번호 세칙 미준수 발신번호",
		"7052": "이통사 발신번호 변작으로 등록된 발신번호",
		"7053": "이통사 번호도용문자 차단서비스에 가입된 발신번호",
		"7054": "이통사 발신번호 기타",
		"7059": "이통사 기타",
		"7060": "컨텐츠 크기 오류(초과 등)",
		"7061": "잘못된 메시지 타입",
		"7069": "컨텐츠 기타",
		"7074": "[Agent] 중복발송 차단 (동일한 수신번호와 메시지 발송 - 기본off, 설정필요)",
		"7075": "[Agent] 발송 Timeout",
		"7076": "[Agent] 유효하지않은 발신번호",
		"7077": "[Agent] 유효하지않은 수신번호",
		"7078": "[Agent] 컨텐츠 오류 (MMS파일없음 등)",
		"7079": "[Agent] 기타",
		"7080": "고객필터링 차단 (발신번호, 수신번호, 메시지 등)",
		"7081": "080 수신거부",
		"7084": "중복발송 차단",
		"7086": "유효하지 않은 수신번호",
		"7087": "유효하지 않은 발신번호",
		"7088": "발신번호 미등록 차단",
		"7089": "시스템필터링 기타",
		"7090": "발송제한 시간 초과",
		"7092": "잔액부족",
		"7093": "월 발송량 초과",
		"7094": "일 발송량 초과",
		"7095": "초당 발송량 초과 (재전송 필요)",
		"7096": "발송시스템 일시적인 부하 (재전송 필요)",
		"7097": "전송 네트워크 오류 (재전송 필요)",
		"7098": "외부발송시스템 장애 (재전송 필요)",
		"7099": "발송시스템 장애 (재전송 필요)",
	}
	val, ex := errmsg[code]
	if !ex {
		val = "기타 오류"
	}
	return val
}
