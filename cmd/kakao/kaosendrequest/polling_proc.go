package kaosendrequest

import (
	"time"
	"sync"
	"context"
	s "strings"
	"database/sql"
	
	config "mycs/configs"
	databasepool "mycs/configs/databasepool"
)

func ResultProc(ctx context.Context) {
	var wg sync.WaitGroup

	for {
			select {
				case <- ctx.Done():
			
			    config.Stdlog.Println("Polling Result process가 10초 후에 종료 됨.")
			    time.Sleep(10 * time.Second)
			    config.Stdlog.Println("Polling Result process 종료 완료")
			    return
			default:	
			
				wg.Add(1)
		
				go resPollingProcess(&wg)
		
				wg.Wait()
			}
	}
}

func resPollingProcess(wg *sync.WaitGroup) {

	defer wg.Done()
 	defer func(){
		if r := recover(); r != nil {
			config.Stdlog.Println("resPollingProcess panic 발생 원인 : ", r)
			if err, ok := r.(error); ok {
				if s.Contains(err.Error(), "connection refused") {
					for {
						config.Stdlog.Println("resPollingProcess send ping to DB")
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

	pollingsql := "SELECT dpr.msg_id, dpr.type FROM DHN_POLLING_RESULT dpr INNER JOIN DHN_RESULT dr ON dpr.msg_id = dr.msgid WHERE dr.result = 'N' AND dr.sync = 'N'"

	resrows, err := db.Query(pollingsql)
	if err != nil {
		errlog.Println("resPollingProcess 쿼리 에러 query : ", pollingsql)
		errlog.Println("resPollingProcess 쿼리 에러 : ", err)
		panic(err)
	}

	supmsgids := []interface{}{}
	fupmsgids := []interface{}{}
	
	for resrows.Next() {
		var msg_id sql.NullString
		var restype sql.NullString

		resrows.Scan(&msg_id, &restype)
		
		if s.EqualFold(restype.String, "S") {
			supmsgids = append(supmsgids, msg_id.String)
		} else {
			fupmsgids = append(fupmsgids, msg_id.String)
		}
		
	
		if len(supmsgids) >= 1000 {

			var commastr = "update DHN_RESULT set result = 'Y' where sync = 'N' and msgid in ("

			for i := 1; i < len(supmsgids); i++ {
				commastr = commastr + "?,"
			}

			commastr = commastr + "?)"

			_, err1 := db.Exec(commastr, supmsgids...)

			if err1 != nil {
				errlog.Println("Result Table S Update 처리 중 오류 발생 ")
			}
			
			commastr = "delete from DHN_POLLING_RESULT where msg_id in ("
			for i := 1; i < len(supmsgids); i++ {
				commastr = commastr + "?,"
			}
			commastr = commastr + "?)"
			_, err1 = db.Exec(commastr, supmsgids...)
			if err1 != nil {
				errlog.Println("Result Table S Delete 처리 중 오류 발생 ")
			}
			
			supmsgids = nil
		}

		if len(fupmsgids) >= 1000 {

			var commastr = "update DHN_RESULT set result = 'Y',code = '9999', message = 'ME09' where sync = 'N' and msgid in ("
			for i := 1; i < len(fupmsgids); i++ {
				commastr = commastr + "?,"
			}
			commastr = commastr + "?)"
			_, err1 := db.Exec(commastr, fupmsgids...)
			if err1 != nil {
				errlog.Println("Result Table F Update 처리 중 오류 발생 ")
			}

			commastr = "delete from DHN_POLLING_RESULT where msg_id in ("
			for i := 1; i < len(fupmsgids); i++ {
				commastr = commastr + "?,"
			}
			commastr = commastr + "?)"
			_, err1 = db.Exec(commastr, fupmsgids...)
			if err1 != nil {
				errlog.Println("Result Table F Delete 처리 중 오류 발생 ")
			}
			fupmsgids = nil
		}
	
	}

	if len(supmsgids) > 0 {

		var commastr = "update DHN_RESULT set result = 'Y' where sync = 'N' and msgid in ("
	
		for i := 1; i < len(supmsgids); i++ {
			commastr = commastr + "?,"
		}
	
		commastr = commastr + "?)"
	
		_, err1 := db.Exec(commastr, supmsgids...)
	
		if err1 != nil {
			errlog.Println("Result Table S Update 처리 중 오류 발생 ")
		}
		
		commastr = "delete from DHN_POLLING_RESULT where msg_id in ("
		for i := 1; i < len(supmsgids); i++ {
			commastr = commastr + "?,"
		}
		commastr = commastr + "?)"
		_, err1 = db.Exec(commastr, supmsgids...)
		if err1 != nil {
			errlog.Println("Result Table S Delete 처리 중 오류 발생 ")
		}
 
		supmsgids = nil
	}
	
	if len(fupmsgids) > 0 {
	
		var commastr = "update DHN_RESULT set result = 'Y',code = '9999', message = 'ME09' where sync = 'N' and msgid in ("
	
		for i := 1; i < len(fupmsgids); i++ {
			commastr = commastr + "?,"
		}
	
		commastr = commastr + "?)"
	
		_, err1 := db.Exec(commastr, fupmsgids...)
	
		if err1 != nil {
			errlog.Println("Result Table F Update 처리 중 오류 발생 ")
		}
	
		commastr = "delete from DHN_POLLING_RESULT where msg_id in ("
		for i := 1; i < len(fupmsgids); i++ {
			commastr = commastr + "?,"
		}
		commastr = commastr + "?)"
		_, err1 = db.Exec(commastr, fupmsgids...)
		if err1 != nil {
			errlog.Println("Result Table F Delete 처리 중 오류 발생 ")
		}
		fupmsgids = nil 
	}
	  
}
