package nanoproc

import (
	"database/sql"
	"fmt"
	config "mycs/src/kaoconfig"
	databasepool "mycs/src/kaodatabasepool"

	//"strconv"

	//	"log"
	s "strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"context"
)

func NanoLMSProcess(ctx context.Context) {
	var wg sync.WaitGroup

	for {
	
			select {
				case <- ctx.Done():
			
			    config.Stdlog.Println("Nano LMS process가 20초 후에 종료 됨.")
			    time.Sleep(20 * time.Second)
			    config.Stdlog.Println("Nano LMS process 종료 완료")
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
	var db = databasepool.DB
	var errlog = config.Stdlog

	var isProc = true
	var t = time.Now()
	var monthStr = fmt.Sprintf("%d%02d", t.Year(), t.Month())

	var MMSTable = "MMS_LOG_" + monthStr

	//발송 6시간 지난 메세지는 응답과 상관 없이 성공 처리 함.

	var groupQuery = "select etc9 as cb_msg_id, rslt as SendResult, File_Path1,sentdate, msgkey as MsgID, telcoinfo as telecom, etc10 as userid  from " + MMSTable + " a where a.status = '3' and a.ETC8 = 'Y' "

	groupRows, err := db.Query(groupQuery)
	if err != nil {
		//errlog.Println("Nano MMS 조회 중 오류 발생", groupQuery)
		errcode := err.Error()
		errlog.Println("Nano MMS 조회 중 오류 발생", groupQuery, errcode)

		if s.Index(errcode, "1146") > 0 {
			db.Exec("Create Table IF NOT EXISTS " + MMSTable + " like MMS_LOG")
			errlog.Println(MMSTable + " 생성 !!")
		} else {
			//errlog.Fatal(groupQuery)
		}

		isProc = false
		return
	}
	defer groupRows.Close()

	if isProc {

		for groupRows.Next() {
			var cb_msg_id, sendresult, file_path1, senddt, msgid, telecom, userid sql.NullString

			groupRows.Scan(&cb_msg_id, &sendresult, &file_path1, &senddt, &msgid, &telecom, &userid)

			tr_net := telecom.String
 
			/*
				var msg_type = "LMS"

				if len(file_path1.String) > 1 {
					msg_type = "MMS"
				}
			*/
			resultCode := NanoCode(sendresult.String)

			if !s.EqualFold(resultCode, "7006") {

				var errcode = resultCode

				val := CodeMessage(resultCode)
		
				db.Exec("update DHN_RESULT dr set dr.message_type = 'PH', dr.result = 'Y', dr.code = '" + errcode + "', dr.message = concat(dr.message, '," + val + "'), dr.remark1 = '" + telecom.String + "', dr.remark2 = '" + senddt.String + "' where userid='" + userid.String + "' and msgid = '" + cb_msg_id.String + "'")
			} else {
				db.Exec("update DHN_RESULT dr set dr.message_type = 'PH', dr.result = 'Y', dr.code = '0000', dr.message = '', dr.remark1 = '" + tr_net + "', dr.remark2 = '" + senddt.String + "' where userid='" + userid.String + "' and  msgid = '" + cb_msg_id.String + "'")
			}

			db.Exec("update " + MMSTable + " set etc8 = 'N' where msgkey = '" + msgid.String + "'")
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

	var groupQuery = "select etc9 as cb_msg_id, rslt as SendResult, File_Path1,sentdate, msgkey as MsgID, telcoinfo as telecom, etc10 as userid  from " + MMSTable + " a where a.status = '3' and a.ETC8 = 'Y' "

	groupRows, err := db.Query(groupQuery)
	if err != nil {
		isProc = false
		return
	}
	defer groupRows.Close()

	if isProc {

		for groupRows.Next() {
			var cb_msg_id, sendresult, file_path1, senddt, msgid, telecom, userid sql.NullString

			groupRows.Scan(&cb_msg_id, &sendresult, &file_path1, &senddt, &msgid, &telecom, &userid)

			/*
				var msg_type = "LMS"

				if len(file_path1.String) > 1 {
					msg_type = "MMS"
				}
			*/
			resultCode := NanoCode(sendresult.String)

			if !s.EqualFold(resultCode, "7006") {
				var errcode = resultCode

				val := CodeMessage(resultCode)
				

				db.Exec("update DHN_RESULT dr set dr.message_type = 'PH', dr.result = 'Y', dr.code = '" + errcode + "', dr.message = concat(dr.message, '," + val + "'), dr.remark1 = '" + telecom.String + "', dr.remark2 = '" + senddt.String + "' where userid='" + userid.String + "' and  msgid = '" + cb_msg_id.String + "'")
			} else {
				db.Exec("update DHN_RESULT dr set dr.message_type = 'PH', dr.result = 'Y', dr.code = '0000', dr.message = '', dr.remark1 = '" + telecom.String + "', dr.remark2 = '" + senddt.String + "' where  userid='" + userid.String + "' and msgid = '" + cb_msg_id.String + "'")
			}

			db.Exec("update " + MMSTable + " set etc8 = 'N' where msgkey = '" + msgid.String + "'")
		}
	}

}
