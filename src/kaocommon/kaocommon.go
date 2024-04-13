package kaocommon

import(
	"crypto/aes"
	"crypto/cipher"
	config "mycs/src/kaoconfig"
	databasepool "mycs/src/kaodatabasepool"
	s "strings"
)

var errlog = config.Stdlog

func init(){
	
}

//친구톡, 알림톡 공통 컬럼(알림톡 칼럼)
func GetAtColumn() []string {
	atColumn := []string{
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
		"header",
		"carousel",
	}
	return atColumn
}

func GetFtColumn() []string {
	ftColumn = SetAtColumn()
	ftColumn = append(ftColumn, "att_items")
	ftColumn = append(ftColumn, "att_coupon")
	return ftColumn
}

//DHN_RESULT, DHN_RESULT_TEMP 테이블의 컬럼
func GetMsgColumn() []string {
	msgColumn = []string{
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
	}
	return msgColumn
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
	stmt := fmt.Sprintf(query, s.Join(insStrs, ","))
	_, err := databasepool.DB.Exec(stmt, insValues...)

	if err != nil {
		errlog.Println("Result Table Insert 처리 중 오류 발생 ", err.Error())
		errlog.Println("table : ", query)
	}
	return nil, nil
}

func InsMsgTemp(query string, insStrs []string, insValues []interface{}, tempFlag bool, tempQuery string) ([]string, []interface{}){
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