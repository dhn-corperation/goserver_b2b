package structs

type Alimtalk struct {
	Message_type    string      `json:"message_type"`
	Serial_number   string      `json:"serial_number"`
	Sender_key      string      `json:"sender_key"`
	Phone_number    string      `json:"phone_number,omitempty"`
	App_user_id     string      `json:"app_user_id,omitempty"`
	Template_code   string      `json:"template_code"`
	Message         string      `json:"message"`
	Title           string      `json:"title,omitempty"`
	Header          string      `json:"header,omitempty"`
	Response_method string      `json:"response_method"`
	Timeout         int         `json:"timeout,omitempty"`
	Attachment      AttachmentB `json:"attachment,omitempty"`
	Supplement      Supplement  `json:"supplement,omitempty"`
	Link            *Link       `json:"link,omitempty"`
	Channel_key     string      `json:"channel_key,omitempty"`
	Price           int64       `json:"price,omitempty"`
	Currency_type   string      `json:"currency_type,omitempty"`
}

type Friendtalk struct {
	Message_type  string     `json:"message_type"`
	Serial_number string     `json:"serial_number"`
	Sender_key    string     `json:"sender_key"`
	Phone_number  string     `json:"phone_number,omitempty"`
	App_user_id   string     `json:"app_user_id,omitempty"`
	User_key      string     `json:"user_key,omitempty"`
	Message       string     `json:"message"`
	Ad_flag       string     `json:"ad_flag,omitempty"`
	Attachment    Attachment `json:"attachment,omitempty"`
	Header        string     `json:"header,omitempty"`
	Carousel      *FCarousel  `json:"carousel,omitempty"`
}

type Button struct {
	Name           string `json:"name"`
	Type           string `json:"type"`
	Scheme_android string `json:"scheme_android,omitempty"`
	Scheme_ios     string `json:"scheme_ios,omitempty"`
	Url_mobile     string `json:"url_mobile,omitempty"`
	Url_pc         string `json:"url_pc,omitempty"`
	Chat_extra     string `json:"chat_extra,omitempty"`
	Chat_event     string `json:"chat_event,omitempty"`
	Plugin_id      string `json:"plugin_id,omitempty"`
	Relay_id       string `json:"relay_id,omitempty"`
	Oneclick_id    string `json:"oneclick_id,omitempty"`
	Product_id     string `json:"product_id,omitempty"`
}

type CButton struct {
	Name         string `json:"name"`
	LinkType     string `json:"linktype"`
	LinkTypeName string `json:"linktypeName"`
	Ordering     string `json:"ordering"`
	LinkMo       string `json:"linkMo"`
	LinkPc       string `json:"linkPc"`
	LinkAnd      string `json:"linkAnd"`
	inkIos       string `json:"linkIos"`
	Pluginid     string `json:"pluginid"`
}

type Quickreply struct {
	Name           string `json:"name"`
	Type           string `json:"type"`
	Scheme_android string `json:"scheme_android,omitempty"`
	Scheme_ios     string `json:"scheme_ios,omitempty"`
	Url_mobile     string `json:"url_mobile,omitempty"`
	Url_pc         string `json:"url_pc,omitempty"`
	Chat_extra     string `json:"chat_extra,omitempty"`
	Chat_event     string `json:"chat_event,omitempty"`
}

type CQuickreply struct {
	Name     string `json:"name"`
	LinkType string `json:"linkType"`
	LinkMo   string `json:"linkMo"`
	LinkPc   string `json:"linkPc"`
	LinkAnd  string `json:"linkAnd"`
	inkIos   string `json:"linkIos"`
}

type Attachment struct {
	Buttons []Button     `json:"button,omitempty"`
	Ftimage *Image        `json:"image,omitempty"`
	Item    *AttItem      `json:"item,omitempty"`
	Coupon  *AttCoupon    `json:"coupon,omitempty"`
}

type AttachmentB struct {
	Buttons []Button `json:"button,omitempty"`
	Item_highlights *Item_highlight `json:"item_highlight,omitempty"`
	Items *Item `json:"item,omitempty"`
}

type AttachmentC struct {
	Item_highlights *Item_highlight `json:"item_highlight,omitempty"`
	Items *Item `json:"item,omitempty"`
}

type Item_highlight struct {
	Title string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type Item struct {
	Lists *[]AtItemList `json:"list,omitempty"`
	Summary *Summary `json:"summary,omitempty"`
}

type AtItemList struct {
	Title string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type Summary struct {
	Title string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type Supplement struct {
	Quick_reply []Quickreply `json:"quick_reply,omitempty"`
}

type Link struct {
	Url_mobile *string `json:"url_mobile,omitempty"`
	Url_pc *string `json:"url_pc,omitempty"`
	Scheme_android *string `json:"scheme_android,omitempty"`
	Scheme_ios *string `json:"scheme_ios,omitempty"`
}

type Image struct {
	Img_url  string `json:"img_url"`
	Img_link string `json:"img_link,omitempty"`
}

type AttCoupon struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	Url_pc          string `json:"url_pc,omitempty"`
	Url_mobile      string `json:"url_mobile,omitempty"`
	Scheme_android  string `json:"scheme_android,omitempty"`
	Scheme_ios      string `json:"scheme_ios,omitempty"`
}

type KakaoResponse struct {
	Code        string
	Received_at string
	Message     string
}

type KakaoResponse2 struct {
	Code        string
	Message     string
}

type PollingResponse struct {
	Code         string
	Response_id  int
	Response     PResponse
	Responsed_at string
	Message      string
}

type PResponse struct {
	Success []PResult
	Fail    []PResult
}

type PResult struct {
	Serial_number string
	Status        string
	Received_at   string
}

type FCarousel struct {
	List     []CarouselList `json:"list,omitempty"`
	Tail     CarouselTail `json:"tail,omitempty"`
}

type CarouselList struct {
	Header        string              `json:"header"`
	Message       string              `json:"message"`
	Attachment    CarouselAttachment  `json:"attachment"`
}

type TCarousel struct {
	List     []TCarouselList `json:"list"`
	Tail     CarouselTail `json:"tail,omitempty"`
}

type TCarouselList struct {
	Header        string              `json:"header"`
	Message       string              `json:"message"`
	Attachment    string              `json:"attachment,omitempty"`
}

type CarouselAttachment struct {
    Button  []CarouselButton `json:"button"` 
	Image     CarouselImage  `json:"image,omitempty"`
}

type CarouselTail struct {
	Url_pc          string              `json:"url_pc,omitempty"`
	Url_mobile      string              `json:"url_mobile,omitempty"`
	Scheme_ios      string              `json:"scheme_ios,omitempty"`
	Scheme_android  string              `json:"scheme_android,omitempty"`
}

type CarouselImage struct {
	Img_url  string `json:"img_url"`
	Img_link string `json:"img_link"`
}

type CarouselButton struct {
	Name           string `json:"name"`
	Type           string `json:"type"`
	Scheme_android string `json:"scheme_android,omitempty"`
	Scheme_ios     string `json:"scheme_ios,omitempty"`
	Url_mobile     string `json:"url_mobile,omitempty"`
	Url_pc         string `json:"url_pc,omitempty"`
}

type AttItem struct {
	AList   []ItemLists `json:"list,omitempty"`
}

type ItemLists struct {
	Title          string `json:"title"`
	Img_url        string `json:"img_url"`
	Scheme_android string `json:"scheme_android,omitempty"`
	Scheme_ios     string `json:"scheme_ios,omitempty"`
	Url_mobile     string `json:"url_mobile"`
	Url_pc         string `json:"url_pc,omitempty"`
}

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
	Attachments  string `json:"attachments"`
	Link  		 string `json:"link"`
	Carousel     string `json:"carousel"`
	Att_items    string `json:"att_items"`
	Att_coupon   string `json:"att_coupon"`
	MmsImageId   string `json:"mmsimageid"`
	Crypto       string `json:"crypto"`
}

type ReceiveRes struct {
	Code string `json:"code"`
	Message string `json:"message"`
	AtCnt *int16 `json:"atcnt,omitempty"`
	FtCnt *int16 `json:"ftcnt,omitempty"`
	MsgCnt *int16 `json:"msgcnt,omitempty"`
	DuplCnt *int16 `json:"duplcnt,omitempty"`
	DuplMsgId *[]string `json:"duplMsgId,omitempty"`
}















////////////////////////////////////////////////////NPS AREA////////////////////////////////////////////////////

type CtReq17 struct {
	SenderKey *string `json:"senderKey,omitempty"`
	SenderKeyType *string `json:"senderKeyType,omitempty"`
	TemplateCode *string `json:"templateCode,omitempty"`
	TemplateName *string `json:"templateName,omitempty"`
	TemplateContent *string `json:"templateContent,omitempty"`
	ButtonType *string `json:"buttonType,omitempty"`
	ButtonName *string `json:"buttonName,omitempty"`
	ButtonUrl *string `json:"buttonUrl,omitempty"`
	Buttons *string `json:"buttons,omitempty"`
}

type CtReqButton17 struct {
	Name           *string `json:"name"`
	Type           *string `json:"type"`
	Ordering       *string `json:"ordering"`
	SchemeAndroid  *string `json:"scheme_android,omitempty"`
	SchemeIos      *string `json:"scheme_ios,omitempty"`
	UrlMobile      *string `json:"url_mobile,omitempty"`
	UrlPc          *string `json:"url_pc,omitempty"`
	PluginId       *string `json:"plugin_id,omitempty"`
}

type CtRes17 struct {
	Code string `json:"code"`
	Data *StKakaoRes `json:"data,omitempty"` 
	Message *string `json:"message,omitempty"`
}

type CtKakaoReq struct {
	SenderKey *string `json:"senderKey,omitempty"`
	SenderKeyType *string `json:"senderKeyType,omitempty"`
	TemplateCode *interface{} `json:"templateCode,omitempty"`
	TemplateName *interface{} `json:"templateName,omitempty"`
	TemplateMessageType string `json:"templateMessageType,omitempty"`
	TemplateEmphasizeType string `json:"templateEmphasizeType,omitempty"`
	TemplateContent *interface{} `json:"templateContent,omitempty"`
	TemplatePreviewMessage *string `json:"templatePreviewMessage,omitempty"`
	TemplateExtra *string `json:"templateExtra,omitempty"`
	TemplateImageName *string `json:"templateImageName,omitempty"`
	TemplateImageUrl *string `json:"templateImageUrl,omitempty"`
	TemplateTitle *string `json:"templateTitle,omitempty"`
	TemplateSubtitle *string `json:"templateSubtitle,omitempty"`
	TemplateHeader *string `json:"templateHeader,omitempty"`
	TemplateItemHighlight *KakaoTemplateItemHighlightNps `json:"templateItemHighlight,omitempty"`
	TemplateRepresentLink *KakaoTemplateRepresentLinkNps `json:"templateRepresentLink,omitempty"`
	CategoryCode *string `json:"categoryCode,omitempty"`
	SecurityFlag *bool `json:"securityFlag,omitempty"`
	Buttons *[]KakaoButtonsNps `json:"buttons,omitempty"`
	QuickReplies *[]KakaoQuickRepliesNps `json:"quickReplies,omitempty"`
}

type KakaoTemplateItemHighlightNps struct {
	Title *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	ImageUrl *string `json:"imageUrl,omitempty"`
}

type KakaoTemplateRepresentLinkNps struct {
	LinkAnd *string `json:"linkAnd,omitempty"`
	LinkIos *string `json:"linkIos,omitempty"`
	LinkPc *string `json:"linkPc,omitempty"`
	LinkMo *string `json:"linkMo,omitempty"`
}

type KakaoButtonsNps struct {
	Name         *interface{} `json:"name,omitempty"`
	LinkType     *interface{} `json:"linkType,omitempty"`
	Ordering     *interface{}    `json:"ordering,omitempty"`
	LinkMo       *interface{} `json:"linkMo,omitempty"`
	LinkPc       *interface{} `json:"linkPc,omitempty"`
	LinkAnd      *interface{} `json:"linkAnd,omitempty"`
	LinkIos      *interface{} `json:"linkIos,omitempty"`
	PluginId     *interface{} `json:"pluginId,omitempty"`
	BizFormId    *interface{} `json:"bizFormId,omitempty"`
}

type KakaoQuickRepliesNps struct {
	Name         *string `json:"name,omitempty"`
	LinkType     *string `json:"linktype,omitempty"`
	LinkMo       *string `json:"linkMo,omitempty"`
	LinkPc       *string `json:"linkPc,omitempty"`
	LinkAnd      *string `json:"linkAnd,omitempty"`
	inkIos       *string `json:"linkIos,omitempty"`
	PluginId     *string `json:"pluginid,omitempty"`
	BizFormId    *string `json:"bizFormId,omitempty"`
}

type StKakaoRes struct {
	Code *string `json:"code,omitempty"`
	Data *StKakaoData `json:"data,omitempty"`
	Message *string `json:"message,omitempty"`
}

type StKakaoData struct {
	SenderKey *string `json:"senderKey,omitempty"`
	SenderKeyType *string `json:"senderKeyType,omitempty"`
	TemplateCode *string `json:"templateCode,omitempty"`
	TemplateName *string `json:"templateName,omitempty"`
	TemplateMessageType *string `json:"templateMessageType,omitempty"`
	TemplateEmphasizeType *string `json:"templateEmphasizeType,omitempty"`
	TemplateContent *string `json:"templateContent,omitempty"`
	TemplatePreviewMessage *string `json:"templatePreviewMessage,omitempty"`
	TemplateExtra *string `json:"templateExtra,omitempty"`
	TemplateImageName *string `json:"templateImageName,omitempty"`
	TemplateImageUrl *string `json:"templateImageUrl,omitempty"`
	TemplateTitle *string `json:"templateTitle,omitempty"`
	TemplateSubtitle *string `json:"templateSubtitle,omitempty"`
	TemplateHeader *string `json:"templateHeader,omitempty"`
	TemplateItemHighlight *KakaoTemplateItemHighlightNps `json:"templateItemHighlight,omitempty"`
	TemplateRepresentLink *KakaoTemplateRepresentLinkNps `json:"templateRepresentLink,omitempty"`
	InspectionStatus string `json:"inspectionStatus,omitempty"`
	CreatedAt *string `json:"createdAt,omitempty"`
	ModifiedAt *string `json:"modifiedAt,omitempty"`
	Status *string `json:"status,omitempty"`
	Block *bool `json:"block,omitempty"`
	Dormant *bool `json:"dormant,omitempty"`
	CategoryCode *string `json:"categoryCode,omitempty"`
	SecurityFlag *bool `json:"securityFlag,omitempty"`
	Comments *[]KakaoCommentsNps `json:"comments,omitempty"`
	Buttons *[]KakaoButtonsNps `json:"buttons,omitempty"`
	QuickReplies *[]KakaoQuickRepliesNps `json:"quickReplies,omitempty"`
}

type KakaoCommentsNps struct {
	Id *int `json:"id,omitempty"`
	Content *string `json:"content,omitempty"`
	UserName *string `json:"userName,omitempty"`
	CreatedAt *string `json:"createdAt,omitempty"`
	Status *string `json:"status,omitempty"`
	Attachment *[]AttachmentNps `json:"attachment,omitempty"`
}

type AttachmentNps struct {
	OriginalFileName *string `json:"originalFileName,omitempty"`
	FilePath *string `json:"filePath,omitempty"`
}


type ApiResultNps struct {
	Code string `json:"code"`
	Message string `json:"message"`
}

type KsReqNps struct {
	SenderKey *string `json:"senderKey,omitempty"`
	SenderKeyType *string `json:"senderKeyType,omitempty"`
	TemplateCode *string `json:"templateCode,omitempty"`
	Comment *string `json:"comment,omitempty"`
}

type KtrResNps struct {
	Code string `json:"code"`
	Message *string `json:"message"`
}

type UtReq17 struct {
	SenderKey *string `json:"senderKey,omitempty"`
	SenderKeyType *string `json:"senderKeyType,omitempty"`
	TemplateCode *string `json:"templateCode,omitempty"`

	NewSenderKey *string `json:"newSenderKey,omitempty"`
	NewSenderKeyType *string `json:"newSenderKeyType,omitempty"`
	NewTemplateCode *string `json:"newTemplateCode,omitempty"`
	NewTemplateName *string `json:"newTemplateName,omitempty"`
	NewTemplateContent *string `json:"newTemplateContent,omitempty"`
	ButtonType *string `json:"buttonType,omitempty"`
	ButtonName *string `json:"buttonName,omitempty"`
	ButtonUrl *string `json:"buttonUrl,omitempty"`
	Buttons *string `json:"buttons,omitempty"`
}

type UtRes17 struct {
	Code string `json:"code"`
	Data *StKakaoRes `json:"data,omitempty"` 
}

type UtKakaoReq struct {
	SenderKey *string `json:"senderKey,omitempty"`
	SenderKeyType *string `json:"senderKeyType,omitempty"`
	TemplateCode *interface{} `json:"templateCode,omitempty"`

	NewSenderKey *string `json:"newSenderKey,omitempty"`
	NewSenderKeyType *string `json:"newSenderKeyType,omitempty"`
	NewTemplateCode *interface{} `json:"newTemplateCode,omitempty"`
	NewTemplateName *interface{} `json:"newTemplateName,omitempty"`
	NewTemplateMessageType string `json:"newTemplateMessageType,omitempty"`
	NewTemplateEmphasizeType string `json:"newTemplateEmphasizeType,omitempty"`
	NewTemplateContent *interface{} `json:"newTemplateContent,omitempty"`
	NewTemplatePreviewMessage *string `json:"newTemplatePreviewMessage,omitempty"`
	NewTemplateExtra *string `json:"newTemplateExtra,omitempty"`
	NewTemplateImageName *string `json:"newTemplateImageName,omitempty"`
	NewTemplateImageUrl *string `json:"newTemplateImageUrl,omitempty"`
	NewTemplateTitle *string `json:"newTemplateTitle,omitempty"`
	NewTemplateSubtitle *string `json:"newTemplateSubtitle,omitempty"`
	NewTemplateHeader *string `json:"newTemplateHeader,omitempty"`
	NewTemplateItemHighlight *KakaoTemplateItemHighlightNps `json:"newTemplateItemHighlight,omitempty"`
	NewTemplateRepresentLink *KakaoTemplateRepresentLinkNps `json:"newTemplateRepresentLink,omitempty"`
	NewCategoryCode *string `json:"newCategoryCode,omitempty"`
	SecurityFlag *bool `json:"securityFlag,omitempty"`
	Buttons *[]KakaoButtonsNps `json:"buttons,omitempty"`
	QuickReplies *[]KakaoQuickRepliesNps `json:"quickReplies,omitempty"`
}

type CmtRes17 struct {
	Code string `json:"code"`
}


////////////////////////////////////////////////////NPS AREA////////////////////////////////////////////////////










