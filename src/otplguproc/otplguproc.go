package otplguproc

import (
	"database/sql"
	"fmt"
	//"strconv"

	//"sync"
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

func LguProcess(ctx context.Context) {
	//var wg sync.WaitGroup
	config.Stdlog.Println("Lgu OTP Process 시작 됨.")
	procCnt = 0
	for {

		if procCnt < 5 {

			select {
			case <-ctx.Done():

				config.Stdlog.Println("Lgu OTP process가 10초 후에 종료 됨.")
				time.Sleep(10 * time.Second)
				config.Stdlog.Println("Lgu OTP process 종료 완료")
				return
			default:

				var count sql.NullInt64
				tickSql := `
				select
					length(msgid) as cnt
				from
					DHN_RESULT dr
				where
					dr.result = 'O'
					and dr.send_group is null
					and ifnull(dr.reserve_dt, '00000000000000') <= date_format(now(), '%Y%m%d%H%i%S')
				limit 1
					`
				cnterr := databasepool.DB.QueryRowContext(ctx, tickSql).Scan(&count)

				if cnterr != nil && cnterr != sql.ErrNoRows {
					config.Stdlog.Println("DHN_RESULT Table - select 오류 : " + cnterr.Error())
					time.Sleep(10 * time.Second)
				} else {
					if count.Int64 > 0 {
						var startNow = time.Now()
						var group_no = fmt.Sprintf("%02d%02d%02d%02d%06d", startNow.Day(), startNow.Hour(), startNow.Minute(), startNow.Second(), (startNow.Nanosecond() / 1000))

						upError := updateReqeust(ctx, group_no)
						if upError != nil {
							config.Stdlog.Println("Group No Update 오류", group_no)
						} else {
							go resProcess(ctx, group_no)
						}
					}
				}
			}

		}
	}
}

func updateReqeust(ctx context.Context, group_no string) error {

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

	config.Stdlog.Println("Lgu OTP Group No Update 시작", group_no)

	gudQuery := `
	update DHN_RESULT dr
	set	send_group = ?
	where result = 'O'
	  and send_group is null
	  and ifnull(reserve_dt, '00000000000000') <= date_format(now(), '%Y%m%d%H%i%S')
	LIMIT 500
	`
	_, err = tx.ExecContext(ctx, gudQuery, group_no)

	if err != nil {
		config.Stdlog.Println("Lgu OTP Group NO Update - Select error : ( group_no : "+group_no+" ) : "+err.Error())
		config.Stdlog.Println(gudQuery)
		return err
	}

	return nil
}

func resProcess(ctx context.Context, group_no string) {
	defer func(){
		if r := recover(); r != nil {
			config.Stdlog.Println("OTPLGU mmsProcess panic 발생 원인 : ", r)
			procCnt--
			if err, ok := r.(error); ok {
				if s.Contains(err.Error(), "connection refused") {
					for {
						config.Stdlog.Println("OTPLGU mmsProcess send ping to DB")
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

	procCnt++
	var db = databasepool.DB
	var stdlog = config.Stdlog

	var msgid, code, message, message_type, msg_sms, phn, remark1, remark2, result, sms_lms_tit, sms_kind, sms_sender, res_dt, reserve_dt, mms_file1, mms_file2, mms_file3, userid, sms_len_check sql.NullString
	var msgLen sql.NullInt64
	var phnstr string
	ossmsStrs := []string{}
	ossmsValues := []interface{}{}

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
		(select ifnull(file1_path, '') from api_mms_images aa where aa.user_id = drr.userid and aa.mms_id = drr.mms_image_id) as mms_file1, 
		(select ifnull(file2_path, '') from api_mms_images aa where aa.user_id = drr.userid and aa.mms_id = drr.mms_image_id) as mms_file2, 
		(select ifnull(file3_path, '') from api_mms_images aa where aa.user_id = drr.userid and aa.mms_id = drr.mms_image_id) as mms_file3,
		(case when sms_kind = 'S' then length(convert(REMOVE_WS(msg_sms) using euckr)) else 100 end) as msg_len,
		userid,
		'Y' as sms_len_check
	FROM DHN_RESULT drr 
	WHERE send_group = ?
	  and result = 'O'
	order by drr.reg_dt
	`
	resrows, err := db.QueryContext(ctx, resquery, group_no)

	if err != nil {
		stdlog.Println("Lgu OTP Result Table 조회 중 오류 발생")
		stdlog.Println(err)
		stdlog.Println(resquery)
	}
	defer resrows.Close()

	scnt := 0
	fcnt := 0
	smscnt := 0
	lmscnt := 0
	tcnt := 0
	reg, err := regexp.Compile("[^0-9]+")
	for resrows.Next() {
		resrows.Scan(&msgid, &code, &message, &message_type, &msg_sms, &phn, &remark1, &remark2, &result, &sms_lms_tit, &sms_kind, &sms_sender, &res_dt, &reserve_dt, &mms_file1, &mms_file2, &mms_file3, &msgLen, &userid, &sms_len_check)

		phnstr = phn.String

		if tcnt == 0 {
			stdlog.Println(group_no, " Lgu OTP 문자발송 처리 시작 : ", "Process cnt : ", procCnt)
		}

		tcnt++
		if len(ossmsStrs) > 500 {
			stmt := fmt.Sprintf("insert into LG_OTP_SC_TRAN(TR_SENDDATE,TR_PHONE,TR_CALLBACK, TR_MSG, TR_ETC1, TR_ETC2, TR_ETC3, TR_KISAORIGCODE) values %s", s.Join(ossmsStrs, ","))
			_, err := db.ExecContext(ctx, stmt, ossmsValues...)

			if err != nil {
				for i := 0; i < len(ossmsValues); i = i + 7 {
					eQuery := fmt.Sprintf("insert into LG_OTP_SC_TRAN(TR_SENDDATE,TR_PHONE,TR_CALLBACK, TR_MSG, TR_ETC1, TR_ETC2, TR_ETC3, TR_KISAORIGCODE) "+
						"values('%v','%v','%v','%v','%v','%v','%v')", ossmsValues[i], ossmsValues[i+1], ossmsValues[i+2], ossmsValues[i+3], ossmsValues[i+4], ossmsValues[i+5], ossmsValues[i+6], ossmsValues[i+7])
					_, err := db.ExecContext(ctx, eQuery)
					if err != nil {
						msgKey := fmt.Sprintf("%v", ossmsValues[i+4])
						useridt := fmt.Sprintf("%v", ossmsValues[i+5])
						stdlog.Println("Lgu OTP SMS Table Insert 처리 중 오류 발생 : "+err.Error(), " - DHN Msg Key : ", msgKey)
						errcodemsg := err.Error()
						if s.Index(errcodemsg, "1366") > 0 {
							db.ExecContext(ctx, "update DHN_RESULT dr set dr.result = 'Y', dr.code='7069', dr.message = concat(dr.message, ',부적절한 문자사용'),dr.remark2 = date_format(now(), '%Y-%m-%d %H:%i:%S') where userid = ? msgid = ?", useridt, msgKey)
						}
					}
				}
			} else {
				stdlog.Println("Lgu OTP SMS Table Insert 처리 : ", len(ossmsStrs), " - LG_OTP_SC_TRAN")
			}
			ossmsStrs = nil
			ossmsValues = nil
		}

		if len(osmmsStrs) > 500 {
			stmt := fmt.Sprintf("insert into LG_OTP_MMS_MSG(SUBJECT, PHONE, CALLBACK, REQDATE, MSG, FILE_CNT, FILE_PATH1, FILE_PATH2, FILE_PATH3, ETC1, ETC2, ETC3, KISA_ORIGCODE) values %s", s.Join(osmmsStrs, ","))
			_, err := db.Exec(stmt, osmmsValues...)

			if err != nil {
				for i := 0; i < len(osmmsValues); i = i + 12 {
					eQuery := fmt.Sprintf("SUBJECT, PHONE, CALLBACK, REQDATE, MSG, FILE_CNT, FILE_PATH1, FILE_PATH2, FILE_PATH3, ETC1, ETC2, ETC3, KISA_ORIGCODE) "+
						"values('%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v')", osmmsValues[i], osmmsValues[i+1], osmmsValues[i+2], osmmsValues[i+3], osmmsValues[i+4], osmmsValues[i+5], osmmsValues[i+6], osmmsValues[i+7], osmmsValues[i+8], osmmsValues[i+9], osmmsValues[i+10], osmmsValues[i+11], osmmsValues[i+12])
					_, err := db.Exec(eQuery)
					if err != nil {
						msgKey := fmt.Sprintf("%v", osmmsValues[i+9])
						useridt := fmt.Sprintf("%v", osmmsValues[i+10])
						stdlog.Println("Lgu OTP LMS Table Insert 처리 중 오류 발생 : "+err.Error(), " - DHN Msg Key : ", msgKey)
						errcodemsg := err.Error()
						if s.Index(errcodemsg, "1366") > 0 {
							db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.code='7069', dr.message = concat(dr.message, ',부적절한 문자사용'),dr.remark2 = date_format(now(), '%Y-%m-%d %H:%i:%S') where userid = '" + useridt + "' and msgid = '" + msgKey + "'")
						}
					}
				}
			} else {
				stdlog.Println("Lgu OTP MMS Table Insert 처리 : ", len(osmmsStrs), " - LG_OTP_MMS_MSG")
			}
			osmmsStrs = nil
			osmmsValues = nil
		}

		// 알림톡 발송 성공 혹은 문자 발송이 아니면
		// API_RESULT 성공 처리 함.
		if len(msg_sms.String) > 0 && len(sms_sender.String) > 0 { // msg_sms 가 와 sms_sender 에 값이 있으면 LG 발송 함.
			phnstr = reg.ReplaceAllString(phnstr, "")
			if s.HasPrefix(phnstr, "82") {
				phnstr = "0" + phnstr[2:len(phnstr)]
			}
			if s.EqualFold(sms_kind.String, "S") {
				if msgLen.Int64 <= 90 || s.EqualFold(sms_len_check.String, "N") {
					ossmsStrs = append(ossmsStrs, "(?,?,?,?,?,?,?,?)")
					ossmsValues = append(ossmsValues, time.Now().Format("2006-01-02 15:04:05"))
					ossmsValues = append(ossmsValues, phnstr)
					ossmsValues = append(ossmsValues, sms_sender.String)
					ossmsValues = append(ossmsValues, msg_sms.String)
					ossmsValues = append(ossmsValues, msgid.String)
					ossmsValues = append(ossmsValues, userid.String)
					ossmsValues = append(ossmsValues, group_no)
					ossmsValues = append(ossmsValues, config.Conf.KISA_CODE)
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
				osmmsStrs = append(osmmsStrs, "(?,?,?,?,?,?,?,?,?,?,?,?,?)")
				osmmsValues = append(osmmsValues, sms_lms_tit.String)
				osmmsValues = append(osmmsValues, phnstr)
				osmmsValues = append(osmmsValues, sms_sender.String)
				osmmsValues = append(osmmsValues, time.Now().Format("2006-01-02 15:04:05"))
				osmmsValues = append(osmmsValues, msg_sms.String)
				osmmsValues = append(osmmsValues, file_cnt)
				osmmsValues = append(osmmsValues, mms_file1.String)
				osmmsValues = append(osmmsValues, mms_file2.String)
				osmmsValues = append(osmmsValues, mms_file3.String)
				osmmsValues = append(osmmsValues, msgid.String)
				osmmsValues = append(osmmsValues, userid.String)
				osmmsValues = append(osmmsValues, group_no)
				osmmsValues = append(osmmsValues, config.Conf.KISA_CODE)
				lmscnt++
			}

		} else {
			db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.code='7011', dr.message = concat(dr.message, ',문자 발송 정보 누락'),dr.remark2 = date_format(now(), '%Y-%m-%d %H:%i:%S') where userid = '" + userid.String + "' and msgid = '" + msgid.String + "'")
		}

	}

	if len(ossmsStrs) > 0 {
		stmt := fmt.Sprintf("insert into LG_OTP_SC_TRAN(TR_SENDDATE,TR_PHONE,TR_CALLBACK, TR_MSG, TR_ETC1, TR_ETC2, TR_ETC3, TR_KISAORIGCODE) values %s", s.Join(ossmsStrs, ","))
		_, err := db.ExecContext(ctx, stmt, ossmsValues...)

		if err != nil {
			for i := 0; i < len(ossmsValues); i = i + 7 {
				eQuery := fmt.Sprintf("insert into LG_OTP_SC_TRAN(TR_SENDDATE,TR_PHONE,TR_CALLBACK, TR_MSG, TR_ETC1, TR_ETC2, TR_ETC3, TR_KISAORIGCODE) "+
					"values('%v','%v','%v','%v','%v','%v','%v')", ossmsValues[i], ossmsValues[i+1], ossmsValues[i+2], ossmsValues[i+3], ossmsValues[i+4], ossmsValues[i+5], ossmsValues[i+6], ossmsValues[i+7])
				_, err := db.ExecContext(ctx, eQuery)
				if err != nil {
					msgKey := fmt.Sprintf("%v", ossmsValues[i+4])
					useridt := fmt.Sprintf("%v", ossmsValues[i+5])
					stdlog.Println("Lgu OTP SMS Table Insert 처리 중 오류 발생 : "+err.Error(), " - DHN Msg Key : ", msgKey)
					errcodemsg := err.Error()
					if s.Index(errcodemsg, "1366") > 0 {
						db.ExecContext(ctx, "update DHN_RESULT dr set dr.result = 'Y', dr.code='7069', dr.message = concat(dr.message, ',부적절한 문자사용'),dr.remark2 = date_format(now(), '%Y-%m-%d %H:%i:%S') where userid = ? and  msgid = ?", useridt, msgKey)
					}
				}
			}
		} else {
			stdlog.Println("Lgu OTP SMS Table Insert 처리 : ", len(ossmsStrs), " - LG_OTP_SC_TRAN")
		}

	}

	if len(osmmsStrs) > 0 {
		stmt := fmt.Sprintf("insert into LG_OTP_MMS_MSG(SUBJECT, PHONE, CALLBACK, REQDATE, MSG, FILE_CNT, FILE_PATH1, FILE_PATH2, FILE_PATH3, ETC1, ETC2, ETC3, KISA_ORIGCODE) values %s", s.Join(osmmsStrs, ","))
		_, err := db.Exec(stmt, osmmsValues...)

		if err != nil {
			for i := 0; i < len(osmmsValues); i = i + 12 {
				eQuery := fmt.Sprintf("SUBJECT, PHONE, CALLBACK, REQDATE, MSG, FILE_CNT, FILE_PATH1, FILE_PATH2, FILE_PATH3, ETC1, ETC2, ETC3, KISA_ORIGCODE) "+
					"values('%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v')", osmmsValues[i], osmmsValues[i+1], osmmsValues[i+2], osmmsValues[i+3], osmmsValues[i+4], osmmsValues[i+5], osmmsValues[i+6], osmmsValues[i+7], osmmsValues[i+8], osmmsValues[i+9], osmmsValues[i+10], osmmsValues[i+11], osmmsValues[i+12])
				_, err := db.Exec(eQuery)
				if err != nil {
					msgKey := fmt.Sprintf("%v", osmmsValues[i+9])
					useridt := fmt.Sprintf("%v", osmmsValues[i+10])
					stdlog.Println("Lgu OTP LMS Table Insert 처리 중 오류 발생 : "+err.Error(), " - DHN Msg Key : ", msgKey)
					errcodemsg := err.Error()
					if s.Index(errcodemsg, "1366") > 0 {
						db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.code='7069', dr.message = concat(dr.message, ',부적절한 문자사용'),dr.remark2 = date_format(now(), '%Y-%m-%d %H:%i:%S') where userid = '" + useridt + "' and msgid = '" + msgKey + "'")
					}
				}
			}
		} else {
			stdlog.Println("Lgu OTP MMS Table Insert 처리 : ", len(osmmsStrs), " - LG_OTP_MMS_MSG")
		}
	}

	if scnt > 0 || smscnt > 0 || lmscnt > 0 || fcnt > 0 {
		stdlog.Println(group_no, " Lgu OTP 문자 발송 처리 완료 ( ", tcnt, " ) : 성공 -", scnt, " , SMS -", smscnt, " , LMS -", lmscnt, ", 실패 - ", fcnt, "  >> Process cnt : ", procCnt)
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

func utf8TOeuckr(str string) string {
	sText := []byte(str)
	eText := make([]byte, hex.EncodedLen(len(sText)))
	hex.Encode(eText, sText)

	temp := string(eText)
	temp = s.Replace(temp, "e2808b", "", -1)

	bs, _ := hex.DecodeString(temp)

	return string(bs)
}
