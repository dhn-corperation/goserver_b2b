package ktproc

var account = []map[string]string{
	{
		"apiKey" : "4E5ECFD4F879B9114BC616D8A09DA63F",
		"apiPw" : "a15227985!@",
		"userKey" : "dhn7137985a",
	},
	{
		"apiKey" : "5BC39E42ED3C4F8D9691A5B79F5642CC",
		"apiPw" : "a15227985!@",
		"userKey" : "dhn7137985b",
	},
	{
		"apiKey" : "3B7B8D9C7DD3F23FB866D219437CDEC0",
		"apiPw" : "a15227985!@",
		"userKey" : "dhn7137985c",
	},
	{
		"apiKey" : "AFEC83C51F429B4D11B1FCED6A3EB618",
		"apiPw" : "a15227985!@",
		"userKey" : "dhn7137985d",
	},
	{
		"apiKey" : "0AED21F2E6202FB354A31EF1AB0C7727",
		"apiPw" : "a15227985!@",
		"userKey" : "dhn7137985e",
	},
	{
		"apiKey" : "FB7852BA5CCBD0175E53DF44F51E190B",
		"apiPw" : "a15227985!@",
		"userKey" : "dhn7137985f",
	},
	{
		"apiKey" : "4FF095019BAF6C6A9066DE56896C5ACB",
		"apiPw" : "a15227985!@",
		"userKey" : "dhn7137985g",
	},
	{
		"apiKey" : "3B1E2E6C338DB2A80B179DD8E47C7A16",
		"apiPw" : "a15227985!@",
		"userKey" : "dhn7137985h",
	},
	{
		"apiKey" : "B3486336892CB575F6B04E75E1F6613B",
		"apiPw" : "a15227985!@",
		"userKey" : "dhn7137985i",
	},
	{
		"apiKey" : "AFA5A54C3C3AF87A66E9599BA9718403",
		"apiPw" : "a15227985!@",
		"userKey" : "dhn7137985j",
	},
}

type SendReqTable struct {
	MessageSubType int   		`json:"MessageSubType,omitempty"`
	CallbackNumber string		`json:"CallbackNumber,omitempty"`
	SendNumber string			`json:"SendNumber,omitempty"`
	ReserveType int 			`json:"ReserveType,omitempty"`
	ReserveTime string			`json:"ReserveTime,omitempty"`
	ReserveDTime string			`json:"ReserveDTime,omitempty"`
	CustomMessageID string		`json:"CustomMessageID,omitempty"`
	CDRID string				`json:"CDRID,omitempty"`
	CDRTime string				`json:"CDRTime,omitempty"`
	CallbackURL string			`json:"CallbackURL,omitempty"`
	ConvertType string			`json:"ConvertType,omitempty"`
	KisaOrigCode uint64			`json:"KisaOrigCode,omitempty"`
	Bundle []Bundle				`json:"Bundle,omitempty"`
}

type Bundle struct {
	Seq int 					`json:"Seq,omitempty"`
	Number string				`json:"Number,omitempty"`
	Content string				`json:"Content,omitempty"`
	Attachment []Attachment		`json:"Attachment,omitempty"`
	Subject string				`json:"Subject,omitempty"`
	CallbackURL string			`json:"CallbackURL,omitempty"`
}

type Attachment struct {
	attachID int 				`json:"attachID,omitempty"`
	Path string					`json:"Path,omitempty"`
}

type SendResTable struct {
	MsgID string				`json:"MsgID,omitempty"`
	SendReqTable SendReqTable	`json:"SendReqTable,omitempty"`
	FileParam []string			`json:"ImageParam,omitempty"`
	MessageType string			`json:"MassageType,omitempty"`
	ResCode int 				`json:"ResCode,omitempty"`
	BodyData []byte 			`json:"BodyData,omitempty"`
	Seq int   					`json:"Seq,omitempty"`
}

type SendResDetileTable struct {
	CustomMessageID string		`json:"CustomMessageID,omitempty"`
	Time string					`json:"Time,omitempty"`
	// GrpID int64					`json:"GrpID,omitempty"`
	SubmitTime string			`json:"SubmitTime,omitempty"`
	Result int 					`json:"Result,omitempty"`
	Count int 					`json:"Count,omitempty"`
	JobIDs []JobIDs				`json:"JobIDs,omitempty"`
}

type JobIDs struct {
	Index int 					`json:"Index,omitempty"`
	JobID int64 				`json:"JobID,omitempty"`
}

type SearchReqTable struct {
	JobIDs []int64				`json:"JobIDs,omitempty"`
	SendDay string				`json:"SendDay,omitempty"`
}

type SearchResDatailTable struct {
	JobIDs []SearchResDatailTable 	`json:"JobIDs,omitempty"`
}

type SearchResDatailTable struct {
	ServiceProviderID string 	`json:"ServiceProviderID,omitempty"`
	EndUserID string			`json:"EndUserID,omitempty"`
	Result int 					`json:"Result,omitempty"`
	Time string 				`json:"Time,omitempty"`
	FinishPage int 				`json:"FinishPage,omitempty"`
	Duration int 				`json:"Duration,omitempty"`
	CustomMessageID string 		`json:"CustomMessageID,omitempty"`
	SequenceNumber int 			`json:"SequenceNumber,omitempty"`
	JobID int64 				`json:"JobID,omitempty"`
	GroupID int64 				`json:"GroupID,omitempty"`
	MessageType int 			`json:"MessageType,omitempty"`
	SendNumber string			`json:"SendNumber,omitempty"`
	ReceiveNumber string 		`json:"ReceiveNumber,omitempty"`
	CallbackNumber string 		`json:"CallbackNumber,omitempty"`
	ReplyInfo string 			`json:"ReplyInfo,omitempty"`
	TelconInfo int              `json:"TelconInfo,omitempty"`
	Fee int 					`json:"Fee,omitempty"`
	Rtime string				`json:"Rtime,omitempty"`
	SubmitTime string			`json:"SubmitTime,omitempty"`
	StatusText string			`json:"StatusText,omitempty"`
}