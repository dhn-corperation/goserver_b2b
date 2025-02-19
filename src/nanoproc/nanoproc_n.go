package nanoproc

import (
	"database/sql"
	"fmt"
	"regexp"
	s "strings"
	"time"
	"context"

	config "mycs/src/kaoconfig"
	databasepool "mycs/src/kaodatabasepool"
)

func NanoProcess_N(user_id string, ctx context.Context) {
	config.Stdlog.Println(user_id, " - Nano Process_N 시작 됨.")
	procCnt_N := 0

	for {
		if procCnt_N < 5 {

			select {
			case <-ctx.Done():

				uid := ctx.Value("user_id")
				config.Stdlog.Println(uid, " - Nano process_N가 10초 후에 종료 됨.")
				time.Sleep(10 * time.Second)
				config.Stdlog.Println(uid, " - Nano process_N 종료 완료")
				return
			default:

				var count int
				tickSql := `select
								count(msgid) as cnt
							from
								DHN_RESULT dr
							where
								dr.result = 'P'
								and dr.send_group is null
								and ifnull(dr.reserve_dt, '00000000000000') <= date_format(now(), '%Y%m%d%H%i%S')
								and userid = '` + user_id + `'
								and sms_sender not like '010%'
							limit 1 `
				cnterr := databasepool.DB.QueryRow(tickSql).Scan(&count)

				if cnterr != nil {
					config.Stdlog.Println(user_id, " - Nano Process_N DHN_RESULT Table - select error : " + cnterr.Error())
					time.Sleep(10 * time.Second)
				} else {

					if count > 0 {

						var startNow = time.Now()
						var group_no = fmt.Sprintf("%02d%02d%02d%02d%06d", startNow.Day(), startNow.Hour(), startNow.Minute(), startNow.Second(), (startNow.Nanosecond() / 1000))

						upError := updateReqeust_N(group_no, user_id)
						if upError != nil {
							config.Stdlog.Println(user_id, " - Nano Process_N Group_N No Update error : ", upError, " / group_no : ", group_no)
						} else {
							go func() {
								procCnt_N++
								config.Stdlog.Println(user_id, " - Nano Process_N 발송 처리 시작 ( ", group_no, " ) : ( Proc Cnt :", procCnt_N, ") - START")
								defer func() {
									procCnt_N--
								}()
								resProcess_N(group_no, user_id, procCnt_N)
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

func updateReqeust_N(group_no, user_id string) error {

	tx, err := databasepool.DB.Begin()
	if err != nil {
		return err
	}

	defer func() error {
		//config.Stdlog.Println("Group No Update End", group_no)
		if err != nil {
			tx.Rollback()
			return err
		}
		err = tx.Commit()
		return err
	}()

	config.Stdlog.Println(user_id, " - Nano Process_N Nano Group_N No Update 시작", group_no)

	gudQuery := `update	DHN_RESULT dr
						set	send_group = '` + group_no + `'
					where result = 'P'
						and send_group is null
						and ifnull(reserve_dt, '00000000000000') <= date_format(now(), '%Y%m%d%H%i%S')
						and userid = '` + user_id + `'
						and sms_sender not like '010%' 
					LIMIT 500
						`
	_, err = tx.Query(gudQuery)

	if err != nil {
		config.Stdlog.Println(user_id, " - Nano Process_N Group_N NO Update - Select error : ( "+group_no+" ) : "+err.Error())
		config.Stdlog.Println(gudQuery)
		return err
	}

	return nil
}

func resProcess_N(group_no, user_id string, pc int) {
	defer func(){
		if r := recover(); r != nil {
			config.Stdlog.Println(user_id, " - Nano Process_N resProcess_N panic error : ", r)
			if err, ok := r.(error); ok {
				if s.Contains(err.Error(), "connection refused") {
					for {
						config.Stdlog.Println(user_id, " - Nano Process_N resProcess_N send ping to DB")
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

	var msgid, code, message, message_type, msg_sms, phn, remark1, remark2, result, sms_lms_tit, sms_kind, sms_sender, res_dt, reserve_dt, mms_file1, mms_file2, mms_file3, userid, sms_len_check sql.NullString
	var msgLen sql.NullInt64
	var phnstr string

	ossmsStrs := []string{}
	ossmsValues := []interface{}{}

	osmmsStrs := []string{}
	osmmsValues := []interface{}{}

	var resquery = `SELECT msgid, 
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
						(case when reserve_dt = '00000000000000' then 
								now()
							when reserve_dt is null  then 
								now()
							when length(trim(reserve_dt)) < 4 then 
								now()
							else
								STR_TO_DATE(reserve_dt, '%Y%m%d%H%i%S')
						end) as  reserve_dt, 
						(select file1_path from api_mms_images aa where aa.user_id = drr.userid and aa.mms_id = drr.mms_image_id) as mms_file1, 
						(select file2_path from api_mms_images aa where aa.user_id = drr.userid and aa.mms_id = drr.mms_image_id) as mms_file2, 
						(select file3_path from api_mms_images aa where aa.user_id = drr.userid and aa.mms_id = drr.mms_image_id) as mms_file3
						,(case when sms_kind = 'S' then length(convert(REMOVE_WS(msg_sms) using euckr)) else 100 end) as msg_len
						,userid
						,(select max(sms_len_check) from DHN_CLIENT_LIST dcl where dcl.user_id = drr.userid) as sms_len_check 
					FROM DHN_RESULT drr 
					WHERE send_group = '` + group_no + `' 
						and result = 'P'
						and userid = '` + user_id + `'
						and sms_sender not like '010%' 
					order by userid
					`

	resrows, err := db.Query(resquery)

	if err != nil {
		stdlog.Println(user_id, " - Nano Process_N Result Table select error : ", err)
		stdlog.Println(resquery)
	}
	defer resrows.Close()
	scnt := 0
	smscnt := 0
	lmscnt := 0
	tcnt := 0
	reg, err := regexp.Compile("[^0-9]+")
	preOshot := ""

	for resrows.Next() {
		resrows.Scan(&msgid, &code, &message, &message_type, &msg_sms, &phn, &remark1, &remark2, &result, &sms_lms_tit, &sms_kind, &sms_sender, &res_dt, &reserve_dt, &mms_file1, &mms_file2, &mms_file3, &msgLen, &userid, &sms_len_check)

		phnstr = phn.String

		tcnt++

		if len(ossmsStrs) > 500 {
			stmt := fmt.Sprintf("insert into SMS_MSG(TR_CALLBACK,TR_PHONE,TR_MSG,TR_SENDDATE,TR_SENDSTAT,TR_MSGTYPE,TR_ETC9,TR_ETC10,TR_IDENTIFICATION_CODE,TR_ETC8) values %s", s.Join(ossmsStrs, ","))
			_, err := db.Exec(stmt, ossmsValues...)

			if err != nil {
				//stdlog.Println("Nano SMS Table Insert 처리 중 오류 발생 " + err.Error())
				for i := 0; i < len(ossmsValues); i = i + 8 {
					eQuery := fmt.Sprintf("insert into SMS_MSG(TR_CALLBACK,TR_PHONE,TR_MSG,TR_SENDDATE,TR_SENDSTAT,TR_MSGTYPE,TR_ETC9,TR_ETC10,TR_IDENTIFICATION_CODE,TR_ETC8) "+
						"values('%v','%v','%v','%v', '%v', '%v','%v', '%v', 'Y')", ossmsValues[i], ossmsValues[i+1], ossmsValues[i+2], ossmsValues[i+3], ossmsValues[i+4], ossmsValues[i+5], ossmsValues[i+6], ossmsValues[i+7], ossmsValues[i+8])
					_, err := db.Exec(eQuery)
					if err != nil {
						msgKey := fmt.Sprintf("%v", ossmsValues[i+6])
						useridt := fmt.Sprintf("%v", ossmsValues[i+7])
						stdlog.Println(user_id, " - Nano Process_N SMS Table Insert error : "+err.Error(), " - DHN Msg Key : ", msgKey)
						errcodemsg := err.Error()
						if s.Index(errcodemsg, "1366") > 0 {
							db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.code='7069', dr.message = concat(dr.message, ',부적절한 문자사용'),dr.remark2 = date_format(now(), '%Y-%m-%d %H:%i:%S') where userid = '" + useridt + "' and  msgid = '" + msgKey + "'")
						}
					}
				}
				//db.Exec("update API_RESULT ar set ar.msg_type = '" + sms_kind.String + "', result_code = '9999', error_text = '기타오류', report_time = date_format(now(), '%Y-%m-%d %H:%i:%S') where dhn_msg_id = '" + msgid.String + "'")
			} else {
				stdlog.Println(user_id, " - Nano Process_N SMS Table Insert 처리 : ", len(ossmsStrs), " - ", preOshot)
			}
			ossmsStrs = nil
			ossmsValues = nil
		}

		if len(osmmsStrs) > 500 {
			stmt := fmt.Sprintf("insert into MMS_MSG(CALLBACK,PHONE,SUBJECT,MSG,REQDATE,STATUS,FILE_CNT,FILE_PATH1,FILE_PATH2,FILE_PATH3,ETC9,ETC10,IDENTIFICATION_CODE,ETC8) values %s", s.Join(osmmsStrs, ","))
			_, err := db.Exec(stmt, osmmsValues...)

			if err != nil {
				//stdlog.Println("Nano SMS Table Insert 처리 중 오류 발생 " + err.Error())
				for i := 0; i < len(osmmsValues); i = i + 12 {
					eQuery := fmt.Sprintf("insert into MMS_MSG(CALLBACK,PHONE,SUBJECT,MSG,REQDATE,STATUS,FILE_CNT,FILE_PATH1,FILE_PATH2,FILE_PATH3,ETC9,ETC10,IDENTIFICATION_CODE,ETC8) "+
						"values('%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','Y')", osmmsValues[i], osmmsValues[i+1], osmmsValues[i+2], osmmsValues[i+3], osmmsValues[i+4], osmmsValues[i+5], osmmsValues[i+6], osmmsValues[i+7], osmmsValues[i+8], osmmsValues[i+9], osmmsValues[i+10], osmmsValues[i+11], osmmsValues[i+12])
					_, err := db.Exec(eQuery)
					if err != nil {
						msgKey := fmt.Sprintf("%v", osmmsValues[i+10])
						useridt := fmt.Sprintf("%v", osmmsValues[i+11])
						stdlog.Println(user_id, " - Nano Process_N MMS Table Insert error : "+err.Error(), " - DHN Msg Key : ", msgKey)
						errcodemsg := err.Error()
						if s.Index(errcodemsg, "1366") > 0 {
							db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.code='7069', dr.message = concat(dr.message, ',부적절한 문자사용'),dr.remark2 = date_format(now(), '%Y-%m-%d %H:%i:%S') where userid = '" + useridt + "' and msgid = '" + msgKey + "'")
						}
					}
				}
				//db.Exec("update API_RESULT ar set ar.msg_type = '" + sms_kind.String + "', result_code = '9999', error_text = '기타오류', report_time = date_format(now(), '%Y-%m-%d %H:%i:%S') where dhn_msg_id = '" + msgid.String + "'")
			} else {
				stdlog.Println(user_id, " - Nano Process_N MMS Table Insert 처리 : ", len(osmmsStrs), " - ", preOshot)
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

			if s.EqualFold(sms_kind.String, "S") {

				if msgLen.Int64 <= 90 || s.EqualFold(sms_len_check.String, "N") {

					ossmsStrs = append(ossmsStrs, "(?,?,?,?,?,?,?,?,?,'Y')")
					ossmsValues = append(ossmsValues, sms_sender.String)
					ossmsValues = append(ossmsValues, phnstr)
					ossmsValues = append(ossmsValues, msg_sms.String)
					ossmsValues = append(ossmsValues, reserve_dt.String)
					ossmsValues = append(ossmsValues, "0")
					ossmsValues = append(ossmsValues, "0")
					ossmsValues = append(ossmsValues, msgid.String)
					ossmsValues = append(ossmsValues, userid.String)
					ossmsValues = append(ossmsValues, config.Conf.KISA_CODE)
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

				osmmsStrs = append(osmmsStrs, "( ?,?,?,?,?,?,?,?,?,?,?,?,?,'Y')")

				osmmsValues = append(osmmsValues, sms_sender.String)
				osmmsValues = append(osmmsValues, phnstr)
				osmmsValues = append(osmmsValues, sms_lms_tit.String)
				osmmsValues = append(osmmsValues, msg_sms.String)
				osmmsValues = append(osmmsValues, reserve_dt.String)
				osmmsValues = append(osmmsValues, "0")
				osmmsValues = append(osmmsValues, filecnt)
				osmmsValues = append(osmmsValues, mms_file1.String)
				osmmsValues = append(osmmsValues, mms_file2.String)
				osmmsValues = append(osmmsValues, mms_file3.String)

				osmmsValues = append(osmmsValues, msgid.String)
				osmmsValues = append(osmmsValues, userid.String)
				osmmsValues = append(osmmsValues, config.Conf.KISA_CODE)
				lmscnt++
			}

		} else {
			db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.code='7011', dr.message = concat(dr.message, ',문자 발송 정보 누락'),dr.remark2 = date_format(now(), '%Y-%m-%d %H:%i:%S') where userid = '" + userid.String + "' and msgid = '" + msgid.String + "'")
		}

	}

	if len(ossmsStrs) > 0 {
		stmt := fmt.Sprintf("insert into SMS_MSG(TR_CALLBACK,TR_PHONE,TR_MSG,TR_SENDDATE,TR_SENDSTAT,TR_MSGTYPE,TR_ETC9,TR_ETC10,TR_IDENTIFICATION_CODE,TR_ETC8) values %s", s.Join(ossmsStrs, ","))
		_, err := db.Exec(stmt, ossmsValues...)

		if err != nil {
			//stdlog.Println("Nano SMS Table Insert 처리 중 오류 발생 " + err.Error())
			for i := 0; i < len(ossmsValues); i = i + 8 {
				eQuery := fmt.Sprintf("insert into SMS_MSG(TR_CALLBACK,TR_PHONE,TR_MSG,TR_SENDDATE,TR_SENDSTAT,TR_MSGTYPE,TR_ETC9,TR_ETC10,TR_IDENTIFICATION_CODE,TR_ETC8) "+
					"values('%v','%v','%v','%v', '%v', '%v','%v', '%v', 'Y')", ossmsValues[i], ossmsValues[i+1], ossmsValues[i+2], ossmsValues[i+3], ossmsValues[i+4], ossmsValues[i+5], ossmsValues[i+6], ossmsValues[i+7], ossmsValues[i+8])
				_, err := db.Exec(eQuery)
				if err != nil {
					msgKey := fmt.Sprintf("%v", ossmsValues[i+6])
					useridt := fmt.Sprintf("%v", ossmsValues[i+7])
					stdlog.Println(user_id, " - Nano Process_N SMS Table Insert error : "+err.Error(), " - DHN Msg Key : ", msgKey)
					errcodemsg := err.Error()
					if s.Index(errcodemsg, "1366") > 0 {
						db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.code='7069', dr.message = concat(dr.message, ',부적절한 문자사용'),dr.remark2 = date_format(now(), '%Y-%m-%d %H:%i:%S') where userid = '" + useridt + "' and  msgid = '" + msgKey + "'")
					}
				}
			}
			//db.Exec("update API_RESULT ar set ar.msg_type = '" + sms_kind.String + "', result_code = '9999', error_text = '기타오류', report_time = date_format(now(), '%Y-%m-%d %H:%i:%S') where dhn_msg_id = '" + msgid.String + "'")
		} else {
			stdlog.Println(user_id, " - Nano Process_N SMS Table Insert 처리 : ", len(ossmsStrs), " - ", preOshot)
		}

	}

	if len(osmmsStrs) > 0 {
		stmt := fmt.Sprintf("insert into MMS_MSG(CALLBACK,PHONE,SUBJECT,MSG,REQDATE,STATUS,FILE_CNT,FILE_PATH1,FILE_PATH2,FILE_PATH3,ETC9,ETC10,IDENTIFICATION_CODE,ETC8) values %s", s.Join(osmmsStrs, ","))
		_, err := db.Exec(stmt, osmmsValues...)

		if err != nil {
			//stdlog.Println("Nano SMS Table Insert 처리 중 오류 발생 " + err.Error())
			for i := 0; i < len(osmmsValues); i = i + 12 {
				eQuery := fmt.Sprintf("insert into MMS_MSG(CALLBACK,PHONE,SUBJECT,MSG,REQDATE,STATUS,FILE_CNT,FILE_PATH1,FILE_PATH2,FILE_PATH3,ETC9,ETC10,IDENTIFICATION_CODE,ETC8) "+
					"values('%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','Y')", osmmsValues[i], osmmsValues[i+1], osmmsValues[i+2], osmmsValues[i+3], osmmsValues[i+4], osmmsValues[i+5], osmmsValues[i+6], osmmsValues[i+7], osmmsValues[i+8], osmmsValues[i+9], osmmsValues[i+10], osmmsValues[i+11], osmmsValues[i+12])
				_, err := db.Exec(eQuery)
				if err != nil {
					msgKey := fmt.Sprintf("%v", osmmsValues[i+10])
					useridt := fmt.Sprintf("%v", osmmsValues[i+11])
					stdlog.Println(user_id, " - Nano Process_N MMS Table Insert error : "+err.Error(), " - DHN Msg Key : ", msgKey)
					errcodemsg := err.Error()
					if s.Index(errcodemsg, "1366") > 0 {
						db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.code='7069', dr.message = concat(dr.message, ',부적절한 문자사용'),dr.remark2 = date_format(now(), '%Y-%m-%d %H:%i:%S') where userid = '" + useridt + "' and msgid = '" + msgKey + "'")
					}
				}
			}
			//db.Exec("update API_RESULT ar set ar.msg_type = '" + sms_kind.String + "', result_code = '9999', error_text = '기타오류', report_time = date_format(now(), '%Y-%m-%d %H:%i:%S') where dhn_msg_id = '" + msgid.String + "'")
		} else {
			stdlog.Println(user_id, " - Nano Process_N MMS Table Insert 처리 : ", len(osmmsStrs), " - ", preOshot)
		}

	}

	if scnt > 0 || smscnt > 0 || lmscnt > 0 {
		stdlog.Println(user_id, " - Nano Process_N 발송 처리 완료 ( ", group_no, " ) : 성공 - ", scnt, " , SMS - ", smscnt, " , LMS - ", lmscnt, ", 총 - ", tcnt, " : ( Proc Cnt :", pc, ") - END")
	}
}

/*
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
*/
