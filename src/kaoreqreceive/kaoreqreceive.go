package kaoreqreceive

import (
	//"encoding/json"
	"fmt"
	config "mycs/src/kaoconfig"
	databasepool "mycs/src/kaodatabasepool"
	"mycs/src/kaoreqtable"
	"mycs/src/kaocommon"
	"strconv"
	s "strings"
	"time"
	"database/sql"

	"github.com/gin-gonic/gin"
	// _ "github.com/go-sql-driver/mysql"
	pq "github.com/lib/pq"
)

//언젠가는 다른곳으로 위치를 옮겨야 함
var SecretKey = "9b4dabe9d4fed126a58f8639846143c7"

func ReqReceive(c *gin.Context) {
	execFlag := false
	ctx := c.Request.Context()
	errlog := config.Stdlog

	userid := c.Request.Header.Get("userid")
	userip := c.ClientIP()
	isValidation := false

	// 허가된 userid 인지 테이블에서 확인
	sqlstr := `
		select 
			count(1) as cnt 
		from
			DHN_CLIENT_LIST
		where
			user_id = $1
			and ip = $2
			and use_flag = 'Y'`
	var cnt sql.NullInt64
	err := databasepool.DB.QueryRowContext(ctx, sqlstr, userid, userip).Scan(&cnt)
	if err != nil { errlog.Println("DHN_CLIENT_LIST 쿼리 에러 ", err) }

	if cnt.Valid && cnt.Int64 > 0 { 
		isValidation = true 
	} else {
		errlog.Println("허용되지 않은 사용자 및 아이피에서 발송 요청!! (userid : ", userid, "/ ip : ", userip, ")")
	}

	var startNow = time.Now()
	var startTime = fmt.Sprintf("%02d:%02d:%02d", startNow.Hour(), startNow.Minute(), startNow.Second())

	if isValidation {

		var msg []kaoreqtable.Reqtable
		//전달온 데이터 kaoreqtable.Reqtable에 맵핑
		err1 := c.ShouldBindJSON(&msg)

		if err1 != nil { errlog.Println(err1) }

		errlog.Println("발송 메세지 수신 시작 ( ", userid, ") : ", len(msg), startTime)

		tx, err := databasepool.DB.Begin()
		if err != nil {
			errlog.Println(err)
		}
		defer tx.Rollback()

		ftStmt, err := tx.Prepare(pq.CopyIn("dhn_request", kaocommon.GetReqColumnPq(kaocommon.FtReqColumn{})...))
		if err != nil {
			errlog.Println("ftStmt 초기화 실패 ", err)
		}
		defer ftStmt.Close()

		atStmt, err := tx.Prepare(pq.CopyIn("dhn_request_at", kaocommon.GetReqColumnPq(kaocommon.AtReqColumn{})...))
		if err != nil {
			errlog.Println("atStmt 초기화 실패 ", err)
		}
		defer atStmt.Close()

		msgStmt, err := tx.Prepare(pq.CopyIn("dhn_result", kaocommon.GetReqColumnPq(kaocommon.MsgReqColumn{})...))
		if err != nil {
			errlog.Println("msgStmt 초기화 실패 ", err)
		}
		defer atStmt.Close()

		msgTempStmt, _ := tx.Prepare(pq.CopyIn("dhn_result_temp", kaocommon.GetReqColumnPq(kaocommon.MsgReqColumn{})...))
		if err != nil {
			errlog.Println("msgTempStmt 초기화 실패 ", err)
		}
		defer msgTempStmt.Close()

		ftValues := []kaocommon.FtReqColumn{}
		atValues := []kaocommon.AtReqColumn{}
		msgValues := []kaocommon.MsgReqColumn{}

		//맵핑한 데이터 row 처리
		for i, _ := range msg {
			var nonce string
			if len(msg[i].Crypto) > 0 {
				nonce = s.Split(msg[i].Crypto, ",")[0]
			}
			//친구톡 insert values 만들기
			if s.HasPrefix(s.ToUpper(msg[i].Messagetype), "F") {
				ftValue := kaocommon.FtReqColumn{}

				ftValue.Msgid = msg[i].Msgid
				ftValue.Userid = userid
				ftValue.Ad_flag = msg[i].Adflag
				ftValue.Button1 = msg[i].Button1
				ftValue.Button2 = msg[i].Button2
				ftValue.Button3 = msg[i].Button3
				ftValue.Button4 = msg[i].Button4
				ftValue.Button5 = msg[i].Button5
				ftValue.Image_link = msg[i].Imagelink
				ftValue.Image_url = msg[i].Imageurl
				ftValue.Message_type = msg[i].Messagetype
				if s.Contains(msg[i].Crypto, "MSG") {
					ftValue.Msg = kaocommon.AES256GSMDecrypt([]byte(SecretKey), msg[i].Msg, nonce)
				} else {
					ftValue.Msg = msg[i].Msg
				}
				if s.Contains(msg[i].Crypto, "Msgsms") && len(msg[i].Msgsms) > 0 {
					ftValue.Msg_sms = kaocommon.AES256GSMDecrypt([]byte(SecretKey), msg[i].Msgsms, nonce)
				} else {
					ftValue.Msg_sms = msg[i].Msgsms
				}
				ftValue.Only_sms = msg[i].Onlysms
				ftValue.Phn = kaocommon.AES256GSMDecrypt([]byte(SecretKey), msg[i].Phn, nonce)
				if s.Contains(msg[i].Crypto, "Profile") && len(msg[i].Profile) > 0 {
					ftValue.Profile = kaocommon.AES256GSMDecrypt([]byte(SecretKey), msg[i].Profile, nonce)
				} else {
					ftValue.Profile = msg[i].Profile
				}
				ftValue.P_com = msg[i].Pcom
				ftValue.P_invoice = msg[i].Pinvoice
				ftValue.Reg_dt = msg[i].Regdt
				ftValue.Remark1 = msg[i].Remark1
				ftValue.Remark2 = msg[i].Remark2
				ftValue.Remark3 = msg[i].Remark3
				ftValue.Remark4 = msg[i].Remark4
				ftValue.Remark5 = msg[i].Remark5
				ftValue.Reserve_dt = msg[i].Reservedt
				ftValue.Sms_kind = msg[i].Smskind
				if s.Contains(msg[i].Crypto, "Smslmstit") && len(msg[i].Smslmstit) > 0 {
					ftValue.Sms_lms_tit = kaocommon.AES256GSMDecrypt([]byte(SecretKey), msg[i].Smslmstit, nonce)
				} else {
					ftValue.Sms_lms_tit = msg[i].Smslmstit
				}
				if s.Contains(msg[i].Crypto, "Smssender") && len(msg[i].Smssender) > 0 {
					ftValue.Sms_sender = kaocommon.AES256GSMDecrypt([]byte(SecretKey), msg[i].Smssender, nonce)
				} else {
					ftValue.Sms_sender = msg[i].Smssender
				}
				ftValue.S_code = msg[i].Scode
				ftValue.Tmpl_id = msg[i].Tmplid
				ftValue.Wide = msg[i].Wide
				ftValue.Send_group = "\\n"
				ftValue.Supplement = msg[i].Supplement
				if len(msg[i].Price) > 0 {
					price, _ := strconv.Atoi(msg[i].Price)
					ftValue.Price = sql.NullInt64{Int64: int64(price), Valid: true}
				} else {
					ftValue.Price = sql.NullInt64{Valid: false}
				}

				ftValue.Currency_type = msg[i].Currencytype
				ftValue.Title = msg[i].Title
				ftValue.Header = msg[i].Header
				ftValue.Carousel = msg[i].Carousel
				ftValue.Att_items = msg[i].Att_items
				ftValue.Att_coupon = msg[i].Att_coupon
				ftValue.Attachments = msg[i].Attachments

				ftValues = append(ftValues, ftValue)

			//문자 insert values 만들기
			} else if s.EqualFold(msg[i].Messagetype, "PH") {
				var resdt = time.Now()
				var resdtstr = fmt.Sprintf("%4d-%02d-%02d %02d:%02d:%02d", resdt.Year(), resdt.Month(), resdt.Day(), resdt.Hour(), resdt.Minute(), resdt.Second())

				msgValue := kaocommon.MsgReqColumn{}

				msgValue.Msgid = msg[i].Msgid
				msgValue.Userid = userid
				msgValue.Ad_flag = msg[i].Adflag
				msgValue.Button1 = msg[i].Button1
				msgValue.Button2 = msg[i].Button2
				msgValue.Button3 = msg[i].Button3
				msgValue.Button4 = msg[i].Button4
				msgValue.Button5 = msg[i].Button5
				msgValue.Code = "9999"
				msgValue.Image_link = msg[i].Imagelink
				msgValue.Image_url = msg[i].Imageurl
				msgValue.Kind = "\\n"
				msgValue.Message = ""
				msgValue.Message_type = msg[i].Messagetype
				if s.Contains(msg[i].Crypto, "MSG") {
					msgValue.Msg = kaocommon.AES256GSMDecrypt([]byte(SecretKey), msg[i].Msg, nonce)
				} else {
					msgValue.Msg = msg[i].Msg
				}
				if s.Contains(msg[i].Crypto, "Msgsms") && len(msg[i].Msgsms) > 0 {
					msgValue.Msg_sms = kaocommon.AES256GSMDecrypt([]byte(SecretKey), msg[i].Msgsms, nonce)
				} else {
					msgValue.Msg_sms = msg[i].Msgsms
				}
				msgValue.Only_sms = msg[i].Onlysms
				msgValue.Phn = kaocommon.AES256GSMDecrypt([]byte(SecretKey), msg[i].Phn, nonce)
				if s.Contains(msg[i].Crypto, "Profile") && len(msg[i].Profile) > 0 {
					msgValue.Profile = kaocommon.AES256GSMDecrypt([]byte(SecretKey), msg[i].Profile, nonce)
				} else {
					msgValue.Profile = msg[i].Profile
				}
				msgValue.P_com = msg[i].Pcom
				msgValue.P_invoice = msg[i].Pinvoice
				msgValue.Reg_dt = msg[i].Regdt
				msgValue.Remark1 = msg[i].Remark1
				msgValue.Remark2 = msg[i].Remark2
				msgValue.Remark3 = msg[i].Remark3
				msgValue.Remark4 = msg[i].Remark4
				msgValue.Remark5 = msg[i].Remark5
				msgValue.Res_dt = resdtstr
				msgValue.Reserve_dt = msg[i].Reservedt
				msgValue.Result = "P"  // sms_kind 가 SMS / LMS / MMS 이면 문자 발송 시도
				msgValue.S_code = msg[i].Scode
				msgValue.Sms_kind = msg[i].Smskind
				if s.Contains(msg[i].Crypto, "Smslmstit") && len(msg[i].Smslmstit) > 0 {
					msgValue.Sms_lms_tit = kaocommon.AES256GSMDecrypt([]byte(SecretKey), msg[i].Smslmstit, nonce)
				} else {
					msgValue.Sms_lms_tit = msg[i].Smslmstit
				}
				if s.Contains(msg[i].Crypto, "Smssender") && len(msg[i].Smssender) > 0 {
					msgValue.Sms_sender = kaocommon.AES256GSMDecrypt([]byte(SecretKey), msg[i].Smssender, nonce)
				} else {
					msgValue.Sms_sender = msg[i].Smssender
				}
				msgValue.Sync = "N"
				msgValue.Tmpl_id = msg[i].Tmplid
				msgValue.Wide = msg[i].Wide
				msgValue.Send_group = "\\n"
				msgValue.Supplement = msg[i].Supplement
				msgValue.Price = sql.NullInt64{Valid: false}
				msgValue.Currency_type = "\\n"
				msgValue.Header = msg[i].Header
				msgValue.Carousel = msg[i].Carousel

				msgValues = append(msgValues, msgValue)

			//알림톡 insert values 만들기
			} else {
				atValue := kaocommon.AtReqColumn{}

				atValue.Msgid = msg[i].Msgid
				atValue.Userid = userid
				atValue.Ad_flag = msg[i].Adflag
				atValue.Button1 = msg[i].Button1
				atValue.Button2 = msg[i].Button2
				atValue.Button3 = msg[i].Button3
				atValue.Button4 = msg[i].Button4
				atValue.Button5 = msg[i].Button5
				atValue.Image_link = msg[i].Imagelink
				atValue.Image_url = msg[i].Imageurl
				atValue.Message_type = msg[i].Messagetype
				if s.Contains(msg[i].Crypto, "MSG") {
					atValue.Msg = kaocommon.AES256GSMDecrypt([]byte(SecretKey), msg[i].Msg, nonce)
				} else {
					atValue.Msg = msg[i].Msg
				}
				if s.Contains(msg[i].Crypto, "Msgsms") && len(msg[i].Msgsms) > 0 {
					atValue.Msg_sms = kaocommon.AES256GSMDecrypt([]byte(SecretKey), msg[i].Msgsms, nonce)
				} else {
					atValue.Msg_sms = msg[i].Msgsms
				}
				atValue.Only_sms = msg[i].Onlysms
				atValue.Phn = kaocommon.AES256GSMDecrypt([]byte(SecretKey), msg[i].Phn, nonce)
				if s.Contains(msg[i].Crypto, "Profile") && len(msg[i].Profile) > 0 {
					atValue.Profile = kaocommon.AES256GSMDecrypt([]byte(SecretKey), msg[i].Profile, nonce)
				} else {
					atValue.Profile = msg[i].Profile
				}
				atValue.P_com = msg[i].Pcom
				atValue.P_invoice = msg[i].Pinvoice
				atValue.Reg_dt = msg[i].Regdt
				atValue.Remark1 = msg[i].Remark1
				atValue.Remark2 = msg[i].Remark2
				atValue.Remark3 = msg[i].Remark3
				atValue.Remark4 = msg[i].Remark4
				atValue.Remark5 = msg[i].Remark5
				atValue.Reserve_dt = msg[i].Reservedt
				atValue.Sms_kind = msg[i].Smskind
				if s.Contains(msg[i].Crypto, "Smslmstit") && len(msg[i].Smslmstit) > 0 {
					atValue.Sms_lms_tit = kaocommon.AES256GSMDecrypt([]byte(SecretKey), msg[i].Smslmstit, nonce)
				} else {
					atValue.Sms_lms_tit = msg[i].Smslmstit
				}
				if s.Contains(msg[i].Crypto, "Smssender") && len(msg[i].Smssender) > 0 {
					atValue.Sms_sender = kaocommon.AES256GSMDecrypt([]byte(SecretKey), msg[i].Smssender, nonce)
				} else {
					atValue.Sms_sender = msg[i].Smssender
				}
				atValue.S_code = msg[i].Scode
				atValue.Tmpl_id = msg[i].Tmplid
				atValue.Wide = msg[i].Wide
				atValue.Send_group = "\\n"
				atValue.Supplement = msg[i].Supplement
				if len(msg[i].Price) > 0 {
					price, _ := strconv.Atoi(msg[i].Price)
					atValue.Price = sql.NullInt64{Int64: int64(price), Valid: true}
				} else {
					atValue.Price = sql.NullInt64{Valid: false}
				}

				atValue.Currency_type = msg[i].Currencytype
				atValue.Title = msg[i].Title

				atValues = append(atValues, atValue)
			}

			// 500건 단위로 처리한다(클라이언트에서 1000건씩 전송하더라도 지정한 단위의 건수로 insert한다.)
			saveCount := 500

			if len(ftValues) >= saveCount {
				for _, data := range ftValues {
					_, err := ftStmt.Exec(data.Msgid,data.Userid,data.Ad_flag,data.Button1,data.Button2,data.Button3,data.Button4,data.Button5,data.Image_link,data.Image_url,data.Message_type,data.Msg,data.Msg_sms,data.Only_sms,data.P_com,data.P_invoice,data.Phn,data.Profile,data.Reg_dt,data.Remark1,data.Remark2,data.Remark3,data.Remark4,data.Remark5,data.Reserve_dt,data.S_code,data.Sms_kind,data.Sms_lms_tit,data.Sms_sender,data.Tmpl_id,data.Wide,data.Send_group,data.Supplement,data.Price,data.Currency_type,data.Header,data.Carousel,data.Att_coupon,data.Attachments)
					if err != nil {
						errlog.Println(err)
					}
				}
				ftValues = []kaocommon.FtReqColumn{}
				_, err = ftStmt.Exec()
				if err != nil {
					errlog.Println(err)
				}
			}

			if len(atValues) >= saveCount {
				for _, data := range atValues {
					_, err := atStmt.Exec(data.Msgid,data.Userid,data.Ad_flag,data.Button1,data.Button2,data.Button3,data.Button4,data.Button5,data.Image_link,data.Image_url,data.Message_type,data.Msg,data.Msg_sms,data.Only_sms,data.P_com,data.P_invoice,data.Phn,data.Profile,data.Reg_dt,data.Remark1,data.Remark2,data.Remark3,data.Remark4,data.Remark5,data.Reserve_dt,data.S_code,data.Sms_kind,data.Sms_lms_tit,data.Sms_sender,data.Tmpl_id,data.Wide,data.Send_group,data.Supplement,data.Price,data.Currency_type)
					if err != nil {
						errlog.Println(err)
					}
				}
				atValues = []kaocommon.AtReqColumn{}
				execFlag = true
				_, err = atStmt.Exec()
				if err != nil {
					errlog.Println(err)
				}
			}

			if len(msgValues) >= saveCount {
				for _, data := range msgValues {
					_, err := msgStmt.Exec(data.Msgid,data.Userid,data.Ad_flag,data.Button1,data.Button2,data.Button3,data.Button4,data.Button5,data.Code,data.Image_link,data.Image_url,data.Kind,data.Message,data.Message_type,data.Msg,data.Msg_sms,data.Only_sms,data.Phn,data.Profile,data.P_com,data.P_invoice,data.Reg_dt,data.Remark1,data.Remark2,data.Remark3,data.Remark4,data.Remark5,data.Res_dt,data.Reserve_dt,data.Result,data.S_code,data.Sms_kind,data.Sms_lms_tit,data.Sms_sender,data.Sync,data.Tmpl_id,data.Wide,data.Send_group,data.Supplement,data.Price,data.Currency_type,data.Header,data.Carousel)
					if err != nil {
						errlog.Println(err)
					}
				}
				msgValues = []kaocommon.MsgReqColumn{}
				execFlag = true
				_, err = msgStmt.Exec()
				if err != nil {
					errlog.Println(err)
				}
			}
		}
		
		// 나머지 건수를 저장하기 위해 다시한번 정의
		if len(ftValues) > 0 {
			for _, data := range ftValues {
				_, err := ftStmt.Exec(data.Msgid,data.Userid,data.Ad_flag,data.Button1,data.Button2,data.Button3,data.Button4,data.Button5,data.Image_link,data.Image_url,data.Message_type,data.Msg,data.Msg_sms,data.Only_sms,data.P_com,data.P_invoice,data.Phn,data.Profile,data.Reg_dt,data.Remark1,data.Remark2,data.Remark3,data.Remark4,data.Remark5,data.Reserve_dt,data.S_code,data.Sms_kind,data.Sms_lms_tit,data.Sms_sender,data.Tmpl_id,data.Wide,data.Send_group,data.Supplement,data.Price,data.Currency_type,data.Header,data.Carousel,data.Att_coupon,data.Attachments)
				if err != nil {
					errlog.Println(err)
				}
			}
			ftValues = []kaocommon.FtReqColumn{}
			execFlag = true
			_, err = ftStmt.Exec()
			if err != nil {
				errlog.Println(err)
			}
		}

		if len(atValues) > 0 {
			for _, data := range atValues {
				_, err := atStmt.Exec(data.Msgid,data.Userid,data.Ad_flag,data.Button1,data.Button2,data.Button3,data.Button4,data.Button5,data.Image_link,data.Image_url,data.Message_type,data.Msg,data.Msg_sms,data.Only_sms,data.P_com,data.P_invoice,data.Phn,data.Profile,data.Reg_dt,data.Remark1,data.Remark2,data.Remark3,data.Remark4,data.Remark5,data.Reserve_dt,data.S_code,data.Sms_kind,data.Sms_lms_tit,data.Sms_sender,data.Tmpl_id,data.Wide,data.Send_group,data.Supplement,data.Price,data.Currency_type)
				if err != nil {
					errlog.Println(err)
				}
			}
			atValues = []kaocommon.AtReqColumn{}
			execFlag = true
			_, err = atStmt.Exec()
			if err != nil {
				errlog.Println(err)
			}
		}

		if len(msgValues) > 0 {
			for _, data := range msgValues {
				_, err := msgStmt.Exec(data.Msgid,data.Userid,data.Ad_flag,data.Button1,data.Button2,data.Button3,data.Button4,data.Button5,data.Code,data.Image_link,data.Image_url,data.Kind,data.Message,data.Message_type,data.Msg,data.Msg_sms,data.Only_sms,data.Phn,data.Profile,data.P_com,data.P_invoice,data.Reg_dt,data.Remark1,data.Remark2,data.Remark3,data.Remark4,data.Remark5,data.Res_dt,data.Reserve_dt,data.Result,data.S_code,data.Sms_kind,data.Sms_lms_tit,data.Sms_sender,data.Sync,data.Tmpl_id,data.Wide,data.Send_group,data.Supplement,data.Price,data.Currency_type,data.Header,data.Carousel)
				if err != nil {
					errlog.Println(err)
				}
			}
			msgValues = []kaocommon.MsgReqColumn{}
			execFlag = true
			_, err = msgStmt.Exec()
			if err != nil {
				errlog.Println(err)
			}
		}

		if execFlag {
			err = tx.Commit()
			if err != nil {
				errlog.Println(err)
			}
		}

		errlog.Println("발송 메세지 수신 끝 ( ", userid, ") : ", len(msg), startTime)

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

func ReqPqTest(c *gin.Context){
	// errlog := config.Stdlog
	// execFlag := false

	// tx, err := databasepool.DB.Begin()
	// if err != nil {
	// 	errlog.Println(err)
	// }
	// defer tx.Rollback()

	// ftStmt, err := tx.Prepare(pq.CopyIn("dhn_request", kaocommon.GetReqColumnPq(kaocommon.FtReqColumn{})...))
	// if err != nil {
	// 	errlog.Println(err)
	// }
	// defer ftStmt.Close()

	// atStmt, err := tx.Prepare(pq.CopyIn("dhn_request_at", kaocommon.GetReqColumnPq(kaocommon.AtReqColumn{})...))
	// if err != nil {
	// 	errlog.Println(err)
	// }
	// defer atStmt.Close()

	// msgStmt, err := tx.Prepare(pq.CopyIn("dhn_result", kaocommon.GetReqColumnPq(kaocommon.MsgReqColumn{})...))
	// if err != nil {
	// 	errlog.Println(err)
	// }
	// defer atStmt.Close()

	// msgTempStmt, _ := tx.Prepare(pq.CopyIn("dhn_result_temp", kaocommon.GetReqColumnPq(kaocommon.MsgReqColumn{})...))
	// if err != nil {
	// 	errlog.Println(err)
	// }
	// defer msgTempStmt.Close()

	// ftValues := []kaocommon.FtReqColumn{}
	// atValues := []kaocommon.AtReqColumn{}
	// msgValues := []kaocommon.MsgReqColumn{}

	// if len(ftValues) > 0 {
	// 	for _, data := range ftValues {
	// 		_, err := ftStmt.Exec(data.Msgid,data.Userid,data.Ad_flag,data.Button1,data.Button2,data.Button3,data.Button4,data.Button5,data.Image_link,data.Image_url,data.Message_type,data.Msg,data.Msg_sms,data.Only_sms,data.P_com,data.P_invoice,data.Phn,data.Profile,data.Reg_dt,data.Remark1,data.Remark2,data.Remark3,data.Remark4,data.Remark5,data.Reserve_dt,data.S_code,data.Sms_kind,data.Sms_lms_tit,data.Sms_sender,data.Tmpl_id,data.Wide,data.Send_group,data.Supplement,data.Price,data.Currency_type,data.Header,data.Carousel,data.Att_coupon,data.Attachments)
	// 		if err != nil {
	// 			errlog.Println(err)
	// 		}
	// 	}
	// 	execFlag = true
	// }

	// if len(atValues) > 0 {
	// 	for _, data := range atValues {
	// 		_, err := atStmt.Exec(data.Msgid,data.Userid,data.Ad_flag,data.Button1,data.Button2,data.Button3,data.Button4,data.Button5,data.Image_link,data.Image_url,data.Message_type,data.Msg,data.Msg_sms,data.Only_sms,data.P_com,data.P_invoice,data.Phn,data.Profile,data.Reg_dt,data.Remark1,data.Remark2,data.Remark3,data.Remark4,data.Remark5,data.Reserve_dt,data.S_code,data.Sms_kind,data.Sms_lms_tit,data.Sms_sender,data.Tmpl_id,data.Wide,data.Send_group,data.Supplement,data.Price,data.Currency_type)
	// 		if err != nil {
	// 			errlog.Println(err)
	// 		}
	// 	}
	// 	execFlag = true
	// }
	

	// if (execFlag){
	// 	_, err = ftStmt.Exec()
	// 	if err != nil {
	// 		errlog.Println(err)
	// 	}

	// 	err = tx.Commit()
	// 	if err != nil {
	// 		errlog.Println(err)
	// 	}
	// }
}






