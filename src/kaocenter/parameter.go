package kaocenter

type SenderCreate struct {
	YellowId             string `json:"yellowId" binding:"required"`
	CategoryCode         string `json:"categoryCode" binding:"required"`
	ChannelKey           string `json:"channelKey,omitempty"`
	Bizchat              string `json:"bizchat,omitempty"`
	CommittalCompanyName string `json:"committalCompanyName,omitempty"`
}

type SenderDelete struct {
	SenderKey string `json:"senderKey" binding:"required"`
}

type TemplateCreate struct {
	SenderKey             string `json:"senderKey" binding:"required"`
	TemplateCode          string `json:"templateCode" binding:"required"`
	TemplateName          string `json:"templateName" binding:"required"`
	TemplateMessageType   string `json:"templateMessageType" binding:"required"`
	TemplateEmphasizeType string `json:"templateEmphasizeType" binding:"required"`
	TemplateContent       string `json:"templateContent" binding:"required"`
	TemplateExtra         string `json:"templateExtra,omitempty"`
	//TemplateAd             string                 `json:"templateAd,omitempty"`
	TemplateImageName      string                 `json:"templateImageName,omitempty"`
	TemplateImageUrl       string                 `json:"templateImageUrl,omitempty"`
	TemplateTitle          string                 `json:"templateTitle,omitempty"`
	TemplateSubtitle       string                 `json:"templateSubtitle,omitempty"`
	TemplateHeader         string                 `json:"templateHeader,omitempty"`
	TemplateItemHighlight  TemplateItemHighlights `json:"templateItemHighlight,omitempty"`
	TemplateItem           TemplateItems          `json:"templateItem,omitempty"`
	SenderKeyType          string                 `json:"senderKeyType,omitempty"`
	CategoryCode           string                 `json:"categoryCode,omitempty"`
	SecurityFlag           bool                   `json:"securityFlag,omitempty"`
	Buttons                []Button               `json:"buttons,omitempty"`
	QuickReplies           []Quickreply           `json:"quickReplies,omitempty"`
	TemplatePreviewMessage string                 `json:"templatePreviewMessage,omitempty"`
	TemplateRepresentLink  templateRepresentLinks `json:"templateRepresentLink,omitempty"`
}

type Button struct {
	Name      string `json:"name"`
	LinkType  string `json:"linkType"`
	Ordering  int    `json:"ordering,omitempty"`
	LinkMo    string `json:"linkMo,omitempty"`
	LinkPc    string `json:"linkPc,omitempty"`
	LinkAnd   string `json:"linkAnd,omitempty"`
	LinkIos   string `json:"linkIos,omitempty"`
	PluginId  string `json:"pluginId,omitempty"`
	BizFormId int    `json:"bizFormId,omitempty"`
}

type Quickreply struct {
	Name      string `json:"name"`
	LinkType  string `json:"linkType"`
	LinkMo    string `json:"linkMo,omitempty"`
	LinkPc    string `json:"linkPc,omitempty"`
	LinkAnd   string `json:"linkAnd,omitempty"`
	LinkIos   string `json:"linkIos,omitempty"`
	BizFormId int    `json:"bizFormId,omitempty"`
}

type TemplateRequest struct {
	SenderKey     string `json:"senderKey" binding:"required"`
	TemplateCode  string `json:"templateCode" binding:"required"`
	SenderKeyType string `json:"senderKeyType,omitempty"`
}

type TemplateUpdate struct {
	SenderKey                string `json:"senderKey" binding:"required"`
	TemplateCode             string `json:"templateCode" binding:"required"`
	SenderKeyType            string `json:"senderKeyType,omitempty"`
	NewSenderKey             string `json:"newSenderKey" binding:"required"`
	NewTemplateCode          string `json:"newTemplateCode" binding:"required"`
	NewTemplateName          string `json:"newTemplateName" binding:"required"`
	NewTemplateMessageType   string `json:"newTemplateMessageType" binding:"required"`
	NewTemplateEmphasizeType string `json:"newTemplateEmphasizeType" binding:"required"`
	NewTemplateContent       string `json:"newTemplateContent" binding:"required"`
	NewTemplateExtra         string `json:"newTemplateExtra,omitempty"`
	//NewTemplateAd            string       `json:"newTemplateAd,omitempty"`
	NewTemplateTitle          string                 `json:"newTemplateTitle,omitempty"`
	NewTemplateSubtitle       string                 `json:"newTemplateSubtitle,omitempty"`
	NewSenderKeyType          string                 `json:"newSenderKeyType,omitempty"`
	NewCategoryCode           string                 `json:"newCategoryCode,omitempty"`
	SecurityFlag              bool                   `json:"securityFlag,omitempty"`
	Buttons                   []Button               `json:"buttons,omitempty"`
	QuickReplies              []Quickreply           `json:"quickReplies,omitempty"`
	NewTemplatePreviewMessage string                 `json:"newTemplatePreviewMessage,omitempty"`
	NewTemplateImageName      string                 `json:"newTemplateImageName,omitempty"`
	NewTemplateImageUrl       string                 `json:"newTemplateImageUrl,omitempty"`
	NewTemplateHeader         string                 `json:"newTemplateHeader,omitempty"`
	NewTemplateItemHighlight  TemplateItemHighlights `json:"newTemplateItemHighlight,omitempty"`
	NewTemplateItem           TemplateItems          `json:"newTemplateItem,omitempty"`
	NewTemplateRepresentLink  templateRepresentLinks `json:"newTemplateRepresentLink,omitempty"`
}

type TemplateItemHighlights struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	ImageUrl    string `json:"imageUrl,omitempty"`
}

type TemplateItems struct {
	List    []TemplateItemList `json:"list,omitempty"`
	Summary TemplateItemList   `json:"summary,omitempty"`
}

type TemplateItemList struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type TemplateComment struct {
	SenderKey     string `json:"senderKey" binding:"required"`
	TemplateCode  string `json:"templateCode" binding:"required"`
	SenderKeyType string `json:"senderKeyType,omitempty"`
	Comment       string `json:"comment,omitempty"`
}

type TemplateCategoryUpdate struct {
	SenderKey     string `json:"senderKey" binding:"required"`
	TemplateCode  string `json:"templateCode" binding:"required"`
	SenderKeyType string `json:"senderKeyType,omitempty"`
	CategoryCode  string `json:"comment" binding:"required"`
}

type GroupSenderAdd struct {
	GroupKey  string `json:"groupKey" binding:"required"`
	SenderKey string `json:"senderKey" binding:"required"`
}

type ChannelCreate struct {
	ChannelKey string `json:"channelKey" binding:"required"`
	Desc       string `json:"desc,omitempty"`
}

// 기존 코드 json:"groupKey" 로 되어있어서 파라미터 오류 나옴
/*
type ChannelSenders struct {
	ChannelKey string `json:"groupKey" binding:"required"`
	SenderKeys string `json:"senderKeys" binding:"required"`
}
*/

type ChannelSenders struct {
	ChannelKey string `json:"channelKey" binding:"required"`
	SenderKeys string `json:"senderKeys" binding:"required"`
}

// 위와 같은 문제
/*
type ChannelDelete struct {
	ChannelKey string `json:"groupKey" binding:"required"`
}
*/

type ChannelDelete struct {
	ChannelKey string `json:"channelKey" binding:"required"`
}

type PluginCallbackUrlCreate struct {
	SenderKey   string `json:"senderKey" binding:"required"`
	PluginType  string `json:"pluginType" binding:"required"`
	PluginId    string `json:"pluginId" binding:"required"`
	CallbackUrl string `json:"callbackUrl" binding:"required"`
}

type PluginCallbackUrlUpdate struct {
	SenderKey   string `json:"senderKey" binding:"required"`
	PluginId    string `json:"pluginId" binding:"required"`
	CallbackUrl string `json:"callbackUrl" binding:"required"`
}

type PluginCallbackUrlDelete struct {
	SenderKey string `json:"senderKey" binding:"required"`
	PluginId  string `json:"pluginId" binding:"required"`
}

// 추가 구조체

type Bizform_upload struct {
	BizFormId string `json:"bizFormId" binding:"required"`
	SenderKey string `json:"senderKey" binding:"required"`
}

type Ft_possible struct {
	Sender_key    string `json:"sender_key" binding:"required"`
	Phone_numbers string `json:"phone_numbers" binding:"required"`
}

type Direct_convert struct {
	KakaoAccount string `json:"kakaoAccount" binding:"required"`
	SenderKey    string `json:"senderKey" binding:"required"`
	BizWalletId  int    `json:"bizWalletId" binding:"required"`
}

type Template_convertAddCh struct {
	SenderKey    string `json:"senderKey" binding:"required"`
	TemplateCode string `json:"templateCode" binding:"required"`
}

type Channel_sender struct {
	ChannelKey string `json:"channelKey" binding:"required"`
	SenderKey  string `json:"senderKey" binding:"required"`
}

// 그룹 태그 관련

type Group_Tag_create struct {
	SenderKey    string `json:"senderKey" binding:"required"`
	GroupTagName string `json:"groupTagName" binding:"required"`
}

type Group_Tag_update struct {
	SenderKey       string `json:"senderKey" binding:"required"`
	GroupTagKey     string `json:"groupTagKey" binding:"required"`
	NewGroupTagName string `json:"newGroupTagName" binding:"required"`
}

type Group_Tag_delete struct {
	SenderKey   string `json:"senderKey" binding:"required"`
	GroupTagKey string `json:"groupTagKey" binding:"required"`
}

type Direct_template_create struct {
	SenderKey         string       `json:"senderKey,omitempty"`
	SenderGroupKey    string       `json:"senderGroupKey,omitempty"`
	Name              string       `json:"name" binding:"required"`
	MessageType       string       `json:"messageType" binding:"required"`
	Content           string       `json:"content,omitempty"`
	Adult             bool         `json:"adult"`
	ImageUrl          string       `json:"imageUrl,omitempty"`
	ImageName         string       `json:"imageName,omitempty"`
	ImageLink         string       `json:"imageLink,omitempty"`
	Header            string       `json:"header,omitempty"`
	AdditionalContent string       `json:"additionalContent,omitempty"`
	Carousel          Carousels    `json:"carousel,omitempty"`
	MainWideItem      WideItems    `json:"mainWideItem,omitempty"`
	SubWideItemList   []WideItems  `json:"subWideItemList,omitempty"`
	Video             Videos       `json:"video,omitempty"`
	Commerce          Commerces    `json:"commerce,omitempty"`
	Buttons           []Quickreply `json:"buttons,omitempty"`
	Coupon            Coupons      `json:"coupon,omitempty"`
}

type Carousels struct {
	Head Heads   `json:"head,omitempty"`
	List []Lists `json:"list,omitempty"`
	Tail Tails   `json:"tail,omitempty"`
}

type Heads struct {
	Header    string `json:"header,omitempty"`
	Content   string `json:"content,omitempty"`
	ImageURL  string `json:"imageUrl,omitempty"`
	ImageName string `json:"imageName,omitempty"`
	LinkMo    string `json:"linkMo,omitempty"`
	LinkPc    string `json:"linkPc,omitempty"`
	LinkAnd   string `json:"linkAnd,omitempty"`
	LinkIos   string `json:"linkIos,omitempty"`
}

type Lists struct {
	Header            string       `json:"header,omitempty"`
	Content           string       `json:"content,omitempty"`
	AdditionalContent string       `json:"additionalContent,omitempty"`
	ImageURL          string       `json:"imageUrl,omitempty"`
	ImageName         string       `json:"imageName,omitempty"`
	ImageLink         string       `json:"imageLink,omitempty"`
	Commerce          Commerces    `json:"commerce,omitempty"`
	Buttons           []Quickreply `json:"buttons,omitempty"`
	Coupon            Coupons      `json:"coupon,omitempty"`
}

type Tails struct {
	LinkMo  string `json:"linkMo,omitempty"`
	LinkPc  string `json:"linkPc,omitempty"`
	LinkAnd string `json:"linkAnd,omitempty"`
	LinkIos string `json:"linkIos,omitempty"`
}

type WideItems struct {
	Title     string `json:"title,omitempty"`
	ImageUrl  string `json:"imageUrl,omitempty"`
	ImageName string `json:"imageName,omitempty"`
	LinkMo    string `json:"linkMo,omitempty"`
	LinkPc    string `json:"linkPc,omitempty"`
	LinkAnd   string `json:"linkAnd,omitempty"`
	LinkIos   string `json:"linkIos,omitempty"`
}

type Videos struct {
	VideoUrl     string `json:"videoUrl,omitempty"`
	ThumbnailUrl string `json:"thumbnailUrl,omitempty"`
}

type Commerces struct {
	Title         string `json:"title,omitempty"`
	RegularPrice  int    `json:"regularPrice,omitempty"`
	DiscountPrice int    `json:"discountPrice,omitempty"`
	DiscountRate  int    `json:"discountRate,omitempty"`
	DiscountFixed int    `json:"discountFixed,omitempty"`
}

/*
type Button_Dir struct {
	Name     string  `json:"name" binding:"required"`
	LinkType string  `json:"linkType" binding:"required"`
	LinkMo   *string `json:"linkMo,omitempty"`
	LinkPc   *string `json:"linkPc,omitempty"`
	LinkAnd  *string `json:"linkAnd,omitempty"`
	LinkIos  *string `json:"linkIos,omitempty"`
}
*/

type Coupons struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	LinkMo      string `json:"linkMo,omitempty"`
	LinkPc      string `json:"linkPc,omitempty"`
	LinkAnd     string `json:"linkAnd,omitempty"`
	LinkIos     string `json:"linkIos,omitempty"`
}

type Direct_bizWallet_change struct {
	KakaoAccount      string `json:"kakaoAccount" binding:"required"`
	SenderKey         string `json:"senderKey" binding:"required"`
	TargetBizWalletId int    `json:"targetBizWalletId" binding:"required"`
}

type templateRepresentLinks struct {
	LinkAnd string `json:"linkAnd,omitempty"`
	LinkIos string `json:"linkIos,omitempty"`
	LinkPc  string `json:"linkPc,omitempty"`
	LinkMo  string `json:"linkMo,omitempty"`
}