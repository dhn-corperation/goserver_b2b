package kaocenter

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"
	"database/sql"

	config "mycs/src/kaoconfig"
	db "mycs/src/kaodatabasepool"
	cm "mycs/src/kaocommon"
	kj "mycs/src/kakaojson"

	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"github.com/goccy/go-json"
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

func Sender_token(c *fasthttp.RequestCtx) {
	conf := config.Conf

	yellowId := string(c.QueryArgs().Peek("yellowId"))
	phoneNumber := string(c.QueryArgs().Peek("phoneNumber"))

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/sender/token?yellowId="+yellowId+"&phoneNumber="+phoneNumber, nil)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Category_all(c *fasthttp.RequestCtx) {
	conf := config.Conf

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/category/all", nil)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Category_(c *fasthttp.RequestCtx) {
	conf := config.Conf

	categoryCode := string(c.QueryArgs().Peek("categoryCode"))
	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/category?categoryCode="+categoryCode, nil)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Sender_Create(c *fasthttp.RequestCtx) {
	conf := config.Conf

	token := string(c.Request.Header.Peek("token"))
	phoneNumber := string(c.Request.Header.Peek("phoneNumber"))

	buff := bytes.NewBuffer(c.PostBody())
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v3/"+conf.PROFILE_KEY+"/sender/create", buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("token", token)
	req.Header.Add("phoneNumber", phoneNumber)

	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Sender_(c *fasthttp.RequestCtx) {
	conf := config.Conf

	senderKey := string(c.QueryArgs().Peek("senderKey"))

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v3/"+conf.PROFILE_KEY+"/sender?senderKey="+senderKey, nil)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Sender_Delete(c *fasthttp.RequestCtx) {
	conf := config.Conf

	buff := bytes.NewBuffer(c.PostBody())
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/sender/delete", buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Sender_Recover(c *fasthttp.RequestCtx) {
	conf := config.Conf

	buff := bytes.NewBuffer(c.PostBody())
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/sender/recover", buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Template_Create(c *fasthttp.RequestCtx) {
	conf := config.Conf

	buff := bytes.NewBuffer(c.PostBody())
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/create", buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Template_(c *fasthttp.RequestCtx) {
	conf := config.Conf

	senderKey := string(c.QueryArgs().Peek("senderKey"))
	templateCode := string(c.QueryArgs().Peek("templateCode"))
	senderKeyType := string(c.QueryArgs().Peek("senderKeyType"))

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template?senderKey="+senderKey+"&templateCode="+url.QueryEscape(templateCode)+"&senderKeyType="+senderKeyType, nil)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	
	req.Header.Add("Accept-Charset", "utf-8")
	
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Template_Request(c *fasthttp.RequestCtx) {
	conf := config.Conf

	buff := bytes.NewBuffer(c.PostBody())
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/request", buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Template_Cancel_Request(c *fasthttp.RequestCtx) {
	conf := config.Conf

	buff := bytes.NewBuffer(c.PostBody())
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/cancel_request", buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Template_Update(c *fasthttp.RequestCtx) {
	conf := config.Conf
	
	buff := bytes.NewBuffer(c.PostBody())
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/update", buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Template_Stop(c *fasthttp.RequestCtx) {
	conf := config.Conf

	buff := bytes.NewBuffer(c.PostBody())
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/stop", buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Template_Reuse(c *fasthttp.RequestCtx) {
	conf := config.Conf

	buff := bytes.NewBuffer(c.PostBody())
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/reuse", buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Template_Delete(c *fasthttp.RequestCtx) {
	conf := config.Conf

	buff := bytes.NewBuffer(c.PostBody())
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/delete", buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Template_Last_Modified(c *fasthttp.RequestCtx) {
	conf := config.Conf

	senderKey := string(c.QueryArgs().Peek("senderKey"))
	senderKeyType := string(c.QueryArgs().Peek("senderKeyType"))
	since := string(c.QueryArgs().Peek("since"))
	page := string(c.QueryArgs().Peek("page"))
	count := string(c.QueryArgs().Peek("count"))

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v3/"+conf.PROFILE_KEY+"/alimtalk/template/last_modified?senderKey="+senderKey+"&senderKeyType="+senderKeyType+"&since="+since+"&page="+page+"&count="+count, nil)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	resp, err := centerClient.Do(req)

	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Template_Category_all(c *fasthttp.RequestCtx) {
	conf := config.Conf

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/category/all", nil)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Template_Category_(c *fasthttp.RequestCtx) {
	conf := config.Conf

	categoryCode := string(c.QueryArgs().Peek("categoryCode"))

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/category?categoryCode="+categoryCode, nil)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	resp, err := centerClient.Do(req)

	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Template_Dormant_Release(c *fasthttp.RequestCtx) {
	conf := config.Conf

	buff := bytes.NewBuffer(c.PostBody())
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/dormant/release", buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)

	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Group_(c *fasthttp.RequestCtx) {
	conf := config.Conf

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/group", nil)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Group_Sender(c *fasthttp.RequestCtx) {
	conf := config.Conf

	groupKey := string(c.QueryArgs().Peek("groupKey"))

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v3/"+conf.PROFILE_KEY+"/group/sender?groupKey="+groupKey, nil)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Group_Sender_Add(c *fasthttp.RequestCtx) {
	conf := config.Conf

	buff := bytes.NewBuffer(c.PostBody())
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/group/sender/add", buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Group_Sender_Remove(c *fasthttp.RequestCtx) {
	conf := config.Conf

	buff := bytes.NewBuffer(c.PostBody())
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/group/sender/remove", buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)

	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Channel_Create_(c *fasthttp.RequestCtx) {
	conf := config.Conf

	buff := bytes.NewBuffer(c.PostBody())
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/channel/create", buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Channel_all(c *fasthttp.RequestCtx) {
	conf := config.Conf

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/channel/all", nil)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Channel_(c *fasthttp.RequestCtx) {
	conf := config.Conf

	channelKey := string(c.QueryArgs().Peek("channelKey"))

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/channel?channelKey="+channelKey, nil)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Channel_Update_(c *fasthttp.RequestCtx) {
	conf := config.Conf

	buff := bytes.NewBuffer(c.PostBody())
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/channel/update", buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Channel_Senders_(c *fasthttp.RequestCtx) {
	conf := config.Conf

	buff := bytes.NewBuffer(c.PostBody())
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/channel/senders", buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Channel_Delete_(c *fasthttp.RequestCtx) {
	conf := config.Conf

	buff := bytes.NewBuffer(c.PostBody())
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/channel/delete", buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Plugin_CallbackUrls_List(c *fasthttp.RequestCtx) {
	conf := config.Conf

	senderKey := string(c.QueryArgs().Peek("senderKey"))

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/plugin/callbackUrl/list?senderKey="+senderKey, nil)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Plugin_callbackUrl_Create(c *fasthttp.RequestCtx) {
	conf := config.Conf

	buff := bytes.NewBuffer(c.PostBody())
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/plugin/callbackUrl/create", buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Plugin_callbackUrl_Update(c *fasthttp.RequestCtx) {
	conf := config.Conf

	buff := bytes.NewBuffer(c.PostBody())
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/plugin/callbackUrl/update", buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()
	defer resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Plugin_callbackUrl_Delete(c *fasthttp.RequestCtx) {
	conf := config.Conf

	buff := bytes.NewBuffer(c.PostBody())
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/plugin/callbackUrl/delete", buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func FT_Upload(c *fasthttp.RequestCtx) {
	conf := config.Conf

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)

	form, err := c.MultipartForm()
	if err != nil {
		c.SetBody(cm.SetCenterResult("error", "MultipartForm 파싱 실패"))
		return
	}

	if form.File == nil || form.File["image"] == nil{
		c.SetBody(cm.SetCenterResult("error", "업로드 파일 존재하지 않음"))
		return
	}

	files := form.File["image"]

	if len(files) == 0 {
		c.SetBody(cm.SetCenterResult("error", "이미지 파일 존재하지 않음"))
		return
	}

	file := files[0]
	extension := filepath.Ext(file.Filename)
	newFileName := uuid.New().String() + extension

	err = saveUploadedFile(file, config.BasePath+"upload/"+newFileName)
	if err != nil {
		c.SetBody(cm.SetCenterResult("error", "이미지 파일 업로드 실패"))
		return
	}

	param := map[string]io.Reader{
		"image": mustOpen(config.BasePath+"upload/" +newFileName),
	}

	resp, err := upload(conf.IMAGE_SERVER+"v1/"+conf.PROFILE_KEY+"/image/friendtalk", param)
	if err != nil {
		c.SetBody(cm.SetCenterResult("error", "카카오톡 api 전송 실패"))
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}


func FT_Wide_Upload(c *fasthttp.RequestCtx) {
	conf := config.Conf

	form, err := c.MultipartForm()
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
	}

	files := form.File["image"]
	if len(files) == 0 {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
	}

	file := files[0]

	extension := filepath.Ext(file.Filename)
	newFileName := uuid.New().String() + extension

	err = saveUploadedFile(file, config.BasePath+"upload/" + newFileName)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	param := map[string]io.Reader{
		"image": mustOpen(config.BasePath+"upload/" + newFileName),
	}

	resp, err := upload(conf.IMAGE_SERVER+"v1/"+conf.PROFILE_KEY+"/image/friendtalk/wide", param)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func AT_Image(c *fasthttp.RequestCtx) {
	conf := config.Conf

	form, err := c.MultipartForm()
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
	}

	files := form.File["image"]
	if len(files) == 0 {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
	}

	file := files[0]

	extension := filepath.Ext(file.Filename)
	newFileName := uuid.New().String() + extension

	err = saveUploadedFile(file, config.BasePath+"upload/" + newFileName)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	param := map[string]io.Reader{
		"image": mustOpen(config.BasePath+"upload/" + newFileName),
	}

	resp, err := upload(conf.IMAGE_SERVER+ "v1/"+conf.PROFILE_KEY+"/image/alimtalk/template", param)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func MMS_Image(c *fasthttp.RequestCtx) {
	//conf := config.Conf
	var newFileName1, newFileName2, newFileName3 string
	imageKeys := []string{"image1", "image2", "image3"}
	userID := string(c.FormValue("userid"))

	form, err := c.MultipartForm()
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
	}
	
	var startNow = time.Now()
	var group_no = fmt.Sprintf("%04d%02d%02d%02d%02d%02d%09d", startNow.Year(), startNow.Month(), startNow.Day(), startNow.Hour(), startNow.Minute(), startNow.Second(), startNow.Nanosecond())
	
	seq := 1
	for _, key := range imageKeys {
		files := form.File[key]
		if len(files) != 0 {
			extension := filepath.Ext(files[0].Filename)
			nfn := config.BasePath+"upload/mms/" + uuid.New().String() + extension
			switch seq {
			case 1:
				newFileName1 = nfn
			case 2:
				newFileName2 = nfn
			case 3:
				newFileName3 = nfn
			}
			err := saveUploadedFile(files[0], nfn)
			if err != nil {
				config.Stdlog.Println("File ", seq," 저장 오류 : ", newFileName1, err)
				switch seq {
				case 1:
					newFileName1 = ""
				case 2:
					newFileName2 = ""
				case 3:
					newFileName3 = ""
				}
			}
		} else {
			switch seq {
			case 1:
				newFileName1 = ""
			case 2:
				newFileName2 = ""
			case 3:
				newFileName3 = ""
			}
		}
		seq++
	}
 
	if len(newFileName1) > 0 || len(newFileName2) > 0 || len(newFileName3) > 0  {
	
		mmsinsQuery := `insert IGNORE into api_mms_images(
  user_id,
  mms_id,             
  origin1_path,
  origin2_path,
  origin3_path,
  file1_path,
  file2_path,
  file3_path   
) values %s
	`
		mmsinsStrs := []string{}
		mmsinsValues := []interface{}{}
	
		mmsinsStrs = append(mmsinsStrs, "(?,?,null,null,null,?,?,?)")
		mmsinsValues = append(mmsinsValues, userID)
		mmsinsValues = append(mmsinsValues, group_no)
		mmsinsValues = append(mmsinsValues, newFileName1)
		mmsinsValues = append(mmsinsValues, newFileName2)
		mmsinsValues = append(mmsinsValues, newFileName3)
		
		if len(mmsinsStrs) >= 1 {
			stmt := fmt.Sprintf(mmsinsQuery, strings.Join(mmsinsStrs, ","))
			_, err := db.DB.Exec(stmt, mmsinsValues...)
	
			if err != nil {
				config.Stdlog.Println("API MMS Insert 처리 중 오류 발생 " + err.Error())
			}
	
			mmsinsStrs = nil
			mmsinsValues = nil
		} 

		res, err := json.Marshal(map[string]string{
			"image_group":group_no,
		})

		if err != nil {
			c.Error(err.Error(), fasthttp.StatusBadRequest)
			return
		}

		c.SetContentType("application/json")
		c.SetStatusCode(fasthttp.StatusOK)
		c.SetBody(res)
	} else {
		res, _ := json.Marshal(map[string]string{
			"message":"no content",
		})
		c.SetContentType("application/json")
		c.SetStatusCode(fasthttp.StatusBadRequest)
		c.SetBody(res)
	}
}


func Image_wideItemList(c *fasthttp.RequestCtx) {
	conf := config.Conf

	param, err := imageMapping(c, map[string]io.Reader{}, "image_", 4)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	} else if len(param) == 0 {
		c.Error("이미지 저장 실패", fasthttp.StatusBadRequest)
		return
	}

	resp, err := upload(conf.IMAGE_SERVER+ "v1/"+conf.PROFILE_KEY+"/image/friendtalk/wideItemList", param)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)
	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Image_carousel(c *fasthttp.RequestCtx) {
	conf := config.Conf

	param, err := imageMapping(c, map[string]io.Reader{}, "image_", 11)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	} else if len(param) == 0 {
		c.Error("이미지 저장 실패", fasthttp.StatusBadRequest)
		return
	}

	resp, err := upload(conf.IMAGE_SERVER+ "v1/"+conf.PROFILE_KEY+"/image/friendtalk/carousel", param)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)
	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Get_Polling_Id(c *fasthttp.RequestCtx) {
	conf := config.Conf
	respId := c.UserValue("respid").(string)

	buff := bytes.NewBuffer([]byte("{}"))
	req, err := http.NewRequest("POST", conf.API_SERVER+"/v3/"+conf.PROFILE_KEY+"/response/"+respId, buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)
	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)

}

func AT_Highlight_Image(c *fasthttp.RequestCtx) {
	conf := config.Conf

	param, err := imageMapping(c, map[string]io.Reader{}, "image", 0)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	} else if len(param) == 0 {
		c.Error("이미지 저장 실패", fasthttp.StatusBadRequest)
		return
	}

	resp, err := upload(conf.IMAGE_SERVER+"v1/"+conf.PROFILE_KEY+"/image/alimtalk/itemHighlight", param)
	if err != nil {
		config.Stdlog.Println("File upload 오류 : ", err)
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func FT_Carousel_Feed_Image(c *fasthttp.RequestCtx) {
	conf := config.Conf

	param, err := imageMapping(c, map[string]io.Reader{}, "image_", 10)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	} else if len(param) == 0 {
		c.Error("이미지 저장 실패", fasthttp.StatusBadRequest)
		return
	}

	resp, err := upload(conf.IMAGE_SERVER+"v1/"+conf.PROFILE_KEY+"/image/friendtalk/carousel", param)
	if err != nil {
		config.Stdlog.Println("File upload 오류 : ", err)
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func FT_Carousel_Commerce_Image(c *fasthttp.RequestCtx) {
	conf := config.Conf

	param, err := imageMapping(c, map[string]io.Reader{}, "image_", 11)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	} else if len(param) == 0 {
		c.Error("이미지 저장 실패", fasthttp.StatusBadRequest)
		return
	}

	resp, err := upload(conf.IMAGE_SERVER+"v1/"+conf.PROFILE_KEY+"/image/friendtalk/carouselCommerce", param)
	if err != nil {
		config.Stdlog.Println("File upload 오류 : ", err)
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func DM_Default_Image(c *fasthttp.RequestCtx) {
	conf := config.Conf

	param, err := imageMapping(c, map[string]io.Reader{}, "image", 0)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	} else if len(param) == 0 {
		c.Error("이미지 저장 실패", fasthttp.StatusBadRequest)
		return
	}

	resp, err := upload(conf.IMAGE_SERVER+"v2/"+conf.PROFILE_KEY+"/image/default", param)
	if err != nil {
		config.Stdlog.Println("File upload 오류 : ", err)
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func DM_Wide_Image(c *fasthttp.RequestCtx) {
	conf := config.Conf

	param, err := imageMapping(c, map[string]io.Reader{}, "image", 0)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	} else if len(param) == 0 {
		c.Error("이미지 저장 실패", fasthttp.StatusBadRequest)
		return
	}

	resp, err := upload(conf.IMAGE_SERVER+"v2/"+conf.PROFILE_KEY+"/image/wide", param)
	if err != nil {
		config.Stdlog.Println("File upload 오류 : ", err)
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func DM_Widelist_First_image(c *fasthttp.RequestCtx) {
	conf := config.Conf

	param, err := imageMapping(c, map[string]io.Reader{}, "image", 0)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	} else if len(param) == 0 {
		c.Error("이미지 저장 실패", fasthttp.StatusBadRequest)
		return
	}

	resp, err := upload(conf.IMAGE_SERVER+"v2/"+conf.PROFILE_KEY+"/image/wideItemList/first", param)
	if err != nil {
		config.Stdlog.Println("File upload 오류 : ", err)
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func DM_Widelist_Image(c *fasthttp.RequestCtx) {
	conf := config.Conf

	param, err := imageMapping(c, map[string]io.Reader{}, "image_", 3)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	} else if len(param) == 0 {
		c.Error("이미지 저장 실패", fasthttp.StatusBadRequest)
		return
	}

	resp, err := upload(conf.IMAGE_SERVER+"v2/"+conf.PROFILE_KEY+"/image/wideItemList", param)
	if err != nil {
		config.Stdlog.Println("File upload 오류 : ", err)
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func DM_Carousel_Feed_Image(c *fasthttp.RequestCtx) {
	conf := config.Conf

	param, err := imageMapping(c, map[string]io.Reader{}, "image_", 10)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	} else if len(param) == 0 {
		c.Error("이미지 저장 실패", fasthttp.StatusBadRequest)
		return
	}

	resp, err := upload(conf.IMAGE_SERVER+"v2/"+conf.PROFILE_KEY+"/image/carouselFeed", param)
	if err != nil {
		config.Stdlog.Println("File upload 오류 : ", err)
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func DM_Carousel_Commerce_Image(c *fasthttp.RequestCtx) {
	conf := config.Conf

	param, err := imageMapping(c, map[string]io.Reader{}, "image_", 11)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	} else if len(param) == 0 {
		c.Error("이미지 저장 실패", fasthttp.StatusBadRequest)
		return
	}

	resp, err := upload(conf.IMAGE_SERVER+"v2/"+conf.PROFILE_KEY+"/image/carouselCommerce", param)
	if err != nil {
		config.Stdlog.Println("File upload 오류 : ", err)
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

// 친구톡 API
// 별첨1 - 비즈폼 업로드 요청
func Bizform_upload_(c *fasthttp.RequestCtx) {
	conf := config.Conf

	buff := bytes.NewBuffer(c.PostBody())
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/bizform/upload", buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

// 별첨2 - 친구톡 발송 가능 모수 확인
func Ft_possible_(c *fasthttp.RequestCtx) {
	conf := config.Conf

	buff := bytes.NewBuffer(c.PostBody())
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/friendtalk/possible", buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

// 센터 API
// 발신 프로필 조회2 (톡 채널 키로 조회)
func Sender_channel(c *fasthttp.RequestCtx) {
	conf := config.Conf

	talkChannelKey := c.UserValue("talkChannelKey").(string)

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v3/"+conf.PROFILE_KEY+"/sender/"+talkChannelKey, nil)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

// 최근 변경 발신 프로필 조회
func Sender_modified(c *fasthttp.RequestCtx) {
	conf := config.Conf

	args := fasthttp.AcquireArgs()
	defer fasthttp.ReleaseArgs(args)

	args.Add("since", string(c.QueryArgs().Peek("since")))
	args.Add("page", string(c.QueryArgs().Peek("page")))
	args.Add("count", string(c.QueryArgs().Peek("count")))

	queryString := args.String()

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v3/"+conf.PROFILE_KEY+"/sender/last_modified?"+queryString, nil)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

// 검수요청 (파일첨부)
func Template_request_with_file(c *fasthttp.RequestCtx) {
	conf := config.Conf

	param, err := imageMapping(c, map[string]io.Reader{}, "attachment", 0)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	} else if len(param) == 0 {
		c.Error("이미지 저장 실패", fasthttp.StatusBadRequest)
		return
	}

	param["senderKey"] = strings.NewReader(string(c.FormValue("senderKey")))
	param["templateCode"] = strings.NewReader(string(c.FormValue("templateCode")))
	param["senderKeyType"] = strings.NewReader(string(c.FormValue("senderKeyType")))
	param["comment"] = strings.NewReader(string(c.FormValue("comment")))

	resp, err := upload(conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/request_with_file", param)
	if err != nil {
		config.Stdlog.Println("File upload 오류 : ", err)
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

// 검수 승인 취소
func Template_cancel_approval_(c *fasthttp.RequestCtx) {
	conf := config.Conf

	buff := bytes.NewBuffer(c.PostBody())
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/cancel_approval", buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

// 기등록된 템플릿 (타입 : BA, EX) 을 채널추가버튼 및 채널추가안내문구가 포함된 템플릿으로 전환 /template/convertAddCh
func Template_convertAddCh_(c *fasthttp.RequestCtx) {
	conf := config.Conf

	buff := bytes.NewBuffer(c.PostBody())
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/convertAddCh", buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

// 채널에 발신 프로필 추가
func Channel_sender_add_(c *fasthttp.RequestCtx) {
	conf := config.Conf

	buff := bytes.NewBuffer(c.PostBody())
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/channel/sender/add", buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)

}

// 채널에 발신 프로필 삭제
func Channel_sender_remove_(c *fasthttp.RequestCtx) {
	conf := config.Conf

	buff := bytes.NewBuffer(c.PostBody())
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/channel/sender/remove", buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

// 알림톡, 친구톡 발송 일별 통계
func Stat_daily(c *fasthttp.RequestCtx) {
	conf := config.Conf

	args := fasthttp.AcquireArgs()
	defer fasthttp.ReleaseArgs(args)

	args.Add("beginDate", string(c.QueryArgs().Peek("beginDate")))
	args.Add("endDate", string(c.QueryArgs().Peek("endDate")))
	args.Add("productType", string(c.QueryArgs().Peek("productType")))
	args.Add("page", string(c.QueryArgs().Peek("page")))

	queryString := args.String()

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/stat/daily?"+queryString, nil)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

// 그룹 태그 생성
func GroupTag_create(c *fasthttp.RequestCtx) {
	conf := config.Conf

	buff := bytes.NewBuffer(c.PostBody())
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/groupTag/create", buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

// 그룹 태그 조회
func GroupTag_(c *fasthttp.RequestCtx) {
	conf := config.Conf

	args := fasthttp.AcquireArgs()
	defer fasthttp.ReleaseArgs(args)

	args.Add("senderKey", string(c.QueryArgs().Peek("senderKey")))
	args.Add("groupTagKey", string(c.QueryArgs().Peek("groupTagKey")))

	queryString := args.String()

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/groupTag?"+queryString, nil)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

// 그룹 태그 전체 조회
func GroupTag_list(c *fasthttp.RequestCtx) {
	conf := config.Conf

	args := fasthttp.AcquireArgs()
	defer fasthttp.ReleaseArgs(args)

	args.Add("senderKey", string(c.QueryArgs().Peek("senderKey")))

	queryString := args.String()

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/groupTag/list?"+queryString, nil)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

// 그룹 태그 수정
func GroupTag_update(c *fasthttp.RequestCtx) {
	conf := config.Conf

	buff := bytes.NewBuffer(c.PostBody())
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/groupTag/update", buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

// 그룹 태그 삭제
func GroupTag_delete(c *fasthttp.RequestCtx) {
	conf := config.Conf

	buff := bytes.NewBuffer(c.PostBody())
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/groupTag/delete", buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

// 광고성 메시지(다이렉트) 템플릿 등록
func Direct_template_create_(c *fasthttp.RequestCtx) {
	conf := config.Conf

	buff := bytes.NewBuffer(c.PostBody())
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v3/"+conf.PROFILE_KEY+"/direct/template/create", buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

// 광고성메시지(다이렉트) 템플릿 조회
func Direct_template_(c *fasthttp.RequestCtx) {
	conf := config.Conf

	code := c.UserValue("code").(string)

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/direct/template/"+code, nil)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

// 광고성메시지(다이렉트) 템플릿 수정
func Direct_template_update_(c *fasthttp.RequestCtx) {
	conf := config.Conf

	code := c.UserValue("code").(string)

	buff := bytes.NewBuffer(c.PostBody())
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v3/"+conf.PROFILE_KEY+"/direct/template/update/"+code, buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

// 광고성메시지(다이렉트) 템플릿 삭제
func Direct_template_delete_(c *fasthttp.RequestCtx) {
	conf := config.Conf
	
	code := c.UserValue("code").(string)

	buff := bytes.NewBuffer([]byte(`{}`))
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/direct/template/delete/"+code, buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

// 발신채널 전환
func Direct_convert_(c *fasthttp.RequestCtx) {
	conf := config.Conf

	buff := bytes.NewBuffer(c.PostBody())
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/sender/direct/convert", buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

// 발신채널 전환 상태 확인
func Direct_convert_result(c *fasthttp.RequestCtx) {
	conf := config.Conf

	args := fasthttp.AcquireArgs()
	defer fasthttp.ReleaseArgs(args)

	args.Add("senderKey", string(c.QueryArgs().Peek("senderKey")))

	queryString := args.String()

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/sender/direct/convert/result?"+queryString, nil)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

// 발신채널에 연결된 비즈월렛 변경
func Direct_bizWallet_change_(c *fasthttp.RequestCtx) {
	conf := config.Conf

	buff := bytes.NewBuffer(c.PostBody())
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/sender/direct/bizWallet/change", buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func TestFunc(c *fasthttp.RequestCtx) {
	res, _ := json.Marshal(map[string]string{
		"code": "success",
		"message": "test okay",
	})

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(res)
}

func imageMapping(c *fasthttp.RequestCtx, param map[string]io.Reader, prefix string, seq int) (map[string]io.Reader, error) {
	var rtnErr error = nil
	form, err := c.MultipartForm()
	if err != nil {
		rtnErr = err
	} else {
		if seq > 0 {
			for a := 1;a <= seq; a++ {
				files := form.File[prefix+strconv.Itoa(a)]
				if len(files) != 0 {
					extension := filepath.Ext(files[0].Filename)
					nfn := uuid.New().String() + extension
					err := saveUploadedFile(files[0], config.BasePath+"upload/" + nfn)
					if err != nil {
						config.Stdlog.Println("File 저장 오류 err : ", err)
					}
					param[prefix+strconv.Itoa(a)] = mustOpen(config.BasePath + "upload/" + nfn)
				}
			}
		} else {
			files := form.File[prefix]
			if len(files) != 0 {
				extension := filepath.Ext(files[0].Filename)
				nfn := uuid.New().String() + extension
				err := saveUploadedFile(files[0], config.BasePath+"upload/" + nfn)
				if err != nil {
					config.Stdlog.Println("File 저장 오류 err : ", err)
				}
				param[prefix] = mustOpen(config.BasePath + "upload/" + nfn)
			}
		}
	}
	
	return param, rtnErr
}

func upload(url string, values map[string]io.Reader) (*http.Response, error) {

	var buff bytes.Buffer
	w := multipart.NewWriter(&buff)

	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}

		if x, ok := r.(*os.File); ok {
			fw, _ = w.CreateFormFile(key, x.Name())
		} else {

			fw, _ = w.CreateFormField(key)
		}
		_, err := io.Copy(fw, r)

		if err != nil {
			return nil, err
		}

	}

	w.Close()

	req, err := http.NewRequest("POST", url, &buff)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := centerClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		//pwd, _ := os.Getwd()
		//fmt.Println("PWD: ", pwd)
		return nil
	}
	return r
}

func saveUploadedFile(fileHeader *multipart.FileHeader, dst string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := out.ReadFrom(src); err != nil {
		return err
	}

	return nil
}


////////////////////////////////////////////////////NPS AREA////////////////////////////////////////////////////

func checkAuthSiteId(auth, siteId string) string {
	var pk sql.NullString
	sql := `
	select
		profile_key
	from
		DHN_AUTH
	where
		auth_key = ? and site_id = ?
		`
		
	db.DB.QueryRow(sql, auth, siteId).Scan(&pk)

	if pk.Valid{
		return pk.String
	} else {
		return ""
	}
}

// 템플릿 등록 API
func CreateTemplateNps(c *fasthttp.RequestCtx) {
	// checkAuthSiteId(c)
	config.Stdlog.Println("CreateTemplateNps 시작")

	conf := config.Conf

	var kakakoReqParam kj.CtKakaoReq
	var kakakoResParam kj.StKakaoRes
	var kakakoResParam2 kj.StKakaoRes

	rawData := string(c.PostBody())

	values, err := url.ParseQuery(rawData)
	if err != nil {
		fmt.Println("Error parsing query:", err)
		return
	}

	result := make(map[string]interface{})

	for key, val := range values {
		if len(val) == 1 {
			// JSON 형태인지 확인 후 변환
			var jsonData interface{}
			if json.Unmarshal([]byte(val[0]), &jsonData) == nil {
				result[key] = jsonData
			} else {
				result[key] = val[0]
			}
		} else {
			result[key] = val
		}
	}

	senderKey := result["senderKey"]
	senderKeyType := result["senderKeyType"]
	templateCode := result["templateCode"]
	templateName := result["templateName"]
	templateContent := result["templateContent"]
	buttons := result["buttons"]

	if senderKey != nil {
		strSenderKey := senderKey.(string)
		kakakoReqParam.SenderKey = &strSenderKey
	} else {
		pk := checkAuthSiteId(string(c.Request.Header.Peek("auth_key")), string(c.Request.Header.Peek("siteid")))
		if len(pk) <= 0 {
			pk = "4e0114103341bd98f65ae0f5fe2acd6df7a6ffe3"
		}
		kakakoReqParam.SenderKey = &pk
	}

	if senderKeyType != nil {
		strSenderKeyType := senderKeyType.(string)
		kakakoReqParam.SenderKeyType = &strSenderKeyType
	} else {
		strSenderKeyType := "S"
		kakakoReqParam.SenderKeyType = &strSenderKeyType
	}

	if templateCode != nil {
		strTemplateCode := templateCode.(string)
		kakakoReqParam.TemplateCode = &strTemplateCode
	}

	if templateName != nil {
		strTemplateName := templateName.(string)
		kakakoReqParam.TemplateName = &strTemplateName
	}

	if templateContent != nil {
		strTemplateContent := templateContent.(string)
		kakakoReqParam.TemplateContent = &strTemplateContent
	}
	
	kakakoReqParam.TemplateMessageType = "BA"
	kakakoReqParam.TemplateEmphasizeType = "NONE"

	if buttons != nil {

		var tempBts []kj.KakaoButtonsNps

	    btnList, ok := buttons.([]interface{})
	    if !ok {
	        fmt.Println("buttons is not an array")
	        return
	    }

	    for _, item := range btnList {
	        v, ok := item.(map[string]interface{})
	        if !ok {
	            fmt.Println("Invalid button format")
	            continue
	        }

	        var tempBt kj.KakaoButtonsNps

	        if name, exists := v["name"].(string); exists {
	            tempBt.Name = &name
	        }
	        if linkType, exists := v["type"].(string); exists {
	            tempBt.LinkType = &linkType
	        }
	        if ordering, exists := v["ordering"].(float64); exists { // JSON 숫자는 float64
	            intValue := int(ordering)
	            tempBt.Ordering = &intValue
	        }
	        if linkAnd, exists := v["schemeAndroid"].(string); exists {
	            tempBt.LinkAnd = &linkAnd
	        }
	        if linkIos, exists := v["schemeIos"].(string); exists {
	            tempBt.LinkIos = &linkIos
	        }
	        if linkMo, exists := v["url_mobile"].(string); exists {
	            tempBt.LinkMo = &linkMo
	        }
	        if linkPc, exists := v["url_pc"].(string); exists {
	            tempBt.LinkPc = &linkPc
	        }
	        if pluginId, exists := v["pluginId"].(string); exists {
	            tempBt.PluginId = &pluginId
	        }

	        tempBts = append(tempBts, tempBt)
	    }

		kakakoReqParam.Buttons = &tempBts
	}

	jsonData, _ := json.Marshal(kakakoReqParam)
	buff := bytes.NewBuffer(jsonData)

	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/create", buff)
	if err != nil {
		config.Stdlog.Println(err)
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		config.Stdlog.Println(err)
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(bytes, &kakakoResParam)
	if strings.EqualFold(*kakakoResParam.Code, "200") {
		var reqData kj.KsReqNps
		var reqRes kj.KtrResNps
		reqData.SenderKey = kakakoResParam.Data.SenderKey
		reqData.SenderKeyType = kakakoResParam.Data.SenderKeyType
		reqData.TemplateCode = kakakoResParam.Data.TemplateCode
		reqRes = templateRequestNps(reqData)

		if strings.EqualFold(reqRes.Code, "200") {
			kakakoResParam2 = templateNps(reqData)
		} else {
			res, _ := json.Marshal(reqRes)
			c.SetContentType("application/json")
			c.SetStatusCode(fasthttp.StatusOK)
			c.SetBody(res)
			return
		}

		kakakoResParam2.Data.InspectionStatus = inspectionMapper(kakakoResParam2.Data.InspectionStatus)

		res, _ := json.Marshal(kakakoResParam2)

		c.SetContentType("application/json")
		c.SetStatusCode(fasthttp.StatusOK)
		c.SetBody(res)
		return
	} else {
		res, _ := json.Marshal(kakakoResParam)

		c.SetContentType("application/json")
		c.SetStatusCode(fasthttp.StatusOK)
		c.SetBody(res)
		return
	}

}

// 템플릿 조회 API
func SearchTemplateNps(c *fasthttp.RequestCtx) {
	// checkAuthSiteId(c)
	config.Stdlog.Println("SearchTemplateNps 시작")

	senderKey := string(c.QueryArgs().Peek("senderKey"))
	templateCode := string(c.QueryArgs().Peek("templateCode"))
	senderKeyType := string(c.QueryArgs().Peek("senderKeyType"))

	var reqData kj.KsReqNps
	var kakakoResParam kj.StKakaoRes

	reqData.SenderKey = &senderKey
	reqData.SenderKeyType = &senderKeyType
	reqData.TemplateCode = &templateCode

	kakakoResParam = templateNps(reqData)

	kakakoResParam.Data.InspectionStatus = inspectionMapper(kakakoResParam.Data.InspectionStatus)

	res, _ := json.Marshal(kakakoResParam)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(res)
	return
}

// 템플릿 수정 API
func UpdateTemplateNps(c *fasthttp.RequestCtx) {
	// checkAuthSiteId(c)
	config.Stdlog.Println("UpdateTemplateNps 시작")
	conf := config.Conf

	// var npsReqParam kj.UtReq17
	// var npsResParam kj.UtRes17
	var kakakoReqParam kj.UtKakaoReq
	var kakakoResParam kj.StKakaoRes
	var kakakoResParam2 kj.StKakaoRes

	rawData := string(c.PostBody())

	values, err := url.ParseQuery(rawData)
	if err != nil {
		fmt.Println("Error parsing query:", err)
		return
	}

	result := make(map[string]interface{})

	for key, val := range values {
		if len(val) == 1 {
			// JSON 형태인지 확인 후 변환
			var jsonData interface{}
			if json.Unmarshal([]byte(val[0]), &jsonData) == nil {
				result[key] = jsonData
			} else {
				result[key] = val[0]
			}
		} else {
			result[key] = val
		}
	}

	senderKey := result["senderKey"]
	templateCode := result["templateCode"]

	newSenderKey := result["newSenderKey"]
	newSenderKeyType := result["newSenderKeyType"]
	newTemplateCode := result["newTemplateCode"]
	newTemplateName := result["newTemplateName"]
	newTemplateContent := result["newTemplateContent"]

	buttons := result["buttons"]

	if senderKey != nil {
		strSenderKey := senderKey.(string)
		kakakoReqParam.SenderKey = &strSenderKey
	} else {
		pk := checkAuthSiteId(string(c.Request.Header.Peek("auth_key")), string(c.Request.Header.Peek("siteid")))
		if len(pk) <= 0 {
			pk = "4e0114103341bd98f65ae0f5fe2acd6df7a6ffe3"
		}
		kakakoReqParam.SenderKey = &pk
	}

	if templateCode != nil {
		strTemplateCode := templateCode.(string)
		kakakoReqParam.TemplateCode = &strTemplateCode
	}

	if newSenderKey != nil {
		strNewSenderKey := newSenderKey.(string)
		kakakoReqParam.NewSenderKey = &strNewSenderKey
	} else {
		pk := checkAuthSiteId(string(c.Request.Header.Peek("auth_key")), string(c.Request.Header.Peek("siteid")))
		if len(pk) <= 0 {
			pk = "4e0114103341bd98f65ae0f5fe2acd6df7a6ffe3"
		}
		kakakoReqParam.NewSenderKey = &pk
	}

	if newSenderKeyType != nil {
		strNewSenderKeyType := newSenderKeyType.(string)
		kakakoReqParam.NewSenderKeyType = &strNewSenderKeyType
	} else {
		strNewSenderKeyType := "S"
		kakakoReqParam.NewSenderKeyType = &strNewSenderKeyType
	}

	if newTemplateCode != nil {
		strNewTemplateCode := newTemplateCode.(string)
		kakakoReqParam.NewTemplateCode = &strNewTemplateCode
	}

	if newTemplateName != nil {
		strNewTemplateName := newTemplateName.(string)
		kakakoReqParam.NewTemplateName = &strNewTemplateName
	}

	if newTemplateContent != nil {
		strNewTemplateContent := newTemplateContent.(string)
		kakakoReqParam.NewTemplateContent = &strNewTemplateContent
	}

	kakakoReqParam.NewTemplateMessageType = "BA"
	kakakoReqParam.NewTemplateEmphasizeType = "NONE"

	if buttons != nil {
		var tempBts []kj.KakaoButtonsNps
		btnList, ok := buttons.([]interface{})
	    if !ok {
	        fmt.Println("buttons is not an array")
	        return
	    }

	    for _, item := range btnList {
	        v, ok := item.(map[string]interface{})
	        if !ok {
	            fmt.Println("Invalid button format")
	            continue
	        }

	        var tempBt kj.KakaoButtonsNps

	        if name, exists := v["name"].(string); exists {
	            tempBt.Name = &name
	        }
	        if linkType, exists := v["type"].(string); exists {
	            tempBt.LinkType = &linkType
	        }
	        if ordering, exists := v["ordering"].(float64); exists { // JSON 숫자는 float64
	            intValue := int(ordering)
	            tempBt.Ordering = &intValue
	        }
	        if linkAnd, exists := v["schemeAndroid"].(string); exists {
	            tempBt.LinkAnd = &linkAnd
	        }
	        if linkIos, exists := v["schemeIos"].(string); exists {
	            tempBt.LinkIos = &linkIos
	        }
	        if linkMo, exists := v["url_mobile"].(string); exists {
	            tempBt.LinkMo = &linkMo
	        }
	        if linkPc, exists := v["url_pc"].(string); exists {
	            tempBt.LinkPc = &linkPc
	        }
	        if pluginId, exists := v["pluginId"].(string); exists {
	            tempBt.PluginId = &pluginId
	        }

	        tempBts = append(tempBts, tempBt)
	    }

		kakakoReqParam.Buttons = &tempBts
	}

	jsonData, _ := json.Marshal(kakakoReqParam)
	buff := bytes.NewBuffer(jsonData)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/update", buff)
	if err != nil {
		config.Stdlog.Println(err)
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		config.Stdlog.Println(err)
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(bytes, &kakakoResParam)

	if kakakoResParam.Message != nil {
		config.Stdlog.Println(*kakakoResParam.Message)
	}

	if strings.EqualFold(*kakakoResParam.Code, "200") {
		var reqData kj.KsReqNps
		var reqRes kj.KtrResNps
		reqData.SenderKey = kakakoResParam.Data.SenderKey
		reqData.SenderKeyType = kakakoResParam.Data.SenderKeyType
		reqData.TemplateCode = kakakoResParam.Data.TemplateCode
		reqRes = templateRequestNps(reqData)

		if strings.EqualFold(reqRes.Code, "200") {
			kakakoResParam2 = templateNps(reqData)
		} else {
			res, _ := json.Marshal(reqRes)
			c.SetContentType("application/json")
			c.SetStatusCode(fasthttp.StatusOK)
			c.SetBody(res)
			return
		}

		kakakoResParam2.Data.InspectionStatus = inspectionMapper(kakakoResParam2.Data.InspectionStatus)

		res, _ := json.Marshal(kakakoResParam2)

		c.SetContentType("application/json")
		c.SetStatusCode(fasthttp.StatusOK)
		c.SetBody(res)
		return
	} else {
		res, _ := json.Marshal(kakakoResParam)

		c.SetContentType("application/json")
		c.SetStatusCode(fasthttp.StatusOK)
		c.SetBody(res)
		return
	}

}

// 템플릿 삭제 API
func DeleteTemplateNps(c *fasthttp.RequestCtx) {
	// checkAuthSiteId(c)

	conf := config.Conf

	var reqData kj.KsReqNps

	rawData := string(c.PostBody())

	values, err := url.ParseQuery(rawData)
	if err != nil {
		fmt.Println("Error parsing query:", err)
		return
	}

	result := make(map[string]interface{})

	for key, val := range values {
		if len(val) == 1 {
			// JSON 형태인지 확인 후 변환
			var jsonData interface{}
			if json.Unmarshal([]byte(val[0]), &jsonData) == nil {
				result[key] = jsonData
			} else {
				result[key] = val[0]
			}
		} else {
			result[key] = val
		}
	}

	senderKey := result["senderKey"]
	senderKeyType := result["senderKeyType"]
	templateCode := result["templateCode"]

	if senderKey != nil {
		strSenderKey := senderKey.(string)
		reqData.SenderKey = &strSenderKey
	} else {
		pk := checkAuthSiteId(string(c.Request.Header.Peek("auth_key")), string(c.Request.Header.Peek("siteid")))
		reqData.SenderKey = &pk
	}

	if senderKeyType != nil {
		strSenderKeyType := senderKeyType.(string)
		reqData.SenderKeyType = &strSenderKeyType
	} else {
		strSenderKeyType := "S"
		reqData.SenderKeyType = &strSenderKeyType
	}

	if templateCode != nil {
		strTemplateCode := templateCode.(string)
		reqData.TemplateCode = &strTemplateCode
	}

	jsonData, _ := json.Marshal(reqData)
	buff := bytes.NewBuffer(jsonData)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/delete", buff)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)

}

// 템플릿 코멘트 등록 API
func SetComment(c *fasthttp.RequestCtx) {
	// checkAuthSiteId(c)
	config.Stdlog.Println("SetComment 시작")

	conf := config.Conf

	var reqData kj.KsReqNps
	var npsResParam kj.UtRes17
	var kakakoResParam kj.StKakaoRes

	rawData := string(c.PostBody())

	values, err := url.ParseQuery(rawData)
	if err != nil {
		fmt.Println("Error parsing query:", err)
		return
	}

	result := make(map[string]interface{})

	for key, val := range values {
		if len(val) == 1 {
			// JSON 형태인지 확인 후 변환
			var jsonData interface{}
			if json.Unmarshal([]byte(val[0]), &jsonData) == nil {
				result[key] = jsonData
			} else {
				result[key] = val[0]
			}
		} else {
			result[key] = val
		}
	}

	senderKey := result["senderKey"]
	senderKeyType := result["senderKeyType"]
	templateCode := result["templateCode"]
	comment := result["comment"]

	if senderKey != nil {
		strSenderKey := senderKey.(string)
		reqData.SenderKey = &strSenderKey
	} else {
		pk := checkAuthSiteId(string(c.Request.Header.Peek("auth_key")), string(c.Request.Header.Peek("siteid")))
		if len(pk) <= 0 {
			pk = "4e0114103341bd98f65ae0f5fe2acd6df7a6ffe3"
		}
		reqData.SenderKey = &pk
	}

	if senderKeyType != nil {
		strSenderKeyType := senderKeyType.(string)
		reqData.SenderKeyType = &strSenderKeyType
	} else {
		strSenderKeyType := "S"
		reqData.SenderKeyType = &strSenderKeyType
	}

	if templateCode != nil {
		strTemplateCode := templateCode.(string)
		reqData.TemplateCode = &strTemplateCode
	}

	if comment != nil {
		var commentStr string
		if comment != nil {
			switch v := comment.(type) {
			case string:
				commentStr = v
			case []string:
				commentStr = strings.Join(v, "\n") // 여러 줄 처리
			default:
				fmt.Println("Unexpected type for comment")
			}
		}
		reqData.Comment = &commentStr
	}

	kakakoResParam = templateNps(reqData)

	if strings.EqualFold(*kakakoResParam.Code, "200") {
		if strings.EqualFold(kakakoResParam.Data.InspectionStatus, "REG") {
			var reqRes kj.KtrResNps
			reqRes = templateRequestNps(reqData)

			res, _ := json.Marshal(reqRes)

			if strings.EqualFold(reqRes.Code, "200") {
				reqRes = cancelTemplateRequestNps(reqData)
			}

			c.SetContentType("application/json")
			c.SetStatusCode(fasthttp.StatusOK)
			c.SetBody(res)
			return
		} else if strings.EqualFold(kakakoResParam.Data.InspectionStatus, "REQ") {
			var reqRes kj.KtrResNps
			reqRes = cancelTemplateRequestNps(reqData)

			res, _ := json.Marshal(reqRes)

			if strings.EqualFold(reqRes.Code, "200") {
				reqRes = templateRequestNps(reqData)
			}

			c.SetContentType("application/json")
			c.SetStatusCode(fasthttp.StatusOK)
			c.SetBody(res)
			return
		} else if strings.EqualFold(kakakoResParam.Data.InspectionStatus, "REJ") {
			var kakakoReqParam kj.UtKakaoReq

			kakakoReqParam.SenderKey = kakakoResParam.Data.SenderKey
			if kakakoResParam.Data.SenderKeyType != nil {
				kakakoReqParam.SenderKeyType = kakakoResParam.Data.SenderKeyType	
			}
			kakakoReqParam.TemplateCode = kakakoResParam.Data.TemplateCode

			kakakoReqParam.NewSenderKey = kakakoResParam.Data.SenderKey
			if kakakoResParam.Data.SenderKey != nil {
				kakakoReqParam.NewSenderKeyType = kakakoResParam.Data.SenderKeyType	
			}
			kakakoReqParam.NewTemplateCode = kakakoResParam.Data.TemplateCode
			kakakoReqParam.NewTemplateName = kakakoResParam.Data.TemplateName
			kakakoReqParam.NewTemplateContent = kakakoResParam.Data.TemplateContent
			kakakoReqParam.NewTemplateMessageType = "BA"
			kakakoReqParam.NewTemplateEmphasizeType = "NONE"
			kakakoReqParam.Buttons = kakakoResParam.Data.Buttons

			jsonData, _ := json.Marshal(kakakoReqParam)
			buff := bytes.NewBuffer(jsonData)
			req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/update", buff)
			if err != nil {
				config.Stdlog.Println(err)
				c.Error(err.Error(), fasthttp.StatusBadRequest)
				return
			}

			req.Header.Add("Content-Type", "application/json")
			resp, err := centerClient.Do(req)
			if err != nil {
				config.Stdlog.Println(err)
				c.Error(err.Error(), fasthttp.StatusBadRequest)
				return
			}
			defer resp.Body.Close()

			bytes, _ := ioutil.ReadAll(resp.Body)

			json.Unmarshal(bytes, &kakakoResParam)

			if strings.EqualFold(*kakakoResParam.Code, "200") {
				var reqRes kj.KtrResNps
				reqRes = templateRequestNps(reqData)

				res, _ := json.Marshal(reqRes)

				if strings.EqualFold(reqRes.Code, "200") {
					reqRes = cancelTemplateRequestNps(reqData)
				}

				c.SetContentType("application/json")
				c.SetStatusCode(fasthttp.StatusOK)
				c.SetBody(res)
				return
			} else {
				npsResParam.Code = *kakakoResParam.Code
				res, _ := json.Marshal(npsResParam)
				c.SetContentType("application/json")
				c.SetStatusCode(fasthttp.StatusBadRequest)
				c.SetBody(res)
				return
			}
			
		}
	} else {
		npsResParam.Code = *kakakoResParam.Code
		res, _ := json.Marshal(npsResParam)
		c.SetContentType("application/json")
		c.SetStatusCode(fasthttp.StatusBadRequest)
		c.SetBody(res)
		return
	}

}

// 템플릿 검수 함수
func templateRequestNps(data kj.KsReqNps) kj.KtrResNps {
	var reqRes kj.KtrResNps

	jsonData, _ := json.Marshal(data)
	buff := bytes.NewBuffer(jsonData)
	req, _ := http.NewRequest("POST", config.Conf.CENTER_SERVER+"api/v2/"+config.Conf.PROFILE_KEY+"/alimtalk/template/request", buff)

	req.Header.Add("Content-Type", "application/json")
	resp, _ := centerClient.Do(req)
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	json.Unmarshal(bytes, &reqRes)

	if reqRes.Message != nil {
		config.Stdlog.Println(*reqRes.Message)
	}

	return reqRes
}

// 템플릿 조회 함수
func templateNps(data kj.KsReqNps) kj.StKakaoRes {
	var reqRes kj.StKakaoRes 

	req, _ := http.NewRequest("GET", config.Conf.CENTER_SERVER+"api/v2/"+config.Conf.PROFILE_KEY+"/alimtalk/template?senderKey="+*data.SenderKey+"&templateCode="+url.QueryEscape(*data.TemplateCode)+"&senderKeyType="+*data.SenderKeyType, nil)
	
	req.Header.Add("Accept-Charset", "utf-8")
	
	resp, _ := centerClient.Do(req)

	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	json.Unmarshal(bytes, &reqRes)

	if reqRes.Message != nil {
		config.Stdlog.Println(*reqRes.Message)
	}

	return reqRes
}

// 템플릿 검수 취소 함수
func cancelTemplateRequestNps(data kj.KsReqNps) kj.KtrResNps {
	var reqRes kj.KtrResNps

	jsonData, _ := json.Marshal(data)
	buff := bytes.NewBuffer(jsonData)
	req, _ := http.NewRequest("POST", config.Conf.CENTER_SERVER+"api/v2/"+config.Conf.PROFILE_KEY+"/alimtalk/template/cancel_request", buff)

	req.Header.Add("Content-Type", "application/json")
	resp, _ := centerClient.Do(req)
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	json.Unmarshal(bytes, &reqRes)

	if reqRes.Message != nil {
		config.Stdlog.Println(*reqRes.Message)
	}

	return reqRes
}

func inspectionMapper(stat string) string {
	ret := ""
	if strings.EqualFold(stat, "REG") {
		ret = "REG"
	} else if strings.EqualFold(stat, "REQ"){
		ret = "INS"
	} else if strings.EqualFold(stat, "APR"){
		ret = "COM"
	} else if strings.EqualFold(stat, "REJ"){
		ret = "REJ"
	}
	return ret
}

////////////////////////////////////////////////////NPS AREA////////////////////////////////////////////////////







