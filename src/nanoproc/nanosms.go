package nanoproc

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

func NanoSMSProcess(ctx context.Context, gFlag bool) {
	var wg sync.WaitGroup

	var db = databasepool.DB
	var errlog = config.Stdlog
	var nanoTable [][]string
	var ntable sql.NullString
	tail := ""
	if gFlag {
		tail = "_G"
	}

	var nanoQuery = `
	select 
		distinct dest
	from
		DHN_CLIENT_LIST
	where
		use_flag = 'Y'
		and dest ilike 'nano%'`

	NanoTable, err := db.Query(nanoQuery)

	if err != nil {
		errlog.Fatal("nanosms.go / NanoSMSProcess / DHN CLIENT LIST 조회 오류 ")
	}
	defer NanoTable.Close()

	for NanoTable.Next() {
		NanoTable.Scan(&ntable)
		nanoTable = append(nanoTable, []string{ntable.String})
	}
	errlog.Println("Nano SMS" , tail, " length : ", len(nanoTable))
	for {
			select {
				case <- ctx.Done():
			
			    config.Stdlog.Println("Nano SMS" , tail, " result process가 20초 후에 종료 됨.")
			    time.Sleep(20 * time.Second)
			    config.Stdlog.Println("Nano SMS" , tail, " result process 종료 완료")
			    return
			default:	
			
				for _, tableName := range nanoTable {
					var t = time.Now()
		
					if t.Day() < 3 {
						wg.Add(1)
						go smsProcess(&wg, tableName[0], true, tail)
					}
		
					wg.Add(1)
					go smsProcess(&wg, tableName[0], false, tail)
				}
		
				wg.Wait()
			}
	}

}

func smsProcess(wg *sync.WaitGroup, tablename string, nsFlag bool, tail string) {

	defer wg.Done()
	var db = databasepool.DB
	var errlog = config.Stdlog

	var isProc = true
	var t = time.Now()
	var monthStr = fmt.Sprintf("%d%02d", t.Year(), t.Month())
	if nsFlag {
		t = t.Add(time.Hour * -96)
	}

	var SMSTable = "SMS"+tail+"_LOG_" + monthStr

	//발송 6시간 지난 메세지는 응답과 상관 없이 성공 처리 함.

	var groupQuery = "select tr_etc9 as cb_msg_id, tr_rsltstat as SendResult, TR_REALSENDDATE as SendDT, tr_num as MsgID, tr_net as telecom,tr_etc10 as userid  from " + SMSTable + " a where a.TR_SENDSTAT = '2' and  a.tr_etc8 ='Y'"

	groupRows, err := db.Query(groupQuery)
	if err != nil {
		errcode := err.Error()
		errlog.Println("Nano SMS" , tail, " 조회 중 오류 발생", groupQuery, errcode)

		if s.Index(errcode, "relation") > 0 {
			db.Exec("Create Table IF NOT EXISTS " + SMSTable + "(LIKE mms_log INCLUDING ALL)")
			errlog.Println("nano "+SMSTable + " 생성 !!")

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
