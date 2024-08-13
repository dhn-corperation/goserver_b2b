package kaosendrequest

import (
	//"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	kakao "kakaojson"
	config "kaoconfig"
	databasepool "kaodatabasepool"

	//"io/ioutil"
	//"net"
	//"net/http"
	s "strings"

	"strconv"
	"sync"
	"time"

	//"github.com/go-resty/resty"
	r "reflect"
)

var ftprocCnt int
var FisRunning bool
var isStoping bool

//var limitcnt int = config.Conf.SENDLIMIT

type resultStr struct {
	Statuscode int
	BodyData   []byte
	Result     map[string]string
}

func FriendtalkProc() {
	ftprocCnt = 1

	for {
		if ftprocCnt <= 5 {
			var count int

			cnterr := databasepool.DB.QueryRow("select length(msgid) as cnt from DHN_REQUEST  where send_group is null and ifnull(reserve_dt,'00000000000000') <= date_format(now(), '%Y%m%d%H%i%S') limit 1").Scan(&count)

			if cnterr != nil {
				//config.Stdlog.Println("DHN_REQUEST Table - select 오류 : " + cnterr.Error())
			} else {

				if count > 0 {
					var startNow = time.Now()
					var group_no = fmt.Sprintf("%02d%02d%02d%09d", startNow.Hour(), startNow.Minute(), startNow.Second(), startNow.Nanosecond())

					updateRows, err := databasepool.DB.Exec("update DHN_REQUEST set send_group = '" + group_no + "' where send_group is null and ifnull(reserve_dt,'00000000000000') <= date_format(now(), '%Y%m%d%H%i%S') limit " + strconv.Itoa(config.Conf.SENDLIMIT))

					if err != nil {
						config.Stdlog.Println("Request Table - send_group Update 오류")
					}

					rowcnt, _ := updateRows.RowsAffected()

					if rowcnt > 0 {
						config.Stdlog.Println("친구톡 발송 처리 시작 ( ", group_no, " ) : ", rowcnt, " 건 ")
						ftprocCnt++
						go ftsendProcess(group_no)
					}
				}
			}
		}
	}

}

func ftsendProcess(group_no string) {

	var db = databasepool.DB
	var conf = config.Conf
	var stdlog = config.Stdlog
	var errlog = config.Stdlog

	reqsql := "select * from DHN_REQUEST where send_group = '" + group_no + "' and message_type like 'f%' "

	reqrows, err := db.Query(reqsql)
	if err != nil {
		errlog.Fatal(err)
	}

	columnTypes, err := reqrows.ColumnTypes()
	if err != nil {
		errlog.Fatal(err)
	}
	count := len(columnTypes)

	var procCount int
	procCount = 0
	var startNow = time.Now()
	var serial_number = fmt.Sprintf("%04d%02d%02d-", startNow.Year(), startNow.Month(), startNow.Day())

	resinsStrs := []string{}
	resinsValues := []interface{}{}
	resinsquery := `insert IGNORE into DHN_RESULT(
msgid ,
userid ,
ad_flag ,
button1 ,
button2 ,
button3 ,
button4 ,
button5 ,
code ,
image_link ,
image_url ,
kind ,
message ,
message_type ,
msg ,
msg_sms ,
only_sms ,
p_com ,
p_invoice ,
phn ,
profile ,
reg_dt ,
remark1 ,
remark2 ,
remark3 ,
remark4 ,
remark5 ,
res_dt ,
reserve_dt ,
result ,
s_code ,
sms_kind ,
sms_lms_tit ,
sms_sender ,
sync ,
tmpl_id ,
wide ,
send_group ,
supplement ,
price ,
currency_type,
header,
carousel,
attachments      
) values %s`

	// friendClient := &http.Client{
	// 	Timeout: time.Second * 20,
	// }

	resultChan := make(chan resultStr, config.Conf.SENDLIMIT)
	var reswg sync.WaitGroup

	for reqrows.Next() {
		scanArgs := make([]interface{}, count)

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

		err := reqrows.Scan(scanArgs...)
		if err != nil {
			errlog.Fatal(err)
		}

		var friendtalk kakao.Friendtalk

		var attachment kakao.FTAttachment
		var carousel kakao.FTCarousel

		result := map[string]string{}

		for i, v := range columnTypes {

			switch s.ToLower(v.Name()) {
			case "msgid":
				if z, ok := (scanArgs[i]).(*sql.NullString); ok {
					friendtalk.Serial_number = serial_number + z.String
				}

			case "message_type":
				if z, ok := (scanArgs[i]).(*sql.NullString); ok {
					friendtalk.Message_type = s.ToUpper(z.String)
				}

			case "profile":
				if z, ok := (scanArgs[i]).(*sql.NullString); ok {
					friendtalk.Sender_key = z.String
				}

			case "phn":
				if z, ok := (scanArgs[i]).(*sql.NullString); ok {
					friendtalk.Phone_number = z.String
				}

			case "msg":
				if z, ok := (scanArgs[i]).(*sql.NullString); ok {
					friendtalk.Message = z.String
				}

			case "ad_flag":
				if z, ok := (scanArgs[i]).(*sql.NullString); ok {
					friendtalk.Ad_flag = z.String
				}

			case "header":
				if z, ok := (scanArgs[i]).(*sql.NullString); ok {
					friendtalk.Header = z.String
				}

			case "user_key":
				if z, ok := (scanArgs[i]).(*sql.NullString); ok {
					friendtalk.User_key = z.String
				}

			case "carousel":
				if z, ok := (scanArgs[i]).(*sql.NullString); ok {

					if len(z.String) > 0 {
						json.Unmarshal([]byte(z.String), &carousel)
						if s.EqualFold(conf.DEBUG, "Y") {
							jsonstr, _ := json.Marshal(carousel)
							stdlog.Println("Carousel")
							stdlog.Println(z.String)
							stdlog.Println(string(jsonstr))
						}
					}
				}
			case "attachments":
				if z, ok := (scanArgs[i]).(*sql.NullString); ok {

					if len(z.String) > 0 {
						json.Unmarshal([]byte(z.String), &attachment)

						if s.EqualFold(conf.DEBUG, "Y") {
							jsonstr, _ := json.Marshal(attachment)
							stdlog.Println("Attachment")
							stdlog.Println(z.String)
							stdlog.Println(string(jsonstr))
						}
					}
				}

			}

			if z, ok := (scanArgs[i]).(*sql.NullString); ok {
				result[s.ToLower(v.Name())] = z.String
			}

			if z, ok := (scanArgs[i]).(*sql.NullInt32); ok {
				result[s.ToLower(v.Name())] = string(z.Int32)
			}

			if z, ok := (scanArgs[i]).(*sql.NullInt64); ok {
				result[s.ToLower(v.Name())] = string(z.Int64)
			}

		}

		if !r.ValueOf(attachment).IsZero() && attachment.Image != nil && len(attachment.Image.ImgURL) > 0 {
			if s.EqualFold(result["wide"], "Y") {
				friendtalk.Message_type = "FW"
			} else {
				friendtalk.Message_type = "FI"
			}
		} else {
			friendtalk.Message_type = "FT"
		}

		result["message_type"] = friendtalk.Message_type

		if !r.ValueOf(attachment).IsZero() {
			friendtalk.Attachment = &attachment
		}

		if !r.ValueOf(carousel).IsZero() {
			friendtalk.Carousel = &carousel
		}

		if s.EqualFold(conf.DEBUG, "Y") {
			jsonstr, _ := json.Marshal(friendtalk)
			stdlog.Println(string(jsonstr))
		}

		var temp resultStr
		temp.Result = result

		reswg.Add(1)
		go sendKakao(&reswg, resultChan, friendtalk, temp)

	}
	reswg.Wait()
	//fmt.Println("Size :", len(resultChan))
	chanCnt := len(resultChan)
	for i := 0; i < chanCnt; i++ {
		resChan := <-resultChan
		//resp := resChan.Statuscode
		result := resChan.Result
		//fmt.Println(result["msgid"], " / ", len(resultChan), " / ", i, "/ ", resChan.Statuscode)
		if resChan.Statuscode == 200 {
			//str := string(resChan.BodyData)
			var kakaoResp kakao.KakaoResponse
			json.Unmarshal(resChan.BodyData, &kakaoResp)
			var resdt = time.Now()
			var resdtstr = fmt.Sprintf("%4d-%02d-%02d %02d:%02d:%02d", resdt.Year(), resdt.Month(), resdt.Day(), resdt.Hour(), resdt.Minute(), resdt.Second())

			var code_ string

			if errcode, ok := config.ErrCode[kakaoResp.Code]; ok {
				code_ = errcode.Key
			} else {

				code_ = "7300"
			}

			resinsStrs = append(resinsStrs, "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
			resinsValues = append(resinsValues, result["msgid"])
			resinsValues = append(resinsValues, result["userid"])
			resinsValues = append(resinsValues, result["ad_flag"])
			resinsValues = append(resinsValues, result["button1"])
			resinsValues = append(resinsValues, result["button2"])
			resinsValues = append(resinsValues, result["button3"])
			resinsValues = append(resinsValues, result["button4"])
			resinsValues = append(resinsValues, result["button5"])
			resinsValues = append(resinsValues, code_) // 결과 code 매핑값
			resinsValues = append(resinsValues, result["image_link"])
			resinsValues = append(resinsValues, result["image_url"])
			resinsValues = append(resinsValues, nil)               // kind
			resinsValues = append(resinsValues, kakaoResp.Message) // 결과 Message
			resinsValues = append(resinsValues, result["message_type"])
			resinsValues = append(resinsValues, result["msg"])
			resinsValues = append(resinsValues, result["msg_sms"])
			resinsValues = append(resinsValues, result["only_sms"])
			resinsValues = append(resinsValues, result["p_com"])
			resinsValues = append(resinsValues, result["p_invoice"])
			resinsValues = append(resinsValues, result["phn"])
			resinsValues = append(resinsValues, result["profile"])
			resinsValues = append(resinsValues, result["reg_dt"])
			resinsValues = append(resinsValues, "KKO")
			resinsValues = append(resinsValues, result["remark2"])
			resinsValues = append(resinsValues, result["remark3"])
			resinsValues = append(resinsValues, result["remark4"])
			resinsValues = append(resinsValues, result["remark5"])
			resinsValues = append(resinsValues, resdtstr) // res_dt
			resinsValues = append(resinsValues, result["reserve_dt"])

			if s.EqualFold(kakaoResp.Code, "0000") {
				resinsValues = append(resinsValues, "Y") //
			} else if len(result["sms_kind"]) >= 1 && s.EqualFold(config.Conf.PHONE_MSG_FLAG, "YES") {
				resinsValues = append(resinsValues, "P") // sms_kind 가 SMS / LMS / MMS 이면 문자 발송 시도
			} else {
				resinsValues = append(resinsValues, "Y") //
			}

			resinsValues = append(resinsValues, kakaoResp.Code) // API Response Code
			resinsValues = append(resinsValues, result["sms_kind"])
			resinsValues = append(resinsValues, result["sms_lms_tit"])
			resinsValues = append(resinsValues, result["sms_sender"])
			resinsValues = append(resinsValues, "N")
			resinsValues = append(resinsValues, result["tmpl_id"])
			resinsValues = append(resinsValues, result["wide"])
			resinsValues = append(resinsValues, nil) // send group
			resinsValues = append(resinsValues, result["supplement"])
			resinsValues = append(resinsValues, result["price"])
			resinsValues = append(resinsValues, result["currency_type"])
			resinsValues = append(resinsValues, result["header"])
			resinsValues = append(resinsValues, result["carousel"])
			resinsValues = append(resinsValues, result["attachments"])

			if len(resinsStrs) >= 500 {
				stmt := fmt.Sprintf(resinsquery, s.Join(resinsStrs, ","))
				//fmt.Println(stmt)
				_, err := databasepool.DB.Exec(stmt, resinsValues...)

				if err != nil {
					stdlog.Println("Result Table Insert 처리 중 오류 발생 " + err.Error())
				}

				resinsStrs = nil
				resinsValues = nil
			}

		} else {
			stdlog.Println("친구톡 서버 처리 오류 : ( ", string(resChan.BodyData), " )", result["msgid"])
			db.Exec("update DHN_REQUEST set send_group = null where msgid = '" + result["msgid"] + "'")
		}
		//}

		//}

		procCount++
	}

	if len(resinsStrs) > 0 {
		stmt := fmt.Sprintf(resinsquery, s.Join(resinsStrs, ","))

		_, err := databasepool.DB.Exec(stmt, resinsValues...)

		if err != nil {
			if s.EqualFold(conf.DEBUG, "Y") {
				stdlog.Println(string(stmt))
				jsonstr, _ := json.Marshal(resinsValues)
				stdlog.Println(string(jsonstr))
			}
			stdlog.Println("Result Table Insert 처리 중 오류 발생 ", err)
		}

		resinsStrs = nil
		resinsValues = nil
	}

	db.Exec("delete from DHN_REQUEST where send_group = '" + group_no + "'")
	stdlog.Println("친구톡 발송 처리 완료 ( ", group_no, " ) : ", procCount, " 건  ( Proc Cnt :", ftprocCnt, ")")

	ftprocCnt--

}

func sendKakao(reswg *sync.WaitGroup, c chan<- resultStr, friendtalk kakao.Friendtalk, temp resultStr) {
	defer reswg.Done()

	resp, err := config.Client.R().
		SetHeaders(map[string]string{"Content-Type": "application/json"}).
		SetBody(friendtalk).
		Post(config.Conf.API_SERVER + "v3/" + config.Conf.PROFILE_KEY + "/friendtalk/send")

	//fmt.Println("SEND :", resp, err)

	if err != nil {
		config.Stdlog.Println("친구톡 메시지 서버 호출 오류 : ", err)
		//	return nil
	} else {
		temp.Statuscode = resp.StatusCode()
		temp.BodyData = resp.Body()
	}
	c <- temp

}
