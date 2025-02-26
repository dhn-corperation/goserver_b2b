package otplguproc

import (
	"fmt"
	"sync"
	"time"
	"context"
	s "strings"
	"database/sql"

	"mycs/cmd/lguproc"
	config "mycs/cmd/kaoconfig"
	databasepool "mycs/cmd/kaodatabasepool"

	_ "github.com/go-sql-driver/mysql"
)

func LMSProcess(ctx context.Context) {
	var wg sync.WaitGroup

	for {
		select {
			case <- ctx.Done():
		
		    config.Stdlog.Println("Lgu OTP LMS process가 10초 후에 종료 됨.")
		    time.Sleep(10 * time.Second)
		    config.Stdlog.Println("Lgu OTP LMS process 종료 완료")
		    return
		default:
			var t = time.Now()

			if t.Day() < 3 {
				wg.Add(1)
				go pre_mmsProcess(&wg)
			}

			wg.Add(1)
			go mmsProcess(&wg)
			
			wg.Wait()
		}
	}

}

func mmsProcess(wg *sync.WaitGroup) {
	defer wg.Done()
	defer func(){
		if r := recover(); r != nil {
			config.Stdlog.Println("OTPLGU mmsProcess panic 발생 원인 : ", r)
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

	var db = databasepool.DB
	var errlog = config.Stdlog

	var isProc = true
	var t = time.Now()
	var monthStr = fmt.Sprintf("%d%02d", t.Year(), t.Month())

	var MMSTable = "LG_OTP_MMS_LOG_" + monthStr

	//발송 6시간 지난 메세지는 응답과 상관 없이 성공 처리 함.
	// db.Exec("update LG_MMS_MSG set STATUS='3', RSLT='9018', RSLTDATE=now() where REQDATE < DATE_SUB(now(), INTERVAL 6 HOUR)")

	var groupQuery = "select ETC1 as cb_msg_id, RSLT as sendresult, RSLTDATE as senddt, MSGKEY as msgid, TELCOINFO as telecom, ETC2 as userid  from " + MMSTable + " a where a.status = '3' and a.ETC4 is null "

	groupRows, err := db.Query(groupQuery)
	if err != nil {
		errcode := err.Error()
		errlog.Println("Lgu OTP MMS 조회 중 오류 발생", groupQuery, errcode)

		if s.Index(errcode, "1146") > 0 {
			db.Exec("Create Table IF NOT EXISTS " + MMSTable + " like LG_OTP_MMS_MSG")
			errlog.Println(MMSTable + " 생성 !!")
		}
		time.Sleep(10 * time.Second)
		isProc = false
		return
	}
	defer groupRows.Close()

	if isProc {

		for groupRows.Next() {
			var cb_msg_id, sendresult, senddt, msgid, telecom, userid sql.NullString
			var sendDt string

			groupRows.Scan(&cb_msg_id, &sendresult, &senddt, &msgid, &telecom, &userid)

			tr_net := telecom.String
			
			if s.EqualFold(tr_net, "KT") {
				tr_net = "KTF"
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
		
				db.Exec("update DHN_RESULT dr set dr.message_type = 'PH', dr.result = 'Y', dr.code = '" + errcode + "', dr.message = concat(dr.message, '," + val + "'), dr.remark1 = '" + telecom.String + "', dr.remark2 = '" + sendDt + "' where userid='" + userid.String + "' and msgid = '" + cb_msg_id.String + "'")
			} else {
				db.Exec("update DHN_RESULT dr set dr.message_type = 'PH', dr.result = 'Y', dr.code = '0000', dr.message = '', dr.remark1 = '" + tr_net + "', dr.remark2 = '" + sendDt + "' where userid='" + userid.String + "' and  msgid = '" + cb_msg_id.String + "'")
			}

			db.Exec("update " + MMSTable + " set ETC4 = '1' where msgkey = '" + msgid.String + "'")
		}
	}
}

func pre_mmsProcess(wg *sync.WaitGroup) {
	defer wg.Done()
	defer func(){
		if r := recover(); r != nil {
			config.Stdlog.Println("OTPLGU mmsProcess panic 발생 원인 : ", r)
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
	
	var db = databasepool.DB

	var isProc = true
	var t = time.Now().Add(time.Hour * -96)
	var monthStr = fmt.Sprintf("%d%02d", t.Year(), t.Month())

	var MMSTable = "LG_OTP_MMS_LOG_" + monthStr

	//발송 6시간 지난 메세지는 응답과 상관 없이 성공 처리 함.
	// db.Exec("update LG_MMS_MSG set STATUS='3', RSLT='9018', RSLTDATE=now() where REQDATE < DATE_SUB(now(), INTERVAL 1 HOUR)")

	var groupQuery = "select ETC1 as cb_msg_id, RSLT as sendresult, RSLTDATE as senddt, MSGKEY as msgid, TELCOINFO as telecom, ETC2 as userid  from " + MMSTable + " a where a.status = '3' and a.ETC4 is null "

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

			tr_net := telecom.String
 
			resultCode := lguproc.LguCode(sendresult.String)

			if !senddt.Valid {
				sendDt = time.Now().Format("2006-01-02 15:04:05")
			} else {
				sendDt = senddt.String
			}

			if !s.EqualFold(resultCode, "7006") {

				var errcode = resultCode

				val := lguproc.CodeMessage(resultCode)
		
				db.Exec("update DHN_RESULT dr set dr.message_type = 'PH', dr.result = 'Y', dr.code = '" + errcode + "', dr.message = concat(dr.message, '," + val + "'), dr.remark1 = '" + telecom.String + "', dr.remark2 = '" + sendDt + "' where userid='" + userid.String + "' and msgid = '" + cb_msg_id.String + "'")
			} else {
				db.Exec("update DHN_RESULT dr set dr.message_type = 'PH', dr.result = 'Y', dr.code = '0000', dr.message = '', dr.remark1 = '" + tr_net + "', dr.remark2 = '" + sendDt + "' where userid='" + userid.String + "' and  msgid = '" + cb_msg_id.String + "'")
			}

			db.Exec("update " + MMSTable + " set ETC4 = '1' where msgkey = '" + msgid.String + "'")
		}
	}
}
