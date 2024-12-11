package kaocommon

import(
	"fmt"
	"crypto/aes"
	"crypto/cipher"
	config "mycs/src/kaoconfig"
	databasepool "mycs/src/kaodatabasepool"
	s "strings"
	"strconv"
	"database/sql"
	"encoding/hex"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
	"io/ioutil"
	"unicode/utf16"
	"github.com/goccy/go-json"
)

type SpecialCharacter struct {
	OriginHex string
	DestStr string
}

var specialCharacters []SpecialCharacter

func init(){
	
}

func SetCenterResult(code, msg string) ([]byte) {
	res, _ := json.Marshal(map[string]string{
		"code": code,
		"message": msg,
	})
	return res
}

func RemoveWs(msg string) (string, error){
	if specialCharacters == nil {
		rows, err := databasepool.DB.Query("select orgin_hex_code, dest_str from SPECIAL_CHARACTER where enabled = 'Y' and dest_str is not null")
		if err != nil {
			config.Stdlog.Println("특수단어 습득 쿼리 에러1 err : ", err)
		}
		defer rows.Close()
		for rows.Next() {
			var sc SpecialCharacter
			err := rows.Scan(&sc.OriginHex, &sc.DestStr)
			if err != nil {
				config.Stdlog.Println("특수단어 습득 쿼리 에러2 err : ", err)
			}
			specialCharacters = append(specialCharacters, sc)
		}
	}
	vMsg := msg

	for _, sc := range specialCharacters {
		decodedOriginHex, err := hexToUTF8(sc.OriginHex)
		if err != nil {
			return "", err
		}
		vMsg = s.ReplaceAll(vMsg, decodedOriginHex, sc.DestStr)
	}
	return vMsg, nil
}

func hexToUTF8(hexStr string) (string, error) {
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func LengthInEUCKR(msg string) (int, error){
	encoder := korean.EUCKR.NewEncoder()
	euckrBytes, err := ioutil.ReadAll(transform.NewReader(s.NewReader(msg), encoder))
	if err != nil {
		return 0, err
	}
	return len(euckrBytes), nil
}

func LengthInUTF16(msg string) int {
	runes := []rune(msg)
	utf16Encoded := utf16.Encode(runes)

	return len(utf16Encoded) * 2
}

//발송 전 친구톡, 알림톡 공통 컬럼(알림톡 칼럼, 알림톡의 삽입 테이블 DHN_REQUEST_AT)
func getCommonColumn() []string {
	reqColumn := []string{
		"msgid",
		"userid",
		"ad_flag",
		"button1",
		"button2",
		"button3",
		"button4",
		"button5",
		"image_link",
		"image_url",
		"message_type",
		"msg",
		"msg_sms",
		"only_sms",
		"phn",
		"profile",
		"p_com",
		"p_invoice",
		"reg_dt",
		"remark1",
		"remark2",
		"remark3",
		"remark4",
		"remark5",
		"reserve_dt",
		"sms_kind",
		"sms_lms_tit",
		"sms_sender",
		"s_code",
		"tmpl_id",
		"wide",
		"send_group",
		"supplement",
		"price",
		"currency_type",
		"title",
		"mms_image_id",
		"header",
		"attachments",
	}
	return reqColumn
}

func GetReqAtColumn() []string {
	reqAtColumn := getCommonColumn()
	reqAtColumn = append(reqAtColumn, "link")
	return reqAtColumn
}

//발송 전 친구톡 추가 칼럼(사입 테이블 : DHN_REQUEST)
func GetReqFtColumn() []string {
	reqFtColumn := getCommonColumn()
	reqFtColumn = append(reqFtColumn, "carousel")
	reqFtColumn = append(reqFtColumn, "att_items")
	reqFtColumn = append(reqFtColumn, "att_coupon")
	return reqFtColumn
}

//재발송 전 알림톡 추가 칼럼(사입 테이블 : DHN_REQUEST)
func GetResendReqAtColumn() []string {
	resendReqAtColumn := GetReqAtColumn()
	resendReqAtColumn = append(resendReqAtColumn, "real_msgid")
	return resendReqAtColumn
}

//재발송 전 알림톡 추가 칼럼(사입 테이블 : DHN_REQUEST)
func GetResendReqFtColumn() []string {
	resendReqFtColumn := GetReqFtColumn()
	resendReqFtColumn = append(resendReqFtColumn, "real_msgid")
	return resendReqFtColumn
}

//발송 전 메시지 삽입 데이터 컬럼(삽입 테이블 DHN_RESULT)
func GetReqMsgColumn() []string {
	msgColumn := []string{
		"msgid",
		"userid",
		"ad_flag",
		"button1",
		"button2",
		"button3",
		"button4",
		"button5",
		"code",
		"image_link",
		"image_url",
		"kind",
		"message",
		"message_type",
		"msg",
		"msg_sms",
		"only_sms",
		"p_com",
		"p_invoice",
		"phn",
		"profile",
		"reg_dt",
		"remark1",
		"remark2",
		"remark3",
		"remark4",
		"remark5",
		"res_dt",
		"reserve_dt",
		"result",
		"s_code",
		"sms_kind",
		"sms_lms_tit",
		"sms_sender",
		"sync",
		"tmpl_id",
		"wide",
		"send_group",
		"supplement",
		"price",
		"currency_type",
		"header",
		"carousel",
		"mms_image_id",
	}
	return msgColumn
}

func GetResAtColumn() []string {
	atResColumn := []string{
		"msgid",
		"userid",
		"ad_flag",
		"button1",
		"button2",
		"button3",
		"button4",
		"button5",
		"code",
		"image_link",
		"image_url",
		"kind",
		"message",
		"message_type",
		"msg",
		"msg_sms",
		"only_sms",
		"p_com",
		"p_invoice",
		"phn",
		"profile",
		"reg_dt",
		"remark1",
		"remark2",
		"remark3",
		"remark4",
		"remark5",
		"res_dt",
		"reserve_dt",
		"result",
		"s_code",
		"sms_kind",
		"sms_lms_tit",
		"sms_sender",
		"sync",
		"tmpl_id",
		"wide",
		"send_group",
		"supplement",
		"price",
		"currency_type",
		"title",
		"mms_image_id",
		"header",
		"attachments",
		"link",
	}
	return atResColumn
}

func GetResFtColumn() []string {
	ftResColumn := []string{
		"msgid",
		"userid",
		"ad_flag",
		"button1",
		"button2",
		"button3",
		"button4",
		"button5",
		"code",
		"image_link",
		"image_url",
		"kind",
		"message",
		"message_type",
		"msg",
		"msg_sms",
		"only_sms",
		"p_com",
		"p_invoice",
		"phn",
		"profile",
		"reg_dt",
		"remark1",
		"remark2",
		"remark3",
		"remark4",
		"remark5",
		"res_dt",
		"reserve_dt",
		"result",
		"s_code",
		"sms_kind",
		"sms_lms_tit",
		"sms_sender",
		"sync",
		"tmpl_id",
		"wide",
		"send_group",
		"supplement",
		"price",
		"currency_type",
		"header",
		"carousel",
		"mms_image_id",
	}
	return ftResColumn
}

//물음표 컬럼 개수만큼 조인
func GetQuestionMark(column []string) string {
	var placeholders []string
	numPlaceholders := len(column) // 원하는 물음표 수
	for i := 0; i < numPlaceholders; i++ {
	    placeholders = append(placeholders, "?")
	}
	return s.Join(placeholders, ",")
}

//테이블 insert 처리
func InsMsg(query string, insStrs []string, insValues []interface{}) ([]string, []interface{}){
	var errlog = config.Stdlog
	stmt := fmt.Sprintf(query, s.Join(insStrs, ","))
	_, err := databasepool.DB.Exec(stmt, insValues...)

	if err != nil {
		errlog.Println("Result Table Insert 처리 중 오류 발생 ", err.Error())
		errlog.Println("table : ", query)
	}
	return nil, nil
}

func InsMsgTemp(query string, insStrs []string, insValues []interface{}, tempFlag bool, tempQuery string) ([]string, []interface{}){
	var errlog = config.Stdlog
	stmt := fmt.Sprintf(query, s.Join(insStrs, ","))
	_, err := databasepool.DB.Exec(stmt, insValues...)

	if err != nil {
		errlog.Println("Result Table Insert 처리 중 오류 발생 ", err.Error())
		errlog.Println("table : ", query)
		if tempFlag {
			errlog.Println("Result Temp Table Insert 시작")
			stmtt := fmt.Sprintf(tempQuery, s.Join(insStrs, ","))
			_, errt := databasepool.DB.Exec(stmtt, insValues...)
			if errt != nil {
				errlog.Println("Result Temp Table Insert 처리 중 오류 발생 ", errt.Error())
			}
		}
	}
	return nil, nil
}

//AES 복호화
func AES256GSMDecrypt(secretKey []byte, ciphertext_ string, nonce_ string) string {
	var errlog = config.Stdlog
	ciphertext, _ := convertByte(ciphertext_)
	nonce, _ := convertByte(nonce_)

	if len(secretKey) != 32 {
		return ""
	}

	// prepare AES-256-GSM cipher
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		errlog.Println("암호화 블록 초기화 실패 : ", err.Error())
		return ""
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		errlog.Println("GCM 암호화기를 초기화 실패 : ", err.Error())
		return ""
	}

	// decrypt ciphertext
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		errlog.Println("복호화 실패 : ", err.Error())
		return ""
	}

	return string(plaintext)
}

//바이트 생성
func convertByte(src string) ([]byte, error) {
	ba := make([]byte, len(src)/2)
	idx := 0
	for i := 0; i < len(src); i = i + 2 {
		b, err := strconv.ParseInt(src[i:i+2], 16, 10)
		if err != nil {
			return nil, err
		}
		ba[idx] = byte(b)
		idx++
	}

	return ba, nil
}

//데이터베이스 default 값 초기화
func InitDatabaseColumn(columnTypes []*sql.ColumnType, length int) []interface{} {
	scanArgs := make([]interface{}, length)

	for i, v := range columnTypes {

		switch v.DatabaseTypeName() {
		case "VARCHAR", "TEXT", "UUID", "TIMESTAMP":
			scanArgs[i] = new(sql.NullString)
			break
		case "BOOL":
			scanArgs[i] = new(sql.NullBool)
			break
		case "INT4":
			scanArgs[i] = new(sql.NullInt64)
			break
		default:
			scanArgs[i] = new(sql.NullString)
		}
	}

	return scanArgs
}

func RemoveValueInPlace(slice []string, value string) []string {
    i := 0
    for _, v := range slice {
        if v != value {
            slice[i] = v
            i++
        }
    }
    return slice[:i]
}