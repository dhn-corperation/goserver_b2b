package kakaoatjson

type Alimtalk struct {
	Message_type    string        `json:"message_type"`
	Serial_number   string        `json:"serial_number"`
	Sender_key      string        `json:"sender_key"`
	Phone_number    string        `json:"phone_number,omitempty"`
	App_user_id     string        `json:"app_user_id,omitempty"`
	Template_code   string        `json:"template_code"`
	Message         string        `json:"message"`
	Title           string        `json:"title,omitempty"`
	Header          string        `json:"header,omitempty"`
	Link            *ATLink       `json:"link,omitempty"`
	Response_method string        `json:"response_method"`
	Timeout         int           `json:"timeout,omitempty"`
	Attachment      *ATAttachment `json:"attachment,omitempty"`
	Supplement      *ATSupplement `json:"supplement,omitempty"`
	Channel_key     string        `json:"channel_key,omitempty"`
	Price           int64         `json:"price,omitempty"`
	Currency_type   string        `json:"currency_type,omitempty"`
}

type ATLink struct {
	URLMobile     string `json:"url_mobile,omitempty"`
	URLPc         string `json:"url_pc,omitempty"`
	SchemeIos     string `json:"scheme_ios,omitempty"`
	SchemeAndroid string `json:"scheme_android,omitempty"`
}

type ATAButton struct {
	Type          string `json:"type,omitempty"`
	Name          string `json:"name,omitempty"`
	SchemeAndroid string `json:"scheme_android,omitempty"`
	SchemeIos     string `json:"scheme_ios,omitempty"`
	ChatExtra     string `json:"chat_extra,omitempty"`
	ChatEvent     string `json:"chat_event,omitempty"`
	PluginID      string `json:"plugin_id,omitempty"`
	RelayID       string `json:"relay_id,omitempty"`
	OneclickID    string `json:"oneclick_id,omitempty"`
	ProductID     string `json:"product_id,omitempty"`
	URLMobile     string `json:"url_mobile,omitempty"`
	URLPc         string `json:"url_pc,omitempty"`
}

type ATAItemHighlight struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type ATAItem struct {
	List    []*ATAIList  `json:"list,omitempty"`
	Summary *ATAISummary `json:"summary,omitempty"`
}

type ATAIList struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type ATAISummary struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type ATAttachment struct {
	Button        []*ATAButton      `json:"button,omitempty"`
	ItemHighlight *ATAItemHighlight `json:"item_highlight,omitempty"`
	Item          *ATAItem          `json:"item,omitempty"`
}

type ATSQuickReply struct {
	Name          string `json:"name,omitempty"`
	Type          string `json:"type,omitempty"`
	URLMobile     string `json:"url_mobile,omitempty"`
	URLPc         string `json:"url_pc,omitempty"`
	SchemeIos     string `json:"scheme_ios,omitempty"`
	SchemeAndroid string `json:"scheme_android,omitempty"`
	ChatExtra     string `json:"chat_extra,omitempty"`
	ChatEvent     string `json:"chat_event,omitempty"`
}

type ATSupplement struct {
	QuickReply *[]ATSQuickReply `json:"quick_reply,omitempty"`
}

type KakaoResponse struct {
	Code        string
	Received_at string
	Message     string
}
