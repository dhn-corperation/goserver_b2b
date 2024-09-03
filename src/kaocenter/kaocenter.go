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
	// "strconv"
	"strings"
	"time"

	config "mycs/src/kaoconfig"
	db "mycs/src/kaodatabasepool"

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

	//client := &http.Client{}
	resp, err := centerClient.Do(req)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Sender_Create(c *fasthttp.RequestCtx) {
	conf := config.Conf

	token := string(c.Request.Header.Peek("token"))
	phoneNumber := string(c.Request.Header.Peek("phoneNumber"))

	jsonstr, err := json.Marshal(c.PostBody())
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	buff := bytes.NewBuffer(jsonstr)
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

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Sender_Delete(c *fasthttp.RequestCtx) {
	conf := config.Conf

	jsonstr, err := json.Marshal(c.PostBody())
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	buff := bytes.NewBuffer(jsonstr)
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

	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Sender_Recover(c *fasthttp.RequestCtx) {
	conf := config.Conf

	jsonstr, err := json.Marshal(c.PostBody())
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	buff := bytes.NewBuffer(jsonstr)
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
	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Template_Create(c *fasthttp.RequestCtx) {
	conf := config.Conf

	jsonstr, err := json.Marshal(c.PostBody())
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	buff := bytes.NewBuffer(jsonstr)
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
	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Template_Request(c *fasthttp.RequestCtx) {
	conf := config.Conf

	jsonstr, err := json.Marshal(c.PostBody())
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	buff := bytes.NewBuffer(jsonstr)
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
	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Template_Cancel_Request(c *fasthttp.RequestCtx) {
	conf := config.Conf

	jsonstr, err := json.Marshal(c.PostBody())
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	buff := bytes.NewBuffer(jsonstr)
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
	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Template_Update(c *fasthttp.RequestCtx) {
	conf := config.Conf
	
	jsonstr, err := json.Marshal(c.PostBody())
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	buff := bytes.NewBuffer(jsonstr)
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
	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Template_Stop(c *fasthttp.RequestCtx) {
	conf := config.Conf

	jsonstr, err := json.Marshal(c.PostBody())
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	buff := bytes.NewBuffer(jsonstr)
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
	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Template_Reuse(c *fasthttp.RequestCtx) {
	conf := config.Conf

	jsonstr, err := json.Marshal(c.PostBody())
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	buff := bytes.NewBuffer(jsonstr)
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
	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Template_Delete(c *fasthttp.RequestCtx) {
	conf := config.Conf

	jsonstr, err := json.Marshal(c.PostBody())
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	buff := bytes.NewBuffer(jsonstr)
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
	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Template_Dormant_Release(c *fasthttp.RequestCtx) {
	conf := config.Conf

	jsonstr, err := json.Marshal(c.PostBody())
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	buff := bytes.NewBuffer(jsonstr)
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
	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Group_Sender_Add(c *fasthttp.RequestCtx) {
	conf := config.Conf

	jsonstr, err := json.Marshal(c.PostBody())
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	buff := bytes.NewBuffer(jsonstr)
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
	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Group_Sender_Remove(c *fasthttp.RequestCtx) {
	conf := config.Conf

	jsonstr, err := json.Marshal(c.PostBody())
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	buff := bytes.NewBuffer(jsonstr)
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
	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Channel_Create_(c *fasthttp.RequestCtx) {
	conf := config.Conf

	jsonstr, err := json.Marshal(c.PostBody())
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	buff := bytes.NewBuffer(jsonstr)
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
	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Channel_Update_(c *fasthttp.RequestCtx) {
	conf := config.Conf

	jsonstr, err := json.Marshal(c.PostBody())
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	buff := bytes.NewBuffer(jsonstr)
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
	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Channel_Senders_(c *fasthttp.RequestCtx) {
	conf := config.Conf

	jsonstr, err := json.Marshal(c.PostBody())
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	buff := bytes.NewBuffer(jsonstr)
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
	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Channel_Delete_(c *fasthttp.RequestCtx) {
	conf := config.Conf

	jsonstr, err := json.Marshal(c.PostBody())
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	buff := bytes.NewBuffer(jsonstr)
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
	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Plugin_callbackUrl_Create(c *fasthttp.RequestCtx) {
	conf := config.Conf

	jsonstr, err := json.Marshal(c.PostBody())
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	buff := bytes.NewBuffer(jsonstr)
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
	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Plugin_callbackUrl_Update(c *fasthttp.RequestCtx) {
	conf := config.Conf

	jsonstr, err := json.Marshal(c.PostBody())
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	buff := bytes.NewBuffer(jsonstr)
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
	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func Plugin_callbackUrl_Delete(c *fasthttp.RequestCtx) {
	conf := config.Conf

	jsonstr, err := json.Marshal(c.PostBody())
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	buff := bytes.NewBuffer(jsonstr)
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
	bytes, _ := ioutil.ReadAll(resp.Body)

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	c.SetBody(bytes)
}

func FT_Upload(c *fasthttp.RequestCtx) {
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

	err = saveUploadedFile(file, config.BasePath+"upload/"+newFileName)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	param := map[string]io.Reader{
		"image": mustOpen(config.BasePath+"upload/" +newFileName),
	}

	resp, err := upload(conf.IMAGE_SERVER+"v1/"+conf.PROFILE_KEY+"/image/friendtalk", param)
	if err != nil {
		c.Error(err.Error(), fasthttp.StatusBadRequest)
	}
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
	}
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
	
	seq := 1
	var startNow = time.Now()
	var group_no = fmt.Sprintf("%04d%02d%02d%02d%02d%02d%09d", startNow.Year(), startNow.Month(), startNow.Day(), startNow.Hour(), startNow.Minute(), startNow.Second(), startNow.Nanosecond())
	
	for _, key := range imageKeys {
		files := form.File[key]
		if len(files) != 0 {
			extension := filepath.Ext(files[0].Filename)
			switch seq {
			case 1:
				newFileName1 = config.BasePath+"upload/mms/" + uuid.New().String() + extension
			case 2:
				newFileName2 = config.BasePath+"upload/mms/" + uuid.New().String() + extension
			case 3:
				newFileName3 = config.BasePath+"upload/mms/" + uuid.New().String() + extension
			}
			err := saveUploadedFile(files[0], config.BasePath+"upload/mms/" + uuid.New().String() + extension)
			if err != nil {
				config.Stdlog.Println("File ", seq," 저장 오류 : ", newFileName1, err)
				newFileName1 = ""
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

		res, _ := json.Marshal(map[string]string{
			"image_group":group_no,
		})
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


// func Image_wideItemList(c *fasthttp.RequestCtx) {
// 	conf := config.Conf
// 	config.Stdlog.Println("Call ")
	
// 	var newFileName1,newFileName2,newFileName3,newFileName4 string
	
// 	file1, err1 := c.FormFile("image_1")
// 	if err1 != nil {
// 		config.Stdlog.Println(err1.Error())
// 		c.String(http.StatusBadRequest, fmt.Sprintf("File 1 - get form err: %s", err1.Error()))
// 		return
// 	}

// 	extension := filepath.Ext(file1.Filename)
// 	newFileName1 = uuid.New().String() + extension

// 	err1 = c.SaveUploadedFile(file1, config.BasePath+"upload/" + newFileName1)
// 	if err1 != nil {
// 		c.String(http.StatusBadRequest, fmt.Sprintf("File 1 - get form err: %s", err1.Error()))
// 		return
// 	}

// 	file2, err2 := c.FormFile("image_2")
// 	if err2 == nil {
// 		extension := filepath.Ext(file2.Filename)
// 		newFileName2 = uuid.New().String() + extension
			
// 		err2 = c.SaveUploadedFile(file2, config.BasePath+"upload/" + newFileName2)
// 		if err2 != nil {
// 			newFileName2 = "_"
// 		}
// 	} else {
// 		newFileName2 = "_"
// 	}

// 	file3, err3 := c.FormFile("image_3")
// 	if err3 == nil {
// 		extension := filepath.Ext(file3.Filename)
// 		newFileName3 = uuid.New().String() + extension
	
// 		err3 = c.SaveUploadedFile(file3, config.BasePath+"upload/" + newFileName3)
// 		if err3 != nil {
// 			newFileName3 = "_"
// 		}
// 	} else {
// 		newFileName3 = "_"
// 	}

// 	file4, err4 := c.FormFile("image_4")
// 	if err4 == nil {
// 		extension := filepath.Ext(file4.Filename)
// 		newFileName4 = uuid.New().String() + extension
	
// 		err4 = c.SaveUploadedFile(file4, config.BasePath+"upload/" + newFileName4)
// 		if err4 != nil {
// 			newFileName4 = "_"
// 		}
// 	} else {
// 		newFileName4 = "_"
// 	}
	
// 	param := map[string]io.Reader{
// 		"image_1": mustOpen(config.BasePath+"upload/" + newFileName1),
// 		"image_2": mustOpen(config.BasePath+"upload/" + newFileName2),
// 		"image_3": mustOpen(config.BasePath+"upload/" + newFileName3),
// 		"image_4": mustOpen(config.BasePath+"upload/" + newFileName4),
// 	}

// 	if newFileName4 == "_" {
// 		delete(param, "image_4")
// 	}

// 	if newFileName3 == "_" {
// 		delete(param, "image_3")
// 	}

// 	if newFileName2 == "_" {
// 		delete(param, "image_2")
// 	}

// 	resp, err1 := upload(conf.IMAGE_SERVER+ "v1/"+conf.PROFILE_KEY+"/image/friendtalk/wideItemList", param)
// 	bytes, _ := ioutil.ReadAll(resp.Body)
// 	c.SetContentType("application/json")
// 	c.SetStatusCode(fasthttp.StatusOK)
// 	c.SetBody(bytes)
// }

// func Image_carousel(c *fasthttp.RequestCtx) {
// 	conf := config.Conf
// 	config.Stdlog.Println("Call ")
	
// 	var newFileName1,newFileName2,newFileName3,newFileName4,newFileName5,newFileName6 string
	
// 	file1, err1 := c.FormFile("image_1")
// 	if err1 != nil {
// 		config.Stdlog.Println(err1.Error())
// 		c.String(http.StatusBadRequest, fmt.Sprintf("File 1 - get form err: %s", err1.Error()))
// 		return
// 	}

// 	extension := filepath.Ext(file1.Filename)
// 	newFileName1 = uuid.New().String() + extension

// 	err1 = c.SaveUploadedFile(file1, config.BasePath+"upload/" + newFileName1)
// 	if err1 != nil {
// 		c.String(http.StatusBadRequest, fmt.Sprintf("File 1 - get form err: %s", err1.Error()))
// 		return
// 	}

// 	file2, err2 := c.FormFile("image_2")
// 	if err2 == nil {
// 		extension := filepath.Ext(file2.Filename)
// 		newFileName2 = uuid.New().String() + extension
			
// 		err2 = c.SaveUploadedFile(file2, config.BasePath+"upload/" + newFileName2)
// 		if err2 != nil {
// 			newFileName2 = "_"
// 		}
// 	} else {
// 		newFileName2 = "_"
// 	}

// 	file3, err3 := c.FormFile("image_3")
// 	if err3 == nil {
// 		extension := filepath.Ext(file3.Filename)
// 		newFileName3 = uuid.New().String() + extension
	
// 		err3 = c.SaveUploadedFile(file3, config.BasePath+"upload/" + newFileName3)
// 		if err3 != nil {
// 			newFileName3 = "_"
// 		}
// 	} else {
// 		newFileName3 = "_"
// 	}

// 	file4, err4 := c.FormFile("image_4")
// 	if err4 == nil {
// 		extension := filepath.Ext(file4.Filename)
// 		newFileName4 = uuid.New().String() + extension
	
// 		err4 = c.SaveUploadedFile(file4, config.BasePath+"upload/" + newFileName4)
// 		if err4 != nil {
// 			newFileName4 = "_"
// 		}
// 	} else {
// 		newFileName4 = "_"
// 	}
	
// 	file5, err5 := c.FormFile("image_5")
// 	if err5 == nil {
// 		extension := filepath.Ext(file5.Filename)
// 		newFileName5 = uuid.New().String() + extension
	
// 		err5 = c.SaveUploadedFile(file5, config.BasePath+"upload/" + newFileName5)
// 		if err5 != nil {
// 			newFileName5 = "_"
// 		}
// 	} else {
// 		newFileName5 = "_"
// 	}
	
// 	file6, err6 := c.FormFile("image_6")
// 	if err6 == nil {
// 		extension := filepath.Ext(file6.Filename)
// 		newFileName6 = uuid.New().String() + extension
	
// 		err6 = c.SaveUploadedFile(file6, config.BasePath+"upload/" + newFileName6)
// 		if err6 != nil {
// 			newFileName6 = "_"
// 		}
// 	} else {
// 		newFileName6 = "_"
// 	}
		
// 	param := map[string]io.Reader{
// 		"image_1": mustOpen(config.BasePath+"upload/" + newFileName1),
// 		"image_2": mustOpen(config.BasePath+"upload/" + newFileName2),
// 		"image_3": mustOpen(config.BasePath+"upload/" + newFileName3),
// 		"image_4": mustOpen(config.BasePath+"upload/" + newFileName4),
// 		"image_5": mustOpen(config.BasePath+"upload/" + newFileName5),
// 		"image_6": mustOpen(config.BasePath+"upload/" + newFileName6),
// 	}
// 	if newFileName6 == "_" {
// 		delete(param, "image_6")
// 	}
	
// 	if newFileName5 == "_" {
// 		delete(param, "image_5")
// 	}
	
// 	if newFileName4 == "_" {
// 		delete(param, "image_4")
// 	}

// 	if newFileName3 == "_" {
// 		delete(param, "image_3")
// 	}

// 	if newFileName2 == "_" {
// 		delete(param, "image_2")
// 	}

// 	resp, err1 := upload(conf.IMAGE_SERVER+ "v1/"+conf.PROFILE_KEY+"/image/friendtalk/wideItemList", param)
// 	bytes, _ := ioutil.ReadAll(resp.Body)
// 	c.SetContentType("application/json")
// 	c.SetStatusCode(fasthttp.StatusOK)
// 	c.SetBody(bytes)
// }

// func Get_Polling_Id(c *fasthttp.RequestCtx) {
// 	conf := config.Conf
// 	respId := c.Param("respid")

// 	buff := bytes.NewBuffer([]byte("{}"))
// 	req, err := http.NewRequest("POST", conf.API_SERVER+"/v3/"+conf.PROFILE_KEY+"/response/"+respId, buff)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	req.Header.Add("Content-Type", "application/json")
// 	resp, err2 := centerClient.Do(req)
// 	if err2 != nil {
// 		c.JSON(http.StatusBadRequest, err2.Error())
// 		return
// 	}

// 	defer resp.Body.Close()

// 	bytes, _ := io.ReadAll(resp.Body)
// 	c.SetContentType("application/json")
// 	c.SetStatusCode(fasthttp.StatusOK)
// 	c.SetBody(bytes)

// }

// func AT_Highlight_Image(c *fasthttp.RequestCtx) {
// 	conf := config.Conf

// 	param, err := image_Seq_Mapping(c, map[string]io.Reader{}, 0, "image")
// 	if err != nil {
// 		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
// 		return
// 	}

// 	// file, err := c.FormFile("image")
// 	// if err != nil {
// 	// 	c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
// 	// 	return
// 	// }

// 	// extension := filepath.Ext(file.Filename)
// 	// newFileName := uuid.New().String() + extension

// 	// err = c.SaveUploadedFile(file, config.BasePath+"upload/" + newFileName)
// 	// if err != nil {
// 	// 	c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
// 	// 	return
// 	// }

// 	// param := map[string]io.Reader{
// 	// 	"image": mustOpen(config.BasePath+"upload/" + newFileName),
// 	// }

// 	resp, err := upload(conf.IMAGE_SERVER+"v1/"+conf.PROFILE_KEY+"/image/alimtalk/itemHighlight", param)
// 	if err != nil {
// 		config.Stdlog.Println("File upload 오류 : ", err)
// 	}

// 	bytes, _ := io.ReadAll(resp.Body)
// 	c.SetContentType("application/json")
// 	c.SetStatusCode(fasthttp.StatusOK)
// 	c.SetBody(bytes)
// }

// func FT_Carousel_Feed_Image(c *fasthttp.RequestCtx) {
// 	conf := config.Conf

// 	param, err := image_Seq_Mapping(c, map[string]io.Reader{}, 10, "image")
// 	if err != nil {
// 		config.Stdlog.Println("image Mapping 오류 : ", err)
// 	}

// 	resp, _ := upload(conf.IMAGE_SERVER+"v1/"+conf.PROFILE_KEY+"/image/friendtalk/carousel", param)
// 	bytes, _ := io.ReadAll(resp.Body)
// 	c.SetContentType("application/json")
// 	c.SetStatusCode(fasthttp.StatusOK)
// 	c.SetBody(bytes)
// }

// func FT_Carousel_Commerce_Image(c *fasthttp.RequestCtx) {
// 	conf := config.Conf

// 	param, err := image_Seq_Mapping(c, map[string]io.Reader{}, 11, "image")
// 	if err != nil {
// 		config.Stdlog.Println("image Mapping 오류 : ", err)
// 	}

// 	resp, _ := upload(conf.IMAGE_SERVER+"v1/"+conf.PROFILE_KEY+"/image/friendtalk/carouselCommerce", param)
// 	bytes, _ := io.ReadAll(resp.Body)
// 	c.SetContentType("application/json")
// 	c.SetStatusCode(fasthttp.StatusOK)
// 	c.SetBody(bytes)
// }

// func DM_Default_Image(c *fasthttp.RequestCtx) {
// 	conf := config.Conf

// 	param, err := image_Seq_Mapping(c, map[string]io.Reader{}, 0, "image")
// 	if err != nil {
// 		config.Stdlog.Println("image Mapping 오류 : ", err)
// 	}

// 	resp, _ := upload(conf.IMAGE_SERVER+"v2/"+conf.PROFILE_KEY+"/image/default", param)
// 	bytes, _ := io.ReadAll(resp.Body)
// 	c.SetContentType("application/json")
// 	c.SetStatusCode(fasthttp.StatusOK)
// 	c.SetBody(bytes)
// }

// func DM_Wide_Image(c *fasthttp.RequestCtx) {
// 	conf := config.Conf

// 	param, err := image_Seq_Mapping(c, map[string]io.Reader{}, 0, "image")
// 	if err != nil {
// 		config.Stdlog.Println("image Mapping 오류 : ", err)
// 	}

// 	resp, _ := upload(conf.IMAGE_SERVER+"v2/"+conf.PROFILE_KEY+"/image/wide", param)
// 	bytes, _ := io.ReadAll(resp.Body)
// 	c.SetContentType("application/json")
// 	c.SetStatusCode(fasthttp.StatusOK)
// 	c.SetBody(bytes)
// }

// func DM_Widelist_First_image(c *fasthttp.RequestCtx) {
// 	conf := config.Conf

// 	param, err := image_Seq_Mapping(c, map[string]io.Reader{}, 0, "image")
// 	if err != nil {
// 		config.Stdlog.Println("image Mapping 오류 : ", err)
// 	}

// 	resp, _ := upload(conf.IMAGE_SERVER+"v2/"+conf.PROFILE_KEY+"/image/wideItemList/first", param)
// 	bytes, _ := io.ReadAll(resp.Body)
// 	c.SetContentType("application/json")
// 	c.SetStatusCode(fasthttp.StatusOK)
// 	c.SetBody(bytes)
// }

// func DM_Widelist_Image(c *fasthttp.RequestCtx) {
// 	conf := config.Conf

// 	param, err := image_Seq_Mapping(c, map[string]io.Reader{}, 3, "image")
// 	if err != nil {
// 		config.Stdlog.Println("image Mapping 오류 : ", err)
// 	}

// 	resp, _ := upload(conf.IMAGE_SERVER+"v2/"+conf.PROFILE_KEY+"/image/wideItemList", param)
// 	bytes, _ := io.ReadAll(resp.Body)
// 	c.SetContentType("application/json")
// 	c.SetStatusCode(fasthttp.StatusOK)
// 	c.SetBody(bytes)
// }

// func DM_Carousel_Feed_Image(c *fasthttp.RequestCtx) {
// 	conf := config.Conf

// 	param, err := image_Seq_Mapping(c, map[string]io.Reader{}, 10, "image")
// 	if err != nil {
// 		config.Stdlog.Println("image Mapping 오류 : ", err)
// 	}

// 	resp, _ := upload(conf.IMAGE_SERVER+"v2/"+conf.PROFILE_KEY+"/image/carouselFeed", param)
// 	bytes, _ := io.ReadAll(resp.Body)
// 	c.SetContentType("application/json")
// 	c.SetStatusCode(fasthttp.StatusOK)
// 	c.SetBody(bytes)
// }

// func DM_Carousel_Commerce_Image(c *fasthttp.RequestCtx) {
// 	conf := config.Conf

// 	param, err := image_Seq_Mapping(c, map[string]io.Reader{}, 11, "image")
// 	if err != nil {
// 		config.Stdlog.Println("image Mapping 오류 : ", err)
// 	}

// 	resp, _ := upload(conf.IMAGE_SERVER+"v2/"+conf.PROFILE_KEY+"/image/carouselCommerce", param)
// 	bytes, _ := io.ReadAll(resp.Body)
// 	c.SetContentType("application/json")
// 	c.SetStatusCode(fasthttp.StatusOK)
// 	c.SetBody(bytes)
// }

// // 친구톡 API
// // 별첨1 - 비즈폼 업로드 요청
// func Bizform_upload_(c *fasthttp.RequestCtx) {
// 	conf := config.Conf

// 	param := &Bizform_upload{}
// 	err := c.Bind(param)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	jsonstr, _ := json.Marshal(param)
// 	buff := bytes.NewBuffer(jsonstr)
// 	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/bizform/upload", buff)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	req.Header.Add("Content-Type", "application/json")
// 	resp, err := centerClient.Do(req)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	bytes, _ := io.ReadAll(resp.Body)
// 	c.SetContentType("application/json")
// 	c.SetStatusCode(fasthttp.StatusOK)
// 	c.SetBody(bytes)

// }

// // 별첨2 - 친구톡 발송 가능 모수 확인
// func Ft_possible_(c *fasthttp.RequestCtx) {
// 	conf := config.Conf

// 	param := &Ft_possible{}
// 	err := c.Bind(param)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	jsonstr, _ := json.Marshal(param)
// 	buff := bytes.NewBuffer(jsonstr)
// 	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/friendtalk/possible", buff)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	req.Header.Add("Content-Type", "application/json")
// 	resp, err := centerClient.Do(req)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	bytes, _ := io.ReadAll(resp.Body)
// 	c.SetContentType("application/json")
// 	c.SetStatusCode(fasthttp.StatusOK)
// 	c.SetBody(bytes)

// }

// // 센터 API
// // 발신 프로필 조회2 (톡 채널 키로 조회)
// func Sender_channel(c *fasthttp.RequestCtx) {
// 	conf := config.Conf

// 	talkChannelKey := c.Param("talkChannelKey")

// 	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v3/"+conf.PROFILE_KEY+"/sender/"+talkChannelKey, nil)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}

// 	resp, err := centerClient.Do(req)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	bytes, _ := io.ReadAll(resp.Body)
// 	c.SetContentType("application/json")
// 	c.SetStatusCode(fasthttp.StatusOK)
// 	c.SetBody(bytes)
// }

// // 최근 변경 발신 프로필 조회
// func Sender_modified(c *fasthttp.RequestCtx) {
// 	conf := config.Conf

// 	//since := c.Query("since")
// 	//page := c.Query("page")
// 	//count := c.Query("count")

// 	params := map[string]string{
// 		"since": c.Query("since"),
// 		"page":  c.Query("page"),
// 		"count": c.Query("count"),
// 	}

// 	query := c.Request.URL.Query()
// 	for key, value := range params {
// 		if value != "" {
// 			query.Set(key, value)
// 		}
// 	}

// 	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v3/"+conf.PROFILE_KEY+"/sender/last_modified?"+query.Encode(), nil)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}

// 	resp, err := centerClient.Do(req)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	bytes, _ := io.ReadAll(resp.Body)
// 	c.SetContentType("application/json")
// 	c.SetStatusCode(fasthttp.StatusOK)
// 	c.SetBody(bytes)
// }

// // 검수요청 (파일첨부)
// func Template_request_with_file(c *fasthttp.RequestCtx) {
// 	conf := config.Conf

// 	param := map[string]io.Reader{
// 		"senderKey":     strings.NewReader(c.PostForm("senderKey")),
// 		"templateCode":  strings.NewReader(c.PostForm("templateCode")),
// 		"senderKeyType": strings.NewReader(c.PostForm("senderKeyType")),
// 		"comment":       strings.NewReader(c.PostForm("comment")),
// 	}

// 	file, err := c.FormFile("attachment")
// 	if err == nil { 
// 		extension := filepath.Ext(file.Filename)
// 		newFileName := uuid.New().String() + extension

// 		err = c.SaveUploadedFile(file, config.BasePath+"upload/"+newFileName)
// 		if err != nil {
// 			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
// 			return
// 		}
// 		param["attachment"] = mustOpen(config.BasePath + "upload/" + newFileName)
// 	}

// 	resp, err := upload(conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/request_with_file", param)
// 	if err != nil {
// 		config.Stdlog.Println("File upload 오류 : ", err)
// 	}
// 	bytes, _ := io.ReadAll(resp.Body)
// 	c.SetContentType("application/json")
// 	c.SetStatusCode(fasthttp.StatusOK)
// 	c.SetBody(bytes)
// }

// // 검수 승인 취소
// func Template_cancel_approval_(c *fasthttp.RequestCtx) {
// 	conf := config.Conf

// 	param := &TemplateRequest{}
// 	err := c.Bind(param)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	jsonstr, _ := json.Marshal(param)
// 	buff := bytes.NewBuffer(jsonstr)
// 	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/cancel_approval", buff)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	req.Header.Add("Content-Type", "application/json")
// 	resp, err := centerClient.Do(req)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	bytes, _ := io.ReadAll(resp.Body)
// 	c.SetContentType("application/json")
// 	c.SetStatusCode(fasthttp.StatusOK)
// 	c.SetBody(bytes)

// }

// // 기등록된 템플릿 (타입 : BA, EX) 을 채널추가버튼 및 채널추가안내문구가 포함된 템플릿으로 전환 /template/convertAddCh
// func Template_convertAddCh_(c *fasthttp.RequestCtx) {
// 	conf := config.Conf

// 	param := &Template_convertAddCh{}
// 	err := c.Bind(param)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	jsonstr, _ := json.Marshal(param)
// 	buff := bytes.NewBuffer(jsonstr)
// 	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/convertAddCh", buff)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	req.Header.Add("Content-Type", "application/json")
// 	resp, err := centerClient.Do(req)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	bytes, _ := io.ReadAll(resp.Body)
// 	c.SetContentType("application/json")
// 	c.SetStatusCode(fasthttp.StatusOK)
// 	c.SetBody(bytes)

// }

// // 채널에 발신 프로필 추가
// func Channel_sender_add_(c *fasthttp.RequestCtx) {
// 	conf := config.Conf

// 	param := &Channel_sender{}
// 	err := c.Bind(param)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	jsonstr, _ := json.Marshal(param)
// 	buff := bytes.NewBuffer(jsonstr)
// 	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/channel/sender/add", buff)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	req.Header.Add("Content-Type", "application/json")
// 	resp, err := centerClient.Do(req)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	bytes, _ := io.ReadAll(resp.Body)
// 	c.SetContentType("application/json")
// 	c.SetStatusCode(fasthttp.StatusOK)
// 	c.SetBody(bytes)

// }

// // 채널에 발신 프로필 삭제
// func Channel_sender_remove_(c *fasthttp.RequestCtx) {
// 	conf := config.Conf

// 	param := &Channel_sender{}
// 	err := c.Bind(param)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	jsonstr, _ := json.Marshal(param)
// 	buff := bytes.NewBuffer(jsonstr)
// 	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/channel/sender/remove", buff)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	req.Header.Add("Content-Type", "application/json")
// 	resp, err := centerClient.Do(req)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	bytes, _ := io.ReadAll(resp.Body)
// 	c.SetContentType("application/json")
// 	c.SetStatusCode(fasthttp.StatusOK)
// 	c.SetBody(bytes)

// }

// // 알림톡, 친구톡 발송 일별 통계
// func Stat_daily(c *fasthttp.RequestCtx) {
// 	conf := config.Conf

// 	//beginDate := c.Query("beginDate")
// 	//endDate := c.Query("endDate")
// 	//productType := c.Query("productType")
// 	//page := c.Query("page")

// 	params := map[string]string{
// 		"beginDate":   c.Query("beginDate"),
// 		"endDate":     c.Query("endDate"),
// 		"productType": c.Query("productType"),
// 		"page":        c.Query("page"),
// 	}

// 	if !MissingParams(c, params) {
// 		return
// 	}

// 	query := c.Request.URL.Query()
// 	for key, value := range params {
// 		if value != "" {
// 			query.Set(key, value)
// 		}
// 	}

// 	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/stat/daily?"+query.Encode(), nil)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}

// 	resp, err := centerClient.Do(req)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	bytes, _ := io.ReadAll(resp.Body)
// 	c.SetContentType("application/json")
// 	c.SetStatusCode(fasthttp.StatusOK)
// 	c.SetBody(bytes)
// }

// // 그룹 태그 생성
// func GroupTag_create(c *fasthttp.RequestCtx) {
// 	conf := config.Conf

// 	param := &Group_Tag_create{}
// 	err := c.Bind(param)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	jsonstr, _ := json.Marshal(param)
// 	buff := bytes.NewBuffer(jsonstr)
// 	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/groupTag/create", buff)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	req.Header.Add("Content-Type", "application/json")
// 	resp, err := centerClient.Do(req)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	bytes, _ := io.ReadAll(resp.Body)
// 	c.SetContentType("application/json")
// 	c.SetStatusCode(fasthttp.StatusOK)
// 	c.SetBody(bytes)

// }

// // 그룹 태그 조회
// func GroupTag_(c *fasthttp.RequestCtx) {
// 	conf := config.Conf

// 	//senderKey := c.Query("senderKey")
// 	//groupTagKey := c.Query("groupTagKey")

// 	params := map[string]string{
// 		"senderKey":   c.Query("senderKey"),
// 		"groupTagKey": c.Query("groupTagKey"),
// 	}

// 	if !MissingParams(c, params) {
// 		return
// 	}

// 	query := c.Request.URL.Query()
// 	for key, value := range params {
// 		if value != "" {
// 			query.Set(key, value)
// 		}
// 	}

// 	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/groupTag?"+query.Encode(), nil)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}

// 	resp, err := centerClient.Do(req)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	bytes, _ := io.ReadAll(resp.Body)
// 	c.SetContentType("application/json")
// 	c.SetStatusCode(fasthttp.StatusOK)
// 	c.SetBody(bytes)
// }

// // 그룹 태그 전체 조회
// func GroupTag_list(c *fasthttp.RequestCtx) {
// 	conf := config.Conf

// 	//senderKey := c.Query("senderKey")

// 	params := map[string]string{
// 		"senderKey": c.Query("senderKey"),
// 	}

// 	if !MissingParams(c, params) {
// 		return
// 	}

// 	query := c.Request.URL.Query()
// 	for key, value := range params {
// 		if value != "" {
// 			query.Set(key, value)
// 		}
// 	}

// 	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/groupTag/list?"+query.Encode(), nil)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}

// 	resp, err := centerClient.Do(req)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	bytes, _ := io.ReadAll(resp.Body)
// 	c.SetContentType("application/json")
// 	c.SetStatusCode(fasthttp.StatusOK)
// 	c.SetBody(bytes)

// }

// // 그룹 태그 수정
// func GroupTag_update(c *fasthttp.RequestCtx) {
// 	conf := config.Conf
// 	param := &Group_Tag_update{}
// 	err := c.Bind(param)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	jsonstr, _ := json.Marshal(param)
// 	buff := bytes.NewBuffer(jsonstr)
// 	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/groupTag/update", buff)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	req.Header.Add("Content-Type", "application/json")
// 	resp, err := centerClient.Do(req)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	bytes, _ := io.ReadAll(resp.Body)
// 	c.SetContentType("application/json")
// 	c.SetStatusCode(fasthttp.StatusOK)
// 	c.SetBody(bytes)

// }

// // 그룹 태그 삭제
// func GroupTag_delete(c *fasthttp.RequestCtx) {
// 	conf := config.Conf
// 	param := &Group_Tag_delete{}
// 	err := c.Bind(param)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	jsonstr, _ := json.Marshal(param)
// 	buff := bytes.NewBuffer(jsonstr)
// 	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/groupTag/delete", buff)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	req.Header.Add("Content-Type", "application/json")
// 	resp, err := centerClient.Do(req)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	bytes, _ := io.ReadAll(resp.Body)
// 	c.SetContentType("application/json")
// 	c.SetStatusCode(fasthttp.StatusOK)
// 	c.SetBody(bytes)

// }

// // 광고성 메시지(다이렉트) 템플릿 등록
// func Direct_template_create_(c *fasthttp.RequestCtx) {
// 	conf := config.Conf
// 	param := &Direct_template_create{}
// 	err := c.Bind(param)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	jsonstr, _ := json.Marshal(param)
// 	buff := bytes.NewBuffer(jsonstr)
// 	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v3/"+conf.PROFILE_KEY+"/direct/template/create", buff)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	req.Header.Add("Content-Type", "application/json")
// 	resp, err := centerClient.Do(req)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	bytes, _ := io.ReadAll(resp.Body)
// 	c.SetContentType("application/json")
// 	c.SetStatusCode(fasthttp.StatusOK)
// 	c.SetBody(bytes)
// }

// // 광고성메시지(다이렉트) 템플릿 조회
// func Direct_template_(c *fasthttp.RequestCtx) {
// 	conf := config.Conf
// 	code := c.Param("code")

// 	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/direct/template/"+code, nil)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}

// 	resp, err := centerClient.Do(req)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	bytes, _ := io.ReadAll(resp.Body)
// 	c.SetContentType("application/json")
// 	c.SetStatusCode(fasthttp.StatusOK)
// 	c.SetBody(bytes)
// }

// // 광고성메시지(다이렉트) 템플릿 수정
// func Direct_template_update_(c *fasthttp.RequestCtx) {
// 	conf := config.Conf
// 	code := c.Param("code")
// 	param := &Direct_template_create{}
// 	err := c.Bind(param)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	jsonstr, _ := json.Marshal(param)
// 	buff := bytes.NewBuffer(jsonstr)
// 	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v3/"+conf.PROFILE_KEY+"/direct/template/update/"+code, buff)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	req.Header.Add("Content-Type", "application/json")
// 	resp, err := centerClient.Do(req)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	bytes, _ := io.ReadAll(resp.Body)
// 	c.SetContentType("application/json")
// 	c.SetStatusCode(fasthttp.StatusOK)
// 	c.SetBody(bytes)
// }

// // 광고성메시지(다이렉트) 템플릿 삭제
// func Direct_template_delete_(c *fasthttp.RequestCtx) {
// 	conf := config.Conf
// 	code := c.Param("code")

// 	buff := bytes.NewBuffer([]byte(`{}`))
// 	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/direct/template/delete/"+code, buff)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	req.Header.Add("Content-Type", "application/json")
// 	resp, err := centerClient.Do(req)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	bytes, _ := io.ReadAll(resp.Body)
// 	c.SetContentType("application/json")
// 	c.SetStatusCode(fasthttp.StatusOK)
// 	c.SetBody(bytes)
// }

// // 발신채널 전환
// func Direct_convert_(c *fasthttp.RequestCtx) {
// 	conf := config.Conf
// 	param := &Direct_convert{}
// 	err := c.Bind(param)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	jsonstr, _ := json.Marshal(param)
// 	buff := bytes.NewBuffer(jsonstr)
// 	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/sender/direct/convert", buff)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	req.Header.Add("Content-Type", "application/json")
// 	resp, err := centerClient.Do(req)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	bytes, _ := io.ReadAll(resp.Body)
// 	c.SetContentType("application/json")
// 	c.SetStatusCode(fasthttp.StatusOK)
// 	c.SetBody(bytes)
// }

// // 발신채널 전환 상태 확인
// func Direct_convert_result(c *fasthttp.RequestCtx) {
// 	conf := config.Conf

// 	params := map[string]string{
// 		"senderKey": c.Query("senderKey"),
// 	}

// 	if !MissingParams(c, params) {
// 		return
// 	}

// 	query := c.Request.URL.Query()
// 	for key, value := range params {
// 		if value != "" {
// 			query.Set(key, value)
// 		}
// 	}

// 	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/sender/direct/convert/result?"+query.Encode(), nil)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}

// 	resp, err := centerClient.Do(req)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	bytes, _ := io.ReadAll(resp.Body)
// 	c.SetContentType("application/json")
// 	c.SetStatusCode(fasthttp.StatusOK)
// 	c.SetBody(bytes)
// }

// // 발신채널에 연결된 비즈월렛 변경
// func Direct_bizWallet_change_(c *fasthttp.RequestCtx) {
// 	conf := config.Conf
// 	param := &Direct_bizWallet_change{}
// 	err := c.Bind(param)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	jsonstr, _ := json.Marshal(param)
// 	buff := bytes.NewBuffer(jsonstr)
// 	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/sender/direct/bizWallet/change", buff)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	req.Header.Add("Content-Type", "application/json")
// 	resp, err := centerClient.Do(req)
// 	if err != nil {
// 		c.Error(err.Error(), fasthttp.StatusBadRequest)
// 		return
// 	}
// 	bytes, _ := io.ReadAll(resp.Body)
// 	c.SetContentType("application/json")
// 	c.SetStatusCode(fasthttp.StatusOK)
// 	c.SetBody(bytes)
// }

// func image_Seq_Mapping(c *fasthttp.RequestCtx, param map[string]io.Reader, max int, filename string) (map[string]io.Reader, error) {
// 	var retErr error = nil
// 	if max > 0 {
// 		for a := 1; a <= max; a++ {
// 			file, err := c.FormFile(filename + "_" + strconv.Itoa(a))
// 			newFileName := ""
// 			if err == nil {
// 				extension := filepath.Ext(file.Filename)
// 				newFileName = uuid.New().String() + extension
// 				err2 := c.SaveUploadedFile(file, config.BasePath+"upload/"+newFileName)
// 				if err2 != nil {
// 					newFileName = "_"
// 					retErr = err2
// 				}
// 			} else {
// 				newFileName = "_"
// 				retErr = err
// 			}

// 			if newFileName != "_" {
// 				param[filename+"_"+strconv.Itoa(a)] = mustOpen(config.BasePath + "upload/" + newFileName)
// 			}
// 		}
// 	} else {
// 		file, err := c.FormFile(filename)
// 		newFileName := ""
// 		if err == nil {
// 			extension := filepath.Ext(file.Filename)
// 			newFileName = uuid.New().String() + extension
// 			err2 := c.SaveUploadedFile(file, config.BasePath+"upload/"+newFileName)
// 			if err2 != nil {
// 				newFileName = "_"
// 				retErr = err2
// 			}
// 		} else {
// 			newFileName = "_"
// 			retErr = err
// 		}

// 		if newFileName != "_" {
// 			param[filename] = mustOpen(config.BasePath + "upload/" + newFileName)
// 		}
// 	}
// 	return param, retErr
// }

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

// func MissingParams(c *fasthttp.RequestCtx, params map[string]string) bool {
// 	var missingParams []string

// 	for param, value := range params {
// 		if value == "" {
// 			missingParams = append(missingParams, param)
// 		}
// 	}

// 	if len(missingParams) > 0 {
// 		message := "필수값이 부족합니다. ( "
// 		for i, param := range missingParams {
// 			if i != 0 {
// 				message += ", "
// 			}
// 			message += param
// 		}
// 		message += " )"
// 		c.JSON(999, gin.H{
// 			"message": message,
// 		})
// 		return false
// 	}

// 	return true
// }
