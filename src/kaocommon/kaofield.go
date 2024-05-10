package kaocommon

import (
	"context"
)

type AtReqColumn struct {
	Msgid			interface{}	`db:"msgid"`
	Userid			interface{}	`db:"userid"`
	Ad_flag			interface{}	`db:"ad_flag"`
	Button1			interface{}	`db:"button1"`
	Button2			interface{}	`db:"button2"`
	Button3			interface{}	`db:"button3"`
	Button4			interface{}	`db:"button4"`
	Button5			interface{}	`db:"button5"`
	Image_link		interface{}	`db:"image_link"`
	Image_url		interface{}	`db:"image_url"`
	Message_type    interface{}	`db:"message_type"`
	Msg    			interface{}	`db:"msg"`
	Msg_sms    		interface{}	`db:"msg_sms"`
	Only_sms    	interface{}	`db:"only_sms"`
	Phn    			interface{}	`db:"phn"`
	Profile    		interface{}	`db:"profile"`
	P_com    		interface{}	`db:"p_com"`
	P_invoice    	interface{}	`db:"p_invoice"`
	Reg_dt   		interface{}	`db:"reg_dt"`
	Remark1   		interface{}	`db:"remark1"`
	Remark2    		interface{}	`db:"remark2"`
	Remark3   		interface{}	`db:"remark3"`
	Remark4    		interface{}	`db:"remark4"`
	Remark5    		interface{}	`db:"remark5"`
	Reserve_dt    	interface{}	`db:"reserve_dt"`
	Sms_kind    	interface{}	`db:"sms_kind"`
	Sms_lms_tit     interface{}	`db:"sms_lms_tit"`
	Sms_sender      interface{}	`db:"sms_sender"`
	S_code    		interface{}	`db:"s_code"`
	Tmpl_id    		interface{}	`db:"tmpl_id"`
	Wide    		interface{}	`db:"wide"`
	Send_group    	interface{}	`db:"send_group"`
	Supplement    	interface{}	`db:"supplement"`
	Price           interface{}	`db:"price"`
	Currency_type   interface{}	`db:"currency_type"`
	Title           interface{}	`db:"title"`
}

type FtReqColumn struct {
	Msgid			interface{}	`db:"msgid"`
	Userid			interface{}	`db:"userid"`
	Ad_flag			interface{}	`db:"ad_flag"`
	Button1			interface{}	`db:"button1"`
	Button2			interface{}	`db:"button2"`
	Button3			interface{}	`db:"button3"`
	Button4			interface{}	`db:"button4"`
	Button5			interface{}	`db:"button5"`
	Image_link		interface{}	`db:"image_link"`
	Image_url		interface{}	`db:"image_url"`
	Message_type    interface{}	`db:"message_type"`
	Msg    			interface{}	`db:"msg"`
	Msg_sms    		interface{}	`db:"msg_sms"`
	Only_sms    	interface{}	`db:"only_sms"`
	Phn    			interface{}	`db:"phn"`
	Profile    		interface{}	`db:"profile"`
	P_com    		interface{}	`db:"p_com"`
	P_invoice    	interface{}	`db:"p_invoice"`
	Reg_dt   		interface{}	`db:"reg_dt"`
	Remark1   		interface{}	`db:"remark1"`
	Remark2    		interface{}	`db:"remark2"`
	Remark3   		interface{}	`db:"remark3"`
	Remark4    		interface{}	`db:"remark4"`
	Remark5    		interface{}	`db:"remark5"`
	Reserve_dt    	interface{}	`db:"reserve_dt"`
	Sms_kind    	interface{}	`db:"sms_kind"`
	Sms_lms_tit     interface{}	`db:"sms_lms_tit"`
	Sms_sender      interface{}	`db:"sms_sender"`
	S_code    		interface{}	`db:"s_code"`
	Tmpl_id    		interface{}	`db:"tmpl_id"`
	Wide    		interface{}	`db:"wide"`
	Send_group    	interface{}	`db:"send_group"`
	Supplement    	interface{}	`db:"supplement"`
	Price           interface{}	`db:"price"`
	Currency_type   interface{}	`db:"currency_type"`
	Title           interface{}  `db:"title"`
	Header    		interface{}	`db:"header"`
	Carousel    	interface{}	`db:"carousel"`
	Att_items       interface{}  `db:"att_items"`
	Att_coupon      interface{}  `db:"att_coupon"`
	Attachments     interface{}  `db:"attachments"`
}

type MsgReqColumn struct {
	Msgid			interface{}	`db:"msgid"`
	Userid			interface{}	`db:"userid"`
	Ad_flag			interface{}	`db:"ad_flag"`
	Button1			interface{}	`db:"button1"`
	Button2			interface{}	`db:"button2"`
	Button3			interface{}	`db:"button3"`
	Button4			interface{}	`db:"button4"`
	Button5			interface{}	`db:"button5"`
	Code     		interface{}	`db:"code"`
	Image_link		interface{}	`db:"image_link"`
	Image_url		interface{}	`db:"image_url"`
	Kind            interface{}	`db:"kind"`
	Message         interface{}	`db:"message"`
	Message_type    interface{}	`db:"message_type"`
	Msg    			interface{}	`db:"msg"`
	Msg_sms    		interface{}	`db:"msg_sms"`
	Only_sms    	interface{}	`db:"only_sms"`
	Phn    			interface{}	`db:"phn"`
	Profile    		interface{}	`db:"profile"`
	P_com    		interface{}	`db:"p_com"`
	P_invoice    	interface{}	`db:"p_invoice"`
	Reg_dt   		interface{}	`db:"reg_dt"`
	Remark1   		interface{}	`db:"remark1"`
	Remark2    		interface{}	`db:"remark2"`
	Remark3   		interface{}	`db:"remark3"`
	Remark4    		interface{}	`db:"remark4"`
	Remark5    		interface{}	`db:"remark5"`
	Res_dt    		interface{}	`db:"res_dt"`
	Reserve_dt    	interface{}	`db:"reserve_dt"`
	Result    		interface{}	`db:"result"`
	S_code    		interface{}	`db:"s_code"`
	Sms_kind    	interface{}	`db:"sms_kind"`
	Sms_lms_tit     interface{}	`db:"sms_lms_tit"`
	Sms_sender      interface{}	`db:"sms_sender"`
	Sync    		interface{}	`db:"sync"`
	Tmpl_id    		interface{}	`db:"tmpl_id"`
	Wide    		interface{}	`db:"wide"`
	Send_group    	interface{}	`db:"send_group"`
	Supplement    	interface{}	`db:"supplement"`
	Price           interface{}	`db:"price"`
	Currency_type   interface{}	`db:"currency_type"`
	Header    		interface{}	`db:"header"`
	Carousel    	interface{}	`db:"carousel"`
}

type AtResColumn struct {
	Msgid			interface{}	`db:"msgid"`
	Userid			interface{}	`db:"userid"`
	Ad_flag			interface{}	`db:"ad_flag"`
	Button1			interface{}	`db:"button1"`
	Button2			interface{}	`db:"button2"`
	Button3			interface{}	`db:"button3"`
	Button4			interface{}	`db:"button4"`
	Button5			interface{}	`db:"button5"`
	Code			interface{}	`db:"code"`
	Image_link		interface{}	`db:"image_link"`
	Image_url		interface{}	`db:"image_url"`
	Kind            interface{}	`db:"kind"`
	Message         interface{}	`db:"message"`
	Message_type    interface{}	`db:"message_type"`
	Msg    			interface{}	`db:"msg"`
	Msg_sms    		interface{}	`db:"msg_sms"`
	Only_sms    	interface{}	`db:"only_sms"`
	P_com    		interface{}	`db:"p_com"`
	P_invoice    	interface{}	`db:"p_invoice"`
	Phn    			interface{}	`db:"phn"`
	Profile    		interface{}	`db:"profile"`
	Reg_dt   		interface{}	`db:"reg_dt"`
	Remark1   		interface{}	`db:"remark1"`
	Remark2    		interface{}	`db:"remark2"`
	Remark3   		interface{}	`db:"remark3"`
	Remark4    		interface{}	`db:"remark4"`
	Remark5    		interface{}	`db:"remark5"`
	Res_dt    		interface{}	`db:"res_dt"`
	Reserve_dt    	interface{}	`db:"reserve_dt"`
	Result    		interface{}	`db:"result"`
	S_code    		interface{}	`db:"s_code"`
	Sms_kind    	interface{}	`db:"sms_kind"`
	Sms_lms_tit     interface{}	`db:"sms_lms_tit"`
	Sms_sender      interface{}	`db:"sms_sender"`
	Sync    		interface{}	`db:"sync"`
	Tmpl_id    		interface{}	`db:"tmpl_id"`
	Wide    		interface{}	`db:"wide"`
	Send_group    	interface{}	`db:"send_group"`
	Supplement    	interface{}	`db:"supplement"`
	Price           interface{}	`db:"price"`
	Currency_type   interface{}	`db:"currency_type"`
	Title           interface{}	`db:"title"`
}

type FtResColumn struct {
	Msgid			interface{}	`db:"msgid"`
	Userid			interface{}	`db:"userid"`
	Ad_flag			interface{}	`db:"ad_flag"`
	Button1			interface{}	`db:"button1"`
	Button2			interface{}	`db:"button2"`
	Button3			interface{}	`db:"button3"`
	Button4			interface{}	`db:"button4"`
	Button5			interface{}	`db:"button5"`
	Code			interface{}	`db:"code"`
	Image_link		interface{}	`db:"image_link"`
	Image_url		interface{}	`db:"image_url"`
	Kind            interface{}	`db:"kind"`
	Message         interface{}	`db:"message"`
	Message_type    interface{}	`db:"message_type"`
	Msg    			interface{}	`db:"msg"`
	Msg_sms    		interface{}	`db:"msg_sms"`
	Only_sms    	interface{}	`db:"only_sms"`
	P_com    		interface{}	`db:"p_com"`
	P_invoice    	interface{}	`db:"p_invoice"`
	Phn    			interface{}	`db:"phn"`
	Profile    		interface{}	`db:"profile"`
	Reg_dt   		interface{}	`db:"reg_dt"`
	Remark1   		interface{}	`db:"remark1"`
	Remark2    		interface{}	`db:"remark2"`
	Remark3   		interface{}	`db:"remark3"`
	Remark4    		interface{}	`db:"remark4"`
	Remark5    		interface{}	`db:"remark5"`
	Res_dt    		interface{}	`db:"res_dt"`
	Reserve_dt    	interface{}	`db:"reserve_dt"`
	Result    		interface{}	`db:"result"`
	S_code    		interface{}	`db:"s_code"`
	Sms_kind    	interface{}	`db:"sms_kind"`
	Sms_lms_tit     interface{}	`db:"sms_lms_tit"`
	Sms_sender      interface{}	`db:"sms_sender"`
	Sync    		interface{}	`db:"sync"`
	Tmpl_id    		interface{}	`db:"tmpl_id"`
	Wide    		interface{}	`db:"wide"`
	Send_group    	interface{}	`db:"send_group"`
	Supplement    	interface{}	`db:"supplement"`
	Price           interface{}	`db:"price"`
	Currency_type   interface{}	`db:"currency_type"`
	Header          interface{}	`db:"header"`
	Carousel        interface{}	`db:"carousel"`
}

type AtPollingResColumn struct {
	Msgid			interface{}	`db:"msgid"`
	Type			interface{}	`db:"type"`
}


type CheckUserReturnField struct {
	Validation	bool
	Userid 	   	string
	Userip		string
	Ctx 		context.Context
	SendLimit   string
}

var ResultTempMigrationColumn = []string{
	"msgid",
    "userid",
    "ad_flag",
    "button1",
    "button2",
    "button3",
    "button4",
    "button5",
    "code",
   	"image_link",
    "image_url",
    "kind",
    "message",
    "message_type",
    "msg",
    "msg_sms",
    "only_sms",
    "p_com",
    "p_invoice",
    "phn",
    "profile",
    "reg_dt",
    "remark1",
    "remark2",
    "remark3",
    "remark4",
    "remark5",
    "res_dt",
    "reserve_dt",
    "result",
    "s_code",
    "sms_kind",
    "sms_lms_tit",
    "sms_sender",
    "sync",
    "tmpl_id",
    "wide",
    "null",
    "supplement",
    "price",
    "currency_type",
    "title",
    "header",
    "carousel",
    "attachments",
    "user_key",
    "response_method",
    "timeout",
}

type OshotReqColumn struct {
	MsgGroupID string
	Sender string
	Receiver string
	Subject string
	Msg string
	Url string
	FilePath1 string
	FilePath2 string
	FilePath3 string
	CbMsgId string
	UserId string
}

type NanoReqColumn struct {
	CALLBACK string
	PHONE string
	SUBJECT string
	MSG string
	REQDATE string
	TR_SENDDATE string
	TR_SENDSTAT string
	TR_MSGTYPE string
	STATUS string
	FILE_CNT string
	FILE_PATH1 string
	FILE_PATH2 string
	FILE_PATH3 string
	ETC9 string
	ETC10 string
	IDENTIFICATION_CODE string
	ETC8 string
}

type NanoSMSReqColumn struct {
	TR_CALLBACK string
	TR_PHONE string
	TR_MSG string
	TR_SENDDATE string
	TR_SENDSTAT string
	TR_MSGTYPE string
	TR_ETC9 string
	TR_ETC10 string
	TR_IDENTIFICATION_CODE string
	TR_ETC8 string
}