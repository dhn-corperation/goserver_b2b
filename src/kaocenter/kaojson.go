package kaocenter

type CenterUser struct {
	BizId  string `json:"bizId"`
	ApiKey string `json:"apiKey"`
}

type CenterResponse struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type TemplateInfo struct {
	BizId          string `json:"bizId,omitempty"`
	ApiKey         string `json:"apiKey,omitempty"`
	SenderKey      string `json:"senderKey,omitempty"`
	SenderKeyType  string `json:"senderKeyType,omitempty"`
	TemplateCode   string `json:"templateCode,omitempty"`
	Message        string `json:"message,omitempty"`
	CategoryCode   string `json:"categoryCode,omitempty"`
	Comment        string `json:"comment,omitempty"`
	Keyword        string `json:"keyword,omitempty"`
	StartDate      string `json:"startDate,omitempty"`
	EndDate        string `json:"endDate,omitempty"`
	TemplateStatus string `json:"templateStatus,omitempty"`
}

type ProfileInfo struct {
	BizId        string `json:"bizId,omitempty"`
	ApiKey       string `json:"apiKey,omitempty"`
	PhoneNumber  string `json:"phoneNumber,omitempty"`
	YellowId     string `json:"yellowId,omitempty"`
	CategoryCode string `json:"categoryCode,omitempty"`
	Token        string `json:"token,omitempty"`
	ChannelKey   string `json:"channelKey,omitempty"`
	SenderKey    string `json:"senderKey,omitempty"`
}

type ProfileInfoMulti struct {
	BizId     string   `json:"bizId,omitempty"`
	ApiKey    string   `json:"apiKey,omitempty"`
	SenderKey []string `json:"senderKey,omitempty"`
}

type GroupInfo struct {
	BizId     string `json:"bizId,omitempty"`
	ApiKey    string `json:"apiKey,omitempty"`
	GroupKey  string `json:"groupKey,omitempty"`
	SenderKey string `json:"senderKey,omitempty"`
}

type TCTemplateItemHighlight struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	ImageURL    string `json:"imageUrl,omitempty"`
}

type TCTemplateItem struct {
	List    *[]TCItemList  `json:"list,omitempty"`
	Summary *TCItemSummary `json:"summary,omitempty"`
}

type TCItemList struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type TCItemSummary struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type TCTemplateRepresentLink struct {
	LinkAnd string `json:"linkAnd,omitempty"`
	LinkIos string `json:"linkIos,omitempty"`
	LinkPc  string `json:"linkPc,omitempty"`
	LinkMo  string `json:"linkMo,omitempty"`
}

type TCButtons struct {
	Name      string `json:"name,omitempty"`
	LinkType  string `json:"linkType,omitempty"`
	Ordering  string `json:"ordering,omitempty"`
	LinkMo    string `json:"linkMo,omitempty"`
	LinkPc    string `json:"linkPc,omitempty"`
	LinkAnd   string `json:"linkAnd,omitempty"`
	LinkIos   string `json:"linkIos,omitempty"`
	PluginID  string `json:"pluginId,omitempty"`
	BizFormID int    `json:"bizFormId,omitempty"`
}

type TCQuickReplies struct {
	Name      string `json:"name,omitempty"`
	LinkType  string `json:"linkType,omitempty"`
	LinkMo    string `json:"linkMo,omitempty"`
	LinkPc    string `json:"linkPc,omitempty"`
	LinkAnd   string `json:"linkAnd,omitempty"`
	LinkIos   string `json:"linkIos,omitempty"`
	BizFormID int    `json:"bizFormId,omitempty"`
}

type TemplateCreateJson struct {
	SenderKey              string                   `json:"senderKey,omitempty"`
	SenderKeyType          string                   `json:"senderKeyType,omitempty"`
	TemplateCode           string                   `json:"templateCode,omitempty"`
	TemplateName           string                   `json:"templateName,omitempty"`
	TemplateMessageType    string                   `json:"templateMessageType,omitempty"`
	TemplateEmphasizeType  string                   `json:"templateEmphasizeType,omitempty"`
	TemplateContent        string                   `json:"templateContent,omitempty"`
	TemplatePreviewMessage string                   `json:"templatePreviewMessage,omitempty"`
	TemplateExtra          string                   `json:"templateExtra,omitempty"`
	TemplateImageName      string                   `json:"templateImageName,omitempty"`
	TemplateImageURL       string                   `json:"templateImageUrl,omitempty"`
	TemplateTitle          string                   `json:"templateTitle,omitempty"`
	TemplateSubtitle       string                   `json:"templateSubtitle,omitempty"`
	TemplateHeader         string                   `json:"templateHeader,omitempty"`
	TemplateItemHighlight  *TCTemplateItemHighlight `json:"templateItemHighlight,omitempty"`
	TemplateItem           *TCTemplateItem          `json:"templateItem,omitempty"`
	TemplateRepresentLink  *TCTemplateRepresentLink `json:"templateRepresentLink,omitempty"`
	CategoryCode           string                   `json:"categoryCode,omitempty"`
	SecurityFlag           bool                     `json:"securityFlag,omitempty"`
	Buttons                *[]TCButtons             `json:"buttons,omitempty"`
	QuickReplies           *[]TCQuickReplies        `json:"quickReplies,omitempty"`
	BizId                  string                   `json:"bizId,omitempty"`
	ApiKey                 string                   `json:"apiKey,omitempty"`
}

type TemplateUpdateJson struct {
	SenderKey             string                   `json:"senderKey,omitempty"`
	SenderKeyType         string                   `json:"senderKeyType,omitempty"`
	TemplateCode          string                   `json:"templateCode,omitempty"`
	NewTemplateCode       string                   `json:"newTemplateCode,omitempty"`
	TemplateName          string                   `json:"templateName,omitempty"`
	TemplateMessageType   string                   `json:"templateMessageType,omitempty"`
	TemplateEmphasizeType string                   `json:"templateEmphasizeType,omitempty"`
	TemplateContent       string                   `json:"templateContent,omitempty"`
	TemplateExtra         string                   `json:"templateExtra,omitempty"`
	TemplateImageName     string                   `json:"templateImageName,omitempty"`
	TemplateImageURL      string                   `json:"templateImageUrl,omitempty"`
	TemplateTitle         string                   `json:"templateTitle,omitempty"`
	TemplateSubtitle      string                   `json:"templateSubtitle,omitempty"`
	TemplateHeader        string                   `json:"templateHeader,omitempty"`
	TemplateItemHighlight *TCTemplateItemHighlight `json:"templateItemHighlight,omitempty"`
	TemplateItem          *TCTemplateItem          `json:"templateItem,omitempty"`
	CategoryCode          string                   `json:"categoryCode,omitempty"`
	SecurityFlag          bool                     `json:"securityFlag,omitempty"`
	Buttons               *[]TCButtons             `json:"buttons,omitempty"`
	QuickReplies          *[]TCQuickReplies        `json:"quickReplies,omitempty"`
	BizId                 string                   `json:"bizId,omitempty"`
	ApiKey                string                   `json:"apiKey,omitempty"`
}

type TemplateUpdateKakao struct {
	SenderKey                 string                   `json:"senderKey,omitempty"`
	SenderKeyType             string                   `json:"senderKeyType,omitempty"`
	TemplateCode              string                   `json:"templateCode,omitempty"`
	NewSenderKey              string                   `json:"newSenderKey,omitempty"`
	NewSenderKeyType          string                   `json:"newSenderKeyType,omitempty"`
	NewTemplateCode           string                   `json:"newTemplateCode,omitempty"`
	NewTemplateName           string                   `json:"newTemplateName,omitempty"`
	NewTemplateMessageType    string                   `json:"newTemplateMessageType,omitempty"`
	NewTemplateEmphasizeType  string                   `json:"newTemplateEmphasizeType,omitempty"`
	NewTemplateContent        string                   `json:"newTemplateContent,omitempty"`
	NewTemplatePreviewMessage string                   `json:"newTemplatePreviewMessage,omitempty"`
	NewTemplateExtra          string                   `json:"newTemplateExtra,omitempty"`
	NewTemplateImageName      string                   `json:"newTemplateImageName,omitempty"`
	NewTemplateImageURL       string                   `json:"newTemplateImageUrl,omitempty"`
	NewTemplateTitle          string                   `json:"newTemplateTitle,omitempty"`
	NewTemplateSubtitle       string                   `json:"newTemplateSubtitle,omitempty"`
	NewTemplateHeader         string                   `json:"newTemplateHeader,omitempty"`
	NewTemplateItemHighlight  *TCTemplateItemHighlight `json:"newTemplateItemHighlight,omitempty"`
	NewTemplateItem           *TCTemplateItem          `json:"newTemplateItem,omitempty"`
	NewTemplateRepresentLink  *TCTemplateRepresentLink `json:"newTemplateRepresentLink,omitempty"`
	NewCategoryCode           string                   `json:"newCategoryCode,omitempty"`
	SecurityFlag              bool                     `json:"securityFlag,omitempty"`
	Buttons                   *[]TCButtons             `json:"buttons,omitempty"`
	QuickReplies              *[]TCQuickReplies        `json:"quickReplies,omitempty"`
}

type TemplateResponse struct {
	Code     string               `json:"code,omitempty"`
	Message  string               `json:"message,omitempty"`
	Template *TemplateReponseJson `json:"data,omitempty"`
}

type TemplateReponseJson struct {
	SenderKey             string `json:"senderKey,omitempty"`
	SenderKeyType         string `json:"senderKeyType,omitempty"`
	TemplateCode          string `json:"templateCode,omitempty"`
	TemplateName          string `json:"templateName,omitempty"`
	TemplateMessageType   string `json:"templateMessageType,omitempty"`
	TemplateEmphasizeType string `json:"templateEmphasizeType,omitempty"`
	InspectionStatus      string `json:"inspectionStatus,omitempty"`
	CreatedAt             string `json:"createdAt,omitempty"`
	ModifiedAt            string `json:"modifiedAt,omitempty"`
	Status                string `json:"status,omitempty"`
	Block                 bool   `json:"block,omitempty"`
	Dormant               bool   `json:"dormant,omitempty"`
	SecurityFlag          bool   `json:"securityFlag,omitempty"`
	CategoryCode          string `json:"categoryCode,omitempty"`
}

type ProfileResponse struct {
	Code    string       `json:"code,omitempty"`
	Message string       `json:"message,omitempty"`
	Profile *ProfileJson `json:"data,omitempty"`
}

type ProfileJson struct {
	SenderKey               string `json:"senderKey,omitempty"`
	UUID                    string `json:"uuid,omitempty"`
	Name                    string `json:"name,omitempty"`
	Status                  string `json:"status,omitempty"`
	Block                   bool   `json:"block,omitempty"`
	Dormant                 bool   `json:"dormant,omitempty"`
	ProfileStatus           string `json:"profileStatus,omitempty"`
	CreatedAt               string `json:"createdAt,omitempty"`
	ModifiedAt              string `json:"modifiedAt,omitempty"`
	CategoryCode            string `json:"categoryCode,omitempty"`
	Alimtalk                bool   `json:"alimtalk,omitempty"`
	Bizchat                 bool   `json:"bizchat,omitempty"`
	Brandtalk               bool   `json:"brandtalk,omitempty"`
	CommitalCompanyName     string `json:"commitalCompanyName,omitempty"`
	ChannelKey              string `json:"channelKey,omitempty"`
	BusinessProfile         bool   `json:"businessProfile,omitempty"`
	BusinessType            string `json:"businessType,omitempty"`
	ProfileSpamLevel        string `json:"profileSpamLevel,omitempty"`
	ProfileMessageSpamLevel string `json:"profileMessageSpamLevel,omitempty"`
	ClearBlockURL           string `json:"clearBlockUrl,omitempty"`
}

type ProfileInfoSenderS struct {
	SenderKey           string      `json:"senderKey,omitempty"`
	UUID                string      `json:"uuid,omitempty"`
	Name                string      `json:"name,omitempty"`
	Status              string      `json:"status,omitempty"`
	Block               string      `json:"block,omitempty"`
	Dormant             string      `json:"dormant,omitempty"`
	ProfileStatus       string      `json:"profileStatus,omitempty"`
	CreatedAt           string      `json:"createdAt,omitempty"`
	ModifiedAt          string      `json:"modifiedAt,omitempty"`
	CategoryCode        string      `json:"categoryCode,omitempty"`
	Alimtalk            string      `json:"alimtalk,omitempty"`
	Bizchat             string      `json:"bizchat,omitempty"`
	Brandtalk           string      `json:"brandtalk,omitempty"`
	CommitalCompanyName string      `json:"commitalCompanyName,omitempty"`
	ChannelKey          string      `json:"channelKey,omitempty"`
	BusinessProfile     string      `json:"businessProfile,omitempty"`
	BusinessType        string      `json:"businessType,omitempty"`
	Groups              []GroupJson `json:"groups,omitempty"`
}

type ProfileInfoSender struct {
	Code    string          `json:"code,omitempty"`
	Message string          `json:"message,omitempty"`
	Data    ProfileInfoData `json:"data,omitempty"`
}

type ProfileInfoSenderF struct {
	Code      string          `json:"code,omitempty"`
	Message   string          `json:"message,omitempty"`
	SenderKey ProfileInfoData `json:"senderKey,omitempty"`
}

type ProfileInfoData struct {
	Success []ProfileInfoSenderS `json:"success,omitempty"`
	Fail    []ProfileInfoSenderF `json:"fail,omitempty"`
}

type GroupResponse struct {
	Code    string      `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Group   []GroupJson `json:"data,omitempty"`
}

type GroupJson struct {
	GroupKey  string `json:"groupKey,omitempty"`
	Name      string `json:"name,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
}

type GroupProfileResponse struct {
	Code    string             `json:"code,omitempty"`
	Message string             `json:"message,omitempty"`
	Profile []GroupProfileJson `json:"data,omitempty"`
}

type GroupProfileJson struct {
	SenderKey     string `json:"senderKey,omitempty"`
	UUID          string `json:"uuid,omitempty"`
	Name          string `json:"name,omitempty"`
	Status        string `json:"status,omitempty"`
	Block         bool   `json:"block,omitempty"`
	Dormant       bool   `json:"dormant,omitempty"`
	ProfileStatus string `json:"profileStatus,omitempty"`
	CreatedAt     string `json:"createdAt,omitempty"`
	ModifiedAt    string `json:"modifiedAt,omitempty"`
	CategoryCode  string `json:"categoryCode,omitempty"`
}

type TemplateList_ struct {
	Message     string             `json:"message,omitempty"`
	TotalCount  int                `json:"totalCount,omitempty"`
	TotalPage   int                `json:"totalPage,omitempty"`
	CurrentPage int                `json:"currentPage,omitempty"`
	Data        []TemplateListJson `json:"data,omitempty"`
}

type TemplateListJson struct {
	SenderKey     string `json:"senderKey,omitempty"`
	SenderKeyType string `json:"senderKeyType,omitempty"`
	TemplateCode  string `json:"templateCode,omitempty"`
	TemplateName  string `json:"templateName,omitempty"`
	CategoryCode  string `json:"categoryCode,omitempty"`
	CreatedAt     string `json:"createdAt,omitempty"`
	ModifiedAt    string `json:"modifiedAt,omitempty"`
	Status        string `json:"serviceStatus,omitempty"`
}
