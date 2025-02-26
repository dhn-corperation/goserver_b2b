package kaoresulttable

type ResultTable struct {
	Msgid       []string `json:"msgid"`
	Regdt       string   `json:"regdt"`
}

type ResultStr struct {
	Statuscode int
	BodyData   []byte
	Result     map[string]string
}