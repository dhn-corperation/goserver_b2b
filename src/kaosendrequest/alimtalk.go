package kaosendrequest

import (
	//"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	kakao "mycs/src/kakaojson"
	config "mycs/src/kaoconfig"
	databasepool "mycs/src/kaodatabasepool"
	"mycs/src/kaocommon"

	//"io/ioutil"
	//	"net"
	//"net/http"
	"strconv"
	s "strings"
	"sync"
	"time"
	"context"
	"github.com/lib/pq"
)

var atprocCnt int

func AlimtalkProc( user_id string, ctx context.Context ) {
	atprocCnt = 1
	config.Stdlog.Println(user_id, " - 알림톡 프로세스 시작 됨 ") 
	for {
		if atprocCnt <=5 {
			
			select {
				case <- ctx.Done():
			    config.Stdlog.Println(user_id, " - Alimtalk process가 20초 후에 종료 됨.")
			    time.Sleep(20 * time.Second)
			    config.Stdlog.Println(user_id, " - Alimtalk process 종료 완료")
			    return
			default:
				var count sql.NullInt64
				cnterr := databasepool.DB.QueryRowContext(ctx, "SELECT count(1) AS cnt FROM DHN_REQUEST_AT WHERE send_group IS NULL AND (reserve_dt IS NULL OR to_timestamp(coalesce(reserve_dt,'00000000000000'), 'YYYYMMDDHH24MISS') <= NOW()) AND userid='"+user_id+"'").Scan(&count)
				if cnterr != nil {
					config.Stdlog.Println("DHN_REQUEST Table - select 오류 " + cnterr.Error())
				} else {
					if count.Valid && count.Int64 > 0 {		
						var startNow = time.Now()
						var group_no = fmt.Sprintf("%02d%02d%02d%09d", startNow.Hour(), startNow.Minute(), startNow.Second(), startNow.Nanosecond())
						
						updateRows, err := databasepool.DB.ExecContext(ctx, "update DHN_REQUEST_AT set send_group = '"+group_no+"' where send_group is null and (reserve_dt IS NULL OR to_timestamp(coalesce(reserve_dt,'00000000000000'), 'YYYYMMDDHH24MISS') <= NOW()) and userid = '"+user_id+"'  limit "+strconv.Itoa(config.Conf.SENDLIMIT))
				
						if err != nil {
							config.Stdlog.Println(user_id,"알림톡 send_group Update 오류 : ", err)
						}
				
						rowcnt, _ := updateRows.RowsAffected()
				
						if rowcnt > 0 {
							config.Stdlog.Println(user_id, "알림톡 발송 처리 시작 ( ", group_no, " ) : ", rowcnt, " 건 ")
							atprocCnt++
							go atsendProcess(group_no, user_id)
				
						}
					}
				}
			}
		}
	}
}

func atsendProcess(group_no string, user_id string) {

	var db = databasepool.DB
	var conf = config.Conf
	var stdlog = config.Stdlog
	var errlog = config.Stdlog

	reqsql := "select * from DHN_REQUEST_AT where send_group = '" + group_no + "' and userid = '" + user_id + "'"

	reqrows, err := db.Query(reqsql)
	if err != nil {
		errlog.Println("atsendProcess 쿼리 에러 query : ", reqsql)
		errlog.Println("atsendProcess 쿼리 에러 : ", err)
		errlog.Fatal(err)
	}

	columnTypes, err := reqrows.ColumnTypes()
	if err != nil {
		errlog.Println("atsendProcess 컬럼 초기화 에러 group_no : ", group_no, " / userid  : ", user_id)
		errlog.Println("atsendProcess 컬럼 초기화 에러 : ", err)
		errlog.Fatal(err)
	}
	count := len(columnTypes)
	initScanArgs := kaocommon.InitDatabaseColumn(columnTypes, count)

	var procCount int
	procCount = 0
	var startNow = time.Now()
	var serial_number = fmt.Sprintf("%04d%02d%02d-", startNow.Year(), startNow.Month(), startNow.Day())

	resultChan := make(chan resultStr, config.Conf.SENDLIMIT) // resultStr 은 friendtalk에 정의 됨
	var reswg sync.WaitGroup

	for reqrows.Next() {
		scanArgs := initScanArgs

		err := reqrows.Scan(scanArgs...)
		if err != nil {
			errlog.Println("atsendProcess 컬럼 스캔 에러 group_no : ", group_no, " / userid  : ", user_id)
			errlog.Println("atsendProcess 컬럼 스캔 에러 : ", err)
			errlog.Fatal(err)
		}

		var alimtalk kakao.Alimtalk
		var attache kakao.AttachmentB
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
					alimtalk.Phone_number = z.String
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

		if s.EqualFold(result["message_type"], "at") {
			alimtalk.Response_method = "realtime"
		} else if s.EqualFold(result["message_type"], "ai") {
			
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

		var temp resultStr
		temp.Result = result
		reswg.Add(1)
		go sendKakaoAlimtalk(&reswg, resultChan, alimtalk, temp)
	}
	reswg.Wait()
	chanCnt := len(resultChan)

	atValues := []kaocommon.AtResColumn{}

	for i := 0; i < chanCnt; i++ {

		resChan := <-resultChan
		result := resChan.Result
		if resChan.Statuscode == 200 {

			// var kakaoResp kakao.KakaoResponse
			// json.Unmarshal(resChan.BodyData, &kakaoResp)

			var resdt = time.Now()
			var resdtstr = fmt.Sprintf("%4d-%02d-%02d %02d:%02d:%02d", resdt.Year(), resdt.Month(), resdt.Day(), resdt.Hour(), resdt.Minute(), resdt.Second())
			
			// var resCode = kakaoResp.Code
			// var resMessage = kakaoResp.Message

			var resCode = "0000"
			var resMessage = ""
			
			if s.EqualFold(resCode, "3005") {
				resCode = "0000"
				resMessage = ""
			} 

			atValue := kaocommon.AtResColumn{}

			atValue.Msgid = result["msgid"]
			atValue.Userid = result["userid"]
			atValue.Ad_flag = result["ad_flag"]
			atValue.Button1 = result["button1"]
			atValue.Button2 = result["button2"]
			atValue.Button3 = result["button3"]
			atValue.Button4 = result["button4"]
			atValue.Button5 = result["button5"]
			atValue.Code = resCode
			atValue.Image_link = result["image_link"]
			atValue.Image_url = result["image_url"]
			atValue.Kind = nil
			atValue.Message = resMessage
			atValue.Message_type = result["message_type"]
			atValue.Msg = result["msg"]
			atValue.Msg_sms = result["msg_sms"]
			atValue.Only_sms = result["only_sms"]
			atValue.P_com = result["p_com"]
			atValue.P_invoice = result["p_invoice"]
			atValue.Profile = result["profile"]
			atValue.Reg_dt = result["reg_dt"]
			atValue.Remark1 = result["remark1"]
			atValue.Remark2 = result["remark2"]
			atValue.Remark3 = result["remark3"]
			atValue.Remark4 = result["remark4"]
			atValue.Remark5 = result["remark5"]
			atValue.Res_dt = resdtstr
			atValue.Reserve_dt = result["reserve_dt"]
			//result 컬럼 처리
			if s.EqualFold(result["message_type"], "at") || !s.EqualFold(resCode, "0000") || (s.EqualFold(result["message_type"], "ai") && s.EqualFold(conf.RESPONSE_METHOD, "push")) {

				if s.EqualFold(resCode,"0000") {
					atValue.Result = "Y" // 
				// 1차 카카오 발송 실패 후 2차 발송을 바로 하기 위해서는 이 조건을 맞춰야함
				} else if len(result["sms_kind"])>=1 && s.EqualFold(config.Conf.PHONE_MSG_FLAG, "YES") {
					atValue.Result = "P" // sms_kind 가 SMS / LMS / MMS 이면 문자 발송 시도
				} else {
					atValue.Result = "Y" // 
				} 
				
			} else if s.EqualFold(result["message_type"], "ai") {
				atValue.Result = "N" // result
			}
			atValue.S_code = resCode
			atValue.Sms_kind = result["sms_kind"]
			atValue.Sms_lms_tit = result["sms_lms_tit"]
			atValue.Sms_sender = result["sms_sender"]
			atValue.Sync = "N"
			atValue.Tmpl_id = result["tmpl_id"]
			atValue.Wide = result["wide"]
			//send_group 컬럼 처리
			if s.EqualFold(result["message_type"], "at") || !s.EqualFold(resCode, "0000") || (s.EqualFold(result["message_type"], "ai") && s.EqualFold(conf.RESPONSE_METHOD, "push")) {

				if s.EqualFold(resCode,"0000") {
					atValue.Send_group = nil //send_group
				} else if len(result["sms_kind"])>=1 && s.EqualFold(config.Conf.PHONE_MSG_FLAG, "YES") {
					atValue.Send_group = nil //send_group
				} else {
					atValue.Send_group = nil //send_group
				} 
				
			} else if s.EqualFold(result["message_type"], "ai") {
				atValue.Send_group = nil //send_group
			}
			atValue.Supplement = result["supplement"]
			atValue.Price = result["price"]
			atValue.Currency_type = result["currency_type"]
			atValue.Title = result["title"]

			atValues = append(atValues, atValue)

			//Center에서도 사용하고 있는 함수이므로 공용 라이브러리 생성이 필요함
			if len(atValues) >= 500 {
				tx, err := databasepool.DB.Begin()
				if err != nil {
					errlog.Println(err)
				}
				defer tx.Rollback()
				atStmt, err := tx.Prepare(pq.CopyIn("dhn_result", kaocommon.GetReqColumnPq(kaocommon.AtResColumn{})...))
				if err != nil {
					errlog.Println("atStmt 초기화 실패 ", err)
					return
				}
				for _, data := range atValues {
					_, err := atStmt.Exec(data.Msgid,data.Userid,data.Ad_flag,data.Button1,data.Button2,data.Button3,data.Button4,data.Button5,data.Code,data.Image_link,data.Image_url,data.Kind,data.Message,data.Message_type,data.Msg,data.Msg_sms,data.Only_sms,data.P_com,data.P_invoice,data.Phn,data.Profile,data.Reg_dt,data.Remark1,data.Remark2,data.Remark3,data.Remark4,data.Remark5,data.Res_dt,data.Reserve_dt,data.Result,data.S_code,data.Sms_kind,data.Sms_lms_tit,data.Sms_sender,data.Sync,data.Tmpl_id,data.Wide,data.Send_group,data.Supplement,data.Price,data.Currency_type,data.Title)
					if err != nil {
						errlog.Println(err)
					}
				}
				atValues = []kaocommon.AtResColumn{}
				_, err = atStmt.Exec()
				if err != nil {
					atStmt.Close()
					errlog.Println(err)
				}
				atStmt.Close()
				err = tx.Commit()
				if err != nil {
					errlog.Println(err)
				}
			}
		} else {
			// stdlog.Println(user_id, "알림톡 서버 처리 오류 !! ( ", string(resChan.BodyData), " )", result["msgid"])
			db.Exec("update DHN_REQUEST_AT set send_group = null where msgid = '" + result["msgid"] + "'")
		}

		procCount++
	}

	//Center에서도 사용하고 있는 함수이므로 공용 라이브러리 생성이 필요함
	if len(atValues) > 0 {
		tx, err := databasepool.DB.Begin()
		if err != nil {
			errlog.Println(err)
		}
		defer tx.Rollback()
		atStmt, err := tx.Prepare(pq.CopyIn("dhn_result", kaocommon.GetReqColumnPq(kaocommon.AtResColumn{})...))
		if err != nil {
			errlog.Println("atStmt 초기화 실패 ", err)
			return
		}
		for _, data := range atValues {
			_, err := atStmt.Exec(data.Msgid,data.Userid,data.Ad_flag,data.Button1,data.Button2,data.Button3,data.Button4,data.Button5,data.Code,data.Image_link,data.Image_url,data.Kind,data.Message,data.Message_type,data.Msg,data.Msg_sms,data.Only_sms,data.P_com,data.P_invoice,data.Phn,data.Profile,data.Reg_dt,data.Remark1,data.Remark2,data.Remark3,data.Remark4,data.Remark5,data.Res_dt,data.Reserve_dt,data.Result,data.S_code,data.Sms_kind,data.Sms_lms_tit,data.Sms_sender,data.Sync,data.Tmpl_id,data.Wide,data.Send_group,data.Supplement,data.Price,data.Currency_type,data.Title)
			if err != nil {
				errlog.Println(err)
			}
		}
		atValues = []kaocommon.AtResColumn{}
		_, err = atStmt.Exec()
		if err != nil {
			atStmt.Close()
			errlog.Println(err)
		}
		atStmt.Close()
		err = tx.Commit()
		if err != nil {
			errlog.Println(err)
		}
	}

	//알림톡 발송 후 DHN_REQUEST_AT 테이블의 데이터는 제거한다.
	db.Exec("delete from DHN_REQUEST_AT where send_group = '" + group_no + "'")

	stdlog.Println(user_id, "알림톡 발송 처리 완료 ( ", group_no, " ) : ", procCount, " 건 ( Proc Cnt :", atprocCnt, ")")
	var resdt2 = time.Now()
	stdlog.Println(fmt.Sprintf("%4d-%02d-%02d %02d:%02d:%02d", resdt2.Year(), resdt2.Month(), resdt2.Day(), resdt2.Hour(), resdt2.Minute(), resdt2.Second()))
	
	atprocCnt--
}

//카카오 서버에 발송을 요청한다.
func sendKakaoAlimtalk(reswg *sync.WaitGroup, c chan<- resultStr, alimtalk kakao.Alimtalk, temp resultStr) {
	defer reswg.Done()

	// resp, err := config.Client.R().
	// 	SetHeaders(map[string]string{"Content-Type": "application/json"}).
	// 	SetBody(alimtalk).
	// 	Post(config.Conf.API_SERVER + "/v3/" + config.Conf.PROFILE_KEY + "/alimtalk/send")

	// if err != nil {
	// 	config.Stdlog.Println("알림톡 메시지 서버 호출 오류 : ", err)
	// } else {
	// 	temp.Statuscode = resp.StatusCode()
	// 	temp.BodyData = resp.Body()
	// }
	temp.Statuscode = 200
	c <- temp

}
