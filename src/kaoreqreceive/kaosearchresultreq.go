package kaoreqreceive

import (
	"fmt"
	"time"
	"database/sql"
	s "strings"

	"mycs/src/kaoresulttable"
	cm "mycs/src/kaocommon"
	config "mycs/src/kaoconfig"
	databasepool "mycs/src/kaodatabasepool"

	_ "github.com/go-sql-driver/mysql"
	"github.com/valyala/fasthttp"
	"github.com/goccy/go-json"
)


func SearchResultReq(c *fasthttp.RequestCtx) {
	errlog := config.Stdlog
	db := databasepool.DB

	userid := string(c.Request.Header.Peek("userid"))
	userip := c.RemoteIP().String()
	isValidation := false

	sqlstr := `
		select 
			count(1) as cnt 
		from
			DHN_CLIENT_LIST
		where
			user_id = ?
			and ip = ?
			and use_flag = 'Y'`

	var cnt int
	err := db.QueryRow(sqlstr, userid, userip).Scan(&cnt)
	if err != nil { errlog.Println(err) }

	if cnt > 0 { 
		isValidation = true 
	} else {
		errlog.Println("허용되지 않은 사용자 및 아이피에서 발송 요청!! (userid : ", userid, "/ ip : ", userip, ")")
	}

	var startNow = time.Now()
	var startTime = fmt.Sprintf("%02d:%02d:%02d", startNow.Hour(), startNow.Minute(), startNow.Second())

	if isValidation {
		db2json := map[string]string{
			"msgid":         "msgid",
			"userid":        "userid",
			"ad_flag":       "ad_flag",
			"button1":       "button1",
			"button2":       "button2",
			"button3":       "button3",
			"button4":       "button4",
			"button5":       "button5",
			"code":          "code",
			"image_link":    "image_link",
			"image_url":     "image_url",
			"kind":          "kind",
			"message":       "message",
			"message_type":  "message_type",
			"msg":           "msg",
			"msg_sms":       "msg_sms",
			"only_sms":      "only_sms",
			"p_com":         "p_com",
			"p_invoice":     "p_invoice",
			"phn":           "phn",
			"profile":       "profile",
			"reg_dt":        "reg_dt",
			"remark1":       "remark1",
			"remark2":       "remark2",
			"remark3":       "remark3",
			"remark4":       "remark4",
			"remark5":       "remark5",
			"res_dt":        "res_dt",
			"reserve_dt":    "reserve_dt",
			"result":        "result",
			"s_code":        "s_code",
			"sms_kind":      "sms_kind",
			"sms_lms_tit":   "sms_lms_tit",
			"sms_sender":    "sms_sender",
			"sync":          "sync",
			"tmpl_id":       "tmpl_id",
			"wide":          "wide",
			"send_group":    "send_group",
			"supplement":    "supplement",
			"price":         "price",
			"currency_type": "currency_type",
			"title":         "title",
		}
		
		var reqData kaoresulttable.ResultTable
		if err1 := json.Unmarshal(c.PostBody(), &reqData); err1 != nil {
			errlog.Println(userid, " - 발송 결과 재수신 data decoding error : ", err1)
		}
		errlog.Println(userid, " - 발송 결과 재수신 init Msgid 건수 : ", len(reqData.Msgid), " / 수신 시각 : ", startTime)

		if len(reqData.Msgid) == 0 {
			res, _ := json.Marshal(map[string]string{
				"code":    "00",
				"message": "msgid가 존재하지 않습니다.",
			})
			c.SetContentType("application/json")
			c.SetStatusCode(fasthttp.StatusOK)
			c.SetBody(res)
			return
		}

		ids := []string{}

		for i, v := range reqData.Msgid {
			reqData.Msgid[i] = fmt.Sprintf("'%s'", v)
			ids = append(ids, v)
		}

		msgids := s.Join(reqData.Msgid, ", ")

		joinSql := ""
		joinTable := "DHN_RESULT_" + reqData.Regdt

		exists, err2 := checkTableExists(db, joinTable)
		if err2 != nil {
			errlog.Println(userid, " - 발송 결과 재수신 table 존재유무 조회 오류 err : ", err2)
		}

		if exists {
			joinSql = " union all select * from " + joinTable + " where userid = '" + userid + "' and msgid in (" + msgids + ")"
		} else {
			joinSql = " "
		}

		resultSql := "select * from DHN_RESULT where userid = '" + userid + "' and msgid in (" + msgids + ")" + joinSql

		reqrows, err := db.Query(resultSql)
		if err != nil {
			errlog.Fatal(resultSql, err)
		}

		columnTypes, err := reqrows.ColumnTypes()
		if err != nil {
			errlog.Fatal(err)
		}
		
		count := len(columnTypes)
		scanArgs := make([]interface{}, count)
		
		finalRows := []interface{}{}
		upmsgids := []interface{}{}

		var isContinue bool
		
		isFirstRow := true

		errlog.Println(userid, " - 발송 결과 재수신 결과 전송 시작 건수 : ", len(reqData.Msgid))
		for reqrows.Next() {
			
			if isFirstRow {
				for i, v := range columnTypes {
	
					switch v.DatabaseTypeName() {
					case "VARCHAR", "TEXT", "UUID", "TIMESTAMP":
						scanArgs[i] = new(sql.NullString)
						break
					case "BOOL":
						scanArgs[i] = new(sql.NullBool)
						break
					case "INT4":
						scanArgs[i] = new(sql.NullInt64)
						break
					default:
						scanArgs[i] = new(sql.NullString)
					}
				}
				isFirstRow = false				
			}

			err := reqrows.Scan(scanArgs...)
			if err != nil {
				errlog.Fatal(err)
			}

			masterData := map[string]interface{}{}
			mid := ""

			for i, v := range columnTypes {

				isContinue = false

				if z, ok := (scanArgs[i]).(*sql.NullBool); ok {
					masterData[db2json[s.ToLower(v.Name())]] = z.Bool
					isContinue = true
				}

				if z, ok := (scanArgs[i]).(*sql.NullString); ok {
					masterData[db2json[s.ToLower(v.Name())]] = z.String
					isContinue = true
				}

				if z, ok := (scanArgs[i]).(*sql.NullInt64); ok {
					masterData[db2json[s.ToLower(v.Name())]] = z.Int64
					isContinue = true
				}

				if z, ok := (scanArgs[i]).(*sql.NullFloat64); ok {
					masterData[db2json[s.ToLower(v.Name())]] = z.Float64
					isContinue = true
				}

				if z, ok := (scanArgs[i]).(*sql.NullInt32); ok {
					masterData[db2json[s.ToLower(v.Name())]] = z.Int32
					isContinue = true
				}
				if !isContinue {
					masterData[db2json[s.ToLower(v.Name())]] = scanArgs[i]
				}

				if s.EqualFold(v.Name(), "MSGID") {
					upmsgids = append(upmsgids, masterData[db2json[s.ToLower(v.Name())]])
					mid, _ = masterData[db2json[s.ToLower(v.Name())]].(string)
				}
			}

			finalRows = append(finalRows, masterData)
			ids = cm.RemoveValueInPlace(ids, mid)
		}

		errlog.Println(userid, " - 발송 결과 재수신 결과 전송 끝 건수 : ", len(finalRows))

		var commastr = "update DHN_RESULT set sync='Y' where userid = '" + userid + "' and msgid in (" + msgids + ")"
		db.Exec(commastr)

		unmarshalRes := map[string]interface{}{
			"code": "00",
			"message": "성공",
			"data": map[string]interface{}{
				"count": len(finalRows),
				"detail": finalRows,
			},
		}

		if len(ids) > 0 {
			searchIds := []string{}
			for _, v := range ids {
				searchIds = append(searchIds, fmt.Sprintf("'%s'", v))
			}

			searchSql := `
				select msgid
				from DHN_RECEPTION
				where msgid in (` + s.Join(searchIds, ", ") + `)
					and userid = '` + userid + `'
			`
			searchProc := true
			searchRow, err := db.Query(searchSql)
			if err != nil {
				errcode := err.Error()
				errlog.Println(userid, " - DHN_RECEPTION 테이블 조회 중 오류 발생 sql : ", searchSql , " / err : ", errcode)
				searchProc = false
			}
			defer searchRow.Close()

			if searchProc {
				for searchRow.Next() {
					var msgid sql.NullString
					searchRow.Scan(&msgid)
					ids = cm.RemoveValueInPlace(ids, msgid.String)
				}
				unmarshalRes["norec"] = ids
			}
		}

		res, _ := json.Marshal(unmarshalRes)

		c.SetContentType("application/json")
		c.SetStatusCode(fasthttp.StatusOK)
		c.SetBody(res)
	} else {
		res, _ := json.Marshal(map[string]string{
			"code":    "error",
			"message": "허용되지 않은 사용자 입니다 / userid : " + userid + " / ip : " + userip,
		})
		c.SetContentType("application/json")
		c.SetStatusCode(fasthttp.StatusNotAcceptable)
		c.SetBody(res)
	}
}