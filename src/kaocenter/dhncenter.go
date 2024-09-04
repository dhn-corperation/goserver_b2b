package kaocenter

import (
	config "mycs/src/kaoconfig"
	db "mycs/src/kaodatabasepool"
	"database/sql"
	"github.com/valyala/fasthttp"
)

func Get_crypto(c *fasthttp.RequestCtx){
	userid := string(c.Request.Header.Peek("userid"))
	userip := c.RemoteIP().String()
	var crypto sql.NullString
	err := db.DB.QueryRow("select crypto from DHN_CLIENT_LIST where user_id = '"+userid+"' and ip = '"+userip+"'").Scan(&crypto)
	if err != nil { config.Stdlog.Println(err) }
	if crypto.Valid {
		c.SetStatusCode(fasthttp.StatusOK)
		c.SetBodyString(crypto.String)
	} else {
		c.SetStatusCode(fasthttp.StatusOK)
		c.SetBodyString("")
	}
	
}