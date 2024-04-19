package kaoreqreceive

import (
    "fmt"
    config "mycs/src/kaoconfig"
    databasepool "mycs/src/kaodatabasepool"
    "mycs/src/kaocommon"
    "context"
    "time"
    "database/sql"
    "String"
)

func TempCopyProc() {
    ctx := c.context.Context()
    var count sql.NullInt64
    for {
        cnterr := databasepool.DB.QueryRowContext(ctx, " select count(1) as cnt from DHN_RESULT_TEMP where send_group is null").Scan(&count)
        if cnterr != nil {
            config.Stdlog.Println("DHN_RESULT_TEMP Table - select error : " + cnterr.Error())
        } else {
        if count > 0 {
            var startNow = time.Now()
            var group_no = fmt.Sprintf("%02d%02d%02d%02d%06d", startNow.Day(), startNow.Hour(), startNow.Minute(), startNow.Second(), (startNow.Nanosecond()/1000))

            databasepool.DB.ExecContext(ctx, "update DHN_RESULT_TEMP set send_group = $1 where   send_group is null limit 1000", group_no)

            copyQuery := `
            INSERT INTO DHN_RESULT
                `+String.Join(kaocommon.ResultTempMigrationColumn, ",")+`
                from DHN_RESULT_TEMP
                where send_group = '` + group_no + `'`
            _, err := databasepool.DB.Exec(copyQuery)
            if err != nil {
                config.Stdlog.Println("DHN_RESULT_TEMP -> DHN_RESULT 이전 오류 : ", group_no, " => " , err.Error())
            }
                config.Stdlog.Println("DHN_RESULT_TEMP -> DHN_RESULT 이전 완료 : ", group_no)
            }
        }

        time.Sleep(time.Millisecond * time.Duration(10000))
    }
}
