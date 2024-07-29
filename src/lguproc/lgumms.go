package lguproc

import (
	"database/sql"
	"fmt"
	config "mycs/src/kaoconfig"
	databasepool "mycs/src/kaodatabasepool"

	s "strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"context"
)

func LMSProcess(ctx context.Context) {
	var wg sync.WaitGroup

	var db = databasepool.DB
	var errlog = config.Stdlog
	var lguTable []string
	var ltable sql.NullString

	var LguQuery = "select distinct ifnull(a.dest, '') as table_name from DHN_CLIENT_LIST a where a.use_flag = 'Y' and a.dest = 'LGU' "

	LguTable, err := db.Query(LguQuery)

	if err != nil {
		errlog.Fatal("lgusms / LMSProcess / DHN CLIENT LIST 조회 오류 ")
	}
	defer LguTable.Close()

	for LguTable.Next() {
		LguTable.Scan(&ltable)
		lguTable = append(lguTable, ltable.String)
	}
	errlog.Println("Lgu LMS length : ", len(lguTable))
	for {
		select {
			case <- ctx.Done():
		
		    config.Stdlog.Println("Lgu LMS process가 20초 후에 종료 됨.")
		    time.Sleep(20 * time.Second)
		    config.Stdlog.Println("Lgu LMS process 종료 완료")
		    return
		default:	
		
			for range lguTable {
				var t = time.Now()
	
				if t.Day() < 3 {
					wg.Add(1)
					go pre_mmsProcess(&wg)
				}
	
				wg.Add(1)
				go mmsProcess(&wg)
			}
			wg.Wait()
		}
	}

}

func mmsProcess(wg *sync.WaitGroup) {

	defer wg.Done()
	var db = databasepool.DB
	var errlog = config.Stdlog

	var isProc = true
	var t = time.Now()
	var monthStr = fmt.Sprintf("%d%02d", t.Year(), t.Month())

	var MMSTable = "LG_MMS_LOG_" + monthStr

	//발송 6시간 지난 메세지는 응답과 상관 없이 성공 처리 함.
	db.Exec("update LG_MMS_MSG set STATUS='3', RSLT='9018', RSLTDATE=now() where REQDATE >= DATE_SUB(now(), INTERVAL 1 MINUTE)")

	var groupQuery = "select ETC1 as cb_msg_id, RSLT as sendresult, SENTDATE as senddt, MSGKEY as msgid, TELCOINFO as telecom, ETC2 as userid  from " + MMSTable + " a where a.status = '3' and a.ETC4 is null "

	groupRows, err := db.Query(groupQuery)
	if err != nil {
		errcode := err.Error()
		errlog.Println("Lgu MMS 조회 중 오류 발생", groupQuery, errcode)

		if s.Index(errcode, "1146") > 0 {
			db.Exec("Create Table IF NOT EXISTS " + MMSTable + " like LG_MMS_MSG")
			errlog.Println(MMSTable + " 생성 !!")
		}

		isProc = false
		return
	}
	defer groupRows.Close()

	if isProc {

		for groupRows.Next() {
			var cb_msg_id, sendresult, senddt, msgid, telecom, userid sql.NullString

			groupRows.Scan(&cb_msg_id, &sendresult, &senddt, &msgid, &telecom, &userid)

			tr_net := telecom.String
 
			resultCode := LguCode(sendresult.String)

			if !s.EqualFold(resultCode, "7006") {

				var errcode = resultCode

				val := CodeMessage(resultCode)
		
				db.Exec("update DHN_RESULT dr set dr.message_type = 'PH', dr.result = 'Y', dr.code = '" + errcode + "', dr.message = concat(dr.message, '," + val + "'), dr.remark1 = '" + telecom.String + "', dr.remark2 = '" + senddt.String + "' where userid='" + userid.String + "' and msgid = '" + cb_msg_id.String + "'")
			} else {
				db.Exec("update DHN_RESULT dr set dr.message_type = 'PH', dr.result = 'Y', dr.code = '0000', dr.message = '', dr.remark1 = '" + tr_net + "', dr.remark2 = '" + senddt.String + "' where userid='" + userid.String + "' and  msgid = '" + cb_msg_id.String + "'")
			}

			db.Exec("update " + MMSTable + " set ETC4 = '1' where msgkey = '" + msgid.String + "'")
		}
	}
}

func pre_mmsProcess(wg *sync.WaitGroup) {

	defer wg.Done()
	var db = databasepool.DB

	var isProc = true
	var t = time.Now().Add(time.Hour * -96)
	var monthStr = fmt.Sprintf("%d%02d", t.Year(), t.Month())

	var MMSTable = "MMS_LOG_" + monthStr

	//발송 6시간 지난 메세지는 응답과 상관 없이 성공 처리 함.
	db.Exec("update LG_MMS_MSG set STATUS='3', RSLT='9018', RSLTDATE=now() where REQDATE >= DATE_SUB(now(), INTERVAL 1 MINUTE)")

	var groupQuery = "select ETC1 as cb_msg_id, RSLT as sendresult, SENTDATE as senddt, MSGKEY as msgid, TELCOINFO as telecom, ETC2 as userid  from " + MMSTable + " a where a.status = '3' and a.ETC4 is null "

	groupRows, err := db.Query(groupQuery)
	if err != nil {
		isProc = false
		return
	}
	defer groupRows.Close()

	if isProc {

		for groupRows.Next() {
			var cb_msg_id, sendresult, senddt, msgid, telecom, userid sql.NullString

			groupRows.Scan(&cb_msg_id, &sendresult, &senddt, &msgid, &telecom, &userid)

			tr_net := telecom.String
 
			resultCode := LguCode(sendresult.String)

			if !s.EqualFold(resultCode, "7006") {

				var errcode = resultCode

				val := CodeMessage(resultCode)
		
				db.Exec("update DHN_RESULT dr set dr.message_type = 'PH', dr.result = 'Y', dr.code = '" + errcode + "', dr.message = concat(dr.message, '," + val + "'), dr.remark1 = '" + telecom.String + "', dr.remark2 = '" + senddt.String + "' where userid='" + userid.String + "' and msgid = '" + cb_msg_id.String + "'")
			} else {
				db.Exec("update DHN_RESULT dr set dr.message_type = 'PH', dr.result = 'Y', dr.code = '0000', dr.message = '', dr.remark1 = '" + tr_net + "', dr.remark2 = '" + senddt.String + "' where userid='" + userid.String + "' and  msgid = '" + cb_msg_id.String + "'")
			}

			db.Exec("update " + MMSTable + " set ETC4 = '1' where msgkey = '" + msgid.String + "'")
		}
	}
}
