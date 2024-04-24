package kaosendrequest

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	kakao "mycs/src/kakaojson"
	config "mycs/src/kaoconfig"
	databasepool "mycs/src/kaodatabasepool"
	"mycs/src/kaocommon"

	"io/ioutil"
	"net"
	"net/http"
	s "strings"

	"strconv"
	"sync"
	"time"
	"context"
	"github.com/lib/pq"
)

var ftprocCnt int
var FisRunning bool
var isStoping bool

type resultStr struct {
	Statuscode int
	BodyData   []byte
	Result     map[string]string
}

func FriendtalkProc(ctx context.Context) {
	ftprocCnt = 1
	
	for {
		if ftprocCnt <=5 {
		
			select {
				case <- ctx.Done():
			
			    config.Stdlog.Println("Friendtalk process가 20초 후에 종료 됨.")
			    time.Sleep(20 * time.Second)
			    config.Stdlog.Println("Friendtalk process 종료 완료")
			    return
			default:
						
				var count sql.NullInt64
	
				cnterr := databasepool.DB.QueryRow("select count(1) as cnt from DHN_REQUEST  where send_group is null and (reserve_dt IS NULL OR to_timestamp(coalesce(reserve_dt,'00000000000000'), 'YYYYMMDDHH24MISS') <= NOW()) limit 1").Scan(&count)
	
				if cnterr != nil {
					//config.Stdlog.Println("DHN_REQUEST Table - select 오류 : " + cnterr.Error())
				} else {
	
					if count.Valid && count.Int64 > 0 {
						var startNow = time.Now()
						var group_no = fmt.Sprintf("%02d%02d%02d%09d", startNow.Hour(), startNow.Minute(), startNow.Second(), startNow.Nanosecond())
				
						updateRows, err := databasepool.DB.Exec("update DHN_REQUEST set send_group = '" + group_no + "' where id in (select id from dhn_request where send_group is null and (reserve_dt IS NULL OR to_timestamp(coalesce(reserve_dt,'00000000000000'), 'YYYYMMDDHH24MISS') <= NOW()) limit "+strconv.Itoa(config.Conf.SENDLIMIT)+")")
				
						if err != nil {
							config.Stdlog.Println("Request Table - send_group Update 오류")
						}
				
						rowcnt, _ := updateRows.RowsAffected()
				
						if rowcnt > 0 {
							config.Stdlog.Println("친구톡 발송 처리 시작 ( ", group_no, " ) : ", rowcnt, " 건 ")
							ftprocCnt ++
							go ftsendProcess(group_no)
						}
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
		errlog.Println("ftsendProcess 쿼리 에러 query : ", reqsql)
		errlog.Println("ftsendProcess 쿼리 에러 : ", err)
		errlog.Fatal(err)
	}

	columnTypes, err := reqrows.ColumnTypes()
	if err != nil {
		errlog.Println("ftsendProcess 컬럼 초기화 에러 group_no : ", group_no)
		errlog.Println("ftsendProcess 컬럼 초기화 에러 : ", err)
		errlog.Fatal(err)
	}
	count := len(columnTypes)
	initScanArgs := kaocommon.InitDatabaseColumn(columnTypes, count)

	var procCount int
	procCount = 0
	var startNow = time.Now()
	var serial_number = fmt.Sprintf("%04d%02d%02d-", startNow.Year(), startNow.Month(), startNow.Day())

	resultChan := make(chan resultStr, config.Conf.SENDLIMIT)
	var reswg sync.WaitGroup

	for reqrows.Next() {
		scanArgs := initScanArgs

		err := reqrows.Scan(scanArgs...)
		if err != nil {
			errlog.Println("ftsendProcess 컬럼 스캔 에러 group_no : ", group_no)
			errlog.Println("ftsendProcess 컬럼 스캔 에러 : ", err)
			errlog.Fatal(err)
		}

		var friendtalk kakao.Friendtalk
		var attache kakao.Attachment
		var tcarousel kakao.TCarousel
		var carousel kakao.FCarousel
		var button []kakao.Button
		var image kakao.Image
		var coupon kakao.AttCoupon
		var itemList kakao.AttItem
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

			case "carousel":
				if z, ok := (scanArgs[i]).(*sql.NullString); ok {
				    
					json.Unmarshal([]byte(z.String), &tcarousel)
					carousel.Tail = tcarousel.Tail
					  
					for ci, _ := range tcarousel.List {
						var catt kakao.CarouselAttachment
						var tcarlist kakao.CarouselList
						
						json.Unmarshal([]byte(tcarousel.List[ci].Attachment), &catt)
						
						tcarlist.Header = tcarousel.List[ci].Header
						tcarlist.Message = tcarousel.List[ci].Message
						tcarlist.Attachment = catt
						carousel.List = append(carousel.List, tcarlist)
					}
					//fmt.Println(len(carousel.List))
					if len(carousel.List) > 0 {
						//fmt.Println(carousel)
						friendtalk.Carousel = &carousel
					}  
				}

			case "image_url":
				if z, ok := (scanArgs[i]).(*sql.NullString); ok {
					image.Img_url = z.String
				}

			case "image_link":
				if z, ok := (scanArgs[i]).(*sql.NullString); ok {
					image.Img_link = z.String
				}

			case "att_items":
				if z, ok := (scanArgs[i]).(*sql.NullString); ok {
					error := json.Unmarshal([]byte(z.String), &itemList)
					if error == nil {
						attache.Item    = &itemList
					}
				}

			case "att_coupon":
				if z, ok := (scanArgs[i]).(*sql.NullString); ok {
					error := json.Unmarshal([]byte(z.String), &coupon)
					if error == nil {
						attache.Coupon    = &coupon
					}
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
			}

			if z, ok := (scanArgs[i]).(*sql.NullString); ok {
				config.Stdlog.Println(z.String)
				result[s.ToLower(v.Name())] = z.String
			}

			if z, ok := (scanArgs[i]).(*sql.NullInt32); ok {
				config.Stdlog.Println(string(z.Int32))
				result[s.ToLower(v.Name())] = string(z.Int32)
			}

			if z, ok := (scanArgs[i]).(*sql.NullInt64); ok {
				config.Stdlog.Println(string(z.Int64))
				result[s.ToLower(v.Name())] = string(z.Int64)
			}

		}

		if len(result["image_url"]) > 0 && s.EqualFold(result["message_type"], "FT") {
			friendtalk.Message_type = "FI"
			if s.EqualFold(result["wide"], "Y") {
				friendtalk.Message_type = "FW"
			}
		}

		attache.Buttons = button
		if len(image.Img_url) > 0 {
			attache.Ftimage = &image
		}
		friendtalk.Attachment = attache
		
		if s.EqualFold(conf.DEBUG,"Y") {
		  	jsonstr, _ := json.Marshal(friendtalk)
			stdlog.Println(string(jsonstr))
		}

		var temp resultStr
		temp.Result = result
		
		
		reswg.Add(1)
		go sendKakao(&reswg, resultChan, friendtalk, temp)

	}
	reswg.Wait()
	chanCnt := len(resultChan)

	ftValues := []kaocommon.FtResColumn{}

	for i := 0; i < chanCnt; i++ {

		resChan := <-resultChan
		result := resChan.Result

		if resChan.Statuscode == 200 {

			var kakaoResp kakao.KakaoResponse
			json.Unmarshal(resChan.BodyData, &kakaoResp)

			var resdt = time.Now()
			var resdtstr = fmt.Sprintf("%4d-%02d-%02d %02d:%02d:%02d", resdt.Year(), resdt.Month(), resdt.Day(), resdt.Hour(), resdt.Minute(), resdt.Second())

			ftValue := kaocommon.FtResColumn{}

			ftValue.Msgid = result["msgid"]
			ftValue.Userid = result["userid"]
			ftValue.Ad_flag = result["ad_flag"]
			ftValue.Button1 = result["button1"]
			ftValue.Button2 = result["button2"]
			ftValue.Button3 = result["button3"]
			ftValue.Button4 = result["button4"]
			ftValue.Button5 = result["button5"]
			ftValue.Code = kakaoResp.Code // 결과 code
			ftValue.Image_link = result["image_link"]
			ftValue.Image_url = result["image_url"]
			ftValue.Kind = nil // kind
			ftValue.Message = kakaoResp.Message // 결과 Message
			ftValue.Message_type = result["message_type"]
			ftValue.Msg = result["msg"]
			ftValue.Msg_sms = result["msg_sms"]
			ftValue.Only_sms = result["only_sms"]
			ftValue.P_com = result["p_com"]
			ftValue.P_invoice = result["p_invoice"]
			ftValue.Phn = result["phn"]
			ftValue.Profile = result["profile"]
			ftValue.Reg_dt = result["reg_dt"]
			ftValue.Remark1 = result["remark1"]
			ftValue.Remark2 = result["remark2"]
			ftValue.Remark3 = result["remark3"]
			ftValue.Remark4 = result["remark4"]
			ftValue.Remark5 = result["remark5"]
			ftValue.Res_dt = resdtstr // res_dt
			ftValue.Reserve_dt = result["reserve_dt"]
			if s.EqualFold(kakaoResp.Code,"0000") {
				ftValue.Result = "Y"
			} else if len(result["sms_kind"])>=1 && s.EqualFold(config.Conf.PHONE_MSG_FLAG, "YES") {
				ftValue.Result = "P" // sms_kind 가 SMS / LMS / MMS 이면 문자 발송 시도
			} else {
				ftValue.Result = "Y"
			}
			ftValue.S_code = kakaoResp.Code
			ftValue.Sms_kind = result["sms_kind"]
			ftValue.Sms_lms_tit = result["sms_lms_tit"]
			ftValue.Sms_sender = result["sms_sender"]
			ftValue.Sync = "N"
			ftValue.Tmpl_id = result["tmpl_id"]
			ftValue.Wide = result["wide"]
			if s.EqualFold(kakaoResp.Code,"0000") {
				ftValue.Send_group = nil
			} else if len(result["sms_kind"])>=1 && s.EqualFold(config.Conf.PHONE_MSG_FLAG, "YES") {
				ftValue.Send_group = nil
			} else {
				ftValue.Send_group = nil
			}
			ftValue.Supplement = result["supplement"]
			ftValue.Price = result["price"]
			ftValue.Currency_type = result["currency_type"]
			ftValue.Header = result["header"]
			ftValue.Carousel = result["carousel"]

			if len(ftValues) >= 500 {
				insertFtResData(ftValues)
				ftValues = []kaocommon.FtResColumn{}
			}

		} else {
			stdlog.Println("친구톡 서버 처리 오류 : ( ", string(resChan.BodyData), " )", result["msgid"])
			db.Exec("update DHN_REQUEST set send_group = null where msgid = '" + result["msgid"] + "'")
		}

		procCount++
	}

	if len(ftValues) > 0 {
		insertFtResData(ftValues)
	}

	db.Exec("delete from DHN_REQUEST where send_group = '" + group_no + "'")
	stdlog.Println("친구톡 발송 처리 완료 ( ", group_no, " ) : ", procCount, " 건  ( Proc Cnt :", ftprocCnt, ")" )
	
	ftprocCnt--

}

func insertFtResData(ftValues []kaocommon.FtResColumn) {
	tx, err := databasepool.DB.Begin()
	if err != nil {
		config.Stdlog.Println("friendtalk.go / insertFtResData / dhn_result / 트랜젝션 초기화 실패 ", err)
	}
	defer tx.Rollback()
	ftStmt, err := tx.Prepare(pq.CopyIn("dhn_result", kaocommon.GetReqColumnPq(kaocommon.FtResColumn{})...))
	if err != nil {
		config.Stdlog.Println("friendtalk.go / insertFtResData / dhn_result / ftStmt 초기화 실패 ", err)
		return
	}
	for _, data := range ftValues {
		_, err := ftStmt.Exec(data.Msgid,data.Userid,data.Ad_flag,data.Button1,data.Button2,data.Button3,data.Button4,data.Button5,data.Code,data.Image_link,data.Image_url,data.Kind,data.Message,data.Message_type,data.Msg,data.Msg_sms,data.Only_sms,data.P_com,data.P_invoice,data.Phn,data.Profile,data.Reg_dt,data.Remark1,data.Remark2,data.Remark3,data.Remark4,data.Remark5,data.Res_dt,data.Reserve_dt,data.Result,data.S_code,data.Sms_kind,data.Sms_lms_tit,data.Sms_sender,data.Sync,data.Tmpl_id,data.Wide,data.Send_group,data.Supplement,data.Price,data.Currency_type,data.Header,data.Carousel)
		if err != nil {
			config.Stdlog.Println("friendtalk.go / insertFtResData / dhn_result / ftStmt personal Exec ", err)
		}
	}
	
	_, err = ftStmt.Exec()
	if err != nil {
		atStmt.Close()
		config.Stdlog.Println("friendtalk.go / insertFtResData / dhn_result / ftStmt Exec ", err)
	}
	ftStmt.Close()
	err = tx.Commit()
	if err != nil {
		config.Stdlog.Println("friendtalk.go / insertFtResData / dhn_result / ftStmt commit ", err)
	}
}

func sendKakao(reswg *sync.WaitGroup, c chan<- resultStr, friendtalk kakao.Friendtalk, temp resultStr) {
	defer reswg.Done()

	jsonData, _ := json.Marshal(friendtalk)
	req, err := http.NewRequest("POST", config.Conf.API_SERVER + "v3/" + config.Conf.PROFILE_KEY + "/friendtalk/send", bytes.NewBuffer(jsonData))
	if err != nil {
		config.Stdlog.Println("친구톡 발송 에러 request 만들기 실패 ", err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := config.GoClient.Do(req)
	if err != nil {
		// 에러가 발생한 경우 처리
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			// 타임아웃 오류 처리
			config.Stdlog.Println("친구톡 발송 타임아웃 Serial_number : ", alimtalk.Serial_number, " / error : ", err.Error())
		} else {
			// 기타 오류 처리
			config.Stdlog.Println("친구톡 발송 실패 Serial_number : ", alimtalk.Serial_number, " / error : ", err.Error())
		}
		return
	} else {
		bodyData, _ := ioutil.ReadAll(resp.Body)
		temp.Statuscode = resp.StatusCode
		temp.BodyData = bodyData
	}

	resp.Body.Close()

	c <- temp

}
