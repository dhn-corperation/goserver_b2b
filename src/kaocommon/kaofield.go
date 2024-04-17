package kaocommon

import(
	"database/sql"
)

type AtReqColumn struct {
	Msgid			sql.NullString	`db:"msgid"`
	Userid			sql.NullString	`db:"userid"`
	Ad_flag			sql.NullString	`db:"ad_flag"`
	Button1			sql.NullString	`db:"button1"`
	Button2			sql.NullString	`db:"button2"`
	Button3			sql.NullString	`db:"button3"`
	Button4			sql.NullString	`db:"button4"`
	Button5			sql.NullString	`db:"button5"`
	Image_link		sql.NullString	`db:"image_link"`
	Image_url		sql.NullString	`db:"image_url"`
	Message_type    sql.NullString	`db:"message_type"`
	Msg    			sql.NullString	`db:"msg"`
	Msg_sms    		sql.NullString	`db:"msg_sms"`
	Only_sms    	sql.NullString	`db:"only_sms"`
	Phn    			sql.NullString	`db:"phn"`
	Profile    		sql.NullString	`db:"profile"`
	P_com    		sql.NullString	`db:"p_com"`
	P_invoice    	sql.NullString	`db:"p_invoice"`
	Reg_dt   		sql.NullString	`db:"reg_dt"`
	Remark1   		sql.NullString	`db:"remark1"`
	Remark2    		sql.NullString	`db:"remark2"`
	Remark3   		sql.NullString	`db:"remark3"`
	Remark4    		sql.NullString	`db:"remark4"`
	Remark5    		sql.NullString	`db:"remark5"`
	Reserve_dt    	sql.NullString	`db:"reserve_dt"`
	Sms_kind    	sql.NullString	`db:"sms_kind"`
	Sms_lms_tit     sql.NullString	`db:"sms_lms_tit"`
	Sms_sender      sql.NullString	`db:"sms_sender"`
	S_code    		sql.NullString	`db:"s_code"`
	Tmpl_id    		sql.NullString	`db:"tmpl_id"`
	Wide    		sql.NullString	`db:"wide"`
	Send_group    	sql.NullString	`db:"send_group"`
	Supplement    	sql.NullString	`db:"supplement"`
	Price           sql.NullInt64	`db:"price"`
	Currency_type   sql.NullString	`db:"currency_type"`
	Title           sql.NullString  `db:"title"`
}

type FtReqColumn struct {
	Msgid			sql.NullString	`db:"msgid"`
	Userid			sql.NullString	`db:"userid"`
	Ad_flag			sql.NullString	`db:"ad_flag"`
	Button1			sql.NullString	`db:"button1"`
	Button2			sql.NullString	`db:"button2"`
	Button3			sql.NullString	`db:"button3"`
	Button4			sql.NullString	`db:"button4"`
	Button5			sql.NullString	`db:"button5"`
	Image_link		sql.NullString	`db:"image_link"`
	Image_url		sql.NullString	`db:"image_url"`
	Message_type    sql.NullString	`db:"message_type"`
	Msg    			sql.NullString	`db:"msg"`
	Msg_sms    		sql.NullString	`db:"msg_sms"`
	Only_sms    	sql.NullString	`db:"only_sms"`
	Phn    			sql.NullString	`db:"phn"`
	Profile    		sql.NullString	`db:"profile"`
	P_com    		sql.NullString	`db:"p_com"`
	P_invoice    	sql.NullString	`db:"p_invoice"`
	Reg_dt   		sql.NullString	`db:"reg_dt"`
	Remark1   		sql.NullString	`db:"remark1"`
	Remark2    		sql.NullString	`db:"remark2"`
	Remark3   		sql.NullString	`db:"remark3"`
	Remark4    		sql.NullString	`db:"remark4"`
	Remark5    		sql.NullString	`db:"remark5"`
	Reserve_dt    	sql.NullString	`db:"reserve_dt"`
	Sms_kind    	sql.NullString	`db:"sms_kind"`
	Sms_lms_tit     sql.NullString	`db:"sms_lms_tit"`
	Sms_sender      sql.NullString	`db:"sms_sender"`
	S_code    		sql.NullString	`db:"s_code"`
	Tmpl_id    		sql.NullString	`db:"tmpl_id"`
	Wide    		sql.NullString	`db:"wide"`
	Send_group    	sql.NullString	`db:"send_group"`
	Supplement    	sql.NullString	`db:"supplement"`
	Price           sql.NullInt64		`db:"price"`
	Currency_type   sql.NullString	`db:"currency_type"`
	Title           sql.NullString  `db:"title"`
	Header    		sql.NullString	`db:"header"`
	Carousel    	sql.NullString	`db:"carousel"`
	Att_items       sql.NullString  `db:"att_items"`
	Att_coupon      sql.NullString  `db:"att_coupon"`
	Attachments     sql.NullString  `db:"attachments"`
}

type MsgReqColumn struct {
	Msgid			sql.NullString	`db:"msgid"`
	Userid			sql.NullString	`db:"userid"`
	Ad_flag			sql.NullString	`db:"ad_flag"`
	Button1			sql.NullString	`db:"button1"`
	Button2			sql.NullString	`db:"button2"`
	Button3			sql.NullString	`db:"button3"`
	Button4			sql.NullString	`db:"button4"`
	Button5			sql.NullString	`db:"button5"`
	Code     		sql.NullString	`db:"code"`
	Image_link		sql.NullString	`db:"image_link"`
	Image_url		sql.NullString	`db:"image_url"`
	Kind            sql.NullString	`db:"kind"`
	Message         sql.NullString	`db:"message"`
	Message_type    sql.NullString	`db:"message_type"`
	Msg    			sql.NullString	`db:"msg"`
	Msg_sms    		sql.NullString	`db:"msg_sms"`
	Only_sms    	sql.NullString	`db:"only_sms"`
	Phn    			sql.NullString	`db:"phn"`
	Profile    		sql.NullString	`db:"profile"`
	P_com    		sql.NullString	`db:"p_com"`
	P_invoice    	sql.NullString	`db:"p_invoice"`
	Reg_dt   		sql.NullString	`db:"reg_dt"`
	Remark1   		sql.NullString	`db:"remark1"`
	Remark2    		sql.NullString	`db:"remark2"`
	Remark3   		sql.NullString	`db:"remark3"`
	Remark4    		sql.NullString	`db:"remark4"`
	Remark5    		sql.NullString	`db:"remark5"`
	Res_dt    		sql.NullString	`db:"res_dt"`
	Reserve_dt    	sql.NullString	`db:"reserve_dt"`
	Result    		sql.NullString	`db:"result"`
	S_code    		sql.NullString	`db:"s_code"`
	Sms_kind    	sql.NullString	`db:"sms_kind"`
	Sms_lms_tit     sql.NullString	`db:"sms_lms_tit"`
	Sms_sender      sql.NullString	`db:"sms_sender"`
	Sync    		sql.NullString	`db:"sync"`
	Tmpl_id    		sql.NullString	`db:"tmpl_id"`
	Wide    		sql.NullString	`db:"wide"`
	Send_group    	sql.NullString	`db:"send_group"`
	Supplement    	sql.NullString	`db:"supplement"`
	Price           sql.NullInt64		`db:"price"`
	Currency_type   sql.NullString	`db:"currency_type"`
	Header    		sql.NullString	`db:"header"`
	Carousel    	sql.NullString	`db:"carousel"`
}