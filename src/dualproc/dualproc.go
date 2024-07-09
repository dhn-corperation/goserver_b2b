package dualproc

import (
	"database/sql"
	"fmt"
	"strconv"

	config "mycs/src/kaoconfig"
	databasepool "mycs/src/kaodatabasepool"

	// "encoding/hex"
	"regexp"
	s "strings"
	"time"
	// "unicode/utf8"

	"context"
)

var procCnt int
var Rate1 int
var Rate2 int
var applyRate1 int
var applyRate2 int
var seq int

func DualProcess(user_id string, rate string, ctx context.Context) {
	//var wg sync.WaitGroup
	config.Stdlog.Println(user_id, " Lgu Process 시작 됨.")
	procCnt = 0

	parts := s.Split(rate, "/")
	Rate1, err1 := strconv.Atoi(parts[0])
	Rate2, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		Rate1 = 1
		Rate2 = 1
	}
	applyRate1 = Rate1
	applyRate2 = Rate2
	seq = 1

	for {

		if procCnt < 5 {

			select {
			case <-ctx.Done():

				uid := ctx.Value("user_id")
				config.Stdlog.Println(uid, " - Dual process가 10초 후에 종료 됨.")
				time.Sleep(10 * time.Second)
				config.Stdlog.Println(uid, " - Dual process 종료 완료")
				return
			default:

				var count sql.NullInt64
				tickSql := `
				select
					length(msgid) as cnt
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
					config.Stdlog.Println("DHN_RESULT Table - select 오류 : " + cnterr.Error())
				} else {
					if count.Int64 > 0 {
						var startNow = time.Now()
						var group_no = fmt.Sprintf("%02d%02d%02d%02d%06d", startNow.Day(), startNow.Hour(), startNow.Minute(), startNow.Second(), (startNow.Nanosecond() / 1000))

						upError := updateReqeust(ctx, group_no, user_id)
						if upError != nil {
							config.Stdlog.Println(user_id, "Group No Update 오류", group_no)
						} else {
							go resProcess(ctx, group_no, user_id, rate)
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

	config.Stdlog.Println(user_id, "- Dual Group No Update 시작", group_no)

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
		config.Stdlog.Println(user_id, "- Group NO Update - Select error : ( group_no : "+group_no+" / user_id : "+user_id+" ) : "+err.Error())
		config.Stdlog.Println(gudQuery)
		return err
	}

	return nil
}

func resProcess(ctx context.Context, group_no string, user_id string, rate string) {
	//defer wg.Done()
	procCnt++
	var db = databasepool.DB
	var stdlog = config.Stdlog

	defer func() {
		if err := recover(); err != nil {
			procCnt--
			stdlog.Println(user_id, "- ", group_no, "Dual 처리 중 오류 발생 : ", err)
		}
	}()

	var msgid, code, message, message_type, msg_sms, phn, remark1, remark2, result, sms_lms_tit, sms_kind, sms_sender, res_dt, reserve_dt, mms_file1, mms_file2, mms_file3, userid, sms_len_check sql.NullString
	var msgLen sql.NullInt64
	var phnstr string

	lgsmsStrs := []string{}
	lgsmsValues := []interface{}{}

	lgmmsStrs := []string{}
	lgmmsValues := []interface{}{}

	nnsmsStrs := []string{}
	nnsmsValues := []interface{}{}

	nnmmsStrs := []string{}
	nnmmsValues := []interface{}{}

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
		(select ifnull(file1_path, '') from api_mms_images aa where aa.user_id = drr.userid and aa.mms_id = drr.p_invoice) as mms_file1, 
		(select ifnull(file2_path, '') from api_mms_images aa where aa.user_id = drr.userid and aa.mms_id = drr.p_invoice) as mms_file2, 
		(select ifnull(file3_path, '') from api_mms_images aa where aa.user_id = drr.userid and aa.mms_id = drr.p_invoice) as mms_file3
		,(case when sms_kind = 'S' then length(convert(REMOVE_WS(msg_sms) using euckr)) else 100 end) as msg_len
		,userid
		,(select max(sms_len_check) from DHN_CLIENT_LIST dcl where dcl.user_id = drr.userid) as sms_len_check
	FROM DHN_RESULT drr 
	WHERE send_group = ?
	  and result = 'P'
      and userid = ?
	order by userid
	`
	resrows, err := db.QueryContext(ctx, resquery, group_no, user_id)

	if err != nil {
		stdlog.Println("Result Table 조회 중 오류 발생")
		stdlog.Println(err)
		stdlog.Println(resquery)
	}
	defer resrows.Close()

	smscnt := 0
	lmscnt := 0
	tcnt := 0
	reg, err := regexp.Compile("[^0-9]+")
	for resrows.Next() {
		resrows.Scan(&msgid, &code, &message, &message_type, &msg_sms, &phn, &remark1, &remark2, &result, &sms_lms_tit, &sms_kind, &sms_sender, &res_dt, &reserve_dt, &mms_file1, &mms_file2, &mms_file3, &msgLen, &userid, &sms_len_check)

		phnstr = phn.String

		if tcnt == 0 {
			stdlog.Println(user_id, "-", group_no, "문자발송 처리 시작 : ", " Process cnt : ", procCnt)
		}

		tcnt++

		if len(lgsmsStrs) > 500 {
			stmt := fmt.Sprintf("insert into LG_SC_TRAN(TR_SENDDATE,TR_PHONE,TR_CALLBACK, TR_MSG, TR_ETC1, TR_ETC2, TR_ETC2, TR_KISAORIGCODE) values %s", s.Join(lgsmsStrs, ","))
			stdlog.Println(stmt)
			_, err := db.ExecContext(ctx, stmt, lgsmsValues...)

			if err != nil {
				for i := 0; i < len(lgsmsValues); i = i + 7 {
					eQuery := fmt.Sprintf("insert into LG_SC_TRAN(TR_SENDDATE,TR_PHONE,TR_CALLBACK, TR_MSG, TR_ETC1, TR_ETC2, TR_ETC3, TR_KISAORIGCODE) "+
						"values('%v','%v','%v','%v','%v','%v','%v', '%v')", lgsmsValues[i], lgsmsValues[i+1], lgsmsValues[i+2], lgsmsValues[i+3], lgsmsValues[i+4], lgsmsValues[i+5], lgsmsValues[i+6], lgsmsValues[i+7])
					_, err := db.ExecContext(ctx, eQuery)
					if err != nil {
						msgKey := fmt.Sprintf("%v", lgsmsValues[i+4])
						useridt := fmt.Sprintf("%v", lgsmsValues[i+5])
						stdlog.Println(user_id, "- Lgu SMS Table Insert 처리 중 오류 발생 : "+err.Error(), " - DHN Msg Key : ", msgKey)
						errcodemsg := err.Error()
						if s.Index(errcodemsg, "1366") > 0 {
							db.ExecContext(ctx, "update DHN_RESULT dr set dr.result = 'Y', dr.code='7069', dr.message = concat(dr.message, ',부적절한 문자사용'),dr.remark2 = date_format(now(), '%Y-%m-%d %H:%i:%S') where userid = ? and  msgid = ?", useridt, msgKey)
						}
					}
				}
			} else {
				stdlog.Println(user_id, "- Lgu SMS Table Insert 처리 : ", len(lgsmsStrs), " - LG_SC_TRAN")
			}
			lgsmsStrs = nil
			lgsmsValues = nil
		}

		if len(lgmmsStrs) > 500 {
			stmt := fmt.Sprintf("insert into LG_MMS_MSG(SUBJECT, PHONE, CALLBACK, REQDATE, MSG, FILE_CNT, FILE_PATH1, FILE_PATH2, FILE_PATH3, ETC1, ETC2, ETC3, KISA_ORIGCODE) values %s", s.Join(lgmmsStrs, ","))
			stdlog.Println(stmt)
			_, err := db.Exec(stmt, lgmmsValues...)

			if err != nil {
				for i := 0; i < len(lgmmsValues); i = i + 12 {
					eQuery := fmt.Sprintf("SUBJECT, PHONE, CALLBACK, REQDATE, MSG, FILE_CNT, FILE_PATH1, FILE_PATH2, FILE_PATH3, ETC1, ETC2, ETC3, KISA_ORIGCODE) "+
						"values('%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v')", lgmmsValues[i], lgmmsValues[i+1], lgmmsValues[i+2], lgmmsValues[i+3], lgmmsValues[i+4], lgmmsValues[i+5], lgmmsValues[i+6], lgmmsValues[i+7], lgmmsValues[i+8], lgmmsValues[i+9], lgmmsValues[i+10], lgmmsValues[i+11], lgmmsValues[i+12])
					_, err := db.Exec(eQuery)
					if err != nil {
						msgKey := fmt.Sprintf("%v", lgmmsValues[i+9])
						useridt := fmt.Sprintf("%v", lgmmsValues[i+10])
						stdlog.Println(user_id, "- Lgu LMS Table Insert 처리 중 오류 발생 : "+err.Error(), " - DHN Msg Key : ", msgKey)
						errcodemsg := err.Error()
						if s.Index(errcodemsg, "1366") > 0 {
							db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.code='7069', dr.message = concat(dr.message, ',부적절한 문자사용'),dr.remark2 = date_format(now(), '%Y-%m-%d %H:%i:%S') where userid = '" + useridt + "' and msgid = '" + msgKey + "'")
						}
					}
				}
			} else {
				stdlog.Println(user_id, "- Lgu MMS Table Insert 처리 : ", len(lgmmsStrs), " - LG_MMS_MSG")
			}
			lgmmsStrs = nil
			lgmmsValues = nil
		}

		if len(nnsmsStrs) > 500 {
			stmt := fmt.Sprintf("insert into SMS_MSG(TR_CALLBACK,TR_PHONE,TR_MSG,TR_SENDDATE,TR_SENDSTAT,TR_MSGTYPE,TR_ETC9,TR_ETC10,TR_IDENTIFICATION_CODE,TR_ETC8) values %s", s.Join(nnsmsStrs, ","))
			stdlog.Println(stmt)
			_, err := db.Exec(stmt, nnsmsValues...)

			if err != nil {
				//stdlog.Println("Nano SMS Table Insert 처리 중 오류 발생 " + err.Error())
				for i := 0; i < len(nnsmsValues); i = i + 8 {
					eQuery := fmt.Sprintf("insert into SMS_MSG(TR_CALLBACK,TR_PHONE,TR_MSG,TR_SENDDATE,TR_SENDSTAT,TR_MSGTYPE,TR_ETC9,TR_ETC10,TR_IDENTIFICATION_CODE,TR_ETC8) "+
						"values('%v','%v','%v','%v','%v','%v','%v', '%v', 'Y')", nnsmsValues[i], nnsmsValues[i+1], nnsmsValues[i+2], nnsmsValues[i+3], nnsmsValues[i+4], nnsmsValues[i+5], nnsmsValues[i+6], nnsmsValues[i+7], nnsmsValues[i+8])
					_, err := db.Exec(eQuery)
					if err != nil {
						msgKey := fmt.Sprintf("%v", nnsmsValues[i+6])
						useridt := fmt.Sprintf("%v", nnsmsValues[i+7])
						stdlog.Println(user_id, "- Nano SMS Table Insert 처리 중 오류 발생 : "+err.Error(), " - DHN Msg Key : ", msgKey)
						errcodemsg := err.Error()
						if s.Index(errcodemsg, "1366") > 0 {
							db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.code='7069', dr.message = concat(dr.message, ',부적절한 문자사용'),dr.remark2 = date_format(now(), '%Y-%m-%d %H:%i:%S') where userid = '" + useridt + "' and  msgid = '" + msgKey + "'")
						}
					}
				}
				//db.Exec("update API_RESULT ar set ar.msg_type = '" + sms_kind.String + "', result_code = '9999', error_text = '기타오류', report_time = date_format(now(), '%Y-%m-%d %H:%i:%S') where dhn_msg_id = '" + msgid.String + "'")
			} else {
				stdlog.Println(user_id, "- Nano SMS Table Insert 처리 : ", len(nnsmsStrs))
			}
			nnsmsStrs = nil
			nnsmsValues = nil
		}

		if len(nnmmsStrs) > 500 {
			stmt := fmt.Sprintf("insert into MMS_MSG(CALLBACK,PHONE,SUBJECT,MSG,REQDATE,STATUS,FILE_CNT,FILE_PATH1,FILE_PATH2,FILE_PATH3,ETC9,ETC10,IDENTIFICATION_CODE,ETC8) values %s", s.Join(nnmmsStrs, ","))
			stdlog.Println(stmt)
			_, err := db.Exec(stmt, nnmmsValues...)

			if err != nil {
				//stdlog.Println("Nano SMS Table Insert 처리 중 오류 발생 " + err.Error())
				for i := 0; i < len(nnmmsValues); i = i + 12 {
					eQuery := fmt.Sprintf("insert into MMS_MSG(CALLBACK,PHONE,SUBJECT,MSG,REQDATE,STATUS,FILE_CNT,FILE_PATH1,FILE_PATH2,FILE_PATH3,ETC9,ETC10,IDENTIFICATION_CODE,ETC8) "+
						"values('%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','Y')", nnmmsValues[i], nnmmsValues[i+1], nnmmsValues[i+2], nnmmsValues[i+3], nnmmsValues[i+4], nnmmsValues[i+5], nnmmsValues[i+6], nnmmsValues[i+7], nnmmsValues[i+8], nnmmsValues[i+9], nnmmsValues[i+10], nnmmsValues[i+11], nnmmsValues[i+12])
					_, err := db.Exec(eQuery)
					if err != nil {
						msgKey := fmt.Sprintf("%v", nnmmsValues[i+10])
						useridt := fmt.Sprintf("%v", nnmmsValues[i+11])
						stdlog.Println(user_id, "- Nano LMS Table Insert 처리 중 오류 발생 : "+err.Error(), " - DHN Msg Key : ", msgKey)
						errcodemsg := err.Error()
						if s.Index(errcodemsg, "1366") > 0 {
							db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.code='7069', dr.message = concat(dr.message, ',부적절한 문자사용'),dr.remark2 = date_format(now(), '%Y-%m-%d %H:%i:%S') where userid = '" + useridt + "' and msgid = '" + msgKey + "'")
						}
					}
				}
				//db.Exec("update API_RESULT ar set ar.msg_type = '" + sms_kind.String + "', result_code = '9999', error_text = '기타오류', report_time = date_format(now(), '%Y-%m-%d %H:%i:%S') where dhn_msg_id = '" + msgid.String + "'")
			} else {
				stdlog.Println(user_id, "- Nano MMS Table Insert 처리 : ", len(nnmmsStrs))
			}
			nnmmsStrs = nil
			nnmmsValues = nil
		}

		// 알림톡 발송 성공 혹은 문자 발송이 아니면
		// API_RESULT 성공 처리 함.
		if len(msg_sms.String) > 0 && len(sms_sender.String) > 0 { // msg_sms 가 와 sms_sender 에 값이 있으면 LG 발송 함.
			phnstr = reg.ReplaceAllString(phnstr, "")
			if s.HasPrefix(phnstr, "82") {
				phnstr = "0" + phnstr[2:len(phnstr)]
			}

			if seq == 1 {
				stdlog.Println("여기오냐1")
				if s.EqualFold(sms_kind.String, "S") {
					if msgLen.Int64 <= 90 || s.EqualFold(sms_len_check.String, "N") {
						lgsmsStrs = append(lgsmsStrs, "(?,?,?,?,?,?,?,?)")
						lgsmsValues = append(lgsmsValues, time.Now().Format("2006-01-02 15:04:05"))
						lgsmsValues = append(lgsmsValues, phnstr)
						lgsmsValues = append(lgsmsValues, sms_sender.String)
						lgsmsValues = append(lgsmsValues, msg_sms.String)
						lgsmsValues = append(lgsmsValues, msgid.String)
						lgsmsValues = append(lgsmsValues, userid.String)
						lgsmsValues = append(lgsmsValues, group_no)
						lgsmsValues = append(lgsmsValues, "302190001")
						smscnt++
					} else {
						db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.code = '7003', dr.message = '메세지 길이 오류', dr.remark2 = date_format(now(), '%Y-%m-%d %H:%i:%S') where userid = '" + userid.String + "' and msgid = '" + msgid.String + "'")
					}
				} else if s.EqualFold(sms_kind.String, "L") || s.EqualFold(sms_kind.String, "M") {
					file_cnt  := 0
					if mms_file1.String != "" {
						file_cnt++
					}
					if mms_file2.String != "" {
						file_cnt++
					}
					if mms_file3.String != "" {
						file_cnt++
					}
					lgmmsStrs = append(lgmmsStrs, "(?,?,?,?,?,?,?,?,?,?,?,?,?)")
					lgmmsValues = append(lgmmsValues, sms_lms_tit.String)
					lgmmsValues = append(lgmmsValues, phnstr)
					lgmmsValues = append(lgmmsValues, sms_sender.String)
					lgmmsValues = append(lgmmsValues, time.Now().Format("2006-01-02 15:04:05"))
					lgmmsValues = append(lgmmsValues, msg_sms.String)
					lgmmsValues = append(lgmmsValues, file_cnt)
					lgmmsValues = append(lgmmsValues, mms_file1.String)
					lgmmsValues = append(lgmmsValues, mms_file2.String)
					lgmmsValues = append(lgmmsValues, mms_file3.String)
					lgmmsValues = append(lgmmsValues, msgid.String)
					lgmmsValues = append(lgmmsValues, userid.String)
					lgmmsValues = append(lgmmsValues, group_no)
					lgmmsValues = append(lgmmsValues, "302190001")
					lmscnt++
				}
				applyRate1--
				if applyRate1 <= 0 {
					applyRate1 = Rate1
					seq = 2
				}
			} else if seq == 2 {
				stdlog.Println("여기오냐2")
				if s.EqualFold(sms_kind.String, "S") {

					if msgLen.Int64 <= 90 || s.EqualFold(sms_len_check.String, "N") {

						nnsmsStrs = append(nnsmsStrs, "(?,?,?,?,?,?,?,?,?,'Y')")
						nnsmsValues = append(nnsmsValues, sms_sender.String)
						nnsmsValues = append(nnsmsValues, phnstr)
						nnsmsValues = append(nnsmsValues, msg_sms.String)
						nnsmsValues = append(nnsmsValues, reserve_dt.String)
						nnsmsValues = append(nnsmsValues, "0")
						nnsmsValues = append(nnsmsValues, "0")
						nnsmsValues = append(nnsmsValues, msgid.String)
						nnsmsValues = append(nnsmsValues, userid.String)
						nnsmsValues = append(nnsmsValues, config.Conf.NANO_IDENTI_CODE)
						smscnt++
					} else {
						db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.code = '7003', dr.message = '메세지 길이 오류', dr.remark2 = date_format(now(), '%Y-%m-%d %H:%i:%S') where userid = '" + userid.String + "' and msgid = '" + msgid.String + "'")
					}
				} else if s.EqualFold(sms_kind.String, "L") || s.EqualFold(sms_kind.String, "M") {

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

					nnmmsStrs = append(nnmmsStrs, "( ?,?,?,?,?,?,?,?,?,?,?,?,?,'Y')")

					nnmmsValues = append(nnmmsValues, sms_sender.String)
					nnmmsValues = append(nnmmsValues, phnstr)
					nnmmsValues = append(nnmmsValues, sms_lms_tit.String)
					nnmmsValues = append(nnmmsValues, msg_sms.String)
					nnmmsValues = append(nnmmsValues, reserve_dt.String)
					nnmmsValues = append(nnmmsValues, "0")
					nnmmsValues = append(nnmmsValues, filecnt)
					nnmmsValues = append(nnmmsValues, mms_file1.String)
					nnmmsValues = append(nnmmsValues, mms_file2.String)
					nnmmsValues = append(nnmmsValues, mms_file3.String)

					nnmmsValues = append(nnmmsValues, msgid.String)
					nnmmsValues = append(nnmmsValues, userid.String)
					nnmmsValues = append(nnmmsValues, config.Conf.NANO_IDENTI_CODE)
					lmscnt++
				}
				applyRate2--
				if applyRate2 <= 0 {
					applyRate2 = Rate2
					seq = 1
				}
			}
		} else {
			db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.code='7011', dr.message = concat(dr.message, ',문자 발송 정보 누락'),dr.remark2 = date_format(now(), '%Y-%m-%d %H:%i:%S') where userid = '" + userid.String + "' and msgid = '" + msgid.String + "'")
		}

	}

	if len(lgsmsStrs) > 0 {
		stmt := fmt.Sprintf("insert into LG_SC_TRAN(TR_SENDDATE,TR_PHONE,TR_CALLBACK, TR_MSG, TR_ETC1, TR_ETC2, TR_ETC2, TR_KISAORIGCODE) values %s", s.Join(lgsmsStrs, ","))
		stdlog.Println(stmt)
		_, err := db.ExecContext(ctx, stmt, lgsmsValues...)

		if err != nil {
			for i := 0; i < len(lgsmsValues); i = i + 7 {
				eQuery := fmt.Sprintf("insert into LG_SC_TRAN(TR_SENDDATE,TR_PHONE,TR_CALLBACK, TR_MSG, TR_ETC1, TR_ETC2, TR_ETC3, TR_KISAORIGCODE) "+
					"values('%v','%v','%v','%v','%v','%v','%v','%v')", lgsmsValues[i], lgsmsValues[i+1], lgsmsValues[i+2], lgsmsValues[i+3], lgsmsValues[i+4], lgsmsValues[i+5], lgsmsValues[i+6], lgsmsValues[i+7])
				_, err := db.ExecContext(ctx, eQuery)
				if err != nil {
					msgKey := fmt.Sprintf("%v", lgsmsValues[i+4])
					useridt := fmt.Sprintf("%v", lgsmsValues[i+5])
					stdlog.Println(user_id, "- Lgu SMS Table Insert 처리 중 오류 발생 : "+err.Error(), " - DHN Msg Key : ", msgKey)
					errcodemsg := err.Error()
					if s.Index(errcodemsg, "1366") > 0 {
						db.ExecContext(ctx, "update DHN_RESULT dr set dr.result = 'Y', dr.code='7069', dr.message = concat(dr.message, ',부적절한 문자사용'),dr.remark2 = date_format(now(), '%Y-%m-%d %H:%i:%S') where userid = ? and  msgid = ?", useridt, msgKey)
					}
				}
			}
		} else {
			stdlog.Println(user_id, "- Lgu SMS Table Insert 처리 : ", len(lgsmsStrs), " - LG_SC_TRAN")
		}

	}

	if len(lgmmsStrs) > 0 {
		stdlog.Println("여기오냐3")
		stmt := fmt.Sprintf("insert into LG_MMS_MSG(SUBJECT, PHONE, CALLBACK, REQDATE, MSG, FILE_CNT, FILE_PATH1, FILE_PATH2, FILE_PATH3, ETC1, ETC2, ETC3, KISA_ORIGCODE) values %s", s.Join(lgmmsStrs, ","))
		stdlog.Println(stmt)
		_, err := db.Exec(stmt, lgmmsValues...)

		if err != nil {
			for i := 0; i < len(lgmmsValues); i = i + 12 {
				eQuery := fmt.Sprintf("SUBJECT, PHONE, CALLBACK, REQDATE, MSG, FILE_CNT, FILE_PATH1, FILE_PATH2, FILE_PATH3, ETC1, ETC2, ETC3, KISA_ORIGCODE) "+
					"values('%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v')", lgmmsValues[i], lgmmsValues[i+1], lgmmsValues[i+2], lgmmsValues[i+3], lgmmsValues[i+4], lgmmsValues[i+5], lgmmsValues[i+6], lgmmsValues[i+7], lgmmsValues[i+8], lgmmsValues[i+9], lgmmsValues[i+10], lgmmsValues[i+11], lgmmsValues[i+12])
				_, err := db.Exec(eQuery)
				if err != nil {
					msgKey := fmt.Sprintf("%v", lgmmsValues[i+9])
					useridt := fmt.Sprintf("%v", lgmmsValues[i+10])
					stdlog.Println(user_id, "- Lgu LMS Table Insert 처리 중 오류 발생 : "+err.Error(), " - DHN Msg Key : ", msgKey)
					errcodemsg := err.Error()
					if s.Index(errcodemsg, "1366") > 0 {
						db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.code='7069', dr.message = concat(dr.message, ',부적절한 문자사용'),dr.remark2 = date_format(now(), '%Y-%m-%d %H:%i:%S') where userid = '" + useridt + "' and msgid = '" + msgKey + "'")
					}
				}
			}
		} else {
			stdlog.Println(user_id, "- Lgu MMS Table Insert 처리 : ", len(lgmmsStrs), " - LG_MMS_MSG")
		}
	}

	if len(nnsmsStrs) > 0 {
		stmt := fmt.Sprintf("insert into SMS_MSG(TR_CALLBACK,TR_PHONE,TR_MSG,TR_SENDDATE,TR_SENDSTAT,TR_MSGTYPE,TR_ETC9,TR_ETC10,TR_IDENTIFICATION_CODE,TR_ETC8) values %s", s.Join(nnsmsStrs, ","))
		stdlog.Println(stmt)
		_, err := db.Exec(stmt, nnsmsValues...)

		if err != nil {
			//stdlog.Println("Nano SMS Table Insert 처리 중 오류 발생 " + err.Error())
			for i := 0; i < len(nnsmsValues); i = i + 8 {
				eQuery := fmt.Sprintf("insert into SMS_MSG(TR_CALLBACK,TR_PHONE,TR_MSG,TR_SENDDATE,TR_SENDSTAT,TR_MSGTYPE,TR_ETC9,TR_ETC10,TR_IDENTIFICATION_CODE,TR_ETC8) "+
					"values('%v','%v','%v','%v','%v', '%v','%v', '%v', 'Y')", nnsmsValues[i], nnsmsValues[i+1], nnsmsValues[i+2], nnsmsValues[i+3], nnsmsValues[i+4], nnsmsValues[i+5], nnsmsValues[i+6], nnsmsValues[i+7], nnsmsValues[i+8])
				_, err := db.Exec(eQuery)
				if err != nil {
					msgKey := fmt.Sprintf("%v", nnsmsValues[i+6])
					useridt := fmt.Sprintf("%v", nnsmsValues[i+7])
					stdlog.Println(user_id, "- Nano SMS Table Insert 처리 중 오류 발생 : "+err.Error(), " - DHN Msg Key : ", msgKey)
					errcodemsg := err.Error()
					if s.Index(errcodemsg, "1366") > 0 {
						db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.code='7069', dr.message = concat(dr.message, ',부적절한 문자사용'),dr.remark2 = date_format(now(), '%Y-%m-%d %H:%i:%S') where userid = '" + useridt + "' and  msgid = '" + msgKey + "'")
					}
				}
			}
			//db.Exec("update API_RESULT ar set ar.msg_type = '" + sms_kind.String + "', result_code = '9999', error_text = '기타오류', report_time = date_format(now(), '%Y-%m-%d %H:%i:%S') where dhn_msg_id = '" + msgid.String + "'")
		} else {
			stdlog.Println(user_id, "- Nano SMS Table Insert 처리 : ", len(nnsmsStrs))
		}

	}

	if len(nnmmsStrs) > 0 {
		stdlog.Println("여기오냐4")
		stmt := fmt.Sprintf("insert into MMS_MSG(CALLBACK,PHONE,SUBJECT,MSG,REQDATE,STATUS,FILE_CNT,FILE_PATH1,FILE_PATH2,FILE_PATH3,ETC9,ETC10,IDENTIFICATION_CODE,ETC8) values %s", s.Join(nnmmsStrs, ","))
		stdlog.Println(stmt)
		_, err := db.Exec(stmt, nnmmsValues...)

		if err != nil {
			//stdlog.Println("Nano SMS Table Insert 처리 중 오류 발생 " + err.Error())
			for i := 0; i < len(nnmmsValues); i = i + 12 {
				eQuery := fmt.Sprintf("insert into MMS_MSG(CALLBACK,PHONE,SUBJECT,MSG,REQDATE,STATUS,FILE_CNT,FILE_PATH1,FILE_PATH2,FILE_PATH3,ETC9,ETC10,IDENTIFICATION_CODE,ETC8) "+
					"values('%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','Y')", nnmmsValues[i], nnmmsValues[i+1], nnmmsValues[i+2], nnmmsValues[i+3], nnmmsValues[i+4], nnmmsValues[i+5], nnmmsValues[i+6], nnmmsValues[i+7], nnmmsValues[i+8], nnmmsValues[i+9], nnmmsValues[i+10], nnmmsValues[i+11], nnmmsValues[i+12])
				_, err := db.Exec(eQuery)
				if err != nil {
					msgKey := fmt.Sprintf("%v", nnmmsValues[i+10])
					useridt := fmt.Sprintf("%v", nnmmsValues[i+11])
					stdlog.Println(user_id, "- Nano LMS Table Insert 처리 중 오류 발생 : "+err.Error(), " - DHN Msg Key : ", msgKey)
					errcodemsg := err.Error()
					if s.Index(errcodemsg, "1366") > 0 {
						db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.code='7069', dr.message = concat(dr.message, ',부적절한 문자사용'),dr.remark2 = date_format(now(), '%Y-%m-%d %H:%i:%S') where userid = '" + useridt + "' and msgid = '" + msgKey + "'")
					}
				}
			}
			//db.Exec("update API_RESULT ar set ar.msg_type = '" + sms_kind.String + "', result_code = '9999', error_text = '기타오류', report_time = date_format(now(), '%Y-%m-%d %H:%i:%S') where dhn_msg_id = '" + msgid.String + "'")
		} else {
			stdlog.Println(user_id, "- Nano MMS Table Insert 처리 : ", len(nnmmsStrs))
		}

	}

	if smscnt > 0 || lmscnt > 0 {
		stdlog.Println(user_id, "-", group_no, "문자 발송 처리 완료 ( ", tcnt, " ) : SMS -", smscnt, " , LMS -", lmscnt, " >> Process cnt : ", procCnt)
	}
	procCnt--
}