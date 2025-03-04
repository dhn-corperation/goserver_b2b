package kaoreqreceive

import (
	"fmt"
	"time"
	"strconv"
	s "strings"

	ss "mycs/internal/structs"
	cm "mycs/internal/commons"
	config "mycs/configs"
	databasepool "mycs/configs/databasepool"

	_ "github.com/go-sql-driver/mysql"
	"github.com/valyala/fasthttp"
	"github.com/goccy/go-json"
)

//언젠가는 다른곳으로 위치를 옮겨야 함
var SecretKey = "9b4dabe9d4fed126a58f8639846143c7"

func ReqReceive(c *fasthttp.RequestCtx) {
	errlog := config.Stdlog

	userid := string(c.Request.Header.Peek("userid"))
	userip := c.RemoteIP().String()
	isValidation := false

	// 허가된 userid 인지 테이블에서 확인
	sqlstr := `
		select 
			count(1) as cnt 
		from
			DHN_CLIENT_LIST
		where
			user_id = ?
			and ip = ?
			and use_flag = 'Y'`

	var cnt int
	err := databasepool.DB.QueryRow(sqlstr, userid, userip).Scan(&cnt)
	if err != nil { errlog.Println(err) }

	if cnt > 0 { 
		isValidation = true 
	} else {
		errlog.Println("허용되지 않은 사용자 및 아이피에서 발송 요청!! (userid : ", userid, "/ ip : ", userip, ")")
	}

	var startNow = time.Now()
	var startTime = fmt.Sprintf("%02d:%02d:%02d", startNow.Hour(), startNow.Minute(), startNow.Second())

	if isValidation {

		var msg []ss.Reqtable
		var resp ss.ReceiveRes
		
		if err1 := json.Unmarshal(c.PostBody(), &msg); err1 != nil {
			errlog.Println(err1)
			resp := ss.ReceiveRes{
				Code : "01",
				Message : "데이터 맵핑 실패",
			}
			respJson, _ := json.Marshal(resp)
			c.SetContentType("application/json")
			c.SetStatusCode(fasthttp.StatusBadRequest)
			c.SetBody(respJson)
			return
		}

		errlog.Println("발송 메세지 수신 시작 ( ", userid, ") : ", len(msg), startTime)

		resp = receiveProc(msg, userid)

		respJson, _ := json.Marshal(resp)

		c.SetContentType("application/json")
		c.SetStatusCode(fasthttp.StatusOK)
		c.SetBody(respJson)
	} else {
		resp := ss.ReceiveRes{
			Code : "01",
			Message : "허용되지 않은 사용자 입니다 / userid : " + userid + " / ip : " + userip,
		}
		respJson, _ := json.Marshal(resp)

		c.SetContentType("application/json")
		c.SetStatusCode(fasthttp.StatusNotAcceptable)
		c.SetBody(respJson)
	}
}


func ReqReceiveSocket(dt []byte, cIp string) []byte {
	errlog := config.Stdlog

	var msg []ss.Reqtable
	var resp ss.ReceiveRes
	
	if err := json.Unmarshal(dt, &msg); err != nil {
		errlog.Println(err)
		resp = ss.ReceiveRes{
			Code : "01",
			Message : "데이터 맵핑 실패",
		}
		respJson, _ := json.Marshal(resp)
		return respJson
	}

	resp = receiveProc(msg, cIp)

	respJson, _ := json.Marshal(resp)
	return respJson
}

func receiveProc(msg []ss.Reqtable, userid string) ss.ReceiveRes{
	var resp ss.ReceiveRes

	var cnt int

	var startNow = time.Now()
	var startTime = fmt.Sprintf("%02d:%02d:%02d", startNow.Hour(), startNow.Minute(), startNow.Second())

	var atCnt int16 = 0
	var ftCnt int16 = 0
	var msgCnt int16 = 0
	var duplCnt int16 = 0

	errlog := config.Stdlog

	ftColumn := cm.GetReqFtColumn()
	atColumn := cm.GetReqAtColumn()
	msgColumn := cm.GetReqMsgColumn()

	ftColumnStr := s.Join(ftColumn, ",")
	atColumnStr := s.Join(atColumn, ",")
	msgColumnStr := s.Join(msgColumn, ",")

	ftQmarkStr := cm.GetQuestionMark(ftColumn)
	atQmarkStr := cm.GetQuestionMark(atColumn)
	msgQmarkStr := cm.GetQuestionMark(msgColumn)

	reqinsStrs := []string{}
	//친구톡 value interface 배열 생성
	reqinsValues := []interface{}{}

	atreqinsStrs := []string{}
	//알림톡 value interface 배열 생성
	atreqinsValues := []interface{}{}

	resinsStrs := []string{}
	//문자 value interface 배열 생성
	resinsValues := []interface{}{}

	insStrs := []string{}
	//수신 value interface 배열 생성
	insValues := []interface{}{}

	//중복 아이디 값 배열 생성
	duplMsgid := []string{}

	//친구톡 insert 컬럼 셋팅
	reqinsQuery := `insert into DHN_REQUEST(`+ftColumnStr+`) values %s`

	//알림톡 insert 컬럼 셋팅
	atreqinsQuery := `insert into DHN_REQUEST_AT(`+atColumnStr+`) values %s`

	//문자 insert 컬럼 셋팅
	resinsQuery := `insert IGNORE into DHN_RESULT(`+msgColumnStr+`) values %s`

	//수신 insert 컬럼 셋팅
	insQuery := `insert IGNORE into DHN_RECEPTION(msgid, userid) values %s`

	//temp 테이블 컬럼 셋팅(DHN_RESULT_TEMP : 에러 시 데이터 유실을 막기 위한 테이블)
	resinstempquery := `insert IGNORE into DHN_RESULT_TEMP(`+msgColumnStr+`) values %s`
	
	errlog.Println("발송 메세지 수신 시작 ( ", userid, ") : ", len(msg), startTime)

	//맵핑한 데이터 row 처리
	for i, _ := range msg {
		msgidsSql := "select count(1) as cnt from DHN_RECEPTION where userid = '" + userid + "' and msgid = '" + msg[i].Msgid + "' and insert_date > date_sub(now(), interval 1 day)"

		err := databasepool.DB.QueryRow(msgidsSql).Scan(&cnt)
		if err != nil { 
			errlog.Println("DHN_RECEPTION 테이블 조회 에러 : ", err)
			errlog.Println("DHN_RECEPTION 테이블 조회 에러 (userid : ", userid, "/ msgid : ", msg[i].Msgid, ")")
			continue
		}

		if cnt > 0 { 
			errlog.Println("중복된 발송아이디 발송 요청!! (userid : ", userid, "/ msgid : ", msg[i].Msgid, ")")
			duplMsgid = append(duplMsgid, msg[i].Msgid)
			duplCnt++
			continue
		}

		insStrs = append(insStrs, "(?,?)")
		insValues = append(insValues, msg[i].Msgid)
		insValues = append(insValues, userid)

		var nonce string
		if len(msg[i].Crypto) > 0 {
			nonce = s.Split(msg[i].Crypto, ",")[0]
		}

		// var processedMsg string
		// var err error
		var smsKind = msg[i].Smskind
		// if s.Contains(s.ToLower(msg[i].Crypto), "msg") && len(msg[i].Msgsms) > 0 {
		// 	processedMsg, err = cm.RemoveWs(cm.AES256GSMDecrypt([]byte(SecretKey), msg[i].Msgsms, nonce))
		// } else {
		// 	processedMsg, err = cm.RemoveWs(msg[i].Msgsms)
		// }
		// if err != nil {
		// 	errlog.Println("RemoveWs 에러 : ", err)
		// } else {
		// 	euckrLength, err := cm.LengthInEUCKR(processedMsg)
		// 	if err != nil {
		// 		errlog.Println("LengthInEUCKR 에러 : ", err)
		// 	}
		// 	if euckrLength <= 90 {
		// 		smsKind = "S"
		// 	} else if euckrLength > 90 && msg[i].Pinvoice == "" {
		// 		smsKind = "L"
		// 	} else {
		// 		smsKind = "M"
		// 	}
		// }

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
			if s.Contains(s.ToLower(msg[i].Crypto), "msg") {
				reqinsValues = append(reqinsValues, cm.AES256GSMDecrypt([]byte(SecretKey), msg[i].Msg, nonce))
			} else {
				reqinsValues = append(reqinsValues, msg[i].Msg)
			}
			if s.Contains(s.ToLower(msg[i].Crypto), "msg") && len(msg[i].Msgsms) > 0 {
				reqinsValues = append(reqinsValues, cm.AES256GSMDecrypt([]byte(SecretKey), msg[i].Msgsms, nonce))
			} else {
				reqinsValues = append(reqinsValues, msg[i].Msgsms)
			}
			reqinsValues = append(reqinsValues, msg[i].Onlysms)
			if s.Contains(s.ToLower(msg[i].Crypto), "phn") && msg[i].Phn != "" {
				reqinsValues = append(reqinsValues, cm.AES256GSMDecrypt([]byte(SecretKey), msg[i].Phn, nonce))
			} else {
				reqinsValues = append(reqinsValues, msg[i].Phn)
			}
			
			if s.Contains(s.ToLower(msg[i].Crypto), "profile") && len(msg[i].Profile) > 0 {
				reqinsValues = append(reqinsValues, cm.AES256GSMDecrypt([]byte(SecretKey), msg[i].Profile, nonce))
			} else {
				reqinsValues = append(reqinsValues, msg[i].Profile)
			}
			reqinsValues = append(reqinsValues, msg[i].Pcom)
			reqinsValues = append(reqinsValues, msg[i].Pinvoice)
			reqinsValues = append(reqinsValues, msg[i].Regdt)
			reqinsValues = append(reqinsValues, msg[i].Remark1)
			reqinsValues = append(reqinsValues, msg[i].Remark2)
			reqinsValues = append(reqinsValues, msg[i].Remark3)
			reqinsValues = append(reqinsValues, msg[i].Remark4)
			reqinsValues = append(reqinsValues, msg[i].Remark5)
			reqinsValues = append(reqinsValues, msg[i].Reservedt)
			reqinsValues = append(reqinsValues, smsKind)
			if s.Contains(s.ToLower(msg[i].Crypto), "smslmstit") && len(msg[i].Smslmstit) > 0 {
				reqinsValues = append(reqinsValues, cm.AES256GSMDecrypt([]byte(SecretKey), msg[i].Smslmstit, nonce))
			} else {
				reqinsValues = append(reqinsValues, msg[i].Smslmstit)
			}
			if s.Contains(s.ToLower(msg[i].Crypto), "smssender") && len(msg[i].Smssender) > 0 {
				reqinsValues = append(reqinsValues, cm.AES256GSMDecrypt([]byte(SecretKey), msg[i].Smssender, nonce))
			} else {
				reqinsValues = append(reqinsValues, msg[i].Smssender)
			}
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
			reqinsValues = append(reqinsValues, msg[i].Attachments)
			reqinsValues = append(reqinsValues, msg[i].Carousel)
			reqinsValues = append(reqinsValues, msg[i].Att_items)
			reqinsValues = append(reqinsValues, msg[i].Att_coupon)
			reqinsValues = append(reqinsValues, msg[i].MmsImageId)
			ftCnt++
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
			if s.Contains(s.ToLower(msg[i].Crypto), "msg") {
				resinsValues = append(resinsValues, cm.AES256GSMDecrypt([]byte(SecretKey), msg[i].Msg, nonce))
			} else {
				resinsValues = append(resinsValues, msg[i].Msg)
			}

			if s.Contains(s.ToLower(msg[i].Crypto), "msg") && len(msg[i].Msgsms) > 0 {
				resinsValues = append(resinsValues, cm.AES256GSMDecrypt([]byte(SecretKey), msg[i].Msgsms, nonce))
			} else {
				resinsValues = append(resinsValues, msg[i].Msgsms)
			}
			resinsValues = append(resinsValues, msg[i].Onlysms)
			resinsValues = append(resinsValues, msg[i].Pcom)
			resinsValues = append(resinsValues, msg[i].Pinvoice)
			if s.Contains(s.ToLower(msg[i].Crypto), "phn") && msg[i].Phn != "" {
				resinsValues = append(resinsValues, cm.AES256GSMDecrypt([]byte(SecretKey), msg[i].Phn, nonce))
			} else {
				resinsValues = append(resinsValues, msg[i].Phn)
			}
			if s.Contains(s.ToLower(msg[i].Crypto), "profile") && len(msg[i].Profile) > 0 {
				resinsValues = append(resinsValues, cm.AES256GSMDecrypt([]byte(SecretKey), msg[i].Profile, nonce))
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
			resinsValues = append(resinsValues, smsKind)
			if s.Contains(s.ToLower(msg[i].Crypto), "smslmstit") && len(msg[i].Smslmstit) > 0 {
				resinsValues = append(resinsValues, cm.AES256GSMDecrypt([]byte(SecretKey), msg[i].Smslmstit, nonce))
			} else {
				resinsValues = append(resinsValues, msg[i].Smslmstit)
			}

			if s.Contains(s.ToLower(msg[i].Crypto), "smssender") && len(msg[i].Smssender) > 0 {
				resinsValues = append(resinsValues, cm.AES256GSMDecrypt([]byte(SecretKey), msg[i].Smssender, nonce))
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
			resinsValues = append(resinsValues, nil) //attachments
			resinsValues = append(resinsValues, msg[i].Carousel)
			resinsValues = append(resinsValues, msg[i].MmsImageId)
			msgCnt++
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
			if s.Contains(s.ToLower(msg[i].Crypto), "msg") {
				atreqinsValues = append(atreqinsValues, cm.AES256GSMDecrypt([]byte(SecretKey), msg[i].Msg, nonce))
			} else {
				atreqinsValues = append(atreqinsValues, msg[i].Msg)
			}

			if s.Contains(s.ToLower(msg[i].Crypto), "msg") && len(msg[i].Msgsms) > 0 {
				atreqinsValues = append(atreqinsValues, cm.AES256GSMDecrypt([]byte(SecretKey), msg[i].Msgsms, nonce))
			} else {
				atreqinsValues = append(atreqinsValues, msg[i].Msgsms)
			}
			atreqinsValues = append(atreqinsValues, msg[i].Onlysms)
			if s.Contains(s.ToLower(msg[i].Crypto), "phn") && msg[i].Phn != "" {
				atreqinsValues = append(atreqinsValues, cm.AES256GSMDecrypt([]byte(SecretKey), msg[i].Phn, nonce))
			} else {
				atreqinsValues = append(atreqinsValues, msg[i].Phn)
			}
			if s.Contains(s.ToLower(msg[i].Crypto), "profile") && len(msg[i].Profile) > 0 {
				atreqinsValues = append(atreqinsValues, cm.AES256GSMDecrypt([]byte(SecretKey), msg[i].Profile, nonce))
			} else {
				atreqinsValues = append(atreqinsValues, msg[i].Profile)
			}
			atreqinsValues = append(atreqinsValues, msg[i].Pcom)
			atreqinsValues = append(atreqinsValues, msg[i].Pinvoice)
			atreqinsValues = append(atreqinsValues, msg[i].Regdt)
			atreqinsValues = append(atreqinsValues, msg[i].Remark1)
			atreqinsValues = append(atreqinsValues, msg[i].Remark2)
			atreqinsValues = append(atreqinsValues, msg[i].Remark3)
			atreqinsValues = append(atreqinsValues, msg[i].Remark4)
			atreqinsValues = append(atreqinsValues, msg[i].Remark5)
			atreqinsValues = append(atreqinsValues, msg[i].Reservedt)
			atreqinsValues = append(atreqinsValues, smsKind)
			if s.Contains(s.ToLower(msg[i].Crypto), "smslmstit") && len(msg[i].Smslmstit) > 0 {
				atreqinsValues = append(atreqinsValues, cm.AES256GSMDecrypt([]byte(SecretKey), msg[i].Smslmstit, nonce))
			} else {
				atreqinsValues = append(atreqinsValues, msg[i].Smslmstit)
			}

			if s.Contains(s.ToLower(msg[i].Crypto), "smssender") && len(msg[i].Smssender) > 0 {
				atreqinsValues = append(atreqinsValues, cm.AES256GSMDecrypt([]byte(SecretKey), msg[i].Smssender, nonce))
			} else {
				atreqinsValues = append(atreqinsValues, msg[i].Smssender)
			}
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
			atreqinsValues = append(atreqinsValues, msg[i].MmsImageId)
			atreqinsValues = append(atreqinsValues, msg[i].Header)
			atreqinsValues = append(atreqinsValues, msg[i].Attachments)
			atreqinsValues = append(atreqinsValues, msg[i].Link)
			atCnt++
		}

		// 500건 단위로 처리한다(클라이언트에서 1000건씩 전송하더라도 지정한 단위의 건수로 insert한다.)
		saveCount := 500
		if len(reqinsStrs) >= saveCount {
			reqinsStrs, reqinsValues = cm.InsMsg(reqinsQuery, reqinsStrs, reqinsValues)
		}

		if len(atreqinsStrs) >= saveCount {
			atreqinsStrs, atreqinsValues = cm.InsMsg(atreqinsQuery, atreqinsStrs, atreqinsValues)
		}

		if len(resinsStrs) >= saveCount {
			resinsStrs, resinsValues = cm.InsMsgTemp(resinsQuery, resinsStrs, resinsValues, true, resinstempquery)
		}

		if len(insStrs) >= saveCount {
			insStrs, insValues = cm.InsMsg(insQuery, insStrs, insValues)
		}
	}
	
	// 나머지 건수를 저장하기 위해 다시한번 정의
	if len(reqinsStrs) > 0 {
		reqinsStrs, reqinsValues = cm.InsMsg(reqinsQuery, reqinsStrs, reqinsValues)
	}

	if len(atreqinsStrs) > 0 {
		atreqinsStrs, atreqinsValues = cm.InsMsg(atreqinsQuery, atreqinsStrs, atreqinsValues)
	}

	if len(resinsStrs) > 0 {
		resinsStrs, resinsValues = cm.InsMsgTemp(resinsQuery, resinsStrs, resinsValues, true, resinstempquery)
	}

	if len(insStrs) > 0 {
		insStrs, insValues = cm.InsMsg(insQuery, insStrs, insValues)
	}

	errlog.Println("발송 메세지 수신 끝 ( ", userid, ") : ", len(msg), startTime)

	resp = ss.ReceiveRes{
		Code : "00",
		Message : "발송 요청이 완료되었습니다.",
		AtCnt : &atCnt,
		FtCnt : &ftCnt,
		MsgCnt : &msgCnt,
		DuplCnt : &duplCnt,
		DuplMsgId : &duplMsgid,
	}

	return resp
}