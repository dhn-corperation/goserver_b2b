package oshotproc

import (
	"database/sql"
	"fmt"
	config "mycs/src/kaoconfig"
	databasepool "mycs/src/kaodatabasepool"

	"strconv"
	s "strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"context"
)

func SMSProcess(ctx context.Context) {
	var wg sync.WaitGroup

	var db = databasepool.DB
	var errlog = config.Stdlog
	var oshotTable [][]string
	var otable sql.NullString

	var OshotQuery = `
	select 
		distinct lower(dest) 
	from 
		DHN_CLIENT_LIST
	where 
		use_flag = 'Y' 
		and dest ilike 'oshot%'
		and length(dest) > 1
		and dest is not null`

	OshotTable, err := db.Query(OshotQuery)

	if err != nil {
		errlog.Println(err.Error())
		errlog.Fatal("webcsms.go / SMSProcess / DHN CLIENT LIST 조회 오류 ")
	}
	defer OshotTable.Close()

	for OshotTable.Next() {
		OshotTable.Scan(&otable)
		oshotTable = append(oshotTable, []string{otable.String})
	}
	errlog.Println("Oshot SMS length : ", len(oshotTable))
	for {
			select {
				case <- ctx.Done():
			
			    config.Stdlog.Println("Oshot SMS process가 20초 후에 종료 됨.")
			    time.Sleep(20 * time.Second)
			    config.Stdlog.Println("Oshot SMS process 종료 완료")
			    return
			default:	
			
				for _, tableName := range oshotTable {
					var t = time.Now()
		
					if t.Day() < 3 {
						wg.Add(1)
						go smsProcess(&wg, tableName[0], true)
					}
		
					wg.Add(1)
					go smsProcess(&wg, tableName[0], false)
				}
		
				wg.Wait()
			}
	}

}

func smsProcess(wg *sync.WaitGroup, ostable string, nsFlag bool) {

	defer wg.Done()
	var db = databasepool.DB
	var errlog = config.Stdlog

	var isProc = true
	var t = time.Now()
	if nsFlag {
		t = t.Add(time.Hour * -96)
	}
	var monthStr = fmt.Sprintf("%d%02d", t.Year(), t.Month())

	var SMSTable = ostable + "SMS_" + monthStr

	//db.Exec("UPDATE OShotSMS SET SendDT=now(), SendResult='6', Telecom='000' WHERE SendResult=1 and date_add(insertdt, interval 6 HOUR) < now()")
	//db.Exec("insert into " + SMSTable + " SELECT * FROM OShotSMS WHERE SendResult>1 AND SendDT is not null and telecom = '000'")
	//db.Exec("delete FROM OShotSMS WHERE SendResult>1 AND SendDT is not null and telecom = '000'")

	//발송 6시간 지난 메세지는 응답과 상관 없이 성공 처리 함.

	var groupQuery = "select cb_msg_id, SendResult, SendDT, MsgID, telecom,userid from " + SMSTable + " a where a.proc_flag = 'Y' "

	groupRows, err := db.Query(groupQuery)
	if err != nil {
		errlog.Println("webcsms.go / smsProcess / 스마트미 SMS 조회 중 오류 발생")
		errcode := err.Error()

		if s.Index(errcode, "relation") > 0 {
			db.Exec("Create Table IF NOT EXISTS " + SMSTable + "(LIKE "+ostable+"sms INCLUDING ALL)")
			errlog.Println("oshot "+SMSTable + " 생성 !!")

		}

		isProc = false
		return
	}
	defer groupRows.Close()

	if isProc {

		for groupRows.Next() {
			var cb_msg_id, sendresult, senddt, msgid, telecom, userid sql.NullString

			groupRows.Scan(&cb_msg_id, &sendresult, &senddt, &msgid, &telecom, &userid)

			tr_net := "ETC"

			if s.EqualFold(telecom.String, "011") {
				tr_net = "SKT"
			} else if s.EqualFold(telecom.String, "016") {
				tr_net = "KTF"
			} else if s.EqualFold(telecom.String, "019") {
				tr_net = "LGT"
			}

			if !s.EqualFold(sendresult.String, "6") {

				numcode, _ := strconv.Atoi(sendresult.String)
				var errcode = fmt.Sprintf("%d%03d", 7, numcode)

				val, exists := errmsg[sendresult.String]
				if !exists {
					val = "기타 오류"
				}

				db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.message_type = 'PH', dr.code = '" + errcode + "', dr.message = concat(dr.message, '," + val + "'), dr.remark1 = '" + telecom.String + "', dr.remark2 = '" + senddt.String + "' where  userid='" + userid.String + "' and msgid = '" + cb_msg_id.String + "'")
			} else {
				db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.message_type = 'PH', dr.code = '0000', dr.message = '', dr.remark1 = '" + tr_net + "', dr.remark2 = '" + senddt.String + "' where  userid='" + userid.String + "' and msgid = '" + cb_msg_id.String + "'")
			}

			db.Exec("update " + SMSTable + " set proc_flag = 'N' where msgid = '" + msgid.String + "'")
		}
	}
}
