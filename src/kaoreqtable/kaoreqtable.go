package kaoreqtable

//"database/sql"

type Reqtable struct {
	Msgid        string `json:"msgid"`
	Adflag       string `json:"adflag"`
	Button1      string `json:"button1"`
	Button2      string `json:"button2"`
	Button3      string `json:"button3"`
	Button4      string `json:"button4"`
	Button5      string `json:"button5"`
	Imagelink    string `json:"imagelink"`
	Imageurl     string `json:"imageurl"`
	Messagetype  string `json:"messagetype"`
	Msg          string `json:"msg"`
	Msgsms       string `json:"msgsms"`
	Onlysms      string `json:"onlysms"`
	Pcom         string `json:"pcom"`
	Pinvoice     string `json:"pinvoice"`
	Phn          string `json:"phn"`
	Profile      string `json:"profile"`
	Regdt        string `json:"regdt"`
	Remark1      string `json:"remark1"`
	Remark2      string `json:"remark2"`
	Remark3      string `json:"remark3"`
	Remark4      string `json:"remark4"`
	Remark5      string `json:"remark5"`
	Reservedt    string `json:"reservedt"`
	Scode        string `json:"scode"`
	Smskind      string `json:"smskind"`
	Smslmstit    string `json:"smslmstit"`
	Smssender    string `json:"smssender"`
	Tmplid       string `json:"tmplid"`
	Wide         string `json:"wide"`
	Supplement   string `json:"supplement"`
	Price        string `json:"price"`
	Currencytype string `json:"currencytype"`
	Title        string `json:"title"`
	Header       string `json:"header"`
	Carousel     string `json:"carousel"`
	Att_items    string `json:"att_items"`
	Att_coupon   string `json:"att_coupon"`
	Crypto       string `json:"crypto"`
}

type NanoReqTable struct {
	Msgid          string         `json:"msgid,omitempty"`
	Adflag         string         `json:"adflag,omitempty"`
	Button1        string         `json:"button1,omitempty"`
	Button2        string         `json:"button2,omitempty"`
	Button3        string         `json:"button3,omitempty"`
	Button4        string         `json:"button4,omitempty"`
	Button5        string         `json:"button5,omitempty"`
	Imagelink      string         `json:"imagelink,omitempty"`
	Imageurl       string         `json:"imageurl,omitempty"`
	Messagetype    string         `json:"messagetype,omitempty"`
	Msg            string         `json:"msg,omitempty"`
	Msgsms         string         `json:"msgsms,omitempty"`
	Onlysms        string         `json:"onlysms,omitempty"`
	Pcom           string         `json:"pcom,omitempty"`
	Pinvoice       string         `json:"pinvoice,omitempty"`
	Phn            string         `json:"phn,omitempty"`
	Profile        string         `json:"profile,omitempty"`
	Regdt          string         `json:"regdt,omitempty"`
	Remark1        string         `json:"remark1,omitempty"`
	Remark2        string         `json:"remark2,omitempty"`
	Remark3        string         `json:"remark3,omitempty"`
	Remark4        string         `json:"remark4,omitempty"`
	Remark5        string         `json:"remark5,omitempty"`
	Reservedt      string         `json:"reservedt,omitempty"`
	Scode          string         `json:"scode,omitempty"`
	Smskind        string         `json:"smskind,omitempty"`
	Smslmstit      string         `json:"smslmstit,omitempty"`
	Smssender      string         `json:"smssender,omitempty"`
	Tmplid         string         `json:"tmplid,omitempty"`
	Wide           string         `json:"wide,omitempty"`
	Supplement     string         `json:"supplement,omitempty"`
	Price          string         `json:"price,omitempty"`
	Currencytype   string         `json:"currencytype,omitempty"`
	Title          string         `json:"title,omitempty"`
	Header         string         `json:"header,omitempty"`
	Carousel       string         `json:"carousel,omitempty"`
	Att_items      string         `json:"att_items,omitempty"`
	Att_coupon     string         `json:"att_coupon,omitempty"`
	AttFilePath    string         `json:"att_file_path,omitempty"`
	Crypto         string         `json:"crypto,omitempty"`
	Userkey        string         `json:"user_key,omitempty"`
	ResponseMethod string         `json:"response_method,omitempty"`
	Timeout        string         `json:"timeout,omitempty"`
	Attachments    *NRAttachments `json:"attachments,omitempty"`
}

type NRAButton struct {
	Name          string `json:"name,omitempty"`
	Type          string `json:"type,omitempty"`
	URLMobile     string `json:"url_mobile,omitempty"`
	URLPc         string `json:"url_pc,omitempty"`
	SchemeIos     string `json:"scheme_ios,omitempty"`
	SchemeAndroid string `json:"scheme_android,omitempty"`
	ChatExtra     string `json:"chat_extra,omitempty"`
	ChatEvent     string `json:"chat_event,omitempty"`
	PluginID      string `json:"plugin_id,omitempty"`
	RelayID       string `json:"relay_id,omitempty"`
	OneclickID    string `json:"oneclick_id,omitempty"`
	ProductID     string `json:"product_id,omitempty"`
}

type NRAImage struct {
	ImgURL  string `json:"img_url,omitempty"`
	ImgLink string `json:"img_link,omitempty"`
}

type NRAItemHighlight struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type NRAItem struct {
	List    *[]NRAIList  `json:"list,omitempty"`
	Summary *NRAISummary `json:"summary,omitempty"`
}

type NRAIList struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type NRAISummary struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type NRAExtra struct {
	MsgType      string          `json:"msg_type,omitempty"`
	Title        string          `json:"title,omitempty"`
	Supplement   *NRAESupplement `json:"supplement,omitempty"`
	Price        int             `json:"price,omitempty"`
	CurrencyType string          `json:"currency_type,omitempty"`
	Header       string          `json:"header,omitempty"`
	Wide         string          `json:"wide,omitempty"`
	Link         *NRAELink       `json:"link,omitempty"`
}

type NRAESupplement struct {
	QuickReply *[]NRAESQuickReply `json:"quick_reply,omitempty"`
}

type NRAELink struct {
	URLMobile     string `json:"url_mobile,omitempty"`
	URLPc         string `json:"url_pc,omitempty"`
	SchemeIos     string `json:"scheme_ios,omitempty"`
	SchemeAndroid string `json:"scheme_android,omitempty"`
}

type NRAESQuickReply struct {
	Name          string `json:"name,omitempty"`
	Type          string `json:"type,omitempty"`
	SchemeAndroid string `json:"scheme_android,omitempty"`
	SchemeIos     string `json:"scheme_ios,omitempty"`
	URLMobile     string `json:"url_mobile,omitempty"`
	URLPc         string `json:"url_pc,omitempty"`
	ChatExtra     string `json:"chat_extra,omitempty"`
	ChatEvent     string `json:"chat_event,omitempty"`
}

type NRAttachments struct {
	Button        *[]NRAButton      `json:"button,omitempty"`
	Image         *NRAImage         `json:"image,omitempty"`
	ItemHighlight *NRAItemHighlight `json:"item_highlight,omitempty"`
	Item          *NRAItem          `json:"item,omitempty"`
	Extra         *NRAExtra         `json:"extra,omitempty"`
}
