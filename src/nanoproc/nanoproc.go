package nanoproc

import (
	"database/sql"
	"fmt"

	config "mycs/src/kaoconfig"
	databasepool "mycs/src/kaodatabasepool"
	"mycs/src/kaocommon"

	"regexp"
	s "strings"
	"time"
	"context"
)

var procCnt int

func NanoProcess(user_id string, ci string, ctx context.Context, gFlag int) {
	config.Stdlog.Println(user_id, " Nano request Process 시작 됨.")
	procCnt = 0

	for {
		if procCnt < 5 {

			select {
			case <-ctx.Done():
				config.Stdlog.Println(user_id, " - Nano request process가 10초 후에 종료 됨.")
				time.Sleep(10 * time.Second)
				config.Stdlog.Println(user_id, " - Nano request process 종료 완료")
				return
			default:
				var count sql.NullInt64
				tickSql := `
				select
					count(1) as cnt
				from
					DHN_RESULT
				where
					result = 'P'
					and send_group is null
					and (reserve_dt IS NULL OR to_timestamp(coalesce(reserve_dt,'00000000000000'), 'YYYYMMDDHH24MISS') <= NOW())
				 	and userid = $1 `
				subQuery := ""
				tail := ""
				switch gFlag {
				case 2:	// 전화번호 010 일때만
					subQuery = " and sms_sender like '010%' "
					tail = "_G"
				case 3: // 전화번호 010 아닌 것들
					subQuery = " and sms_sender not like '010%' "
					tail = "_G"
				}
				tickSql = tickSql + subQuery + ` limit 1`

				cnterr := databasepool.DB.QueryRowContext(ctx, tickSql, user_id).Scan(&count)

				if cnterr != nil && cnterr != sql.ErrNoRows {
					config.Stdlog.Println("nanoproc.go / NanoProcess / DHN_RESULT Table - select 오류 : " + cnterr.Error())
				} else {

					if count.Int64 > 0 {

						var startNow = time.Now()
						var group_no = fmt.Sprintf("%02d%02d%02d%02d%06d", startNow.Day(), startNow.Hour(), startNow.Minute(), startNow.Second(), (startNow.Nanosecond() / 1000))

						upError := updateReqeust(ctx, group_no, user_id, subQuery)
						if upError != nil {
							config.Stdlog.Println(user_id, " nanoproc.go / NanoProcess / Nano Group No Update 오류", group_no)
						} else {
							config.Stdlog.Println(user_id, " nanoproc.go / NanoProcess / 문자 발송 처리 시작 ( ", group_no, " )")
							procCnt++
							go resProcess(ctx, group_no, user_id, tail, ci)
						}
					}
				}
			}
		}
	}
}

func updateReqeust(ctx context.Context, group_no string, user_id string, subQuery string) error {

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

	config.Stdlog.Println(user_id, "- Nano Group No Update 시작", group_no)

	gudQuery := `
	update
		DHN_RESULT dr
	set
		send_group = $1
	where (msgid, userid) in (
		select msgid, userid
		from dhn_result dr
		where result = 'P'
	  	and send_group is null
	  	and (dr.reserve_dt IS NULL OR to_timestamp(coalesce(dr.reserve_dt,'00000000000000'), 'YYYYMMDDHH24MISS') <= NOW())
	  	and userid = $2
	  	limit 500
	)`

	gudQuery = gudQuery + subQuery

	_, err = tx.ExecContext(ctx, gudQuery, group_no, user_id)

	if err != nil {
		config.Stdlog.Println("nanoproc.go / updateReqeust / ", user_id, "- Group NO Update - Select error : ( group_no : ", group_no, " / user_id : ", user_id, " ) error : ", err)
		config.Stdlog.Println(gudQuery)
		return err
	}

	err = tx.Commit()
	if err != nil {
		config.Stdlog.Println("nanoproc.go / updateReqeust / ", user_id, "- Group No Update - Commit error : ( group_no : ", group_no, " / user_id : ", user_id, " ) error : ", err)
		config.Stdlog.Println(gudQuery)
		return err
	}

	return nil
}

func resProcess(ctx context.Context, group_no string, user_id string, tail string, ci string) {
	procCnt++
	var db = databasepool.DB
	var stdlog = config.Stdlog

	defer func() {
		if err := recover(); err != nil {
			procCnt--
			stdlog.Println("nanoproc.go / resProcess / ", user_id, "-", group_no, " recover() Nano 문자 처리 중 오류 발생 : ", err)
		}
	}()

	var msgid, msg_sms, phn, sms_lms_tit, sms_kind, sms_sender, reserve_dt, mms_file1, mms_file2, mms_file3, userid, sms_len_check sql.NullString
	var phnstr string

	var resquery = `
	SELECT 
		msgid,
		msg_sms, 
		phn,
		sms_lms_tit, 
		sms_kind, 
		sms_sender,
		(case
			when reserve_dt = '00000000000000'  then 
		        now()
		    when reserve_dt is null  then 
		        now()
		    when length(trim(reserve_dt)) < 4  then 
	        	now()
		    else
		        to_timestamp(reserve_dt, 'YYYYMMDDHH24MISS')
	     end) as  reserve_dt, 
		(select file1_path from api_mms_images aa where aa.user_id = drr.userid and aa.mms_id = drr.p_invoice) as mms_file1, 
		(select file2_path from api_mms_images aa where aa.user_id = drr.userid and aa.mms_id = drr.p_invoice) as mms_file2, 
		(select file3_path from api_mms_images aa where aa.user_id = drr.userid and aa.mms_id = drr.p_invoice) as mms_file3
		,userid
		,(select max(sms_len_check) from DHN_CLIENT_LIST dcl where dcl.user_id = drr.userid) as sms_len_check 
	FROM DHN_RESULT drr 
	WHERE send_group = $1
	    and result = 'P'
        and userid = $2`

	resrows, err := db.QueryContext(ctx, resquery, group_no, user_id)

	if err != nil {
		stdlog.Println("nanoproc.go / resProcess / Result Table 조회 중 오류 발생")
		stdlog.Println(err)
		stdlog.Println(resquery)
	}
	defer resrows.Close()

	smsValues := []kaocommon.NanoReqColumn{}
	mmsValues := []kaocommon.NanoReqColumn{}

	smscnt := 0
	lmscnt := 0
	tcnt := 0
	reg, err := regexp.Compile("[^0-9]+")

	for resrows.Next() {
		resrows.Scan(&msgid, &msg_sms, &phn, &sms_lms_tit, &sms_kind, &sms_sender, &reserve_dt, &mms_file1, &mms_file2, &mms_file3, &userid, &sms_len_check)

		phnstr = phn.String

		if tcnt == 0 {
			stdlog.Println("nanoproc.go / resProcess / ", user_id, "-", group_no, "문자발송 처리 시작 : ", " Process cnt : ", procCnt)
		}

		tcnt++

		if len(smsValues) > 500 {
			insertNanoReqData(smsValues, "sms_msg" + tail)
			smsValues = []kaocommon.NanoReqColumn{}
		}

		if len(mmsValues) > 500 {
			insertNanoReqData(mmsValues, "mms_msg" + tail)
			mmsValues = []kaocommon.NanoReqColumn{}
		}

		if len(msg_sms.String) > 0 && len(sms_sender.String) > 0 { // msg_sms 와 sms_sender 에 값이 있으면 nano 발송 함.

			phnstr = reg.ReplaceAllString(phnstr, "")
			if s.HasPrefix(phnstr, "82") {
				phnstr = "0" + phnstr[2:len(phnstr)]
			}

			ms := kaocommon.RemoveSpecialChar(msg_sms.String)
			ml := len(ms)

			if s.EqualFold(sms_kind.String, "S") {
				ms = ms[1:100]
				if ml <= 90 || s.EqualFold(sms_len_check.String, "N") {
					smsValues = append(smsValues, kaocommon.NanoReqColumn{
						CALLBACK : sms_sender.String,
						PHONE : phnstr,
						MSG : ms,
						TR_SENDDATE : reserve_dt.String,
						TR_SENDSTAT : "0",
						TR_MSGTYPE : "0",
						ETC9 : msgid.String,
						ETC10 : userid.String,
						IDENTIFICATION_CODE : ci,
						ETC8 : "Y",
					})
					smscnt++
				} else {
					db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.code = '7003', dr.message = '메세지 길이 오류', dr.remark2 = date_format(now(), '%Y-%m-%d %H:%i:%S') where userid = '" + userid.String + "' and msgid = '" + msgid.String + "'")
				}
			} else if s.EqualFold(sms_kind.String, "L") || s.EqualFold(sms_kind.String, "M") {
				ml = 100
				filecnt := 0

				if len(mms_file1.String) > 0 {
					filecnt = filecnt + 1
				}

				if len(mms_file2.String) > 0 {
					filecnt = filecnt + 1
				}

				if len(mms_file3.String) > 0 {
					filecnt = filecnt + 1
				}

				slt := kaocommon.RemoveSpecialChar(sms_lms_tit.String)

				mmsValues = append(mmsValues, kaocommon.NanoReqColumn{
					CALLBACK : sms_sender.String,
					PHONE : phnstr,
					SUBJECT : slt,
					MSG : ms,
					REQDATE : reserve_dt.String,
					STATUS : "0",
					FILE_CNT : filecnt,
					FILE_PATH1 : mms_file1.String,
					FILE_PATH2 : mms_file2.String,
					FILE_PATH3 : mms_file3.String,
					ETC9 : msgid.String,
					ETC10 : userid.String,
					IDENTIFICATION_CODE : ci,
					ETC8 : "Y",
				})
				lmscnt++
			}

		} else {
			db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.code='7011', dr.message = concat(dr.message, ',문자 발송 정보 누락'),dr.remark2 = date_format(now(), '%Y-%m-%d %H:%i:%S') where userid = '" + userid.String + "' and msgid = '" + msgid.String + "'")
		}

	}

	if len(smsValues) > 0 {
		insertNanoReqData(smsValues, "sms_msg" + tail)
	}

	if len(mmsValues) > 0 {
		insertNanoReqData(mmsValues, "mms_msg" + tail)

	}

	if smscnt > 0 || lmscnt > 0 {
		stdlog.Println("nanoproc.go / resProcess / ", user_id, "-", group_no, "문자 발송 처리 완료 ( ", tcnt, " ) : SMS -", smscnt, " , LMS -", lmscnt, "  >> Process cnt : ", procCnt)
	}
	procCnt--
}

func insertNanoReqData(msgValues []kaocommon.NanoReqColumn, tableName string) {
	tx, err := databasepool.DB.Begin()
	if err != nil {
		config.Stdlog.Println("nanoproc.go / insertNanoReqData / ", tableName, " / 트랜젝션 초기화 실패 ", err)
	}
	defer tx.Rollback()

	tableName = s.ToLower(tableName)

	var stmt *sql.Stmt
	var stmtSql string

	if s.Contains(tableName, "sms") {
		stmtSql = "insert into "+tableName+"(tr_callback,tr_phone,tr_msg,tr_senddate,tr_sendstat,tr_msgtype,tr_etc9,tr_etc10,tr_identification_code,tr_etc8) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"
		stmt, err = tx.Prepare(stmtSql)
		if err != nil {
			config.Stdlog.Println("nanoproc.go / insertNanoReqData / ", tableName, "/ stmt 초기화 실패 ", err)
			return
		}

		for _, data := range msgValues {
			_, err = stmt.Exec(data.CALLBACK, data.PHONE, data.MSG, data.TR_SENDDATE, data.TR_SENDSTAT, data.TR_MSGTYPE, data.ETC9, data.ETC10, data.IDENTIFICATION_CODE, data.ETC8)
			if err != nil {
				config.Stdlog.Println("nanoproc.go / insertNanoReqData / ", tableName," / stmt personal Exec ", err)
			}
		}
	} else if s.Contains(tableName, "mms") {
		stmtSql = "insert into "+tableName+"(callback,phone,subject,msg,reqdate,status,file_cnt,file_path1,file_path2,file_path3,etc9,etc10,identification_code,etc8) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)"
		stmt, err := tx.Prepare(stmtSql)
		if err != nil {
			config.Stdlog.Println("nanoproc.go / insertNanoReqData / ", tableName, "/ stmt 초기화 실패 ", err)
			return
		}

		for _, data := range msgValues {
			_, err = stmt.Exec(data.CALLBACK, data.PHONE, data.SUBJECT, data.MSG, data.REQDATE, data.STATUS, data.FILE_CNT, data.FILE_PATH1, data.FILE_PATH2, data.FILE_PATH3, data.ETC9, data.ETC10, data.IDENTIFICATION_CODE, data.ETC8)
			if err != nil {
				config.Stdlog.Println("nanoproc.go / insertNanoReqData / ", tableName," / stmt personal Exec ", err)
			}
		}
	}
	
	
	stmt.Close()
	err = tx.Commit()

	if err != nil {
		if s.Contains(tableName, "sms") {
			for _, data := range msgValues {
				_, err = databasepool.DB.Exec(stmtSql, data.CALLBACK, data.PHONE, data.MSG, data.TR_SENDDATE, data.TR_SENDSTAT, data.TR_MSGTYPE, data.ETC9, data.ETC10, data.IDENTIFICATION_CODE, data.ETC8)
				if err != nil {
					checkErr(err, data.ETC9, data.ETC10, tableName)
				}
			}
		} else if s.Contains(tableName, "mms") {
			for _, data := range msgValues {
				_, err = databasepool.DB.Exec(stmtSql, data.CALLBACK, data.PHONE, data.SUBJECT, data.MSG, data.REQDATE, data.STATUS, data.FILE_CNT, data.FILE_PATH1, data.FILE_PATH2, data.FILE_PATH3, data.ETC9, data.ETC10, data.IDENTIFICATION_CODE, data.ETC8)
				if err != nil {
					checkErr(err, data.ETC9, data.ETC10, tableName)
				}
			}
		}
	} else {
		config.Stdlog.Println(msgValues[0].ETC10, "- 나노 MMS Table Insert 처리 : ", len(msgValues), " - ", tableName)
	}
}

func checkErr(err error, cbMsgId string, userId string, tableName string){
	config.Stdlog.Println("nanoproc.go / insertNanoReqData / ", tableName, " / stmt commit ", err)
	config.Stdlog.Println(userId, "- 나노 ", tableName, " Insert 처리 중 오류 발생 : "+err.Error(), " - DHN Msg Key : ", cbMsgId)
	errcodemsg := err.Error()
	if s.Index(errcodemsg, "invalid byte sequence") > 0 {
		databasepool.DB.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.code='7069', dr.message = concat(dr.message, ',부적절한 문자사용'),dr.remark2 = TO_CHAR(now(), 'YYYY-MM-DD H:i:s') where userid = $1 and msgid = $2", userId, cbMsgId)
	}
}