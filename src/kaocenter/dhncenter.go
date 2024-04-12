package kaocenter

import (
	config "mycs/src/kaoconfig"
	db "mycs/src/kaodatabasepool"
	"database/sql"
)

var centerClient *http.Client = &http.Client{
	Timeout: time.Second * 30,
	Transport: &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	},
}

var conf = config.Conf

func Get_crypto(c *gin.Context){
	userid := c.Request.Header.Get("userid")
	userip := c.ClientIP()
	var crypto sql.NullString
	err := db.DB.QueryRow("select crypto from DHN_CLIENT_LIST where user_id = '"+userid+"' and ip = '"+userip+"'").Scan(&crypto)
	if err != nil { conf.Println(err) }
	if crypto.Valid {
		c.String(200, crypto)
	} else {
		c.String(404, "조건에 맞는 결과값이 없습니다.")
	}
	
}