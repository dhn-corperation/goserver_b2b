package kaocommon

import(
	"database/sql"
)

type AtReqColumn struct {
	Msgid			string	`db:"msgid"`
	Userid			string	`db:"userid"`
	Ad_flag			string	`db:"ad_flag"`
	Button1			string	`db:"button1"`
	Button2			string	`db:"button2"`
	Button3			string	`db:"button3"`
	Button4			string	`db:"button4"`
	Button5			string	`db:"button5"`
	Image_link		string	`db:"image_link"`
	Image_url		string	`db:"image_url"`
	Message_type    string	`db:"message_type"`
	Msg    			string	`db:"msg"`
	Msg_sms    		string	`db:"msg_sms"`
	Only_sms    	string	`db:"only_sms"`
	Phn    			string	`db:"phn"`
	Profile    		string	`db:"profile"`
	P_com    		string	`db:"p_com"`
	P_invoice    	string	`db:"p_invoice"`
	Reg_dt   		string	`db:"reg_dt"`
	Remark1   		string	`db:"remark1"`
	Remark2    		string	`db:"remark2"`
	Remark3   		string	`db:"remark3"`
	Remark4    		string	`db:"remark4"`
	Remark5    		string	`db:"remark5"`
	Reserve_dt    	string	`db:"reserve_dt"`
	Sms_kind    	string	`db:"sms_kind"`
	Sms_lms_tit     string	`db:"sms_lms_tit"`
	Sms_sender      string	`db:"sms_sender"`
	S_code    		string	`db:"s_code"`
	Tmpl_id    		string	`db:"tmpl_id"`
	Wide    		string	`db:"wide"`
	Send_group    	sql.NullString	`db:"send_group"`
	Supplement    	string	`db:"supplement"`
	Price           sql.NullInt64		`db:"price"`
	Currency_type   string	`db:"currency_type"`
	Title           string  `db:"title"`
}

type FtReqColumn struct {
	Msgid			string	`db:"msgid"`
	Userid			string	`db:"userid"`
	Ad_flag			string	`db:"ad_flag"`
	Button1			string	`db:"button1"`
	Button2			string	`db:"button2"`
	Button3			string	`db:"button3"`
	Button4			string	`db:"button4"`
	Button5			string	`db:"button5"`
	Image_link		string	`db:"image_link"`
	Image_url		string	`db:"image_url"`
	Message_type    string	`db:"message_type"`
	Msg    			string	`db:"msg"`
	Msg_sms    		string	`db:"msg_sms"`
	Only_sms    	string	`db:"only_sms"`
	Phn    			string	`db:"phn"`
	Profile    		string	`db:"profile"`
	P_com    		string	`db:"p_com"`
	P_invoice    	string	`db:"p_invoice"`
	Reg_dt   		string	`db:"reg_dt"`
	Remark1   		string	`db:"remark1"`
	Remark2    		string	`db:"remark2"`
	Remark3   		string	`db:"remark3"`
	Remark4    		string	`db:"remark4"`
	Remark5    		string	`db:"remark5"`
	Reserve_dt    	string	`db:"reserve_dt"`
	Sms_kind    	string	`db:"sms_kind"`
	Sms_lms_tit     string	`db:"sms_lms_tit"`
	Sms_sender      string	`db:"sms_sender"`
	S_code    		string	`db:"s_code"`
	Tmpl_id    		string	`db:"tmpl_id"`
	Wide    		string	`db:"wide"`
	Send_group    	sql.NullString	`db:"send_group"`
	Supplement    	string	`db:"supplement"`
	Price           sql.NullInt64		`db:"price"`
	Currency_type   string	`db:"currency_type"`
	Title           string  `db:"title"`
	Header    		string	`db:"header"`
	Carousel    	string	`db:"carousel"`
	Att_items       string  `db:"att_items"`
	Att_coupon      string  `db:"att_coupon"`
	Attachments     string  `db:"attachments"`
}

type MsgReqColumn struct {
	Msgid			string	`db:"msgid"`
	Userid			string	`db:"userid"`
	Ad_flag			string	`db:"ad_flag"`
	Button1			string	`db:"button1"`
	Button2			string	`db:"button2"`
	Button3			string	`db:"button3"`
	Button4			string	`db:"button4"`
	Button5			string	`db:"button5"`
	Code     		string	`db:"code"`
	Image_link		string	`db:"image_link"`
	Image_url		string	`db:"image_url"`
	Kind            sql.NullString	`db:"kind"`
	Message         string	`db:"message"`
	Message_type    string	`db:"message_type"`
	Msg    			string	`db:"msg"`
	Msg_sms    		string	`db:"msg_sms"`
	Only_sms    	string	`db:"only_sms"`
	Phn    			string	`db:"phn"`
	Profile    		string	`db:"profile"`
	P_com    		string	`db:"p_com"`
	P_invoice    	string	`db:"p_invoice"`
	Reg_dt   		string	`db:"reg_dt"`
	Remark1   		string	`db:"remark1"`
	Remark2    		string	`db:"remark2"`
	Remark3   		string	`db:"remark3"`
	Remark4    		string	`db:"remark4"`
	Remark5    		string	`db:"remark5"`
	Res_dt    		string	`db:"res_dt"`
	Reserve_dt    	string	`db:"reserve_dt"`
	Result    		string	`db:"result"`
	S_code    		string	`db:"s_code"`
	Sms_kind    	string	`db:"sms_kind"`
	Sms_lms_tit     string	`db:"sms_lms_tit"`
	Sms_sender      string	`db:"sms_sender"`
	Sync    		string	`db:"sync"`
	Tmpl_id    		string	`db:"tmpl_id"`
	Wide    		string	`db:"wide"`
	Send_group    	sql.NullString	`db:"send_group"`
	Supplement    	string	`db:"supplement"`
	Price           sql.NullInt64		`db:"price"`
	Currency_type   sql.NullString	`db:"currency_type"`
	Header    		string	`db:"header"`
	Carousel    	string	`db:"carousel"`
}