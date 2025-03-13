package dhnm

import (
	"database/sql"
	"github.com/valyala/fasthttp"
	
	config "mycs/configs"
	db "mycs/configs/databasepool"
)

func get_crypto(c *fasthttp.RequestCtx){
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