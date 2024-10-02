package otplguproc

import (
	"database/sql"
	"fmt"
	config "mycs/src/kaoconfig"
	databasepool "mycs/src/kaodatabasepool"

	s "strings"
	"sync"
	"time"

	"mycs/src/lguproc"

	_ "github.com/go-sql-driver/mysql"
	"context"
)

func SMSProcess(ctx context.Context) {
	var wg sync.WaitGroup

	for {
		select {
			case <- ctx.Done():
		
		    config.Stdlog.Println("Lgu OTP SMS process가 10초 후에 종료 됨.")
		    time.Sleep(10 * time.Second)
		    config.Stdlog.Println("Lgu OTP SMS process 종료 완료")
		    return
		default:	
			var t = time.Now()

			if t.Day() < 3 {
				wg.Add(1)
				go pre_smsProcess(&wg)
			}

			wg.Add(1)
			go smsProcess(&wg)
	
			wg.Wait()
		}
	}

}

func smsProcess(wg *sync.WaitGroup) {
	defer wg.Done()
	defer func(){
		if r := recover(); r != nil {
			config.Stdlog.Println("OTPLGU smsProcess panic 발생 원인 : ", r)
			if err, ok := r.(error); ok {
				if s.Contains(err.Error(), "connection refused") {
					for {
						config.Stdlog.Println("OTPLGU smsProcess send ping to DB")
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
	var errlog = config.Stdlog

	var isProc = true
	var t = time.Now()
	var monthStr = fmt.Sprintf("%d%02d", t.Year(), t.Month())

	var SMSTable = "LG_OTP_SC_LOG_" + monthStr

	//발송 6시간 지난 메세지는 응답과 상관 없이 성공 처리 함.
	// db.Exec("update LG_SC_TRAN set TR_SENDSTAT='2', TR_RSLTSTAT='98', TR_RSLTDATE=now() where TR_SENDDATE < DATE_SUB(now(), INTERVAL 6 HOUR)")

	var groupQuery = "select TR_ETC1 as cb_msg_id, TR_RSLTSTAT as sendResult, TR_RSLTDATE as senddt, TR_NUM as msgid, TR_NET as telecom, a.TR_ETC2 as userid  from " + SMSTable + " a where a.TR_SENDSTAT = '2' and  a.TR_ETC4 is null"

	groupRows, err := db.Query(groupQuery)
	if err != nil {
		errcode := err.Error()
		errlog.Println("Lgu OTP SMS 조회 중 오류 발생", groupQuery, errcode)

		if s.Index(errcode, "1146") > 0 {
			db.Exec("Create Table IF NOT EXISTS " + SMSTable + " like LG_OTP_SC_TRAN")
			errlog.Println(SMSTable + " 생성 !!")

		}

		isProc = false
		return
	}
	defer groupRows.Close()

	if isProc {

		for groupRows.Next() {
			var cb_msg_id, sendresult, senddt, msgid, telecom, userid sql.NullString
			var sendDt string

			groupRows.Scan(&cb_msg_id, &sendresult, &senddt, &msgid, &telecom, &userid)

			tr_net := "ETC"
			
			if s.EqualFold(telecom.String, "011") {
				tr_net = "SKT"
			} else if s.EqualFold(telecom.String, "016") {
				tr_net = "KTF"
			} else if s.EqualFold(telecom.String, "019") {
				tr_net = "LGT"
			}

			resultCode := lguproc.LguCode(sendresult.String)

			if !senddt.Valid {
				sendDt = time.Now().Format("2006-01-02 15:04:05")
			} else {
				sendDt = senddt.String
			}

			if !s.EqualFold(resultCode, "7006") {

				var errcode = resultCode

				val := lguproc.CodeMessage(resultCode)
				
				db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.message_type = 'PH', dr.code = '" + errcode + "', dr.message = concat(dr.message, '," + val + "'), dr.remark1 = '" + telecom.String + "', dr.remark2 = '" + sendDt + "' where  userid='" + userid.String + "' and msgid = '" + cb_msg_id.String + "'")
			} else {
				db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.message_type = 'PH', dr.code = '0000', dr.message = '', dr.remark1 = '" + tr_net + "', dr.remark2 = '" + sendDt + "' where  userid='" + userid.String + "' and msgid = '" + cb_msg_id.String + "'")
			}

			db.Exec("update " + SMSTable + " set TR_ETC4 = '1' where TR_NUM = '" + msgid.String + "'")
		}
	}
}

func pre_smsProcess(wg *sync.WaitGroup) {
	defer wg.Done()
	defer func(){
		if r := recover(); r != nil {
			config.Stdlog.Println("OTPLGU smsProcess panic 발생 원인 : ", r)
			if err, ok := r.(error); ok {
				if s.Contains(err.Error(), "connection refused") {
					for {
						config.Stdlog.Println("OTPLGU smsProcess send ping to DB")
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

	var isProc = true
	var t = time.Now().Add(time.Hour * -96)
	var monthStr = fmt.Sprintf("%d%02d", t.Year(), t.Month())

	var SMSTable = "LG_OTP_SC_LOG_" + monthStr

	//발송 6시간 지난 메세지는 응답과 상관 없이 성공 처리 함.

	var groupQuery = "select TR_ETC1 as cb_msg_id, TR_RSLTSTAT as sendResult, TR_RSLTDATE as senddt, TR_NUM as msgid, TR_NET as telecom, a.TR_ETC2 as userid  from " + SMSTable + " a where a.TR_SENDSTAT = '2' and  a.TR_ETC4 is null"

	groupRows, err := db.Query(groupQuery)
	if err != nil {
		isProc = false
		return
	}
	defer groupRows.Close()

	if isProc {

		for groupRows.Next() {
			var cb_msg_id, sendresult, senddt, msgid, telecom, userid sql.NullString
			var sendDt string

			groupRows.Scan(&cb_msg_id, &sendresult, &senddt, &msgid, &telecom, &userid)

			tr_net := "ETC"
			
			if s.EqualFold(telecom.String, "011") {
				tr_net = "SKT"
			} else if s.EqualFold(telecom.String, "016") {
				tr_net = "KTF"
			} else if s.EqualFold(telecom.String, "019") {
				tr_net = "LGT"
			}

			resultCode := lguproc.LguCode(sendresult.String)

			if !senddt.Valid {
				sendDt = time.Now().Format("2006-01-02 15:04:05")
			} else {
				sendDt = senddt.String
			}

			if !s.EqualFold(resultCode, "7006") {

				var errcode = resultCode

				val := lguproc.CodeMessage(resultCode)
				
				db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.message_type = 'PH', dr.code = '" + errcode + "', dr.message = concat(dr.message, '," + val + "'), dr.remark1 = '" + telecom.String + "', dr.remark2 = '" + sendDt + "' where  userid='" + userid.String + "' and msgid = '" + cb_msg_id.String + "'")
			} else {
				db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.message_type = 'PH', dr.code = '0000', dr.message = '', dr.remark1 = '" + tr_net + "', dr.remark2 = '" + sendDt + "' where  userid='" + userid.String + "' and msgid = '" + cb_msg_id.String + "'")
			}

			db.Exec("update " + SMSTable + " set TR_ETC4 = '1' where TR_NUM = '" + msgid.String + "'")
		}
	}
}
