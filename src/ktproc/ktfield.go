package ktproc

type SendReqTable struct {
	MessageSubType int   		`json:"MessageSubType,omitempty"`
	CallbackNumber string		`json:"CallbackNumber,omitempty"`
	SendNumber string			`json:"SendNumber,omitempty"`
	ReserveType int 			`json:"ReserveType,omitempty"`
	ReserveTime string			`json:"ReserveTime,omitempty"`
	ReserveDTime string			`json:"ReserveDTime,omitempty"`
	CustomMessageID string		`json:"Custom MessageID,omitempty"`
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
	ResCode int
	BodyData []byte
}

type SendResDetileTable struct {
	CustomMessageID string		`json:"CustomMessageID,omitempty"`
	Time string					`json:"Time,omitempty"`
	GrpID int64					`json:"GrpID,omitempty"`
	SubmitTime string			`json:"SubmitTime,omitempty"`
	Result int 					`json:"Result,omitempty"`
	Count int 					`json:"Count,omitempty"`
	JobIDs []JobIDs				`json:"JobIDs,omitempty"`
}

type JobIDs struct {
	Index int 					`json:"Index,omitempty"`
	JobID int64 				`json:"JobID,omitempty"`
}