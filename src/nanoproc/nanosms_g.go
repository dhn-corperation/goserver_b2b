package nanoproc

import (
	"database/sql"
	"fmt"
	config "mycs/src/kaoconfig"
	databasepool "mycs/src/kaodatabasepool"

	//"strconv"
	s "strings"
	"sync"
	"time"

	"context"

	_ "github.com/go-sql-driver/mysql"
)

func NanoSMSProcess_G(ctx context.Context) {
	var wg sync.WaitGroup

	var db = databasepool.DB
	var errlog = config.Stdlog
	var oshotTable [][]string
	var otable sql.NullString

	var OshotQuery = "select distinct ifnull(a.oshot, '') as table_name from DHN_CLIENT_LIST a where a.use_flag = 'Y' and ifnull(a.dest,'OSHOT') = 'NANO' "

	OshotTable, err := db.Query(OshotQuery)

	if err != nil {
		errlog.Fatal("DHN CLIENT LIST 조회 오류 ")
	}
	defer OshotTable.Close()

	for OshotTable.Next() {
		OshotTable.Scan(&otable)
		oshotTable = append(oshotTable, []string{otable.String})
	}
	errlog.Println("Nano SMS_G length : ", len(oshotTable))
	for {
		select {
		case <-ctx.Done():

			config.Stdlog.Println("Nano SMS process_G가 20초 후에 종료 됨.")
			time.Sleep(20 * time.Second)
			config.Stdlog.Println("Nano SMS process_G 종료 완료")
			return
		default:

			for _, tableName := range oshotTable {
				var t = time.Now()

				if t.Day() < 3 {
					wg.Add(1)
					go pre_smsProcess_G(&wg, tableName[0])
				}

				wg.Add(1)
				go smsProcess_G(&wg, tableName[0])
			}

			wg.Wait()
		}
	}

}

func smsProcess_G(wg *sync.WaitGroup, tablename string) {

	defer wg.Done()
	var db = databasepool.DB
	var errlog = config.Stdlog

	var isProc = true
	var t = time.Now()
	var monthStr = fmt.Sprintf("%d%02d", t.Year(), t.Month())

	var SMSTable = "SMS_G_LOG_" + monthStr

	//발송 6시간 지난 메세지는 응답과 상관 없이 성공 처리 함.

	var groupQuery = "select tr_etc9 as cb_msg_id, tr_rsltstat as SendResult, TR_REALSENDDATE as SendDT, tr_num as MsgID, tr_net as telecom,tr_etc10 as userid  from " + SMSTable + " a where a.TR_SENDSTAT = '2' and  a.tr_etc8 ='Y'"

	groupRows, err := db.Query(groupQuery)
	if err != nil {
		errlog.Println("Nano SMS_G 조회 중 오류 발생")
		errcode := err.Error()

		if s.Index(errcode, "1146") > 0 {
			db.Exec("Create Table IF NOT EXISTS " + SMSTable + " like SMS_LOG")
			errlog.Println(SMSTable + " 생성 !!")

		} else {
			//errlog.Fatal(err)
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

			resultCode := NanoCode(sendresult.String)

			if !s.EqualFold(resultCode, "7006") {

				var errcode = resultCode

				val := CodeMessage(resultCode)

				db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.message_type = 'PH', dr.code = '" + errcode + "', dr.message = concat(dr.message, '," + val + "'), dr.remark1 = '" + telecom.String + "', dr.remark2 = '" + senddt.String + "' where  userid='" + userid.String + "' and msgid = '" + cb_msg_id.String + "'")
			} else {
				db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.message_type = 'PH', dr.code = '0000', dr.message = '', dr.remark1 = '" + tr_net + "', dr.remark2 = '" + senddt.String + "' where  userid='" + userid.String + "' and msgid = '" + cb_msg_id.String + "'")
			}

			db.Exec("update " + SMSTable + " set tr_etc8 = 'N' where tr_num = '" + msgid.String + "'")
		}
	}
}

func pre_smsProcess_G(wg *sync.WaitGroup, tablename string) {

	defer wg.Done()
	var db = databasepool.DB
	var errlog = config.Stdlog

	var isProc = true
	var t = time.Now().Add(time.Hour * -96)
	var monthStr = fmt.Sprintf("%d%02d", t.Year(), t.Month())

	var SMSTable = "SMS_G_LOG_" + monthStr

	//발송 6시간 지난 메세지는 응답과 상관 없이 성공 처리 함.

	var groupQuery = "select tr_etc9 as cb_msg_id, tr_rsltstat as SendResult, TR_REALSENDDATE as SendDT, tr_num as MsgID, tr_net as telecom,tr_etc10 as userid  from " + SMSTable + " a where a.TR_SENDSTAT = '2' and  a.tr_etc8 ='Y'"

	groupRows, err := db.Query(groupQuery)
	if err != nil {
		errlog.Println("Nano SMS_G 조회 중 오류 발생")
		errcode := err.Error()

		if s.Index(errcode, "1146") > 0 {
			db.Exec("Create Table IF NOT EXISTS " + SMSTable + " like SMS_LOG")
			errlog.Println(SMSTable + " 생성 !!")

		} else {
			//errlog.Fatal(err)
		}

		isProc = false
		return
	}
	defer groupRows.Close()

	if isProc {

		for groupRows.Next() {
			var cb_msg_id, sendresult, senddt, msgid, telecom, userid sql.NullString

			groupRows.Scan(&cb_msg_id, &sendresult, &senddt, &msgid, &telecom, &userid)

			resultCode := NanoCode(sendresult.String)

			if !s.EqualFold(resultCode, "7006") {

				//numcode, _ := strconv.Atoi(sendresult.String)
				var errcode = resultCode

				val := CodeMessage(resultCode)

				db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.message_type = 'PH', dr.code = '" + errcode + "', dr.message = concat(dr.message, '," + val + "'), dr.remark1 = '" + telecom.String + "', dr.remark2 = '" + senddt.String + "' where  userid='" + userid.String + "' and msgid = '" + cb_msg_id.String + "'")
			} else {
				db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.message_type = 'PH', dr.code = '0000', dr.message = '', dr.remark1 = '" + telecom.String + "', dr.remark2 = '" + senddt.String + "' where  userid='" + userid.String + "' and msgid = '" + cb_msg_id.String + "'")
			}

			db.Exec("update " + SMSTable + " set tr_etc8 = 'N' where tr_num = '" + msgid.String + "'")
		}
	}
}
