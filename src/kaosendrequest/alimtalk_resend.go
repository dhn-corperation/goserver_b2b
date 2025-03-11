package kaosendrequest

import (
	"fmt"
	"time"
	"sync"
	"context"
	"strconv"
	"database/sql"
	"encoding/json"
	s "strings"

	kakao "mycs/src/kakaojson"
	config "mycs/src/kaoconfig"
	databasepool "mycs/src/kaodatabasepool"
	cm "mycs/src/kaocommon"
	krt "mycs/src/kaoresulttable"
)

func AlimtalkResendProc(ctx context.Context) {
	procCnt := 0
	config.Stdlog.Println("Alimtalk 9999 resend - 프로세스 시작 됨 ")

	for {
		if procCnt < 3 {
			select {
			case <- ctx.Done():
			    config.Stdlog.Println("Alimtalk 9999 resend - process가 10초 후에 종료 됨.")
			    time.Sleep(10 * time.Second)
			    config.Stdlog.Println("Alimtalk 9999 resend - process 종료 완료")
			    return
			default:
				tx, err := databasepool.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
				if err != nil {
					config.Stdlog.Println("Alimtalk 9999 resend 트랜잭션 초기화 실패 : ", err)
					continue
				}

				var startNow = time.Now()
				var group_no = fmt.Sprintf("%02d%02d%02d%09d", startNow.Hour(), startNow.Minute(), startNow.Second(), startNow.Nanosecond()) + strconv.Itoa(procCnt)

				updateRows, err := tx.Exec("update DHN_REQUEST_AT_RESEND as a join (select id from DHN_REQUEST_AT_RESEND where send_group is null limit ?) as b on a.id = b.id set send_group = ?", strconv.Itoa(config.Conf.SENDLIMIT), group_no)

				if err != nil {
					config.Stdlog.Println("Alimtalk 9999 resend send_group update 오류 : ", err)
					tx.Rollback()
					continue
				}
				rowCount, err := updateRows.RowsAffected()

				if err != nil {
					config.Stdlog.Println("Alimtalk 9999 resend RowsAffected 확인 오류 : ", err)
					tx.Rollback()
					continue
				}

				if rowCount == 0 {
					tx.Rollback()
					time.Sleep(10 * time.Second)
					continue
				}
				if err := tx.Commit(); err != nil {
					config.Stdlog.Println("Alimtalk 9999 resend tx Commit 오류 : ", err)
					tx.Rollback()
					continue
				}

				procCnt++
				config.Stdlog.Println("Alimtalk 9999 resend 발송 처리 시작 ( ", group_no, " ) : ", rowCount, " 건  ( Proc Cnt :", procCnt, ") - START")

				go func() {
					defer func() {
						procCnt--
					}()
					atResendProcess(group_no, procCnt)
				}()
			}
		}
	}
}

func atResendProcess(group_no string, pc int) {
	defer func(){
		if r := recover(); r != nil {
			config.Stdlog.Println("Alimtalk 9999 resend - atResendProcess panic error : ", r, " / group_no : ", group_no)
			if err, ok := r.(error); ok {
				if s.Contains(err.Error(), "connection refused") {
					for {
						config.Stdlog.Println("Alimtalk 9999 resend - atResendProcess send ping to DB / group_no : ", group_no)
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

	reqsql := "select * from DHN_REQUEST_AT_RESEND where send_group = '" + group_no + "'"

	reqrows, err := db.Query(reqsql)
	if err != nil {
		errlog.Fatal(err)
	}
	defer reqrows.Close()

	columnTypes, err := reqrows.ColumnTypes()
	if err != nil {
		errlog.Fatal(err)
	}
	count := len(columnTypes)
	initScanArgs := cm.InitDatabaseColumn(columnTypes, count)

	var procCount int
	procCount = 0
	var startNow = time.Now()
	var serial_number = fmt.Sprintf("%04d%02d%02d-", startNow.Year(), startNow.Month(), startNow.Day())

	resinsStrs := []string{}
	resinsValues := []interface{}{}
	resinsQuery := `insert IGNORE into DHN_RESULT(`+atColumnStr+`) values %s`

	resultChan := make(chan krt.ResultStr, config.Conf.SENDLIMIT) // resultStr 은 friendtalk에 정의 됨
	var reswg sync.WaitGroup

	for reqrows.Next() {
		scanArgs := initScanArgs

		err := reqrows.Scan(scanArgs...)
		if err != nil {
			errlog.Println("atResendProcess column scan error : ", err, " / group_no : ", group_no)
			time.Sleep(5 * time.Second)
		}

		var alimtalk kakao.Alimtalk
		var attache kakao.AttachmentB
		var attache2 kakao.AttachmentC
		var link *kakao.Link
		var supplement kakao.Supplement
		var button []kakao.Button
		var quickreply []kakao.Quickreply
		result := map[string]string{}

		for i, v := range columnTypes {

			switch s.ToLower(v.Name()) {
			case "msgid":
				if z, ok := (scanArgs[i]).(*sql.NullString); ok {
					alimtalk.Serial_number = serial_number + z.String
				}

			case "message_type":
				if z, ok := (scanArgs[i]).(*sql.NullString); ok {
					alimtalk.Message_type = s.ToUpper(z.String)
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
						var btn kakao.Button

						json.Unmarshal([]byte(z.String), &btn)
						button = append(button, btn)
					}
				}
			case "supplement":
				if z, ok := (scanArgs[i]).(*sql.NullString); ok {
					if len(z.String) > 0 {
						var qrp []kakao.Quickreply

						json.Unmarshal([]byte(z.String), &qrp)
						for i, _ := range qrp {
							quickreply = append(quickreply, qrp[i])
						}
					}
				}
			case "header":
				if z, ok := (scanArgs[i]).(*sql.NullString); ok {
					alimtalk.Header = z.String
				}
			case "attachments":
				if z, ok := (scanArgs[i]).(*sql.NullString); ok {
					if len(z.String) > 0 {
						json.Unmarshal([]byte(z.String), &attache2)
					}
				}
			case "link":
				if z, ok := (scanArgs[i]).(*sql.NullString); ok {
					if len(z.String) > 0 {
						json.Unmarshal([]byte(z.String), &link)
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

		if s.EqualFold(s.ToUpper(result["message_type"]), "AT") {
			alimtalk.Response_method = "realtime"
		} else if s.EqualFold(s.ToUpper(result["message_type"]), "AI") {
			
			alimtalk.Response_method = conf.RESPONSE_METHOD
			alimtalk.Channel_key = conf.CHANNEL
			
			if  s.EqualFold(conf.RESPONSE_METHOD, "polling") {
				alimtalk.Timeout = 600
			}
			
		}

		attache.Buttons = button
		if attache2.Item_highlights != nil {
			attache.Item_highlights = attache2.Item_highlights
		}

		if attache2.Items != nil {
			attache.Items = attache2.Items
		}

		supplement.Quick_reply = quickreply

		alimtalk.Attachment = attache
		alimtalk.Supplement = supplement

		if link != nil {
			alimtalk.Link = link
		}

		var temp krt.ResultStr
		temp.Result = result
		reswg.Add(1)
		go resendKakaoAlimtalk(&reswg, resultChan, alimtalk, temp)
	}


	reswg.Wait()
	chanCnt := len(resultChan)

	atQmarkStr := cm.GetQuestionMark(atColumn)

	for i := 0; i < chanCnt; i++ {

		resChan := <-resultChan
		result := resChan.Result

		if resChan.Statuscode == 200 {

			var kakaoResp kakao.KakaoResponse
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
			resinsValues = append(resinsValues, result["real_msgid"])
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
			resinsValues = append(resinsValues, 2)
			resinsValues = append(resinsValues, result["remark4"])
			resinsValues = append(resinsValues, result["remark5"])
			resinsValues = append(resinsValues, resdtstr) // res_dt
			resinsValues = append(resinsValues, result["reserve_dt"])

			messageType := s.ToUpper(result["message_type"])

			//result 컬럼 처리
			if s.EqualFold(messageType, "AT") || !s.EqualFold(resCode, "0000") || (s.EqualFold(messageType, "AI") && s.EqualFold(conf.RESPONSE_METHOD, "push")) {

				if s.EqualFold(resCode,"0000") {
					resinsValues = append(resinsValues, "Y") // 
				// 1차 카카오 발송 실패 후 2차 발송을 바로 하기 위해서는 이 조건을 맞춰야함
				} else if len(result["sms_kind"])>=1 {
					resinsValues = append(resinsValues, "P") // sms_kind 가 SMS / LMS / MMS 이면 문자 발송 시도
				} else {
					resinsValues = append(resinsValues, "Y") // 
				} 
				
			} else if s.EqualFold(messageType, "AI") {
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
			resinsValues = append(resinsValues, result["header"])
			resinsValues = append(resinsValues, result["attachments"])
			resinsValues = append(resinsValues, result["link"])

			if len(resinsStrs) >= 500 {
				resinsStrs, resinsValues = cm.InsMsg(resinsQuery, resinsStrs, resinsValues)
			}
		} else {
			stdlog.Println("Alimtalk 9999 resend - 알림톡 서버 처리 오류 !! ( status : ", resChan.Statuscode, " / body : ", string(resChan.BodyData), " )", result["msgid"])
			db.Exec("update DHN_REQUEST_AT_RESEND set send_group = null where msgid = '" + result["msgid"] + "'")
		}

		procCount++
	}

	if len(resinsStrs) > 0 {
		resinsStrs, resinsValues = cm.InsMsg(resinsQuery, resinsStrs, resinsValues)
	}
	
	stdlog.Println("Alimtalk 9999 resend - 발송 처리 완료 ( ", group_no, " ) : ", procCount, " 건  ( Proc Cnt :", pc, ") - END")
	
}

func resendKakaoAlimtalk(reswg *sync.WaitGroup, c chan<- krt.ResultStr, alimtalk kakao.Alimtalk, temp krt.ResultStr) {
	defer reswg.Done()

	var seq int = 1

	for {
		alimtalk.Serial_number = alimtalk.Serial_number + strconv.Itoa(seq)
		resp, err := config.Client.R().
			SetHeaders(map[string]string{"Content-Type": "application/json"}).
			SetBody(alimtalk).
			Post(config.Conf.API_SERVER + "/v3/" + config.Conf.PROFILE_KEY + "/alimtalk/send")

		if err != nil {
			config.Stdlog.Println("Alimtalk 9999 resend - 알림톡 메시지 서버 호출 오류 : ", err)
		} else {
			temp.Statuscode = resp.StatusCode()
			if temp.Statuscode != 500 {
				temp.BodyData = resp.Body()
				break
			}
		}
		seq++
	}
	databasepool.DB.Exec("update DHN_REQUEST_AT_RESEND set try_cnt = " + strconv.Itoa(seq) + " where msgid = '" + temp.Result["msgid"] + "'")
	c <- temp

}
