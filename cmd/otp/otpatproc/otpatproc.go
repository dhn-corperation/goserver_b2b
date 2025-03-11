package otpatproc

import (
	"fmt"
	"sync"
	"time"
	"strconv"
	"context"
	s "strings"
	"database/sql"
	"encoding/json"

	ss "mycs/internal/structs"
	config "mycs/configs"
	databasepool "mycs/configs/databasepool"
	cm "mycs/internal/commons"
)

func AlimtalkProc(ctx context.Context) {
	atprocCnt := 0
	config.Stdlog.Println("알림톡 OTP 프로세스 시작 됨 ") 
	for {
		if atprocCnt < 5 {
			
			select {
				case <- ctx.Done():
			    config.Stdlog.Println("Alimtalk OTP process가 10초 후에 종료 됨.")
			    time.Sleep(10 * time.Second)
			    config.Stdlog.Println("Alimtalk OTP process 종료 완료")
			    return
			default:
				var count sql.NullInt64
				cnterr := databasepool.DB.QueryRowContext(ctx, "SELECT count(msgid) AS cnt FROM DHN_REQUEST_AT WHERE (upper(message_type) = 'AO' or upper(message_type) = 'IO') and send_group IS NULL AND IFNULL(reserve_dt,'00000000000000') <= DATE_FORMAT(NOW(), '%Y%m%d%H%i%S')").Scan(&count)
				
				if cnterr != nil && cnterr != sql.ErrNoRows {
					config.Stdlog.Println("Alimtalk OTP DHN_REQUEST Table - select 오류 : " + cnterr.Error())
					time.Sleep(10 * time.Second)
				} else {
					if count.Valid && count.Int64 > 0 {
						var startNow = time.Now()
						var group_no = fmt.Sprintf("%02d%02d%02d%09d", startNow.Hour(), startNow.Minute(), startNow.Second(), startNow.Nanosecond())
						
						updateRows, err := databasepool.DB.ExecContext(ctx, "update DHN_REQUEST_AT set send_group = ? where (upper(message_type) = 'AO' or upper(message_type) = 'IO') and send_group is null and ifnull(reserve_dt,'00000000000000') <= date_format(now(), '%Y%m%d%H%i%S') limit ?", group_no, strconv.Itoa(config.Conf.SENDLIMIT))
				
						if err != nil {
							config.Stdlog.Println("Alimtalk OTP send_group Update error : ", err, " / group_no : ", group_no)
						}
				
						rowcnt, _ := updateRows.RowsAffected()
				
						if rowcnt > 0 {
							go func() {
								atprocCnt++
								config.Stdlog.Println("Alimtalk 발송 처리 시작 ( ", group_no, " ) : ( Proc Cnt :", atprocCnt, ") - START")
								defer func() {
									atprocCnt--
								}()
								atsendProcess(group_no, atprocCnt)
							}()
						}
					} else {
						time.Sleep(50 * time.Millisecond)
					}
				}
			}
		}
	}

}

func atsendProcess(group_no string, pc int) {
	defer func(){
		if r := recover(); r != nil {
			config.Stdlog.Println("Alimtalk OTP atsendProcess OTP panic error : ", r)
			if err, ok := r.(error); ok {
				if s.Contains(err.Error(), "connection refused") {
					for {
						config.Stdlog.Println("Alimtalk OTP atsendProcess OTP send ping to DB")
						err := databasepool.DB.Ping()
						if err == nil {
							break
						}
						time.Sleep(10 * time.Second)
					}
				}
			}
		}
	}()
	atColumn := cm.GetResAtColumn()
	atColumnStr := s.Join(atColumn, ",")

	var db = databasepool.DB
	var conf = config.Conf
	var stdlog = config.Stdlog
	var errlog = config.Stdlog

	reqsql := "select * from DHN_REQUEST_AT where send_group = '" + group_no + "'"

	reqrows, err := db.Query(reqsql)
	if err != nil {
		errlog.Println("Alimtalk OTP atsendProcess select error : ", err, " / group_no : ", group_no, " / query : ", reqsql)
		panic(err)
	}

	columnTypes, err := reqrows.ColumnTypes()
	if err != nil {
		errlog.Println("Alimtalk OTP atsendProcess column init error : ", err, " / group_no : ", group_no)
		time.Sleep(5 * time.Second)
	}
	count := len(columnTypes)
	initScanArgs := cm.InitDatabaseColumn(columnTypes, count)

	var procCount int
	procCount = 0
	var startNow = time.Now()
	var serial_number = fmt.Sprintf("%04d%02d%02d-", startNow.Year(), startNow.Month(), startNow.Day())

	resinsStrs := []string{}
	resinsValues := []interface{}{}
	resinsquery := `insert IGNORE into DHN_RESULT(`+atColumnStr+`) values %s`

	resultChan := make(chan ss.ResultStr, config.Conf.SENDLIMIT)
	var reswg sync.WaitGroup

	for reqrows.Next() {
		scanArgs := initScanArgs

		err := reqrows.Scan(scanArgs...)
		if err != nil {
			errlog.Println("Alimtalk OTP atsendProcess column scan error : ", err, " / group_no : ", group_no)
			time.Sleep(5 * time.Second)
		}

		var alimtalk ss.Alimtalk
		var attache ss.AttachmentB
		var supplement ss.Supplement
		var button []ss.Button
		var quickreply []ss.Quickreply
		result := map[string]string{}

		for i, v := range columnTypes {

			switch s.ToLower(v.Name()) {
			case "msgid":
				if z, ok := (scanArgs[i]).(*sql.NullString); ok {
					alimtalk.Serial_number = serial_number + z.String
				}

			case "message_type":
				if z, ok := (scanArgs[i]).(*sql.NullString); ok {
					mt := ""
					if s.ToUpper(z.String) == "AO" {
						mt = "AT"
					} else if s.ToUpper(z.String) == "IO" {
						mt = "AI"
					}
					alimtalk.Message_type = mt
				}

			case "profile":
				if z, ok := (scanArgs[i]).(*sql.NullString); ok {
					alimtalk.Sender_key = z.String
				}

			case "phn":
				if z, ok := (scanArgs[i]).(*sql.NullString); ok {
					var cPhn string
					if s.HasPrefix(z.String, "0"){
						cPhn = s.Replace(z.String, "0", "82", 1)
					} else {
						cPhn = z.String
					}
					alimtalk.Phone_number = cPhn
				}

			case "tmpl_id":
				if z, ok := (scanArgs[i]).(*sql.NullString); ok {
					alimtalk.Template_code = z.String
				}

			case "msg":
				if z, ok := (scanArgs[i]).(*sql.NullString); ok {
					alimtalk.Message = z.String
				}

			case "price":
				if z, ok := (scanArgs[i]).(*sql.NullInt64); ok {
					alimtalk.Price = z.Int64
				}

			case "currency_type":
				if z, ok := (scanArgs[i]).(*sql.NullString); ok {
					alimtalk.Currency_type = z.String
				}

			case "title":
				if z, ok := (scanArgs[i]).(*sql.NullString); ok {
					alimtalk.Title = z.String
				}

			case "button1":
				fallthrough
			case "button2":
				fallthrough
			case "button3":
				fallthrough
			case "button4":
				fallthrough
			case "button5":
				if z, ok := (scanArgs[i]).(*sql.NullString); ok {
					if len(z.String) > 0 {
						var btn ss.Button

						json.Unmarshal([]byte(z.String), &btn)
						button = append(button, btn)
					}
				}
			case "supplement":
				if z, ok := (scanArgs[i]).(*sql.NullString); ok {
					if len(z.String) > 0 {
						var qrp []ss.Quickreply

						json.Unmarshal([]byte(z.String), &qrp)
						for i, _ := range qrp {
							quickreply = append(quickreply, qrp[i])
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

		if s.EqualFold(s.ToUpper(result["message_type"]), "AO") {
			alimtalk.Response_method = "realtime"
		} else if s.EqualFold(s.ToUpper(result["message_type"]), "IO") {
			
			alimtalk.Response_method = conf.RESPONSE_METHOD
			alimtalk.Channel_key = conf.CHANNEL
			
			if  s.EqualFold(conf.RESPONSE_METHOD, "polling") {
				alimtalk.Timeout = 600
			}
			
		}

		attache.Buttons = button
		supplement.Quick_reply = quickreply

		alimtalk.Attachment = attache
		alimtalk.Supplement = supplement

		var temp ss.ResultStr
		temp.Result = result
		reswg.Add(1)
		go sendKakaoAlimtalk(&reswg, resultChan, alimtalk, temp)
	}
	reswg.Wait()
	chanCnt := len(resultChan)

	atQmarkStr := cm.GetQuestionMark(atColumn)

	for i := 0; i < chanCnt; i++ {

		resChan := <-resultChan
		result := resChan.Result
		if resChan.Statuscode == 200 {
			var kakaoResp ss.KakaoResponse
			json.Unmarshal(resChan.BodyData, &kakaoResp)

			var resdt = time.Now()
			var resdtstr = fmt.Sprintf("%4d-%02d-%02d %02d:%02d:%02d", resdt.Year(), resdt.Month(), resdt.Day(), resdt.Hour(), resdt.Minute(), resdt.Second())
			
			var resCode = kakaoResp.Code
			var resMessage = kakaoResp.Message
			
			if s.EqualFold(resCode, "3005") {
				resCode = "0000"
				resMessage = ""
			} 
			
			resinsStrs = append(resinsStrs, "("+atQmarkStr+")")
			resinsValues = append(resinsValues, result["msgid"])
			resinsValues = append(resinsValues, result["userid"])
			resinsValues = append(resinsValues, result["ad_flag"])
			resinsValues = append(resinsValues, result["button1"])
			resinsValues = append(resinsValues, result["button2"])
			resinsValues = append(resinsValues, result["button3"])
			resinsValues = append(resinsValues, result["button4"])
			resinsValues = append(resinsValues, result["button5"])
			resinsValues = append(resinsValues, resCode) // 결과 code
			resinsValues = append(resinsValues, result["image_link"])
			resinsValues = append(resinsValues, result["image_url"])
			resinsValues = append(resinsValues, nil)//kind
			resinsValues = append(resinsValues, resMessage) // 결과 Message
			resinsValues = append(resinsValues, result["message_type"])
			resinsValues = append(resinsValues, result["msg"])
			resinsValues = append(resinsValues, result["msg_sms"])
			resinsValues = append(resinsValues, result["only_sms"])
			resinsValues = append(resinsValues, result["p_com"])
			resinsValues = append(resinsValues, result["p_invoice"])
			resinsValues = append(resinsValues, result["phn"])
			resinsValues = append(resinsValues, result["profile"])
			resinsValues = append(resinsValues, result["reg_dt"])
			resinsValues = append(resinsValues, result["remark1"])
			resinsValues = append(resinsValues, result["remark2"])
			resinsValues = append(resinsValues, result["remark3"])
			resinsValues = append(resinsValues, result["remark4"])
			resinsValues = append(resinsValues, result["remark5"])
			resinsValues = append(resinsValues, resdtstr) // res_dt
			resinsValues = append(resinsValues, result["reserve_dt"])

			messageType := s.ToUpper(result["message_type"])

			//result 컬럼 처리
			if s.EqualFold(messageType, "AO") || !s.EqualFold(resCode, "0000") || (s.EqualFold(messageType, "IO") && s.EqualFold(conf.RESPONSE_METHOD, "push")) {

				if s.EqualFold(resCode,"0000") {
					resinsValues = append(resinsValues, "Y") // 
				// 1차 카카오 발송 실패 후 2차 발송을 바로 하기 위해서는 이 조건을 맞춰야함
				} else if len(result["sms_kind"])>=1 {
					resinsValues = append(resinsValues, "O") // sms_kind 가 SMS / LMS / MMS 이면 문자 발송 시도
				} else {
					resinsValues = append(resinsValues, "Y") // 
				} 
				
			} else if s.EqualFold(messageType, "IO") {
				resinsValues = append(resinsValues, "N") // result
			}
			resinsValues = append(resinsValues, resCode)
			resinsValues = append(resinsValues, result["sms_kind"])
			resinsValues = append(resinsValues, result["sms_lms_tit"])
			resinsValues = append(resinsValues, result["sms_sender"])
			resinsValues = append(resinsValues, "N") //sync
			resinsValues = append(resinsValues, result["tmpl_id"])
			resinsValues = append(resinsValues, result["wide"])
			resinsValues = append(resinsValues, nil) //send_group
			resinsValues = append(resinsValues, result["supplement"])
			resinsValues = append(resinsValues, result["price"])
			resinsValues = append(resinsValues, result["currency_type"])
			resinsValues = append(resinsValues, result["title"])
			resinsValues = append(resinsValues, result["mms_image_id"])

			//Center에서도 사용하고 있는 함수이므로 공용 라이브러리 생성이 필요함
			if len(resinsStrs) >= 500 {
				resinsStrs, resinsValues = cm.InsMsg(resinsquery, resinsStrs, resinsValues, 0)
			}
		} else {
			stdlog.Println("Alimtalk OTP server process error : ( ", string(resChan.BodyData), " )", result["msgid"])
			db.Exec("update DHN_REQUEST_AT set send_group = null where msgid = '" + result["msgid"] + "'")
		}

		procCount++
	}

	//Center에서도 사용하고 있는 함수이므로 공용 라이브러리 생성이 필요함
	if len(resinsStrs) > 0 {
		resinsStrs, resinsValues = cm.InsMsg(resinsquery, resinsStrs, resinsValues, 0)
	}

	//알림톡 발송 후 DHN_REQUEST_AT 테이블의 데이터는 제거한다.
	db.Exec("delete from DHN_REQUEST_AT where send_group = '" + group_no + "'")

	stdlog.Println("Alimtalk OTP 발송 처리 완료 ( ", group_no, " ) : ", procCount, " 건 ( Proc Cnt :", pc, ")")
}

//카카오 서버에 발송을 요청한다.
func sendKakaoAlimtalk(reswg *sync.WaitGroup, c chan<- ss.ResultStr, alimtalk ss.Alimtalk, temp ss.ResultStr) {
	defer reswg.Done()

	resp, err := config.Client.R().
		SetHeaders(map[string]string{"Content-Type": "application/json"}).
		SetBody(alimtalk).
		Post(config.Conf.API_SERVER + "/v3/" + config.Conf.PROFILE_KEY + "/alimtalk/send")

	if err != nil {
		config.Stdlog.Println("alimtalk OTP server request error : ", err, " / serial_number : ", alimtalk.Serial_number)
	} else {
		temp.Statuscode = resp.StatusCode()
		temp.BodyData = resp.Body()
	}
	c <- temp

}
