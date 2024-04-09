package kaoreqreceive

import (
	//"encoding/json"
	"fmt"
	config "kaoconfig"
	databasepool "kaodatabasepool"
	"kaoreqtable"
	"strconv"
	s "strings"
	"time"

	"crypto/aes"
	"crypto/cipher"

	//"crypto/rand"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var SecretKey = "9b4dabe9d4fed126a58f8639846143c7"

//친구톡, 알림톡 공통 컬럼
var atColumn = []string{
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

var ftColumn = atColumn
//친구톡 추가 칼럼
ftColumn = append(ftColumn, "att_items")
ftColumn = append(ftColumn, "att_coupon")

//DHN_RESULT, DHN_RESULT_TEMP 테이블의 컬럼
var msgColumn = []string{
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

ftColumnStr := s.Join(ftColumn, ",")
atColumnStr := s.Join(atColumn, ",")
msgColumnStr := s.Join(msgColumn, ",")


func ReqReceive(c *gin.Context) {
	ctx := c.Request.Context()
	errlog := config.Stdlog

	userid := c.Request.Header.Get("userid")
	userip := c.ClientIP()
	isValidation := false

	// 허가된 userid 인지 테이블에서 확인
	sqlstr := "
		select 
			count(1) as cnt 
		from
			DHN_CLIENT_LIST
		where
			user_id = ?
			and ip = ?
			and use_flag = 'Y'"
	// val, verr := databasepool.DB.Query(sqlstr)

	var cnt int
	err := databasepool.DB.QueryRowContext(ctx, sqlstr, userid, userip).Scan(&cnt)
	if err != nil { errlog.Println(err) }

	// val.Next()
	// val.Scan(&cnt)

	defer val.Close()

	if cnt > 0 { isValidation = true }

	var startNow = time.Now()
	var startTime = fmt.Sprintf("%02d:%02d:%02d", startNow.Hour(), startNow.Minute(), startNow.Second())

	errlog.Println("메세지 발송 정보 수신 시작!! ", startTime)

	if isValidation {

		var msg []kaoreqtable.Reqtable
		//전달온 데이터 kaoreqtable.Reqtable에 맵핑
		err1 := c.ShouldBindJSON(&msg)
		if err1 != nil { errlog.Println(err1) }

		errlog.Println("발송 메세지 수신 ( ", userid, ") : ", len(msg), startTime)

		reqinsStrs := []string{}
		//친구톡 value interface 배열 생성
		reqinsValues := []interface{}{}

		atreqinsStrs := []string{}
		//알림톡 value interface 배열 생성
		atreqinsValues := []interface{}{}

		resinsStrs := []string{}
		//문자 value interface 배열 생성
		resinsValues := []interface{}{}

		//친구톡 insert 컬럼 셋팅
		reqinsQuery := `insert IGNORE into DHN_REQUEST(`+ftColumnStr+`) values %s`

		//알림톡 insert 컬럼 셋팅
		atreqinsQuery := `insert IGNORE into DHN_REQUEST_AT(`+atColumnStr+`) values %s`

		//문자 insert 컬럼 셋팅
		resinsquery := `insert IGNORE into DHN_RESULT(`+msgColumnStr+`) values %s`

		//temp 테이블 컬럼 셋팅(DHN_RESULT_TEMP : 에러 시 데이터 유실을 막기 위한 테이블)
		resinstempquery := `insert IGNORE into DHN_RESULT_TEMP(`+msgColumnStr+`) values %s`

		ftQmarkStr := setQuestionMark(ftColumn)
		atQmarkStr := setQuestionMark(atColumn)
		msgQmarkStr := setQuestionMark(msgColumn)

		//맵핑한 데이터 row 처리
		for i, _ := range msg {
			//친구톡 insert values 만들기
			if s.HasPrefix(s.ToUpper(msg[i].Messagetype), "F") {
				reqinsStrs = append(reqinsStrs, "("+ftQmarkStr+")")
				reqinsValues = append(reqinsValues, msg[i].Msgid)
				reqinsValues = append(reqinsValues, userid)
				reqinsValues = append(reqinsValues, msg[i].Adflag)
				reqinsValues = append(reqinsValues, msg[i].Button1)
				reqinsValues = append(reqinsValues, msg[i].Button2)
				reqinsValues = append(reqinsValues, msg[i].Button3)
				reqinsValues = append(reqinsValues, msg[i].Button4)
				reqinsValues = append(reqinsValues, msg[i].Button5)
				reqinsValues = append(reqinsValues, msg[i].Imagelink)
				reqinsValues = append(reqinsValues, msg[i].Imageurl)
				reqinsValues = append(reqinsValues, msg[i].Messagetype)
				reqinsValues = append(reqinsValues, msg[i].Msg)
				reqinsValues = append(reqinsValues, msg[i].Msgsms)
				reqinsValues = append(reqinsValues, msg[i].Onlysms)
				reqinsValues = append(reqinsValues, msg[i].Phn)
				reqinsValues = append(reqinsValues, msg[i].Profile)
				reqinsValues = append(reqinsValues, msg[i].Pcom)
				reqinsValues = append(reqinsValues, msg[i].Pinvoice)
				reqinsValues = append(reqinsValues, msg[i].Regdt)
				reqinsValues = append(reqinsValues, msg[i].Remark1)
				reqinsValues = append(reqinsValues, msg[i].Remark2)
				reqinsValues = append(reqinsValues, msg[i].Remark3)
				reqinsValues = append(reqinsValues, msg[i].Remark4)
				reqinsValues = append(reqinsValues, msg[i].Remark5)
				reqinsValues = append(reqinsValues, msg[i].Reservedt)
				reqinsValues = append(reqinsValues, msg[i].Smskind)
				reqinsValues = append(reqinsValues, msg[i].Smslmstit)
				reqinsValues = append(reqinsValues, msg[i].Smssender)
				reqinsValues = append(reqinsValues, msg[i].Scode)
				reqinsValues = append(reqinsValues, msg[i].Tmplid)
				reqinsValues = append(reqinsValues, msg[i].Wide)
				reqinsValues = append(reqinsValues, nil)
				reqinsValues = append(reqinsValues, msg[i].Supplement)
				if len(msg[i].Price) > 0 {
					price, _ := strconv.Atoi(msg[i].Price)
					reqinsValues = append(reqinsValues, price)
				} else {
					reqinsValues = append(reqinsValues, nil)
				}

				reqinsValues = append(reqinsValues, msg[i].Currencytype)
				reqinsValues = append(reqinsValues, msg[i].Title)
				reqinsValues = append(reqinsValues, msg[i].Header)
				reqinsValues = append(reqinsValues, msg[i].Carousel)
				reqinsValues = append(reqinsValues, msg[i].Att_items)
				reqinsValues = append(reqinsValues, msg[i].Att_coupon)
			//문자 insert values 만들기
			} else if s.EqualFold(msg[i].Messagetype, "PH") {
				var resdt = time.Now()
				var resdtstr = fmt.Sprintf("%4d-%02d-%02d %02d:%02d:%02d", resdt.Year(), resdt.Month(), resdt.Day(), resdt.Hour(), resdt.Minute(), resdt.Second())
				resinsStrs = append(resinsStrs, "("+msgQmarkStr+")")
				resinsValues = append(resinsValues, msg[i].Msgid)
				resinsValues = append(resinsValues, userid)
				resinsValues = append(resinsValues, msg[i].Adflag)
				resinsValues = append(resinsValues, msg[i].Button1)
				resinsValues = append(resinsValues, msg[i].Button2)
				resinsValues = append(resinsValues, msg[i].Button3)
				resinsValues = append(resinsValues, msg[i].Button4)
				resinsValues = append(resinsValues, msg[i].Button5)
				resinsValues = append(resinsValues, "9999") // 결과 code
				resinsValues = append(resinsValues, msg[i].Imagelink)
				resinsValues = append(resinsValues, msg[i].Imageurl)
				resinsValues = append(resinsValues, nil) // kind
				resinsValues = append(resinsValues, "")  // 결과 Message
				resinsValues = append(resinsValues, msg[i].Messagetype)

				if s.EqualFold(msg[i].Crypto, "Y") {
					resinsValues = append(resinsValues, AES256GSMDecrypt([]byte(SecretKey), msg[i].Msg, msg[i].Profile))
				} else {
					resinsValues = append(resinsValues, msg[i].Msg)
				}

				if s.EqualFold(msg[i].Crypto, "Y") {
					resinsValues = append(resinsValues, AES256GSMDecrypt([]byte(SecretKey), msg[i].Msgsms, msg[i].Profile))
				} else {
					resinsValues = append(resinsValues, msg[i].Msgsms)
				}
				resinsValues = append(resinsValues, msg[i].Onlysms)
				resinsValues = append(resinsValues, msg[i].Pcom)
				resinsValues = append(resinsValues, msg[i].Pinvoice)

				if s.EqualFold(msg[i].Crypto, "Y") {
					resinsValues = append(resinsValues, AES256GSMDecrypt([]byte(SecretKey), msg[i].Phn, msg[i].Profile))
				} else {
					resinsValues = append(resinsValues, msg[i].Phn)
				}

				if s.EqualFold(msg[i].Crypto, "Y") {
					resinsValues = append(resinsValues, nil)
				} else {
					resinsValues = append(resinsValues, msg[i].Profile)
				}
				resinsValues = append(resinsValues, msg[i].Regdt)
				resinsValues = append(resinsValues, msg[i].Remark1)
				resinsValues = append(resinsValues, msg[i].Remark2)
				resinsValues = append(resinsValues, msg[i].Remark3)
				resinsValues = append(resinsValues, msg[i].Remark4)
				resinsValues = append(resinsValues, msg[i].Remark5)
				resinsValues = append(resinsValues, resdtstr) // res_dt
				resinsValues = append(resinsValues, msg[i].Reservedt)
				resinsValues = append(resinsValues, "P") // sms_kind 가 SMS / LMS / MMS 이면 문자 발송 시도
				resinsValues = append(resinsValues, msg[i].Scode)
				resinsValues = append(resinsValues, msg[i].Smskind)

				if s.EqualFold(msg[i].Crypto, "Y") {
					resinsValues = append(resinsValues, AES256GSMDecrypt([]byte(SecretKey), msg[i].Smslmstit, msg[i].Profile))
				} else {
					resinsValues = append(resinsValues, msg[i].Smslmstit)
				}

				if s.EqualFold(msg[i].Crypto, "Y") {
					resinsValues = append(resinsValues, AES256GSMDecrypt([]byte(SecretKey), msg[i].Smssender, msg[i].Profile))
				} else {
					resinsValues = append(resinsValues, msg[i].Smssender)
				}
				resinsValues = append(resinsValues, "N")
				resinsValues = append(resinsValues, msg[i].Tmplid)
				resinsValues = append(resinsValues, msg[i].Wide)
				resinsValues = append(resinsValues, nil) // send_group
				resinsValues = append(resinsValues, msg[i].Supplement)
				resinsValues = append(resinsValues, nil) //price
				resinsValues = append(resinsValues, nil) //currency_type
				resinsValues = append(resinsValues, msg[i].Header)
				resinsValues = append(resinsValues, msg[i].Carousel)
			//알림톡 insert values 만들기
			} else {
				atreqinsStrs = append(atreqinsStrs, "("+atQmarkStr+")")
				atreqinsValues = append(atreqinsValues, msg[i].Msgid)
				atreqinsValues = append(atreqinsValues, userid)
				atreqinsValues = append(atreqinsValues, msg[i].Adflag)
				atreqinsValues = append(atreqinsValues, msg[i].Button1)
				atreqinsValues = append(atreqinsValues, msg[i].Button2)
				atreqinsValues = append(atreqinsValues, msg[i].Button3)
				atreqinsValues = append(atreqinsValues, msg[i].Button4)
				atreqinsValues = append(atreqinsValues, msg[i].Button5)
				atreqinsValues = append(atreqinsValues, msg[i].Imagelink)
				atreqinsValues = append(atreqinsValues, msg[i].Imageurl)
				atreqinsValues = append(atreqinsValues, msg[i].Messagetype)
				atreqinsValues = append(atreqinsValues, msg[i].Msg)
				atreqinsValues = append(atreqinsValues, msg[i].Msgsms)
				atreqinsValues = append(atreqinsValues, msg[i].Onlysms)
				atreqinsValues = append(atreqinsValues, msg[i].Phn)
				atreqinsValues = append(atreqinsValues, msg[i].Profile)
				atreqinsValues = append(atreqinsValues, msg[i].Pcom)
				atreqinsValues = append(atreqinsValues, msg[i].Pinvoice)
				atreqinsValues = append(atreqinsValues, msg[i].Regdt)
				atreqinsValues = append(atreqinsValues, msg[i].Remark1)
				atreqinsValues = append(atreqinsValues, msg[i].Remark2)
				atreqinsValues = append(atreqinsValues, msg[i].Remark3)
				atreqinsValues = append(atreqinsValues, msg[i].Remark4)
				atreqinsValues = append(atreqinsValues, msg[i].Remark5)
				atreqinsValues = append(atreqinsValues, msg[i].Reservedt)
				atreqinsValues = append(atreqinsValues, msg[i].Smskind)
				atreqinsValues = append(atreqinsValues, msg[i].Smslmstit)
				atreqinsValues = append(atreqinsValues, msg[i].Smssender)
				atreqinsValues = append(atreqinsValues, msg[i].Scode)
				atreqinsValues = append(atreqinsValues, msg[i].Tmplid)
				atreqinsValues = append(atreqinsValues, msg[i].Wide)
				atreqinsValues = append(atreqinsValues, nil) //send_group
				atreqinsValues = append(atreqinsValues, msg[i].Supplement)

				if len(msg[i].Price) > 0 {
					price, _ := strconv.Atoi(msg[i].Price)
					atreqinsValues = append(atreqinsValues, price)
				} else {
					atreqinsValues = append(atreqinsValues, nil)
				}

				atreqinsValues = append(atreqinsValues, msg[i].Currencytype)
				atreqinsValues = append(atreqinsValues, msg[i].Title)
				atreqinsValues = append(atreqinsValues, msg[i].Header)
				atreqinsValues = append(atreqinsValues, msg[i].Carousel)
			}

			// 500건 단위로 처리한다(클라이언트에서 1000건씩 전송하더라도 지정한 단위의 건수로 insert한다.)
			saveCount := 500
			if len(reqinsStrs) >= saveCount {
				reqinsStrs, reqinsValues = setMsg(reqinsQuery, reqinsStrs, reqinsValues)
			}

			if len(atreqinsStrs) >= saveCount {
				atreqinsStrs, atreqinsValues = setMsg(atreqinsQuery, atreqinsStrs, atreqinsValues)
			}

			if len(resinsStrs) >= saveCount {
				resinsStrs, resinsValues = setMsg(resinsquery, resinsStrs, resinsValues, true, resinstempquery)
			}
		}
		
		// 나머지 건수를 저장하기 위해 다시한번 정의
		if len(reqinsStrs) > 0 {
			reqinsStrs, reqinsValues = setMsg(reqinsQuery, reqinsStrs, reqinsValues)
		}

		if len(atreqinsStrs) > 0 {
			atreqinsStrs, atreqinsValues = setMsg(atreqinsQuery, atreqinsStrs, atreqinsValues)
		}

		if len(resinsStrs) > 0 {
			resinsStrs, resinsValues = setMsg(resinsquery, resinsStrs, resinsValues, true, resinstempquery)
		}

		errlog.Println("메세지 발송 정보 수신 끝!! ", startTime)

		c.JSON(200, gin.H{
			"message": "ok",
		})
	} else {
		c.JSON(404, gin.H{
			"code":    "error",
			"message": "허용되지 않은 사용자 입니다",
			"userid":  userid,
			"ip":      userip,
		})
	}
}

//테이블 insert 처리
func setMsg(query string, insStrs []string{}, insValues []interface{}{}, tempFlag bool, tempQuery string) (interface{}, interface{}){
	stmt := fmt.Sprintf(query, s.Join(insStrs, ","))
	_, err := databasepool.DB.Exec(stmt, insValues...)

	if err != nil {
		config.Stdlog.Println("Result Table Insert 처리 중 오류 발생 " + err.Error())
		if tempFlag {
			config.Stdlog.Println("Result Temp Table Insert 시작")
			stmtt := fmt.Sprintf(tempQuery, s.Join(insStrs, ","))
			_, errt := databasepool.DB.Exec(stmtt, insValues...)
			if errt != nil {
				config.Stdlog.Println("Result Temp Table Insert 처리 중 오류 발생 " + errt.Error())
			}
		}
	}
	return nil, nil
}

//물음표 컬럼 개수만큼 조인
func setQuestionMark(column []string) string {
	var placeholders []string
	numPlaceholders := len(column) // 원하는 물음표 수
	for i := 0; i < numPlaceholders; i++ {
	    placeholders = append(placeholders, "?")
	}
	return s.Join(placeholders, ",")
}

//AES 복호화
func AES256GSMDecrypt(secretKey []byte, ciphertext_ string, nonce_ string) string {

	ciphertext, _ := ConvertByte(ciphertext_)
	nonce, _ := ConvertByte(nonce_)

	if len(secretKey) != 32 {
		return ""
	}

	// prepare AES-256-GSM cipher
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return ""
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return ""
	}

	// decrypt ciphertext
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return ""
	}

	return string(plaintext)
}

//바이트 생성
func ConvertByte(src string) ([]byte, error) {
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
