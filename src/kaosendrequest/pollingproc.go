package kaosendrequest

import (
	config "mycs/src/kaoconfig"
	databasepool "mycs/src/kaodatabasepool"
	s "strings"
	"sync"
	"database/sql"
	"context"
	"time"
)

func ResultProc(ctx context.Context) {
	var wg sync.WaitGroup

	for {
			select {
				case <- ctx.Done():
			
			    config.Stdlog.Println("Polling Result process가 20초 후에 종료 됨.")
			    time.Sleep(20 * time.Second)
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
 
	var db = databasepool.DB
	var errlog = config.Stdlog

	pollingsql := `
	SELECT dpr.msg_id, dpr.type 
	FROM DHN_POLLING_RESULT dpr 
	INNER JOIN DHN_RESULT dr ON dpr.msg_id = dr.msgid 
	WHERE dr.result = 'N' AND dr.sync = 'N'`

	resrows, err := db.Query(pollingsql)
	if err != nil {
		errlog.Fatal(err)
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
			insDelResData(supmsgids, "S")
			supmsgids = nil
		}

		if len(fupmsgids) >= 1000 {
			insDelResData(supmsgids, "F")
			fupmsgids = nil
		}
	
	}

	if len(supmsgids) > 0 {
		insDelResData(supmsgids, "S")
	}
	
	if len(fupmsgids) > 0 {
		insDelResData(supmsgids, "F")
	}
	  
}

func insDelResData(idValues []interface{}, res string) {
	tx, err := databasepool.DB.Begin()
	if err != nil {
		confiag.Stdlog.Println("polling_proc.go / getPollingProcess / dhn_result / 트랜젝션 초기화 실패 ", err)
	}
	defer tx.Rollback()

	var stmt *sql.Stmt
	if res == "S" {
		stmt, err = tx.Prepare("update DHN_RESULT set result = 'Y' where sync = 'N' and msgid = $1")
	} else {
		stmt, err = tx.Prepare("update DHN_RESULT set result = 'Y',code = '9999', message = 'ME09' where sync = 'N' and msgid = $1")
	}

	if err != nil {
		confiag.Stdlog.Println("polling_proc.go / getPollingProcess / dhn_result / stmt insert 초기화 실패 ", err)
		return
	}

	stmt2, err := tx.Prepare("delete from DHN_POLLING_RESULT where msg_id = $1")
	
	if err != nil {
		confiag.Stdlog.Println("polling_proc.go / getPollingProcess / dhn_result / delete stmt 초기화 실패 ", err)
		stmt.Close()
		return
	}

	for _, data := range idValues {
		_, err = stmt.Exec(data)
		if err != nil {
			confiag.Stdlog.Println("polling_proc.go / getPollingProcess / dhn_result / stmt insert Exec ", err)
		}
		_, err = stmt2.Exec(data)
		if err != nil {
			confiag.Stdlog.Println("polling_proc.go / getPollingProcess / dhn_result / stmt delete Exec ", err)
		}
	}

	stmt.Close()
	stmt2.Close()

	err = tx.Commit()
	if err != nil {
		confiag.Stdlog.Println("polling_proc.go / getPollingProcess / dhn_result / ftStmt commit ", err)
	}

}
