package kaoreqreceive

import (
	//"encoding/json"
	"database/sql"
	//"fmt"

	_ "github.com/go-sql-driver/mysql"

	config "mycs/src/kaoconfig"
	databasepool "mycs/src/kaodatabasepool"
	"mycs/src/kaocommon"

	//"kaoreqtable"

	"github.com/gin-gonic/gin"

	//"strconv"
	s "strings"
)

var db2json = map[string]string{
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

func Resultreq(c *gin.Context) {

	errlog := config.Stdlog
	db := databasepool.DB
	var checkResult kaocommon.CheckUserReturnField = kaocommon.CheckUser(c)
	userid := checkResult.Userid
	userip := checkResult.Userip
	ctx := checkResult.Ctx


	if checkResult.Validation {

		sqlstr := "select * from DHN_RESULT_PROC where userid = $1 and sync='N' and result = 'Y' limit $2"

		reqrows, err := db.QueryContext(ctx, sqlstr, userid, checkResult.SendLimit)
		if err != nil {
			errlog.Println("Resultreq 쿼리 에러 query : ", sqlstr)
			errlog.Println("Resultreq 쿼리 에러 : ", err)
			errlog.Fatal(sqlstr, err)
		}

		columnTypes, err := reqrows.ColumnTypes()
		if err != nil {
			errlog.Println("Resultreq 컬럼 초기화 에러 userid : ", userid)
			errlog.Println("Resultreq 컬럼 초기화 에러 : ", err)
			errlog.Fatal(err)
		}
		
		count := len(columnTypes)
		scanArgs := kaocommon.InitDatabaseColumn(columnTypes, count)
		
		finalRows := []interface{}{}
		upmsgids := []interface{}{}

		var isContinue bool
		
		isFirstRow := true
		
		for reqrows.Next() {
			
			if isFirstRow {
				errlog.Println("결과 전송 ( ", userid, " ) : 시작 " )
				isFirstRow = false				
			}

			err := reqrows.Scan(scanArgs...)
			if err != nil {
				errlog.Fatal(err)
			}

			masterData := map[string]interface{}{}

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
				}
			}

			finalRows = append(finalRows, masterData)

			if len(upmsgids) >= 500 {
				updateResultData(upmsgids, userid)
				upmsgids = nil
			}
		}
		if len(upmsgids) > 0 {
			updateResultData(upmsgids, userid)
		}
		if len(finalRows) > 0 {
			errlog.Println("결과 전송 ( ", userid, " ) : ", len(finalRows))
		}
		c.JSON(200, finalRows)
	} else {
		c.JSON(404, gin.H{
			"code":    "error",
			"message": "사용자 아이디 확인",
			"userid":  userid,
			"ip":      userip,
		})
	}
}

func updateResultData(upmsgids []interface{}, userid string){
	tx, err := databasepool.DB.Begin()
	if err != nil {
		config.Stdlog.Println(err)
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare("update DHN_RESULT_PROC set sync='Y' where userid = $1 and MSGID = $2")
	if err != nil {
		config.Stdlog.Println("stmt 초기화 실패 ", err)
		return
	}
	for _, data := range upmsgids {
		_, err := stmt.Exec(userid, data)
		if err != nil {
			config.Stdlog.Println("stmt personal Exec  ", err)
		}
	}
	_, err = stmt.Exec()
	if err != nil {
		stmt.Close()
		config.Stdlog.Println("stmt Exec ", err)
	}
	stmt.Close()
	err = tx.Commit()
	if err != nil {
		config.Stdlog.Println("stmt commit ", err)
	}
}
