package ktproc

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	config "mycs/src/kaoconfig"
	databasepool "mycs/src/kaodatabasepool"
	"strconv"

	//	"log"
	// s "strings"
	"sync"
	"time"

	"context"

	_ "github.com/go-sql-driver/mysql"
)

func SMSProcess(ctx context.Context, acc int) {
	var wg sync.WaitGroup

	for {
		select {
		case <-ctx.Done():

			config.Stdlog.Println("Oshot SMS process가 10초 후에 종료 됨.")
			time.Sleep(10 * time.Second)
			config.Stdlog.Println("Oshot SMS process 종료 완료")
			return
		default:

			for i := 1; i < 6; i++ {
				wg.Add(1)
				go smsProcess(&wg, "KT_SMS", i, acc)
			}

			wg.Wait()
		}
	}

}

func smsProcess(wg *sync.WaitGroup, table string, seq int, acc int) {

	defer wg.Done()
	var db = databasepool.DB
	var errlog = config.Stdlog

	defer func() {
		if err := recover(); err != nil {
			errlog.Println("KT크로샷 SMS 결과 처리 중 패닉 발생 : ", err)
		}
	}()

	var isProc = true
	var t time.Time
	t = time.Now()

	var monthStr = fmt.Sprintf("%d%02d", t.Year(), t.Month())

	var SMSTable = table + "_" + monthStr
	
	exists, err := checkTableExists(db, SMSTable)
	if err != nil {
		errlog.Println("KT MMS LOG table 존재유무 조회 오류 err : ", err)
	}
	if !exists {
		db.Exec("Create Table IF NOT EXISTS " + SMSTable + " like " + table)
		errlog.Println(SMSTable + " 생성 !!")
	}

	var searchQuery = "select userid, msgid, resp_JobID, resp_SubmitTime from " + table + " where sep_seq = " + strconv.Itoa(seq)

	searchData, err := db.Query(searchQuery)
	if err != nil {
		errcode := err.Error()
		errlog.Println("KT크로샷 SMS 조회 중 오류 발생", searchQuery, errcode)
		return
	}
	defer searchData.Close()

	if isProc {
		acc := account[acc]
		client := NewMessage(acc["apiKey"], acc["apiPw"], acc["userKey"], false, 3)
		for searchData.Next() {
			var userid, msgid, resp_SubmitTime sql.NullString
			var resp_JobID sql.NullInt64

			searchData.Scan(&userid, &msgid, &resp_JobID, &resp_SubmitTime)

			var st string

			if len(resp_SubmitTime.String) >= 8 {
				st = resp_SubmitTime.String[:8]
			} else {
				continue
			}

			sendData := SearchReqTable{
				JobIDs: []int64{
					resp_JobID.Int64,
				},
				SendDay: st,
			}

			resp, err := client.SearchResult("/inquiry/report/", sendData)

			if err != nil {
				errlog.Println(userid.String, "- msgid : ", msgid.String, " KT크로샷 결과조회 API 발송 중 오류 발생 : err : ", err)
				continue
			}

			if resp.StatusCode != 200 {
				// errlog.Println(userid.String, "- msgid : ", msgid.String, " KT크로샷 결과조회 API 발송 중 오류 발생 : statusCode : ", resp.StatusCode)
				continue
			}

			body, _ := ioutil.ReadAll(resp.Body)
			var decodeBody SearchResTable

			err = json.Unmarshal([]byte(body), &decodeBody)
			if err != nil {
				errlog.Println(userid.String, "- msgid : ", msgid.String, " KT크로샷 결과조 API 결과 변환중 에러 발생 : ", err)
				continue
			}
			first := decodeBody.JobIDs[0]

			if first.Result == 0 {
				continue
			}

			convResult := strconv.Itoa(first.Result)
			resultCode := KTCode(convResult)
			resultMessage := KTCodeMessage(resultCode)

			var telInfo = "ETC"
			var telInfoLog = 0
			if first.TelcoInfo != nil {
				if *first.TelcoInfo == 1 {
					telInfo = "SKT"
				} else if *first.TelcoInfo == 2 {
					telInfo = "KTF"
				} else if *first.TelcoInfo == 3 {
					telInfo = "LGT"
				}
				telInfoLog = *first.TelcoInfo
			}

			parsedTime, err := time.Parse("20060102150405", first.Time)
			if err != nil {
				errlog.Println(userid.String, "- msgid : ", msgid.String, " KT크로샷 결과조회 API 발송 중 시간변환 오류 발생 : ", err, "  /  statusCode : ", first.Time)
				continue
			}

			formattedTime := parsedTime.Format("2006-01-02 15:04:05")

			_, err = db.Exec(`insert into ` + SMSTable + `(userid, msgid, MessageSubType, CallbackNumber, Bundle_Seq, Bundle_Num, Bundle_Content, resp_JobID, resp_Time, resp_SubmitTime, resp_Result, Resp_TelconInfo, resp_EndUserID, resp_ServiceProviderID, sep_seq, dhn_id)
				select userid, msgid, MessageSubType, CallbackNumber, Bundle_Seq, Bundle_Num, Bundle_Content, resp_JobID, '` + first.Time + `', '` + first.SubmitTime + `', '` + strconv.Itoa(first.Result) + `', '` + strconv.Itoa(telInfoLog) + `', '` + first.EndUserID + `', '` + first.ServiceProviderID + `', sep_seq, dhn_id
				from KT_SMS
				WHERE userid = '` + userid.String + `' and msgid = '` + msgid.String + `'`)
			if err != nil {
				errlog.Println(userid.String, "- msgid : ", msgid.String, " KT크로샷 결과조 API 결과 LOG 테이블 입력중 에러 발생 : ", err)
				continue
			}
			if first.Result != 10000 {
				_, err = db.Exec("update DHN_RESULT set message_type = 'PH', result = 'Y', code = '" + resultCode + "', message = concat(message, '," + resultMessage + "'), remark1 = '" + telInfo + "', remark2 = '" + formattedTime + "' where userid='" + userid.String + "' and msgid = '" + msgid.String + "'")
				if err != nil {
					errlog.Println(userid.String, "- msgid : ", msgid.String, " KT크로샷 결과조 API 결과 DHN_RESULT 테이블 반영 실패1 : ", err)
				}
			} else {
				_, err = db.Exec("update DHN_RESULT set message_type = 'PH', result = 'Y', code = '0000', message = '', remark1 = '" + telInfo + "', remark2 = '" + formattedTime + "' where userid='" + userid.String + "' and msgid = '" + msgid.String + "'")
				if err != nil {
					errlog.Println(userid.String, "- msgid : ", msgid.String, " KT크로샷 결과조 API 결과 DHN_RESULT 테이블 반영 실패2 : ", err)
				}
			}

			db.Exec(`delete from KT_SMS where userid = '` + userid.String + `' and msgid = '` + msgid.String + `'`)

		}
	}
}
