package kaosendrequest

import (
	"fmt"
	"sync"
	"time"
	"strconv"
	"context"
	"database/sql"
	"encoding/json"
	s "strings"

	kakao "mycs/src/kakaojson"
	config "mycs/src/kaoconfig"
	databasepool "mycs/src/kaodatabasepool"
	cm "mycs/src/kaocommon"
	krt "mycs/src/kaoresulttable"
)

func AlimtalkProc(user_id string, ctx context.Context) {
	atprocCnt := 0
	config.Stdlog.Println(user_id, " - Alimtalk Process 시작 됨.") 
	for {
			
		select {
		case <- ctx.Done():
		    config.Stdlog.Println(user_id, " - Alimtalk process가 10초 후에 종료 됨.")
		    time.Sleep(10 * time.Second)
		    config.Stdlog.Println(user_id, " - Alimtalk process 종료 완료")
		    return
		default:
			tx, err := databasepool.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})

			if err != nil {
				config.Stdlog.Println(user_id, " - Alimtalk init tx : ", err)
				continue
			}

			var startNow = time.Now()
			var group_no = fmt.Sprintf("%02d%02d%02d%09d", startNow.Hour(), startNow.Minute(), startNow.Second(), startNow.Nanosecond()) + strconv.Itoa(atprocCnt)

			updateRows, err := databasepool.DB.Exec("update DHN_REQUEST_AT as a join (select id from DHN_REQUEST_AT where send_group is null and userid = ? and ifnull(reserve_dt,'00000000000000') <= date_format(now(), '%Y%m%d%H%i%S') limit ?) as b on a.id = b.id set send_group = ?", user_id, strconv.Itoa(config.Conf.SENDLIMIT), group_no)

			if err != nil {
				config.Stdlog.Println(user_id, " - Alimtalk send_group update error : ", err)
				tx.Rollback()
				continue
			}

			rowCount, err := updateRows.RowsAffected()

			if err != nil {
				config.Stdlog.Println(user_id, " - Alimtalk RowsAffected error : ", err)
				tx.Rollback()
				continue
			}

			if rowCount == 0 {
				tx.Rollback()
				continue
			}

			if err := tx.Commit(); err != nil {
				config.Stdlog.Println(user_id, " - Alimtalk tx Commit 오류 : ", err)
				tx.Rollback()
				continue
			}

			atprocCnt++
			config.Stdlog.Println(user_id, " - Alimtalk 발송 처리 시작 ( ", group_no, " ) : ", rowCount, " 건  ( Proc Cnt :", atprocCnt, ") - START")

			go func() {
				defer func() {
					atprocCnt--
				}()
				atsendProcess(group_no, user_id, atprocCnt)
			}()
		}
	}

}

func atsendProcess(group_no, user_id string, pc int) {
	defer func(){
		if r := recover(); r != nil {
			config.Stdlog.Println(user_id, " - atsendProcess panic error : ", r, " / group_no : ", group_no, " / userid  : ", user_id)
			if err, ok := r.(error); ok {
				if s.Contains(err.Error(), "connection refused") {
					for {
						config.Stdlog.Println(user_id, " - atsendProcess send ping to DB / group_no : ", group_no, " / userid  : ", user_id)
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

	resendAtColumn := cm.GetResendReqAtColumn()
	resendAtColumnStr := s.Join(resendAtColumn, ",")

	var db = databasepool.DB
	var conf = config.Conf
	var stdlog = config.Stdlog
	var errlog = config.Stdlog

	reqsql := "select * from DHN_REQUEST_AT where send_group = '" + group_no + "' and userid = '" + user_id + "'"

	reqrows, err := db.Query(reqsql)
	if err != nil {
		errlog.Println(user_id, " - atsendProcess select error : ", err, " / group_no : ", group_no, " / userid  : ", user_id," / query : ", reqsql)
		panic(err)
	}

	columnTypes, err := reqrows.ColumnTypes()
	if err != nil {
		errlog.Println(user_id, " - atsendProcess column init error : ", err, " / group_no : ", group_no, " / userid  : ", user_id)
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
	resinsQuery := `insert IGNORE into DHN_RESULT(`+atColumnStr+`) values %s`

	atreqinsStrs := []string{}
	atreqinsValues := []interface{}{}
	atreqinsQuery := `insert IGNORE into DHN_REQUEST_AT_RESEND(`+resendAtColumnStr+`) values %s`

	resultChan := make(chan krt.ResultStr, config.Conf.SENDLIMIT)
	var reswg sync.WaitGroup

	for reqrows.Next() {
		scanArgs := initScanArgs

		err := reqrows.Scan(scanArgs...)
		if err != nil {
			errlog.Println(user_id, " - atsendProcess column scan error : ", err, " / group_no : ", group_no)
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
		go sendKakaoAlimtalk(&reswg, resultChan, alimtalk, temp)
	}

	reswg.Wait()
	chanCnt := len(resultChan)

	atQmarkStr := cm.GetQuestionMark(atColumn)
	resendAtQmarkStr := cm.GetQuestionMark(resendAtColumn)

	nineErrCnt := 0

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
			
		} else if resChan.Statuscode == 500 {

			var kakaoResp2 kakao.KakaoResponse2
			json.Unmarshal(resChan.BodyData, &kakaoResp2)
			
			var resCode = kakaoResp2.Code

			if s.EqualFold(resCode, "9999"){
				nineErrCnt++
				atreqinsStrs, atreqinsValues = insAtErrResend(result, atreqinsStrs, atreqinsValues, resendAtQmarkStr)

				if len(atreqinsStrs) >= 500 {
					stmt := fmt.Sprintf(atreqinsQuery, s.Join(atreqinsStrs, ","))
					_, err := databasepool.DB.Exec(stmt, atreqinsValues...)

					if err != nil {
						stdlog.Println(user_id, " - Alimtalk 9999 resend - Resend Table Insert 처리 중 오류 발생 ", err)
					}

					atreqinsStrs = nil
					atreqinsValues = nil
				}
			}
		} else {
			stdlog.Println(user_id, " - alimtalk server process error : ( ", string(resChan.BodyData), " )", result["msgid"])
			db.Exec("update DHN_REQUEST_AT set send_group = null where msgid = '" + result["msgid"] + "'")
		}

		procCount++
	}

	if len(atreqinsStrs) > 0 {
		stmt := fmt.Sprintf(atreqinsQuery, s.Join(atreqinsStrs, ","))
		_, err := databasepool.DB.Exec(stmt, atreqinsValues...)

		if err != nil {
			stdlog.Println(user_id, " - Alimtalk 9999 resend - Resend Table Insert 처리 중 오류 발생 ", err)
		}
	}

	//Center에서도 사용하고 있는 함수이므로 공용 라이브러리 생성이 필요함
	if len(resinsStrs) > 0 {
		resinsStrs, resinsValues = cm.InsMsg(resinsQuery, resinsStrs, resinsValues)
	}

	//알림톡 발송 후 DHN_REQUEST_AT 테이블의 데이터는 제거한다.
	db.Exec("delete from DHN_REQUEST_AT where send_group = '" + group_no + "'")
	
	stdlog.Println(user_id, " - Alimtalk 발송 처리 완료 ( ", group_no, " ) : ", procCount, " 건 ( Proc Cnt :", pc, ") - END")
	
}

//카카오 서버에 발송을 요청한다.
func sendKakaoAlimtalk(reswg *sync.WaitGroup, c chan<- krt.ResultStr, alimtalk kakao.Alimtalk, temp krt.ResultStr) {
	defer reswg.Done()

	for {
		if config.RL > 0 {
			config.RL--
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	resp, err := config.Client.R().
		SetHeaders(map[string]string{"Content-Type": "application/json"}).
		SetBody(alimtalk).
		Post(config.Conf.API_SERVER + "/v3/" + config.Conf.PROFILE_KEY + "/alimtalk/send")

	if err != nil {
		config.Stdlog.Println("alimtalk server request error : ", err, " / serial_number : ", alimtalk.Serial_number)
	} else {
		temp.Statuscode = resp.StatusCode()
		temp.BodyData = resp.Body()
	}
	c <- temp
}

func insAtErrResend(result map[string]string, rs []string, rv []interface{}, qm string) ([]string, []interface{}) {
	rs = append(rs, "(" + qm + ")")
	rv = append(rv, result["msgid"] + "_rs")
	rv = append(rv, result["userid"])
	rv = append(rv, result["ad_flag"])
	rv = append(rv, result["button1"])
	rv = append(rv, result["button2"])
	rv = append(rv, result["button3"])
	rv = append(rv, result["button4"])
	rv = append(rv, result["button5"])
	rv = append(rv, result["image_link"])
	rv = append(rv, result["image_url"])
	rv = append(rv, result["message_type"])
	rv = append(rv, result["msg"])
	rv = append(rv, result["msg_sms"])
	rv = append(rv, result["only_sms"])
	rv = append(rv, result["phn"])
	rv = append(rv, result["profile"])
	rv = append(rv, result["p_com"])
	rv = append(rv, result["p_invoice"])
	rv = append(rv, result["reg_dt"])
	rv = append(rv, result["remark1"])
	rv = append(rv, result["remark2"])
	rv = append(rv, 2)
	rv = append(rv, result["remark4"])
	rv = append(rv, result["remark5"])
	rv = append(rv, result["reserve_dt"])
	rv = append(rv, result["sms_kind"])
	rv = append(rv, result["sms_lms_tit"])
	rv = append(rv, result["sms_sender"])
	rv = append(rv, result["s_code"])
	rv = append(rv, result["tmpl_id"])
	rv = append(rv, result["wide"])
	rv = append(rv, nil) //send_group
	rv = append(rv, result["supplement"])

	if len(result["price"]) > 0 {
		price, _ := strconv.Atoi(result["price"])
		rv = append(rv, price)
	} else {
		rv = append(rv, nil)
	}

	rv = append(rv, result["currency_type"])
	rv = append(rv, result["title"])
	rv = append(rv, result["mms_image_id"])
	rv = append(rv, result["header"])
	rv = append(rv, result["attachments"])
	rv = append(rv, result["link"])
    rv = append(rv, result["msgid"])

	return rs, rv
}
