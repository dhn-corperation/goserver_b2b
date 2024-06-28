package kaocenter

import (
	config "mycs/src/kaoconfig"
	db "mycs/src/kaodatabasepool"
	"database/sql"
	"github.com/gin-gonic/gin"
)

func Get_crypto(c *gin.Context){
	userid := c.Request.Header.Get("userid")
	userip := c.ClientIP()
	var crypto sql.NullString
	err := db.DB.QueryRow("select crypto from DHN_CLIENT_LIST where user_id = '"+userid+"' and ip = '"+userip+"'").Scan(&crypto)
	if err != nil { config.Stdlog.Println(err) }
	if crypto.Valid {
		c.String(200, crypto.String)
	} else {
		c.String(200, "")
	}
	
}