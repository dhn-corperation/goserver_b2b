package structs

import (
	"time"
)

// TABLE
// NAME : DHN_CLIENT_LIST
type ClientList struct {
	User_id string 			`gorm:"size:50;not null"`
	Ip string 				`gorm:"size:50;not null"`
	Use_flag string 		`gorm:"type:char(1);default:'Y'"`
	Send_limit string 		`gorm:"size:50;default:'500'"`
	Sms_len_check string 	`gorm:"type:char(1);default:'N'"`
	Oshot string 			`gorm:"size:50;default:null"`
	Dest string 			`gorm:"size:50;default:null"`
	Alimtalk string 		`gorm:"type:char(1);default:'Y'"`
	Friendtalk string 		`gorm:"type:char(1);default:'N'"`
	Description string 		`gorm:"size:500;default:null"`
}

func (ClientList) TableName() string {
	return "DHN_CLIENT_LIST"
}

// TABLE
// NAME : DHN_RECEPTION
type Reception struct {
	Id uint 				`gorm:"type:bigint;primaryKey;size:20;autoIncrement"`
	Msgid string 			`gorm:"size:20;default:null;index:idx_msgid_userid,priority:2"`
	Userid string 			`gorm:"size:20;default:null;index:idx_msgid_userid,priority:1"`
	Insert_date time.Time 	`gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (Reception) TableName() string {
	return "DHN_RECEPTION"
}

// TABLE
// NAME : DHN_REQUEST
type Request struct {
	Id uint 				`gorm:"type:bigint;primaryKey;size:20;autoIncrement"`
	Msgid string 			`gorm:"size:20;not null;index:idx_msgid"`
	Userid string 			`gorm:"size:20;not null"`
	Ad_flag string 			`gorm:"type:char(1);default:null"`
	Button1 string 			`gorm:"type:text;default:null"`
	Button2 string 			`gorm:"type:text;default:null"`
	Button3 string 			`gorm:"type:text;default:null"`
	Button4 string 			`gorm:"type:text;default:null"`
	Button5 string 			`gorm:"type:text;default:null"`
	Image_link string 		`gorm:"type:text;default:null"`
	Image_url string 		`gorm:"type:text;default:null"`
	Message_type string 	`gorm:"type:char(2);default:null"`
	Msg string 				`gorm:"type:text;not null"`
	Msg_sms string 			`gorm:"type:text;default:null"`
	Only_sms string 		`gorm:"type:char(1);default:null"`
	Phn string 				`gorm:"size:15;not null"`
	Profile string 			`gorm:"size:50;default:null"`
	P_com string 			`gorm:"size:2;default:null"`
	P_invoice string 		`gorm:"size:100;default:null"`
	Reg_dt string 			`gorm:"size:20;not null"`
	Remark1 string 			`gorm:"size:50;default:null"`
	Remark2 string 			`gorm:"size:50;default:null"`
	Remark3 string 			`gorm:"size:50;default:null"`
	Remark4 string 			`gorm:"size:50;default:null"`
	Remark5 string 			`gorm:"size:50;default:null"`
	Reserve_dt string 		`gorm:"size:14;not null;index:idx_sendgroup_reservedt,priority:2"`
	Sms_kind string 		`gorm:"type:char(1);default:null"`
	Sms_lms_tit string 		`gorm:"size:120;default:null"`
	Sms_sender string 		`gorm:"size:15;default:null"`
	S_code string 			`gorm:"size:4;default:null"`
	Tmpl_id string 			`gorm:"size:30;default:null"`
	Wide string 			`gorm:"type:char(1);default:null"`
	Send_group string 		`gorm:"type:char(20);default:null;index:idx_sendgroup_reservedt,priority:1;index:idx_sendgroup"`
	Supplement string 		`gorm:"type:mediumtext;default:null"`
	Price uint 				`gorm:"type:int(11);default:0"`
	Currency_type string 	`gorm:"type:char(3);default:null"`
	Title string 			`gorm:"size:50;default:null"`
	Header string 			`gorm:"size:100;default:null"`
	Carousel string 		`gorm:"type:text;default:null"`
	Attachments string 		`gorm:"type:text;default:null"`
	Att_items string 		`gorm:"type:text;default:null"`
	Att_coupon string 		`gorm:"type:text;default:null"`
	Mms_image_id string 	`gorm:"size:100;default:null"`
}

func (Request) TableName() string {
	return "DHN_REQUEST"
}

// TABLE
// NAME : DHN_REQUEST_AT
type RequestAt struct {
	Id uint 				`gorm:"type:bigint;primaryKey;size:20;autoIncrement"`
	Msgid string 			`gorm:"size:20;not null;index:idx_msgid"`
	Userid string 			`gorm:"size:20;not null"`
	Ad_flag string 			`gorm:"type:char(1);default:null"`
	Button1 string 			`gorm:"type:text;default:null"`
	Button2 string 			`gorm:"type:text;default:null"`
	Button3 string 			`gorm:"type:text;default:null"`
	Button4 string 			`gorm:"type:text;default:null"`
	Button5 string 			`gorm:"type:text;default:null"`
	Image_link string 		`gorm:"type:text;default:null"`
	Image_url string 		`gorm:"type:text;default:null"`
	Message_type string 	`gorm:"type:char(2);default:null"`
	Msg string 				`gorm:"type:text;not null"`
	Msg_sms string 			`gorm:"type:text;default:null"`
	Only_sms string 		`gorm:"type:char(1);default:null"`
	Phn string 				`gorm:"size:15;not null"`
	Profile string 			`gorm:"size:50;default:null"`
	P_com string 			`gorm:"size:2;default:null"`
	P_invoice string 		`gorm:"size:100;default:null"`
	Reg_dt string 			`gorm:"size:20;not null"`
	Remark1 string 			`gorm:"size:50;default:null"`
	Remark2 string 			`gorm:"size:50;default:null"`
	Remark3 string 			`gorm:"size:50;default:null"`
	Remark4 string 			`gorm:"size:50;default:null"`
	Remark5 string 			`gorm:"size:50;default:null"`
	Reserve_dt string 		`gorm:"size:14;not null;index:idx_sendgroup_reservedt,priority:2"`
	Sms_kind string 		`gorm:"type:char(1);default:null"`
	Sms_lms_tit string 		`gorm:"size:120;default:null"`
	Sms_sender string 		`gorm:"size:15;default:null"`
	S_code string 			`gorm:"size:4;default:null"`
	Tmpl_id string 			`gorm:"size:30;default:null"`
	Wide string 			`gorm:"type:char(1);default:null"`
	Send_group string 		`gorm:"type:char(20);default:null;index:idx_sendgroup_reservedt,priority:1;index:idx_sendgroup"`
	Supplement string 		`gorm:"type:mediumtext;default:null"`
	Price uint 				`gorm:"type:int(11);;default:0"`
	Currency_type string 	`gorm:"type:char(3);default:null"`
	Title string 			`gorm:"size:50;default:null"`
	Header string 			`gorm:"size:100;default:null"`
	Carousel string 		`gorm:"type:text;default:null"`
	Attachments string 		`gorm:"type:text;default:null"`
	Link string 			`gorm:"type:text;default:null"`
	Mms_image_id string 	`gorm:"size:100;default:null"`
}

func (RequestAt) TableName() string {
	return "DHN_REQUEST_AT"
}

// TABLE
// NAME : DHN_REQUEST_RESEND
type RequestResend struct {
	Id uint 				`gorm:"type:bigint;primaryKey;size:20;autoIncrement"`
	Msgid string 			`gorm:"size:20;not null;index:idx_msgid"`
	Userid string 			`gorm:"size:20;not null"`
	Ad_flag string 			`gorm:"type:char(1);default:null"`
	Button1 string 			`gorm:"type:text;default:null"`
	Button2 string 			`gorm:"type:text;default:null"`
	Button3 string 			`gorm:"type:text;default:null"`
	Button4 string 			`gorm:"type:text;default:null"`
	Button5 string 			`gorm:"type:text;default:null"`
	Image_link string 		`gorm:"type:text;default:null"`
	Image_url string 		`gorm:"type:text;default:null"`
	Message_type string 	`gorm:"type:char(2);default:null"`
	Msg string 				`gorm:"type:text;not null"`
	Msg_sms string 			`gorm:"type:text;default:null"`
	Only_sms string 		`gorm:"type:char(1);default:null"`
	Phn string 				`gorm:"size:15;not null"`
	Profile string 			`gorm:"size:50;default:null"`
	P_com string 			`gorm:"size:2;default:null"`
	P_invoice string 		`gorm:"size:100;default:null"`
	Reg_dt string 			`gorm:"size:20;not null"`
	Remark1 string 			`gorm:"size:50;default:null"`
	Remark2 string 			`gorm:"size:50;default:null"`
	Remark3 string 			`gorm:"size:50;default:null"`
	Remark4 string 			`gorm:"size:50;default:null"`
	Remark5 string 			`gorm:"size:50;default:null"`
	Reserve_dt string 		`gorm:"size:14;not null;index:idx_sendgroup_reservedt,priority:2"`
	Sms_kind string 		`gorm:"type:char(1);default:null"`
	Sms_lms_tit string 		`gorm:"size:120;default:null"`
	Sms_sender string 		`gorm:"size:15;default:null"`
	S_code string 			`gorm:"size:4;default:null"`
	Tmpl_id string 			`gorm:"size:30;default:null"`
	Wide string 			`gorm:"type:char(1);default:null"`
	Send_group string 		`gorm:"type:char(20);default:null;index:idx_sendgroup_reservedt,priority:1;index:idx_sendgroup"`
	Supplement string 		`gorm:"type:mediumtext;default:null"`
	Price uint 				`gorm:"type:int(11);default:0"`
	Currency_type string 	`gorm:"type:char(3);default:null"`
	Title string 			`gorm:"size:50;default:null"`
	Header string 			`gorm:"size:100;default:null"`
	Carousel string 		`gorm:"type:text;default:null"`
	Attachments string 		`gorm:"type:text;default:null"`
	Att_items string 		`gorm:"type:text;default:null"`
	Att_coupon string 		`gorm:"type:text;default:null"`
	Mms_image_id string 	`gorm:"size:100;default:null"`
}

func (RequestResend) TableName() string {
	return "DHN_REQUEST_RESEND"
}

// TABLE
// NAME : DHN_REQUEST_AT_RESEND
type RequestAtResend struct {
	Id uint 				`gorm:"type:bigint;primaryKey;size:20;autoIncrement"`
	Msgid string 			`gorm:"size:20;not null;index:idx_msgid"`
	Userid string 			`gorm:"size:20;not null"`
	Ad_flag string 			`gorm:"type:char(1);default:null"`
	Button1 string 			`gorm:"type:text;default:null"`
	Button2 string 			`gorm:"type:text;default:null"`
	Button3 string 			`gorm:"type:text;default:null"`
	Button4 string 			`gorm:"type:text;default:null"`
	Button5 string 			`gorm:"type:text;default:null"`
	Image_link string 		`gorm:"type:text;default:null"`
	Image_url string 		`gorm:"type:text;default:null"`
	Message_type string 	`gorm:"type:char(2);default:null"`
	Msg string 				`gorm:"type:text;not null"`
	Msg_sms string 			`gorm:"type:char(2);default:null"`
	Only_sms string 		`gorm:"type:char(1);default:null"`
	Phn string 				`gorm:"size:15;not null"`
	Profile string 			`gorm:"size:50;default:null"`
	P_com string 			`gorm:"size:2;default:null"`
	P_invoice string 		`gorm:"size:100;default:null"`
	Reg_dt string 			`gorm:"size:20;not null"`
	Remark1 string 			`gorm:"size:50;default:null"`
	Remark2 string 			`gorm:"size:50;default:null"`
	Remark3 string 			`gorm:"size:50;default:null"`
	Remark4 string 			`gorm:"size:50;default:null"`
	Remark5 string 			`gorm:"size:50;default:null"`
	Reserve_dt string 		`gorm:"size:14;not null;index:idx_sendgroup_reservedt,priority:2"`
	Sms_kind string 		`gorm:"type:char(1);default:null"`
	Sms_lms_tit string 		`gorm:"size:120;default:null"`
	Sms_sender string 		`gorm:"size:15;default:null"`
	S_code string 			`gorm:"size:4;default:null"`
	Tmpl_id string 			`gorm:"size:30;default:null"`
	Wide string 			`gorm:"type:char(1);default:null"`
	Send_group string 		`gorm:"type:char(20);default:null;index:idx_sendgroup_reservedt,priority:1;index:idx_sendgroup"`
	Supplement string 		`gorm:"type:mediumtext;default:null"`
	Price uint 				`gorm:"type:int(11);;default:0"`
	Currency_type string 	`gorm:"type:char(3);default:null"`
	Title string 			`gorm:"size:50;default:null"`
	Header string 			`gorm:"size:100;default:null"`
	Carousel string 		`gorm:"type:text;default:null"`
	Attachments string 		`gorm:"type:text;default:null"`
	Link string 			`gorm:"type:text;default:null"`
	Mms_image_id string 	`gorm:"size:100;default:null"`
}

func (RequestAtResend) TableName() string {
	return "DHN_REQUEST_AT_RESEND"
}

// TABLE
// NAME : DHN_RESULT
type Result struct {
	Msgid string 			`gorm:"size:20;primaryKey;index:idx_userid_msgid,priority:2"`
	Userid string 			`gorm:"size:20;primaryKey;not null;index:idx_userid_msgid,priority:1;index:idx_userid_regdt,priority:1;index:idx_userid_sync_result,priority:1;index:idx_userid_result_sendgroup,priority:1"`
	Ad_flag string 			`gorm:"type:char(1);default:null"`
	Button1 string 			`gorm:"type:text;default:null"`
	Button2 string 			`gorm:"type:text;default:null"`
	Button3 string 			`gorm:"type:text;default:null"`
	Button4 string 			`gorm:"type:text;default:null"`
	Button5 string 			`gorm:"type:text;default:null"`
	Image_link string 		`gorm:"type:text;default:null"`
	Image_url string 		`gorm:"type:text;default:null"`
	Kind string 			`gorm:"type:char(1);default:null"`
	Message string 			`gorm:"type:text;default:null"`
	Message_type string 	`gorm:"type:char(2);default:null"`
	Msg string 				`gorm:"type:text;not null"`
	Msg_sms string 			`gorm:"type:text;default:null"`
	Only_sms string 		`gorm:"type:char(1);default:null"`
	Phn string 				`gorm:"size:15;not null"`
	Profile string 			`gorm:"size:50;default:null"`
	P_com string 			`gorm:"size:2;default:null"`
	P_invoice string 		`gorm:"size:100;default:null"`
	Reg_dt string 			`gorm:"size:20;not null;index:idx_userid_regdt,priority:2"`
	Res_dt string 			`gorm:"size:20;default:null"`
	Result string 			`gorm:"type:char(1);default:null;index:idx_userid_sync_result,priority:3;index:idx_userid_result_sendgroup,priority:2"`
	Code string 			`gorm:"size:4;default:null"`
	S_code string 			`gorm:"size:2;default:null"`
	Sync string 			`gorm:"type:char(1);not null;index:idx_userid_sync_result,priority:2"`
	Remark1 string 			`gorm:"size:50;default:null"`
	Remark2 string 			`gorm:"size:50;default:null"`
	Remark3 string 			`gorm:"size:50;default:null"`
	Remark4 string 			`gorm:"size:50;default:null"`
	Remark5 string 			`gorm:"size:50;default:null"`
	Reserve_dt string 		`gorm:"size:14;not null"`
	Sms_kind string 		`gorm:"type:char(1);default:null"`
	Sms_lms_tit string 		`gorm:"size:120;default:null"`
	Sms_sender string 		`gorm:"size:15;default:null"`
	Tmpl_id string 			`gorm:"size:30;default:null"`
	Wide string 			`gorm:"type:char(1);default:null"`
	Send_group string 		`gorm:"type:char(20);default:null;index:idx_userid_result_sendgroup,priority:3"`
	Supplement string 		`gorm:"type:mediumtext;default:null"`
	Price uint 				`gorm:"type:int(11);default:0"`
	Currency_type string 	`gorm:"type:char(3);default:null"`
	Title string 			`gorm:"size:50;default:null"`
	Header string 			`gorm:"size:100;default:null"`
	Carousel string 		`gorm:"type:text;default:null"`
	Attachments string 		`gorm:"type:text;default:null"`
	Att_items string 		`gorm:"type:text;default:null"`
	Att_coupon string 		`gorm:"type:text;default:null"`
	Link string 			`gorm:"type:text;default:null"`
	Mms_image_id string 	`gorm:"size:100;default:null"`
}

func (Result) TableName() string {
	return "DHN_RESULT"
}

// TABLE
// NAME : DHN_RESULT_TEMP
type ResultTemp struct {
	Msgid string 			`gorm:"size:20;primaryKey;index:idx_userid_msgid,priority:2"`
	Userid string 			`gorm:"size:20;primaryKey;not null;index:idx_userid_msgid,priority:1;index:idx_userid_regdt,priority:1;index:idx_userid_sync_result,priority:1;index:idx_userid_result_sendgroup,priority:1"`
	Ad_flag string 			`gorm:"type:char(1);default:null"`
	Button1 string 			`gorm:"type:text;default:null"`
	Button2 string 			`gorm:"type:text;default:null"`
	Button3 string 			`gorm:"type:text;default:null"`
	Button4 string 			`gorm:"type:text;default:null"`
	Button5 string 			`gorm:"type:text;default:null"`
	Image_link string 		`gorm:"type:text;default:null"`
	Image_url string 		`gorm:"type:text;default:null"`
	Kind string 			`gorm:"type:char(1);default:null"`
	Message string 			`gorm:"type:text;default:null"`
	Message_type string 	`gorm:"type:char(2);default:null"`
	Msg string 				`gorm:"type:text;not null"`
	Msg_sms string 			`gorm:"type:text;default:null"`
	Only_sms string 		`gorm:"type:char(1);default:null"`
	Phn string 				`gorm:"size:15;not null"`
	Profile string 			`gorm:"size:50;default:null"`
	P_com string 			`gorm:"size:2;default:null"`
	P_invoice string 		`gorm:"size:100;default:null"`
	Reg_dt string 			`gorm:"size:20;not null;index:idx_userid_regdt,priority:2"`
	Res_dt string 			`gorm:"size:20;default:null"`
	Result string 			`gorm:"type:char(1);default:null;index:idx_userid_sync_result,priority:3;index:idx_userid_result_sendgroup,priority:2"`
	Code string 			`gorm:"size:4;default:null"`
	S_code string 			`gorm:"size:2;default:null"`
	Sync string 			`gorm:"type:char(1);not null;index:idx_userid_sync_result,priority:2"`
	Remark1 string 			`gorm:"size:50;default:null"`
	Remark2 string 			`gorm:"size:50;default:null"`
	Remark3 string 			`gorm:"size:50;default:null"`
	Remark4 string 			`gorm:"size:50;default:null"`
	Remark5 string 			`gorm:"size:50;default:null"`
	Reserve_dt string 		`gorm:"size:14;not null"`
	Sms_kind string 		`gorm:"type:char(1);default:null"`
	Sms_lms_tit string 		`gorm:"size:120;default:null"`
	Sms_sender string 		`gorm:"size:15;default:null"`
	Tmpl_id string 			`gorm:"size:30;default:null"`
	Wide string 			`gorm:"type:char(1);default:null"`
	Send_group string 		`gorm:"type:char(20);default:null;index:idx_userid_result_sendgroup,priority:3"`
	Supplement string 		`gorm:"type:mediumtext;default:null"`
	Price uint 				`gorm:"type:int(11);default:0"`
	Currency_type string 	`gorm:"type:char(3);default:null"`
	Title string 			`gorm:"size:50;default:null"`
	Header string 			`gorm:"size:100;default:null"`
	Carousel string 		`gorm:"type:text;default:null"`
	Attachments string 		`gorm:"type:text;default:null"`
	Att_items string 		`gorm:"type:text;default:null"`
	Att_coupon string 		`gorm:"type:text;default:null"`
	Link string 			`gorm:"type:text;default:null"`
	Mms_image_id string 	`gorm:"size:100;default:null"`
}

func (ResultTemp) TableName() string {
	return "DHN_RESULT_TEMP"
}

// TABLE
// NAME : DHN_RESULT_BK_TEMP
type ResultBkTemp struct {
	Msgid string 			`gorm:"size:20;primaryKey;index:idx_userid_msgid,priority:2"`
	Userid string 			`gorm:"size:20;primaryKey;not null;index:idx_userid_msgid,priority:1;index:idx_userid_regdt,priority:1;index:idx_userid_sync_result,priority:1;index:idx_userid_result_sendgroup,priority:1"`
	Ad_flag string 			`gorm:"type:char(1);default:null"`
	Button1 string 			`gorm:"type:text;default:null"`
	Button2 string 			`gorm:"type:text;default:null"`
	Button3 string 			`gorm:"type:text;default:null"`
	Button4 string 			`gorm:"type:text;default:null"`
	Button5 string 			`gorm:"type:text;default:null"`
	Image_link string 		`gorm:"type:text;default:null"`
	Image_url string 		`gorm:"type:text;default:null"`
	Kind string 			`gorm:"type:char(1);default:null"`
	Message string 			`gorm:"type:text;default:null"`
	Message_type string 	`gorm:"type:char(2);default:null"`
	Msg string 				`gorm:"type:text;not null"`
	Msg_sms string 			`gorm:"type:text;default:null"`
	Only_sms string 		`gorm:"type:char(1);default:null"`
	Phn string 				`gorm:"size:15;not null"`
	Profile string 			`gorm:"size:50;default:null"`
	P_com string 			`gorm:"size:2;default:null"`
	P_invoice string 		`gorm:"size:100;default:null"`
	Reg_dt string 			`gorm:"size:20;not null;index:idx_userid_regdt,priority:2"`
	Res_dt string 			`gorm:"size:20;default:null"`
	Result string 			`gorm:"type:char(1);default:null;index:idx_userid_sync_result,priority:3;index:idx_userid_result_sendgroup,priority:2"`
	Code string 			`gorm:"size:4;default:null"`
	S_code string 			`gorm:"size:2;default:null"`
	Sync string 			`gorm:"type:char(1);not null;index:idx_userid_sync_result,priority:2"`
	Remark1 string 			`gorm:"size:50;default:null"`
	Remark2 string 			`gorm:"size:50;default:null"`
	Remark3 string 			`gorm:"size:50;default:null"`
	Remark4 string 			`gorm:"size:50;default:null"`
	Remark5 string 			`gorm:"size:50;default:null"`
	Reserve_dt string 		`gorm:"size:14;not null"`
	Sms_kind string 		`gorm:"type:char(1);default:null"`
	Sms_lms_tit string 		`gorm:"size:120;default:null"`
	Sms_sender string 		`gorm:"size:15;default:null"`
	Tmpl_id string 			`gorm:"size:30;default:null"`
	Wide string 			`gorm:"type:char(1);default:null"`
	Send_group string 		`gorm:"type:char(20);default:null;index:idx_userid_result_sendgroup,priority:3"`
	Supplement string 		`gorm:"type:mediumtext;default:null"`
	Price uint 				`gorm:"type:int(11);default:0"`
	Currency_type string 	`gorm:"type:char(3);default:null"`
	Title string 			`gorm:"size:50;default:null"`
	Header string 			`gorm:"size:100;default:null"`
	Carousel string 		`gorm:"type:text;default:null"`
	Attachments string 		`gorm:"type:text;default:null"`
	Att_items string 		`gorm:"type:text;default:null"`
	Att_coupon string 		`gorm:"type:text;default:null"`
	Link string 			`gorm:"type:text;default:null"`
	Mms_image_id string 	`gorm:"size:100;default:null"`
}

func (ResultBkTemp) TableName() string {
	return "DHN_RESULT_BK_TEMP"
}

// TABLE
// NAME : DHN_RESULT_STA
type ResultSta struct {
	Userid string `gorm:"size:20;not null;index:idx_userid_depart_senddate,priority:1"`
	Depart string `gorm:"size:100;default:null;index:idx_userid_depart_senddate,priority:2"`
	Send_date string `gorm:"size:20;default:null;index:idx_userid_depart_senddate,priority:3"`
	Send_cnt string `gorm:"type:int(11);default:0"`
	Ats_cnt string `gorm:"type:int(11);default:0"`
	Ate_cnt string `gorm:"type:int(11);default:0"`
	Fts_cnt string `gorm:"type:int(11);default:0"`
	Fte_cnt string `gorm:"type:int(11);default:0"`
	Ftis_cnt string `gorm:"type:int(11);default:0"`
	Ftie_cnt string `gorm:"type:int(11);default:0"`
	Ftws_cnt string `gorm:"type:int(11);default:0"`
	Ftwe_cnt string `gorm:"type:int(11);default:0"`
	Smss_cnt string `gorm:"type:int(11);default:0"`
	Smsd_cnt string `gorm:"type:int(11);default:0"`
	Smse_cnt string `gorm:"type:int(11);default:0"`
	Lmss_cnt string `gorm:"type:int(11);default:0"`
	Lmsd_cnt string `gorm:"type:int(11);default:0"`
	Lmse_cnt string `gorm:"type:int(11);default:0"`
	Mmss_cnt string `gorm:"type:int(11);default:0"`
	Mmsd_cnt string `gorm:"type:int(11);default:0"`
	Mmse_cnt string `gorm:"type:int(11);default:0"`
}

func (ResultSta) TableName() string {
	return "DHN_RESULT_STA"
}

// TABLE
// NAME : SPECIAL_CHARACTER
type SpecialCharacter struct {
	Origin_hex_code string `gorm:"size:20;uniqueIndex;default:null"`
	Dest_str string `gorm:"size:20;default:null"`
	Enabled string `gorm:"type:char(1);not null;default:'Y'"`
}

func (SpecialCharacter) TableName() string {
	return "SPECIAL_CHARACTER"
}

// TABLE
// NAME : API_MMS_IMAGES
type ApiMmsImages struct {
	Userid string `gorm:"size:20;default:null"`
	Mms_image_id string `gorm:"size:30;default:null"`
	Origin1_path string `gorm:"type:text;default:null"`
	Origin2_path string `gorm:"type:text;default:null"`
	Origin3_path string `gorm:"type:text;default:null"`
	File1_path string `gorm:"type:text;default:null"`
	File2_path string `gorm:"type:text;default:null"`
	File3_path string `gorm:"type:text;default:null"`
}

func (ApiMmsImages) TableName() string {
	return "API_MMS_IMAGES"
}