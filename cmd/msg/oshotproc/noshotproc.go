package oshotproc

import (
	"fmt"
	"time"
	"regexp"
	"strconv"
	"context"
	s "strings"
	"database/sql"

	config "mycs/configs"
	databasepool "mycs/configs/databasepool"
)

func NOshotProcess(user_id string, ctx context.Context) {
	config.Stdlog.Println(user_id, " - NOshot Proc Process 시작 됨.")
	procCnt := 0
	for {

		if procCnt < 5 {

			select {
			case <-ctx.Done():

				uid := ctx.Value("user_id")
				config.Stdlog.Println(uid, " - NOshot Proc process가 10초 후에 종료 됨.")
				time.Sleep(10 * time.Second)
				config.Stdlog.Println(uid, " - NOshot Proc process 종료 완료")
				return
			default:

				var count sql.NullInt64
				tickSql := `
				select
					count(msgid) as cnt
				from
					DHN_RESULT dr
				where
					dr.result = 'P'
					and dr.send_group is null
					and ifnull(dr.reserve_dt, '00000000000000') <= date_format(now(), '%Y%m%d%H%i%S')
					and userid = ?
				limit 1
					`
				cnterr := databasepool.DB.QueryRowContext(ctx, tickSql, user_id).Scan(&count)

				if cnterr != nil && cnterr != sql.ErrNoRows {
					config.Stdlog.Println(user_id, " - NOshot Proc DHN_RESULT Table - select error : " + cnterr.Error())
					time.Sleep(10 * time.Second)
				} else {
					if count.Int64 > 0 {
						var startNow = time.Now()
						var group_no = fmt.Sprintf("%02d%02d%02d%02d%06d", startNow.Day(), startNow.Hour(), startNow.Minute(), startNow.Second(), (startNow.Nanosecond() / 1000))

						upError := nupdateReqeust(ctx, group_no, user_id)
						if upError != nil {
							config.Stdlog.Println(user_id, " - NOshot Proc Group No Update error : ", upError, " / group_no : ", group_no)
						} else {
							go func() {
								procCnt++
								config.Stdlog.Println(user_id, " - NOshot 발송 처리 시작 ( ", group_no, " ) : ( Proc Cnt :", procCnt, ") - START")
								defer func() {
									procCnt--
								}()
								nresProcess(ctx, group_no, user_id, procCnt)
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

func nupdateReqeust(ctx context.Context, group_no, user_id string) error {

	tx, err := databasepool.DB.Begin()
	if err != nil {
		config.Stdlog.Println(err)
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

	config.Stdlog.Println(user_id, " - NOshot Proc Group No Update 시작", group_no)

	gudQuery := `
	update DHN_RESULT dr
	set	send_group = ?
	where result = 'P'
	  and send_group is null
	  and ifnull(reserve_dt, '00000000000000') <= date_format(now(), '%Y%m%d%H%i%S')
	  and userid = ?
	LIMIT 500
	`
	_, err = tx.ExecContext(ctx, gudQuery, group_no, user_id)

	if err != nil {
		config.Stdlog.Println(user_id, " - NOshot Proc Group NO Update - Select error : ( group_no : "+group_no+" / user_id : "+user_id+" ) : "+err.Error())
		config.Stdlog.Println(gudQuery)
		return err
	}

	return nil
}

func nresProcess(ctx context.Context, group_no, user_id string, pc int) {
	defer func(){
		if r := recover(); r != nil {
			config.Stdlog.Println(user_id, " - NOshot Proc resProcess panic error : ", r)
			if err, ok := r.(error); ok {
				if s.Contains(err.Error(), "connection refused") {
					for {
						config.Stdlog.Println(user_id, " - NOshot Proc resProcess send ping to DB")
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

	var db = databasepool.DB
	var stdlog = config.Stdlog

	var msgid, code, message, message_type, msg_sms, phn, remark1, remark2, result, sms_lms_tit, sms_kind, sms_sender, res_dt, reserve_dt, mms_file1, mms_file2, mms_file3, userid, sms_len_check, oshot sql.NullString
	var msgLen sql.NullInt64
	var phnstr string

	osmmsStrs := []string{}
	osmmsValues := []interface{}{}

	var resquery = `
	SELECT
		msgid, 
		code, 
		message, 
		message_type, 
		(case when sms_kind = 'S' then 
			substr(convert(REMOVE_WS(msg_sms) using euckr),1,100)
		 else 
		   convert(REMOVE_WS(msg_sms) using euckr)
	     end) as msg_sms, 
		phn, 
		remark1, 
		remark2,
		result, 
		convert(REMOVE_WS(sms_lms_tit) using euckr) as sms_lms_tit,
		sms_kind,
		sms_sender, 
		res_dt, 
		reserve_dt, 
		(select file1_path from api_mms_images aa where aa.user_id = drr.userid and aa.mms_id = drr.mms_image_id) as mms_file1, 
		(select file2_path from api_mms_images aa where aa.user_id = drr.userid and aa.mms_id = drr.mms_image_id) as mms_file2, 
		(select file3_path from api_mms_images aa where aa.user_id = drr.userid and aa.mms_id = drr.mms_image_id) as mms_file3
		,(case when sms_kind = 'S' then length(convert(REMOVE_WS(msg_sms) using euckr)) else 100 end) as msg_len
		,userid
		,(select max(sms_len_check) from DHN_CLIENT_LIST dcl where dcl.user_id = drr.userid) as sms_len_check
		,(select ifnull(max(oshot), 'OShot') from DHN_CLIENT_LIST dcl where dcl.user_id = drr.userid) as oshot  
	FROM DHN_RESULT drr 
	WHERE send_group = ?
	  and result = 'P'
      and userid = ?
	order by userid
	`
	resrows, err := db.QueryContext(ctx, resquery, group_no, user_id)

	if err != nil {
		stdlog.Println(user_id, " - NOshot Proc Result Table select error : ", err)
		stdlog.Println(resquery)
	}
	defer resrows.Close()

	smscnt := 0
	lmscnt := 0
	tcnt := 0
	reg, err := regexp.Compile("[^0-9]+")
	for resrows.Next() {
		resrows.Scan(&msgid, &code, &message, &message_type, &msg_sms, &phn, &remark1, &remark2, &result, &sms_lms_tit, &sms_kind, &sms_sender, &res_dt, &reserve_dt, &mms_file1, &mms_file2, &mms_file3, &msgLen, &userid, &sms_len_check, &oshot)

		phnstr = phn.String

		tcnt++
		if len(osmmsStrs) > 500 {
			stmt := fmt.Sprintf("insert into OShotMSG(MsgGroupID,SendType,Sender,Receiver,Subject,Msg,ReserveDT,TimeoutDT,SendResult,File_Path1,File_Path2,File_Path3,etc1,etc2) values %s", s.Join(osmmsStrs, ","))
			_, err := db.Exec(stmt, osmmsValues...)

			if err != nil {
				for i := 0; i < len(osmmsValues); i = i + 13 {
					eQuery := fmt.Sprintf("insert into OShotMSG(MsgGroupID,SendType,Sender,Receiver,Subject,Msg,ReserveDT,TimeoutDT,SendResult,File_Path1,File_Path2,File_Path3,etc1,etc2) "+
						"values('%v','%v','%v','%v','%v','%v',null,null,'%v','%v','%v','%v','%v','%v')", osmmsValues[i], osmmsValues[i+1], osmmsValues[i+2], osmmsValues[i+3], osmmsValues[i+4], osmmsValues[i+6], osmmsValues[i+7], osmmsValues[i+8], osmmsValues[i+9], osmmsValues[i+11], osmmsValues[i+12])
					_, err := db.Exec(eQuery)
					if err != nil {
						msgKey := fmt.Sprintf("%v", osmmsValues[i+11])
						useridt := fmt.Sprintf("%v", osmmsValues[i+12])
						stdlog.Println(user_id, " - NOshot Proc MSG Table Insert error : "+err.Error(), " - DHN Msg Key : ", msgKey)
						errcodemsg := err.Error()
						if s.Index(errcodemsg, "1366") > 0 {
							db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.code='7069', dr.message = concat(dr.message, ',부적절한 문자사용'),dr.remark2 = date_format(now(), '%Y-%m-%d %H:%i:%S') where userid = '" + useridt + "' and msgid = '" + msgKey + "'")
						}
					}
				}
			} else {
				stdlog.Println(user_id, " - NOshot Proc MSG Table Insert 처리 : ", len(osmmsStrs))
			}
			osmmsStrs = nil
			osmmsValues = nil
		}

		// 알림톡 발송 성공 혹은 문자 발송이 아니면
		// API_RESULT 성공 처리 함.
		if len(msg_sms.String) > 0 && len(sms_sender.String) > 0 { // msg_sms 가 와 sms_sender 에 값이 있으면 Oshot 발송 함.
			phnstr = reg.ReplaceAllString(phnstr, "")
			if s.HasPrefix(phnstr, "82") {
				phnstr = "0" + phnstr[2:len(phnstr)]
			}

			osmmsStrs = append(osmmsStrs, "(?,?,?,?,?,?,?,null,?,?,?,?,?,?)")
			osmmsValues = append(osmmsValues, group_no+"-"+strconv.Itoa(lmscnt))
			if sms_kind.String == "S" {
				osmmsValues = append(osmmsValues, "SMS")
				smscnt++
			} else {
				osmmsValues = append(osmmsValues, "MMS")
				lmscnt++
			}
			osmmsValues = append(osmmsValues, sms_sender.String)
			osmmsValues = append(osmmsValues, phnstr)
			osmmsValues = append(osmmsValues, sms_lms_tit.String)
			osmmsValues = append(osmmsValues, msg_sms.String)
			if s.EqualFold(reserve_dt.String, "00000000000000") {
				osmmsValues = append(osmmsValues, sql.NullString{})
			} else {
				osmmsValues = append(osmmsValues, sql.NullString{})
			}
			osmmsValues = append(osmmsValues, "0")
			osmmsValues = append(osmmsValues, mms_file1.String)
			osmmsValues = append(osmmsValues, mms_file2.String)
			osmmsValues = append(osmmsValues, mms_file3.String)
			osmmsValues = append(osmmsValues, msgid.String)
			osmmsValues = append(osmmsValues, userid.String)

		} else {
			db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.code='7011', dr.message = concat(dr.message, ',문자 발송 정보 누락'),dr.remark2 = date_format(now(), '%Y-%m-%d %H:%i:%S') where userid = '" + userid.String + "' and msgid = '" + msgid.String + "'")
		}

	}

	if len(osmmsStrs) > 0 {
		stmt := fmt.Sprintf("insert into OShotMSG(MsgGroupID,SendType,Sender,Receiver,Subject,Msg,ReserveDT,TimeoutDT,SendResult,File_Path1,File_Path2,File_Path3,etc1,etc2 ) values %s", s.Join(osmmsStrs, ","))
		_, err := db.Exec(stmt, osmmsValues...)

		if err != nil {
			for i := 0; i < len(osmmsValues); i = i + 13 {
				eQuery := fmt.Sprintf("insert into OShotMSG(MsgGroupID,SendType,Sender,Receiver,Subject,Msg,ReserveDT,TimeoutDT,SendResult,File_Path1,File_Path2,File_Path3,etc1,etc2 ) "+
					"values('%v','%v','%v','%v','%v',null,null,'%v','%v','%v','%v',null,'%v','%v')", osmmsValues[i], osmmsValues[i+1], osmmsValues[i+2], osmmsValues[i+3], osmmsValues[i+4], osmmsValues[i+6], osmmsValues[i+7], osmmsValues[i+8], osmmsValues[i+9], osmmsValues[i+11], osmmsValues[i+12])
				_, err := db.Exec(eQuery)
				if err != nil {
					msgKey := fmt.Sprintf("%v", osmmsValues[i+11])
					useridt := fmt.Sprintf("%v", osmmsValues[i+12])
					stdlog.Println(user_id, " - NOshot Proc MSG Table Insert error : "+err.Error(), " - DHN Msg Key : ", msgKey)
					errcodemsg := err.Error()
					if s.Index(errcodemsg, "1366") > 0 {
						db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.code='7069', dr.message = concat(dr.message, ',부적절한 문자사용'),dr.remark2 = date_format(now(), '%Y-%m-%d %H:%i:%S') where userid = '" + useridt + "' and msgid = '" + msgKey + "'")
					}
				}
			}
		} else {
			stdlog.Println(user_id, " - NOshot Proc MSG Table Insert 처리 : ", len(osmmsStrs))
		}

	}

	if smscnt > 0 || lmscnt > 0 {
		stdlog.Println(user_id, " - NOshot 발송 처리 완료 ( ", group_no, " ) : SMS - ", smscnt, " , LMS - ", lmscnt, ", 총 - ", tcnt, " : ( Proc Cnt :", pc, ") - END")
	}
}
