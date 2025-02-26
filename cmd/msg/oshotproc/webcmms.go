package oshotproc

import (
	"fmt"
	"sync"
	"time"
	"strconv"
	"context"
	s "strings"
	"database/sql"

	config "mycs/configs"
	databasepool "mycs/configs/databasepool"

	_ "github.com/go-sql-driver/mysql"
)

func LMSProcess(ctx context.Context) {
	var wg sync.WaitGroup

	var db = databasepool.DB
	var errlog = config.Stdlog
	var oshotTable [][]string
	var otable sql.NullString

	var OshotQuery = "select distinct a.oshot from DHN_CLIENT_LIST a where a.use_flag = 'Y' and ifnull(a.dest,'OSHOT') = 'OSHOT' and LENGTH(a.oshot) > 1 and a.oshot  is not null"

	OshotTable, err := db.Query(OshotQuery)

	if err != nil {
		errlog.Println("DHN CLIENT LIST 조회 오류 ")
		time.Sleep(10 * time.Second)
	}
	defer OshotTable.Close()

	for OshotTable.Next() {
		OshotTable.Scan(&otable)
		oshotTable = append(oshotTable, []string{otable.String})
	}
	errlog.Println("Oshot MMS length : ", len(oshotTable))
	for {
			select {
				case <- ctx.Done():
			
			    config.Stdlog.Println("Oshot MMS process가 10초 후에 종료 됨.")
			    time.Sleep(10 * time.Second)
			    config.Stdlog.Println("Oshot MMS process 종료 완료")
			    return
			default:	
			
				for _, tableName := range oshotTable {
					var t = time.Now()
		
					if t.Day() < 3 {
						wg.Add(1)
						go pre_mmsProcess(&wg, tableName[0])
					}
		
					wg.Add(1)
					go mmsProcess(&wg, tableName[0])
				}
				wg.Wait()
			}
	}

}

func mmsProcess(wg *sync.WaitGroup, ostable string) {
	defer wg.Done()
	defer func(){
		if r := recover(); r != nil {
			config.Stdlog.Println("OSHOT mmsProcess panic 발생 원인 : ", r)
			if err, ok := r.(error); ok {
				if s.Contains(err.Error(), "connection refused") {
					for {
						config.Stdlog.Println("OSHOT mmsProcess send ping to DB")
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

	var MMSTable = ostable + "MMS_" + monthStr

	//lms 성공 처리
	//db.Exec("UPDATE OShotMMS SET SendDT=now(), SendResult='6', Telecom='000' WHERE SendResult=1 and date_add(insertdt, interval 6 HOUR) < now()")
	//db.Exec("insert into " + MMSTable + " SELECT * FROM OShotMMS WHERE SendResult>1 AND SendDT is not null and telecom = '000'")
	//db.Exec("delete FROM OShotMMS WHERE SendResult>1 AND SendDT is not null and telecom = '000'")

	//발송 6시간 지난 메세지는 응답과 상관 없이 성공 처리 함.

	var groupQuery = "select cb_msg_id, SendResult, File_Path1,SendDT, MsgID, telecom, userid  from " + MMSTable + " a where a.proc_flag = 'Y' "

	groupRows, err := db.Query(groupQuery)
	if err != nil {
		//errlog.Println("스마트미 MMS 조회 중 오류 발생", groupQuery)
		errcode := err.Error()
		errlog.Println("스마트미 MMS 조회 중 오류 발생", groupQuery, errcode)

		if s.Index(errcode, "1146") > 0 {
			db.Exec("Create Table IF NOT EXISTS " + MMSTable + " like " + ostable + "MMS")
			errlog.Println(MMSTable + " 생성 !!")
		}
		time.Sleep(10 * time.Second)
		isProc = false
		return
	}
	defer groupRows.Close()

	if isProc {

		for groupRows.Next() {
			var cb_msg_id, sendresult, file_path1, senddt, msgid, telecom, userid sql.NullString

			groupRows.Scan(&cb_msg_id, &sendresult, &file_path1, &senddt, &msgid, &telecom, &userid)

			tr_net := "ETC"

			if s.EqualFold(telecom.String, "011") {
				tr_net = "SKT"
			} else if s.EqualFold(telecom.String, "016") {
				tr_net = "KTF"
			} else if s.EqualFold(telecom.String, "019") {
				tr_net = "LGT"
			}
			/*
				var msg_type = "LMS"

				if len(file_path1.String) > 1 {
					msg_type = "MMS"
				}
			*/

			if !s.EqualFold(sendresult.String, "6") {
				numcode, _ := strconv.Atoi(sendresult.String)
				var errcode = fmt.Sprintf("%d%03d", 7, numcode)

				val, exists := errmsg[sendresult.String]
				if !exists {
					val = "기타 오류"
				}

				db.Exec("update DHN_RESULT dr set dr.message_type = 'PH', dr.result = 'Y', dr.code = '" + errcode + "', dr.message = concat(dr.message, '," + val + "'), dr.remark1 = '" + telecom.String + "', dr.remark2 = '" + senddt.String + "' where userid='" + userid.String + "' and msgid = '" + cb_msg_id.String + "'")
			} else {
				db.Exec("update DHN_RESULT dr set dr.message_type = 'PH', dr.result = 'Y', dr.code = '0000', dr.message = '', dr.remark1 = '" + tr_net + "', dr.remark2 = '" + senddt.String + "' where userid='" + userid.String + "' and  msgid = '" + cb_msg_id.String + "'")
			}

			db.Exec("update " + MMSTable + " set proc_flag = 'N' where msgid = '" + msgid.String + "'")
		}
	}

}

func pre_mmsProcess(wg *sync.WaitGroup, ostable string) {
	defer wg.Done()
	defer func(){
		if r := recover(); r != nil {
			config.Stdlog.Println("OSHOT mmsProcess panic 발생 원인 : ", r)
			if err, ok := r.(error); ok {
				if s.Contains(err.Error(), "connection refused") {
					for {
						config.Stdlog.Println("OSHOT mmsProcess send ping to DB")
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
	var t = time.Now().Add(time.Hour * -96)
	var monthStr = fmt.Sprintf("%d%02d", t.Year(), t.Month())

	var MMSTable = ostable + "MMS_" + monthStr

	//lms 성공 처리
	//db.Exec("UPDATE OShotMMS SET SendDT=now(), SendResult='6', Telecom='000' WHERE SendResult=1 and date_add(insertdt, interval 6 HOUR) < now()")
	//db.Exec("insert into " + MMSTable + " SELECT * FROM OShotMMS WHERE SendResult>1 AND SendDT is not null and telecom = '000'")
	//db.Exec("delete FROM OShotMMS WHERE SendResult>1 AND SendDT is not null and telecom = '000'")

	//발송 6시간 지난 메세지는 응답과 상관 없이 성공 처리 함.

	var groupQuery = "select cb_msg_id, SendResult, File_Path1,SendDT, MsgID, telecom,userid  from " + MMSTable + " a where a.proc_flag = 'Y' "

	groupRows, err := db.Query(groupQuery)
	if err != nil {
		//errlog.Println("스마트미 MMS 조회 중 오류 발생")
		errcode := err.Error()
		errlog.Println("스마트미 MMS 조회 중 오류 발생", groupQuery, errcode)

		if s.Index(errcode, "1146") > 0 {
			db.Exec("Create Table IF NOT EXISTS " + MMSTable + " like " + ostable + "MMS")
			errlog.Println(MMSTable + " 생성 !!")
		}
		time.Sleep(10 * time.Second)
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

			if !s.EqualFold(sendresult.String, "6") {
				numcode, _ := strconv.Atoi(sendresult.String)
				var errcode = fmt.Sprintf("%d%03d", 7, numcode)

				val, exists := errmsg[sendresult.String]
				if !exists {
					val = "기타 오류"
				}

				db.Exec("update DHN_RESULT dr set dr.message_type = 'PH', dr.result = 'Y', dr.code = '" + errcode + "', dr.message = concat(dr.message, '," + val + "'), dr.remark1 = '" + telecom.String + "', dr.remark2 = '" + senddt.String + "' where userid='" + userid.String + "' and  msgid = '" + cb_msg_id.String + "'")
			} else {
				db.Exec("update DHN_RESULT dr set dr.message_type = 'PH', dr.result = 'Y', dr.code = '0000', dr.message = '', dr.remark1 = '" + telecom.String + "', dr.remark2 = '" + senddt.String + "' where  userid='" + userid.String + "' and msgid = '" + cb_msg_id.String + "'")
			}

			db.Exec("update " + MMSTable + " set proc_flag = 'N' where msgid = '" + msgid.String + "'")
		}
	}

}
