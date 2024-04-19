package kaoreqreceive

import (
	//"encoding/json"
	"database/sql"
	//"fmt"

	_ "github.com/go-sql-driver/mysql"

	config "mycs/src/kaoconfig"
	databasepool "mycs/src/kaodatabasepool"

	//"kaoreqtable"

	"github.com/gin-gonic/gin"

	//"strconv"
	s "strings"
)

func FriendInforeq(c *gin.Context) {

	db := databasepool.DB
	errlog := config.Stdlog
	isValidation, userid, userip := kaocommon.CheckUser(c)
	
	if isValidation {
		sqlstr := "select * from sw_talk_link where partner_send_yn = 'N' and send_user_id = '" + userid + "'"

		reqrows, err := db.Query(sqlstr)
		if err != nil {
			errlog.Fatal(err)
		}

		columnTypes, err := reqrows.ColumnTypes()
		if err != nil {
			errlog.Fatal(err)
		}
		count := len(columnTypes)

		finalRows := []interface{}{}
		upids := []interface{}{}

		var isContinue bool
		for reqrows.Next() {
			scanArgs := make([]interface{}, count)

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

			err := reqrows.Scan(scanArgs...)
			if err != nil {
				errlog.Fatal(err)
			}

			masterData := map[string]interface{}{}

			for i, v := range columnTypes {

				isContinue = false

				if z, ok := (scanArgs[i]).(*sql.NullBool); ok {
					masterData[s.ToLower(v.Name())] = z.Bool
					isContinue = true
				}

				if z, ok := (scanArgs[i]).(*sql.NullString); ok {
					masterData[s.ToLower(v.Name())] = z.String
					isContinue = true
				}

				if z, ok := (scanArgs[i]).(*sql.NullInt64); ok {
					masterData[s.ToLower(v.Name())] = z.Int64
					isContinue = true
				}

				if z, ok := (scanArgs[i]).(*sql.NullFloat64); ok {
					masterData[s.ToLower(v.Name())] = z.Float64
					isContinue = true
				}

				if z, ok := (scanArgs[i]).(*sql.NullInt32); ok {
					masterData[s.ToLower(v.Name())] = z.Int32
					isContinue = true
				}
				if !isContinue {
					masterData[s.ToLower(v.Name())] = scanArgs[i]
				}

				if s.EqualFold(v.Name(), "tl_id") {
					upids = append(upids, masterData[s.ToLower(v.Name())])
				}
			}

			finalRows = append(finalRows, masterData)

			if len(upids) >= 500 {

				var commastr = "update sw_talk_link set partner_send_yn='Y', partner_send_dt = now() where tl_id in ("

				for i := 1; i < len(upids); i++ {
					commastr = commastr + "?,"
				}

				commastr = commastr + "?)"

				_, err1 := db.Exec(commastr, upids...)

				if err1 != nil {
					errlog.Println("Friend Infor Table Update 처리 중 오류 발생 ")
				}

				upids = nil
			}
		}

		if len(upids) > 0 {

				var commastr = "update sw_talk_link set partner_send_yn='Y', partner_send_dt = now() where tl_id in ("

				for i := 1; i < len(upids); i++ {
					commastr = commastr + "?,"
				}

				commastr = commastr + "?)"

				_, err1 := db.Exec(commastr, upids...)

				if err1 != nil {
					errlog.Println("Friend Infor Table Update 처리 중 오류 발생 ")
				}

				upids = nil
		}
		if len(finalRows) > 0 {
			errlog.Println("결과 전송 ( ", userid, " ) : ", len(finalRows))
			//errlog.Println(finalRows)
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
