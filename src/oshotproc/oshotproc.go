package oshotproc

import (
	"database/sql"
	"fmt"

	config "mycs/src/kaoconfig"
	databasepool "mycs/src/kaodatabasepool"

	"encoding/hex"
	"regexp"
	s "strings"
	"time"
	"unicode/utf8"

	"context"
)

var procCnt int

func OshotProcess(user_id string, ctx context.Context) {
	config.Stdlog.Println(user_id, " Oshot Process 시작 됨.")
	procCnt = 0
	for {

		if procCnt < 5 {

			select {
			case <-ctx.Done():

				config.Stdlog.Println(user_id, " - Oshot process가 10초 후에 종료 됨.")
				time.Sleep(10 * time.Second)
				config.Stdlog.Println(user_id, " - Oshot process 종료 완료")
				return
			default:

				var count sql.NullInt64
				tickSql := `
				select
					count(1) as cnt
				from
					DHN_RESULT
				where
					dr.result = 'P'
					and send_group is null
					and (reserve_dt IS NULL OR to_timestamp(coalesce(reserve_dt,'00000000000000'), 'YYYYMMDDHH24MISS') <= NOW())
					and userid = $1
				limit 1`
				cnterr := databasepool.DB.QueryRowContext(ctx, tickSql, user_id).Scan(&count)

				if cnterr != nil && cnterr != sql.ErrNoRows {
					config.Stdlog.Println("oshotproc.go / DHN_RESULT Table - select 오류 : " + cnterr.Error())
				} else {
					if count.Int64 > 0 {
						var startNow = time.Now()
						var group_no = fmt.Sprintf("%02d%02d%02d%02d%06d", startNow.Day(), startNow.Hour(), startNow.Minute(), startNow.Second(), (startNow.Nanosecond() / 1000))

						upError := updateReqeust(ctx, group_no, user_id)
						if upError != nil {
							config.Stdlog.Println(user_id , "oshotproc.go / Group No Update 오류", group_no)
						} else {
							config.Stdlog.Println(user_id, "oshotproc.go / 문자 발송 처리 시작 ( ", group_no, " )")
							procCnt++
							go resProcess(ctx, group_no, user_id)
						}
					}
				}
			}

		}
	}
}

func updateReqeust(ctx context.Context, group_no string, user_id string) error {

	tx, err := databasepool.DB.Begin()
	if err != nil {
		return err
	}

	defer func() error {
		if err != nil {
			tx.Rollback()
			return err
		}
		err = tx.Commit()
		return err
	}()

	config.Stdlog.Println("oshotproc.go / ", user_id, "- 스마트미 Group No Update 시작", group_no)

	gudQuery := `
	update DHN_RESULT dr
	set	send_group = $1
	where result = 'P'
	  and send_group is null
	  and (dr.reserve_dt IS NULL OR to_timestamp(coalesce(dr.reserve_dt,'00000000000000'), 'YYYYMMDDHH24MISS') <= NOW())
	  and userid = $2
	LIMIT 500
	`
	_, err = tx.ExecContext(ctx, gudQuery, group_no, user_id)

	if err != nil {
		config.Stdlog.Println("oshotproc.go / ", user_id, "- Group No Update - Select error : ( group_no : ", group_no, " / user_id : ", user_id, " ) : ", err)
		config.Stdlog.Println(gudQuery)
		return err
	}

	err = tx.Commit()
	if err != nil {
		config.Stdlog.Println("oshotproc.go / ", user_id, "- Group No Update - Commit error : ( group_no : ", group_no, " / user_id : ", user_id, " ) : ", err)
		config.Stdlog.Println(gudQuery)
		return err
	}

	return nil
}

func resProcess(ctx context.Context, group_no string, user_id string) {

	var db = databasepool.DB
	var stdlog = config.Stdlog

	defer func() {
		if err := recover(); err != nil {
			procCnt--
			stdlog.Println(user_id, "- ", group_no, "recover() Oshot 처리 중 오류 발생 : ", err)
		}
	}()

	// var msgid, code, message, message_type, msg_sms, phn, remark1, remark2, result, sms_lms_tit, sms_kind, sms_sender, res_dt, reserve_dt, mms_file1, mms_file2, mms_file3, userid, sms_len_check, oshot sql.NullString
	var msgid, msg_sms, phn, sms_lms_tit, sms_kind, sms_sender, mms_file1, mms_file2, mms_file3, userid, sms_len_check, oshot sql.NullString
	var msgLen sql.NullInt64
	var phnstr string

	ossmsStrs := []string{}
	ossmsValues := []interface{}{}

	osmmsStrs := []string{}
	osmmsValues := []interface{}{}

	var resquery = `
	SELECT
	    msgid,
	    CASE WHEN sms_kind = 'S' THEN 
	        SUBSTRING(trim(msg_sms), 1, 100)
	    ELSE 
	        trim(msg_sms)
	    END AS msg_sms, 
	    phn,
	    trim(sms_lms_tit) AS sms_lms_tit, 
	    sms_kind, 
	    sms_sender,
	    (SELECT file1_path FROM api_mms_images aa WHERE aa.user_id = drr.userid AND aa.mms_id = drr.p_invoice) AS mms_file1, 
	    (SELECT file2_path FROM api_mms_images aa WHERE aa.user_id = drr.userid AND aa.mms_id = drr.p_invoice) AS mms_file2, 
	    (SELECT file3_path FROM api_mms_images aa WHERE aa.user_id = drr.userid AND aa.mms_id = drr.p_invoice) AS mms_file3,
	    CASE WHEN sms_kind = 'S' THEN LENGTH(trim(msg_sms)) ELSE 100 END AS msg_len,
	    userid,
	    (SELECT MAX(sms_len_check) FROM DHN_CLIENT_LIST dcl WHERE dcl.user_id = drr.userid) AS sms_len_check,
	    (SELECT lower(MAX(dest)) FROM DHN_CLIENT_LIST dcl WHERE dcl.user_id = drr.userid) AS oshot  
	FROM DHN_RESULT drr 
	WHERE send_group = $1
	  AND result = 'P'
	  AND userid = $2
	ORDER BY userid
	`

	resrows, err := db.QueryContext(ctx, resquery, group_no, user_id)

	if err != nil {
		stdlog.Println("Result Table 조회 중 오류 발생")
		stdlog.Println(err)
		stdlog.Println(resquery)
		return
	}
	defer resrows.Close()

	smsValues := []kaocommon.OshotReqColumn{}
	mmsValues := []kaocommon.OshotReqColumn{}

	smscnt := 0
	lmscnt := 0
	tcnt := 0
	reg, err := regexp.Compile("[^0-9]+")
	preOshot := ""

	for resrows.Next() {
		// resrows.Scan(&msgid, &code, &message, &message_type, &msg_sms, &phn, &remark1, &remark2, &result, &sms_lms_tit, &sms_kind, &sms_sender, &res_dt, &reserve_dt, &mms_file1, &mms_file2, &mms_file3, &msgLen, &userid, &sms_len_check, &oshot)
		resrows.Scan(&msgid, &msg_sms, &phn, &sms_lms_tit, &sms_kind, &sms_sender, &mms_file1, &mms_file2, &mms_file3, &msgLen, &userid, &sms_len_check, &oshot)

		phnstr = phn.String

		if tcnt == 0 {
			stdlog.Println(user_id, "-", group_no, "문자발송 처리 시작 : ", " Process cnt : ", procCnt)
			preOshot = oshot.String
		}

		tcnt++

		if len(smsValues) > 500 || preOshot != oshot.String {
			insertOshotReqData(smsValues, preOshot+"sms")
			smsValues = []kaocommon.OshotReqColumn{}
		}

		if len(mmsValues) > 500 || preOshot != oshot.String {
			insertOshotReqData(mmsValues, preOshot+"sms")
			mmsValues = []kaocommon.OshotReqColumn{}
		}

		// 알림톡 발송 성공 혹은 문자 발송이 아니면
		// API_RESULT 성공 처리 함.
		if len(msg_sms.String) > 0 && len(sms_sender.String) > 0 { // msg_sms 가 와 sms_sender 에 값이 있으면 Oshot 발송 함.

			phnstr = reg.ReplaceAllString(phnstr, "")
			if s.HasPrefix(phnstr, "82") {
				phnstr = "0" + phnstr[2:len(phnstr)]
			}

			if s.EqualFold(sms_kind.String, "S") {
				if msgLen.Int64 <= 90 || s.EqualFold(sms_len_check.String, "N") {
					smsValues = append(smsValues, kaocommon.OshotReqColumn{
						Sender : sms_sender.String,
						Receiver : phnstr,
						Msg : msg_sms.String,
						Url : "",
						CbMsgId : msgid.String,
						UserId : userid.String,
					})
					smscnt++
				} else {
					db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.code = '7003', dr.message = '메세지 길이 오류', dr.remark2 = date_format(now(), '%Y-%m-%d %H:%i:%S') where userid = '" + userid.String + "' and msgid = '" + msgid.String + "'")
				}
			} else if s.EqualFold(sms_kind.String, "L") || s.EqualFold(sms_kind.String, "M") {
				mmsValues = append(smsValues, kaocommon.OshotReqColumn{
					MsgGroupID : group_no,
					Sender : sms_sender.String,
					Receiver : phnstr,
					Subject : sms_lms_tit.String,
					Msg : msg_sms.String,
					FilePath1 : mms_file1.String,
					FilePath2 : mms_file2.String,
					FilePath3 : mms_file3.String,
					CbMsgId : msgid.String,
					UserId : userid.String,
				})
				lmscnt++
			}

			preOshot = oshot.String
		} else {
			db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.code='7011', dr.message = concat(dr.message, ',문자 발송 정보 누락'),dr.remark2 = date_format(now(), '%Y-%m-%d %H:%i:%S') where userid = '" + userid.String + "' and msgid = '" + msgid.String + "'")
		}

	}

	if len(smsValues) > 0 {
		insertOshotReqData(smsValues, preOshot+"sms")
	}

	if len(mmsValues) > 0 {
		insertOshotReqData(mmsValues, preOshot+"sms")
	}

	if scnt > 0 || smscnt > 0 || lmscnt > 0 || fcnt > 0 {
		stdlog.Println(user_id, "-", group_no, "문자 발송 처리 완료 ( ", tcnt, " ) : SMS -", smscnt, " , LMS -", lmscnt, "  >> Process cnt : ", procCnt)
	}
	procCnt--
}

func stringSplit(str string, lencnt int) string {
	b := []byte(str)
	idx := 0
	for i := 0; i < lencnt; i++ {
		_, size := utf8.DecodeRune(b[idx:])
		idx += size
	}
	return str[:idx]
}

func insertOshotReqData(msgValues []kaocommon.OshotReqColumn, tableName string) {
	tx, err := databasepool.DB.Begin()
	if err != nil {
		config.Stdlog.Println("oshotproc.go / insertOshotReqData / ", tableName, " / 트랜젝션 초기화 실패 ", err)
	}
	defer tx.Rollback()

	tableName = tableName.ToLower(tableName)

	var stmt *sql.Stmt
	var stmtSql string

	if tableName.Contains("sms") {
		stmtSql = "insert into "+tableName+"(Sender,Receiver,Msg,URL,cb_msg_id,userid ) values ($1, $2, $3, $4, $5, $6)"
		stmt, err = tx.Prepare(stmtSql)
		if err != nil {
			config.Stdlog.Println("oshotproc.go / insertOshotReqData / ", tableName, "/ stmt 초기화 실패 ", err)
			return
		}

		for _, data := range msgValues {
			_, err = stmt.Exec(data.Sender, data.Receiver, data.Msg, "", data.CbMsgId, data.UserId)
			if err != nil {
				config.Stdlog.Println("oshotproc.go / insertOshotReqData / ", tableName," / stmt personal Exec ", err)
			}
		}
	} else if tableName.Contains("mms") {
		stmtSql = "insert into "+tableName+"(MsgGroupID, Sender, Receiver, Subject, Msg, File_Path1, File_Path2, File_Path3, cb_msg_id, userid) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"
		stmt, err := tx.Prepare(stmtSql)
		if err != nil {
			config.Stdlog.Println("oshotproc.go / insertOshotReqData / ", tableName, "/ stmt 초기화 실패 ", err)
			return
		}

		for _, data := range msgValues {
			_, err = stmt.Exec(data.MsgGroupID, data.Sender, data.Receiver, data.Subject, data.Msg, data.FilePath1, data.FilePath2, data.FilePath3, data.CbMsgId, data.UserId)
			if err != nil {
				config.Stdlog.Println("oshotproc.go / insertOshotReqData / ", tableName," / stmt personal Exec ", err)
			}
		}
	}
	
	
	stmt.Close()
	err = tx.Commit()

	if err != nil {
		func checkErr(err error, cbMsgId string, userId string, tableName string){
			config.Stdlog.Println("oshotproc.go / insertOshotReqData / ", tableName, " / stmt commit ", err)
			config.Stdlog.Println(userId, "- 스마트미 ", tableName, " Insert 처리 중 오류 발생 : "+err.Error(), " - DHN Msg Key : ", cbMsgId)
			errcodemsg := err.Error()
			if s.Index(errcodemsg, "1366") > 0 {
				db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.code='7069', dr.message = concat(dr.message, ',부적절한 문자사용'),dr.remark2 = TO_CHAR(now(), 'YYYY-MM-DD H:i:s') where userid = $1 and msgid = $2", userId, cbMsgId)
			}
		}
		if tableName.Contains("sms") {
			for _, data := range msgValues {
				_, err = databasepool.DB.Exec(stmtSql, data.Sender, data.Receiver, data.Msg, "", data.CbMsgId, data.UserId)
				if err != nil {
					checkErr(err, data.CbMsgId, data.Userid, UserId)
				}
			}
		} else if tableName.Contains("mms") {
			for _, data := range msgValues {
				_, err = databasepool.DB.Exec(stmtSql, data.MsgGroupID, data.Sender, data.Receiver, data.Subject, data.Msg, data.FilePath1, data.FilePath2, data.FilePath3, data.CbMsgId, data.UserId)
				if err != nil {
					checkErr(err, data.CbMsgId, data.UserId, tableName)
				}
			}
		}
	} else {
		config.Stdlog.Println(msgValues[0].UserId, "- 스마트미 MMS Table Insert 처리 : ", len(msgValues), " - ", tableName)
	}
}
