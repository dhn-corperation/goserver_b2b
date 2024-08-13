package kaoreqreceive

import (
	"fmt"
	config "kaoconfig"
	databasepool "kaodatabasepool"

	"time"
)

func TempCopyProc() {
	var count int
	for {
		cnterr := databasepool.DB.QueryRow(" select count(1) as cnt from DHN_RESULT_TEMP where send_group is null").Scan(&count)
		if cnterr != nil {
			config.Stdlog.Println("DHN_RESULT_TEMP Table - select error : " + cnterr.Error())
		} else {
			if count > 0 {
				var startNow = time.Now()
				var group_no = fmt.Sprintf("%02d%02d%02d%02d%06d", startNow.Day(), startNow.Hour(), startNow.Minute(), startNow.Second(), (startNow.Nanosecond()/1000))
				
				databasepool.DB.Exec("update DHN_RESULT_TEMP set send_group = '" + group_no + "' where   send_group is null limit 1000")
			
				copyQuery := `INSERT INTO DHN_RESULT
select msgid
      ,userid
      ,ad_flag
      ,button1
      ,button2
      ,button3
      ,button4
      ,button5
      ,code
      ,image_link
      ,image_url
      ,kind
      ,message
      ,message_type
      ,msg
      ,msg_sms
      ,only_sms
      ,p_com
      ,p_invoice
      ,phn
      ,profile
      ,reg_dt
      ,remark1
      ,remark2
      ,remark3
      ,remark4
      ,remark5
      ,res_dt
      ,reserve_dt
      ,result
      ,s_code
      ,sms_kind
      ,sms_lms_tit
      ,sms_sender
      ,sync
      ,tmpl_id
      ,wide
      ,null
      ,supplement
      ,price
      ,currency_type
      ,title
      ,header
      ,carousel										
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
